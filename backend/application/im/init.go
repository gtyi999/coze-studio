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
	singleagent "github.com/coze-dev/coze-studio/backend/domain/agent/singleagent/service"
	conversation "github.com/coze-dev/coze-studio/backend/domain/conversation/conversation/service"
	imDomain "github.com/coze-dev/coze-studio/backend/domain/im/service"
	"github.com/coze-dev/coze-studio/backend/infra/cache"
	"github.com/coze-dev/coze-studio/backend/infra/im/factory"
	_ "github.com/coze-dev/coze-studio/backend/infra/im/impl/dingtalk"
	_ "github.com/coze-dev/coze-studio/backend/infra/im/impl/feishu"
	_ "github.com/coze-dev/coze-studio/backend/infra/im/impl/wecom"
	"github.com/coze-dev/coze-studio/backend/infra/im/taskrepo"
)

type ServiceComponents struct {
	ConversationDomainSVC conversation.Conversation
	SingleAgentDomainSVC  singleagent.SingleAgent
	CacheCli              cache.Cmdable
}

var IMSVC *IMApplicationService

func InitService(c *ServiceComponents) *IMApplicationService {
	domainSVC := imDomain.NewService(factory.BuildAll()...)
	taskRepo := taskrepo.New(c.CacheCli)

	IMSVC = &IMApplicationService{
		IMDomainSVC:          domainSVC,
		SingleAgentDomainSVC: c.SingleAgentDomainSVC,
		TaskRepo:             taskRepo,
		Gateway: NewCozeAgentGateway(&CozeAgentGatewayComponents{
			ConversationDomainSVC: c.ConversationDomainSVC,
		}),
	}
	IMSVC.TaskExecutor = NewLocalAsyncTaskExecutor(IMSVC.ExecuteTask)

	return IMSVC
}
