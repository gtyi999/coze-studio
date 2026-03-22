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

package im

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	imEntity "github.com/coze-dev/coze-studio/backend/domain/im/entity"
	imService "github.com/coze-dev/coze-studio/backend/domain/im/service"
)

const (
	defaultTaskMaxRetryCount = int32(3)
	defaultTaskDeadline      = 15 * time.Minute
	taskExecutionLockTTL     = 2 * time.Minute
	taskRetryBaseDelay       = 15 * time.Second
)

func (a *IMApplicationService) CreateTask(ctx context.Context, msg *imEntity.IMessage, req *imEntity.AgentRequest) (*imEntity.TaskRecord, bool, error) {
	taskID := buildTaskID(msg)
	existing, found, err := a.TaskRepo.Get(ctx, taskID)
	if err != nil {
		return nil, false, err
	}
	if found {
		return existing, true, nil
	}

	nowMS := nowMillis()
	task := &imEntity.TaskRecord{
		ID:             taskID,
		Platform:       msg.Platform,
		Status:         imEntity.TaskStatusPending,
		IdempotencyKey: buildIdempotencyKey(msg),
		MaxRetryCount:  defaultTaskMaxRetryCount,
		DeadlineAtMS:   nowMS + int64(defaultTaskDeadline/time.Millisecond),
		CreatedAtMS:    nowMS,
		UpdatedAtMS:    nowMS,
		Message:        cloneMessage(msg),
		Request:        cloneAgentRequest(req),
	}

	return task, false, a.TaskRepo.Save(ctx, task)
}

func (a *IMApplicationService) ExecuteTask(ctx context.Context, taskID string) error {
	task, found, err := a.TaskRepo.Get(ctx, taskID)
	if err != nil {
		return err
	}
	if !found || task == nil {
		return nil
	}
	if task.IsFinal() && task.ResultDelivered {
		return nil
	}

	acquired, err := a.TaskRepo.TryAcquireExecution(ctx, task.ID, taskExecutionLockTTL)
	if err != nil {
		return err
	}
	if !acquired {
		return nil
	}
	defer func() {
		_ = a.TaskRepo.ReleaseExecution(context.WithoutCancel(ctx), task.ID)
	}()

	if expired := task.DeadlineAtMS > 0 && nowMillis() > task.DeadlineAtMS; expired {
		return a.markTaskFailed(ctx, task, "task execution timed out")
	}

	nowMS := nowMillis()
	task.Status = imEntity.TaskStatusRunning
	task.UpdatedAtMS = nowMS
	if task.StartedAtMS == 0 {
		task.StartedAtMS = nowMS
	}
	if err = a.TaskRepo.Save(ctx, task); err != nil {
		return err
	}

	adapter, err := a.IMDomainSVC.GetAdapter(task.Platform)
	if err != nil {
		return err
	}

	if task.Result != nil && !task.ResultDelivered {
		return a.deliverTaskResult(ctx, adapter, task, task.Result)
	}
	if task.Ticket != nil {
		return a.awaitGatewayTask(ctx, adapter, task)
	}
	if task.Request == nil {
		return a.markTaskFailed(ctx, task, "task request is missing")
	}

	resp, err := a.Gateway.Execute(ctx, task.Request)
	if err != nil {
		return a.scheduleRetryOrFail(ctx, task, fmt.Sprintf("execute agent request failed: %v", err), true)
	}

	return a.consumeGatewayResponse(ctx, adapter, task, resp)
}

func (a *IMApplicationService) awaitGatewayTask(ctx context.Context, adapter imService.PlatformAdapter, task *imEntity.TaskRecord) error {
	if task.Request == nil || task.Ticket == nil {
		return a.markTaskFailed(ctx, task, "task ticket is missing")
	}

	resp, err := a.Gateway.Await(ctx, task.Request, task.Ticket)
	if err != nil {
		return a.scheduleRetryOrFail(ctx, task, fmt.Sprintf("await agent result failed: %v", err), true)
	}

	return a.consumeGatewayResponse(ctx, adapter, task, resp)
}

func (a *IMApplicationService) consumeGatewayResponse(ctx context.Context, adapter imService.PlatformAdapter, task *imEntity.TaskRecord, resp *imEntity.AgentResponse) error {
	if resp == nil {
		return a.scheduleRetryOrFail(ctx, task, "empty agent response", true)
	}

	if resp.Task != nil {
		task.Ticket = resp.Task
		task.UpdatedAtMS = nowMillis()
		if err := a.TaskRepo.Save(ctx, task); err != nil {
			return err
		}
	}

	switch resp.Status {
	case imEntity.AgentResponseStatusPending:
		return a.scheduleRetryOrFail(ctx, task, "agent response is still pending", true)
	case imEntity.AgentResponseStatusSuccess, imEntity.AgentResponseStatusFailed:
		task.Result = resp
		task.LastError = ""
		task.NextRetryAtMS = 0
		task.UpdatedAtMS = nowMillis()
		if err := a.TaskRepo.Save(ctx, task); err != nil {
			return err
		}

		return a.deliverTaskResult(ctx, adapter, task, resp)
	default:
		return a.scheduleRetryOrFail(ctx, task, fmt.Sprintf("unsupported agent response status: %s", resp.Status), false)
	}
}

func (a *IMApplicationService) deliverTaskResult(ctx context.Context, adapter imService.PlatformAdapter, task *imEntity.TaskRecord, resp *imEntity.AgentResponse) error {
	if task.Message == nil {
		return a.markTaskFailed(ctx, task, "task message context is missing")
	}

	if err := adapter.SendReply(ctx, task.Message, normalizeTaskReply(resp)); err != nil {
		return a.scheduleRetryOrFail(ctx, task, fmt.Sprintf("push final im reply failed: %v", err), true)
	}

	nowMS := nowMillis()
	task.ResultDelivered = true
	task.DeliveredAtMS = nowMS
	task.FinishedAtMS = nowMS
	task.UpdatedAtMS = nowMS
	task.NextRetryAtMS = 0
	task.LastError = ""

	switch resp.Status {
	case imEntity.AgentResponseStatusSuccess:
		task.Status = imEntity.TaskStatusSuccess
	case imEntity.AgentResponseStatusFailed:
		task.Status = imEntity.TaskStatusFailed
	default:
		task.Status = imEntity.TaskStatusRunning
	}

	return a.TaskRepo.Save(ctx, task)
}

func (a *IMApplicationService) scheduleRetryOrFail(ctx context.Context, task *imEntity.TaskRecord, reason string, retryable bool) error {
	if retryable && task.CanRetry() {
		task.RetryCount++
		task.Status = imEntity.TaskStatusRetrying
		task.LastError = reason
		task.NextRetryAtMS = nowMillis() + int64(nextRetryDelay(task.RetryCount)/time.Millisecond)
		task.UpdatedAtMS = nowMillis()
		if err := a.TaskRepo.Save(ctx, task); err != nil {
			return err
		}

		return a.TaskExecutor.SubmitAfter(ctx, task.ID, nextRetryDelay(task.RetryCount))
	}

	if adapter, err := a.IMDomainSVC.GetAdapter(task.Platform); err == nil && task.Message != nil {
		_ = adapter.SendReply(ctx, task.Message, &imEntity.AgentResponse{
			Status: imEntity.AgentResponseStatusFailed,
			Text:   firstNonEmpty(strings.TrimSpace(reason), "task execution failed"),
		})
	}

	return a.markTaskFailed(ctx, task, reason)
}

func (a *IMApplicationService) markTaskFailed(ctx context.Context, task *imEntity.TaskRecord, reason string) error {
	nowMS := nowMillis()
	task.Status = imEntity.TaskStatusFailed
	task.LastError = reason
	task.NextRetryAtMS = 0
	task.FinishedAtMS = nowMS
	task.UpdatedAtMS = nowMS

	return a.TaskRepo.Save(ctx, task)
}

func buildAcceptedReply(task *imEntity.TaskRecord) *imEntity.AgentResponse {
	return &imEntity.AgentResponse{
		Status: imEntity.AgentResponseStatusPending,
		Text:   fmt.Sprintf("Task accepted, processing.\nTask ID: %s", task.ID),
	}
}

func normalizeTaskReply(resp *imEntity.AgentResponse) *imEntity.AgentResponse {
	if resp == nil {
		return &imEntity.AgentResponse{
			Status: imEntity.AgentResponseStatusFailed,
			Text:   "task execution failed",
		}
	}
	if strings.TrimSpace(resp.Text) != "" {
		return resp
	}

	result := *resp
	switch resp.Status {
	case imEntity.AgentResponseStatusSuccess:
		result.Text = "task completed"
	case imEntity.AgentResponseStatusFailed:
		result.Text = "task execution failed"
	case imEntity.AgentResponseStatusPending:
		result.Text = "task is still running"
	}

	return &result
}

func buildIdempotencyKey(msg *imEntity.IMessage) string {
	parts := []string{
		string(msg.Platform),
		firstNonEmpty(msg.EventID, msg.MessageID),
		msg.ExternalBotID,
		msg.ChatID,
		msg.ThreadID,
		msg.UserID,
	}

	return strings.Join(parts, ":")
}

func buildTaskID(msg *imEntity.IMessage) string {
	sum := sha256.Sum256([]byte(buildIdempotencyKey(msg)))
	return hex.EncodeToString(sum[:])[:24]
}

func nextRetryDelay(retryCount int32) time.Duration {
	delay := taskRetryBaseDelay * time.Duration(1<<(maxInt32(retryCount-1, 0)))
	if delay > time.Minute {
		return time.Minute
	}

	return delay
}

func nowMillis() int64 {
	return time.Now().UnixMilli()
}

func maxInt32(v int32, floor int32) int32 {
	if v < floor {
		return floor
	}

	return v
}

func cloneMessage(msg *imEntity.IMessage) *imEntity.IMessage {
	if msg == nil {
		return nil
	}

	cloned := *msg
	cloned.Metadata = cloneStringMap(msg.Metadata)
	return &cloned
}

func cloneAgentRequest(req *imEntity.AgentRequest) *imEntity.AgentRequest {
	if req == nil {
		return nil
	}

	cloned := *req
	cloned.Metadata = cloneStringMap(req.Metadata)
	if len(req.Parameters) > 0 {
		cloned.Parameters = make(map[string]any, len(req.Parameters))
		for key, value := range req.Parameters {
			cloned.Parameters[key] = value
		}
	}

	return &cloned
}
