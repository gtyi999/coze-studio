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
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/crm/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

func (r *crmRepository) CreateCustomer(ctx context.Context, customer *entity.Customer) (*entity.Customer, error) {
	customerID, err := r.idGen.GenID(ctx)
	if err != nil {
		return nil, wrapOperateError(err, "generate customer id failed")
	}

	model := toCustomerModel(customer)
	model.CustomerID = customerID
	var created *entity.Customer
	if err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(model).Error; err != nil {
			return wrapOperateError(err, "create customer failed")
		}

		created = toCustomerEntity(model)
		return r.appendAuditLog(
			ctx,
			tx,
			created.TenantID,
			created.SpaceID,
			crmAuditResourceCustomer,
			created.CustomerID,
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

func (r *crmRepository) UpdateCustomer(ctx context.Context, customer *entity.Customer) (*entity.Customer, error) {
	var beforeModel crmCustomerModel
	var afterModel crmCustomerModel

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).
			Where("id = ? AND tenant_id = ? AND space_id = ? AND is_deleted = ?", customer.CustomerID, customer.TenantID, customer.SpaceID, false).
			First(&beforeModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorx.New(errno.ErrCRMRecordNotFoundCode)
			}
			return wrapOperateError(err, "get customer before update failed")
		}

		result := tx.WithContext(ctx).
			Model(&crmCustomerModel{}).
			Where("id = ? AND tenant_id = ? AND space_id = ? AND is_deleted = ?", customer.CustomerID, customer.TenantID, customer.SpaceID, false).
			Updates(map[string]any{
				"customer_name":   customer.CustomerName,
				"customer_code":   customer.CustomerCode,
				"industry":        customer.Industry,
				"level":           customer.Level,
				"owner_user_id":   customer.OwnerUserID,
				"owner_user_name": customer.OwnerUserName,
				"status":          customer.Status,
				"mobile":          customer.Mobile,
				"email":           customer.Email,
				"address":         customer.Address,
				"remark":          customer.Remark,
				"updated_by":      customer.UpdatedBy,
				"updated_at":      timeNowMillis(),
			})
		if result.Error != nil {
			return wrapOperateError(result.Error, "update customer failed")
		}
		if result.RowsAffected == 0 {
			return errorx.New(errno.ErrCRMRecordNotFoundCode)
		}

		if err := tx.WithContext(ctx).
			Where("id = ? AND tenant_id = ? AND space_id = ? AND is_deleted = ?", customer.CustomerID, customer.TenantID, customer.SpaceID, false).
			First(&afterModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorx.New(errno.ErrCRMRecordNotFoundCode)
			}
			return wrapOperateError(err, "get customer after update failed")
		}

		return r.appendAuditLog(
			ctx,
			tx,
			customer.TenantID,
			customer.SpaceID,
			crmAuditResourceCustomer,
			customer.CustomerID,
			crmAuditActionUpdate,
			customer.UpdatedBy,
			toCustomerEntity(&beforeModel),
			toCustomerEntity(&afterModel),
		)
	})
	if err != nil {
		return nil, err
	}

	return toCustomerEntity(&afterModel), nil
}

func (r *crmRepository) DeleteCustomer(ctx context.Context, scope *entity.Scope, customerID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var beforeModel crmCustomerModel
		if err := tx.WithContext(ctx).
			Where("id = ? AND tenant_id = ? AND space_id = ? AND is_deleted = ?", customerID, scope.TenantID, scope.SpaceID, false).
			First(&beforeModel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errorx.New(errno.ErrCRMRecordNotFoundCode)
			}
			return wrapOperateError(err, "get customer before delete failed")
		}

		if err := r.softDeleteWithDB(ctx, tx, &crmCustomerModel{}, scope, customerID, "delete customer failed"); err != nil {
			return err
		}

		return r.appendAuditLog(
			ctx,
			tx,
			scope.TenantID,
			scope.SpaceID,
			crmAuditResourceCustomer,
			customerID,
			crmAuditActionDelete,
			resolveOperatorID(ctx, beforeModel.UpdatedBy),
			toCustomerEntity(&beforeModel),
			nil,
		)
	})
}

func (r *crmRepository) GetCustomer(ctx context.Context, scope *entity.Scope, customerID int64) (*entity.Customer, error) {
	var model crmCustomerModel
	if err := r.scopedFirst(ctx, &model, scope, customerID); err != nil {
		return nil, err
	}

	return toCustomerEntity(&model), nil
}

func (r *crmRepository) GetCustomerByID(ctx context.Context, customerID int64) (*entity.Customer, error) {
	var model crmCustomerModel
	if err := r.firstByID(ctx, &model, customerID); err != nil {
		return nil, err
	}

	return toCustomerEntity(&model), nil
}

func (r *crmRepository) ListCustomers(ctx context.Context, filter *entity.CustomerFilter) ([]*entity.Customer, int64, error) {
	query := r.scopedListQuery(ctx, &crmCustomerModel{}, &filter.Scope)
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
	if filter.Keyword != "" {
		like := keywordLike(filter.Keyword)
		query = query.Where(
			"customer_name LIKE ? OR customer_code LIKE ? OR mobile LIKE ? OR email LIKE ?",
			like, like, like, like,
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, wrapOperateError(err, "count customer list failed")
	}

	var models []*crmCustomerModel
	if err := query.
		Order("updated_at DESC").
		Offset(pageOffset(filter.Page, filter.PageSize)).
		Limit(filter.PageSize).
		Find(&models).Error; err != nil {
		return nil, 0, wrapOperateError(err, "list customer failed")
	}

	customers := make([]*entity.Customer, 0, len(models))
	for _, model := range models {
		customers = append(customers, toCustomerEntity(model))
	}

	return customers, total, nil
}

func (r *crmRepository) CountActiveContactsByCustomer(ctx context.Context, scope *entity.Scope, customerID int64) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&crmContactModel{}).
		Where("tenant_id = ? AND space_id = ? AND customer_id = ? AND is_deleted = ?", scope.TenantID, scope.SpaceID, customerID, false).
		Count(&total).Error
	if err != nil {
		return 0, wrapOperateError(err, "count customer contacts failed")
	}
	return total, nil
}

func (r *crmRepository) CountActiveOpportunitiesByCustomer(ctx context.Context, scope *entity.Scope, customerID int64) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&crmOpportunityModel{}).
		Where("tenant_id = ? AND space_id = ? AND customer_id = ? AND is_deleted = ?", scope.TenantID, scope.SpaceID, customerID, false).
		Count(&total).Error
	if err != nil {
		return 0, wrapOperateError(err, "count customer opportunities failed")
	}
	return total, nil
}

func (r *crmRepository) CreateContact(ctx context.Context, contact *entity.Contact) (*entity.Contact, error) {
	contactID, err := r.idGen.GenID(ctx)
	if err != nil {
		return nil, wrapOperateError(err, "generate contact id failed")
	}

	model := toContactModel(contact)
	model.ContactID = contactID
	if err = r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, wrapOperateError(err, "create contact failed")
	}

	return r.GetContact(ctx, &entity.Scope{TenantID: contact.TenantID, SpaceID: contact.SpaceID}, contactID)
}

func (r *crmRepository) UpdateContact(ctx context.Context, contact *entity.Contact) (*entity.Contact, error) {
	result := r.db.WithContext(ctx).
		Model(&crmContactModel{}).
		Where("id = ? AND tenant_id = ? AND space_id = ? AND is_deleted = ?", contact.ContactID, contact.TenantID, contact.SpaceID, false).
		Updates(map[string]any{
			"customer_id":  contact.CustomerID,
			"contact_name": contact.ContactName,
			"mobile":       contact.Mobile,
			"email":        contact.Email,
			"title":        contact.Title,
			"is_primary":   contact.IsPrimary,
			"status":       contact.Status,
			"remark":       contact.Remark,
			"updated_by":   contact.UpdatedBy,
			"updated_at":   timeNowMillis(),
		})
	if result.Error != nil {
		return nil, wrapOperateError(result.Error, "update contact failed")
	}
	if result.RowsAffected == 0 {
		return nil, errorx.New(errno.ErrCRMRecordNotFoundCode)
	}

	return r.GetContact(ctx, &entity.Scope{TenantID: contact.TenantID, SpaceID: contact.SpaceID}, contact.ContactID)
}

func (r *crmRepository) DeleteContact(ctx context.Context, scope *entity.Scope, contactID int64) error {
	return r.softDelete(ctx, &crmContactModel{}, scope, contactID, "delete contact failed")
}

func (r *crmRepository) GetContact(ctx context.Context, scope *entity.Scope, contactID int64) (*entity.Contact, error) {
	var model crmContactModel
	if err := r.scopedFirst(ctx, &model, scope, contactID); err != nil {
		return nil, err
	}

	contact := toContactEntity(&model)
	nameMap, err := r.loadCustomerNameMap(ctx, scope, []int64{contact.CustomerID})
	if err != nil {
		return nil, err
	}
	contact.CustomerName = nameMap[contact.CustomerID]

	return contact, nil
}

func (r *crmRepository) GetContactByID(ctx context.Context, contactID int64) (*entity.Contact, error) {
	var model crmContactModel
	if err := r.firstByID(ctx, &model, contactID); err != nil {
		return nil, err
	}

	return toContactEntity(&model), nil
}

func (r *crmRepository) ListContacts(ctx context.Context, filter *entity.ContactFilter) ([]*entity.Contact, int64, error) {
	query := r.scopedListQuery(ctx, &crmContactModel{}, &filter.Scope)
	if filter.CustomerID != nil && *filter.CustomerID > 0 {
		query = query.Where("customer_id = ?", *filter.CustomerID)
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
	if filter.Keyword != "" {
		like := keywordLike(filter.Keyword)
		query = query.Where(
			"contact_name LIKE ? OR mobile LIKE ? OR email LIKE ? OR title LIKE ?",
			like, like, like, like,
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, wrapOperateError(err, "count contact list failed")
	}

	var models []*crmContactModel
	if err := query.
		Order("updated_at DESC").
		Offset(pageOffset(filter.Page, filter.PageSize)).
		Limit(filter.PageSize).
		Find(&models).Error; err != nil {
		return nil, 0, wrapOperateError(err, "list contact failed")
	}

	contacts := make([]*entity.Contact, 0, len(models))
	customerIDs := make([]int64, 0, len(models))
	for _, model := range models {
		contact := toContactEntity(model)
		contacts = append(contacts, contact)
		customerIDs = append(customerIDs, contact.CustomerID)
	}

	nameMap, err := r.loadCustomerNameMap(ctx, &filter.Scope, customerIDs)
	if err != nil {
		return nil, 0, err
	}
	for _, contact := range contacts {
		contact.CustomerName = nameMap[contact.CustomerID]
	}

	return contacts, total, nil
}

func (r *crmRepository) loadCustomerNameMap(ctx context.Context, scope *entity.Scope, customerIDs []int64) (map[int64]string, error) {
	nameMap := make(map[int64]string)
	ids := uniqueInt64(customerIDs)
	if len(ids) == 0 {
		return nameMap, nil
	}

	type customerNameRow struct {
		CustomerID   int64  `gorm:"column:id"`
		CustomerName string `gorm:"column:customer_name"`
	}

	rows := make([]*customerNameRow, 0, len(ids))
	if err := r.db.WithContext(ctx).
		Model(&crmCustomerModel{}).
		Select("id", "customer_name").
		Where("tenant_id = ? AND space_id = ? AND is_deleted = ? AND id IN ?", scope.TenantID, scope.SpaceID, false, ids).
		Find(&rows).Error; err != nil {
		return nil, wrapOperateError(err, "load customer name failed")
	}

	for _, row := range rows {
		nameMap[row.CustomerID] = row.CustomerName
	}

	return nameMap, nil
}

func timeNowMillis() int64 {
	return time.Now().UnixMilli()
}
