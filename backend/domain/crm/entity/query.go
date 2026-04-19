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

package entity

type IntentType string
type QueryPlanType string
type ChartType string
type QueryLogStatus string
type DataScopeType string

const (
	IntentTypeUnknown         IntentType = "unknown"
	IntentTypeCustomerCount   IntentType = "customer_count"
	IntentTypeTopSales        IntentType = "top_sales_current_quarter"
	IntentTypeSalesTopN       IntentType = "sales_topn_current_quarter"
	IntentTypeProductSalesTop IntentType = "product_sales_topn"
	IntentTypeForecastProduct IntentType = "forecast_hot_product"
)

const (
	QueryPlanTypeUnknown     QueryPlanType = "unknown"
	QueryPlanTypeAggregation QueryPlanType = "aggregation"
	QueryPlanTypeRanking     QueryPlanType = "ranking"
	QueryPlanTypeTrend       QueryPlanType = "trend"
	QueryPlanTypeForecast    QueryPlanType = "forecast"
)

const (
	ChartTypeUnknown ChartType = "unknown"
	ChartTypeStat    ChartType = "stat"
	ChartTypeBar     ChartType = "bar"
	ChartTypeLine    ChartType = "line"
	ChartTypeTable   ChartType = "table"
)

const (
	QueryLogStatusPending QueryLogStatus = "pending"
	QueryLogStatusSuccess QueryLogStatus = "success"
	QueryLogStatusEmpty   QueryLogStatus = "empty"
	QueryLogStatusReject  QueryLogStatus = "reject"
	QueryLogStatusFailed  QueryLogStatus = "failed"
)

const (
	DataScopeTypeSelf       DataScopeType = "self"
	DataScopeTypeDepartment DataScopeType = "department"
	DataScopeTypeTenant     DataScopeType = "tenant"
)

type QueryRequest struct {
	Scope

	UserID         int64
	RequestID      string
	ConversationID string
	Question       string
	Debug          bool
	DryRun         bool
}

type QueryContext struct {
	Scope

	UserID            int64
	RequestID         string
	ConversationID    string
	OriginalQuestion  string
	RewrittenQuestion string
	IntentType        IntentType
	HitTables         []string
	HitColumns        []string
	HitMetrics        []string
	DataScope         *DataScope
	Debug             bool
	DryRun            bool
}

type DataScope struct {
	ScopeType     DataScopeType
	UserIDs       []int64
	MaskSensitive bool
}

type SQLPlan struct {
	PlanType   QueryPlanType
	IntentType IntentType
	MetricKey  string
	Dimensions []string
	Filters    []*PlanFilter
	GroupBy    []string
	OrderBy    []*OrderBy
	TimeRange  *TimeRange
	Limit      int64
	FinalSQL   string
	SQLArgs    []any
	ChartHint  ChartType
	Explain    *QueryCost
}

type PlanFilter struct {
	Field string
	Op    string
	Value any
}

type OrderBy struct {
	Field     string
	Direction string
}

type TimeRange struct {
	Label     string
	StartDate string
	EndDate   string
}

type QueryCost struct {
	EstimatedRows int64
	HitIndex      bool
	RiskLevel     string
	ExplainRows   []string
}

type QueryExecution struct {
	Columns  []string
	Rows     []map[string]any
	RowCount int64
	ExecMS   int64
}

type ChartConfig struct {
	Type   ChartType
	Title  string
	XField string
	YField string
	Series []string
}

type QueryMeta struct {
	RequestID         string
	RewrittenQuestion string
	FinalSQL          string
	RowCount          int64
	ExecMS            int64
}

type QueryResponse struct {
	Answer     string
	IntentType IntentType
	Data       []map[string]any
	Chart      *ChartConfig
	Disclaimer string
	SQLPlan    *SQLPlan
	Forecast   *ForecastResult
	Meta       *QueryMeta
}

type QueryLogFilter struct {
	Scope
	PageOption

	UserID     *int64
	IntentType *IntentType
	Status     *QueryLogStatus
	Question   string
}

type QueryLogRecord struct {
	ID int64
	Scope

	UserID            int64
	RequestID         string
	ConversationID    string
	OriginalQuestion  string
	RewrittenQuestion string
	IntentType        IntentType
	Status            QueryLogStatus
	RejectReason      string
	SQLHash           string
	FinalSQL          string
	RowCount          int64
	ExecMS            int64

	AuditInfo
}
