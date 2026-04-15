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
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	crmapp "github.com/coze-dev/coze-studio/backend/application/crm"
	"github.com/coze-dev/coze-studio/backend/crossdomain/user"
	"github.com/coze-dev/coze-studio/backend/domain/crm/auditctx"
	crmentity "github.com/coze-dev/coze-studio/backend/domain/crm/entity"
	crmdomainservice "github.com/coze-dev/coze-studio/backend/domain/crm/service"
	userentity "github.com/coze-dev/coze-studio/backend/domain/user/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

type crmTestEnvelope struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

type crmCustomerListResp struct {
	List []struct {
		CustomerID string `json:"customer_id"`
		Status     string `json:"status"`
	} `json:"list"`
	Total int64 `json:"total"`
}

type crmOpportunityResp struct {
	OpportunityID string `json:"opportunity_id"`
	Amount        string `json:"amount"`
}

type crmDashboardResp struct {
	CustomerTotal             int64  `json:"customer_total"`
	NewCustomersThisMonth     int64  `json:"new_customers_this_month"`
	OpportunityTotalAmount    string `json:"opportunity_total_amount"`
	NewOpportunitiesThisMonth int64  `json:"new_opportunities_this_month"`
	SalesOrderTotalAmount     string `json:"sales_order_total_amount"`
	RecentOrderTrend          []struct {
		Date        string `json:"date"`
		OrderCount  int64  `json:"order_count"`
		OrderAmount string `json:"order_amount"`
	} `json:"recent_order_trend"`
}

type fakeCrossUserService struct {
	spaces []*userentity.Space
}

func (f *fakeCrossUserService) GetUserSpaceList(ctx context.Context, userID int64) ([]*userentity.Space, error) {
	return f.spaces, nil
}

func (f *fakeCrossUserService) GetUserSpaceBySpaceID(ctx context.Context, spaceID []int64) ([]*userentity.Space, error) {
	return f.spaces, nil
}

type fakeCRMDomainService struct {
	createCustomerFn    func(context.Context, *crmentity.Customer) (*crmentity.Customer, error)
	deleteCustomerFn    func(context.Context, *crmentity.Scope, int64) error
	listCustomersFn     func(context.Context, *crmentity.CustomerFilter) ([]*crmentity.Customer, int64, error)
	createOpportunityFn func(context.Context, *crmentity.Opportunity) (*crmentity.Opportunity, error)
	getDashboardFn      func(context.Context, *crmentity.Scope) (*crmentity.DashboardOverview, error)
}

var _ crmdomainservice.CRM = (*fakeCRMDomainService)(nil)

func (f *fakeCRMDomainService) CreateCustomer(ctx context.Context, customer *crmentity.Customer) (*crmentity.Customer, error) {
	if f.createCustomerFn != nil {
		return f.createCustomerFn(ctx, customer)
	}
	return customer, nil
}

func (f *fakeCRMDomainService) GetDashboardOverview(ctx context.Context, scope *crmentity.Scope) (*crmentity.DashboardOverview, error) {
	if f.getDashboardFn != nil {
		return f.getDashboardFn(ctx, scope)
	}
	return &crmentity.DashboardOverview{}, nil
}

func (f *fakeCRMDomainService) UpdateCustomer(ctx context.Context, customer *crmentity.Customer) (*crmentity.Customer, error) {
	return customer, nil
}

func (f *fakeCRMDomainService) DeleteCustomer(ctx context.Context, scope *crmentity.Scope, customerID int64) error {
	if f.deleteCustomerFn != nil {
		return f.deleteCustomerFn(ctx, scope, customerID)
	}
	return nil
}

func (f *fakeCRMDomainService) GetCustomer(ctx context.Context, scope *crmentity.Scope, customerID int64) (*crmentity.Customer, error) {
	return &crmentity.Customer{CustomerID: customerID, TenantID: scope.TenantID, SpaceID: scope.SpaceID}, nil
}

func (f *fakeCRMDomainService) ListCustomers(ctx context.Context, filter *crmentity.CustomerFilter) ([]*crmentity.Customer, int64, error) {
	if f.listCustomersFn != nil {
		return f.listCustomersFn(ctx, filter)
	}
	return nil, 0, nil
}

func (f *fakeCRMDomainService) CreateContact(ctx context.Context, contact *crmentity.Contact) (*crmentity.Contact, error) {
	return contact, nil
}

func (f *fakeCRMDomainService) UpdateContact(ctx context.Context, contact *crmentity.Contact) (*crmentity.Contact, error) {
	return contact, nil
}

func (f *fakeCRMDomainService) DeleteContact(ctx context.Context, scope *crmentity.Scope, contactID int64) error {
	return nil
}

func (f *fakeCRMDomainService) GetContact(ctx context.Context, scope *crmentity.Scope, contactID int64) (*crmentity.Contact, error) {
	return &crmentity.Contact{ContactID: contactID, TenantID: scope.TenantID, SpaceID: scope.SpaceID}, nil
}

func (f *fakeCRMDomainService) ListContacts(ctx context.Context, filter *crmentity.ContactFilter) ([]*crmentity.Contact, int64, error) {
	return nil, 0, nil
}

func (f *fakeCRMDomainService) CreateOpportunity(ctx context.Context, opportunity *crmentity.Opportunity) (*crmentity.Opportunity, error) {
	if f.createOpportunityFn != nil {
		return f.createOpportunityFn(ctx, opportunity)
	}
	return opportunity, nil
}

func (f *fakeCRMDomainService) UpdateOpportunity(ctx context.Context, opportunity *crmentity.Opportunity) (*crmentity.Opportunity, error) {
	return opportunity, nil
}

func (f *fakeCRMDomainService) DeleteOpportunity(ctx context.Context, scope *crmentity.Scope, opportunityID int64) error {
	return nil
}

func (f *fakeCRMDomainService) GetOpportunity(ctx context.Context, scope *crmentity.Scope, opportunityID int64) (*crmentity.Opportunity, error) {
	return &crmentity.Opportunity{OpportunityID: opportunityID, TenantID: scope.TenantID, SpaceID: scope.SpaceID}, nil
}

func (f *fakeCRMDomainService) ListOpportunities(ctx context.Context, filter *crmentity.OpportunityFilter) ([]*crmentity.Opportunity, int64, error) {
	return nil, 0, nil
}

func (f *fakeCRMDomainService) CreateFollowRecord(ctx context.Context, followRecord *crmentity.FollowRecord) (*crmentity.FollowRecord, error) {
	return followRecord, nil
}

func (f *fakeCRMDomainService) UpdateFollowRecord(ctx context.Context, followRecord *crmentity.FollowRecord) (*crmentity.FollowRecord, error) {
	return followRecord, nil
}

func (f *fakeCRMDomainService) DeleteFollowRecord(ctx context.Context, scope *crmentity.Scope, followRecordID int64) error {
	return nil
}

func (f *fakeCRMDomainService) GetFollowRecord(ctx context.Context, scope *crmentity.Scope, followRecordID int64) (*crmentity.FollowRecord, error) {
	return &crmentity.FollowRecord{FollowRecordID: followRecordID, TenantID: scope.TenantID, SpaceID: scope.SpaceID}, nil
}

func (f *fakeCRMDomainService) ListFollowRecords(ctx context.Context, filter *crmentity.FollowRecordFilter) ([]*crmentity.FollowRecord, int64, error) {
	return nil, 0, nil
}

func (f *fakeCRMDomainService) CreateProduct(ctx context.Context, product *crmentity.Product) (*crmentity.Product, error) {
	return product, nil
}

func (f *fakeCRMDomainService) UpdateProduct(ctx context.Context, product *crmentity.Product) (*crmentity.Product, error) {
	return product, nil
}

func (f *fakeCRMDomainService) DeleteProduct(ctx context.Context, scope *crmentity.Scope, productID int64) error {
	return nil
}

func (f *fakeCRMDomainService) GetProduct(ctx context.Context, scope *crmentity.Scope, productID int64) (*crmentity.Product, error) {
	return &crmentity.Product{ProductID: productID, TenantID: scope.TenantID, SpaceID: scope.SpaceID}, nil
}

func (f *fakeCRMDomainService) ListProducts(ctx context.Context, filter *crmentity.ProductFilter) ([]*crmentity.Product, int64, error) {
	return nil, 0, nil
}

func (f *fakeCRMDomainService) CreateSalesOrder(ctx context.Context, salesOrder *crmentity.SalesOrder) (*crmentity.SalesOrder, error) {
	return salesOrder, nil
}

func (f *fakeCRMDomainService) UpdateSalesOrder(ctx context.Context, salesOrder *crmentity.SalesOrder) (*crmentity.SalesOrder, error) {
	return salesOrder, nil
}

func (f *fakeCRMDomainService) DeleteSalesOrder(ctx context.Context, scope *crmentity.Scope, salesOrderID int64) error {
	return nil
}

func (f *fakeCRMDomainService) GetSalesOrder(ctx context.Context, scope *crmentity.Scope, salesOrderID int64) (*crmentity.SalesOrder, error) {
	return &crmentity.SalesOrder{SalesOrderID: salesOrderID, TenantID: scope.TenantID, SpaceID: scope.SpaceID}, nil
}

func (f *fakeCRMDomainService) ListSalesOrders(ctx context.Context, filter *crmentity.SalesOrderFilter) ([]*crmentity.SalesOrder, int64, error) {
	return nil, 0, nil
}

func newCRMTestServer(t *testing.T, domainSVC crmdomainservice.CRM) *server.Hertz {
	t.Helper()

	oldApp := crmapp.CRMSVC
	oldCross := user.DefaultSVC()
	t.Cleanup(func() {
		crmapp.CRMSVC = oldApp
		user.SetDefaultSVC(oldCross)
	})

	crmapp.CRMSVC = &crmapp.CRMApplicationService{
		DomainSVC: domainSVC,
	}
	user.SetDefaultSVC(&fakeCrossUserService{
		spaces: []*userentity.Space{{ID: 1, Name: "default"}},
	})

	h := server.Default()
	h.Use(func(c context.Context, ctx *app.RequestContext) {
		c = ctxcache.Init(c)
		ctxcache.Store(c, consts.SessionDataKeyInCtx, &userentity.Session{UserID: 123})
		ctx.Next(c)
	})
	RegisterManualCRMRoutes(h)
	return h
}

func TestCRMListCustomers(t *testing.T) {
	var captured *crmentity.CustomerFilter
	h := newCRMTestServer(t, &fakeCRMDomainService{
		listCustomersFn: func(ctx context.Context, filter *crmentity.CustomerFilter) ([]*crmentity.Customer, int64, error) {
			captured = filter
			return []*crmentity.Customer{
				{
					CustomerID:    101,
					TenantID:      filter.TenantID,
					SpaceID:       filter.SpaceID,
					CustomerName:  "Acme",
					CustomerCode:  "CUST-001",
					Status:        "active",
					OwnerUserID:   88,
					OwnerUserName: "owner",
				},
			}, 1, nil
		},
	})

	w := ut.PerformRequest(
		h.Engine,
		http.MethodGet,
		"/api/crm/customer/list?space_id=1&page=2&page_size=10&keyword=Acme&status=active&owner_user_id=88&created_at_start=1710000000000&created_at_end=1710000009999",
		nil,
	)

	require.Equal(t, http.StatusOK, w.Code)
	require.NotNil(t, captured)
	assert.Equal(t, int64(1), captured.TenantID)
	assert.Equal(t, int64(1), captured.SpaceID)
	assert.Equal(t, 2, captured.Page)
	assert.Equal(t, 10, captured.PageSize)
	assert.Equal(t, "Acme", captured.Keyword)
	require.NotNil(t, captured.Status)
	assert.Equal(t, "active", *captured.Status)
	require.NotNil(t, captured.OwnerUserID)
	assert.Equal(t, int64(88), *captured.OwnerUserID)
	require.NotNil(t, captured.CreatedAtStart)
	assert.Equal(t, int64(1710000000000), *captured.CreatedAtStart)
	require.NotNil(t, captured.CreatedAtEnd)
	assert.Equal(t, int64(1710000009999), *captured.CreatedAtEnd)

	var envelope crmTestEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &envelope))
	assert.Equal(t, 0, envelope.Code)

	var resp crmCustomerListResp
	require.NoError(t, json.Unmarshal(envelope.Data, &resp))
	assert.Equal(t, int64(1), resp.Total)
	require.Len(t, resp.List, 1)
	assert.Equal(t, "101", resp.List[0].CustomerID)
	assert.Equal(t, "active", resp.List[0].Status)
}

func TestCRMCreateCustomerInjectsTenantAndAuditFields(t *testing.T) {
	var captured *crmentity.Customer
	h := newCRMTestServer(t, &fakeCRMDomainService{
		createCustomerFn: func(ctx context.Context, customer *crmentity.Customer) (*crmentity.Customer, error) {
			captured = customer
			customer.CustomerID = 1001
			return customer, nil
		},
	})

	payload := []byte(`{"space_id":"1","customer_name":"Acme","owner_user_name":"Alice","status":"active"}`)
	w := ut.PerformRequest(
		h.Engine,
		http.MethodPost,
		"/api/crm/customer/create",
		&ut.Body{Body: bytes.NewBuffer(payload), Len: len(payload)},
		ut.Header{Key: "Content-Type", Value: "application/json"},
	)

	require.Equal(t, http.StatusOK, w.Code)
	require.NotNil(t, captured)
	assert.Equal(t, int64(1), captured.TenantID)
	assert.Equal(t, int64(1), captured.SpaceID)
	assert.Equal(t, int64(123), captured.CreatedBy)
	assert.Equal(t, int64(123), captured.UpdatedBy)
	assert.Equal(t, int64(123), captured.CreatorID)
	assert.Equal(t, int64(123), captured.UpdaterID)
}

func TestCRMGetDashboardOverview(t *testing.T) {
	var captured *crmentity.Scope
	h := newCRMTestServer(t, &fakeCRMDomainService{
		getDashboardFn: func(ctx context.Context, scope *crmentity.Scope) (*crmentity.DashboardOverview, error) {
			captured = scope
			return &crmentity.DashboardOverview{
				CustomerTotal:             36,
				NewCustomersThisMonth:     4,
				OpportunityTotalAmount:    128000.5,
				NewOpportunitiesThisMonth: 7,
				SalesOrderTotalAmount:     86500.25,
				RecentOrderTrend: []*crmentity.DashboardOrderTrendPoint{
					{
						Date:        "2026-04-10",
						OrderCount:  2,
						OrderAmount: 16800,
					},
					{
						Date:        "2026-04-11",
						OrderCount:  0,
						OrderAmount: 0,
					},
				},
			}, nil
		},
	})

	w := ut.PerformRequest(
		h.Engine,
		http.MethodGet,
		"/api/crm/dashboard/overview?space_id=1",
		nil,
	)

	require.Equal(t, http.StatusOK, w.Code)
	require.NotNil(t, captured)
	assert.Equal(t, int64(1), captured.TenantID)
	assert.Equal(t, int64(1), captured.SpaceID)

	var envelope crmTestEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &envelope))
	assert.Equal(t, 0, envelope.Code)

	var resp crmDashboardResp
	require.NoError(t, json.Unmarshal(envelope.Data, &resp))
	assert.Equal(t, int64(36), resp.CustomerTotal)
	assert.Equal(t, int64(4), resp.NewCustomersThisMonth)
	assert.Equal(t, "128000.5", resp.OpportunityTotalAmount)
	assert.Equal(t, int64(7), resp.NewOpportunitiesThisMonth)
	assert.Equal(t, "86500.25", resp.SalesOrderTotalAmount)
	require.Len(t, resp.RecentOrderTrend, 2)
	assert.Equal(t, "2026-04-10", resp.RecentOrderTrend[0].Date)
	assert.Equal(t, int64(2), resp.RecentOrderTrend[0].OrderCount)
	assert.Equal(t, "16800", resp.RecentOrderTrend[0].OrderAmount)
}

func TestCRMDeleteCustomerInjectsAuditActor(t *testing.T) {
	var capturedScope *crmentity.Scope
	var capturedActor *auditctx.Actor
	h := newCRMTestServer(t, &fakeCRMDomainService{
		deleteCustomerFn: func(ctx context.Context, scope *crmentity.Scope, customerID int64) error {
			capturedScope = scope
			actor, ok := auditctx.ActorFromContext(ctx)
			require.True(t, ok)
			capturedActor = actor
			assert.Equal(t, int64(2001), customerID)
			return nil
		},
	})

	payload := []byte(`{"space_id":"1","customer_id":"2001"}`)
	w := ut.PerformRequest(
		h.Engine,
		http.MethodPost,
		"/api/crm/customer/delete",
		&ut.Body{Body: bytes.NewBuffer(payload), Len: len(payload)},
		ut.Header{Key: "Content-Type", Value: "application/json"},
	)

	require.Equal(t, http.StatusOK, w.Code)
	require.NotNil(t, capturedScope)
	require.NotNil(t, capturedActor)
	assert.Equal(t, int64(1), capturedScope.TenantID)
	assert.Equal(t, int64(1), capturedScope.SpaceID)
	assert.Equal(t, int64(123), capturedActor.UserID)
	assert.Equal(t, int64(1), capturedActor.TenantID)
	assert.Equal(t, int64(1), capturedActor.SpaceID)
}

func TestCRMCreateOpportunity(t *testing.T) {
	var captured *crmentity.Opportunity
	h := newCRMTestServer(t, &fakeCRMDomainService{
		createOpportunityFn: func(ctx context.Context, opportunity *crmentity.Opportunity) (*crmentity.Opportunity, error) {
			captured = opportunity
			return &crmentity.Opportunity{
				OpportunityID:     2001,
				TenantID:          opportunity.TenantID,
				SpaceID:           opportunity.SpaceID,
				CustomerID:        opportunity.CustomerID,
				OpportunityName:   opportunity.OpportunityName,
				Amount:            opportunity.Amount,
				ExpectedCloseDate: opportunity.ExpectedCloseDate,
				OwnerUserID:       opportunity.OwnerUserID,
				OwnerUserName:     opportunity.OwnerUserName,
				Status:            opportunity.Status,
				Remark:            opportunity.Remark,
				AuditInfo: crmentity.AuditInfo{
					CreatedBy: opportunity.CreatedBy,
					UpdatedBy: opportunity.UpdatedBy,
				},
			}, nil
		},
	})

	payload := []byte(`{"space_id":"1","customer_id":"66","opportunity_name":"Big Deal","stage":"proposal","amount":"1200.5","expected_close_date":"2026-06-30","owner_user_id":"88","owner_user_name":"alice","status":"open","remark":"important"}`)
	w := ut.PerformRequest(
		h.Engine,
		http.MethodPost,
		"/api/crm/opportunity/create",
		&ut.Body{Body: bytes.NewBuffer(payload), Len: len(payload)},
		ut.Header{Key: "Content-Type", Value: "application/json"},
	)

	require.Equal(t, http.StatusOK, w.Code)
	require.NotNil(t, captured)
	assert.Equal(t, int64(1), captured.TenantID)
	assert.Equal(t, int64(1), captured.SpaceID)
	assert.Equal(t, int64(66), captured.CustomerID)
	assert.Equal(t, 1200.5, captured.Amount)
	assert.Equal(t, int64(88), captured.OwnerUserID)
	assert.Equal(t, int64(123), captured.CreatedBy)
	assert.Equal(t, int64(123), captured.UpdatedBy)

	var envelope crmTestEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &envelope))
	assert.Equal(t, 0, envelope.Code)

	var resp crmOpportunityResp
	require.NoError(t, json.Unmarshal(envelope.Data, &resp))
	assert.Equal(t, "2001", resp.OpportunityID)
	assert.Equal(t, "1200.5", resp.Amount)
}

func TestCRMListSalesOrdersInvalidSpaceID(t *testing.T) {
	h := newCRMTestServer(t, &fakeCRMDomainService{})

	w := ut.PerformRequest(
		h.Engine,
		http.MethodGet,
		"/api/crm/sales_order/list?space_id=invalid",
		nil,
	)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var envelope crmTestEnvelope
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &envelope))
	assert.Equal(t, http.StatusBadRequest, envelope.Code)
	assert.Equal(t, "space_id is invalid", envelope.Msg)
}
