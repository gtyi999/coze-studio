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

//go:generate mockgen -destination ../../../internal/mock/domain/crm/repository_mock.go --package crm -source repository.go
type Repository interface {
	DashboardRepository
	CustomerRepository
	ContactRepository
	OpportunityRepository
	FollowRecordRepository
	ProductRepository
	SalesOrderRepository
	QueryRepository
	SemanticRepository
	QueryLogRepository
	ForecastRepository
}

type DashboardRepository interface {
	GetDashboardOverview(ctx context.Context, scope *entity.Scope) (*entity.DashboardOverview, error)
}

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, customer *entity.Customer) (*entity.Customer, error)
	UpdateCustomer(ctx context.Context, customer *entity.Customer) (*entity.Customer, error)
	DeleteCustomer(ctx context.Context, scope *entity.Scope, customerID int64) error
	GetCustomer(ctx context.Context, scope *entity.Scope, customerID int64) (*entity.Customer, error)
	GetCustomerByID(ctx context.Context, customerID int64) (*entity.Customer, error)
	ListCustomers(ctx context.Context, filter *entity.CustomerFilter) ([]*entity.Customer, int64, error)
	CountActiveContactsByCustomer(ctx context.Context, scope *entity.Scope, customerID int64) (int64, error)
	CountActiveOpportunitiesByCustomer(ctx context.Context, scope *entity.Scope, customerID int64) (int64, error)
}

type ContactRepository interface {
	CreateContact(ctx context.Context, contact *entity.Contact) (*entity.Contact, error)
	UpdateContact(ctx context.Context, contact *entity.Contact) (*entity.Contact, error)
	DeleteContact(ctx context.Context, scope *entity.Scope, contactID int64) error
	GetContact(ctx context.Context, scope *entity.Scope, contactID int64) (*entity.Contact, error)
	GetContactByID(ctx context.Context, contactID int64) (*entity.Contact, error)
	ListContacts(ctx context.Context, filter *entity.ContactFilter) ([]*entity.Contact, int64, error)
}

type OpportunityRepository interface {
	CreateOpportunity(ctx context.Context, opportunity *entity.Opportunity) (*entity.Opportunity, error)
	UpdateOpportunity(ctx context.Context, opportunity *entity.Opportunity) (*entity.Opportunity, error)
	DeleteOpportunity(ctx context.Context, scope *entity.Scope, opportunityID int64) error
	GetOpportunity(ctx context.Context, scope *entity.Scope, opportunityID int64) (*entity.Opportunity, error)
	GetOpportunityByID(ctx context.Context, opportunityID int64) (*entity.Opportunity, error)
	ListOpportunities(ctx context.Context, filter *entity.OpportunityFilter) ([]*entity.Opportunity, int64, error)
}

type FollowRecordRepository interface {
	CreateFollowRecord(ctx context.Context, followRecord *entity.FollowRecord) (*entity.FollowRecord, error)
	UpdateFollowRecord(ctx context.Context, followRecord *entity.FollowRecord) (*entity.FollowRecord, error)
	DeleteFollowRecord(ctx context.Context, scope *entity.Scope, followRecordID int64) error
	GetFollowRecord(ctx context.Context, scope *entity.Scope, followRecordID int64) (*entity.FollowRecord, error)
	GetFollowRecordByID(ctx context.Context, followRecordID int64) (*entity.FollowRecord, error)
	ListFollowRecords(ctx context.Context, filter *entity.FollowRecordFilter) ([]*entity.FollowRecord, int64, error)
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
	UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error)
	DeleteProduct(ctx context.Context, scope *entity.Scope, productID int64) error
	GetProduct(ctx context.Context, scope *entity.Scope, productID int64) (*entity.Product, error)
	GetProductByID(ctx context.Context, productID int64) (*entity.Product, error)
	ListProducts(ctx context.Context, filter *entity.ProductFilter) ([]*entity.Product, int64, error)
}

type SalesOrderRepository interface {
	CreateSalesOrder(ctx context.Context, salesOrder *entity.SalesOrder) (*entity.SalesOrder, error)
	UpdateSalesOrder(ctx context.Context, salesOrder *entity.SalesOrder) (*entity.SalesOrder, error)
	DeleteSalesOrder(ctx context.Context, scope *entity.Scope, salesOrderID int64) error
	GetSalesOrder(ctx context.Context, scope *entity.Scope, salesOrderID int64) (*entity.SalesOrder, error)
	GetSalesOrderByID(ctx context.Context, salesOrderID int64) (*entity.SalesOrder, error)
	ListSalesOrders(ctx context.Context, filter *entity.SalesOrderFilter) ([]*entity.SalesOrder, int64, error)
}
