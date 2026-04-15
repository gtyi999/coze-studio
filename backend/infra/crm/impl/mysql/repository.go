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
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/crm/entity"
	"github.com/coze-dev/coze-studio/backend/domain/crm/repository"
	"github.com/coze-dev/coze-studio/backend/infra/idgen"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

const dateLayout = "2006-01-02"

type crmRepository struct {
	db    *gorm.DB
	idGen idgen.IDGenerator
}

func NewRepository(db *gorm.DB, idGenSVC idgen.IDGenerator) repository.Repository {
	return &crmRepository{
		db:    db,
		idGen: idGenSVC,
	}
}

func (r *crmRepository) softDelete(ctx context.Context, model any, scope *entity.Scope, id int64, message string) error {
	return r.softDeleteWithDB(ctx, r.db, model, scope, id, message)
}

func (r *crmRepository) softDeleteWithDB(ctx context.Context, db *gorm.DB, model any, scope *entity.Scope, id int64, message string) error {
	updates := map[string]any{
		"is_deleted": true,
		"updated_at": time.Now().UnixMilli(),
	}
	if operatorID := resolveOperatorID(ctx, 0); operatorID > 0 {
		updates["updated_by"] = operatorID
	}

	result := db.WithContext(ctx).
		Model(model).
		Where("id = ? AND tenant_id = ? AND space_id = ? AND is_deleted = ?", id, scope.TenantID, scope.SpaceID, false).
		Updates(updates)
	if result.Error != nil {
		return wrapOperateError(result.Error, message)
	}
	if result.RowsAffected == 0 {
		return errorx.New(errno.ErrCRMRecordNotFoundCode)
	}
	return nil
}

func (r *crmRepository) scopedFirst(ctx context.Context, dst any, scope *entity.Scope, id int64) error {
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ? AND space_id = ? AND is_deleted = ?", id, scope.TenantID, scope.SpaceID, false).
		First(dst).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorx.New(errno.ErrCRMRecordNotFoundCode)
		}
		return wrapOperateError(err, "get crm record failed")
	}
	return nil
}

func (r *crmRepository) firstByID(ctx context.Context, dst any, id int64) error {
	err := r.db.WithContext(ctx).
		Where("id = ? AND is_deleted = ?", id, false).
		First(dst).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorx.New(errno.ErrCRMRecordNotFoundCode)
		}
		return wrapOperateError(err, "get crm record failed")
	}
	return nil
}

func (r *crmRepository) scopedListQuery(ctx context.Context, model any, scope *entity.Scope) *gorm.DB {
	return r.db.WithContext(ctx).
		Model(model).
		Where("tenant_id = ? AND space_id = ? AND is_deleted = ?", scope.TenantID, scope.SpaceID, false)
}

func pageOffset(page int, pageSize int) int {
	if page <= 1 {
		return 0
	}
	return (page - 1) * pageSize
}

func keywordLike(keyword string) string {
	return "%" + strings.TrimSpace(keyword) + "%"
}

func uniqueInt64(values []int64) []int64 {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(values))
	result := make([]int64, 0, len(values))
	for _, value := range values {
		if value <= 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func wrapOperateError(err error, message string) error {
	return errorx.WrapByCode(err, errno.ErrCRMOperateCode, errorx.KV("msg", message))
}

func parseDate(value string, ms int64) (*time.Time, error) {
	value = strings.TrimSpace(value)
	switch {
	case value != "":
		t, err := time.Parse(dateLayout, value)
		if err != nil {
			return nil, errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "invalid date format, expected YYYY-MM-DD"))
		}
		return &t, nil
	case ms > 0:
		t := time.UnixMilli(ms).UTC()
		date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		return &date, nil
	default:
		return nil, nil
	}
}

func formatDate(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format(dateLayout)
}

func dateToMillis(value *time.Time) int64 {
	if value == nil {
		return 0
	}
	date := time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
	return date.UnixMilli()
}

func nullableDateTime(ms int64) *time.Time {
	if ms <= 0 {
		return nil
	}
	t := time.UnixMilli(ms).UTC()
	return &t
}

func dateTimeToMillis(value *time.Time) int64 {
	if value == nil {
		return 0
	}
	return value.UnixMilli()
}

func nullableID(id int64) any {
	if id <= 0 {
		return nil
	}
	return id
}
