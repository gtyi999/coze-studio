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
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/coze-dev/coze-studio/backend/api/model/conversation/common"
	messageAPI "github.com/coze-dev/coze-studio/backend/api/model/conversation/message"
	runAPI "github.com/coze-dev/coze-studio/backend/api/model/conversation/run"
	workflowAPI "github.com/coze-dev/coze-studio/backend/api/model/workflow"
	conversationApp "github.com/coze-dev/coze-studio/backend/application/conversation"
	workflowApp "github.com/coze-dev/coze-studio/backend/application/workflow"
	crossmessage "github.com/coze-dev/coze-studio/backend/crossdomain/message/model"
	agentrunEntity "github.com/coze-dev/coze-studio/backend/domain/conversation/agentrun/entity"
	convEntity "github.com/coze-dev/coze-studio/backend/domain/conversation/conversation/entity"
	conversationDomain "github.com/coze-dev/coze-studio/backend/domain/conversation/conversation/service"
	imEntity "github.com/coze-dev/coze-studio/backend/domain/im/entity"
	openapiAuthEntity "github.com/coze-dev/coze-studio/backend/domain/openauth/openapiauth/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/pkg/sonic"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

type CozeAgentGateway struct {
	ConversationDomainSVC conversationDomain.Conversation
}

type CozeAgentGatewayComponents struct {
	ConversationDomainSVC conversationDomain.Conversation
}

func NewCozeAgentGateway(c *CozeAgentGatewayComponents) *CozeAgentGateway {
	return &CozeAgentGateway{
		ConversationDomainSVC: c.ConversationDomainSVC,
	}
}

func (g *CozeAgentGateway) Execute(ctx context.Context, req *imEntity.AgentRequest) (*imEntity.AgentResponse, error) {
	if req == nil {
		return nil, errors.New("agent request is nil")
	}

	execCtx := g.withOpenAPIAuth(ctx, req)

	switch req.TargetType {
	case imEntity.AgentTargetTypeAgent, imEntity.AgentTargetTypeChat:
		return g.executeChat(execCtx, req)
	case imEntity.AgentTargetTypeWorkflow:
		return g.executeWorkflow(execCtx, req)
	default:
		return nil, fmt.Errorf("unsupported agent target type: %s", req.TargetType)
	}
}

func (g *CozeAgentGateway) Await(ctx context.Context, req *imEntity.AgentRequest, ticket *imEntity.TaskTicket) (*imEntity.AgentResponse, error) {
	if req == nil {
		return nil, errors.New("agent request is nil")
	}
	if ticket == nil {
		return nil, errors.New("task ticket is nil")
	}

	execCtx := g.withOpenAPIAuth(ctx, req)

	switch ticket.Kind {
	case imEntity.TaskKindChatRun:
		return g.awaitChat(execCtx, ticket, req.EffectiveWaitTimeoutMS())
	case imEntity.TaskKindWorkflowRun:
		return g.awaitWorkflow(execCtx, ticket, req.EffectiveWaitTimeoutMS())
	default:
		return nil, fmt.Errorf("unsupported task kind: %s", ticket.Kind)
	}
}

func (g *CozeAgentGateway) executeChat(ctx context.Context, req *imEntity.AgentRequest) (*imEntity.AgentResponse, error) {
	// Reuse the existing OpenAPI chat chain so conversation creation and run lifecycle stay consistent.
	conversationID, err := g.resolveChatConversationID(ctx, req)
	if err != nil {
		return nil, err
	}

	resp, err := conversationApp.ConversationOpenAPISVC.OpenapiAgentRunSync(ctx, g.buildChatRequest(req, conversationID))
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.ChatDetail == nil {
		return nil, errors.New("empty chat response")
	}

	ticket := &imEntity.TaskTicket{
		Kind:           imEntity.TaskKindChatRun,
		TargetType:     req.TargetType,
		TargetID:       req.TargetID,
		SessionID:      req.SessionID,
		ConversationID: resp.ChatDetail.ConversationID,
		RunID:          resp.ChatDetail.ID,
		Status:         resp.ChatDetail.Status,
		Metadata:       cloneStringMap(req.Metadata),
	}

	if req.PreferAsync {
		return &imEntity.AgentResponse{
			Status:   imEntity.AgentResponseStatusPending,
			Task:     ticket,
			Metadata: cloneStringMap(req.Metadata),
		}, nil
	}

	return g.awaitChat(ctx, ticket, req.EffectiveWaitTimeoutMS())
}

func (g *CozeAgentGateway) executeWorkflow(ctx context.Context, req *imEntity.AgentRequest) (*imEntity.AgentResponse, error) {
	parameters, err := buildWorkflowParameters(req)
	if err != nil {
		return nil, err
	}

	request := &workflowAPI.OpenAPIRunFlowRequest{
		WorkflowID:  strconv.FormatInt(req.TargetID, 10),
		Ext:         buildWorkflowExt(req),
		IsAsync:     ptr.Of(req.PreferAsync),
		ConnectorID: ptr.Of(strconv.FormatInt(req.ConnectorID, 10)),
	}
	if len(parameters) > 0 {
		rawParameters, marshalErr := sonic.MarshalString(parameters)
		if marshalErr != nil {
			return nil, marshalErr
		}
		request.Parameters = ptr.Of(rawParameters)
	}
	if req.AppID != nil {
		request.AppID = ptr.Of(strconv.FormatInt(*req.AppID, 10))
	}
	if req.ProjectID != nil {
		request.ProjectID = ptr.Of(strconv.FormatInt(*req.ProjectID, 10))
	}

	resp, err := workflowApp.SVC.OpenAPIRun(ctx, request)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("empty workflow response")
	}
	if data := strings.TrimSpace(resp.GetData()); data != "" {
		return &imEntity.AgentResponse{
			Status:   imEntity.AgentResponseStatusSuccess,
			Text:     data,
			Metadata: cloneStringMap(req.Metadata),
		}, nil
	}

	if executeID := strings.TrimSpace(resp.GetExecuteID()); executeID != "" {
		return &imEntity.AgentResponse{
			Status: imEntity.AgentResponseStatusPending,
			Task: &imEntity.TaskTicket{
				Kind:       imEntity.TaskKindWorkflowRun,
				TargetType: req.TargetType,
				TargetID:   req.TargetID,
				SessionID:  req.SessionID,
				ExecuteID:  executeID,
				Status:     "running",
				Metadata:   cloneStringMap(req.Metadata),
			},
			Metadata: cloneStringMap(req.Metadata),
		}, nil
	}

	if msg := strings.TrimSpace(resp.GetMsg()); msg != "" && resp.GetCode() != 0 {
		return &imEntity.AgentResponse{
			Status:   imEntity.AgentResponseStatusFailed,
			Text:     msg,
			Metadata: cloneStringMap(req.Metadata),
		}, nil
	}

	return &imEntity.AgentResponse{
		Status:   imEntity.AgentResponseStatusSuccess,
		Metadata: cloneStringMap(req.Metadata),
	}, nil
}

func (g *CozeAgentGateway) awaitChat(ctx context.Context, ticket *imEntity.TaskTicket, waitTimeoutMS int64) (*imEntity.AgentResponse, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(waitTimeoutMS)*time.Millisecond)
	defer cancel()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		record, err := conversationApp.ConversationOpenAPISVC.RetrieveRunRecord(timeoutCtx, &runAPI.RetrieveChatOpenRequest{
			ConversationID: ticket.ConversationID,
			ChatID:         ticket.RunID,
		})
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(timeoutCtx.Err(), context.DeadlineExceeded) {
				return pendingResponse(ticket), nil
			}
			return nil, err
		}
		if record != nil && record.ChatDetail != nil {
			ticket.Status = record.ChatDetail.Status

			switch record.ChatDetail.Status {
			case string(agentrunEntity.RunStatusCompleted):
				text, loadErr := g.loadChatReplyText(timeoutCtx, ticket)
				if loadErr != nil {
					return nil, loadErr
				}
				return &imEntity.AgentResponse{
					Status: imEntity.AgentResponseStatusSuccess,
					Text:   text,
					Task:   ticket,
				}, nil
			case string(agentrunEntity.RunStatusFailed), string(agentrunEntity.RunStatusExpired), string(agentrunEntity.RunStatusCancelled):
				return &imEntity.AgentResponse{
					Status: imEntity.AgentResponseStatusFailed,
					Text:   "agent run failed",
					Task:   ticket,
				}, nil
			}
		}

		select {
		case <-timeoutCtx.Done():
			return pendingResponse(ticket), nil
		case <-ticker.C:
		}
	}
}

func (g *CozeAgentGateway) awaitWorkflow(ctx context.Context, ticket *imEntity.TaskTicket, waitTimeoutMS int64) (*imEntity.AgentResponse, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(waitTimeoutMS)*time.Millisecond)
	defer cancel()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		resp, err := workflowApp.SVC.OpenAPIGetWorkflowRunHistory(timeoutCtx, &workflowAPI.GetWorkflowRunHistoryRequest{
			WorkflowID: strconv.FormatInt(ticket.TargetID, 10),
			ExecuteID:  ptr.Of(ticket.ExecuteID),
		})
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(timeoutCtx.Err(), context.DeadlineExceeded) {
				return pendingResponse(ticket), nil
			}
			return nil, err
		}
		if len(resp.GetData()) > 0 && resp.GetData()[0] != nil {
			history := resp.GetData()[0]
			ticket.Status = history.GetExecuteStatus()

			switch strings.ToLower(history.GetExecuteStatus()) {
			case "success":
				return &imEntity.AgentResponse{
					Status: imEntity.AgentResponseStatusSuccess,
					Text:   strings.TrimSpace(history.GetOutput()),
					Task:   ticket,
				}, nil
			case "fail", "cancel":
				return &imEntity.AgentResponse{
					Status: imEntity.AgentResponseStatusFailed,
					Text:   firstNonEmpty(strings.TrimSpace(history.GetErrorMsg()), "workflow run failed"),
					Task:   ticket,
				}, nil
			}
		}

		select {
		case <-timeoutCtx.Done():
			return pendingResponse(ticket), nil
		case <-ticker.C:
		}
	}
}

func (g *CozeAgentGateway) resolveChatConversationID(ctx context.Context, req *imEntity.AgentRequest) (*int64, error) {
	if strings.TrimSpace(req.SessionID) == "" {
		return nil, nil
	}

	conversations, _, err := g.ConversationDomainSVC.List(ctx, &convEntity.ListMeta{
		CreatorID:   req.CreatorID,
		UserID:      ptr.Of(req.SessionID),
		ConnectorID: req.ConnectorID,
		Scene:       common.Scene_SceneOpenApi,
		AgentID:     req.TargetID,
		Page:        1,
		Limit:       1,
	})
	if err != nil {
		return nil, err
	}
	if len(conversations) == 0 {
		return nil, nil
	}

	return ptr.Of(conversations[0].ID), nil
}

func (g *CozeAgentGateway) buildChatRequest(req *imEntity.AgentRequest, conversationID *int64) *runAPI.ChatV3Request {
	stream := false
	request := &runAPI.ChatV3Request{
		BotID:       req.TargetID,
		User:        g.runtimeSessionUser(req),
		Stream:      ptr.Of(stream),
		ConnectorID: ptr.Of(req.ConnectorID),
		AdditionalMessages: []*runAPI.EnterMessage{
			{
				Role:        "user",
				Content:     req.Query,
				ContentType: runAPI.ContentTypeText,
				MetaData:    cloneStringMap(req.Metadata),
			},
		},
		ExtraParams: cloneStringMap(req.Metadata),
	}
	if conversationID != nil {
		request.ConversationID = conversationID
	}
	if request.ExtraParams == nil {
		request.ExtraParams = map[string]string{}
	}
	if req.RuntimeUserID != "" {
		request.ExtraParams["runtime_user_id"] = req.RuntimeUserID
	}
	if req.SessionID != "" {
		request.ExtraParams["session_id"] = req.SessionID
	}

	return request
}

func (g *CozeAgentGateway) loadChatReplyText(ctx context.Context, ticket *imEntity.TaskTicket) (string, error) {
	resp, err := conversationApp.ConversationOpenAPISVC.ListChatMessageApi(ctx, &messageAPI.ListChatMessageApiRequest{
		ConversationID: ticket.ConversationID,
		ChatID:         ticket.RunID,
	})
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", nil
	}

	replies := make([]string, 0, len(resp.Messages))
	for _, message := range resp.Messages {
		if message == nil || message.Role != string(agentrunEntity.RoleTypeAssistant) {
			continue
		}

		switch crossmessage.MessageType(message.Type) {
		case crossmessage.MessageTypeAnswer, crossmessage.MessageTypeToolAsAnswer, crossmessage.MessageTypeToolMidAnswer:
			if text := strings.TrimSpace(message.Content); text != "" {
				replies = append(replies, text)
			}
		}
	}

	return strings.TrimSpace(strings.Join(replies, "\n")), nil
}

func (g *CozeAgentGateway) withOpenAPIAuth(ctx context.Context, req *imEntity.AgentRequest) context.Context {
	// Conversation and workflow OpenAPI entrypoints read auth data from ctx cache.
	ctx = ctxcache.Init(ctx)
	ctxcache.Store(ctx, consts.OpenapiAuthKeyInCtx, &openapiAuthEntity.ApiKey{
		UserID:      req.CreatorID,
		ConnectorID: req.ConnectorID,
	})

	return ctx
}

func (g *CozeAgentGateway) runtimeSessionUser(req *imEntity.AgentRequest) string {
	if strings.TrimSpace(req.SessionID) != "" {
		return req.SessionID
	}
	if strings.TrimSpace(req.RuntimeUserID) != "" {
		return req.RuntimeUserID
	}

	return strconv.FormatInt(req.CreatorID, 10)
}

func buildWorkflowParameters(req *imEntity.AgentRequest) (map[string]any, error) {
	parameters := make(map[string]any, len(req.Parameters)+1)
	for key, value := range req.Parameters {
		parameters[key] = value
	}

	if strings.TrimSpace(req.Query) == "" {
		return parameters, nil
	}

	for _, key := range []string{"user_input", "query", "input"} {
		if _, ok := parameters[key]; ok {
			return parameters, nil
		}
	}

	parameters["user_input"] = req.Query
	return parameters, nil
}

func buildWorkflowExt(req *imEntity.AgentRequest) map[string]string {
	ext := cloneStringMap(req.Metadata)
	if ext == nil {
		ext = map[string]string{}
	}
	if req.RuntimeUserID != "" {
		ext["user_id"] = req.RuntimeUserID
	}
	if req.SessionID != "" {
		ext["session_id"] = req.SessionID
	}

	return ext
}

func cloneStringMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}

	out := make(map[string]string, len(in))
	for key, value := range in {
		out[key] = value
	}

	return out
}

func pendingResponse(ticket *imEntity.TaskTicket) *imEntity.AgentResponse {
	return &imEntity.AgentResponse{
		Status: imEntity.AgentResponseStatusPending,
		Task:   ticket,
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}

	return ""
}
