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
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	crmmodel "github.com/coze-dev/coze-studio/backend/api/model/crm"
	crmapp "github.com/coze-dev/coze-studio/backend/application/crm"
	"github.com/coze-dev/coze-studio/backend/domain/crm/entity"
)

func GetDashboardOverview(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.DashboardOverviewQuery
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}

	overview, err := crmapp.CRMSVC.GetDashboardOverview(ctx, spaceID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toDashboardOverviewData(overview))
}

func toDashboardOverviewData(overview *entity.DashboardOverview) *crmmodel.DashboardOverviewData {
	if overview == nil {
		return nil
	}

	data := &crmmodel.DashboardOverviewData{
		CustomerTotal:             overview.CustomerTotal,
		NewCustomersThisMonth:     overview.NewCustomersThisMonth,
		OpportunityTotalAmount:    formatCRMFloat(overview.OpportunityTotalAmount),
		NewOpportunitiesThisMonth: overview.NewOpportunitiesThisMonth,
		SalesOrderTotalAmount:     formatCRMFloat(overview.SalesOrderTotalAmount),
		RecentOrderTrend:          make([]*crmmodel.DashboardOrderTrendData, 0, len(overview.RecentOrderTrend)),
	}
	for _, point := range overview.RecentOrderTrend {
		if point == nil {
			continue
		}
		data.RecentOrderTrend = append(data.RecentOrderTrend, &crmmodel.DashboardOrderTrendData{
			Date:        point.Date,
			OrderCount:  point.OrderCount,
			OrderAmount: formatCRMFloat(point.OrderAmount),
		})
	}

	return data
}
