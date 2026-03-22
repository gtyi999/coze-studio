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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type apiClient struct {
	httpClient *http.Client
}

type commonResponse struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func newClient() *apiClient {
	return &apiClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *apiClient) replyText(ctx context.Context, sessionWebhook string, text string) error {
	if strings.TrimSpace(sessionWebhook) == "" {
		return fmt.Errorf("missing dingtalk session webhook")
	}

	payload, err := json.Marshal(map[string]any{
		"msgtype": "text",
		"text": map[string]string{
			"content": text,
		},
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sessionWebhook, bytes.NewReader(payload))
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
		return fmt.Errorf("dingtalk reply failed, status=%d body=%s", resp.StatusCode, string(respBody))
	}
	if len(respBody) == 0 {
		return nil
	}

	var commonResp commonResponse
	if err = json.Unmarshal(respBody, &commonResp); err != nil {
		return nil
	}
	if commonResp.ErrCode != 0 {
		return fmt.Errorf("dingtalk reply failed: code=%d msg=%s", commonResp.ErrCode, commonResp.ErrMsg)
	}

	return nil
}
