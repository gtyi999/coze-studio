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
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/coze-dev/coze-studio/backend/bizpkg/config"
	saEntity "github.com/coze-dev/coze-studio/backend/domain/agent/singleagent/entity"
	singleagent "github.com/coze-dev/coze-studio/backend/domain/agent/singleagent/service"
	imEntity "github.com/coze-dev/coze-studio/backend/domain/im/entity"
	imRepo "github.com/coze-dev/coze-studio/backend/domain/im/repository"
	imService "github.com/coze-dev/coze-studio/backend/domain/im/service"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

type IMApplicationService struct {
	IMDomainSVC          imService.Service
	SingleAgentDomainSVC singleagent.SingleAgent
	ChannelConfigRepo    imRepo.ChannelConfigRepository
	TaskRepo             imRepo.TaskRepository
	TaskExecutor         AsyncTaskExecutor
	Gateway              *CozeAgentGateway
}

func (a *IMApplicationService) ListPlatformInfo(ctx context.Context) ([]*imEntity.PlatformInfo, error) {
	host, err := config.Base().GetServerHost(ctx)
	if err != nil {
		return nil, err
	}

	adapters := a.IMDomainSVC.ListAdapters()
	result := make([]*imEntity.PlatformInfo, 0, len(adapters))
	for _, adapter := range adapters {
		result = append(result, adapter.PlatformInfo(host))
	}

	return result, nil
}

func (a *IMApplicationService) HandleFeishuEvent(ctx context.Context, req *imEntity.CallbackRequest) *imEntity.RawResponse {
	resp, err := a.handleEvent(ctx, imEntity.PlatformFeishu, req)
	if err != nil {
		logs.CtxErrorf(ctx, "handle feishu event failed: %v", err)
		return imEntity.ErrorResponse(http.StatusInternalServerError, "internal server error")
	}
	if resp == nil {
		return imEntity.TextResponse(http.StatusOK, "")
	}

	return resp
}

func (a *IMApplicationService) HandleDingTalkEvent(ctx context.Context, req *imEntity.CallbackRequest) *imEntity.RawResponse {
	resp, err := a.handleEvent(ctx, imEntity.PlatformDingTalk, req)
	if err != nil {
		logs.CtxErrorf(ctx, "handle dingtalk event failed: %v", err)
		return imEntity.ErrorResponse(http.StatusInternalServerError, "internal server error")
	}
	if resp == nil {
		return imEntity.TextResponse(http.StatusOK, "")
	}

	return resp
}

func (a *IMApplicationService) HandleWeComEvent(ctx context.Context, req *imEntity.CallbackRequest) *imEntity.RawResponse {
	resp, err := a.handleEvent(ctx, imEntity.PlatformWeCom, req)
	if err != nil {
		logs.CtxErrorf(ctx, "handle wecom event failed: %v", err)
		return imEntity.ErrorResponse(http.StatusInternalServerError, "internal server error")
	}
	if resp == nil {
		return imEntity.TextResponse(http.StatusOK, "")
	}

	return resp
}

func (a *IMApplicationService) GetTask(ctx context.Context, taskID string) (*imEntity.TaskRecord, bool, error) {
	return a.TaskRepo.Get(ctx, taskID)
}

func (a *IMApplicationService) RetryTask(ctx context.Context, taskID string) (*imEntity.TaskRecord, error) {
	task, found, err := a.TaskRepo.Get(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("im task not found")
	}
	if task.Status != imEntity.TaskStatusFailed {
		return task, nil
	}

	nowMS := nowMillis()
	task.Status = imEntity.TaskStatusRetrying
	task.RetryCount = 0
	task.LastError = ""
	task.NextRetryAtMS = nowMS
	task.DeadlineAtMS = nowMS + int64(defaultTaskDeadline/time.Millisecond)
	task.Result = nil
	task.ResultDelivered = false
	task.DeliveredAtMS = 0
	task.FinishedAtMS = 0
	task.Ticket = nil
	task.UpdatedAtMS = nowMS
	if err = a.TaskRepo.Save(ctx, task); err != nil {
		return nil, err
	}

	return task, a.TaskExecutor.Submit(ctx, task.ID)
}

func (a *IMApplicationService) handleEvent(ctx context.Context, platform imEntity.Platform, req *imEntity.CallbackRequest) (*imEntity.RawResponse, error) {
	adapter, err := a.IMDomainSVC.GetAdapter(platform)
	if err != nil {
		return imEntity.ErrorResponse(http.StatusNotFound, "unsupported im platform"), nil
	}

	msg, rawResp, err := adapter.ParseCallback(ctx, req)
	if err != nil {
		return a.toRawResponse(err), nil
	}
	if rawResp != nil {
		return rawResp, nil
	}
	if msg == nil || !msg.IsAtBot || strings.TrimSpace(msg.Text) == "" {
		return adapter.SuccessResponse(), nil
	}

	agent, err := a.resolveAgent(ctx, adapter, msg)
	if err != nil {
		return nil, err
	}

	agentReq := a.buildAgentRequest(adapter, agent, msg)
	task, existed, err := a.CreateTask(ctx, msg, agentReq)
	if err != nil {
		return nil, err
	}
	if existed {
		return adapter.SuccessResponse(), nil
	}

	if sendErr := adapter.SendReply(ctx, msg, buildAcceptedReply(task)); sendErr != nil {
		logs.CtxWarnf(ctx, "send accepted im reply failed, task_id=%s err=%v", task.ID, sendErr)
	}
	if err = a.TaskExecutor.Submit(ctx, task.ID); err != nil {
		return nil, err
	}

	return adapter.SuccessResponse(), nil
}

func (a *IMApplicationService) resolveAgent(ctx context.Context, adapter imService.PlatformAdapter, msg *imEntity.IMessage) (*saEntity.SingleAgent, error) {
	agent, err := a.SingleAgentDomainSVC.ObtainAgentByIdentity(ctx, &saEntity.AgentIdentity{
		AgentID:     msg.AgentID,
		IsDraft:     false,
		ConnectorID: adapter.ConnectorID(),
	})
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return nil, errors.New("published agent not found for im message")
	}

	return agent, nil
}

func (a *IMApplicationService) buildAgentRequest(adapter imService.PlatformAdapter, agent *saEntity.SingleAgent, msg *imEntity.IMessage) *imEntity.AgentRequest {
	metadata := map[string]string{
		"source":     "im",
		"platform":   string(msg.Platform),
		"event_id":   msg.EventID,
		"message_id": msg.MessageID,
		"chat_id":    msg.ChatID,
		"thread_id":  msg.ThreadID,
		"user_id":    msg.UserID,
		"user_name":  msg.UserName,
	}
	for key, value := range msg.Metadata {
		if strings.TrimSpace(value) == "" {
			continue
		}
		metadata["metadata_"+key] = value
	}

	return &imEntity.AgentRequest{
		TargetType:    imEntity.AgentTargetTypeAgent,
		SessionID:     msg.SessionID,
		ConnectorID:   adapter.ConnectorID(),
		CreatorID:     agent.CreatorID,
		TargetID:      agent.AgentID,
		RuntimeUserID: msg.UserID,
		Query:         msg.Text,
		Metadata:      metadata,
		PreferAsync:   true,
		WaitTimeoutMS: 60000,
	}
}

func (a *IMApplicationService) toRawResponse(err error) *imEntity.RawResponse {
	var statusErr *imEntity.StatusError
	if errors.As(err, &statusErr) {
		return imEntity.ErrorResponse(statusErr.StatusCode, statusErr.Message)
	}

	return imEntity.ErrorResponse(http.StatusInternalServerError, "internal server error")
}
