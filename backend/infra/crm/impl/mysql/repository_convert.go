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

import "github.com/coze-dev/coze-studio/backend/domain/crm/entity"

func toCustomerModel(customer *entity.Customer) *crmCustomerModel {
	customer.Normalize()
	return &crmCustomerModel{
		CustomerID:    customer.CustomerID,
		TenantID:      customer.TenantID,
		SpaceID:       customer.SpaceID,
		CustomerName:  customer.CustomerName,
		CustomerCode:  customer.CustomerCode,
		Industry:      customer.Industry,
		Level:         customer.Level,
		OwnerUserID:   customer.OwnerUserID,
		OwnerUserName: customer.OwnerUserName,
		Status:        customer.Status,
		Mobile:        customer.Mobile,
		Email:         customer.Email,
		Address:       customer.Address,
		Remark:        customer.Remark,
		CreatedBy:     customer.CreatedBy,
		UpdatedBy:     customer.UpdatedBy,
		CreatedAt:     customer.CreatedAt,
		UpdatedAt:     customer.UpdatedAt,
		IsDeleted:     customer.IsDeleted,
	}
}

func toCustomerEntity(model *crmCustomerModel) *entity.Customer {
	if model == nil {
		return nil
	}
	return &entity.Customer{
		CustomerID:    model.CustomerID,
		TenantID:      model.TenantID,
		SpaceID:       model.SpaceID,
		CustomerName:  model.CustomerName,
		CustomerCode:  model.CustomerCode,
		Industry:      model.Industry,
		Level:         model.Level,
		CustomerLevel: model.Level,
		OwnerUserID:   model.OwnerUserID,
		OwnerUserName: model.OwnerUserName,
		Status:        model.Status,
		Mobile:        model.Mobile,
		Email:         model.Email,
		Address:       model.Address,
		Remark:        model.Remark,
		Description:   model.Remark,
		AuditInfo: entity.AuditInfo{
			CreatedBy: model.CreatedBy,
			UpdatedBy: model.UpdatedBy,
			CreatorID: model.CreatedBy,
			UpdaterID: model.UpdatedBy,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			IsDeleted: model.IsDeleted,
		},
	}
}

func toContactModel(contact *entity.Contact) *crmContactModel {
	contact.Normalize()
	return &crmContactModel{
		ContactID:   contact.ContactID,
		TenantID:    contact.TenantID,
		SpaceID:     contact.SpaceID,
		CustomerID:  contact.CustomerID,
		ContactName: contact.ContactName,
		Mobile:      contact.Mobile,
		Email:       contact.Email,
		Title:       contact.Title,
		IsPrimary:   contact.IsPrimary,
		Status:      contact.Status,
		Remark:      contact.Remark,
		CreatedBy:   contact.CreatedBy,
		UpdatedBy:   contact.UpdatedBy,
		CreatedAt:   contact.CreatedAt,
		UpdatedAt:   contact.UpdatedAt,
		IsDeleted:   contact.IsDeleted,
	}
}

func toContactEntity(model *crmContactModel) *entity.Contact {
	if model == nil {
		return nil
	}
	return &entity.Contact{
		ContactID:   model.ContactID,
		TenantID:    model.TenantID,
		SpaceID:     model.SpaceID,
		CustomerID:  model.CustomerID,
		ContactName: model.ContactName,
		Mobile:      model.Mobile,
		Email:       model.Email,
		Title:       model.Title,
		Position:    model.Title,
		IsPrimary:   model.IsPrimary,
		Status:      model.Status,
		Remark:      model.Remark,
		Description: model.Remark,
		AuditInfo: entity.AuditInfo{
			CreatedBy: model.CreatedBy,
			UpdatedBy: model.UpdatedBy,
			CreatorID: model.CreatedBy,
			UpdaterID: model.UpdatedBy,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			IsDeleted: model.IsDeleted,
		},
	}
}

func toOpportunityModel(opportunity *entity.Opportunity) (*crmOpportunityModel, error) {
	opportunity.Normalize()
	expectedCloseDate, err := parseDate(opportunity.ExpectedCloseDate, opportunity.ExpectedCloseTime)
	if err != nil {
		return nil, err
	}
	return &crmOpportunityModel{
		OpportunityID:     opportunity.OpportunityID,
		TenantID:          opportunity.TenantID,
		SpaceID:           opportunity.SpaceID,
		CustomerID:        opportunity.CustomerID,
		OpportunityName:   opportunity.OpportunityName,
		Stage:             opportunity.Stage,
		Amount:            opportunity.Amount,
		ExpectedCloseDate: expectedCloseDate,
		OwnerUserID:       opportunity.OwnerUserID,
		OwnerUserName:     opportunity.OwnerUserName,
		Status:            opportunity.Status,
		Remark:            opportunity.Remark,
		CreatedBy:         opportunity.CreatedBy,
		UpdatedBy:         opportunity.UpdatedBy,
		CreatedAt:         opportunity.CreatedAt,
		UpdatedAt:         opportunity.UpdatedAt,
		IsDeleted:         opportunity.IsDeleted,
	}, nil
}

func toOpportunityEntity(model *crmOpportunityModel) *entity.Opportunity {
	if model == nil {
		return nil
	}
	return &entity.Opportunity{
		OpportunityID:     model.OpportunityID,
		TenantID:          model.TenantID,
		SpaceID:           model.SpaceID,
		CustomerID:        model.CustomerID,
		OpportunityName:   model.OpportunityName,
		Stage:             model.Stage,
		Amount:            model.Amount,
		ExpectedCloseDate: formatDate(model.ExpectedCloseDate),
		ExpectedCloseTime: dateToMillis(model.ExpectedCloseDate),
		OwnerUserID:       model.OwnerUserID,
		OwnerUserName:     model.OwnerUserName,
		Status:            model.Status,
		Remark:            model.Remark,
		Description:       model.Remark,
		AuditInfo: entity.AuditInfo{
			CreatedBy: model.CreatedBy,
			UpdatedBy: model.UpdatedBy,
			CreatorID: model.CreatedBy,
			UpdaterID: model.UpdatedBy,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			IsDeleted: model.IsDeleted,
		},
	}
}

func toFollowRecordModel(followRecord *entity.FollowRecord) *crmFollowRecordModel {
	followRecord.Normalize()
	return &crmFollowRecordModel{
		FollowRecordID: followRecord.FollowRecordID,
		TenantID:       followRecord.TenantID,
		SpaceID:        followRecord.SpaceID,
		CustomerID:     followRecord.CustomerID,
		ContactID:      followRecord.ContactID,
		FollowType:     followRecord.FollowType,
		Content:        followRecord.Content,
		NextFollowTime: nullableDateTime(followRecord.NextFollowTime),
		OwnerUserID:    followRecord.OwnerUserID,
		OwnerUserName:  followRecord.OwnerUserName,
		Status:         followRecord.Status,
		CreatedBy:      followRecord.CreatedBy,
		UpdatedBy:      followRecord.UpdatedBy,
		CreatedAt:      followRecord.CreatedAt,
		UpdatedAt:      followRecord.UpdatedAt,
		IsDeleted:      followRecord.IsDeleted,
	}
}

func toFollowRecordEntity(model *crmFollowRecordModel) *entity.FollowRecord {
	if model == nil {
		return nil
	}
	return &entity.FollowRecord{
		FollowRecordID: model.FollowRecordID,
		TenantID:       model.TenantID,
		SpaceID:        model.SpaceID,
		CustomerID:     model.CustomerID,
		ContactID:      model.ContactID,
		FollowType:     model.FollowType,
		Content:        model.Content,
		FollowContent:  model.Content,
		NextFollowTime: dateTimeToMillis(model.NextFollowTime),
		OwnerUserID:    model.OwnerUserID,
		OwnerUserName:  model.OwnerUserName,
		Status:         model.Status,
		AuditInfo: entity.AuditInfo{
			CreatedBy: model.CreatedBy,
			UpdatedBy: model.UpdatedBy,
			CreatorID: model.CreatedBy,
			UpdaterID: model.UpdatedBy,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			IsDeleted: model.IsDeleted,
		},
	}
}

func toProductModel(product *entity.Product) *crmProductModel {
	product.Normalize()
	return &crmProductModel{
		ProductID:   product.ProductID,
		TenantID:    product.TenantID,
		SpaceID:     product.SpaceID,
		ProductName: product.ProductName,
		ProductCode: product.ProductCode,
		Category:    product.Category,
		UnitPrice:   product.UnitPrice,
		Status:      product.Status,
		Remark:      product.Remark,
		CreatedBy:   product.CreatedBy,
		UpdatedBy:   product.UpdatedBy,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		IsDeleted:   product.IsDeleted,
	}
}

func toProductEntity(model *crmProductModel) *entity.Product {
	if model == nil {
		return nil
	}
	return &entity.Product{
		ProductID:   model.ProductID,
		TenantID:    model.TenantID,
		SpaceID:     model.SpaceID,
		ProductName: model.ProductName,
		ProductCode: model.ProductCode,
		Category:    model.Category,
		UnitPrice:   model.UnitPrice,
		Status:      model.Status,
		Remark:      model.Remark,
		Description: model.Remark,
		AuditInfo: entity.AuditInfo{
			CreatedBy: model.CreatedBy,
			UpdatedBy: model.UpdatedBy,
			CreatorID: model.CreatedBy,
			UpdaterID: model.UpdatedBy,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			IsDeleted: model.IsDeleted,
		},
	}
}

func toSalesOrderModel(salesOrder *entity.SalesOrder) (*crmSalesOrderModel, error) {
	salesOrder.Normalize()
	orderDate, err := parseDate(salesOrder.OrderDate, salesOrder.SignTime)
	if err != nil {
		return nil, err
	}
	return &crmSalesOrderModel{
		SalesOrderID:  salesOrder.SalesOrderID,
		TenantID:      salesOrder.TenantID,
		SpaceID:       salesOrder.SpaceID,
		CustomerID:    salesOrder.CustomerID,
		OpportunityID: salesOrder.OpportunityID,
		ProductID:     salesOrder.ProductID,
		ProductName:   salesOrder.ProductName,
		SalesUserID:   salesOrder.SalesUserID,
		SalesUserName: salesOrder.SalesUserName,
		Quantity:      salesOrder.Quantity,
		Amount:        salesOrder.Amount,
		OrderDate:     orderDate,
		Status:        salesOrder.Status,
		Remark:        salesOrder.Remark,
		CreatedBy:     salesOrder.CreatedBy,
		UpdatedBy:     salesOrder.UpdatedBy,
		CreatedAt:     salesOrder.CreatedAt,
		UpdatedAt:     salesOrder.UpdatedAt,
		IsDeleted:     salesOrder.IsDeleted,
	}, nil
}

func toSalesOrderEntity(model *crmSalesOrderModel) *entity.SalesOrder {
	if model == nil {
		return nil
	}
	return &entity.SalesOrder{
		SalesOrderID:   model.SalesOrderID,
		TenantID:       model.TenantID,
		SpaceID:        model.SpaceID,
		CustomerID:     model.CustomerID,
		OpportunityID:  model.OpportunityID,
		ProductID:      model.ProductID,
		ProductName:    model.ProductName,
		ProductSummary: model.ProductName,
		SalesUserID:    model.SalesUserID,
		SalesUserName:  model.SalesUserName,
		Quantity:       model.Quantity,
		Amount:         model.Amount,
		TotalAmount:    model.Amount,
		OrderDate:      formatDate(model.OrderDate),
		SignTime:       dateToMillis(model.OrderDate),
		Status:         model.Status,
		Remark:         model.Remark,
		Description:    model.Remark,
		AuditInfo: entity.AuditInfo{
			CreatedBy: model.CreatedBy,
			UpdatedBy: model.UpdatedBy,
			CreatorID: model.CreatedBy,
			UpdaterID: model.UpdatedBy,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
			IsDeleted: model.IsDeleted,
		},
	}
}
