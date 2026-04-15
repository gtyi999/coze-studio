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

func (r *crmRepository) CreateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	productID, err := r.idGen.GenID(ctx)
	if err != nil {
		return nil, wrapOperateError(err, "generate product id failed")
	}

	model := toProductModel(product)
	model.ProductID = productID
	if err = r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, wrapOperateError(err, "create product failed")
	}

	return toProductEntity(model), nil
}

func (r *crmRepository) UpdateProduct(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	result := r.db.WithContext(ctx).
		Model(&crmProductModel{}).
		Where("id = ? AND tenant_id = ? AND space_id = ? AND is_deleted = ?", product.ProductID, product.TenantID, product.SpaceID, false).
		Updates(map[string]any{
			"product_name": product.ProductName,
			"product_code": product.ProductCode,
			"category":     product.Category,
			"unit_price":   product.UnitPrice,
			"status":       product.Status,
			"remark":       product.Remark,
			"updated_by":   product.UpdatedBy,
			"updated_at":   timeNowMillis(),
		})
	if result.Error != nil {
		return nil, wrapOperateError(result.Error, "update product failed")
	}
	if result.RowsAffected == 0 {
		return nil, errorx.New(errno.ErrCRMRecordNotFoundCode)
	}

	return r.GetProduct(ctx, &entity.Scope{TenantID: product.TenantID, SpaceID: product.SpaceID}, product.ProductID)
}

func (r *crmRepository) DeleteProduct(ctx context.Context, scope *entity.Scope, productID int64) error {
	return r.softDelete(ctx, &crmProductModel{}, scope, productID, "delete product failed")
}

func (r *crmRepository) GetProduct(ctx context.Context, scope *entity.Scope, productID int64) (*entity.Product, error) {
	var model crmProductModel
	if err := r.scopedFirst(ctx, &model, scope, productID); err != nil {
		return nil, err
	}

	return toProductEntity(&model), nil
}

func (r *crmRepository) GetProductByID(ctx context.Context, productID int64) (*entity.Product, error) {
	var model crmProductModel
	if err := r.firstByID(ctx, &model, productID); err != nil {
		return nil, err
	}

	return toProductEntity(&model), nil
}

func (r *crmRepository) ListProducts(ctx context.Context, filter *entity.ProductFilter) ([]*entity.Product, int64, error) {
	query := r.scopedListQuery(ctx, &crmProductModel{}, &filter.Scope)
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
			"product_name LIKE ? OR product_code LIKE ? OR category LIKE ?",
			like, like, like,
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, wrapOperateError(err, "count product list failed")
	}

	var models []*crmProductModel
	if err := query.
		Order("updated_at DESC").
		Offset(pageOffset(filter.Page, filter.PageSize)).
		Limit(filter.PageSize).
		Find(&models).Error; err != nil {
		return nil, 0, wrapOperateError(err, "list product failed")
	}

	products := make([]*entity.Product, 0, len(models))
	for _, model := range models {
		products = append(products, toProductEntity(model))
	}

	return products, total, nil
}

func (r *crmRepository) CreateSalesOrder(ctx context.Context, salesOrder *entity.SalesOrder) (*entity.SalesOrder, error) {
	salesOrderID, err := r.idGen.GenID(ctx)
	if err != nil {
		return nil, wrapOperateError(err, "generate sales order id failed")
	}

	model, err := toSalesOrderModel(salesOrder)
	if err != nil {
		return nil, err
	}
	model.SalesOrderID = salesOrderID
	var created *entity.SalesOrder
	if err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(model).Error; err != nil {
			return wrapOperateError(err, "create sales order failed")
		}

		created = toSalesOrderEntity(model)
		return r.appendAuditLog(
			ctx,
			tx,
			created.TenantID,
			created.SpaceID,
			crmAuditResourceSalesOrder,
			created.SalesOrderID,
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

func (r *crmRepository) UpdateSalesOrder(ctx context.Context, salesOrder *entity.SalesOrder) (*entity.SalesOrder, error) {
	orderDate, err := parseDate(salesOrder.OrderDate, salesOrder.SignTime)
	if err != nil {
		return nil, err
	}

	result := r.db.WithContext(ctx).
		Model(&crmSalesOrderModel{}).
		Where("id = ? AND tenant_id = ? AND space_id = ? AND is_deleted = ?", salesOrder.SalesOrderID, salesOrder.TenantID, salesOrder.SpaceID, false).
		Updates(map[string]any{
			"customer_id":     salesOrder.CustomerID,
			"opportunity_id":  nullableID(salesOrder.OpportunityID),
			"product_id":      salesOrder.ProductID,
			"product_name":    salesOrder.ProductName,
			"sales_user_id":   salesOrder.SalesUserID,
			"sales_user_name": salesOrder.SalesUserName,
			"quantity":        salesOrder.Quantity,
			"amount":          salesOrder.Amount,
			"order_date":      orderDate,
			"status":          salesOrder.Status,
			"remark":          salesOrder.Remark,
			"updated_by":      salesOrder.UpdatedBy,
			"updated_at":      timeNowMillis(),
		})
	if result.Error != nil {
		return nil, wrapOperateError(result.Error, "update sales order failed")
	}
	if result.RowsAffected == 0 {
		return nil, errorx.New(errno.ErrCRMRecordNotFoundCode)
	}

	return r.GetSalesOrder(ctx, &entity.Scope{TenantID: salesOrder.TenantID, SpaceID: salesOrder.SpaceID}, salesOrder.SalesOrderID)
}

func (r *crmRepository) DeleteSalesOrder(ctx context.Context, scope *entity.Scope, salesOrderID int64) error {
	return r.softDelete(ctx, &crmSalesOrderModel{}, scope, salesOrderID, "delete sales order failed")
}

func (r *crmRepository) GetSalesOrder(ctx context.Context, scope *entity.Scope, salesOrderID int64) (*entity.SalesOrder, error) {
	var model crmSalesOrderModel
	if err := r.scopedFirst(ctx, &model, scope, salesOrderID); err != nil {
		return nil, err
	}

	return toSalesOrderEntity(&model), nil
}

func (r *crmRepository) GetSalesOrderByID(ctx context.Context, salesOrderID int64) (*entity.SalesOrder, error) {
	var model crmSalesOrderModel
	if err := r.firstByID(ctx, &model, salesOrderID); err != nil {
		return nil, err
	}

	return toSalesOrderEntity(&model), nil
}

func (r *crmRepository) ListSalesOrders(ctx context.Context, filter *entity.SalesOrderFilter) ([]*entity.SalesOrder, int64, error) {
	query := r.scopedListQuery(ctx, &crmSalesOrderModel{}, &filter.Scope)
	if filter.CustomerID != nil && *filter.CustomerID > 0 {
		query = query.Where("customer_id = ?", *filter.CustomerID)
	}
	if filter.OpportunityID != nil && *filter.OpportunityID > 0 {
		query = query.Where("opportunity_id = ?", *filter.OpportunityID)
	}
	if filter.ProductID != nil && *filter.ProductID > 0 {
		query = query.Where("product_id = ?", *filter.ProductID)
	}
	if filter.SalesUserID != nil && *filter.SalesUserID > 0 {
		query = query.Where("sales_user_id = ?", *filter.SalesUserID)
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
	if filter.OrderDateStart != nil && *filter.OrderDateStart != "" {
		query = query.Where("order_date >= ?", *filter.OrderDateStart)
	}
	if filter.OrderDateEnd != nil && *filter.OrderDateEnd != "" {
		query = query.Where("order_date <= ?", *filter.OrderDateEnd)
	}
	if filter.Keyword != "" {
		like := keywordLike(filter.Keyword)
		query = query.Where("product_name LIKE ? OR sales_user_name LIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, wrapOperateError(err, "count sales order list failed")
	}

	var models []*crmSalesOrderModel
	if err := query.
		Order("updated_at DESC").
		Offset(pageOffset(filter.Page, filter.PageSize)).
		Limit(filter.PageSize).
		Find(&models).Error; err != nil {
		return nil, 0, wrapOperateError(err, "list sales order failed")
	}

	salesOrders := make([]*entity.SalesOrder, 0, len(models))
	for _, model := range models {
		salesOrders = append(salesOrders, toSalesOrderEntity(model))
	}

	return salesOrders, total, nil
}
