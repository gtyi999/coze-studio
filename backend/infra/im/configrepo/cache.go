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

package configrepo

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	imEntity "github.com/coze-dev/coze-studio/backend/domain/im/entity"
	imRepo "github.com/coze-dev/coze-studio/backend/domain/im/repository"
	"github.com/coze-dev/coze-studio/backend/infra/cache"
	"github.com/coze-dev/coze-studio/backend/pkg/sonic"
)

const (
	configKeyPattern      = "im:channel_config:%s"
	configSpaceKeyPattern = "im:channel_config:space:%s"
	configCacheTTL        = 30 * 24 * time.Hour
)

type cacheRepository struct {
	cacheCli cache.Cmdable
}

func New(cacheCli cache.Cmdable) imRepo.ChannelConfigRepository {
	return &cacheRepository{
		cacheCli: cacheCli,
	}
}

func (r *cacheRepository) Save(ctx context.Context, cfg *imEntity.ChannelConfig) error {
	if cfg == nil || strings.TrimSpace(cfg.ConfigID) == "" {
		return errors.New("im channel config is invalid")
	}

	payload, err := sonic.MarshalString(cfg)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(configKeyPattern, cfg.ConfigID)
	if err = r.cacheCli.Set(ctx, key, payload, configCacheTTL).Err(); err != nil {
		return err
	}

	if strings.TrimSpace(cfg.SpaceID) == "" {
		return nil
	}

	return r.cacheCli.HSet(
		ctx,
		fmt.Sprintf(configSpaceKeyPattern, cfg.SpaceID),
		cfg.ConfigID,
		fmt.Sprintf("%d", cfg.UpdatedAtMS),
	).Err()
}

func (r *cacheRepository) Get(ctx context.Context, configID string) (*imEntity.ChannelConfig, bool, error) {
	if strings.TrimSpace(configID) == "" {
		return nil, false, nil
	}

	cmd := r.cacheCli.Get(ctx, fmt.Sprintf(configKeyPattern, configID))
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), cache.Nil) {
			return nil, false, nil
		}
		return nil, false, cmd.Err()
	}

	var cfg imEntity.ChannelConfig
	if err := sonic.UnmarshalString(cmd.Val(), &cfg); err != nil {
		return nil, false, err
	}

	return &cfg, true, nil
}

func (r *cacheRepository) ListBySpace(ctx context.Context, spaceID string) ([]*imEntity.ChannelConfig, error) {
	if strings.TrimSpace(spaceID) == "" {
		return nil, nil
	}

	index, err := r.cacheCli.HGetAll(ctx, fmt.Sprintf(configSpaceKeyPattern, spaceID)).Result()
	if err != nil {
		return nil, err
	}
	if len(index) == 0 {
		return nil, nil
	}

	configs := make([]*imEntity.ChannelConfig, 0, len(index))
	for configID := range index {
		cfg, found, getErr := r.Get(ctx, configID)
		if getErr != nil {
			return nil, getErr
		}
		if !found || cfg == nil {
			continue
		}
		configs = append(configs, cfg)
	}

	sort.Slice(configs, func(i, j int) bool {
		return configs[i].UpdatedAtMS > configs[j].UpdatedAtMS
	})

	return configs, nil
}
