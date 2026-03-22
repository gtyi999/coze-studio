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

type AgentTargetType string

const (
	AgentTargetTypeAgent    AgentTargetType = "agent"
	AgentTargetTypeChat     AgentTargetType = "chat"
	AgentTargetTypeWorkflow AgentTargetType = "workflow"
)

type TaskKind string

const (
	TaskKindChatRun     TaskKind = "chat_run"
	TaskKindWorkflowRun TaskKind = "workflow_run"
)

type AgentRequest struct {
	TargetType    AgentTargetType
	SessionID     string
	ConnectorID   int64
	CreatorID     int64
	TargetID      int64
	AppID         *int64
	ProjectID     *int64
	RuntimeUserID string
	Query         string
	Parameters    map[string]any
	Metadata      map[string]string
	PreferAsync   bool
	WaitTimeoutMS int64
}

func (r *AgentRequest) EffectiveWaitTimeoutMS() int64 {
	if r == nil || r.WaitTimeoutMS <= 0 {
		return 8000
	}

	return r.WaitTimeoutMS
}
