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

func ListProducts(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.ProductListQuery
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}

	list, total, err := crmapp.CRMSVC.ListProducts(ctx, &entity.ProductFilter{
		Scope: entity.Scope{SpaceID: spaceID},
		PageOption: entity.PageOption{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
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

	resp := &crmmodel.ProductListData{
		List:  make([]*crmmodel.ProductData, 0, len(list)),
		Total: total,
	}
	for _, item := range list {
		resp.List = append(resp.List, toProductData(item))
	}

	writeCRMSuccess(c, resp)
}

func GetProduct(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.GetProductRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}
	productID, ok := parseRequiredInt64Param(c, req.ProductID, "product_id")
	if !ok {
		return
	}

	product, err := crmapp.CRMSVC.GetProduct(ctx, spaceID, productID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toProductData(product))
}

func CreateProduct(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.CreateProductRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}
	unitPrice, ok := parseDecimalParam(c, req.UnitPrice, "unit_price", false)
	if !ok {
		return
	}

	product, err := crmapp.CRMSVC.CreateProduct(ctx, &entity.Product{
		SpaceID:     spaceID,
		ProductName: req.ProductName,
		ProductCode: req.ProductCode,
		Category:    req.Category,
		UnitPrice:   unitPrice,
		Status:      req.Status,
		Remark:      req.Remark,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toProductData(product))
}

func UpdateProduct(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.UpdateProductRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}
	productID, ok := parseRequiredInt64Param(c, req.ProductID, "product_id")
	if !ok {
		return
	}
	unitPrice, ok := parseDecimalParam(c, req.UnitPrice, "unit_price", false)
	if !ok {
		return
	}

	product, err := crmapp.CRMSVC.UpdateProduct(ctx, &entity.Product{
		ProductID:   productID,
		SpaceID:     spaceID,
		ProductName: req.ProductName,
		ProductCode: req.ProductCode,
		Category:    req.Category,
		UnitPrice:   unitPrice,
		Status:      req.Status,
		Remark:      req.Remark,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toProductData(product))
}

func DeleteProduct(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.DeleteProductRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}
	productID, ok := parseRequiredInt64Param(c, req.ProductID, "product_id")
	if !ok {
		return
	}

	if err := crmapp.CRMSVC.DeleteProduct(ctx, spaceID, productID); err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, map[string]any{})
}

func ListSalesOrders(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.SalesOrderListQuery
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}

	list, total, err := crmapp.CRMSVC.ListSalesOrders(ctx, &entity.SalesOrderFilter{
		Scope: entity.Scope{SpaceID: spaceID},
		PageOption: entity.PageOption{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		CustomerID:     parseOptionalInt64Param(c, req.CustomerID, "customer_id"),
		OpportunityID:  parseOptionalInt64Param(c, req.OpportunityID, "opportunity_id"),
		ProductID:      parseOptionalInt64Param(c, req.ProductID, "product_id"),
		SalesUserID:    parseOptionalInt64Param(c, req.SalesUserID, "sales_user_id"),
		Keyword:        strings.TrimSpace(req.Keyword),
		Status:         parseOptionalString(req.Status),
		CreatedAtStart: parseOptionalTimestampParam(c, req.CreatedAtStart, "created_at_start"),
		CreatedAtEnd:   parseOptionalTimestampParam(c, req.CreatedAtEnd, "created_at_end"),
		OrderDateStart: parseOptionalDate(req.OrderDateStart),
		OrderDateEnd:   parseOptionalDate(req.OrderDateEnd),
	})
	if c.IsAborted() {
		return
	}
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	resp := &crmmodel.SalesOrderListData{
		List:  make([]*crmmodel.SalesOrderData, 0, len(list)),
		Total: total,
	}
	for _, item := range list {
		resp.List = append(resp.List, toSalesOrderData(item))
	}

	writeCRMSuccess(c, resp)
}

func GetSalesOrder(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.GetSalesOrderRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}
	salesOrderID, ok := parseRequiredInt64Param(c, req.SalesOrderID, "sales_order_id")
	if !ok {
		return
	}

	order, err := crmapp.CRMSVC.GetSalesOrder(ctx, spaceID, salesOrderID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toSalesOrderData(order))
}

func CreateSalesOrder(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.CreateSalesOrderRequest
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
	productID, ok := parseRequiredInt64Param(c, req.ProductID, "product_id")
	if !ok {
		return
	}
	quantity, ok := parseDecimalParam(c, req.Quantity, "quantity", false)
	if !ok {
		return
	}
	amount, ok := parseDecimalParam(c, req.Amount, "amount", false)
	if !ok {
		return
	}

	order, err := crmapp.CRMSVC.CreateSalesOrder(ctx, &entity.SalesOrder{
		SpaceID:       spaceID,
		CustomerID:    customerID,
		OpportunityID: parseOptionalInt64Value(req.OpportunityID),
		ProductID:     productID,
		ProductName:   req.ProductName,
		SalesUserID:   parseOptionalInt64Value(req.SalesUserID),
		SalesUserName: req.SalesUserName,
		Quantity:      quantity,
		Amount:        amount,
		OrderDate:     req.OrderDate,
		Status:        req.Status,
		Remark:        req.Remark,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toSalesOrderData(order))
}

func UpdateSalesOrder(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.UpdateSalesOrderRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}
	salesOrderID, ok := parseRequiredInt64Param(c, req.SalesOrderID, "sales_order_id")
	if !ok {
		return
	}
	customerID, ok := parseRequiredInt64Param(c, req.CustomerID, "customer_id")
	if !ok {
		return
	}
	productID, ok := parseRequiredInt64Param(c, req.ProductID, "product_id")
	if !ok {
		return
	}
	quantity, ok := parseDecimalParam(c, req.Quantity, "quantity", false)
	if !ok {
		return
	}
	amount, ok := parseDecimalParam(c, req.Amount, "amount", false)
	if !ok {
		return
	}

	order, err := crmapp.CRMSVC.UpdateSalesOrder(ctx, &entity.SalesOrder{
		SalesOrderID:  salesOrderID,
		SpaceID:       spaceID,
		CustomerID:    customerID,
		OpportunityID: parseOptionalInt64Value(req.OpportunityID),
		ProductID:     productID,
		ProductName:   req.ProductName,
		SalesUserID:   parseOptionalInt64Value(req.SalesUserID),
		SalesUserName: req.SalesUserName,
		Quantity:      quantity,
		Amount:        amount,
		OrderDate:     req.OrderDate,
		Status:        req.Status,
		Remark:        req.Remark,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, toSalesOrderData(order))
}

func DeleteSalesOrder(ctx context.Context, c *app.RequestContext) {
	var req crmmodel.DeleteSalesOrderRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	spaceID, ok := parseRequiredInt64Param(c, req.SpaceID, "space_id")
	if !ok {
		return
	}
	salesOrderID, ok := parseRequiredInt64Param(c, req.SalesOrderID, "sales_order_id")
	if !ok {
		return
	}

	if err := crmapp.CRMSVC.DeleteSalesOrder(ctx, spaceID, salesOrderID); err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	writeCRMSuccess(c, map[string]any{})
}

func toProductData(product *entity.Product) *crmmodel.ProductData {
	if product == nil {
		return nil
	}
	return &crmmodel.ProductData{
		ProductID:   formatCRMInt64(product.ProductID),
		TenantID:    formatCRMInt64(product.TenantID),
		SpaceID:     formatCRMInt64(product.SpaceID),
		ProductName: product.ProductName,
		ProductCode: product.ProductCode,
		Category:    product.Category,
		UnitPrice:   formatCRMFloat(product.UnitPrice),
		Status:      product.Status,
		Remark:      product.Remark,
		CreatedBy:   formatCRMInt64(product.CreatedBy),
		UpdatedBy:   formatCRMInt64(product.UpdatedBy),
		CreatedAt:   formatCRMInt64(product.CreatedAt),
		UpdatedAt:   formatCRMInt64(product.UpdatedAt),
	}
}

func toSalesOrderData(order *entity.SalesOrder) *crmmodel.SalesOrderData {
	if order == nil {
		return nil
	}
	return &crmmodel.SalesOrderData{
		SalesOrderID:  formatCRMInt64(order.SalesOrderID),
		TenantID:      formatCRMInt64(order.TenantID),
		SpaceID:       formatCRMInt64(order.SpaceID),
		CustomerID:    formatCRMInt64(order.CustomerID),
		OpportunityID: formatCRMInt64(order.OpportunityID),
		ProductID:     formatCRMInt64(order.ProductID),
		ProductName:   order.ProductName,
		SalesUserID:   formatCRMInt64(order.SalesUserID),
		SalesUserName: order.SalesUserName,
		Quantity:      formatCRMFloat(order.Quantity),
		Amount:        formatCRMFloat(order.Amount),
		OrderDate:     order.OrderDate,
		Status:        order.Status,
		Remark:        order.Remark,
		CreatedBy:     formatCRMInt64(order.CreatedBy),
		UpdatedBy:     formatCRMInt64(order.UpdatedBy),
		CreatedAt:     formatCRMInt64(order.CreatedAt),
		UpdatedAt:     formatCRMInt64(order.UpdatedAt),
	}
}
