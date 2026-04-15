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
	"encoding/json"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/crm/auditctx"
)

const (
	crmAuditActionCreate = "create"
	crmAuditActionUpdate = "update"
	crmAuditActionDelete = "delete"

	crmAuditResourceCustomer    = "customer"
	crmAuditResourceOpportunity = "opportunity"
	crmAuditResourceSalesOrder  = "sales_order"
)

func (r *crmRepository) appendAuditLog(
	ctx context.Context,
	tx *gorm.DB,
	tenantID int64,
	spaceID int64,
	resourceType string,
	resourceID int64,
	action string,
	operatorID int64,
	before any,
	after any,
) error {
	logID, err := r.idGen.GenID(ctx)
	if err != nil {
		return wrapOperateError(err, "generate crm audit log id failed")
	}

	beforeSnapshot, err := marshalAuditSnapshot(before)
	if err != nil {
		return wrapOperateError(err, "marshal crm audit before snapshot failed")
	}
	afterSnapshot, err := marshalAuditSnapshot(after)
	if err != nil {
		return wrapOperateError(err, "marshal crm audit after snapshot failed")
	}

	model := &crmAuditLogModel{
		AuditLogID:     logID,
		TenantID:       tenantID,
		SpaceID:        spaceID,
		ResourceType:   resourceType,
		ResourceID:     resourceID,
		Action:         action,
		OperatorID:     resolveOperatorID(ctx, operatorID),
		BeforeSnapshot: beforeSnapshot,
		AfterSnapshot:  afterSnapshot,
		OperationAt:    timeNowMillis(),
	}

	if err := tx.WithContext(ctx).Create(model).Error; err != nil {
		return wrapOperateError(err, "create crm audit log failed")
	}

	return nil
}

func marshalAuditSnapshot(value any) (string, error) {
	if value == nil {
		return "", nil
	}

	payload, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func resolveOperatorID(ctx context.Context, fallback int64) int64 {
	actor, ok := auditctx.ActorFromContext(ctx)
	if ok && actor.UserID > 0 {
		return actor.UserID
	}
	return fallback
}
