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

func ListContacts(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.ContactListQuery
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}

	list, total, err := crmapp.CRMSVC.ListContacts(ctx, &entity.ContactFilter{
		Scope: entity.Scope{SpaceID: spaceID},
		PageOption: entity.PageOption{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		CustomerID:     parseOptionalInt64Param(c, req.CustomerID, "customer_id"),
		Keyword:        strings.TrimSpace(req.Keyword),
		Status:         parseOptionalString(req.Status),
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

	resp := &crmmodel.ContactListData{
		List:  make([]*crmmodel.ContactData, 0, len(list)),
		Total: total,
	}
	for _, item := range list {
		resp.List = append(resp.List, toContactData(item))
	}

	writeCRMSuccess(c, resp)
}

func GetContact(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.GetContactRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}
	contactID, ok := parseRequiredInt64Param(c, req.ContactID, "contact_id")
	if !ok {
		return
	}

	contact, err := crmapp.CRMSVC.GetContact(ctx, spaceID, contactID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toContactData(contact))
}

func CreateContact(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.CreateContactRequest
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

	contact, err := crmapp.CRMSVC.CreateContact(ctx, &entity.Contact{
		SpaceID:     spaceID,
		CustomerID:  customerID,
		ContactName: req.ContactName,
		Mobile:      req.Mobile,
		Email:       req.Email,
		Title:       req.Title,
		IsPrimary:   req.IsPrimary,
		Status:      req.Status,
		Remark:      req.Remark,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toContactData(contact))
}

func UpdateContact(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.UpdateContactRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}
	contactID, ok := parseRequiredInt64Param(c, req.ContactID, "contact_id")
	if !ok {
		return
	}
	customerID, ok := parseRequiredInt64Param(c, req.CustomerID, "customer_id")
	if !ok {
		return
	}

	contact, err := crmapp.CRMSVC.UpdateContact(ctx, &entity.Contact{
		ContactID:   contactID,
		SpaceID:     spaceID,
		CustomerID:  customerID,
		ContactName: req.ContactName,
		Mobile:      req.Mobile,
		Email:       req.Email,
		Title:       req.Title,
		IsPrimary:   req.IsPrimary,
		Status:      req.Status,
		Remark:      req.Remark,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toContactData(contact))
}

func DeleteContact(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.DeleteContactRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}
	contactID, ok := parseRequiredInt64Param(c, req.ContactID, "contact_id")
	if !ok {
		return
	}

	if err := crmapp.CRMSVC.DeleteContact(ctx, spaceID, contactID); err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, map[string]any{})
}

func toContactData(contact *entity.Contact) *crmmodel.ContactData {
	if contact == nil {
		return nil
	}
	return &crmmodel.ContactData{
		ContactID:    formatCRMInt64(contact.ContactID),
		TenantID:     formatCRMInt64(contact.TenantID),
		SpaceID:      formatCRMInt64(contact.SpaceID),
		CustomerID:   formatCRMInt64(contact.CustomerID),
		CustomerName: contact.CustomerName,
		ContactName:  contact.ContactName,
		Mobile:       contact.Mobile,
		Email:        contact.Email,
		Title:        contact.Title,
		IsPrimary:    contact.IsPrimary,
		Status:       contact.Status,
		Remark:       contact.Remark,
		CreatedBy:    formatCRMInt64(contact.CreatedBy),
		UpdatedBy:    formatCRMInt64(contact.UpdatedBy),
		CreatedAt:    formatCRMInt64(contact.CreatedAt),
		UpdatedAt:    formatCRMInt64(contact.UpdatedAt),
	}
}
