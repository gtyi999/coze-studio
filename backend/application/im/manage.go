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
	"sort"
	"strings"

	"github.com/google/uuid"

	saEntity "github.com/coze-dev/coze-studio/backend/domain/agent/singleagent/entity"
	imEntity "github.com/coze-dev/coze-studio/backend/domain/im/entity"
)

type ListChannelConfigOptions struct {
	SpaceID  string
	Platform string
	Status   string
	Keyword  string
}

type ListTaskOptions struct {
	SpaceID  string
	Platform string
	Status   string
	ConfigID string
	TaskID   string
}

type ConnectivityTestRequest struct {
	SpaceID     string
	ConfigID    string
	SenderID    string
	SenderName  string
	ChatID      string
	ThreadID    string
	MessageText string
	IsAtBot     bool
	PreferAsync bool
}

func (a *IMApplicationService) GetChannelConfig(ctx context.Context, configID string) (*imEntity.ChannelConfig, bool, error) {
	return a.ChannelConfigRepo.Get(ctx, configID)
}

func (a *IMApplicationService) ListChannelConfigs(ctx context.Context, opt *ListChannelConfigOptions) ([]*imEntity.ChannelConfig, error) {
	if opt == nil || strings.TrimSpace(opt.SpaceID) == "" {
		return nil, nil
	}

	configs, err := a.ChannelConfigRepo.ListBySpace(ctx, opt.SpaceID)
	if err != nil {
		return nil, err
	}

	filtered := make([]*imEntity.ChannelConfig, 0, len(configs))
	for _, cfg := range configs {
		if cfg == nil {
			continue
		}
		if opt.Platform != "" && string(cfg.Platform) != opt.Platform {
			continue
		}
		if opt.Status != "" && string(cfg.Status) != opt.Status {
			continue
		}
		if opt.Keyword != "" {
			keyword := strings.ToLower(opt.Keyword)
			if !strings.Contains(strings.ToLower(cfg.Name), keyword) &&
				!strings.Contains(strings.ToLower(cfg.ConfigID), keyword) {
				continue
			}
		}
		filtered = append(filtered, fillChannelConfigDefaults(cfg))
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].UpdatedAtMS > filtered[j].UpdatedAtMS
	})

	return filtered, nil
}

func (a *IMApplicationService) CreateChannelConfig(ctx context.Context, cfg *imEntity.ChannelConfig) (*imEntity.ChannelConfig, error) {
	if cfg == nil {
		return nil, errors.New("channel config is nil")
	}
	if strings.TrimSpace(cfg.SpaceID) == "" {
		return nil, errors.New("space id is required")
	}

	nowMS := nowMillis()
	cfg.ConfigID = uuid.NewString()
	cfg.CreatedAtMS = nowMS
	cfg.UpdatedAtMS = nowMS
	cfg = fillChannelConfigDefaults(cfg)

	return cfg, a.ChannelConfigRepo.Save(ctx, cfg)
}

func (a *IMApplicationService) UpdateChannelConfig(ctx context.Context, cfg *imEntity.ChannelConfig) (*imEntity.ChannelConfig, error) {
	if cfg == nil || strings.TrimSpace(cfg.ConfigID) == "" {
		return nil, errors.New("config id is required")
	}

	current, found, err := a.ChannelConfigRepo.Get(ctx, cfg.ConfigID)
	if err != nil {
		return nil, err
	}
	if !found || current == nil {
		return nil, errors.New("channel config not found")
	}

	merged := *current
	if cfg.Name != "" {
		merged.Name = cfg.Name
	}
	if cfg.BotID != 0 {
		merged.BotID = cfg.BotID
	}
	if cfg.ConnectorID != 0 {
		merged.ConnectorID = cfg.ConnectorID
	}
	if cfg.TenantKey != "" {
		merged.TenantKey = cfg.TenantKey
	}
	if cfg.AppID != "" {
		merged.AppID = cfg.AppID
	}
	if cfg.BotCode != "" {
		merged.BotCode = cfg.BotCode
	}
	if cfg.SessionScope != "" {
		merged.SessionScope = cfg.SessionScope
	}
	if cfg.Status != "" {
		merged.Status = cfg.Status
	}
	if cfg.PlatformConfig != "" {
		merged.PlatformConfig = cfg.PlatformConfig
	}
	if cfg.Ext != nil {
		merged.Ext = cloneStringMap(cfg.Ext)
	}
	merged.UpdatedAtMS = nowMillis()
	merged.MaskedPlatformConfig = merged.PlatformConfig
	merged = *fillChannelConfigDefaults(&merged)

	return &merged, a.ChannelConfigRepo.Save(ctx, &merged)
}

func (a *IMApplicationService) ListTasks(ctx context.Context, opt *ListTaskOptions) ([]*imEntity.TaskRecord, error) {
	if opt == nil || strings.TrimSpace(opt.SpaceID) == "" {
		return nil, nil
	}

	tasks, err := a.TaskRepo.ListBySpace(ctx, opt.SpaceID)
	if err != nil {
		return nil, err
	}

	filtered := make([]*imEntity.TaskRecord, 0, len(tasks))
	for _, task := range tasks {
		if task == nil {
			continue
		}
		if opt.Platform != "" && string(task.Platform) != opt.Platform {
			continue
		}
		if opt.Status != "" && string(task.Status) != opt.Status {
			continue
		}
		if opt.ConfigID != "" && task.ConfigID != opt.ConfigID {
			continue
		}
		if opt.TaskID != "" && task.ID != opt.TaskID {
			continue
		}
		filtered = append(filtered, task)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].UpdatedAtMS > filtered[j].UpdatedAtMS
	})

	return filtered, nil
}

func (a *IMApplicationService) RunConnectivityTest(ctx context.Context, req *ConnectivityTestRequest) (*imEntity.TaskRecord, error) {
	if req == nil {
		return nil, errors.New("connectivity request is nil")
	}
	if strings.TrimSpace(req.ConfigID) == "" {
		return nil, errors.New("config id is required")
	}
	if strings.TrimSpace(req.MessageText) == "" {
		return nil, errors.New("message text is required")
	}
	if strings.TrimSpace(req.SenderID) == "" {
		return nil, errors.New("sender id is required")
	}

	cfg, found, err := a.ChannelConfigRepo.Get(ctx, req.ConfigID)
	if err != nil {
		return nil, err
	}
	if !found || cfg == nil {
		return nil, errors.New("channel config not found")
	}
	if req.SpaceID != "" && cfg.SpaceID != req.SpaceID {
		return nil, errors.New("channel config does not belong to the current space")
	}
	if cfg.BotID == 0 {
		return nil, errors.New("bot id is required for connectivity test")
	}

	adapter, err := a.IMDomainSVC.GetAdapter(cfg.Platform)
	if err != nil {
		return nil, err
	}

	agent, err := a.SingleAgentDomainSVC.ObtainAgentByIdentity(ctx, &saEntity.AgentIdentity{
		AgentID:     cfg.BotID,
		IsDraft:     false,
		ConnectorID: cfg.ConnectorID,
	})
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return nil, errors.New("published bot not found for the selected connector")
	}

	message := &imEntity.IMessage{
		Platform:      cfg.Platform,
		AgentID:       cfg.BotID,
		ExternalBotID: firstNonEmpty(cfg.BotCode, cfg.AppID),
		EventID:       uuid.NewString(),
		MessageID:     uuid.NewString(),
		SessionID:     buildConnectivitySessionID(cfg.SessionScope, req.SenderID, req.ChatID, req.ThreadID),
		UserID:        req.SenderID,
		UserName:      req.SenderName,
		ChatID:        req.ChatID,
		ThreadID:      req.ThreadID,
		Text:          req.MessageText,
		IsAtBot:       req.IsAtBot,
		Metadata: map[string]string{
			"source":    "im_connectivity_test",
			"space_id":  cfg.SpaceID,
			"config_id": cfg.ConfigID,
		},
	}

	agentReq := a.buildAgentRequest(adapter, agent, message)
	agentReq.PreferAsync = req.PreferAsync
	agentReq.Metadata["space_id"] = cfg.SpaceID
	agentReq.Metadata["config_id"] = cfg.ConfigID

	task, _, err := a.CreateTask(ctx, message, agentReq)
	if err != nil {
		return nil, err
	}

	task.SpaceID = cfg.SpaceID
	task.ConfigID = cfg.ConfigID
	task.PreviewOnly = true
	task.UpdatedAtMS = nowMillis()
	if err = a.TaskRepo.Save(ctx, task); err != nil {
		return nil, err
	}

	if req.PreferAsync {
		return task, a.TaskExecutor.Submit(ctx, task.ID)
	}

	if err = a.ExecuteTask(ctx, task.ID); err != nil {
		return nil, err
	}

	latest, found, err := a.TaskRepo.Get(ctx, task.ID)
	if err != nil {
		return nil, err
	}
	if !found || latest == nil {
		return task, nil
	}

	return latest, nil
}

func fillChannelConfigDefaults(cfg *imEntity.ChannelConfig) *imEntity.ChannelConfig {
	if cfg == nil {
		return nil
	}
	if cfg.SessionScope == "" {
		cfg.SessionScope = imEntity.SessionScopeChat
	}
	if cfg.Status == "" {
		cfg.Status = imEntity.ChannelStatusEnabled
	}
	if cfg.ConnectorID == 0 {
		cfg.ConnectorID = cfg.Platform.ConnectorID()
	}
	if cfg.Ext == nil {
		cfg.Ext = map[string]string{}
	} else {
		cfg.Ext = cloneStringMap(cfg.Ext)
	}
	if cfg.MaskedPlatformConfig == "" {
		cfg.MaskedPlatformConfig = cfg.PlatformConfig
	}
	return cfg
}

func buildConnectivitySessionID(scope imEntity.SessionScope, senderID, chatID, threadID string) string {
	switch scope {
	case imEntity.SessionScopeUser:
		return firstNonEmpty(senderID, uuid.NewString())
	case imEntity.SessionScopeThread:
		return firstNonEmpty(threadID, chatID, senderID, uuid.NewString())
	case imEntity.SessionScopeChatUser:
		return firstNonEmpty(chatID, "chat") + ":" + firstNonEmpty(senderID, "user")
	case imEntity.SessionScopeChat:
		fallthrough
	default:
		return firstNonEmpty(chatID, senderID, uuid.NewString())
	}
}
