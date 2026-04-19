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

package mysql

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/domain/crm/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

func (r *crmRepository) GetSemanticCatalog(ctx context.Context, scope *entity.Scope, req *entity.SemanticCatalogRequest) (*entity.SemanticCatalog, error) {
	_ = ctx
	_ = scope
	_ = req
	return nil, errorx.New(errno.ErrCRMQueryFeaturePendingCode, errorx.KV("msg", "semantic catalog repository is not implemented"))
}

func (r *crmRepository) ExecuteQueryPlan(ctx context.Context, queryCtx *entity.QueryContext, plan *entity.SQLPlan) (*entity.QueryExecution, error) {
	_ = ctx
	_ = queryCtx
	_ = plan
	return nil, errorx.New(errno.ErrCRMQueryFeaturePendingCode, errorx.KV("msg", "query executor repository is not implemented"))
}

func (r *crmRepository) ExplainQueryPlan(ctx context.Context, queryCtx *entity.QueryContext, plan *entity.SQLPlan) (*entity.QueryCost, error) {
	_ = ctx
	_ = queryCtx
	_ = plan
	return nil, errorx.New(errno.ErrCRMQueryFeaturePendingCode, errorx.KV("msg", "query explain repository is not implemented"))
}

func (r *crmRepository) AppendQueryLog(ctx context.Context, record *entity.QueryLogRecord) error {
	_ = ctx
	_ = record
	return errorx.New(errno.ErrCRMQueryFeaturePendingCode, errorx.KV("msg", "query audit log repository is not implemented"))
}

func (r *crmRepository) ListQueryLogs(ctx context.Context, filter *entity.QueryLogFilter) ([]*entity.QueryLogRecord, int64, error) {
	_ = ctx
	_ = filter
	return nil, 0, errorx.New(errno.ErrCRMQueryFeaturePendingCode, errorx.KV("msg", "query log list repository is not implemented"))
}

func (r *crmRepository) ListForecastFeatures(ctx context.Context, scope *entity.Scope, req *entity.ForecastRequest) ([]*entity.ForecastFeature, error) {
	_ = ctx
	_ = scope
	_ = req
	return nil, errorx.New(errno.ErrCRMQueryFeaturePendingCode, errorx.KV("msg", "forecast feature repository is not implemented"))
}
