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

func ListCustomers(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.CustomerListQuery
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}

	list, total, err := crmapp.CRMSVC.ListCustomers(ctx, &entity.CustomerFilter{
		Scope: entity.Scope{SpaceID: spaceID},
		PageOption: entity.PageOption{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		Keyword:        strings.TrimSpace(req.Keyword),
		Status:         parseOptionalString(req.Status),
		OwnerUserID:    parseOptionalInt64Param(c, req.OwnerUserID, "owner_user_id"),
		CreatedAtStart: parseOptionalTimestampParam(c, req.CreatedAtStart, "created_at_start"),
		CreatedAtEnd:   parseOptionalTimestampParam(c, req.CreatedAtEnd, "created_at_end"),
	})
	if c.IsAborted() {
		return
	}
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	resp := &crmmodel.CustomerListData{
		List:  make([]*crmmodel.CustomerData, 0, len(list)),
		Total: total,
	}
	for _, item := range list {
		resp.List = append(resp.List, toCustomerData(item))
	}

	writeCRMSuccess(c, resp)
}

func GetCustomer(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.GetCustomerRequest
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

	customer, err := crmapp.CRMSVC.GetCustomer(ctx, spaceID, customerID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toCustomerData(customer))
}

func CreateCustomer(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.CreateCustomerRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}

	customer, err := crmapp.CRMSVC.CreateCustomer(ctx, &entity.Customer{
		SpaceID:       spaceID,
		CustomerName:  req.CustomerName,
		CustomerCode:  req.CustomerCode,
		Industry:      req.Industry,
		Level:         req.Level,
		OwnerUserID:   parseOptionalInt64Value(req.OwnerUserID),
		OwnerUserName: req.OwnerUserName,
		Status:        req.Status,
		Mobile:        req.Mobile,
		Email:         req.Email,
		Address:       req.Address,
		Remark:        req.Remark,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toCustomerData(customer))
}

func UpdateCustomer(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.UpdateCustomerRequest
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

	customer, err := crmapp.CRMSVC.UpdateCustomer(ctx, &entity.Customer{
		CustomerID:    customerID,
		SpaceID:       spaceID,
		CustomerName:  req.CustomerName,
		CustomerCode:  req.CustomerCode,
		Industry:      req.Industry,
		Level:         req.Level,
		OwnerUserID:   parseOptionalInt64Value(req.OwnerUserID),
		OwnerUserName: req.OwnerUserName,
		Status:        req.Status,
		Mobile:        req.Mobile,
		Email:         req.Email,
		Address:       req.Address,
		Remark:        req.Remark,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toCustomerData(customer))
}

func DeleteCustomer(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.DeleteCustomerRequest
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

	if err := crmapp.CRMSVC.DeleteCustomer(ctx, spaceID, customerID); err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, map[string]any{})
}

func toCustomerData(customer *entity.Customer) *crmmodel.CustomerData {
	if customer == nil {
		return nil
	}
	return &crmmodel.CustomerData{
		CustomerID:    formatCRMInt64(customer.CustomerID),
		TenantID:      formatCRMInt64(customer.TenantID),
		SpaceID:       formatCRMInt64(customer.SpaceID),
		CustomerName:  customer.CustomerName,
		CustomerCode:  customer.CustomerCode,
		Industry:      customer.Industry,
		Level:         customer.Level,
		OwnerUserID:   formatCRMInt64(customer.OwnerUserID),
		OwnerUserName: customer.OwnerUserName,
		Status:        customer.Status,
		Mobile:        customer.Mobile,
		Email:         customer.Email,
		Address:       customer.Address,
		Remark:        customer.Remark,
		CreatedBy:     formatCRMInt64(customer.CreatedBy),
		UpdatedBy:     formatCRMInt64(customer.UpdatedBy),
		CreatedAt:     formatCRMInt64(customer.CreatedAt),
		UpdatedAt:     formatCRMInt64(customer.UpdatedAt),
	}
}
