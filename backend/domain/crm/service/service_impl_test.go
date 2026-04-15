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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/coze-dev/coze-studio/backend/domain/crm/entity"
	mockcrm "github.com/coze-dev/coze-studio/backend/internal/mock/domain/crm"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

func TestCRMServiceDeleteCustomer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		contactCount    int64
		oppCount        int64
		followCount     int64
		salesOrderCount int64
		expectedCode    int32
		expectDelete    bool
	}{
		{
			name:         "block when active contacts exist",
			contactCount: 1,
			expectedCode: errno.ErrCRMOperateCode,
		},
		{
			name:         "block when active opportunities exist",
			contactCount: 0,
			oppCount:     1,
			expectedCode: errno.ErrCRMOperateCode,
		},
		{
			name:         "block when follow records exist",
			contactCount: 0,
			oppCount:     0,
			followCount:  1,
			expectedCode: errno.ErrCRMOperateCode,
		},
		{
			name:            "block when sales orders exist",
			contactCount:    0,
			oppCount:        0,
			followCount:     0,
			salesOrderCount: 1,
			expectedCode:    errno.ErrCRMOperateCode,
		},
		{
			name:         "delete when no active dependencies",
			expectedCode: 0,
			expectDelete: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mockcrm.NewMockRepository(ctrl)
			svc := NewService(&Components{Repository: repo})

			scope := &entity.Scope{TenantID: 1, SpaceID: 100}
			repo.EXPECT().GetCustomerByID(gomock.Any(), int64(2001)).Return(&entity.Customer{
				CustomerID: 2001,
				TenantID:   scope.TenantID,
				SpaceID:    scope.SpaceID,
			}, nil)
			repo.EXPECT().CountActiveContactsByCustomer(gomock.Any(), scope, int64(2001)).Return(tt.contactCount, nil)

			if tt.contactCount == 0 {
				repo.EXPECT().CountActiveOpportunitiesByCustomer(gomock.Any(), scope, int64(2001)).Return(tt.oppCount, nil)
			}
			if tt.contactCount == 0 && tt.oppCount == 0 {
				repo.EXPECT().ListFollowRecords(gomock.Any(), gomock.Any()).Return(nil, tt.followCount, nil)
			}
			if tt.contactCount == 0 && tt.oppCount == 0 && tt.followCount == 0 {
				repo.EXPECT().ListSalesOrders(gomock.Any(), gomock.Any()).Return(nil, tt.salesOrderCount, nil)
			}
			if tt.expectDelete {
				repo.EXPECT().DeleteCustomer(gomock.Any(), scope, int64(2001)).Return(nil)
			}

			err := svc.DeleteCustomer(context.Background(), scope, 2001)
			if tt.expectedCode == 0 {
				require.NoError(t, err)
				return
			}

			requireErrorCode(t, err, tt.expectedCode)
		})
	}
}

func TestCRMServiceDeleteContact(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		followCount  int64
		expectedCode int32
		expectDelete bool
	}{
		{
			name:         "block when follow records exist",
			followCount:  1,
			expectedCode: errno.ErrCRMOperateCode,
		},
		{
			name:         "delete when no active dependencies",
			expectedCode: 0,
			expectDelete: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mockcrm.NewMockRepository(ctrl)
			svc := NewService(&Components{Repository: repo})

			scope := &entity.Scope{TenantID: 1, SpaceID: 100}
			repo.EXPECT().GetContactByID(gomock.Any(), int64(2002)).Return(&entity.Contact{
				ContactID: 2002,
				TenantID:  scope.TenantID,
				SpaceID:   scope.SpaceID,
			}, nil)

			repo.EXPECT().ListFollowRecords(gomock.Any(), gomock.Any()).Return(nil, tt.followCount, nil)
			if tt.expectDelete {
				repo.EXPECT().DeleteContact(gomock.Any(), scope, int64(2002)).Return(nil)
			}

			err := svc.DeleteContact(context.Background(), scope, 2002)
			if tt.expectedCode == 0 {
				require.NoError(t, err)
				return
			}

			requireErrorCode(t, err, tt.expectedCode)
		})
	}
}

func TestCRMServiceDeleteOpportunity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		followCount  int64
		orderCount   int64
		expectedCode int32
		expectDelete bool
	}{
		{
			name:         "block when follow records exist",
			followCount:  1,
			expectedCode: errno.ErrCRMOperateCode,
		},
		{
			name:         "block when sales orders exist",
			followCount:  0,
			orderCount:   1,
			expectedCode: errno.ErrCRMOperateCode,
		},
		{
			name:         "delete when no active dependencies",
			expectedCode: 0,
			expectDelete: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mockcrm.NewMockRepository(ctrl)
			svc := NewService(&Components{Repository: repo})

			scope := &entity.Scope{TenantID: 1, SpaceID: 200}
			repo.EXPECT().GetOpportunityByID(gomock.Any(), int64(3002)).Return(&entity.Opportunity{
				OpportunityID: 3002,
				TenantID:      scope.TenantID,
				SpaceID:       scope.SpaceID,
			}, nil)

			repo.EXPECT().ListFollowRecords(gomock.Any(), gomock.Any()).Return(nil, tt.followCount, nil)
			if tt.followCount == 0 {
				repo.EXPECT().ListSalesOrders(gomock.Any(), gomock.Any()).Return(nil, tt.orderCount, nil)
			}
			if tt.expectDelete {
				repo.EXPECT().DeleteOpportunity(gomock.Any(), scope, int64(3002)).Return(nil)
			}

			err := svc.DeleteOpportunity(context.Background(), scope, 3002)
			if tt.expectedCode == 0 {
				require.NoError(t, err)
				return
			}

			requireErrorCode(t, err, tt.expectedCode)
		})
	}
}

func TestCRMServiceDeleteProduct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		orderCount   int64
		expectedCode int32
		expectDelete bool
	}{
		{
			name:         "block when sales orders exist",
			orderCount:   1,
			expectedCode: errno.ErrCRMOperateCode,
		},
		{
			name:         "delete when no active dependencies",
			expectedCode: 0,
			expectDelete: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mockcrm.NewMockRepository(ctrl)
			svc := NewService(&Components{Repository: repo})

			scope := &entity.Scope{TenantID: 1, SpaceID: 300}
			repo.EXPECT().GetProductByID(gomock.Any(), int64(4002)).Return(&entity.Product{
				ProductID: 4002,
				TenantID:  scope.TenantID,
				SpaceID:   scope.SpaceID,
			}, nil)

			repo.EXPECT().ListSalesOrders(gomock.Any(), gomock.Any()).Return(nil, tt.orderCount, nil)
			if tt.expectDelete {
				repo.EXPECT().DeleteProduct(gomock.Any(), scope, int64(4002)).Return(nil)
			}

			err := svc.DeleteProduct(context.Background(), scope, 4002)
			if tt.expectedCode == 0 {
				require.NoError(t, err)
				return
			}

			requireErrorCode(t, err, tt.expectedCode)
		})
	}
}

func TestCRMServiceCreateContact(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		customer      *entity.Customer
		expectedCode  int32
		expectCreate  bool
		expectedTitle string
	}{
		{
			name: "create contact under same tenant customer",
			customer: &entity.Customer{
				CustomerID: 9001,
				TenantID:   1,
				SpaceID:    101,
			},
			expectCreate:  true,
			expectedTitle: "VP Sales",
		},
		{
			name: "reject cross tenant customer write",
			customer: &entity.Customer{
				CustomerID: 9001,
				TenantID:   2,
				SpaceID:    101,
			},
			expectedCode: errno.ErrCRMPermissionCode,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mockcrm.NewMockRepository(ctrl)
			svc := NewService(&Components{Repository: repo})

			contact := &entity.Contact{
				TenantID:    1,
				SpaceID:     101,
				CustomerID:  9001,
				ContactName: " Alice ",
				Title:       " VP Sales ",
				AuditInfo: entity.AuditInfo{
					CreatedBy: 7,
					UpdatedBy: 7,
				},
			}

			repo.EXPECT().GetCustomerByID(gomock.Any(), int64(9001)).Return(tt.customer, nil)
			if tt.expectCreate {
				repo.EXPECT().CreateContact(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, got *entity.Contact) (*entity.Contact, error) {
					assert.Equal(t, tt.expectedTitle, got.Title)
					assert.Equal(t, tt.expectedTitle, got.Position)
					assert.Equal(t, entity.StatusActive, got.Status)
					got.ContactID = 5001
					return got, nil
				})
			}

			result, err := svc.CreateContact(context.Background(), contact)
			if tt.expectedCode != 0 {
				require.Nil(t, result)
				requireErrorCode(t, err, tt.expectedCode)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, int64(5001), result.ContactID)
		})
	}
}

func TestCRMServiceCreateOpportunity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		opportunity  *entity.Opportunity
		customer     *entity.Customer
		expectedCode int32
		expectCreate bool
	}{
		{
			name: "reject negative amount",
			opportunity: &entity.Opportunity{
				TenantID:        1,
				SpaceID:         201,
				CustomerID:      3001,
				OpportunityName: "Renewal",
				Amount:          -1,
				AuditInfo: entity.AuditInfo{
					CreatedBy: 9,
					UpdatedBy: 9,
				},
			},
			expectedCode: errno.ErrCRMInvalidParamCode,
		},
		{
			name: "reject cross tenant customer write",
			opportunity: &entity.Opportunity{
				TenantID:        1,
				SpaceID:         201,
				CustomerID:      3001,
				OpportunityName: "Renewal",
				Amount:          1000,
				AuditInfo: entity.AuditInfo{
					CreatedBy: 9,
					UpdatedBy: 9,
				},
			},
			customer: &entity.Customer{
				CustomerID: 3001,
				TenantID:   2,
				SpaceID:    201,
			},
			expectedCode: errno.ErrCRMPermissionCode,
		},
		{
			name: "create opportunity with defaults",
			opportunity: &entity.Opportunity{
				TenantID:        1,
				SpaceID:         201,
				CustomerID:      3001,
				OpportunityName: "Renewal",
				Amount:          1000,
				AuditInfo: entity.AuditInfo{
					CreatedBy: 9,
					UpdatedBy: 9,
				},
			},
			customer: &entity.Customer{
				CustomerID: 3001,
				TenantID:   1,
				SpaceID:    201,
			},
			expectCreate: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mockcrm.NewMockRepository(ctrl)
			svc := NewService(&Components{Repository: repo})

			if tt.customer != nil {
				repo.EXPECT().GetCustomerByID(gomock.Any(), tt.opportunity.CustomerID).Return(tt.customer, nil)
			}
			if tt.expectCreate {
				repo.EXPECT().CreateOpportunity(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, got *entity.Opportunity) (*entity.Opportunity, error) {
					assert.Equal(t, entity.StageInitial, got.Stage)
					assert.Equal(t, entity.StatusOpen, got.Status)
					assert.Equal(t, got.CreatedBy, got.OwnerUserID)
					got.OpportunityID = 6001
					return got, nil
				})
			}

			result, err := svc.CreateOpportunity(context.Background(), tt.opportunity)
			if tt.expectedCode != 0 {
				require.Nil(t, result)
				requireErrorCode(t, err, tt.expectedCode)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, int64(6001), result.OpportunityID)
		})
	}
}

func TestCRMServiceCreateSalesOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		order        *entity.SalesOrder
		customer     *entity.Customer
		product      *entity.Product
		expectedCode int32
		expectCreate bool
	}{
		{
			name: "reject negative quantity",
			order: &entity.SalesOrder{
				TenantID:   1,
				SpaceID:    301,
				CustomerID: 4001,
				ProductID:  5001,
				Quantity:   -2,
				Amount:     100,
				AuditInfo: entity.AuditInfo{
					CreatedBy: 11,
					UpdatedBy: 11,
				},
			},
			expectedCode: errno.ErrCRMInvalidParamCode,
		},
		{
			name: "reject negative amount",
			order: &entity.SalesOrder{
				TenantID:   1,
				SpaceID:    301,
				CustomerID: 4001,
				ProductID:  5001,
				Quantity:   2,
				Amount:     -100,
				AuditInfo: entity.AuditInfo{
					CreatedBy: 11,
					UpdatedBy: 11,
				},
			},
			expectedCode: errno.ErrCRMInvalidParamCode,
		},
		{
			name: "reject cross tenant product write",
			order: &entity.SalesOrder{
				TenantID:   1,
				SpaceID:    301,
				CustomerID: 4001,
				ProductID:  5001,
				Quantity:   2,
				Amount:     100,
				AuditInfo: entity.AuditInfo{
					CreatedBy: 11,
					UpdatedBy: 11,
				},
			},
			customer: &entity.Customer{
				CustomerID: 4001,
				TenantID:   1,
				SpaceID:    301,
			},
			product: &entity.Product{
				ProductID: 5001,
				TenantID:  2,
				SpaceID:   301,
			},
			expectedCode: errno.ErrCRMPermissionCode,
		},
		{
			name: "create sales order with product snapshot defaults",
			order: &entity.SalesOrder{
				TenantID:   1,
				SpaceID:    301,
				CustomerID: 4001,
				ProductID:  5001,
				Quantity:   2,
				Amount:     100,
				AuditInfo: entity.AuditInfo{
					CreatedBy: 11,
					UpdatedBy: 11,
				},
			},
			customer: &entity.Customer{
				CustomerID: 4001,
				TenantID:   1,
				SpaceID:    301,
			},
			product: &entity.Product{
				ProductID:   5001,
				TenantID:    1,
				SpaceID:     301,
				ProductName: "AI Seats",
			},
			expectCreate: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			repo := mockcrm.NewMockRepository(ctrl)
			svc := NewService(&Components{Repository: repo})

			if tt.customer != nil {
				repo.EXPECT().GetCustomerByID(gomock.Any(), tt.order.CustomerID).Return(tt.customer, nil)
			}
			if tt.product != nil {
				repo.EXPECT().GetProductByID(gomock.Any(), tt.order.ProductID).Return(tt.product, nil)
			}
			if tt.expectCreate {
				repo.EXPECT().CreateSalesOrder(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, got *entity.SalesOrder) (*entity.SalesOrder, error) {
					assert.Equal(t, entity.StatusDraft, got.Status)
					assert.Equal(t, "AI Seats", got.ProductName)
					assert.Equal(t, int64(11), got.SalesUserID)
					got.SalesOrderID = 7001
					return got, nil
				})
			}

			result, err := svc.CreateSalesOrder(context.Background(), tt.order)
			if tt.expectedCode != 0 {
				require.Nil(t, result)
				requireErrorCode(t, err, tt.expectedCode)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, int64(7001), result.SalesOrderID)
		})
	}
}

func requireErrorCode(t *testing.T, err error, code int32) {
	t.Helper()

	require.Error(t, err)
	var statusErr errorx.StatusError
	require.True(t, errors.As(err, &statusErr))
	assert.Equal(t, code, statusErr.Code())
}
