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

package coze

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	crmapp "github.com/coze-dev/coze-studio/backend/application/crm"
	"github.com/coze-dev/coze-studio/backend/domain/crm/entity"
)

type runCRMNLQueryRequest struct {
	SpaceID        string `query:"space_id" form:"space_id" json:"space_id"`
	ConversationID string `query:"conversation_id" form:"conversation_id" json:"conversation_id"`
	Question       string `query:"question" form:"question" json:"question"`
	Debug          bool   `query:"debug" form:"debug" json:"debug"`
	DryRun         bool   `query:"dry_run" form:"dry_run" json:"dry_run"`
}

type getCRMSemanticCatalogRequest struct {
	SpaceID string `query:"space_id" form:"space_id" json:"space_id"`
	Keyword string `query:"keyword" form:"keyword" json:"keyword"`
}

type listCRMQueryLogsRequest struct {
	SpaceID  string `query:"space_id" form:"space_id" json:"space_id"`
	UserID   string `query:"user_id" form:"user_id" json:"user_id"`
	Page     int    `query:"page" form:"page" json:"page"`
	PageSize int    `query:"page_size" form:"page_size" json:"page_size"`
	Question string `query:"question" form:"question" json:"question"`
}

type getCRMForecastResultRequest struct {
	SpaceID   string `query:"space_id" form:"space_id" json:"space_id"`
	MetricKey string `query:"metric_key" form:"metric_key" json:"metric_key"`
	Months    int    `query:"months" form:"months" json:"months"`
	Limit     int    `query:"limit" form:"limit" json:"limit"`
}

func RunCRMNLQuery(ctx context.Context, c *app.RequestContext) {
	var req runCRMNLQueryRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}

	resp, err := crmapp.CRMSVC.RunNLQuery(ctx, &entity.QueryRequest{
		Scope: entity.Scope{
			SpaceID: spaceID,
		},
		ConversationID: strings.TrimSpace(req.ConversationID),
		Question:       strings.TrimSpace(req.Question),
		Debug:          req.Debug,
		DryRun:         req.DryRun,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toCRMQueryResponseData(resp))
}

func GetCRMSemanticCatalog(ctx context.Context, c *app.RequestContext) {
	var req getCRMSemanticCatalogRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}

	resp, err := crmapp.CRMSVC.GetSemanticCatalog(ctx, &entity.SemanticCatalogRequest{
		Scope: entity.Scope{
			SpaceID: spaceID,
		},
		Keyword: strings.TrimSpace(req.Keyword),
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toCRMSemanticCatalogData(resp))
}

func ListCRMQueryLogs(ctx context.Context, c *app.RequestContext) {
	var req listCRMQueryLogsRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}

	list, total, err := crmapp.CRMSVC.ListQueryLogs(ctx, &entity.QueryLogFilter{
		Scope: entity.Scope{
			SpaceID: spaceID,
		},
		PageOption: entity.PageOption{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		UserID:   parseOptionalInt64Param(c, req.UserID, "user_id"),
		Question: strings.TrimSpace(req.Question),
	})
	if c.IsAborted() {
		return
	}
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, map[string]any{
		"list":  toCRMQueryLogListData(list),
		"total": total,
	})
}

func GetCRMForecastResult(ctx context.Context, c *app.RequestContext) {
	var req getCRMForecastResultRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}

	resp, err := crmapp.CRMSVC.GetForecastResult(ctx, &entity.ForecastRequest{
		Scope: entity.Scope{
			SpaceID: spaceID,
		},
		MetricKey: strings.TrimSpace(req.MetricKey),
		Months:    req.Months,
		Limit:     req.Limit,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toCRMForecastResultData(resp))
}

func toCRMQueryResponseData(resp *entity.QueryResponse) map[string]any {
	if resp == nil {
		return nil
	}

	var chart map[string]any
	if resp.Chart != nil {
		chart = map[string]any{
			"type":    string(resp.Chart.Type),
			"title":   resp.Chart.Title,
			"x_field": resp.Chart.XField,
			"y_field": resp.Chart.YField,
			"series":  resp.Chart.Series,
		}
	}

	var meta map[string]any
	if resp.Meta != nil {
		meta = map[string]any{
			"request_id":         resp.Meta.RequestID,
			"rewritten_question": resp.Meta.RewrittenQuestion,
			"final_sql":          resp.Meta.FinalSQL,
			"row_count":          resp.Meta.RowCount,
			"exec_ms":            resp.Meta.ExecMS,
		}
	}

	return map[string]any{
		"answer":      resp.Answer,
		"intent_type": string(resp.IntentType),
		"data":        resp.Data,
		"chart":       chart,
		"disclaimer":  resp.Disclaimer,
		"meta":        meta,
	}
}

func toCRMSemanticCatalogData(catalog *entity.SemanticCatalog) map[string]any {
	if catalog == nil {
		return nil
	}

	tables := make([]map[string]any, 0, len(catalog.Tables))
	for _, item := range catalog.Tables {
		if item == nil {
			continue
		}
		tables = append(tables, map[string]any{
			"id":                  formatCRMInt64(item.ID),
			"table_key":           item.TableKey,
			"table_name":          item.TableName,
			"table_desc":          item.TableDesc,
			"physical_table_name": item.PhysicalTableName,
			"primary_time_column": item.PrimaryTimeColumnKey,
			"status":              item.Status,
			"default_scope_json":  item.DefaultScopeJSON,
			"owner_domain":        item.OwnerDomain,
			"version_no":          item.VersionNo,
		})
	}

	metrics := make([]map[string]any, 0, len(catalog.Metrics))
	for _, item := range catalog.Metrics {
		if item == nil {
			continue
		}
		metrics = append(metrics, map[string]any{
			"id":                      formatCRMInt64(item.ID),
			"metric_key":              item.MetricKey,
			"metric_name":             item.MetricName,
			"table_key":               item.TableKey,
			"metric_type":             item.MetricType,
			"agg_func":                item.AggFunc,
			"measure_column_key":      item.MeasureColumnKey,
			"default_time_column_key": item.DefaultTimeColumnKey,
			"unit":                    item.Unit,
		})
	}

	return map[string]any{
		"tables":  tables,
		"metrics": metrics,
	}
}

func toCRMQueryLogListData(list []*entity.QueryLogRecord) []map[string]any {
	resp := make([]map[string]any, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}
		resp = append(resp, map[string]any{
			"id":                 formatCRMInt64(item.ID),
			"user_id":            formatCRMInt64(item.UserID),
			"request_id":         item.RequestID,
			"conversation_id":    item.ConversationID,
			"original_question":  item.OriginalQuestion,
			"rewritten_question": item.RewrittenQuestion,
			"intent_type":        string(item.IntentType),
			"status":             string(item.Status),
			"reject_reason":      item.RejectReason,
			"sql_hash":           item.SQLHash,
			"final_sql":          item.FinalSQL,
			"row_count":          item.RowCount,
			"exec_ms":            item.ExecMS,
			"created_at":         formatCRMInt64(item.CreatedAt),
		})
	}
	return resp
}

func toCRMForecastResultData(resp *entity.ForecastResult) map[string]any {
	if resp == nil {
		return nil
	}

	features := make([]map[string]any, 0, len(resp.Features))
	for _, item := range resp.Features {
		if item == nil {
			continue
		}
		features = append(features, map[string]any{
			"product_id":      formatCRMInt64(item.ProductID),
			"product_name":    item.ProductName,
			"metric_key":      item.MetricKey,
			"period":          item.Period,
			"metric_value":    formatCRMFloat(item.MetricValue),
			"growth_rate":     formatCRMFloat(item.GrowthRate),
			"trend_slope":     formatCRMFloat(item.TrendSlope),
			"weighted_avg_3m": formatCRMFloat(item.WeightedAvg3M),
			"volatility":      formatCRMFloat(item.Volatility),
			"score":           formatCRMFloat(item.Score),
		})
	}

	return map[string]any{
		"metric_key":       resp.MetricKey,
		"top_product_id":   formatCRMInt64(resp.TopProductID),
		"top_product_name": resp.TopProductName,
		"features":         features,
		"reasons":          resp.Reasons,
		"disclaimer":       resp.Disclaimer,
		"generated_at":     formatCRMInt64(resp.GeneratedAt),
	}
}
