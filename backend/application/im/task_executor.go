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
	"strings"
	"time"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/pkg/safego"
)

type AsyncTaskExecutor interface {
	Submit(ctx context.Context, taskID string) error
	SubmitAfter(ctx context.Context, taskID string, delay time.Duration) error
}

type localAsyncTaskExecutor struct {
	run func(context.Context, string) error
}

func NewLocalAsyncTaskExecutor(run func(context.Context, string) error) AsyncTaskExecutor {
	return &localAsyncTaskExecutor{run: run}
}

func (e *localAsyncTaskExecutor) Submit(ctx context.Context, taskID string) error {
	return e.submit(ctx, taskID, 0)
}

func (e *localAsyncTaskExecutor) SubmitAfter(ctx context.Context, taskID string, delay time.Duration) error {
	return e.submit(ctx, taskID, delay)
}

func (e *localAsyncTaskExecutor) submit(ctx context.Context, taskID string, delay time.Duration) error {
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil
	}

	bgCtx := context.WithoutCancel(ctx)
	safego.Go(bgCtx, func() {
		if delay > 0 {
			timer := time.NewTimer(delay)
			defer timer.Stop()

			select {
			case <-bgCtx.Done():
				return
			case <-timer.C:
			}
		}

		if err := e.run(bgCtx, taskID); err != nil {
			logs.CtxErrorf(bgCtx, "execute im task failed, task_id=%s err=%v", taskID, err)
		}
	})

	return nil
}
