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
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

type crmService struct {
	components *Components
}

func (s *crmService) GetDashboardOverview(ctx context.Context, scope *entity.Scope) (*entity.DashboardOverview, error) {
	if err := validateScope(scope); err != nil {
		return nil, err
	}

	return s.components.Repository.GetDashboardOverview(ctx, scope)
}

func (s *crmService) CreateCustomer(ctx context.Context, customer *entity.Customer) (*entity.Customer, error) {
	if err := validateCustomer(customer); err != nil {
		return nil, err
	}

	customer.Status = normalizeStatus(customer.Status, entity.StatusActive)
	if customer.OwnerUserID == 0 {
		customer.OwnerUserID = customer.CreatedBy
	}

	return s.components.Repository.CreateCustomer(ctx, customer)
}

func (s *crmService) UpdateCustomer(ctx context.Context, customer *entity.Customer) (*entity.Customer, error) {
	if customer == nil || customer.CustomerID <= 0 {
		return nil, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "customer_id is required"))
	}
	if err := validateCustomer(customer); err != nil {
		return nil, err
	}

	current, err := loadCustomerForWrite(ctx, s.components.Repository, &entity.Scope{
		TenantID: customer.TenantID,
		SpaceID:  customer.SpaceID,
	}, customer.CustomerID)
	if err != nil {
		return nil, err
	}

	customer.Status = normalizeStatus(customer.Status, fallbackString(current.Status, entity.StatusActive))
	if customer.OwnerUserID == 0 {
		customer.OwnerUserID = current.OwnerUserID
	}
	if customer.OwnerUserName == "" {
		customer.OwnerUserName = current.OwnerUserName
	}
	inheritAudit(&customer.AuditInfo, current.AuditInfo)

	return s.components.Repository.UpdateCustomer(ctx, customer)
}

func (s *crmService) DeleteCustomer(ctx context.Context, scope *entity.Scope, customerID int64) error {
	if err := validateScope(scope); err != nil {
		return err
	}
	if customerID <= 0 {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "customer_id is required"))
	}

	if _, err := loadCustomerForWrite(ctx, s.components.Repository, scope, customerID); err != nil {
		return err
	}

	contactCount, err := s.components.Repository.CountActiveContactsByCustomer(ctx, scope, customerID)
	if err != nil {
		return err
	}
	if contactCount > 0 {
		return errorx.New(errno.ErrCRMOperateCode, errorx.KV("msg", "customer has undeleted contacts"))
	}

	opportunityCount, err := s.components.Repository.CountActiveOpportunitiesByCustomer(ctx, scope, customerID)
	if err != nil {
		return err
	}
	if opportunityCount > 0 {
		return errorx.New(errno.ErrCRMOperateCode, errorx.KV("msg", "customer has undeleted opportunities"))
	}

	followHasRecords, err := hasAnyFollowRecords(ctx, s.components.Repository, &entity.FollowRecordFilter{
		Scope: *scope,
		PageOption: entity.PageOption{
			Page:     1,
			PageSize: 1,
		},
		CustomerID: &customerID,
	})
	if err != nil {
		return err
	}
	if followHasRecords {
		return errorx.New(errno.ErrCRMOperateCode, errorx.KV("msg", "customer has undeleted follow records"))
	}

	salesOrderHasRecords, err := hasAnySalesOrders(ctx, s.components.Repository, &entity.SalesOrderFilter{
		Scope: *scope,
		PageOption: entity.PageOption{
			Page:     1,
			PageSize: 1,
		},
		CustomerID: &customerID,
	})
	if err != nil {
		return err
	}
	if salesOrderHasRecords {
		return errorx.New(errno.ErrCRMOperateCode, errorx.KV("msg", "customer has undeleted sales orders"))
	}

	return s.components.Repository.DeleteCustomer(ctx, scope, customerID)
}

func (s *crmService) GetCustomer(ctx context.Context, scope *entity.Scope, customerID int64) (*entity.Customer, error) {
	if err := validateScope(scope); err != nil {
		return nil, err
	}
	if customerID <= 0 {
		return nil, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "customer_id is required"))
	}

	return loadCustomerForWrite(ctx, s.components.Repository, scope, customerID)
}

func (s *crmService) ListCustomers(ctx context.Context, filter *entity.CustomerFilter) ([]*entity.Customer, int64, error) {
	if err := validateCustomerFilter(filter); err != nil {
		return nil, 0, err
	}

	return s.components.Repository.ListCustomers(ctx, filter)
}

func (s *crmService) CreateContact(ctx context.Context, contact *entity.Contact) (*entity.Contact, error) {
	if err := validateContact(contact); err != nil {
		return nil, err
	}

	if _, err := loadCustomerForWrite(ctx, s.components.Repository, &entity.Scope{
		TenantID: contact.TenantID,
		SpaceID:  contact.SpaceID,
	}, contact.CustomerID); err != nil {
		return nil, err
	}

	contact.Status = normalizeStatus(contact.Status, entity.StatusActive)
	return s.components.Repository.CreateContact(ctx, contact)
}

func (s *crmService) UpdateContact(ctx context.Context, contact *entity.Contact) (*entity.Contact, error) {
	if contact == nil || contact.ContactID <= 0 {
		return nil, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "contact_id is required"))
	}
	if err := validateContact(contact); err != nil {
		return nil, err
	}

	current, err := loadContactForWrite(ctx, s.components.Repository, &entity.Scope{
		TenantID: contact.TenantID,
		SpaceID:  contact.SpaceID,
	}, contact.ContactID)
	if err != nil {
		return nil, err
	}
	if _, err := loadCustomerForWrite(ctx, s.components.Repository, &entity.Scope{
		TenantID: contact.TenantID,
		SpaceID:  contact.SpaceID,
	}, contact.CustomerID); err != nil {
		return nil, err
	}

	contact.Status = normalizeStatus(contact.Status, fallbackString(current.Status, entity.StatusActive))
	inheritAudit(&contact.AuditInfo, current.AuditInfo)

	return s.components.Repository.UpdateContact(ctx, contact)
}

func (s *crmService) DeleteContact(ctx context.Context, scope *entity.Scope, contactID int64) error {
	if err := validateScope(scope); err != nil {
		return err
	}
	if contactID <= 0 {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "contact_id is required"))
	}

	if _, err := loadContactForWrite(ctx, s.components.Repository, scope, contactID); err != nil {
		return err
	}

	followHasRecords, err := hasAnyFollowRecords(ctx, s.components.Repository, &entity.FollowRecordFilter{
		Scope: *scope,
		PageOption: entity.PageOption{
			Page:     1,
			PageSize: 1,
		},
		ContactID: &contactID,
	})
	if err != nil {
		return err
	}
	if followHasRecords {
		return errorx.New(errno.ErrCRMOperateCode, errorx.KV("msg", "contact has undeleted follow records"))
	}

	return s.components.Repository.DeleteContact(ctx, scope, contactID)
}

func (s *crmService) GetContact(ctx context.Context, scope *entity.Scope, contactID int64) (*entity.Contact, error) {
	if err := validateScope(scope); err != nil {
		return nil, err
	}
	if contactID <= 0 {
		return nil, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "contact_id is required"))
	}

	return loadContactForWrite(ctx, s.components.Repository, scope, contactID)
}

func (s *crmService) ListContacts(ctx context.Context, filter *entity.ContactFilter) ([]*entity.Contact, int64, error) {
	if err := validateContactFilter(filter); err != nil {
		return nil, 0, err
	}

	return s.components.Repository.ListContacts(ctx, filter)
}

func (s *crmService) CreateOpportunity(ctx context.Context, opportunity *entity.Opportunity) (*entity.Opportunity, error) {
	if err := validateOpportunity(opportunity); err != nil {
		return nil, err
	}

	scope := &entity.Scope{TenantID: opportunity.TenantID, SpaceID: opportunity.SpaceID}
	if _, err := loadCustomerForWrite(ctx, s.components.Repository, scope, opportunity.CustomerID); err != nil {
		return nil, err
	}
	if opportunity.ContactID > 0 {
		if _, err := loadContactForWrite(ctx, s.components.Repository, scope, opportunity.ContactID); err != nil {
			return nil, err
		}
	}

	opportunity.Stage = fallbackString(opportunity.Stage, entity.StageInitial)
	opportunity.Status = normalizeStatus(opportunity.Status, entity.StatusOpen)
	if opportunity.OwnerUserID == 0 {
		opportunity.OwnerUserID = opportunity.CreatedBy
	}

	return s.components.Repository.CreateOpportunity(ctx, opportunity)
}

func (s *crmService) UpdateOpportunity(ctx context.Context, opportunity *entity.Opportunity) (*entity.Opportunity, error) {
	if opportunity == nil || opportunity.OpportunityID <= 0 {
		return nil, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "opportunity_id is required"))
	}
	if err := validateOpportunity(opportunity); err != nil {
		return nil, err
	}

	scope := &entity.Scope{TenantID: opportunity.TenantID, SpaceID: opportunity.SpaceID}
	current, err := loadOpportunityForWrite(ctx, s.components.Repository, scope, opportunity.OpportunityID)
	if err != nil {
		return nil, err
	}
	if _, err := loadCustomerForWrite(ctx, s.components.Repository, scope, opportunity.CustomerID); err != nil {
		return nil, err
	}
	if opportunity.ContactID > 0 {
		if _, err := loadContactForWrite(ctx, s.components.Repository, scope, opportunity.ContactID); err != nil {
			return nil, err
		}
	}

	opportunity.Stage = fallbackString(opportunity.Stage, fallbackString(current.Stage, entity.StageInitial))
	opportunity.Status = normalizeStatus(opportunity.Status, fallbackString(current.Status, entity.StatusOpen))
	if opportunity.OwnerUserID == 0 {
		opportunity.OwnerUserID = current.OwnerUserID
	}
	if opportunity.OwnerUserName == "" {
		opportunity.OwnerUserName = current.OwnerUserName
	}
	inheritAudit(&opportunity.AuditInfo, current.AuditInfo)

	return s.components.Repository.UpdateOpportunity(ctx, opportunity)
}

func (s *crmService) DeleteOpportunity(ctx context.Context, scope *entity.Scope, opportunityID int64) error {
	if err := validateScope(scope); err != nil {
		return err
	}
	if opportunityID <= 0 {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "opportunity_id is required"))
	}

	if _, err := loadOpportunityForWrite(ctx, s.components.Repository, scope, opportunityID); err != nil {
		return err
	}

	followHasRecords, err := hasAnyFollowRecords(ctx, s.components.Repository, &entity.FollowRecordFilter{
		Scope: *scope,
		PageOption: entity.PageOption{
			Page:     1,
			PageSize: 1,
		},
		OpportunityID: &opportunityID,
	})
	if err != nil {
		return err
	}
	if followHasRecords {
		return errorx.New(errno.ErrCRMOperateCode, errorx.KV("msg", "opportunity has undeleted follow records"))
	}

	salesOrderHasRecords, err := hasAnySalesOrders(ctx, s.components.Repository, &entity.SalesOrderFilter{
		Scope: *scope,
		PageOption: entity.PageOption{
			Page:     1,
			PageSize: 1,
		},
		OpportunityID: &opportunityID,
	})
	if err != nil {
		return err
	}
	if salesOrderHasRecords {
		return errorx.New(errno.ErrCRMOperateCode, errorx.KV("msg", "opportunity has undeleted sales orders"))
	}

	return s.components.Repository.DeleteOpportunity(ctx, scope, opportunityID)
}

func (s *crmService) GetOpportunity(ctx context.Context, scope *entity.Scope, opportunityID int64) (*entity.Opportunity, error) {
	if err := validateScope(scope); err != nil {
		return nil, err
	}
	if opportunityID <= 0 {
		return nil, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "opportunity_id is required"))
	}

	return loadOpportunityForWrite(ctx, s.components.Repository, scope, opportunityID)
}

func (s *crmService) ListOpportunities(ctx context.Context, filter *entity.OpportunityFilter) ([]*entity.Opportunity, int64, error) {
	if err := validateOpportunityFilter(filter); err != nil {
		return nil, 0, err
	}

	return s.components.Repository.ListOpportunities(ctx, filter)
}

func (s *crmService) CreateFollowRecord(ctx context.Context, followRecord *entity.FollowRecord) (*entity.FollowRecord, error) {
	if err := validateFollowRecord(followRecord); err != nil {
		return nil, err
	}

	scope := &entity.Scope{TenantID: followRecord.TenantID, SpaceID: followRecord.SpaceID}
	if _, err := loadCustomerForWrite(ctx, s.components.Repository, scope, followRecord.CustomerID); err != nil {
		return nil, err
	}
	if followRecord.ContactID > 0 {
		if _, err := loadContactForWrite(ctx, s.components.Repository, scope, followRecord.ContactID); err != nil {
			return nil, err
		}
	}
	if followRecord.OpportunityID > 0 {
		if _, err := loadOpportunityForWrite(ctx, s.components.Repository, scope, followRecord.OpportunityID); err != nil {
			return nil, err
		}
	}

	followRecord.Status = normalizeStatus(followRecord.Status, entity.StatusActive)
	if followRecord.OwnerUserID == 0 {
		followRecord.OwnerUserID = followRecord.CreatedBy
	}

	return s.components.Repository.CreateFollowRecord(ctx, followRecord)
}

func (s *crmService) UpdateFollowRecord(ctx context.Context, followRecord *entity.FollowRecord) (*entity.FollowRecord, error) {
	if followRecord == nil || followRecord.FollowRecordID <= 0 {
		return nil, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "follow_record_id is required"))
	}
	if err := validateFollowRecord(followRecord); err != nil {
		return nil, err
	}

	scope := &entity.Scope{TenantID: followRecord.TenantID, SpaceID: followRecord.SpaceID}
	current, err := loadFollowRecordForWrite(ctx, s.components.Repository, scope, followRecord.FollowRecordID)
	if err != nil {
		return nil, err
	}
	if _, err := loadCustomerForWrite(ctx, s.components.Repository, scope, followRecord.CustomerID); err != nil {
		return nil, err
	}
	if followRecord.ContactID > 0 {
		if _, err := loadContactForWrite(ctx, s.components.Repository, scope, followRecord.ContactID); err != nil {
			return nil, err
		}
	}
	if followRecord.OpportunityID > 0 {
		if _, err := loadOpportunityForWrite(ctx, s.components.Repository, scope, followRecord.OpportunityID); err != nil {
			return nil, err
		}
	}

	followRecord.Status = normalizeStatus(followRecord.Status, fallbackString(current.Status, entity.StatusActive))
	if followRecord.OwnerUserID == 0 {
		followRecord.OwnerUserID = current.OwnerUserID
	}
	if followRecord.OwnerUserName == "" {
		followRecord.OwnerUserName = current.OwnerUserName
	}
	inheritAudit(&followRecord.AuditInfo, current.AuditInfo)

	return s.components.Repository.UpdateFollowRecord(ctx, followRecord)
}

func (s *crmService) DeleteFollowRecord(ctx context.Context, scope *entity.Scope, followRecordID int64) error {
	if err := validateScope(scope); err != nil {
		return err
	}
	if followRecordID <= 0 {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "follow_record_id is required"))
	}

	if _, err := loadFollowRecordForWrite(ctx, s.components.Repository, scope, followRecordID); err != nil {
		return err
	}

	return s.components.Repository.DeleteFollowRecord(ctx, scope, followRecordID)
}

func (s *crmService) GetFollowRecord(ctx context.Context, scope *entity.Scope, followRecordID int64) (*entity.FollowRecord, error) {
	if err := validateScope(scope); err != nil {
		return nil, err
	}
	if followRecordID <= 0 {
		return nil, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "follow_record_id is required"))
	}

	return loadFollowRecordForWrite(ctx, s.components.Repository, scope, followRecordID)
}

func (s *crmService) ListFollowRecords(ctx context.Context, filter *entity.FollowRecordFilter) ([]*entity.FollowRecord, int64, error) {
	if err := validateFollowRecordFilter(filter); err != nil {
		return nil, 0, err
	}

	return s.components.Repository.ListFollowRecords(ctx, filter)
}

func (s *crmService) CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	if err := validateProduct(product); err != nil {
		return nil, err
	}

	product.Status = normalizeStatus(product.Status, entity.StatusActive)
	return s.components.Repository.CreateProduct(ctx, product)
}

func (s *crmService) UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	if product == nil || product.ProductID <= 0 {
		return nil, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "product_id is required"))
	}
	if err := validateProduct(product); err != nil {
		return nil, err
	}

	current, err := loadProductForWrite(ctx, s.components.Repository, &entity.Scope{
		TenantID: product.TenantID,
		SpaceID:  product.SpaceID,
	}, product.ProductID)
	if err != nil {
		return nil, err
	}

	product.Status = normalizeStatus(product.Status, fallbackString(current.Status, entity.StatusActive))
	inheritAudit(&product.AuditInfo, current.AuditInfo)

	return s.components.Repository.UpdateProduct(ctx, product)
}

func (s *crmService) DeleteProduct(ctx context.Context, scope *entity.Scope, productID int64) error {
	if err := validateScope(scope); err != nil {
		return err
	}
	if productID <= 0 {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "product_id is required"))
	}

	if _, err := loadProductForWrite(ctx, s.components.Repository, scope, productID); err != nil {
		return err
	}

	salesOrderHasRecords, err := hasAnySalesOrders(ctx, s.components.Repository, &entity.SalesOrderFilter{
		Scope: *scope,
		PageOption: entity.PageOption{
			Page:     1,
			PageSize: 1,
		},
		ProductID: &productID,
	})
	if err != nil {
		return err
	}
	if salesOrderHasRecords {
		return errorx.New(errno.ErrCRMOperateCode, errorx.KV("msg", "product has undeleted sales orders"))
	}

	return s.components.Repository.DeleteProduct(ctx, scope, productID)
}

func (s *crmService) GetProduct(ctx context.Context, scope *entity.Scope, productID int64) (*entity.Product, error) {
	if err := validateScope(scope); err != nil {
		return nil, err
	}
	if productID <= 0 {
		return nil, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "product_id is required"))
	}

	return loadProductForWrite(ctx, s.components.Repository, scope, productID)
}

func (s *crmService) ListProducts(ctx context.Context, filter *entity.ProductFilter) ([]*entity.Product, int64, error) {
	if err := validateProductFilter(filter); err != nil {
		return nil, 0, err
	}

	return s.components.Repository.ListProducts(ctx, filter)
}

func (s *crmService) CreateSalesOrder(ctx context.Context, salesOrder *entity.SalesOrder) (*entity.SalesOrder, error) {
	if err := validateSalesOrder(salesOrder); err != nil {
		return nil, err
	}

	scope := &entity.Scope{TenantID: salesOrder.TenantID, SpaceID: salesOrder.SpaceID}
	if _, err := loadCustomerForWrite(ctx, s.components.Repository, scope, salesOrder.CustomerID); err != nil {
		return nil, err
	}
	product, err := loadProductForWrite(ctx, s.components.Repository, scope, salesOrder.ProductID)
	if err != nil {
		return nil, err
	}
	if salesOrder.OpportunityID > 0 {
		if _, err := loadOpportunityForWrite(ctx, s.components.Repository, scope, salesOrder.OpportunityID); err != nil {
			return nil, err
		}
	}

	salesOrder.Status = normalizeStatus(salesOrder.Status, entity.StatusDraft)
	if salesOrder.SalesUserID == 0 {
		salesOrder.SalesUserID = salesOrder.CreatedBy
	}
	if salesOrder.ProductName == "" {
		salesOrder.ProductName = product.ProductName
		salesOrder.ProductSummary = product.ProductName
	}

	return s.components.Repository.CreateSalesOrder(ctx, salesOrder)
}

func (s *crmService) UpdateSalesOrder(ctx context.Context, salesOrder *entity.SalesOrder) (*entity.SalesOrder, error) {
	if salesOrder == nil || salesOrder.SalesOrderID <= 0 {
		return nil, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "sales_order_id is required"))
	}
	if err := validateSalesOrder(salesOrder); err != nil {
		return nil, err
	}

	scope := &entity.Scope{TenantID: salesOrder.TenantID, SpaceID: salesOrder.SpaceID}
	current, err := loadSalesOrderForWrite(ctx, s.components.Repository, scope, salesOrder.SalesOrderID)
	if err != nil {
		return nil, err
	}
	if _, err := loadCustomerForWrite(ctx, s.components.Repository, scope, salesOrder.CustomerID); err != nil {
		return nil, err
	}
	product, err := loadProductForWrite(ctx, s.components.Repository, scope, salesOrder.ProductID)
	if err != nil {
		return nil, err
	}
	if salesOrder.OpportunityID > 0 {
		if _, err := loadOpportunityForWrite(ctx, s.components.Repository, scope, salesOrder.OpportunityID); err != nil {
			return nil, err
		}
	}

	salesOrder.Status = normalizeStatus(salesOrder.Status, fallbackString(current.Status, entity.StatusDraft))
	if salesOrder.SalesUserID == 0 {
		salesOrder.SalesUserID = current.SalesUserID
	}
	if salesOrder.SalesUserName == "" {
		salesOrder.SalesUserName = current.SalesUserName
	}
	if salesOrder.ProductName == "" {
		salesOrder.ProductName = fallbackString(current.ProductName, product.ProductName)
		salesOrder.ProductSummary = salesOrder.ProductName
	}
	inheritAudit(&salesOrder.AuditInfo, current.AuditInfo)

	return s.components.Repository.UpdateSalesOrder(ctx, salesOrder)
}

func (s *crmService) DeleteSalesOrder(ctx context.Context, scope *entity.Scope, salesOrderID int64) error {
	if err := validateScope(scope); err != nil {
		return err
	}
	if salesOrderID <= 0 {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "sales_order_id is required"))
	}

	if _, err := loadSalesOrderForWrite(ctx, s.components.Repository, scope, salesOrderID); err != nil {
		return err
	}

	return s.components.Repository.DeleteSalesOrder(ctx, scope, salesOrderID)
}

func (s *crmService) GetSalesOrder(ctx context.Context, scope *entity.Scope, salesOrderID int64) (*entity.SalesOrder, error) {
	if err := validateScope(scope); err != nil {
		return nil, err
	}
	if salesOrderID <= 0 {
		return nil, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "sales_order_id is required"))
	}

	return loadSalesOrderForWrite(ctx, s.components.Repository, scope, salesOrderID)
}

func (s *crmService) ListSalesOrders(ctx context.Context, filter *entity.SalesOrderFilter) ([]*entity.SalesOrder, int64, error) {
	if err := validateSalesOrderFilter(filter); err != nil {
		return nil, 0, err
	}

	return s.components.Repository.ListSalesOrders(ctx, filter)
}

func hasAnyFollowRecords(ctx context.Context, repo repository.Repository, filter *entity.FollowRecordFilter) (bool, error) {
	if filter == nil {
		return false, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "follow_record filter is required"))
	}

	filter.Page = 1
	filter.PageSize = 1
	_, total, err := repo.ListFollowRecords(ctx, filter)
	if err != nil {
		return false, err
	}
	return total > 0, nil
}

func hasAnySalesOrders(ctx context.Context, repo repository.Repository, filter *entity.SalesOrderFilter) (bool, error) {
	if filter == nil {
		return false, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "sales_order filter is required"))
	}

	filter.Page = 1
	filter.PageSize = 1
	_, total, err := repo.ListSalesOrders(ctx, filter)
	if err != nil {
		return false, err
	}
	return total > 0, nil
}

func fallbackString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func inheritAudit(target *entity.AuditInfo, current entity.AuditInfo) {
	if target == nil {
		return
	}

	current.Normalize()
	target.Normalize()
	target.CreatedBy = current.CreatedBy
	target.CreatorID = current.CreatorID
	target.CreatedAt = current.CreatedAt
	target.IsDeleted = current.IsDeleted
	if target.UpdatedBy == 0 {
		target.UpdatedBy = current.UpdatedBy
	}
	if target.UpdaterID == 0 {
		target.UpdaterID = current.UpdaterID
	}
	target.Normalize()
}
