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
	"strings"

	"github.com/coze-dev/coze-studio/backend/domain/crm/entity"
	"github.com/coze-dev/coze-studio/backend/domain/crm/valueobject"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

func validateScope(scope *entity.Scope) error {
	if scope == nil {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "scope is required"))
	}

	return valueobject.TenantScope{
		TenantID: scope.TenantID,
		SpaceID:  scope.SpaceID,
	}.Validate()
}

func validateCustomer(customer *entity.Customer) error {
	if customer == nil {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "customer is required"))
	}
	if err := validateScope(&entity.Scope{TenantID: customer.TenantID, SpaceID: customer.SpaceID}); err != nil {
		return err
	}

	customer.Normalize()
	customer.CustomerName = strings.TrimSpace(customer.CustomerName)
	if customer.CustomerName == "" {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "customer_name is required"))
	}
	customer.CustomerCode = strings.TrimSpace(customer.CustomerCode)
	customer.Industry = strings.TrimSpace(customer.Industry)
	customer.Level = strings.TrimSpace(customer.Level)
	customer.CustomerLevel = customer.Level
	customer.Mobile = strings.TrimSpace(customer.Mobile)
	customer.Email = strings.TrimSpace(customer.Email)
	customer.Address = strings.TrimSpace(customer.Address)
	customer.OwnerUserName = strings.TrimSpace(customer.OwnerUserName)
	customer.Remark = strings.TrimSpace(customer.Remark)
	customer.Description = customer.Remark
	customer.CustomerType = strings.TrimSpace(customer.CustomerType)
	return nil
}

func validateContact(contact *entity.Contact) error {
	if contact == nil {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "contact is required"))
	}
	if err := validateScope(&entity.Scope{TenantID: contact.TenantID, SpaceID: contact.SpaceID}); err != nil {
		return err
	}
	if contact.CustomerID <= 0 {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "customer_id is required"))
	}

	contact.Normalize()
	contact.ContactName = strings.TrimSpace(contact.ContactName)
	if contact.ContactName == "" {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "contact_name is required"))
	}
	contact.Mobile = strings.TrimSpace(contact.Mobile)
	contact.Email = strings.TrimSpace(contact.Email)
	contact.Title = strings.TrimSpace(contact.Title)
	contact.Position = contact.Title
	contact.Gender = strings.TrimSpace(contact.Gender)
	contact.Remark = strings.TrimSpace(contact.Remark)
	contact.Description = contact.Remark
	return nil
}

func validateOpportunity(opportunity *entity.Opportunity) error {
	if opportunity == nil {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "opportunity is required"))
	}
	if err := validateScope(&entity.Scope{TenantID: opportunity.TenantID, SpaceID: opportunity.SpaceID}); err != nil {
		return err
	}
	if opportunity.CustomerID <= 0 {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "customer_id is required"))
	}

	opportunity.Normalize()
	opportunity.OpportunityName = strings.TrimSpace(opportunity.OpportunityName)
	if opportunity.OpportunityName == "" {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "opportunity_name is required"))
	}
	if err := valueobject.EnsureNonNegativeFloat("amount", opportunity.Amount); err != nil {
		return err
	}
	opportunity.Stage = strings.TrimSpace(opportunity.Stage)
	opportunity.ExpectedCloseDate = strings.TrimSpace(opportunity.ExpectedCloseDate)
	opportunity.OwnerUserName = strings.TrimSpace(opportunity.OwnerUserName)
	opportunity.Remark = strings.TrimSpace(opportunity.Remark)
	opportunity.Description = opportunity.Remark
	return nil
}

func validateFollowRecord(followRecord *entity.FollowRecord) error {
	if followRecord == nil {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "follow_record is required"))
	}
	if err := validateScope(&entity.Scope{TenantID: followRecord.TenantID, SpaceID: followRecord.SpaceID}); err != nil {
		return err
	}
	if followRecord.CustomerID <= 0 {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "customer_id is required"))
	}

	followRecord.Normalize()
	followRecord.FollowType = strings.TrimSpace(followRecord.FollowType)
	if followRecord.FollowType == "" {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "follow_type is required"))
	}
	followRecord.Content = strings.TrimSpace(followRecord.Content)
	followRecord.FollowContent = followRecord.Content
	if followRecord.Content == "" {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "content is required"))
	}
	followRecord.OwnerUserName = strings.TrimSpace(followRecord.OwnerUserName)
	return nil
}

func validateProduct(product *entity.Product) error {
	if product == nil {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "product is required"))
	}
	if err := validateScope(&entity.Scope{TenantID: product.TenantID, SpaceID: product.SpaceID}); err != nil {
		return err
	}

	product.Normalize()
	product.ProductName = strings.TrimSpace(product.ProductName)
	if product.ProductName == "" {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "product_name is required"))
	}
	if err := valueobject.EnsureNonNegativeFloat("unit_price", product.UnitPrice); err != nil {
		return err
	}
	product.ProductCode = strings.TrimSpace(product.ProductCode)
	product.Category = strings.TrimSpace(product.Category)
	product.Unit = strings.TrimSpace(product.Unit)
	product.Remark = strings.TrimSpace(product.Remark)
	product.Description = product.Remark
	return nil
}

func validateSalesOrder(salesOrder *entity.SalesOrder) error {
	if salesOrder == nil {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "sales_order is required"))
	}
	if err := validateScope(&entity.Scope{TenantID: salesOrder.TenantID, SpaceID: salesOrder.SpaceID}); err != nil {
		return err
	}
	if salesOrder.CustomerID <= 0 {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "customer_id is required"))
	}
	if salesOrder.ProductID <= 0 {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "product_id is required"))
	}

	salesOrder.Normalize()
	if err := valueobject.EnsureNonNegativeFloat("quantity", salesOrder.Quantity); err != nil {
		return err
	}
	if err := valueobject.EnsureNonNegativeFloat("amount", salesOrder.Amount); err != nil {
		return err
	}
	salesOrder.ProductName = strings.TrimSpace(salesOrder.ProductName)
	salesOrder.SalesUserName = strings.TrimSpace(salesOrder.SalesUserName)
	salesOrder.OrderDate = strings.TrimSpace(salesOrder.OrderDate)
	salesOrder.Remark = strings.TrimSpace(salesOrder.Remark)
	salesOrder.Description = salesOrder.Remark
	return nil
}

func validateCustomerFilter(filter *entity.CustomerFilter) error {
	if filter == nil {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "customer filter is required"))
	}
	if err := validateScope(&filter.Scope); err != nil {
		return err
	}
	normalizePage(&filter.PageOption)
	filter.Keyword = strings.TrimSpace(filter.Keyword)
	return nil
}

func validateContactFilter(filter *entity.ContactFilter) error {
	if filter == nil {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "contact filter is required"))
	}
	if err := validateScope(&filter.Scope); err != nil {
		return err
	}
	normalizePage(&filter.PageOption)
	filter.Keyword = strings.TrimSpace(filter.Keyword)
	return nil
}

func validateOpportunityFilter(filter *entity.OpportunityFilter) error {
	if filter == nil {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "opportunity filter is required"))
	}
	if err := validateScope(&filter.Scope); err != nil {
		return err
	}
	normalizePage(&filter.PageOption)
	filter.Keyword = strings.TrimSpace(filter.Keyword)
	if filter.ExpectedCloseDateStart != nil {
		value := strings.TrimSpace(*filter.ExpectedCloseDateStart)
		filter.ExpectedCloseDateStart = &value
	}
	if filter.ExpectedCloseDateEnd != nil {
		value := strings.TrimSpace(*filter.ExpectedCloseDateEnd)
		filter.ExpectedCloseDateEnd = &value
	}
	return nil
}

func validateFollowRecordFilter(filter *entity.FollowRecordFilter) error {
	if filter == nil {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "follow_record filter is required"))
	}
	if err := validateScope(&filter.Scope); err != nil {
		return err
	}
	normalizePage(&filter.PageOption)
	filter.Keyword = strings.TrimSpace(filter.Keyword)
	return nil
}

func validateProductFilter(filter *entity.ProductFilter) error {
	if filter == nil {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "product filter is required"))
	}
	if err := validateScope(&filter.Scope); err != nil {
		return err
	}
	normalizePage(&filter.PageOption)
	filter.Keyword = strings.TrimSpace(filter.Keyword)
	return nil
}

func validateSalesOrderFilter(filter *entity.SalesOrderFilter) error {
	if filter == nil {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "sales_order filter is required"))
	}
	if err := validateScope(&filter.Scope); err != nil {
		return err
	}
	normalizePage(&filter.PageOption)
	filter.Keyword = strings.TrimSpace(filter.Keyword)
	if filter.OrderDateStart != nil {
		value := strings.TrimSpace(*filter.OrderDateStart)
		filter.OrderDateStart = &value
	}
	if filter.OrderDateEnd != nil {
		value := strings.TrimSpace(*filter.OrderDateEnd)
		filter.OrderDateEnd = &value
	}
	return nil
}

func normalizePage(option *entity.PageOption) {
	if option == nil {
		return
	}
	if option.Page <= 0 {
		option.Page = 1
	}
	if option.PageSize <= 0 {
		option.PageSize = 20
	}
	if option.PageSize > 100 {
		option.PageSize = 100
	}
}

func normalizeStatus(status string, defaultStatus string) string {
	status = strings.TrimSpace(strings.ToLower(status))
	if status == "" {
		return defaultStatus
	}
	return status
}

func loadCustomerForWrite(ctx context.Context, repo interface {
	GetCustomerByID(context.Context, int64) (*entity.Customer, error)
}, scope *entity.Scope, customerID int64) (*entity.Customer, error) {
	customer, err := repo.GetCustomerByID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	if err := (valueobject.TenantScope{TenantID: scope.TenantID, SpaceID: scope.SpaceID}).EnsureSame(customer.TenantID, customer.SpaceID, "customer"); err != nil {
		return nil, err
	}
	return customer, nil
}

func loadContactForWrite(ctx context.Context, repo interface {
	GetContactByID(context.Context, int64) (*entity.Contact, error)
}, scope *entity.Scope, contactID int64) (*entity.Contact, error) {
	contact, err := repo.GetContactByID(ctx, contactID)
	if err != nil {
		return nil, err
	}
	if err := (valueobject.TenantScope{TenantID: scope.TenantID, SpaceID: scope.SpaceID}).EnsureSame(contact.TenantID, contact.SpaceID, "contact"); err != nil {
		return nil, err
	}
	return contact, nil
}

func loadOpportunityForWrite(ctx context.Context, repo interface {
	GetOpportunityByID(context.Context, int64) (*entity.Opportunity, error)
}, scope *entity.Scope, opportunityID int64) (*entity.Opportunity, error) {
	opportunity, err := repo.GetOpportunityByID(ctx, opportunityID)
	if err != nil {
		return nil, err
	}
	if err := (valueobject.TenantScope{TenantID: scope.TenantID, SpaceID: scope.SpaceID}).EnsureSame(opportunity.TenantID, opportunity.SpaceID, "opportunity"); err != nil {
		return nil, err
	}
	return opportunity, nil
}

func loadProductForWrite(ctx context.Context, repo interface {
	GetProductByID(context.Context, int64) (*entity.Product, error)
}, scope *entity.Scope, productID int64) (*entity.Product, error) {
	product, err := repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if err := (valueobject.TenantScope{TenantID: scope.TenantID, SpaceID: scope.SpaceID}).EnsureSame(product.TenantID, product.SpaceID, "product"); err != nil {
		return nil, err
	}
	return product, nil
}

func loadFollowRecordForWrite(ctx context.Context, repo interface {
	GetFollowRecordByID(context.Context, int64) (*entity.FollowRecord, error)
}, scope *entity.Scope, followRecordID int64) (*entity.FollowRecord, error) {
	followRecord, err := repo.GetFollowRecordByID(ctx, followRecordID)
	if err != nil {
		return nil, err
	}
	if err := (valueobject.TenantScope{TenantID: scope.TenantID, SpaceID: scope.SpaceID}).EnsureSame(followRecord.TenantID, followRecord.SpaceID, "follow_record"); err != nil {
		return nil, err
	}
	return followRecord, nil
}

func loadSalesOrderForWrite(ctx context.Context, repo interface {
	GetSalesOrderByID(context.Context, int64) (*entity.SalesOrder, error)
}, scope *entity.Scope, salesOrderID int64) (*entity.SalesOrder, error) {
	salesOrder, err := repo.GetSalesOrderByID(ctx, salesOrderID)
	if err != nil {
		return nil, err
	}
	if err := (valueobject.TenantScope{TenantID: scope.TenantID, SpaceID: scope.SpaceID}).EnsureSame(salesOrder.TenantID, salesOrder.SpaceID, "sales_order"); err != nil {
		return nil, err
	}
	return salesOrder, nil
}
