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
	"fmt"

	"github.com/coze-dev/coze-studio/backend/application/base/ctxutil"
	crossuser "github.com/coze-dev/coze-studio/backend/crossdomain/user"
	"github.com/coze-dev/coze-studio/backend/domain/crm/auditctx"
	"github.com/coze-dev/coze-studio/backend/domain/crm/entity"
	domainservice "github.com/coze-dev/coze-studio/backend/domain/crm/service"
	userentity "github.com/coze-dev/coze-studio/backend/domain/user/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

type CRMApplicationService struct {
	DomainSVC domainservice.CRM
}

func (s *CRMApplicationService) GetDashboardOverview(ctx context.Context, spaceID int64) (*entity.DashboardOverview, error) {
	scope, err := s.authorizeEntityScope(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.GetDashboardOverview(ctx, scope)
}

func (s *CRMApplicationService) ListCustomers(ctx context.Context, filter *entity.CustomerFilter) ([]*entity.Customer, int64, error) {
	if filter == nil {
		return nil, 0, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "customer filter is required"))
	}
	_, tenantID, err := s.authorizeScope(ctx, filter.SpaceID)
	if err != nil {
		return nil, 0, err
	}
	filter.TenantID = tenantID
	return s.DomainSVC.ListCustomers(ctx, filter)
}

func (s *CRMApplicationService) GetCustomer(ctx context.Context, spaceID int64, customerID int64) (*entity.Customer, error) {
	scope, err := s.authorizeEntityScope(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.GetCustomer(ctx, scope, customerID)
}

func (s *CRMApplicationService) CreateCustomer(ctx context.Context, customer *entity.Customer) (*entity.Customer, error) {
	actorCtx, err := s.prepareCreate(ctx, customer.SpaceID, &customer.TenantID, &customer.AuditInfo)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.CreateCustomer(actorCtx, customer)
}

func (s *CRMApplicationService) UpdateCustomer(ctx context.Context, customer *entity.Customer) (*entity.Customer, error) {
	actorCtx, err := s.prepareUpdate(ctx, customer.SpaceID, &customer.TenantID, &customer.AuditInfo)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.UpdateCustomer(actorCtx, customer)
}

func (s *CRMApplicationService) DeleteCustomer(ctx context.Context, spaceID int64, customerID int64) error {
	scope, actorCtx, err := s.authorizeActorScope(ctx, spaceID)
	if err != nil {
		return err
	}
	return s.DomainSVC.DeleteCustomer(actorCtx, scope, customerID)
}

func (s *CRMApplicationService) ListContacts(ctx context.Context, filter *entity.ContactFilter) ([]*entity.Contact, int64, error) {
	if filter == nil {
		return nil, 0, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "contact filter is required"))
	}
	_, tenantID, err := s.authorizeScope(ctx, filter.SpaceID)
	if err != nil {
		return nil, 0, err
	}
	filter.TenantID = tenantID
	return s.DomainSVC.ListContacts(ctx, filter)
}

func (s *CRMApplicationService) GetContact(ctx context.Context, spaceID int64, contactID int64) (*entity.Contact, error) {
	scope, err := s.authorizeEntityScope(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.GetContact(ctx, scope, contactID)
}

func (s *CRMApplicationService) CreateContact(ctx context.Context, contact *entity.Contact) (*entity.Contact, error) {
	actorCtx, err := s.prepareCreate(ctx, contact.SpaceID, &contact.TenantID, &contact.AuditInfo)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.CreateContact(actorCtx, contact)
}

func (s *CRMApplicationService) UpdateContact(ctx context.Context, contact *entity.Contact) (*entity.Contact, error) {
	actorCtx, err := s.prepareUpdate(ctx, contact.SpaceID, &contact.TenantID, &contact.AuditInfo)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.UpdateContact(actorCtx, contact)
}

func (s *CRMApplicationService) DeleteContact(ctx context.Context, spaceID int64, contactID int64) error {
	scope, actorCtx, err := s.authorizeActorScope(ctx, spaceID)
	if err != nil {
		return err
	}
	return s.DomainSVC.DeleteContact(actorCtx, scope, contactID)
}

func (s *CRMApplicationService) ListOpportunities(ctx context.Context, filter *entity.OpportunityFilter) ([]*entity.Opportunity, int64, error) {
	if filter == nil {
		return nil, 0, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "opportunity filter is required"))
	}
	_, tenantID, err := s.authorizeScope(ctx, filter.SpaceID)
	if err != nil {
		return nil, 0, err
	}
	filter.TenantID = tenantID
	return s.DomainSVC.ListOpportunities(ctx, filter)
}

func (s *CRMApplicationService) GetOpportunity(ctx context.Context, spaceID int64, opportunityID int64) (*entity.Opportunity, error) {
	scope, err := s.authorizeEntityScope(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.GetOpportunity(ctx, scope, opportunityID)
}

func (s *CRMApplicationService) CreateOpportunity(ctx context.Context, opportunity *entity.Opportunity) (*entity.Opportunity, error) {
	actorCtx, err := s.prepareCreate(ctx, opportunity.SpaceID, &opportunity.TenantID, &opportunity.AuditInfo)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.CreateOpportunity(actorCtx, opportunity)
}

func (s *CRMApplicationService) UpdateOpportunity(ctx context.Context, opportunity *entity.Opportunity) (*entity.Opportunity, error) {
	actorCtx, err := s.prepareUpdate(ctx, opportunity.SpaceID, &opportunity.TenantID, &opportunity.AuditInfo)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.UpdateOpportunity(actorCtx, opportunity)
}

func (s *CRMApplicationService) DeleteOpportunity(ctx context.Context, spaceID int64, opportunityID int64) error {
	scope, actorCtx, err := s.authorizeActorScope(ctx, spaceID)
	if err != nil {
		return err
	}
	return s.DomainSVC.DeleteOpportunity(actorCtx, scope, opportunityID)
}

func (s *CRMApplicationService) ListFollowRecords(ctx context.Context, filter *entity.FollowRecordFilter) ([]*entity.FollowRecord, int64, error) {
	if filter == nil {
		return nil, 0, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "follow_record filter is required"))
	}
	_, tenantID, err := s.authorizeScope(ctx, filter.SpaceID)
	if err != nil {
		return nil, 0, err
	}
	filter.TenantID = tenantID
	return s.DomainSVC.ListFollowRecords(ctx, filter)
}

func (s *CRMApplicationService) GetFollowRecord(ctx context.Context, spaceID int64, followRecordID int64) (*entity.FollowRecord, error) {
	scope, err := s.authorizeEntityScope(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.GetFollowRecord(ctx, scope, followRecordID)
}

func (s *CRMApplicationService) CreateFollowRecord(ctx context.Context, followRecord *entity.FollowRecord) (*entity.FollowRecord, error) {
	actorCtx, err := s.prepareCreate(ctx, followRecord.SpaceID, &followRecord.TenantID, &followRecord.AuditInfo)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.CreateFollowRecord(actorCtx, followRecord)
}

func (s *CRMApplicationService) UpdateFollowRecord(ctx context.Context, followRecord *entity.FollowRecord) (*entity.FollowRecord, error) {
	actorCtx, err := s.prepareUpdate(ctx, followRecord.SpaceID, &followRecord.TenantID, &followRecord.AuditInfo)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.UpdateFollowRecord(actorCtx, followRecord)
}

func (s *CRMApplicationService) DeleteFollowRecord(ctx context.Context, spaceID int64, followRecordID int64) error {
	scope, actorCtx, err := s.authorizeActorScope(ctx, spaceID)
	if err != nil {
		return err
	}
	return s.DomainSVC.DeleteFollowRecord(actorCtx, scope, followRecordID)
}

func (s *CRMApplicationService) ListProducts(ctx context.Context, filter *entity.ProductFilter) ([]*entity.Product, int64, error) {
	if filter == nil {
		return nil, 0, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "product filter is required"))
	}
	_, tenantID, err := s.authorizeScope(ctx, filter.SpaceID)
	if err != nil {
		return nil, 0, err
	}
	filter.TenantID = tenantID
	return s.DomainSVC.ListProducts(ctx, filter)
}

func (s *CRMApplicationService) GetProduct(ctx context.Context, spaceID int64, productID int64) (*entity.Product, error) {
	scope, err := s.authorizeEntityScope(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.GetProduct(ctx, scope, productID)
}

func (s *CRMApplicationService) CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	actorCtx, err := s.prepareCreate(ctx, product.SpaceID, &product.TenantID, &product.AuditInfo)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.CreateProduct(actorCtx, product)
}

func (s *CRMApplicationService) UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	actorCtx, err := s.prepareUpdate(ctx, product.SpaceID, &product.TenantID, &product.AuditInfo)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.UpdateProduct(actorCtx, product)
}

func (s *CRMApplicationService) DeleteProduct(ctx context.Context, spaceID int64, productID int64) error {
	scope, actorCtx, err := s.authorizeActorScope(ctx, spaceID)
	if err != nil {
		return err
	}
	return s.DomainSVC.DeleteProduct(actorCtx, scope, productID)
}

func (s *CRMApplicationService) ListSalesOrders(ctx context.Context, filter *entity.SalesOrderFilter) ([]*entity.SalesOrder, int64, error) {
	if filter == nil {
		return nil, 0, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "sales_order filter is required"))
	}
	_, tenantID, err := s.authorizeScope(ctx, filter.SpaceID)
	if err != nil {
		return nil, 0, err
	}
	filter.TenantID = tenantID
	return s.DomainSVC.ListSalesOrders(ctx, filter)
}

func (s *CRMApplicationService) GetSalesOrder(ctx context.Context, spaceID int64, salesOrderID int64) (*entity.SalesOrder, error) {
	scope, err := s.authorizeEntityScope(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.GetSalesOrder(ctx, scope, salesOrderID)
}

func (s *CRMApplicationService) CreateSalesOrder(ctx context.Context, salesOrder *entity.SalesOrder) (*entity.SalesOrder, error) {
	actorCtx, err := s.prepareCreate(ctx, salesOrder.SpaceID, &salesOrder.TenantID, &salesOrder.AuditInfo)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.CreateSalesOrder(actorCtx, salesOrder)
}

func (s *CRMApplicationService) UpdateSalesOrder(ctx context.Context, salesOrder *entity.SalesOrder) (*entity.SalesOrder, error) {
	actorCtx, err := s.prepareUpdate(ctx, salesOrder.SpaceID, &salesOrder.TenantID, &salesOrder.AuditInfo)
	if err != nil {
		return nil, err
	}
	return s.DomainSVC.UpdateSalesOrder(actorCtx, salesOrder)
}

func (s *CRMApplicationService) DeleteSalesOrder(ctx context.Context, spaceID int64, salesOrderID int64) error {
	scope, actorCtx, err := s.authorizeActorScope(ctx, spaceID)
	if err != nil {
		return err
	}
	return s.DomainSVC.DeleteSalesOrder(actorCtx, scope, salesOrderID)
}

func (s *CRMApplicationService) prepareCreate(ctx context.Context, spaceID int64, tenantID *int64, audit *entity.AuditInfo) (context.Context, error) {
	userID, resolvedTenantID, err := s.authorizeScope(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	*tenantID = resolvedTenantID
	audit.CreatedBy = userID
	audit.UpdatedBy = userID
	audit.CreatorID = userID
	audit.UpdaterID = userID
	audit.Normalize()
	return withActorContext(ctx, userID, resolvedTenantID, spaceID), nil
}

func (s *CRMApplicationService) prepareUpdate(ctx context.Context, spaceID int64, tenantID *int64, audit *entity.AuditInfo) (context.Context, error) {
	userID, resolvedTenantID, err := s.authorizeScope(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	*tenantID = resolvedTenantID
	audit.UpdatedBy = userID
	audit.UpdaterID = userID
	audit.Normalize()
	return withActorContext(ctx, userID, resolvedTenantID, spaceID), nil
}

func (s *CRMApplicationService) authorizeEntityScope(ctx context.Context, spaceID int64) (*entity.Scope, error) {
	_, tenantID, err := s.authorizeScope(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	return &entity.Scope{
		TenantID: tenantID,
		SpaceID:  spaceID,
	}, nil
}

func (s *CRMApplicationService) authorizeActorScope(ctx context.Context, spaceID int64) (*entity.Scope, context.Context, error) {
	userID, tenantID, err := s.authorizeScope(ctx, spaceID)
	if err != nil {
		return nil, nil, err
	}

	return &entity.Scope{
			TenantID: tenantID,
			SpaceID:  spaceID,
		},
		withActorContext(ctx, userID, tenantID, spaceID),
		nil
}

func (s *CRMApplicationService) authorizeScope(ctx context.Context, spaceID int64) (int64, int64, error) {
	if spaceID <= 0 {
		return 0, 0, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "space_id is required"))
	}

	userID := ctxutil.GetUIDFromCtx(ctx)
	if userID == nil {
		return 0, 0, errorx.New(errno.ErrCRMPermissionCode, errorx.KV("msg", "session required"))
	}
	space, err := loadAuthorizedSpace(ctx, *userID, spaceID)
	if err != nil {
		return 0, 0, errorx.New(errno.ErrCRMPermissionCode, errorx.KV("msg", err.Error()))
	}

	return *userID, resolveTenantID(space), nil
}

func resolveTenantID(space *userentity.Space) int64 {
	if space == nil {
		return 0
	}

	// The current open-source workspace model does not expose an independent
	// tenant object yet, so CRM uses the authorized workspace as the tenant
	// boundary and persists both values for future expansion.
	return space.ID
}

func loadAuthorizedSpace(ctx context.Context, uid int64, spaceID int64) (*userentity.Space, error) {
	spaces, err := crossuser.DefaultSVC().GetUserSpaceList(ctx, uid)
	if err != nil {
		return nil, err
	}

	for _, space := range spaces {
		if space.ID == spaceID {
			return space, nil
		}
	}

	return nil, fmt.Errorf("user %d does not have access to space %d", uid, spaceID)
}

func withActorContext(ctx context.Context, userID int64, tenantID int64, spaceID int64) context.Context {
	return auditctx.WithActor(ctx, &auditctx.Actor{
		UserID:   userID,
		TenantID: tenantID,
		SpaceID:  spaceID,
	})
}
