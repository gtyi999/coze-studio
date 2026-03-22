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

package wecom

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"net/http"
	"strings"

	imEntity "github.com/coze-dev/coze-studio/backend/domain/im/entity"
	imService "github.com/coze-dev/coze-studio/backend/domain/im/service"
	"github.com/coze-dev/coze-studio/backend/infra/im/factory"
)

type Adapter struct {
	cfg    Config
	client *apiClient
}

type encryptedEnvelope struct {
	XMLName xml.Name `xml:"xml"`
	ToUser  string   `xml:"ToUserName"`
	AgentID string   `xml:"AgentID"`
	Encrypt string   `xml:"Encrypt"`
}

type callbackMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	MsgID        string   `xml:"MsgId"`
	AgentID      string   `xml:"AgentID"`
	ChatID       string   `xml:"ChatId"`
	Event        string   `xml:"Event"`
	EventKey     string   `xml:"EventKey"`
}

func init() {
	factory.Register(imEntity.PlatformWeCom, func() imService.PlatformAdapter {
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
	return imEntity.PlatformWeCom
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
		return nil, nil, imEntity.NewStatusError(http.StatusServiceUnavailable, "wecom adapter is not configured")
	}

	timestamp := req.Query["timestamp"]
	nonce := req.Query["nonce"]
	signatureParam := req.Query["msg_signature"]

	if strings.EqualFold(req.Method, http.MethodGet) {
		echoStr := req.Query["echostr"]
		if echoStr == "" {
			return nil, nil, imEntity.NewStatusError(http.StatusBadRequest, "missing wecom echostr")
		}
		if !verifySignature(a.cfg.Token, timestamp, nonce, echoStr, signatureParam) {
			return nil, nil, imEntity.NewStatusError(http.StatusUnauthorized, "invalid wecom signature")
		}

		plainText, corpID, err := decryptMessage(a.cfg.EncodingAESKey, echoStr)
		if err != nil {
			return nil, nil, imEntity.NewStatusError(http.StatusBadRequest, "invalid wecom echostr")
		}
		if corpID != "" && corpID != a.cfg.CorpID {
			return nil, nil, imEntity.NewStatusError(http.StatusUnauthorized, "unexpected wecom corp id")
		}

		return nil, imEntity.TextResponse(http.StatusOK, plainText), nil
	}

	var envelope encryptedEnvelope
	if err := xml.Unmarshal(req.Body, &envelope); err != nil {
		return nil, nil, imEntity.NewStatusError(http.StatusBadRequest, "invalid wecom request body")
	}
	if strings.TrimSpace(envelope.Encrypt) == "" {
		return nil, nil, imEntity.NewStatusError(http.StatusBadRequest, "missing wecom encrypt field")
	}
	if !verifySignature(a.cfg.Token, timestamp, nonce, envelope.Encrypt, signatureParam) {
		return nil, nil, imEntity.NewStatusError(http.StatusUnauthorized, "invalid wecom signature")
	}

	plainText, corpID, err := decryptMessage(a.cfg.EncodingAESKey, envelope.Encrypt)
	if err != nil {
		return nil, nil, imEntity.NewStatusError(http.StatusBadRequest, "invalid wecom encrypted payload")
	}
	if corpID != "" && corpID != a.cfg.CorpID {
		return nil, nil, imEntity.NewStatusError(http.StatusUnauthorized, "unexpected wecom corp id")
	}

	var msg callbackMessage
	if err = xml.Unmarshal([]byte(plainText), &msg); err != nil {
		return nil, nil, imEntity.NewStatusError(http.StatusBadRequest, "invalid wecom event payload")
	}
	if msg.MsgType != "text" {
		return nil, nil, nil
	}

	externalBotID := msg.AgentID
	agentID, ok := a.cfg.AgentID(externalBotID)
	if !ok {
		return nil, nil, imEntity.NewStatusError(http.StatusNotFound, "wecom agent is not bound to a published agent")
	}

	text := normalizeText(msg.Content)
	isAtBot := strings.TrimSpace(msg.ChatID) == "" || strings.Contains(msg.Content, "@")
	if !isAtBot || text == "" {
		return nil, nil, nil
	}
	if strings.TrimSpace(msg.MsgID) == "" {
		return nil, nil, imEntity.NewStatusError(http.StatusBadRequest, "missing wecom message id")
	}

	return &imEntity.IMessage{
		Platform:      a.Platform(),
		AgentID:       agentID,
		ExternalBotID: externalBotID,
		EventID:       firstNonEmpty(msg.MsgID, msg.EventKey),
		MessageID:     msg.MsgID,
		SessionID:     buildSessionID(a.cfg.CorpID, externalBotID, firstNonEmpty(msg.ChatID, msg.FromUserName)),
		UserID:        msg.FromUserName,
		UserName:      msg.FromUserName,
		ChatID:        msg.ChatID,
		ThreadID:      msg.ChatID,
		Text:          text,
		IsAtBot:       true,
		Metadata: map[string]string{
			"agent_id": msg.AgentID,
			"corp_id":  a.cfg.CorpID,
			"chat_id":  msg.ChatID,
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

	return a.client.replyText(ctx, msg.Metadata["agent_id"], msg.UserID, msg.ChatID, text)
}

func (a *Adapter) SuccessResponse() *imEntity.RawResponse {
	return imEntity.TextResponse(http.StatusOK, "success")
}

func normalizeText(text string) string {
	text = strings.ReplaceAll(text, "\u00a0", " ")
	return strings.TrimSpace(strings.Join(strings.Fields(text), " "))
}

func buildSessionID(corpID string, agentID string, scope string) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{
		"wecom",
		corpID,
		agentID,
		scope,
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
