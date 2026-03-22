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

type TaskStatus string

const (
	TaskStatusPending  TaskStatus = "pending"
	TaskStatusRunning  TaskStatus = "running"
	TaskStatusSuccess  TaskStatus = "success"
	TaskStatusFailed   TaskStatus = "failed"
	TaskStatusRetrying TaskStatus = "retrying"
)

type TaskRecord struct {
	ID              string         `json:"id"`
	Platform        Platform       `json:"platform"`
	Status          TaskStatus     `json:"status"`
	IdempotencyKey  string         `json:"idempotency_key,omitempty"`
	RetryCount      int32          `json:"retry_count"`
	MaxRetryCount   int32          `json:"max_retry_count"`
	LastError       string         `json:"last_error,omitempty"`
	NextRetryAtMS   int64          `json:"next_retry_at_ms,omitempty"`
	DeadlineAtMS    int64          `json:"deadline_at_ms,omitempty"`
	CreatedAtMS     int64          `json:"created_at_ms"`
	UpdatedAtMS     int64          `json:"updated_at_ms"`
	StartedAtMS     int64          `json:"started_at_ms,omitempty"`
	FinishedAtMS    int64          `json:"finished_at_ms,omitempty"`
	DeliveredAtMS   int64          `json:"delivered_at_ms,omitempty"`
	ResultDelivered bool           `json:"result_delivered"`
	Message         *IMessage      `json:"message,omitempty"`
	Request         *AgentRequest  `json:"request,omitempty"`
	Ticket          *TaskTicket    `json:"ticket,omitempty"`
	Result          *AgentResponse `json:"result,omitempty"`
}

func (t *TaskRecord) IsFinal() bool {
	if t == nil {
		return false
	}

	return t.Status == TaskStatusSuccess || t.Status == TaskStatusFailed
}

func (t *TaskRecord) CanRetry() bool {
	if t == nil {
		return false
	}

	return t.RetryCount < t.MaxRetryCount
}
