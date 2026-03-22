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

package factory

import (
	"fmt"
	"sync"

	imEntity "github.com/coze-dev/coze-studio/backend/domain/im/entity"
	imService "github.com/coze-dev/coze-studio/backend/domain/im/service"
)

type Builder func() imService.PlatformAdapter

var (
	registryMu sync.RWMutex
	registry   = make(map[imEntity.Platform]Builder)
)

func Register(platform imEntity.Platform, builder Builder) {
	if builder == nil {
		panic(fmt.Sprintf("im adapter builder is nil: %s", platform))
	}

	registryMu.Lock()
	defer registryMu.Unlock()

	if _, exists := registry[platform]; exists {
		panic(fmt.Sprintf("im adapter already registered: %s", platform))
	}

	registry[platform] = builder
}

func BuildAll() []imService.PlatformAdapter {
	registryMu.RLock()
	builders := make(map[imEntity.Platform]Builder, len(registry))
	for platform, builder := range registry {
		builders[platform] = builder
	}
	registryMu.RUnlock()

	adapters := make([]imService.PlatformAdapter, 0, len(builders))
	for _, platform := range imEntity.AllPlatforms() {
		builder, ok := builders[platform]
		if !ok {
			continue
		}

		adapter := builder()
		if adapter == nil {
			continue
		}
		adapters = append(adapters, adapter)
	}

	return adapters
}
