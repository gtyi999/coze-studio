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

type ChannelStatus string

const (
	ChannelStatusEnabled  ChannelStatus = "enabled"
	ChannelStatusDisabled ChannelStatus = "disabled"
)

type SessionScope string

const (
	SessionScopeUser     SessionScope = "user"
	SessionScopeChat     SessionScope = "chat"
	SessionScopeThread   SessionScope = "thread"
	SessionScopeChatUser SessionScope = "chat_user"
)

type ChannelConfig struct {
	ConfigID             string            `json:"config_id"`
	Platform             Platform          `json:"platform"`
	Name                 string            `json:"name"`
	SpaceID              string            `json:"space_id"`
	BotID                int64             `json:"bot_id"`
	ConnectorID          int64             `json:"connector_id"`
	TenantKey            string            `json:"tenant_key"`
	AppID                string            `json:"app_id"`
	BotCode              string            `json:"bot_code"`
	SessionScope         SessionScope      `json:"session_scope"`
	Status               ChannelStatus     `json:"status"`
	PlatformConfig       string            `json:"platform_config"`
	MaskedPlatformConfig string            `json:"masked_platform_config"`
	Ext                  map[string]string `json:"ext,omitempty"`
	CreatedAtMS          int64             `json:"created_at_ms"`
	UpdatedAtMS          int64             `json:"updated_at_ms"`
}
