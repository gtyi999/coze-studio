/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dingtalk

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	imEntity "github.com/coze-dev/coze-studio/backend/domain/im/entity"
	imService "github.com/coze-dev/coze-studio/backend/domain/im/service"
	"github.com/coze-dev/coze-studio/backend/infra/im/factory"
)

var (
	atRegexp         = regexp.MustCompile(`@\S+`)
	multiSpaceRegexp = regexp.MustCompile(`\s+`)
)

type Adapter struct {
	cfg    Config
	client *apiClient
}

type callbackEnvelope struct {
	Challenge string `json:"challenge"`
	Encrypt   string `json:"encrypt"`
}

type callbackMessage struct {
	EventID                   string            `json:"eventId"`
	MsgID                     string            `json:"msgId"`
	MsgType                   string            `json:"msgtype"`
	Text                      *messageText      `json:"text"`
	ConversationID            string            `json:"conversationId"`
	ConversationType          string            `json:"conversationType"`
	ConversationTitle         string            `json:"conversationTitle"`
	SessionWebhook            string            `json:"sessionWebhook"`
	SessionWebhookExpiredTime int64             `json:"sessionWebhookExpiredTime"`
	SenderStaffID             string            `json:"senderStaffId"`
	SenderNick                string            `json:"senderNick"`
	SenderCorpID              string            `json:"senderCorpId"`
	ChatbotUserID             string            `json:"chatbotUserId"`
	ChatbotCorpID             string            `json:"chatbotCorpId"`
	RobotCode                 string            `json:"robotCode"`
	AtUsers                   []*atUser         `json:"atUsers"`
	Extensions                map[string]string `json:"extensions"`
}

type messageText struct {
	Content string `json:"content"`
}

type atUser struct {
	DingTalkID string `json:"dingtalkId"`
	StaffID    string `json:"staffId"`
}

func init() {
	factory.Register(imEntity.PlatformDingTalk, func() imService.PlatformAdapter {
		return NewAdapter()
	})
}

func NewAdapter() *Adapter {
	cfg := LoadConfig()
	return &Adapter{
		cfg:    cfg,
		client: newClient(),
	}
}

func (a *Adapter) Platform() imEntity.Platform {
	return imEntity.PlatformDingTalk
}

func (a *Adapter) ConnectorID() int64 {
	return a.Platform().ConnectorID()
}

func (a *Adapter) PlatformInfo(host string) *imEntity.PlatformInfo {
	return &imEntity.PlatformInfo{
		Platform:     a.Platform(),
		Name:         a.Platform().DisplayName(),
		ConnectorID:  a.ConnectorID(),
		CallbackPath: a.Platform().CallbackPath(),
		CallbackURL:  strings.TrimRight(host, "/") + a.Platform().CallbackPath(),
		Enabled:      a.cfg.Enabled(),
	}
}

func (a *Adapter) ParseCallback(_ context.Context, req *imEntity.CallbackRequest) (*imEntity.IMessage, *imEntity.RawResponse, error) {
	if !a.cfg.Enabled() {
		return nil, nil, imEntity.NewStatusError(http.StatusServiceUnavailable, "dingtalk adapter is not configured")
	}

	body, challenge, err := a.parseBody(req)
	if err != nil {
		return nil, nil, err
	}
	if strings.TrimSpace(challenge) != "" {
		return nil, imEntity.JSONResponse(http.StatusOK, map[string]string{"challenge": challenge}), nil
	}

	var msg callbackMessage
	if err = json.Unmarshal(body, &msg); err != nil {
		return nil, nil, imEntity.NewStatusError(http.StatusBadRequest, "invalid dingtalk event payload")
	}
	if msg.MsgType != "text" || msg.Text == nil {
		return nil, nil, nil
	}
	if msg.ConversationType != "" && msg.ConversationType != "2" {
		return nil, nil, nil
	}

	externalBotID := firstNonEmpty(msg.RobotCode, a.cfg.AppKey, msg.ChatbotUserID)
	agentID, ok := a.cfg.AgentID(externalBotID)
	if !ok {
		return nil, nil, imEntity.NewStatusError(http.StatusNotFound, "dingtalk robot is not bound to a published agent")
	}

	text, isAtBot := extractText(msg.Text.Content, msg.AtUsers)
	if !isAtBot || strings.TrimSpace(text) == "" {
		return nil, nil, nil
	}
	if strings.TrimSpace(msg.ConversationID) == "" || strings.TrimSpace(msg.MsgID) == "" {
		return nil, nil, imEntity.NewStatusError(http.StatusBadRequest, "missing dingtalk message identifiers")
	}

	metadata := map[string]string{
		"session_webhook":              msg.SessionWebhook,
		"conversation_title":           msg.ConversationTitle,
		"chatbot_user_id":              msg.ChatbotUserID,
		"chatbot_corp_id":              msg.ChatbotCorpID,
		"sender_corp_id":               msg.SenderCorpID,
		"session_webhook_expired_time": "",
	}
	if msg.SessionWebhookExpiredTime > 0 {
		metadata["session_webhook_expired_time"] = strconv.FormatInt(msg.SessionWebhookExpiredTime, 10)
	}
	for key, value := range msg.Extensions {
		metadata["extension_"+key] = value
	}

	return &imEntity.IMessage{
		Platform:      a.Platform(),
		AgentID:       agentID,
		ExternalBotID: externalBotID,
		EventID:       firstNonEmpty(msg.EventID, msg.MsgID),
		MessageID:     msg.MsgID,
		SessionID:     buildSessionID(externalBotID, msg.ConversationID),
		UserID:        msg.SenderStaffID,
		UserName:      msg.SenderNick,
		ChatID:        msg.ConversationID,
		ThreadID:      msg.ConversationID,
		Text:          text,
		IsAtBot:       true,
		Metadata:      metadata,
	}, nil, nil
}

func (a *Adapter) SendReply(ctx context.Context, msg *imEntity.IMessage, reply *imEntity.AgentResponse) error {
	if msg == nil || reply == nil {
		return nil
	}

	text := strings.TrimSpace(reply.Text)
	if text == "" {
		return nil
	}

	return a.client.replyText(ctx, msg.Metadata["session_webhook"], text)
}

func (a *Adapter) SuccessResponse() *imEntity.RawResponse {
	return imEntity.JSONResponse(http.StatusOK, map[string]any{
		"errcode": 0,
		"errmsg":  "ok",
	})
}

func (a *Adapter) parseBody(req *imEntity.CallbackRequest) ([]byte, string, error) {
	body := append([]byte(nil), req.Body...)

	var envelope callbackEnvelope
	if len(body) > 0 {
		_ = json.Unmarshal(body, &envelope)
	}
	if strings.TrimSpace(envelope.Challenge) != "" {
		return body, envelope.Challenge, nil
	}
	if strings.TrimSpace(envelope.Encrypt) == "" {
		return body, "", nil
	}

	signature := firstNonEmpty(req.Query["signature"], req.Query["msg_signature"])
	timestamp := req.Query["timestamp"]
	nonce := req.Query["nonce"]
	if a.cfg.Token != "" && signature != "" && !verifySignature(a.cfg.Token, timestamp, nonce, envelope.Encrypt, signature) {
		return nil, "", imEntity.NewStatusError(http.StatusUnauthorized, "invalid dingtalk signature")
	}

	plainBody, err := decryptBody(a.cfg.AESKey, envelope.Encrypt)
	if err != nil {
		return nil, "", imEntity.NewStatusError(http.StatusBadRequest, "invalid dingtalk encrypted payload")
	}

	var plainEnvelope callbackEnvelope
	if err = json.Unmarshal(plainBody, &plainEnvelope); err == nil && strings.TrimSpace(plainEnvelope.Challenge) != "" {
		return plainBody, plainEnvelope.Challenge, nil
	}

	return plainBody, "", nil
}

func extractText(raw string, atUsers []*atUser) (string, bool) {
	isAtBot := len(atUsers) > 0 || strings.Contains(raw, "@")
	text := strings.TrimSpace(atRegexp.ReplaceAllString(raw, " "))
	text = strings.TrimSpace(multiSpaceRegexp.ReplaceAllString(text, " "))
	return text, isAtBot
}

func buildSessionID(externalBotID string, conversationID string) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{
		"dingtalk",
		externalBotID,
		conversationID,
	}, ":")))

	return hex.EncodeToString(sum[:])[:32]
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}

	return ""
}
