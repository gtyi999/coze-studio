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
	"encoding/json"
	"os"
	"strings"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

type Config struct {
	AppKey    string
	AppSecret string
	Token     string
	AESKey    string
	AgentMap  map[string]int64
}

func LoadConfig() Config {
	return Config{
		AppKey:    strings.TrimSpace(os.Getenv(consts.IMDingTalkAppKeyEnv)),
		AppSecret: strings.TrimSpace(os.Getenv(consts.IMDingTalkAppSecretEnv)),
		Token:     strings.TrimSpace(os.Getenv(consts.IMDingTalkTokenEnv)),
		AESKey:    strings.TrimSpace(os.Getenv(consts.IMDingTalkAESKeyEnv)),
		AgentMap:  parseInt64MapEnv(consts.IMDingTalkAgentMapEnv),
	}
}

func (c Config) Enabled() bool {
	return len(c.AgentMap) > 0
}

func (c Config) AgentID(externalBotID string) (int64, bool) {
	agentID, ok := c.AgentMap[externalBotID]
	return agentID, ok
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
