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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	tenantAccessTokenURL = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
	replyMessageURLTmpl  = "https://open.feishu.cn/open-apis/im/v1/messages/%s/reply"
)

type apiClient struct {
	cfg        Config
	httpClient *http.Client

	token     string
	tokenExp  time.Time
	tokenLock chan struct{}
}

type tenantAccessTokenResponse struct {
	Code              int64  `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int64  `json:"expire"`
}

type commonResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func newClient(cfg Config) *apiClient {
	lock := make(chan struct{}, 1)
	lock <- struct{}{}

	return &apiClient{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		tokenLock: lock,
	}
}

func (c *apiClient) replyText(ctx context.Context, messageID string, text string) error {
	accessToken, err := c.getTenantAccessToken(ctx)
	if err != nil {
		return err
	}

	var resp commonResponse
	if err = c.doJSON(ctx, http.MethodPost, fmt.Sprintf(replyMessageURLTmpl, messageID), map[string]any{
		"msg_type": "text",
		"content": map[string]string{
			"text": text,
		},
	}, map[string]string{
		"Authorization": "Bearer " + accessToken,
	}, &resp); err != nil {
		return err
	}
	if resp.Code != 0 {
		return fmt.Errorf("feishu reply failed: code=%d msg=%s", resp.Code, resp.Msg)
	}

	return nil
}

func (c *apiClient) getTenantAccessToken(ctx context.Context) (string, error) {
	<-c.tokenLock
	defer func() {
		c.tokenLock <- struct{}{}
	}()

	if c.token != "" && time.Until(c.tokenExp) > time.Minute {
		return c.token, nil
	}

	var resp tenantAccessTokenResponse
	if err := c.doJSON(ctx, http.MethodPost, tenantAccessTokenURL, map[string]string{
		"app_id":     c.cfg.AppID,
		"app_secret": c.cfg.AppSecret,
	}, nil, &resp); err != nil {
		return "", err
	}
	if resp.Code != 0 || resp.TenantAccessToken == "" {
		return "", fmt.Errorf("get feishu tenant access token failed: code=%d msg=%s", resp.Code, resp.Msg)
	}

	expire := resp.Expire
	if expire <= 0 {
		expire = 7200
	}

	c.token = resp.TenantAccessToken
	c.tokenExp = time.Now().Add(time.Duration(expire) * time.Second)

	return c.token, nil
}

func (c *apiClient) doJSON(ctx context.Context, method string, rawURL string, requestBody any, headers map[string]string, out any) error {
	var bodyReader io.Reader
	if requestBody != nil {
		payload, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, bodyReader)
	if err != nil {
		return err
	}
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

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
		return fmt.Errorf("request failed, status=%d body=%s", resp.StatusCode, string(respBody))
	}
	if out == nil {
		return nil
	}

	return json.Unmarshal(respBody, out)
}
