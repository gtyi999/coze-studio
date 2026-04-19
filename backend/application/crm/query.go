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

package crm

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/domain/crm/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

func (s *CRMApplicationService) RunNLQuery(ctx context.Context, req *entity.QueryRequest) (*entity.QueryResponse, error) {
	if req == nil {
		return nil, errorx.New(errno.ErrCRMQueryInvalidParamCode, errorx.KV("msg", "query request is required"))
	}

	userID, tenantID, err := s.authorizeScope(ctx, req.SpaceID)
	if err != nil {
		return nil, err
	}
	req.UserID = userID
	req.TenantID = tenantID

	return s.DomainSVC.RunNLQuery(ctx, req)
}

func (s *CRMApplicationService) GetSemanticCatalog(ctx context.Context, req *entity.SemanticCatalogRequest) (*entity.SemanticCatalog, error) {
	if req == nil {
		return nil, errorx.New(errno.ErrCRMQueryInvalidParamCode, errorx.KV("msg", "semantic catalog request is required"))
	}

	_, tenantID, err := s.authorizeScope(ctx, req.SpaceID)
	if err != nil {
		return nil, err
	}
	req.TenantID = tenantID

	return s.DomainSVC.GetSemanticCatalog(ctx, req)
}

func (s *CRMApplicationService) ListQueryLogs(ctx context.Context, filter *entity.QueryLogFilter) ([]*entity.QueryLogRecord, int64, error) {
	if filter == nil {
		return nil, 0, errorx.New(errno.ErrCRMQueryInvalidParamCode, errorx.KV("msg", "query log filter is required"))
	}

	_, tenantID, err := s.authorizeScope(ctx, filter.SpaceID)
	if err != nil {
		return nil, 0, err
	}
	filter.TenantID = tenantID

	return s.DomainSVC.ListQueryLogs(ctx, filter)
}

func (s *CRMApplicationService) GetForecastResult(ctx context.Context, req *entity.ForecastRequest) (*entity.ForecastResult, error) {
	if req == nil {
		return nil, errorx.New(errno.ErrCRMQueryInvalidParamCode, errorx.KV("msg", "forecast request is required"))
	}

	userID, tenantID, err := s.authorizeScope(ctx, req.SpaceID)
	if err != nil {
		return nil, err
	}
	req.UserID = userID
	req.TenantID = tenantID

	return s.DomainSVC.GetForecastResult(ctx, req)
}
