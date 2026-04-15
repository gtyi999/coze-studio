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
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/crm/entity"
)

func (r *crmRepository) GetDashboardOverview(ctx context.Context, scope *entity.Scope) (*entity.DashboardOverview, error) {
	overview := &entity.DashboardOverview{
		RecentOrderTrend: make([]*entity.DashboardOrderTrendPoint, 0, 30),
	}

	now := time.Now().UTC()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	nextMonthStart := monthStart.AddDate(0, 1, 0)

	if err := r.loadDashboardCustomerStats(ctx, scope, overview, monthStart.UnixMilli(), nextMonthStart.UnixMilli()); err != nil {
		return nil, err
	}
	if err := r.loadDashboardOpportunityStats(ctx, scope, overview, monthStart.UnixMilli(), nextMonthStart.UnixMilli()); err != nil {
		return nil, err
	}
	if err := r.loadDashboardSalesOrderStats(ctx, scope, overview); err != nil {
		return nil, err
	}

	trend, err := r.loadRecentOrderTrend(ctx, scope, now)
	if err != nil {
		return nil, err
	}
	overview.RecentOrderTrend = trend

	return overview, nil
}

func (r *crmRepository) loadDashboardCustomerStats(ctx context.Context, scope *entity.Scope, overview *entity.DashboardOverview, monthStart int64, nextMonthStart int64) error {
	if err := r.scopedListQuery(ctx, &crmCustomerModel{}, scope).
		Count(&overview.CustomerTotal).Error; err != nil {
		return wrapOperateError(err, "count dashboard customer total failed")
	}

	if err := r.scopedListQuery(ctx, &crmCustomerModel{}, scope).
		Where("created_at >= ? AND created_at < ?", monthStart, nextMonthStart).
		Count(&overview.NewCustomersThisMonth).Error; err != nil {
		return wrapOperateError(err, "count dashboard monthly customers failed")
	}

	return nil
}

func (r *crmRepository) loadDashboardOpportunityStats(ctx context.Context, scope *entity.Scope, overview *entity.DashboardOverview, monthStart int64, nextMonthStart int64) error {
	type aggregateRow struct {
		Total float64 `gorm:"column:total"`
	}

	var totalRow aggregateRow
	if err := r.scopedListQuery(ctx, &crmOpportunityModel{}, scope).
		Select("COALESCE(SUM(amount), 0) AS total").
		Scan(&totalRow).Error; err != nil {
		return wrapOperateError(err, "sum dashboard opportunity amount failed")
	}
	overview.OpportunityTotalAmount = totalRow.Total

	if err := r.scopedListQuery(ctx, &crmOpportunityModel{}, scope).
		Where("created_at >= ? AND created_at < ?", monthStart, nextMonthStart).
		Count(&overview.NewOpportunitiesThisMonth).Error; err != nil {
		return wrapOperateError(err, "count dashboard monthly opportunities failed")
	}

	return nil
}

func (r *crmRepository) loadDashboardSalesOrderStats(ctx context.Context, scope *entity.Scope, overview *entity.DashboardOverview) error {
	type aggregateRow struct {
		Total float64 `gorm:"column:total"`
	}

	var totalRow aggregateRow
	if err := r.scopedListQuery(ctx, &crmSalesOrderModel{}, scope).
		Select("COALESCE(SUM(amount), 0) AS total").
		Scan(&totalRow).Error; err != nil {
		return wrapOperateError(err, "sum dashboard sales order amount failed")
	}
	overview.SalesOrderTotalAmount = totalRow.Total

	return nil
}

func (r *crmRepository) loadRecentOrderTrend(ctx context.Context, scope *entity.Scope, now time.Time) ([]*entity.DashboardOrderTrendPoint, error) {
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -29)
	endDate := startDate.AddDate(0, 0, 30)

	type orderTrendRow struct {
		OrderDate   time.Time `gorm:"column:order_date"`
		OrderCount  int64     `gorm:"column:order_count"`
		OrderAmount float64   `gorm:"column:order_amount"`
	}

	rows := make([]*orderTrendRow, 0, 30)
	if err := r.db.WithContext(ctx).
		Model(&crmSalesOrderModel{}).
		Select("order_date, COUNT(*) AS order_count, COALESCE(SUM(amount), 0) AS order_amount").
		Where(
			"tenant_id = ? AND space_id = ? AND is_deleted = ? AND order_date >= ? AND order_date < ?",
			scope.TenantID,
			scope.SpaceID,
			false,
			startDate,
			endDate,
		).
		Group("order_date").
		Order("order_date ASC").
		Scan(&rows).Error; err != nil {
		return nil, wrapOperateError(err, "list dashboard order trend failed")
	}

	trendMap := make(map[string]*entity.DashboardOrderTrendPoint, len(rows))
	for _, row := range rows {
		date := row.OrderDate.UTC().Format(dateLayout)
		trendMap[date] = &entity.DashboardOrderTrendPoint{
			Date:        date,
			OrderCount:  row.OrderCount,
			OrderAmount: row.OrderAmount,
		}
	}

	trend := make([]*entity.DashboardOrderTrendPoint, 0, 30)
	for i := 0; i < 30; i++ {
		current := startDate.AddDate(0, 0, i).Format(dateLayout)
		point, ok := trendMap[current]
		if !ok {
			point = &entity.DashboardOrderTrendPoint{
				Date:        current,
				OrderCount:  0,
				OrderAmount: 0,
			}
		}
		trend = append(trend, point)
	}

	return trend, nil
}
