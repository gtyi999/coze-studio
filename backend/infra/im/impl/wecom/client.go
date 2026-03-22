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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	getTokenURL    = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	messageSendURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
	appChatSendURL = "https://qyapi.weixin.qq.com/cgi-bin/appchat/send"
)

type apiClient struct {
	cfg        Config
	httpClient *http.Client

	mu        sync.Mutex
	tokenByID map[string]string
	expByID   map[string]time.Time
}

type accessTokenResponse struct {
	ErrCode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type commonResponse struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func newClient(cfg Config) *apiClient {
	return &apiClient{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		tokenByID: map[string]string{},
		expByID:   map[string]time.Time{},
	}
}

func (c *apiClient) replyText(ctx context.Context, agentID string, userID string, chatID string, text string) error {
	accessToken, err := c.getAccessToken(ctx, agentID)
	if err != nil {
		return err
	}

	if strings.TrimSpace(chatID) != "" {
		return c.sendJSON(ctx, withToken(appChatSendURL, accessToken), map[string]any{
			"chatid":  chatID,
			"msgtype": "text",
			"text": map[string]string{
				"content": text,
			},
			"safe": 0,
		})
	}

	if strings.TrimSpace(userID) == "" {
		return fmt.Errorf("missing wecom reply user")
	}

	agentInt, err := strconv.ParseInt(agentID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid wecom agent id: %w", err)
	}

	return c.sendJSON(ctx, withToken(messageSendURL, accessToken), map[string]any{
		"touser":  userID,
		"msgtype": "text",
		"agentid": agentInt,
		"text": map[string]string{
			"content": text,
		},
		"safe": 0,
	})
}

func (c *apiClient) getAccessToken(ctx context.Context, agentID string) (string, error) {
	c.mu.Lock()
	if token := strings.TrimSpace(c.tokenByID[agentID]); token != "" && time.Until(c.expByID[agentID]) > time.Minute {
		c.mu.Unlock()
		return token, nil
	}
	c.mu.Unlock()

	secret, ok := c.cfg.AgentSecret(agentID)
	if !ok || strings.TrimSpace(secret) == "" {
		return "", fmt.Errorf("missing wecom agent secret for agent_id=%s", agentID)
	}

	query := url.Values{}
	query.Set("corpid", c.cfg.CorpID)
	query.Set("corpsecret", secret)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getTokenURL+"?"+query.Encode(), nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("get wecom access token failed, status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var tokenResp accessTokenResponse
	if err = json.Unmarshal(respBody, &tokenResp); err != nil {
		return "", err
	}
	if tokenResp.ErrCode != 0 || strings.TrimSpace(tokenResp.AccessToken) == "" {
		return "", fmt.Errorf("get wecom access token failed: code=%d msg=%s", tokenResp.ErrCode, tokenResp.ErrMsg)
	}

	expire := tokenResp.ExpiresIn
	if expire <= 0 {
		expire = 7200
	}

	c.mu.Lock()
	c.tokenByID[agentID] = tokenResp.AccessToken
	c.expByID[agentID] = time.Now().Add(time.Duration(expire) * time.Second)
	c.mu.Unlock()

	return tokenResp.AccessToken, nil
}

func (c *apiClient) sendJSON(ctx context.Context, rawURL string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("wecom send reply failed, status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var commonResp commonResponse
	if err = json.Unmarshal(respBody, &commonResp); err != nil {
		return err
	}
	if commonResp.ErrCode != 0 {
		return fmt.Errorf("wecom send reply failed: code=%d msg=%s", commonResp.ErrCode, commonResp.ErrMsg)
	}

	return nil
}

func withToken(rawURL, accessToken string) string {
	query := url.Values{}
	query.Set("access_token", accessToken)
	return rawURL + "?" + query.Encode()
}
