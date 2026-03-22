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

import (
	"fmt"

	"github.com/coze-dev/coze-studio/backend/types/consts"
)

type Platform string

const (
	PlatformFeishu   Platform = "feishu"
	PlatformDingTalk Platform = "dingtalk"
	PlatformWeCom    Platform = "wecom"
)

func ParsePlatform(platform string) (Platform, bool) {
	switch Platform(platform) {
	case PlatformFeishu, PlatformDingTalk, PlatformWeCom:
		return Platform(platform), true
	default:
		return "", false
	}
}

func AllPlatforms() []Platform {
	return []Platform{PlatformFeishu, PlatformDingTalk, PlatformWeCom}
}

func (p Platform) CallbackPath() string {
	switch p {
	case PlatformFeishu:
		return "/api/im/feishu/event"
	case PlatformDingTalk:
		return "/api/im/dingtalk/event"
	case PlatformWeCom:
		return "/api/im/wecom/event"
	default:
		return fmt.Sprintf("/api/im/%s/event", p)
	}
}

func (p Platform) DisplayName() string {
	switch p {
	case PlatformFeishu:
		return "Feishu"
	case PlatformDingTalk:
		return "DingTalk"
	case PlatformWeCom:
		return "WeCom"
	default:
		return string(p)
	}
}

func (p Platform) ConnectorID() int64 {
	switch p {
	case PlatformFeishu:
		return consts.FeishuConnectorID
	case PlatformDingTalk:
		return consts.DingTalkConnectorID
	case PlatformWeCom:
		return consts.WeComConnectorID
	default:
		return 0
	}
}
