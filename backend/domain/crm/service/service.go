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

	"github.com/coze-dev/coze-studio/backend/domain/crm/entity"
	"github.com/coze-dev/coze-studio/backend/domain/crm/repository"
)

type Components struct {
	Repository repository.Repository
	Agent      AgentService
	Forecast   ForecastEngine
}

type AgentService interface {
	RewriteQuery(ctx context.Context, queryCtx *entity.QueryContext) (string, error)
	ClassifyIntent(ctx context.Context, queryCtx *entity.QueryContext, catalog *entity.SemanticCatalog) (entity.IntentType, []string, []string, []string, error)
	BuildSQLPlan(ctx context.Context, queryCtx *entity.QueryContext, catalog *entity.SemanticCatalog) (*entity.SQLPlan, error)
	FormatQueryResult(ctx context.Context, queryCtx *entity.QueryContext, plan *entity.SQLPlan, execution *entity.QueryExecution) (*entity.QueryResponse, error)
}

type QueryService interface {
	RunNLQuery(ctx context.Context, req *entity.QueryRequest) (*entity.QueryResponse, error)
	ListQueryLogs(ctx context.Context, filter *entity.QueryLogFilter) ([]*entity.QueryLogRecord, int64, error)
}

type SemanticMetadataService interface {
	GetSemanticCatalog(ctx context.Context, req *entity.SemanticCatalogRequest) (*entity.SemanticCatalog, error)
}

type ForecastService interface {
	GetForecastResult(ctx context.Context, req *entity.ForecastRequest) (*entity.ForecastResult, error)
}

type ForecastEngine interface {
	Analyze(ctx context.Context, req *entity.ForecastRequest, features []*entity.ForecastFeature) (*entity.ForecastResult, error)
}

type CRM interface {
	GetDashboardOverview(ctx context.Context, scope *entity.Scope) (*entity.DashboardOverview, error)
	QueryService
	SemanticMetadataService
	ForecastService

	CreateCustomer(ctx context.Context, customer *entity.Customer) (*entity.Customer, error)
	UpdateCustomer(ctx context.Context, customer *entity.Customer) (*entity.Customer, error)
	DeleteCustomer(ctx context.Context, scope *entity.Scope, customerID int64) error
	GetCustomer(ctx context.Context, scope *entity.Scope, customerID int64) (*entity.Customer, error)
	ListCustomers(ctx context.Context, filter *entity.CustomerFilter) ([]*entity.Customer, int64, error)

	CreateContact(ctx context.Context, contact *entity.Contact) (*entity.Contact, error)
	UpdateContact(ctx context.Context, contact *entity.Contact) (*entity.Contact, error)
	DeleteContact(ctx context.Context, scope *entity.Scope, contactID int64) error
	GetContact(ctx context.Context, scope *entity.Scope, contactID int64) (*entity.Contact, error)
	ListContacts(ctx context.Context, filter *entity.ContactFilter) ([]*entity.Contact, int64, error)

	CreateOpportunity(ctx context.Context, opportunity *entity.Opportunity) (*entity.Opportunity, error)
	UpdateOpportunity(ctx context.Context, opportunity *entity.Opportunity) (*entity.Opportunity, error)
	DeleteOpportunity(ctx context.Context, scope *entity.Scope, opportunityID int64) error
	GetOpportunity(ctx context.Context, scope *entity.Scope, opportunityID int64) (*entity.Opportunity, error)
	ListOpportunities(ctx context.Context, filter *entity.OpportunityFilter) ([]*entity.Opportunity, int64, error)

	CreateFollowRecord(ctx context.Context, followRecord *entity.FollowRecord) (*entity.FollowRecord, error)
	UpdateFollowRecord(ctx context.Context, followRecord *entity.FollowRecord) (*entity.FollowRecord, error)
	DeleteFollowRecord(ctx context.Context, scope *entity.Scope, followRecordID int64) error
	GetFollowRecord(ctx context.Context, scope *entity.Scope, followRecordID int64) (*entity.FollowRecord, error)
	ListFollowRecords(ctx context.Context, filter *entity.FollowRecordFilter) ([]*entity.FollowRecord, int64, error)

	CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
	UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
	DeleteProduct(ctx context.Context, scope *entity.Scope, productID int64) error
	GetProduct(ctx context.Context, scope *entity.Scope, productID int64) (*entity.Product, error)
	ListProducts(ctx context.Context, filter *entity.ProductFilter) ([]*entity.Product, int64, error)

	CreateSalesOrder(ctx context.Context, salesOrder *entity.SalesOrder) (*entity.SalesOrder, error)
	UpdateSalesOrder(ctx context.Context, salesOrder *entity.SalesOrder) (*entity.SalesOrder, error)
	DeleteSalesOrder(ctx context.Context, scope *entity.Scope, salesOrderID int64) error
	GetSalesOrder(ctx context.Context, scope *entity.Scope, salesOrderID int64) (*entity.SalesOrder, error)
	ListSalesOrders(ctx context.Context, filter *entity.SalesOrderFilter) ([]*entity.SalesOrder, int64, error)
}

func NewService(components *Components) CRM {
	if components == nil {
		components = &Components{}
	}
	if components.Agent == nil {
		components.Agent = noopAgentService{}
	}
	if components.Forecast == nil {
		components.Forecast = noopForecastEngine{}
	}
	return &crmService{
		components: components,
	}
}
