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

package feishu

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"html"
	"net/http"
	"regexp"
	"strings"

	imEntity "github.com/coze-dev/coze-studio/backend/domain/im/entity"
	imService "github.com/coze-dev/coze-studio/backend/domain/im/service"
	"github.com/coze-dev/coze-studio/backend/infra/im/factory"
)

var (
	atTagRegexp      = regexp.MustCompile(`(?is)<at\b[^>]*>.*?</at>`)
	multipleWSRegexp = regexp.MustCompile(`\s+`)
)

type Adapter struct {
	cfg    Config
	client *apiClient
}

type callbackEnvelope struct {
	Challenge string          `json:"challenge"`
	Token     string          `json:"token"`
	Encrypt   string          `json:"encrypt"`
	Header    *callbackHeader `json:"header"`
	Event     *callbackEvent  `json:"event"`
}

type callbackHeader struct {
	EventType string `json:"event_type"`
	EventID   string `json:"event_id"`
	AppID     string `json:"app_id"`
	TenantKey string `json:"tenant_key"`
}

type callbackEvent struct {
	Sender  *callbackSender  `json:"sender"`
	Message *callbackMessage `json:"message"`
}

type callbackSender struct {
	SenderID   *callbackSenderID `json:"sender_id"`
	SenderType string            `json:"sender_type"`
	SenderName string            `json:"sender_name"`
}

type callbackSenderID struct {
	UserID  string `json:"user_id"`
	OpenID  string `json:"open_id"`
	UnionID string `json:"union_id"`
}

type callbackMessage struct {
	MessageID   string             `json:"message_id"`
	RootID      string             `json:"root_id"`
	ParentID    string             `json:"parent_id"`
	ChatID      string             `json:"chat_id"`
	ThreadID    string             `json:"thread_id"`
	ChatType    string             `json:"chat_type"`
	MessageType string             `json:"message_type"`
	Content     string             `json:"content"`
	Mentions    []*callbackMention `json:"mentions"`
}

type callbackMention struct {
	Key  string             `json:"key"`
	Name string             `json:"name"`
	ID   *callbackMentionID `json:"id"`
	Type string             `json:"type"`
}

type callbackMentionID struct {
	OpenID  string `json:"open_id"`
	UserID  string `json:"user_id"`
	UnionID string `json:"union_id"`
}

type messageContent struct {
	Text string `json:"text"`
}

func init() {
	factory.Register(imEntity.PlatformFeishu, func() imService.PlatformAdapter {
		return NewAdapter()
	})
}

func NewAdapter() *Adapter {
	cfg := LoadConfig()
	return &Adapter{
		cfg:    cfg,
		client: newClient(cfg),
	}
}

func (a *Adapter) Platform() imEntity.Platform {
	return imEntity.PlatformFeishu
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
		return nil, nil, imEntity.NewStatusError(http.StatusServiceUnavailable, "feishu adapter is not configured")
	}

	payload, err := a.parseEnvelope(req)
	if err != nil {
		return nil, nil, err
	}
	if payload.Challenge != "" {
		return nil, imEntity.JSONResponse(http.StatusOK, map[string]string{"challenge": payload.Challenge}), nil
	}

	if a.cfg.VerificationToken != "" && payload.Token != "" && payload.Token != a.cfg.VerificationToken {
		return nil, nil, imEntity.NewStatusError(http.StatusUnauthorized, "invalid feishu verification token")
	}
	if payload.Header == nil || payload.Event == nil || payload.Event.Message == nil || payload.Event.Sender == nil {
		return nil, nil, nil
	}
	if payload.Header.EventType != "im.message.receive_v1" {
		return nil, nil, nil
	}
	if a.cfg.AppID != "" && payload.Header.AppID != "" && a.cfg.AppID != payload.Header.AppID {
		return nil, nil, imEntity.NewStatusError(http.StatusUnauthorized, "unexpected feishu app id")
	}
	if payload.Event.Message.MessageType != "text" {
		return nil, nil, nil
	}
	if payload.Event.Message.ChatType != "" && payload.Event.Message.ChatType != "group" {
		return nil, nil, nil
	}
	if payload.Event.Sender.SenderType != "" && payload.Event.Sender.SenderType != "user" {
		return nil, nil, nil
	}

	agentID, ok := a.cfg.AgentID(payload.Header.AppID)
	if !ok {
		return nil, nil, imEntity.NewStatusError(http.StatusNotFound, "feishu app id is not bound to a published agent")
	}

	text, isAtBot, err := extractText(payload.Event.Message.Content, payload.Event.Message.Mentions)
	if err != nil {
		return nil, nil, imEntity.NewStatusError(http.StatusBadRequest, "invalid feishu message content")
	}
	if !isAtBot || strings.TrimSpace(text) == "" {
		return nil, nil, nil
	}

	threadID := firstNonEmpty(payload.Event.Message.ThreadID, payload.Event.Message.RootID, payload.Event.Message.ChatID)
	if payload.Event.Message.ChatID == "" || payload.Event.Message.MessageID == "" {
		return nil, nil, imEntity.NewStatusError(http.StatusBadRequest, "missing feishu message identifiers")
	}

	return &imEntity.IMessage{
		Platform:      a.Platform(),
		AgentID:       agentID,
		ExternalBotID: payload.Header.AppID,
		EventID:       payload.Header.EventID,
		MessageID:     payload.Event.Message.MessageID,
		SessionID:     buildSessionID(payload.Header.AppID, payload.Event.Message.ChatID, threadID),
		UserID:        firstNonEmpty(payload.Event.Sender.SenderID.UserID, payload.Event.Sender.SenderID.OpenID, payload.Event.Sender.SenderID.UnionID),
		UserName:      payload.Event.Sender.SenderName,
		ChatID:        payload.Event.Message.ChatID,
		ThreadID:      threadID,
		Text:          text,
		IsAtBot:       true,
		Metadata: map[string]string{
			"tenant_key": payload.Header.TenantKey,
			"app_id":     payload.Header.AppID,
			"chat_type":  payload.Event.Message.ChatType,
			"root_id":    payload.Event.Message.RootID,
			"parent_id":  payload.Event.Message.ParentID,
		},
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

	return a.client.replyText(ctx, msg.MessageID, text)
}

func (a *Adapter) SuccessResponse() *imEntity.RawResponse {
	return imEntity.JSONResponse(http.StatusOK, map[string]any{
		"code": 0,
		"msg":  "success",
	})
}

func (a *Adapter) parseEnvelope(req *imEntity.CallbackRequest) (*callbackEnvelope, error) {
	body := append([]byte(nil), req.Body...)
	if a.cfg.EncryptKey != "" {
		signature := headerValue(req.Headers, "X-Lark-Signature")
		timestamp := headerValue(req.Headers, "X-Lark-Request-Timestamp")
		nonce := headerValue(req.Headers, "X-Lark-Request-Nonce")
		if signature != "" || timestamp != "" || nonce != "" {
			if err := verifySignature(signature, timestamp, nonce, a.cfg.EncryptKey, body); err != nil {
				return nil, imEntity.NewStatusError(http.StatusUnauthorized, err.Error())
			}
		}
	}

	var envelope callbackEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, imEntity.NewStatusError(http.StatusBadRequest, "invalid feishu request body")
	}
	if envelope.Encrypt != "" {
		plainBody, err := decryptBody(a.cfg.EncryptKey, envelope.Encrypt)
		if err != nil {
			return nil, imEntity.NewStatusError(http.StatusBadRequest, "invalid feishu encrypted payload")
		}
		if err = json.Unmarshal(plainBody, &envelope); err != nil {
			return nil, imEntity.NewStatusError(http.StatusBadRequest, "invalid feishu event payload")
		}
	}

	return &envelope, nil
}

func extractText(rawContent string, mentions []*callbackMention) (string, bool, error) {
	var content messageContent
	if err := json.Unmarshal([]byte(rawContent), &content); err != nil {
		return "", false, err
	}

	isAtBot := len(mentions) > 0 || atTagRegexp.MatchString(content.Text)
	plainText := atTagRegexp.ReplaceAllString(content.Text, " ")
	plainText = html.UnescapeString(plainText)
	plainText = multipleWSRegexp.ReplaceAllString(strings.TrimSpace(plainText), " ")

	return plainText, isAtBot, nil
}

func buildSessionID(appID string, chatID string, threadID string) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{
		"feishu",
		appID,
		chatID,
		threadID,
	}, ":")))

	return hex.EncodeToString(sum[:])[:32]
}

func headerValue(headers map[string]string, key string) string {
	for headerKey, value := range headers {
		if strings.EqualFold(headerKey, key) {
			return strings.TrimSpace(value)
		}
	}

	return ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}

	return ""
}
