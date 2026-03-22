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

package service

import (
	"fmt"

	"github.com/coze-dev/coze-studio/backend/domain/im/entity"
)

type serviceImpl struct {
	adapters map[entity.Platform]PlatformAdapter
}

func NewService(adapters ...PlatformAdapter) Service {
	svc := &serviceImpl{
		adapters: make(map[entity.Platform]PlatformAdapter, len(adapters)),
	}

	for _, adapter := range adapters {
		if adapter == nil {
			continue
		}
		svc.adapters[adapter.Platform()] = adapter
	}

	return svc
}

func (s *serviceImpl) GetAdapter(platform entity.Platform) (PlatformAdapter, error) {
	adapter, ok := s.adapters[platform]
	if !ok {
		return nil, fmt.Errorf("unsupported im platform: %s", platform)
	}

	return adapter, nil
}

func (s *serviceImpl) ListAdapters() []PlatformAdapter {
	adapters := make([]PlatformAdapter, 0, len(entity.AllPlatforms()))
	for _, platform := range entity.AllPlatforms() {
		adapter, ok := s.adapters[platform]
		if !ok {
			continue
		}
		adapters = append(adapters, adapter)
	}

	return adapters
}
