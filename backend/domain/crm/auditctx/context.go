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

package auditctx

import "context"

type Actor struct {
	UserID   int64
	TenantID int64
	SpaceID  int64
}

type actorKey struct{}

func WithActor(ctx context.Context, actor *Actor) context.Context {
	if actor == nil {
		return ctx
	}
	return context.WithValue(ctx, actorKey{}, *actor)
}

func ActorFromContext(ctx context.Context) (*Actor, bool) {
	if ctx == nil {
		return nil, false
	}

	actor, ok := ctx.Value(actorKey{}).(Actor)
	if !ok {
		return nil, false
	}
	return &actor, true
}
