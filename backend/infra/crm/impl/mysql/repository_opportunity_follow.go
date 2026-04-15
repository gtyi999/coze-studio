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

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/crm/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

func (r *crmRepository) CreateOpportunity(ctx context.Context, opportunity *entity.Opportunity) (*entity.Opportunity, error) {
	opportunityID, err := r.idGen.GenID(ctx)
	if err != nil {
		return nil, wrapOperateError(err, "generate opportunity id failed")
	}

	model, err := toOpportunityModel(opportunity)
	if err != nil {
		return nil, err
	}
	model.OpportunityID = opportunityID
	var created *entity.Opportunity
	if err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(model).Error; err != nil {
			return wrapOperateError(err, "create opportunity failed")
		}

		created = toOpportunityEntity(model)
		return r.appendAuditLog(
			ctx,
			tx,
			created.TenantID,
			created.SpaceID,
			crmAuditResourceOpportunity,
			created.OpportunityID,
			crmAuditActionCreate,
			created.CreatedBy,
			nil,
			created,
		)
	}); err != nil {
		return nil, err
	}

	return created, nil
}

func (r *crmRepository) UpdateOpportunity(ctx context.Context, opportunity *entity.Opportunity) (*entity.Opportunity, error) {
	expectedCloseDate, err := parseDate(opportunity.ExpectedCloseDate, opportunity.ExpectedCloseTime)
	if err != nil {
		return nil, err
	}

	result := r.db.WithContext(ctx).
		Model(&crmOpportunityModel{}).
		Where("id = ? AND tenant_id = ? AND space_id = ? AND is_deleted = ?", opportunity.OpportunityID, opportunity.TenantID, opportunity.SpaceID, false).
		Updates(map[string]any{
			"customer_id":         opportunity.CustomerID,
			"opportunity_name":    opportunity.OpportunityName,
			"stage":               opportunity.Stage,
			"amount":              opportunity.Amount,
			"expected_close_date": expectedCloseDate,
			"owner_user_id":       opportunity.OwnerUserID,
			"owner_user_name":     opportunity.OwnerUserName,
			"status":              opportunity.Status,
			"remark":              opportunity.Remark,
			"updated_by":          opportunity.UpdatedBy,
			"updated_at":          timeNowMillis(),
		})
	if result.Error != nil {
		return nil, wrapOperateError(result.Error, "update opportunity failed")
	}
	if result.RowsAffected == 0 {
		return nil, errorx.New(errno.ErrCRMRecordNotFoundCode)
	}

	return r.GetOpportunity(ctx, &entity.Scope{TenantID: opportunity.TenantID, SpaceID: opportunity.SpaceID}, opportunity.OpportunityID)
}

func (r *crmRepository) DeleteOpportunity(ctx context.Context, scope *entity.Scope, opportunityID int64) error {
	return r.softDelete(ctx, &crmOpportunityModel{}, scope, opportunityID, "delete opportunity failed")
}

func (r *crmRepository) GetOpportunity(ctx context.Context, scope *entity.Scope, opportunityID int64) (*entity.Opportunity, error) {
	var model crmOpportunityModel
	if err := r.scopedFirst(ctx, &model, scope, opportunityID); err != nil {
		return nil, err
	}

	return toOpportunityEntity(&model), nil
}

func (r *crmRepository) GetOpportunityByID(ctx context.Context, opportunityID int64) (*entity.Opportunity, error) {
	var model crmOpportunityModel
	if err := r.firstByID(ctx, &model, opportunityID); err != nil {
		return nil, err
	}

	return toOpportunityEntity(&model), nil
}

func (r *crmRepository) ListOpportunities(ctx context.Context, filter *entity.OpportunityFilter) ([]*entity.Opportunity, int64, error) {
	query := r.scopedListQuery(ctx, &crmOpportunityModel{}, &filter.Scope)
	if filter.CustomerID != nil && *filter.CustomerID > 0 {
		query = query.Where("customer_id = ?", *filter.CustomerID)
	}
	if filter.OwnerUserID != nil && *filter.OwnerUserID > 0 {
		query = query.Where("owner_user_id = ?", *filter.OwnerUserID)
	}
	if filter.Status != nil && *filter.Status != "" {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.CreatedAtStart != nil && *filter.CreatedAtStart > 0 {
		query = query.Where("created_at >= ?", *filter.CreatedAtStart)
	}
	if filter.CreatedAtEnd != nil && *filter.CreatedAtEnd > 0 {
		query = query.Where("created_at <= ?", *filter.CreatedAtEnd)
	}
	if filter.ExpectedCloseDateStart != nil && *filter.ExpectedCloseDateStart != "" {
		query = query.Where("expected_close_date >= ?", *filter.ExpectedCloseDateStart)
	}
	if filter.ExpectedCloseDateEnd != nil && *filter.ExpectedCloseDateEnd != "" {
		query = query.Where("expected_close_date <= ?", *filter.ExpectedCloseDateEnd)
	}
	if filter.Keyword != "" {
		like := keywordLike(filter.Keyword)
		query = query.Where("opportunity_name LIKE ? OR stage LIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, wrapOperateError(err, "count opportunity list failed")
	}

	var models []*crmOpportunityModel
	if err := query.
		Order("updated_at DESC").
		Offset(pageOffset(filter.Page, filter.PageSize)).
		Limit(filter.PageSize).
		Find(&models).Error; err != nil {
		return nil, 0, wrapOperateError(err, "list opportunity failed")
	}

	opportunities := make([]*entity.Opportunity, 0, len(models))
	for _, model := range models {
		opportunities = append(opportunities, toOpportunityEntity(model))
	}

	return opportunities, total, nil
}

func (r *crmRepository) CreateFollowRecord(ctx context.Context, followRecord *entity.FollowRecord) (*entity.FollowRecord, error) {
	followRecordID, err := r.idGen.GenID(ctx)
	if err != nil {
		return nil, wrapOperateError(err, "generate follow record id failed")
	}

	model := toFollowRecordModel(followRecord)
	model.FollowRecordID = followRecordID
	if err = r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, wrapOperateError(err, "create follow record failed")
	}

	return r.GetFollowRecord(ctx, &entity.Scope{TenantID: followRecord.TenantID, SpaceID: followRecord.SpaceID}, followRecordID)
}

func (r *crmRepository) UpdateFollowRecord(ctx context.Context, followRecord *entity.FollowRecord) (*entity.FollowRecord, error) {
	result := r.db.WithContext(ctx).
		Model(&crmFollowRecordModel{}).
		Where("id = ? AND tenant_id = ? AND space_id = ? AND is_deleted = ?", followRecord.FollowRecordID, followRecord.TenantID, followRecord.SpaceID, false).
		Updates(map[string]any{
			"customer_id":      followRecord.CustomerID,
			"contact_id":       nullableID(followRecord.ContactID),
			"follow_type":      followRecord.FollowType,
			"content":          followRecord.Content,
			"next_follow_time": nullableDateTime(followRecord.NextFollowTime),
			"owner_user_id":    followRecord.OwnerUserID,
			"owner_user_name":  followRecord.OwnerUserName,
			"status":           followRecord.Status,
			"updated_by":       followRecord.UpdatedBy,
			"updated_at":       timeNowMillis(),
		})
	if result.Error != nil {
		return nil, wrapOperateError(result.Error, "update follow record failed")
	}
	if result.RowsAffected == 0 {
		return nil, errorx.New(errno.ErrCRMRecordNotFoundCode)
	}

	return r.GetFollowRecord(ctx, &entity.Scope{TenantID: followRecord.TenantID, SpaceID: followRecord.SpaceID}, followRecord.FollowRecordID)
}

func (r *crmRepository) DeleteFollowRecord(ctx context.Context, scope *entity.Scope, followRecordID int64) error {
	return r.softDelete(ctx, &crmFollowRecordModel{}, scope, followRecordID, "delete follow record failed")
}

func (r *crmRepository) GetFollowRecord(ctx context.Context, scope *entity.Scope, followRecordID int64) (*entity.FollowRecord, error) {
	var model crmFollowRecordModel
	if err := r.scopedFirst(ctx, &model, scope, followRecordID); err != nil {
		return nil, err
	}

	return toFollowRecordEntity(&model), nil
}

func (r *crmRepository) GetFollowRecordByID(ctx context.Context, followRecordID int64) (*entity.FollowRecord, error) {
	var model crmFollowRecordModel
	if err := r.firstByID(ctx, &model, followRecordID); err != nil {
		return nil, err
	}

	return toFollowRecordEntity(&model), nil
}

func (r *crmRepository) ListFollowRecords(ctx context.Context, filter *entity.FollowRecordFilter) ([]*entity.FollowRecord, int64, error) {
	query := r.scopedListQuery(ctx, &crmFollowRecordModel{}, &filter.Scope)
	if filter.CustomerID != nil && *filter.CustomerID > 0 {
		query = query.Where("customer_id = ?", *filter.CustomerID)
	}
	if filter.ContactID != nil && *filter.ContactID > 0 {
		query = query.Where("contact_id = ?", *filter.ContactID)
	}
	if filter.OwnerUserID != nil && *filter.OwnerUserID > 0 {
		query = query.Where("owner_user_id = ?", *filter.OwnerUserID)
	}
	if filter.Status != nil && *filter.Status != "" {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.CreatedAtStart != nil && *filter.CreatedAtStart > 0 {
		query = query.Where("created_at >= ?", *filter.CreatedAtStart)
	}
	if filter.CreatedAtEnd != nil && *filter.CreatedAtEnd > 0 {
		query = query.Where("created_at <= ?", *filter.CreatedAtEnd)
	}
	if filter.NextFollowTimeStart != nil && *filter.NextFollowTimeStart > 0 {
		query = query.Where("next_follow_time >= ?", nullableDateTime(*filter.NextFollowTimeStart))
	}
	if filter.NextFollowTimeEnd != nil && *filter.NextFollowTimeEnd > 0 {
		query = query.Where("next_follow_time <= ?", nullableDateTime(*filter.NextFollowTimeEnd))
	}
	if filter.Keyword != "" {
		like := keywordLike(filter.Keyword)
		query = query.Where("follow_type LIKE ? OR content LIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, wrapOperateError(err, "count follow record list failed")
	}

	var models []*crmFollowRecordModel
	if err := query.
		Order("updated_at DESC").
		Offset(pageOffset(filter.Page, filter.PageSize)).
		Limit(filter.PageSize).
		Find(&models).Error; err != nil {
		return nil, 0, wrapOperateError(err, "list follow record failed")
	}

	followRecords := make([]*entity.FollowRecord, 0, len(models))
	for _, model := range models {
		followRecords = append(followRecords, toFollowRecordEntity(model))
	}

	return followRecords, total, nil
}
