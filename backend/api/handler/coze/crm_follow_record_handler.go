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

func ListFollowRecords(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.FollowRecordListQuery
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}

	list, total, err := crmapp.CRMSVC.ListFollowRecords(ctx, &entity.FollowRecordFilter{
		Scope: entity.Scope{SpaceID: spaceID},
		PageOption: entity.PageOption{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		CustomerID:          parseOptionalInt64Param(c, req.CustomerID, "customer_id"),
		ContactID:           parseOptionalInt64Param(c, req.ContactID, "contact_id"),
		OwnerUserID:         parseOptionalInt64Param(c, req.OwnerUserID, "owner_user_id"),
		Keyword:             strings.TrimSpace(req.Keyword),
		Status:              parseOptionalString(req.Status),
		CreatedAtStart:      parseOptionalTimestampParam(c, req.CreatedAtStart, "created_at_start"),
		CreatedAtEnd:        parseOptionalTimestampParam(c, req.CreatedAtEnd, "created_at_end"),
		NextFollowTimeStart: parseOptionalTimestampParam(c, req.NextFollowTimeStart, "next_follow_time_start"),
		NextFollowTimeEnd:   parseOptionalTimestampParam(c, req.NextFollowTimeEnd, "next_follow_time_end"),
	})
	if c.IsAborted() {
		return
	}
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	resp := &crmmodel.FollowRecordListData{
		List:  make([]*crmmodel.FollowRecordData, 0, len(list)),
		Total: total,
	}
	for _, item := range list {
		resp.List = append(resp.List, toFollowRecordData(item))
	}

	writeCRMSuccess(c, resp)
}

func GetFollowRecord(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.GetFollowRecordRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}
	followRecordID, ok := parseRequiredInt64Param(c, req.FollowRecordID, "follow_record_id")
	if !ok {
		return
	}

	record, err := crmapp.CRMSVC.GetFollowRecord(ctx, spaceID, followRecordID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toFollowRecordData(record))
}

func CreateFollowRecord(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.CreateFollowRecordRequest
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

	record, err := crmapp.CRMSVC.CreateFollowRecord(ctx, &entity.FollowRecord{
		SpaceID:        spaceID,
		CustomerID:     customerID,
		ContactID:      parseOptionalInt64Value(req.ContactID),
		FollowType:     req.FollowType,
		Content:        req.Content,
		NextFollowTime: parseOptionalTimestampValue(req.NextFollowTime),
		OwnerUserID:    parseOptionalInt64Value(req.OwnerUserID),
		OwnerUserName:  req.OwnerUserName,
		Status:         req.Status,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toFollowRecordData(record))
}

func UpdateFollowRecord(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.UpdateFollowRecordRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}
	followRecordID, ok := parseRequiredInt64Param(c, req.FollowRecordID, "follow_record_id")
	if !ok {
		return
	}
	customerID, ok := parseRequiredInt64Param(c, req.CustomerID, "customer_id")
	if !ok {
		return
	}

	record, err := crmapp.CRMSVC.UpdateFollowRecord(ctx, &entity.FollowRecord{
		FollowRecordID: followRecordID,
		SpaceID:        spaceID,
		CustomerID:     customerID,
		ContactID:      parseOptionalInt64Value(req.ContactID),
		FollowType:     req.FollowType,
		Content:        req.Content,
		NextFollowTime: parseOptionalTimestampValue(req.NextFollowTime),
		OwnerUserID:    parseOptionalInt64Value(req.OwnerUserID),
		OwnerUserName:  req.OwnerUserName,
		Status:         req.Status,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toFollowRecordData(record))
}

func DeleteFollowRecord(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.DeleteFollowRecordRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}
	followRecordID, ok := parseRequiredInt64Param(c, req.FollowRecordID, "follow_record_id")
	if !ok {
		return
	}

	if err := crmapp.CRMSVC.DeleteFollowRecord(ctx, spaceID, followRecordID); err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, map[string]any{})
}

func toFollowRecordData(record *entity.FollowRecord) *crmmodel.FollowRecordData {
	if record == nil {
		return nil
	}
	return &crmmodel.FollowRecordData{
		FollowRecordID: formatCRMInt64(record.FollowRecordID),
		TenantID:       formatCRMInt64(record.TenantID),
		SpaceID:        formatCRMInt64(record.SpaceID),
		CustomerID:     formatCRMInt64(record.CustomerID),
		ContactID:      formatCRMInt64(record.ContactID),
		FollowType:     record.FollowType,
		Content:        record.Content,
		NextFollowTime: formatCRMInt64(record.NextFollowTime),
		OwnerUserID:    formatCRMInt64(record.OwnerUserID),
		OwnerUserName:  record.OwnerUserName,
		Status:         record.Status,
		CreatedBy:      formatCRMInt64(record.CreatedBy),
		UpdatedBy:      formatCRMInt64(record.UpdatedBy),
		CreatedAt:      formatCRMInt64(record.CreatedAt),
		UpdatedAt:      formatCRMInt64(record.UpdatedAt),
	}
}
