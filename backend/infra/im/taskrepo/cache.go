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

package taskrepo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	imEntity "github.com/coze-dev/coze-studio/backend/domain/im/entity"
	imRepo "github.com/coze-dev/coze-studio/backend/domain/im/repository"
	"github.com/coze-dev/coze-studio/backend/infra/cache"
	"github.com/coze-dev/coze-studio/backend/pkg/sonic"
)

const (
	taskKeyPattern      = "im:task:%s"
	taskLockKeyPattern  = "im:task:lock:%s"
	defaultTaskCacheTTL = 7 * 24 * time.Hour
)

type cacheRepository struct {
	cacheCli cache.Cmdable
}

func New(cacheCli cache.Cmdable) imRepo.TaskRepository {
	return &cacheRepository{
		cacheCli: cacheCli,
	}
}

func (r *cacheRepository) Save(ctx context.Context, task *imEntity.TaskRecord) error {
	if task == nil || strings.TrimSpace(task.ID) == "" {
		return errors.New("im task is invalid")
	}

	payload, err := sonic.MarshalString(task)
	if err != nil {
		return err
	}

	return r.cacheCli.Set(ctx, fmt.Sprintf(taskKeyPattern, task.ID), payload, defaultTaskCacheTTL).Err()
}

func (r *cacheRepository) Get(ctx context.Context, taskID string) (*imEntity.TaskRecord, bool, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, false, nil
	}

	cmd := r.cacheCli.Get(ctx, fmt.Sprintf(taskKeyPattern, taskID))
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), cache.Nil) {
			return nil, false, nil
		}
		return nil, false, cmd.Err()
	}

	var task imEntity.TaskRecord
	if err := sonic.UnmarshalString(cmd.Val(), &task); err != nil {
		return nil, false, err
	}

	return &task, true, nil
}

func (r *cacheRepository) TryAcquireExecution(ctx context.Context, taskID string, ttl time.Duration) (bool, error) {
	if strings.TrimSpace(taskID) == "" {
		return false, errors.New("task id is empty")
	}

	counter, err := r.cacheCli.Incr(ctx, fmt.Sprintf(taskLockKeyPattern, taskID)).Result()
	if err != nil {
		return false, err
	}
	if _, expErr := r.cacheCli.Expire(ctx, fmt.Sprintf(taskLockKeyPattern, taskID), ttl).Result(); expErr != nil {
		return false, expErr
	}

	return counter == 1, nil
}

func (r *cacheRepository) ReleaseExecution(ctx context.Context, taskID string) error {
	if strings.TrimSpace(taskID) == "" {
		return nil
	}

	_, err := r.cacheCli.Del(ctx, fmt.Sprintf(taskLockKeyPattern, taskID)).Result()
	return err
}
