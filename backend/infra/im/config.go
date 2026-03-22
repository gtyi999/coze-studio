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
	"encoding/json"
	"os"
	"strings"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

type FeishuConfig struct {
	AppID             string
	AppSecret         string
	VerificationToken string
	AgentMap          map[string]int64
}

func (c FeishuConfig) Enabled() bool {
	return c.AppID != "" && c.AppSecret != ""
}

func (c FeishuConfig) AgentID(appID string) (int64, bool) {
	agentID, ok := c.AgentMap[appID]
	return agentID, ok
}

type DingTalkConfig struct {
	AgentMap map[string]int64
}

func (c DingTalkConfig) Enabled() bool {
	return len(c.AgentMap) > 0
}

func (c DingTalkConfig) AgentID(robotCode string) (int64, bool) {
	agentID, ok := c.AgentMap[robotCode]
	return agentID, ok
}

type WeComConfig struct {
	CorpID         string
	Token          string
	EncodingAESKey string
	AgentMap       map[string]int64
	AgentSecretMap map[string]string
}

func (c WeComConfig) Enabled() bool {
	return c.CorpID != "" && c.Token != "" && c.EncodingAESKey != ""
}

func (c WeComConfig) AgentID(externalBotID string) (int64, bool) {
	agentID, ok := c.AgentMap[externalBotID]
	return agentID, ok
}

func (c WeComConfig) AgentSecret(agentID string) (string, bool) {
	secret, ok := c.AgentSecretMap[agentID]
	return secret, ok
}

func LoadFeishuConfig() FeishuConfig {
	return FeishuConfig{
		AppID:             strings.TrimSpace(os.Getenv(consts.IMFeishuAppIDEnv)),
		AppSecret:         strings.TrimSpace(os.Getenv(consts.IMFeishuAppSecretEnv)),
		VerificationToken: strings.TrimSpace(os.Getenv(consts.IMFeishuVerificationTokenEnv)),
		AgentMap:          parseInt64MapEnv(consts.IMFeishuAgentMapEnv),
	}
}

func LoadDingTalkConfig() DingTalkConfig {
	return DingTalkConfig{
		AgentMap: parseInt64MapEnv(consts.IMDingTalkAgentMapEnv),
	}
}

func LoadWeComConfig() WeComConfig {
	return WeComConfig{
		CorpID:         strings.TrimSpace(os.Getenv(consts.IMWeComCorpIDEnv)),
		Token:          strings.TrimSpace(os.Getenv(consts.IMWeComTokenEnv)),
		EncodingAESKey: strings.TrimSpace(os.Getenv(consts.IMWeComEncodingAESKeyEnv)),
		AgentMap:       parseInt64MapEnv(consts.IMWeComAgentMapEnv),
		AgentSecretMap: parseStringMapEnv(consts.IMWeComAgentSecretMapEnv),
	}
}

func parseInt64MapEnv(key string) map[string]int64 {
	result := map[string]int64{}
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return result
	}

	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		logs.Warnf("parse env %s failed: %v", key, err)
		return map[string]int64{}
	}

	return result
}

func parseStringMapEnv(key string) map[string]string {
	result := map[string]string{}
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return result
	}

	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		logs.Warnf("parse env %s failed: %v", key, err)
		return map[string]string{}
	}

	return result
}
