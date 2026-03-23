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

package coze

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	protocolconsts "github.com/cloudwego/hertz/pkg/protocol/consts"

	imApp "github.com/coze-dev/coze-studio/backend/application/im"
	imEntity "github.com/coze-dev/coze-studio/backend/domain/im/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/sonic"
)

type listIMChannelConfigsRequest struct {
	SpaceID  string `query:"space_id"`
	Platform string `query:"platform"`
	Status   string `query:"status"`
	Keyword  string `query:"keyword"`
}

type getIMChannelConfigRequest struct {
	ConfigID string `query:"config_id"`
}

type createIMChannelConfigRequest struct {
	Platform       string            `json:"platform"`
	Name           string            `json:"name"`
	SpaceID        string            `json:"space_id"`
	BotID          string            `json:"bot_id"`
	ConnectorID    string            `json:"connector_id"`
	TenantKey      string            `json:"tenant_key"`
	AppID          string            `json:"app_id"`
	BotCode        string            `json:"bot_code"`
	SessionScope   string            `json:"session_scope"`
	Status         string            `json:"status"`
	PlatformConfig string            `json:"platform_config"`
	Ext            map[string]string `json:"ext"`
}

type updateIMChannelConfigRequest struct {
	ConfigID       string            `json:"config_id"`
	Name           string            `json:"name"`
	BotID          string            `json:"bot_id"`
	ConnectorID    string            `json:"connector_id"`
	TenantKey      string            `json:"tenant_key"`
	AppID          string            `json:"app_id"`
	BotCode        string            `json:"bot_code"`
	SessionScope   string            `json:"session_scope"`
	Status         string            `json:"status"`
	PlatformConfig string            `json:"platform_config"`
	Ext            map[string]string `json:"ext"`
}

type testIMChannelConnectivityRequest struct {
	ConfigID    string `json:"config_id"`
	SpaceID     string `json:"space_id"`
	SenderID    string `json:"sender_id"`
	SenderName  string `json:"sender_name"`
	ChatID      string `json:"chat_id"`
	ThreadID    string `json:"thread_id"`
	MessageText string `json:"message_text"`
	IsAtBot     bool   `json:"is_at_bot"`
	PreferAsync bool   `json:"prefer_async"`
}

type listIMTaskRecordsRequest struct {
	SpaceID  string `query:"space_id"`
	Platform string `query:"platform"`
	Status   string `query:"status"`
	ConfigID string `query:"config_id"`
	TaskID   string `query:"task_id"`
}

type getIMTaskDetailRequest struct {
	TaskID string `query:"task_id"`
}

type retryIMTaskRequest struct {
	TaskID string `json:"task_id"`
}

type imChannelConfigResponse struct {
	ConfigID             string            `json:"config_id,omitempty"`
	Platform             string            `json:"platform,omitempty"`
	Name                 string            `json:"name,omitempty"`
	SpaceID              string            `json:"space_id,omitempty"`
	BotID                string            `json:"bot_id,omitempty"`
	ConnectorID          string            `json:"connector_id,omitempty"`
	TenantKey            string            `json:"tenant_key,omitempty"`
	AppID                string            `json:"app_id,omitempty"`
	BotCode              string            `json:"bot_code,omitempty"`
	CallbackPath         string            `json:"callback_path,omitempty"`
	CallbackURL          string            `json:"callback_url,omitempty"`
	SessionScope         string            `json:"session_scope,omitempty"`
	Status               string            `json:"status,omitempty"`
	PlatformConfig       string            `json:"platform_config,omitempty"`
	MaskedPlatformConfig string            `json:"masked_platform_config,omitempty"`
	Ext                  map[string]string `json:"ext,omitempty"`
	CreatedAt            string            `json:"created_at,omitempty"`
	UpdatedAt            string            `json:"updated_at,omitempty"`
}

type imTaskResponse struct {
	TaskID          string            `json:"task_id,omitempty"`
	Platform        string            `json:"platform,omitempty"`
	ConfigID        string            `json:"config_id,omitempty"`
	TaskType        string            `json:"task_type,omitempty"`
	Status          string            `json:"status,omitempty"`
	ConversationID  string            `json:"conversation_id,omitempty"`
	RunID           string            `json:"run_id,omitempty"`
	BotID           string            `json:"bot_id,omitempty"`
	EventID         string            `json:"event_id,omitempty"`
	SessionID       string            `json:"session_id,omitempty"`
	RetryCount      int32             `json:"retry_count,omitempty"`
	MaxRetryCount   int32             `json:"max_retry_count,omitempty"`
	NextRetryAt     string            `json:"next_retry_at,omitempty"`
	ErrorCode       int32             `json:"error_code,omitempty"`
	ErrorMsg        string            `json:"error_msg,omitempty"`
	TraceID         string            `json:"trace_id,omitempty"`
	Ext             map[string]string `json:"ext,omitempty"`
	MessageSnapshot string            `json:"message_snapshot,omitempty"`
	RequestSnapshot string            `json:"request_snapshot,omitempty"`
	GatewayTicket   string            `json:"gateway_ticket,omitempty"`
	ResultSnapshot  string            `json:"result_snapshot,omitempty"`
	CreatedAt       string            `json:"created_at,omitempty"`
	UpdatedAt       string            `json:"updated_at,omitempty"`
}

type imConnectivityResponse struct {
	Accepted     bool   `json:"accepted"`
	Status       string `json:"status,omitempty"`
	TaskID       string `json:"task_id,omitempty"`
	TraceID      string `json:"trace_id,omitempty"`
	Message      string `json:"message,omitempty"`
	ReplyPreview string `json:"reply_preview,omitempty"`
}

func ListIMChannelConfigs(ctx context.Context, c *app.RequestContext) {
	var req listIMChannelConfigsRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	if strings.TrimSpace(req.SpaceID) == "" {
		writeIMError(c, protocolconsts.StatusBadRequest, "space_id is required")
		return
	}

	configs, err := imApp.IMSVC.ListChannelConfigs(ctx, &imApp.ListChannelConfigOptions{
		SpaceID:  req.SpaceID,
		Platform: req.Platform,
		Status:   req.Status,
		Keyword:  req.Keyword,
	})
	if err != nil {
		writeIMInternalError(c, err)
		return
	}

	platformMap, err := getIMPlatformMap(ctx)
	if err != nil {
		writeIMInternalError(c, err)
		return
	}

	list := make([]*imChannelConfigResponse, 0, len(configs))
	for _, cfg := range configs {
		list = append(list, toIMChannelConfigResponse(cfg, platformMap))
	}

	writeIMSuccess(c, map[string]any{
		"list":  list,
		"total": len(list),
	})
}

func GetIMChannelConfig(ctx context.Context, c *app.RequestContext) {
	var req getIMChannelConfigRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	cfg, found, err := imApp.IMSVC.GetChannelConfig(ctx, req.ConfigID)
	if err != nil {
		writeIMInternalError(c, err)
		return
	}
	if !found || cfg == nil {
		writeIMError(c, protocolconsts.StatusNotFound, "channel config not found")
		return
	}

	platformMap, err := getIMPlatformMap(ctx)
	if err != nil {
		writeIMInternalError(c, err)
		return
	}

	writeIMSuccess(c, toIMChannelConfigResponse(cfg, platformMap))
}

func CreateIMChannelConfig(ctx context.Context, c *app.RequestContext) {
	var req createIMChannelConfigRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	cfg, err := buildCreateChannelConfig(&req)
	if err != nil {
		writeIMError(c, protocolconsts.StatusBadRequest, err.Error())
		return
	}

	created, err := imApp.IMSVC.CreateChannelConfig(ctx, cfg)
	if err != nil {
		writeIMInternalError(c, err)
		return
	}

	platformMap, err := getIMPlatformMap(ctx)
	if err != nil {
		writeIMInternalError(c, err)
		return
	}

	writeIMSuccess(c, toIMChannelConfigResponse(created, platformMap))
}

func UpdateIMChannelConfig(ctx context.Context, c *app.RequestContext) {
	var req updateIMChannelConfigRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	if strings.TrimSpace(req.ConfigID) == "" {
		writeIMError(c, protocolconsts.StatusBadRequest, "config_id is required")
		return
	}

	cfg, err := buildUpdateChannelConfig(&req)
	if err != nil {
		writeIMError(c, protocolconsts.StatusBadRequest, err.Error())
		return
	}

	updated, err := imApp.IMSVC.UpdateChannelConfig(ctx, cfg)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeIMError(c, protocolconsts.StatusNotFound, err.Error())
			return
		}
		writeIMInternalError(c, err)
		return
	}

	platformMap, err := getIMPlatformMap(ctx)
	if err != nil {
		writeIMInternalError(c, err)
		return
	}

	writeIMSuccess(c, toIMChannelConfigResponse(updated, platformMap))
}

func TestIMChannelConnectivity(ctx context.Context, c *app.RequestContext) {
	var req testIMChannelConnectivityRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	task, err := imApp.IMSVC.RunConnectivityTest(ctx, &imApp.ConnectivityTestRequest{
		SpaceID:     req.SpaceID,
		ConfigID:    req.ConfigID,
		SenderID:    req.SenderID,
		SenderName:  req.SenderName,
		ChatID:      req.ChatID,
		ThreadID:    req.ThreadID,
		MessageText: req.MessageText,
		IsAtBot:     req.IsAtBot,
		PreferAsync: req.PreferAsync,
	})
	if err != nil {
		writeIMError(c, protocolconsts.StatusBadRequest, err.Error())
		return
	}

	writeIMSuccess(c, &imConnectivityResponse{
		Accepted:     true,
		Status:       string(task.Status),
		TaskID:       task.ID,
		TraceID:      lookupTaskTraceID(task),
		Message:      firstNonEmpty(task.LastError, connectivityStatusMessage(task)),
		ReplyPreview: replyPreview(task),
	})
}

func ListIMTaskRecords(ctx context.Context, c *app.RequestContext) {
	var req listIMTaskRecordsRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	if strings.TrimSpace(req.SpaceID) == "" {
		writeIMError(c, protocolconsts.StatusBadRequest, "space_id is required")
		return
	}

	tasks, err := imApp.IMSVC.ListTasks(ctx, &imApp.ListTaskOptions{
		SpaceID:  req.SpaceID,
		Platform: req.Platform,
		Status:   req.Status,
		ConfigID: req.ConfigID,
		TaskID:   req.TaskID,
	})
	if err != nil {
		writeIMInternalError(c, err)
		return
	}

	list := make([]*imTaskResponse, 0, len(tasks))
	for _, task := range tasks {
		list = append(list, toIMTaskResponse(task))
	}

	writeIMSuccess(c, map[string]any{
		"list":  list,
		"total": len(list),
	})
}

func GetIMTaskDetail(ctx context.Context, c *app.RequestContext) {
	var req getIMTaskDetailRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	task, found, err := imApp.IMSVC.GetTask(ctx, req.TaskID)
	if err != nil {
		writeIMInternalError(c, err)
		return
	}
	if !found || task == nil {
		writeIMError(c, protocolconsts.StatusNotFound, "task not found")
		return
	}

	writeIMSuccess(c, toIMTaskResponse(task))
}

func RetryIMTask(ctx context.Context, c *app.RequestContext) {
	var req retryIMTaskRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	if strings.TrimSpace(req.TaskID) == "" {
		writeIMError(c, protocolconsts.StatusBadRequest, "task_id is required")
		return
	}

	task, err := imApp.IMSVC.RetryTask(ctx, req.TaskID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeIMError(c, protocolconsts.StatusNotFound, err.Error())
			return
		}
		writeIMInternalError(c, err)
		return
	}

	writeIMSuccess(c, toIMTaskResponse(task))
}

func buildCreateChannelConfig(req *createIMChannelConfigRequest) (*imEntity.ChannelConfig, error) {
	platform, err := parsePlatform(req.Platform)
	if err != nil {
		return nil, err
	}
	botID, err := parseOptionalInt64(req.BotID)
	if err != nil {
		return nil, errors.New("bot_id is invalid")
	}
	connectorID, err := parseOptionalInt64(req.ConnectorID)
	if err != nil {
		return nil, errors.New("connector_id is invalid")
	}

	return &imEntity.ChannelConfig{
		Platform:       platform,
		Name:           req.Name,
		SpaceID:        req.SpaceID,
		BotID:          botID,
		ConnectorID:    connectorID,
		TenantKey:      req.TenantKey,
		AppID:          req.AppID,
		BotCode:        req.BotCode,
		SessionScope:   imEntity.SessionScope(req.SessionScope),
		Status:         imEntity.ChannelStatus(req.Status),
		PlatformConfig: req.PlatformConfig,
		Ext:            req.Ext,
	}, nil
}

func buildUpdateChannelConfig(req *updateIMChannelConfigRequest) (*imEntity.ChannelConfig, error) {
	botID, err := parseOptionalInt64(req.BotID)
	if err != nil {
		return nil, errors.New("bot_id is invalid")
	}
	connectorID, err := parseOptionalInt64(req.ConnectorID)
	if err != nil {
		return nil, errors.New("connector_id is invalid")
	}

	return &imEntity.ChannelConfig{
		ConfigID:       req.ConfigID,
		Name:           req.Name,
		BotID:          botID,
		ConnectorID:    connectorID,
		TenantKey:      req.TenantKey,
		AppID:          req.AppID,
		BotCode:        req.BotCode,
		SessionScope:   imEntity.SessionScope(req.SessionScope),
		Status:         imEntity.ChannelStatus(req.Status),
		PlatformConfig: req.PlatformConfig,
		Ext:            req.Ext,
	}, nil
}

func parsePlatform(value string) (imEntity.Platform, error) {
	platform, ok := imEntity.ParsePlatform(value)
	if !ok {
		return "", errors.New("platform is invalid")
	}
	return platform, nil
}

func parseOptionalInt64(value string) (int64, error) {
	if strings.TrimSpace(value) == "" {
		return 0, nil
	}
	return strconv.ParseInt(value, 10, 64)
}

func getIMPlatformMap(ctx context.Context) (map[string]*imEntity.PlatformInfo, error) {
	platforms, err := imApp.IMSVC.ListPlatformInfo(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*imEntity.PlatformInfo, len(platforms))
	for _, platform := range platforms {
		if platform == nil {
			continue
		}
		result[string(platform.Platform)] = platform
	}

	return result, nil
}

func toIMChannelConfigResponse(cfg *imEntity.ChannelConfig, platformMap map[string]*imEntity.PlatformInfo) *imChannelConfigResponse {
	if cfg == nil {
		return nil
	}

	platformInfo := platformMap[string(cfg.Platform)]
	callbackPath := cfg.Platform.CallbackPath()
	callbackURL := ""
	if platformInfo != nil {
		if platformInfo.CallbackPath != "" {
			callbackPath = platformInfo.CallbackPath
		}
		callbackURL = platformInfo.CallbackURL
	}

	return &imChannelConfigResponse{
		ConfigID:             cfg.ConfigID,
		Platform:             string(cfg.Platform),
		Name:                 cfg.Name,
		SpaceID:              cfg.SpaceID,
		BotID:                formatInt64(cfg.BotID),
		ConnectorID:          formatInt64(cfg.ConnectorID),
		TenantKey:            cfg.TenantKey,
		AppID:                cfg.AppID,
		BotCode:              cfg.BotCode,
		CallbackPath:         callbackPath,
		CallbackURL:          callbackURL,
		SessionScope:         string(cfg.SessionScope),
		Status:               string(cfg.Status),
		PlatformConfig:       cfg.PlatformConfig,
		MaskedPlatformConfig: cfg.MaskedPlatformConfig,
		Ext:                  cfg.Ext,
		CreatedAt:            formatMillis(cfg.CreatedAtMS),
		UpdatedAt:            formatMillis(cfg.UpdatedAtMS),
	}
}

func toIMTaskResponse(task *imEntity.TaskRecord) *imTaskResponse {
	if task == nil {
		return nil
	}

	taskType := ""
	if task.Ticket != nil && task.Ticket.Kind != "" {
		taskType = string(task.Ticket.Kind)
	} else if task.Request != nil {
		taskType = string(task.Request.TargetType)
	}

	runID := ""
	conversationID := ""
	if task.Ticket != nil {
		runID = firstNonEmpty(formatInt64(task.Ticket.RunID), task.Ticket.ExecuteID)
		conversationID = formatInt64(task.Ticket.ConversationID)
	}

	botID := ""
	if task.Request != nil {
		botID = formatInt64(task.Request.TargetID)
	} else if task.Message != nil {
		botID = formatInt64(task.Message.AgentID)
	}

	return &imTaskResponse{
		TaskID:          task.ID,
		Platform:        string(task.Platform),
		ConfigID:        task.ConfigID,
		TaskType:        taskType,
		Status:          string(task.Status),
		ConversationID:  conversationID,
		RunID:           runID,
		BotID:           botID,
		EventID:         messageField(task, func(msg *imEntity.IMessage) string { return msg.EventID }),
		SessionID:       messageField(task, func(msg *imEntity.IMessage) string { return msg.SessionID }),
		RetryCount:      task.RetryCount,
		MaxRetryCount:   task.MaxRetryCount,
		NextRetryAt:     formatMillis(task.NextRetryAtMS),
		ErrorMsg:        task.LastError,
		TraceID:         lookupTaskTraceID(task),
		Ext:             cloneStringMapLocal(messageMetadata(task)),
		MessageSnapshot: marshalSnapshot(task.Message),
		RequestSnapshot: marshalSnapshot(task.Request),
		GatewayTicket:   marshalSnapshot(task.Ticket),
		ResultSnapshot:  marshalSnapshot(task.Result),
		CreatedAt:       formatMillis(task.CreatedAtMS),
		UpdatedAt:       formatMillis(task.UpdatedAtMS),
	}
}

func replyPreview(task *imEntity.TaskRecord) string {
	if task == nil || task.Result == nil {
		return ""
	}
	return task.Result.Text
}

func connectivityStatusMessage(task *imEntity.TaskRecord) string {
	if task == nil {
		return ""
	}
	switch task.Status {
	case imEntity.TaskStatusSuccess:
		return "Connectivity test completed"
	case imEntity.TaskStatusFailed:
		return "Connectivity test failed"
	case imEntity.TaskStatusPending, imEntity.TaskStatusRunning, imEntity.TaskStatusRetrying:
		return "Connectivity test task submitted"
	default:
		return "Connectivity test task submitted"
	}
}

func lookupTaskTraceID(task *imEntity.TaskRecord) string {
	if task == nil {
		return ""
	}
	if task.Result != nil && task.Result.Metadata != nil {
		if value := strings.TrimSpace(task.Result.Metadata["trace_id"]); value != "" {
			return value
		}
	}
	if task.Request != nil && task.Request.Metadata != nil {
		if value := strings.TrimSpace(task.Request.Metadata["trace_id"]); value != "" {
			return value
		}
	}
	if task.Ticket != nil && task.Ticket.Metadata != nil {
		if value := strings.TrimSpace(task.Ticket.Metadata["trace_id"]); value != "" {
			return value
		}
	}
	return ""
}

func messageMetadata(task *imEntity.TaskRecord) map[string]string {
	if task == nil || task.Message == nil {
		return nil
	}
	return task.Message.Metadata
}

func marshalSnapshot(value any) string {
	if value == nil {
		return ""
	}
	text, err := sonic.MarshalString(value)
	if err != nil {
		return ""
	}
	return text
}

func formatMillis(value int64) string {
	if value <= 0 {
		return ""
	}
	return strconv.FormatInt(value, 10)
}

func formatInt64(value int64) string {
	if value <= 0 {
		return ""
	}
	return strconv.FormatInt(value, 10)
}

func messageField(task *imEntity.TaskRecord, fn func(msg *imEntity.IMessage) string) string {
	if task == nil || task.Message == nil {
		return ""
	}
	return fn(task.Message)
}

func cloneStringMapLocal(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func writeIMSuccess(c *app.RequestContext, data any) {
	c.JSON(protocolconsts.StatusOK, map[string]any{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

func writeIMError(c *app.RequestContext, statusCode int, message string) {
	c.JSON(statusCode, map[string]any{
		"code": statusCode,
		"msg":  message,
	})
}

func writeIMInternalError(c *app.RequestContext, err error) {
	writeIMError(c, protocolconsts.StatusInternalServerError, err.Error())
}
