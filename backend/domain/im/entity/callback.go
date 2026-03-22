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

package entity

import "encoding/json"

type CallbackRequest struct {
	Platform Platform
	Method   string
	Headers  map[string]string
	Query    map[string]string
	Body     []byte
}

type RawResponse struct {
	StatusCode  int
	ContentType string
	Body        []byte
}

type IMessage struct {
	Platform      Platform
	AgentID       int64
	ExternalBotID string
	EventID       string
	MessageID     string
	SessionID     string
	UserID        string
	UserName      string
	ChatID        string
	ThreadID      string
	Text          string
	IsAtBot       bool
	Metadata      map[string]string
}

type AgentResponseStatus string

const (
	AgentResponseStatusSuccess AgentResponseStatus = "success"
	AgentResponseStatusPending AgentResponseStatus = "async_pending"
	AgentResponseStatusFailed  AgentResponseStatus = "failed"
)

type TaskTicket struct {
	Kind           TaskKind          `json:"kind"`
	TargetType     AgentTargetType   `json:"target_type"`
	TargetID       int64             `json:"target_id"`
	SessionID      string            `json:"session_id,omitempty"`
	ConversationID int64             `json:"conversation_id"`
	RunID          int64             `json:"run_id"`
	ExecuteID      string            `json:"execute_id,omitempty"`
	Status         string            `json:"status"`
	EventID        string            `json:"event_id,omitempty"`
	MessageID      string            `json:"message_id,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

type AgentResponse struct {
	Status   AgentResponseStatus `json:"status"`
	Text     string              `json:"text,omitempty"`
	Task     *TaskTicket         `json:"task,omitempty"`
	Metadata map[string]string   `json:"metadata,omitempty"`
}

type PlatformInfo struct {
	Platform     Platform
	Name         string
	ConnectorID  int64
	CallbackPath string
	CallbackURL  string
	Enabled      bool
}

type StatusError struct {
	StatusCode int
	Message    string
}

func (e *StatusError) Error() string {
	return e.Message
}

func NewStatusError(statusCode int, message string) error {
	return &StatusError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func JSONResponse(statusCode int, payload any) *RawResponse {
	body, err := json.Marshal(payload)
	if err != nil {
		body = []byte(`{"code":500,"msg":"marshal response failed"}`)
		statusCode = 500
	}

	return &RawResponse{
		StatusCode:  statusCode,
		ContentType: "application/json; charset=utf-8",
		Body:        body,
	}
}

func TextResponse(statusCode int, body string) *RawResponse {
	return &RawResponse{
		StatusCode:  statusCode,
		ContentType: "text/plain; charset=utf-8",
		Body:        []byte(body),
	}
}

func ErrorResponse(statusCode int, message string) *RawResponse {
	return JSONResponse(statusCode, map[string]any{
		"code": statusCode,
		"msg":  message,
	})
}

func AsyncPendingResponse(ticket *TaskTicket) *RawResponse {
	return JSONResponse(200, map[string]any{
		"code": 0,
		"msg":  "async_pending",
		"data": ticket,
	})
}
