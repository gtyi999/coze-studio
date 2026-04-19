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
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/crm/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

type noopAgentService struct{}

func (noopAgentService) RewriteQuery(_ context.Context, queryCtx *entity.QueryContext) (string, error) {
	if queryCtx == nil {
		return "", errorx.New(errno.ErrCRMQueryInvalidParamCode, errorx.KV("msg", "query context is required"))
	}
	return strings.TrimSpace(queryCtx.OriginalQuestion), nil
}

func (noopAgentService) ClassifyIntent(_ context.Context, _ *entity.QueryContext, _ *entity.SemanticCatalog) (entity.IntentType, []string, []string, []string, error) {
	return entity.IntentTypeUnknown, nil, nil, nil, errorx.New(errno.ErrCRMQueryFeaturePendingCode, errorx.KV("msg", "intent classifier is not implemented"))
}

func (noopAgentService) BuildSQLPlan(_ context.Context, _ *entity.QueryContext, _ *entity.SemanticCatalog) (*entity.SQLPlan, error) {
	return nil, errorx.New(errno.ErrCRMQueryFeaturePendingCode, errorx.KV("msg", "sql planner is not implemented"))
}

func (noopAgentService) FormatQueryResult(_ context.Context, queryCtx *entity.QueryContext, plan *entity.SQLPlan, execution *entity.QueryExecution) (*entity.QueryResponse, error) {
	if queryCtx == nil || plan == nil || execution == nil {
		return nil, errorx.New(errno.ErrCRMQueryInvalidParamCode, errorx.KV("msg", "query result formatter input is invalid"))
	}

	return &entity.QueryResponse{
		IntentType: queryCtx.IntentType,
		Data:       execution.Rows,
		SQLPlan:    plan,
		Meta: &entity.QueryMeta{
			RequestID:         queryCtx.RequestID,
			RewrittenQuestion: queryCtx.RewrittenQuestion,
			FinalSQL:          plan.FinalSQL,
			RowCount:          execution.RowCount,
			ExecMS:            execution.ExecMS,
		},
	}, nil
}

type noopForecastEngine struct{}

func (noopForecastEngine) Analyze(_ context.Context, _ *entity.ForecastRequest, _ []*entity.ForecastFeature) (*entity.ForecastResult, error) {
	return nil, errorx.New(errno.ErrCRMQueryFeaturePendingCode, errorx.KV("msg", "forecast engine is not implemented"))
}

func (s *crmService) RunNLQuery(ctx context.Context, req *entity.QueryRequest) (*entity.QueryResponse, error) {
	if err := validateQueryRequest(req); err != nil {
		return nil, err
	}

	queryCtx := &entity.QueryContext{
		Scope: entity.Scope{
			TenantID: req.TenantID,
			SpaceID:  req.SpaceID,
		},
		UserID:           req.UserID,
		RequestID:        strings.TrimSpace(req.RequestID),
		ConversationID:   strings.TrimSpace(req.ConversationID),
		OriginalQuestion: strings.TrimSpace(req.Question),
		Debug:            req.Debug,
		DryRun:           req.DryRun,
	}

	if queryCtx.RequestID == "" {
		queryCtx.RequestID = buildCRMQueryRequestID(req.UserID)
	}

	rewritten, err := s.components.Agent.RewriteQuery(ctx, queryCtx)
	if err != nil {
		return nil, err
	}
	queryCtx.RewrittenQuestion = rewritten

	catalog, err := s.components.Repository.GetSemanticCatalog(ctx, &queryCtx.Scope, &entity.SemanticCatalogRequest{
		Scope: queryCtx.Scope,
	})
	if err != nil {
		return nil, err
	}

	intentType, hitTables, hitColumns, hitMetrics, err := s.components.Agent.ClassifyIntent(ctx, queryCtx, catalog)
	if err != nil {
		return nil, err
	}
	queryCtx.IntentType = intentType
	queryCtx.HitTables = hitTables
	queryCtx.HitColumns = hitColumns
	queryCtx.HitMetrics = hitMetrics

	plan, err := s.components.Agent.BuildSQLPlan(ctx, queryCtx, catalog)
	if err != nil {
		return nil, err
	}

	cost, err := s.components.Repository.ExplainQueryPlan(ctx, queryCtx, plan)
	if err != nil {
		return nil, err
	}
	plan.Explain = cost

	execution, err := s.components.Repository.ExecuteQueryPlan(ctx, queryCtx, plan)
	if err != nil {
		return nil, err
	}

	resp, err := s.components.Agent.FormatQueryResult(ctx, queryCtx, plan, execution)
	if err != nil {
		return nil, err
	}

	if resp != nil && resp.Meta == nil {
		resp.Meta = &entity.QueryMeta{}
	}
	if resp != nil && resp.Meta != nil {
		resp.Meta.RequestID = queryCtx.RequestID
		resp.Meta.RewrittenQuestion = queryCtx.RewrittenQuestion
		resp.Meta.FinalSQL = plan.FinalSQL
		resp.Meta.RowCount = execution.RowCount
		resp.Meta.ExecMS = execution.ExecMS
	}

	_ = s.components.Repository.AppendQueryLog(ctx, &entity.QueryLogRecord{
		Scope: entity.Scope{
			TenantID: queryCtx.TenantID,
			SpaceID:  queryCtx.SpaceID,
		},
		UserID:            queryCtx.UserID,
		RequestID:         queryCtx.RequestID,
		ConversationID:    queryCtx.ConversationID,
		OriginalQuestion:  queryCtx.OriginalQuestion,
		RewrittenQuestion: queryCtx.RewrittenQuestion,
		IntentType:        queryCtx.IntentType,
		Status:            entity.QueryLogStatusSuccess,
		FinalSQL:          plan.FinalSQL,
		RowCount:          execution.RowCount,
		ExecMS:            execution.ExecMS,
		AuditInfo: entity.AuditInfo{
			CreatedBy: queryCtx.UserID,
			UpdatedBy: queryCtx.UserID,
			CreatedAt: time.Now().UnixMilli(),
			UpdatedAt: time.Now().UnixMilli(),
		},
	})

	return resp, nil
}

func (s *crmService) GetSemanticCatalog(ctx context.Context, req *entity.SemanticCatalogRequest) (*entity.SemanticCatalog, error) {
	if err := validateSemanticCatalogRequest(req); err != nil {
		return nil, err
	}

	return s.components.Repository.GetSemanticCatalog(ctx, &req.Scope, req)
}

func (s *crmService) ListQueryLogs(ctx context.Context, filter *entity.QueryLogFilter) ([]*entity.QueryLogRecord, int64, error) {
	if err := validateQueryLogFilter(filter); err != nil {
		return nil, 0, err
	}

	return s.components.Repository.ListQueryLogs(ctx, filter)
}

func (s *crmService) GetForecastResult(ctx context.Context, req *entity.ForecastRequest) (*entity.ForecastResult, error) {
	if err := validateForecastRequest(req); err != nil {
		return nil, err
	}

	features, err := s.components.Repository.ListForecastFeatures(ctx, &req.Scope, req)
	if err != nil {
		return nil, err
	}

	return s.components.Forecast.Analyze(ctx, req, features)
}

func buildCRMQueryRequestID(userID int64) string {
	return strings.Join([]string{
		"crmq",
		time.Now().UTC().Format("20060102150405"),
		strconv.FormatInt(userID, 10),
	}, "-")
}
