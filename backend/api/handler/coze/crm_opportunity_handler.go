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
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	crmmodel "github.com/coze-dev/coze-studio/backend/api/model/crm"
	crmapp "github.com/coze-dev/coze-studio/backend/application/crm"
	"github.com/coze-dev/coze-studio/backend/domain/crm/entity"
)

func ListOpportunities(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.OpportunityListQuery
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}

	list, total, err := crmapp.CRMSVC.ListOpportunities(ctx, &entity.OpportunityFilter{
		Scope: entity.Scope{SpaceID: spaceID},
		PageOption: entity.PageOption{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		CustomerID:             parseOptionalInt64Param(c, req.CustomerID, "customer_id"),
		OwnerUserID:            parseOptionalInt64Param(c, req.OwnerUserID, "owner_user_id"),
		Keyword:                strings.TrimSpace(req.Keyword),
		Status:                 parseOptionalString(req.Status),
		CreatedAtStart:         parseOptionalTimestampParam(c, req.CreatedAtStart, "created_at_start"),
		CreatedAtEnd:           parseOptionalTimestampParam(c, req.CreatedAtEnd, "created_at_end"),
		ExpectedCloseDateStart: parseOptionalDate(req.ExpectedCloseDateStart),
		ExpectedCloseDateEnd:   parseOptionalDate(req.ExpectedCloseDateEnd),
	})
	if c.IsAborted() {
		return
	}
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	resp := &crmmodel.OpportunityListData{
		List:  make([]*crmmodel.OpportunityData, 0, len(list)),
		Total: total,
	}
	for _, item := range list {
		resp.List = append(resp.List, toOpportunityData(item))
	}

	writeCRMSuccess(c, resp)
}

func GetOpportunity(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.GetOpportunityRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}
	opportunityID, ok := parseRequiredInt64Param(c, req.OpportunityID, "opportunity_id")
	if !ok {
		return
	}

	opportunity, err := crmapp.CRMSVC.GetOpportunity(ctx, spaceID, opportunityID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toOpportunityData(opportunity))
}

func CreateOpportunity(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.CreateOpportunityRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}
	customerID, ok := parseRequiredInt64Param(c, req.CustomerID, "customer_id")
	if !ok {
		return
	}
	amount, ok := parseDecimalParam(c, req.Amount, "amount", false)
	if !ok {
		return
	}

	opportunity, err := crmapp.CRMSVC.CreateOpportunity(ctx, &entity.Opportunity{
		SpaceID:           spaceID,
		CustomerID:        customerID,
		OpportunityName:   req.OpportunityName,
		Stage:             req.Stage,
		Amount:            amount,
		ExpectedCloseDate: req.ExpectedCloseDate,
		OwnerUserID:       parseOptionalInt64Value(req.OwnerUserID),
		OwnerUserName:     req.OwnerUserName,
		Status:            req.Status,
		Remark:            req.Remark,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toOpportunityData(opportunity))
}

func UpdateOpportunity(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.UpdateOpportunityRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}
	opportunityID, ok := parseRequiredInt64Param(c, req.OpportunityID, "opportunity_id")
	if !ok {
		return
	}
	customerID, ok := parseRequiredInt64Param(c, req.CustomerID, "customer_id")
	if !ok {
		return
	}
	amount, ok := parseDecimalParam(c, req.Amount, "amount", false)
	if !ok {
		return
	}

	opportunity, err := crmapp.CRMSVC.UpdateOpportunity(ctx, &entity.Opportunity{
		OpportunityID:     opportunityID,
		SpaceID:           spaceID,
		CustomerID:        customerID,
		OpportunityName:   req.OpportunityName,
		Stage:             req.Stage,
		Amount:            amount,
		ExpectedCloseDate: req.ExpectedCloseDate,
		OwnerUserID:       parseOptionalInt64Value(req.OwnerUserID),
		OwnerUserName:     req.OwnerUserName,
		Status:            req.Status,
		Remark:            req.Remark,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toOpportunityData(opportunity))
}

func DeleteOpportunity(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.DeleteOpportunityRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}
	opportunityID, ok := parseRequiredInt64Param(c, req.OpportunityID, "opportunity_id")
	if !ok {
		return
	}

	if err := crmapp.CRMSVC.DeleteOpportunity(ctx, spaceID, opportunityID); err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, map[string]any{})
}

func toOpportunityData(opportunity *entity.Opportunity) *crmmodel.OpportunityData {
	if opportunity == nil {
		return nil
	}
	return &crmmodel.OpportunityData{
		OpportunityID:     formatCRMInt64(opportunity.OpportunityID),
		TenantID:          formatCRMInt64(opportunity.TenantID),
		SpaceID:           formatCRMInt64(opportunity.SpaceID),
		CustomerID:        formatCRMInt64(opportunity.CustomerID),
		OpportunityName:   opportunity.OpportunityName,
		Stage:             opportunity.Stage,
		Amount:            formatCRMFloat(opportunity.Amount),
		ExpectedCloseDate: opportunity.ExpectedCloseDate,
		OwnerUserID:       formatCRMInt64(opportunity.OwnerUserID),
		OwnerUserName:     opportunity.OwnerUserName,
		Status:            opportunity.Status,
		Remark:            opportunity.Remark,
		CreatedBy:         formatCRMInt64(opportunity.CreatedBy),
		UpdatedBy:         formatCRMInt64(opportunity.UpdatedBy),
		CreatedAt:         formatCRMInt64(opportunity.CreatedAt),
		UpdatedAt:         formatCRMInt64(opportunity.UpdatedAt),
	}
}
