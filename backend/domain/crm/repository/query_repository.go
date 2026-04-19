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

	"github.com/coze-dev/coze-studio/backend/domain/crm/entity"
)

type QueryRepository interface {
	ExecuteQueryPlan(ctx context.Context, queryCtx *entity.QueryContext, plan *entity.SQLPlan) (*entity.QueryExecution, error)
	ExplainQueryPlan(ctx context.Context, queryCtx *entity.QueryContext, plan *entity.SQLPlan) (*entity.QueryCost, error)
}

type SemanticRepository interface {
	GetSemanticCatalog(ctx context.Context, scope *entity.Scope, req *entity.SemanticCatalogRequest) (*entity.SemanticCatalog, error)
}

type QueryLogRepository interface {
	AppendQueryLog(ctx context.Context, record *entity.QueryLogRecord) error
	ListQueryLogs(ctx context.Context, filter *entity.QueryLogFilter) ([]*entity.QueryLogRecord, int64, error)
}

type ForecastRepository interface {
	ListForecastFeatures(ctx context.Context, scope *entity.Scope, req *entity.ForecastRequest) ([]*entity.ForecastFeature, error)
}
