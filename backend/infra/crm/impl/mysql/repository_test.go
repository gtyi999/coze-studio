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

package mysql

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/crm/auditctx"
	"github.com/coze-dev/coze-studio/backend/domain/crm/entity"
	mockorm "github.com/coze-dev/coze-studio/backend/internal/mock/infra/orm"
)

type stubIDGenerator struct {
	next int64
}

func (s *stubIDGenerator) GenID(_ context.Context) (int64, error) {
	s.next++
	return s.next, nil
}

func (s *stubIDGenerator) GenMultiIDs(_ context.Context, count int) ([]int64, error) {
	ids := make([]int64, 0, count)
	for i := 0; i < count; i++ {
		s.next++
		ids = append(ids, s.next)
	}
	return ids, nil
}

func TestCRMRepositoryCreateCustomerWritesAuditLog(t *testing.T) {
	t.Parallel()

	repo, db := newCRMRepositoryForTest(t)
	ctx := auditctx.WithActor(context.Background(), &auditctx.Actor{
		UserID:   88,
		TenantID: 1,
		SpaceID:  1,
	})

	created, err := repo.CreateCustomer(ctx, &entity.Customer{
		TenantID:      1,
		SpaceID:       1,
		CustomerName:  "Acme Robotics",
		CustomerCode:  "CUST-ACME",
		Industry:      "Manufacturing",
		OwnerUserID:   88,
		OwnerUserName: "Alice",
		Status:        entity.StatusActive,
		AuditInfo: entity.AuditInfo{
			CreatedBy: 88,
			UpdatedBy: 88,
		},
	})

	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, int64(1001), created.CustomerID)

	var auditRows []crmAuditLogModel
	require.NoError(t, db.Where("resource_type = ?", crmAuditResourceCustomer).Find(&auditRows).Error)
	require.Len(t, auditRows, 1)
	assert.Equal(t, int64(1), auditRows[0].TenantID)
	assert.Equal(t, int64(1), auditRows[0].SpaceID)
	assert.Equal(t, created.CustomerID, auditRows[0].ResourceID)
	assert.Equal(t, crmAuditActionCreate, auditRows[0].Action)
	assert.Equal(t, int64(88), auditRows[0].OperatorID)
	assert.Empty(t, auditRows[0].BeforeSnapshot)
	assert.True(t, strings.Contains(auditRows[0].AfterSnapshot, "Acme Robotics"))
}

func TestCRMRepositoryListCustomersScopesTenant(t *testing.T) {
	t.Parallel()

	repo, _ := newCRMRepositoryForTest(
		t,
		&crmCustomerModel{
			CustomerID:    2001,
			TenantID:      1,
			SpaceID:       1,
			CustomerName:  "Scoped Customer",
			OwnerUserID:   9,
			OwnerUserName: "Owner",
			Status:        entity.StatusActive,
		},
		&crmCustomerModel{
			CustomerID:    2002,
			TenantID:      1,
			SpaceID:       2,
			CustomerName:  "Other Space",
			OwnerUserID:   9,
			OwnerUserName: "Owner",
			Status:        entity.StatusActive,
		},
		&crmCustomerModel{
			CustomerID:    2003,
			TenantID:      2,
			SpaceID:       1,
			CustomerName:  "Other Tenant",
			OwnerUserID:   9,
			OwnerUserName: "Owner",
			Status:        entity.StatusActive,
		},
		&crmCustomerModel{
			CustomerID:    2004,
			TenantID:      1,
			SpaceID:       1,
			CustomerName:  "Soft Deleted",
			OwnerUserID:   9,
			OwnerUserName: "Owner",
			Status:        entity.StatusActive,
			IsDeleted:     true,
		},
	)

	list, total, err := repo.ListCustomers(context.Background(), &entity.CustomerFilter{
		Scope: entity.Scope{
			TenantID: 1,
			SpaceID:  1,
		},
		PageOption: entity.PageOption{
			Page:     1,
			PageSize: 20,
		},
	})

	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, list, 1)
	assert.Equal(t, int64(2001), list[0].CustomerID)
	assert.Equal(t, "Scoped Customer", list[0].CustomerName)
}

func TestCRMRepositoryDeleteCustomerSoftDeletesAndAudits(t *testing.T) {
	t.Parallel()

	originalUpdatedAt := time.Now().Add(-time.Hour).UnixMilli()
	repo, db := newCRMRepositoryForTest(
		t,
		&crmCustomerModel{
			CustomerID:    3001,
			TenantID:      1,
			SpaceID:       1,
			CustomerName:  "Delete Me",
			OwnerUserID:   7,
			OwnerUserName: "Owner",
			Status:        entity.StatusActive,
			UpdatedBy:     7,
			CreatedAt:     time.Now().Add(-24 * time.Hour).UnixMilli(),
			UpdatedAt:     originalUpdatedAt,
		},
	)

	ctx := auditctx.WithActor(context.Background(), &auditctx.Actor{
		UserID:   99,
		TenantID: 1,
		SpaceID:  1,
	})
	err := repo.DeleteCustomer(ctx, &entity.Scope{TenantID: 1, SpaceID: 1}, 3001)
	require.NoError(t, err)

	var customer crmCustomerModel
	require.NoError(t, db.First(&customer, "id = ?", 3001).Error)
	assert.True(t, customer.IsDeleted)
	assert.Equal(t, int64(99), customer.UpdatedBy)
	assert.Greater(t, customer.UpdatedAt, originalUpdatedAt)

	var auditRows []crmAuditLogModel
	require.NoError(t, db.Where("resource_type = ? AND action = ?", crmAuditResourceCustomer, crmAuditActionDelete).Find(&auditRows).Error)
	require.Len(t, auditRows, 1)
	assert.Equal(t, int64(99), auditRows[0].OperatorID)
	assert.True(t, strings.Contains(auditRows[0].BeforeSnapshot, "Delete Me"))
	assert.Empty(t, auditRows[0].AfterSnapshot)
}

func TestCRMRepositoryGetDashboardOverviewAggregatesScope(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastMonth := monthStart.Add(-24 * time.Hour)
	orderDateRecent := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)
	orderDateMid := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -10)

	repo, _ := newCRMRepositoryForTest(
		t,
		&crmCustomerModel{
			CustomerID:    4001,
			TenantID:      1,
			SpaceID:       1,
			CustomerName:  "Current Month Customer",
			OwnerUserID:   1,
			OwnerUserName: "Owner",
			Status:        entity.StatusActive,
			CreatedAt:     monthStart.Add(2 * time.Hour).UnixMilli(),
			UpdatedAt:     monthStart.Add(4 * time.Hour).UnixMilli(),
		},
		&crmCustomerModel{
			CustomerID:    4002,
			TenantID:      1,
			SpaceID:       1,
			CustomerName:  "Last Month Customer",
			OwnerUserID:   1,
			OwnerUserName: "Owner",
			Status:        entity.StatusActive,
			CreatedAt:     lastMonth.UnixMilli(),
			UpdatedAt:     lastMonth.UnixMilli(),
		},
		&crmCustomerModel{
			CustomerID:    4003,
			TenantID:      2,
			SpaceID:       2,
			CustomerName:  "Other Scope Customer",
			OwnerUserID:   1,
			OwnerUserName: "Owner",
			Status:        entity.StatusActive,
			CreatedAt:     monthStart.Add(2 * time.Hour).UnixMilli(),
			UpdatedAt:     monthStart.Add(2 * time.Hour).UnixMilli(),
		},
		&crmOpportunityModel{
			OpportunityID:   5001,
			TenantID:        1,
			SpaceID:         1,
			CustomerID:      4001,
			OpportunityName: "Scope Opportunity A",
			Amount:          100.5,
			Status:          entity.StatusOpen,
			CreatedAt:       monthStart.Add(3 * time.Hour).UnixMilli(),
			UpdatedAt:       monthStart.Add(3 * time.Hour).UnixMilli(),
		},
		&crmOpportunityModel{
			OpportunityID:   5002,
			TenantID:        1,
			SpaceID:         1,
			CustomerID:      4002,
			OpportunityName: "Scope Opportunity B",
			Amount:          200,
			Status:          entity.StatusOpen,
			CreatedAt:       lastMonth.UnixMilli(),
			UpdatedAt:       lastMonth.UnixMilli(),
		},
		&crmOpportunityModel{
			OpportunityID:   5003,
			TenantID:        2,
			SpaceID:         2,
			CustomerID:      4003,
			OpportunityName: "Other Scope Opportunity",
			Amount:          999,
			Status:          entity.StatusOpen,
			CreatedAt:       monthStart.Add(3 * time.Hour).UnixMilli(),
			UpdatedAt:       monthStart.Add(3 * time.Hour).UnixMilli(),
		},
		&crmSalesOrderModel{
			SalesOrderID:  6001,
			TenantID:      1,
			SpaceID:       1,
			CustomerID:    4001,
			ProductID:     7001,
			ProductName:   "Scope Product A",
			SalesUserID:   1,
			SalesUserName: "Owner",
			Quantity:      1,
			Amount:        10,
			OrderDate:     &orderDateRecent,
			Status:        entity.StatusActive,
			CreatedAt:     orderDateRecent.UnixMilli(),
			UpdatedAt:     orderDateRecent.UnixMilli(),
		},
		&crmSalesOrderModel{
			SalesOrderID:  6002,
			TenantID:      1,
			SpaceID:       1,
			CustomerID:    4001,
			ProductID:     7002,
			ProductName:   "Scope Product B",
			SalesUserID:   1,
			SalesUserName: "Owner",
			Quantity:      1,
			Amount:        20,
			OrderDate:     &orderDateMid,
			Status:        entity.StatusActive,
			CreatedAt:     orderDateMid.UnixMilli(),
			UpdatedAt:     orderDateMid.UnixMilli(),
		},
		&crmSalesOrderModel{
			SalesOrderID:  6003,
			TenantID:      1,
			SpaceID:       1,
			CustomerID:    4002,
			ProductID:     7003,
			ProductName:   "Scope Product C",
			SalesUserID:   1,
			SalesUserName: "Owner",
			Quantity:      1,
			Amount:        30,
			OrderDate:     &orderDateMid,
			Status:        entity.StatusActive,
			CreatedAt:     orderDateMid.UnixMilli(),
			UpdatedAt:     orderDateMid.UnixMilli(),
		},
		&crmSalesOrderModel{
			SalesOrderID:  6004,
			TenantID:      2,
			SpaceID:       2,
			CustomerID:    4003,
			ProductID:     7004,
			ProductName:   "Other Scope Product",
			SalesUserID:   1,
			SalesUserName: "Owner",
			Quantity:      1,
			Amount:        999,
			OrderDate:     &orderDateRecent,
			Status:        entity.StatusActive,
			CreatedAt:     orderDateRecent.UnixMilli(),
			UpdatedAt:     orderDateRecent.UnixMilli(),
		},
	)

	overview, err := repo.GetDashboardOverview(context.Background(), &entity.Scope{
		TenantID: 1,
		SpaceID:  1,
	})

	require.NoError(t, err)
	require.NotNil(t, overview)
	assert.Equal(t, int64(2), overview.CustomerTotal)
	assert.Equal(t, int64(1), overview.NewCustomersThisMonth)
	assert.InDelta(t, 300.5, overview.OpportunityTotalAmount, 0.001)
	assert.Equal(t, int64(1), overview.NewOpportunitiesThisMonth)
	assert.InDelta(t, 60, overview.SalesOrderTotalAmount, 0.001)
	require.Len(t, overview.RecentOrderTrend, 30)

	recentPoint := findTrendPoint(overview.RecentOrderTrend, orderDateRecent.Format(dateLayout))
	require.NotNil(t, recentPoint)
	assert.Equal(t, int64(1), recentPoint.OrderCount)
	assert.InDelta(t, 10, recentPoint.OrderAmount, 0.001)

	midPoint := findTrendPoint(overview.RecentOrderTrend, orderDateMid.Format(dateLayout))
	require.NotNil(t, midPoint)
	assert.Equal(t, int64(2), midPoint.OrderCount)
	assert.InDelta(t, 50, midPoint.OrderAmount, 0.001)
}

func newCRMRepositoryForTest(t *testing.T, rows ...any) (*crmRepository, *gorm.DB) {
	t.Helper()

	repo, db, _ := newCRMRepositoryAndDBForTest(t, rows...)
	return repo, db
}

func newCRMRepositoryAndDBForTest(t *testing.T, rows ...any) (*crmRepository, *gorm.DB, *mockorm.MockDB) {
	t.Helper()

	mockDB := mockorm.NewMockDB()
	mockDB.AddTable(&crmAuditLogModel{})
	mockDB.AddTable(&crmCustomerModel{})
	mockDB.AddTable(&crmContactModel{})
	mockDB.AddTable(&crmOpportunityModel{})
	mockDB.AddTable(&crmFollowRecordModel{})
	mockDB.AddTable(&crmProductModel{})
	mockDB.AddTable(&crmSalesOrderModel{})

	for _, row := range rows {
		switch typed := row.(type) {
		case *crmAuditLogModel:
			mockDB.AddTable(&crmAuditLogModel{}).AddRows(typed)
		case *crmCustomerModel:
			mockDB.AddTable(&crmCustomerModel{}).AddRows(typed)
		case *crmContactModel:
			mockDB.AddTable(&crmContactModel{}).AddRows(typed)
		case *crmOpportunityModel:
			mockDB.AddTable(&crmOpportunityModel{}).AddRows(typed)
		case *crmFollowRecordModel:
			mockDB.AddTable(&crmFollowRecordModel{}).AddRows(typed)
		case *crmProductModel:
			mockDB.AddTable(&crmProductModel{}).AddRows(typed)
		case *crmSalesOrderModel:
			mockDB.AddTable(&crmSalesOrderModel{}).AddRows(typed)
		default:
			t.Fatalf("unsupported crm test seed row: %T", row)
		}
	}

	db, err := mockDB.DB()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = mockDB.Close()
	})

	return &crmRepository{
		db:    db,
		idGen: &stubIDGenerator{next: 1000},
	}, db, mockDB
}

func findTrendPoint(points []*entity.DashboardOrderTrendPoint, date string) *entity.DashboardOrderTrendPoint {
	for _, point := range points {
		if point.Date == date {
			return point
		}
	}
	return nil
}
