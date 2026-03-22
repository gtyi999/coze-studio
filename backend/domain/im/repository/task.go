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

package repository

import (
	"context"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/im/entity"
)

type TaskRepository interface {
	Save(ctx context.Context, task *entity.TaskRecord) error
	Get(ctx context.Context, taskID string) (*entity.TaskRecord, bool, error)
	TryAcquireExecution(ctx context.Context, taskID string, ttl time.Duration) (bool, error)
	ReleaseExecution(ctx context.Context, taskID string) error
}
