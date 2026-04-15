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

import "time"

const (
	tableNameCRMCustomer     = "crm_customer"
	tableNameCRMAuditLog     = "crm_audit_log"
	tableNameCRMContact      = "crm_contact"
	tableNameCRMOpportunity  = "crm_opportunity"
	tableNameCRMFollowRecord = "crm_follow_record"
	tableNameCRMProduct      = "crm_product"
	tableNameCRMSalesOrder   = "crm_sales_order"
)

type crmAuditLogModel struct {
	AuditLogID     int64  `gorm:"column:id;primaryKey;comment:Audit Log ID" json:"id"`
	TenantID       int64  `gorm:"column:tenant_id;not null;comment:Tenant ID" json:"tenant_id"`
	SpaceID        int64  `gorm:"column:space_id;not null;comment:Space ID" json:"space_id"`
	ResourceType   string `gorm:"column:resource_type;not null;comment:Resource Type" json:"resource_type"`
	ResourceID     int64  `gorm:"column:resource_id;not null;comment:Resource ID" json:"resource_id"`
	Action         string `gorm:"column:action;not null;comment:Action Type" json:"action"`
	OperatorID     int64  `gorm:"column:operator_id;not null;comment:Operator User ID" json:"operator_id"`
	BeforeSnapshot string `gorm:"column:before_snapshot;type:mediumtext;comment:Before Snapshot JSON" json:"before_snapshot"`
	AfterSnapshot  string `gorm:"column:after_snapshot;type:mediumtext;comment:After Snapshot JSON" json:"after_snapshot"`
	OperationAt    int64  `gorm:"column:operation_at;not null;comment:Operation Time in Milliseconds" json:"operation_at"`
}

func (*crmAuditLogModel) TableName() string {
	return tableNameCRMAuditLog
}

type crmCustomerModel struct {
	CustomerID    int64  `gorm:"column:id;primaryKey;comment:Customer ID" json:"id"`
	TenantID      int64  `gorm:"column:tenant_id;not null;comment:Tenant ID" json:"tenant_id"`
	SpaceID       int64  `gorm:"column:space_id;not null;comment:Space ID" json:"space_id"`
	CustomerName  string `gorm:"column:customer_name;not null;comment:Customer Name" json:"customer_name"`
	CustomerCode  string `gorm:"column:customer_code;comment:Customer Code" json:"customer_code"`
	Industry      string `gorm:"column:industry;comment:Industry" json:"industry"`
	Level         string `gorm:"column:level;comment:Customer Level" json:"level"`
	OwnerUserID   int64  `gorm:"column:owner_user_id;not null;comment:Owner User ID" json:"owner_user_id"`
	OwnerUserName string `gorm:"column:owner_user_name;not null;comment:Owner User Name" json:"owner_user_name"`
	Status        string `gorm:"column:status;not null;comment:Status" json:"status"`
	Mobile        string `gorm:"column:mobile;comment:Mobile" json:"mobile"`
	Email         string `gorm:"column:email;comment:Email" json:"email"`
	Address       string `gorm:"column:address;comment:Address" json:"address"`
	Remark        string `gorm:"column:remark;type:text;comment:Remark" json:"remark"`
	CreatedBy     int64  `gorm:"column:created_by;not null;comment:Created By" json:"created_by"`
	UpdatedBy     int64  `gorm:"column:updated_by;not null;comment:Updated By" json:"updated_by"`
	CreatedAt     int64  `gorm:"column:created_at;not null;autoCreateTime:milli;comment:Create Time in Milliseconds" json:"created_at"`
	UpdatedAt     int64  `gorm:"column:updated_at;not null;autoUpdateTime:milli;comment:Update Time in Milliseconds" json:"updated_at"`
	IsDeleted     bool   `gorm:"column:is_deleted;not null;comment:Soft Delete Flag" json:"is_deleted"`
}

func (*crmCustomerModel) TableName() string {
	return tableNameCRMCustomer
}

type crmContactModel struct {
	ContactID   int64  `gorm:"column:id;primaryKey;comment:Contact ID" json:"id"`
	TenantID    int64  `gorm:"column:tenant_id;not null;comment:Tenant ID" json:"tenant_id"`
	SpaceID     int64  `gorm:"column:space_id;not null;comment:Space ID" json:"space_id"`
	CustomerID  int64  `gorm:"column:customer_id;not null;comment:Customer ID" json:"customer_id"`
	ContactName string `gorm:"column:contact_name;not null;comment:Contact Name" json:"contact_name"`
	Mobile      string `gorm:"column:mobile;comment:Mobile" json:"mobile"`
	Email       string `gorm:"column:email;comment:Email" json:"email"`
	Title       string `gorm:"column:title;comment:Job Title" json:"title"`
	IsPrimary   bool   `gorm:"column:is_primary;not null;comment:Is Primary Contact" json:"is_primary"`
	Status      string `gorm:"column:status;not null;comment:Status" json:"status"`
	Remark      string `gorm:"column:remark;type:text;comment:Remark" json:"remark"`
	CreatedBy   int64  `gorm:"column:created_by;not null;comment:Created By" json:"created_by"`
	UpdatedBy   int64  `gorm:"column:updated_by;not null;comment:Updated By" json:"updated_by"`
	CreatedAt   int64  `gorm:"column:created_at;not null;autoCreateTime:milli;comment:Create Time in Milliseconds" json:"created_at"`
	UpdatedAt   int64  `gorm:"column:updated_at;not null;autoUpdateTime:milli;comment:Update Time in Milliseconds" json:"updated_at"`
	IsDeleted   bool   `gorm:"column:is_deleted;not null;comment:Soft Delete Flag" json:"is_deleted"`
}

func (*crmContactModel) TableName() string {
	return tableNameCRMContact
}

type crmOpportunityModel struct {
	OpportunityID     int64      `gorm:"column:id;primaryKey;comment:Opportunity ID" json:"id"`
	TenantID          int64      `gorm:"column:tenant_id;not null;comment:Tenant ID" json:"tenant_id"`
	SpaceID           int64      `gorm:"column:space_id;not null;comment:Space ID" json:"space_id"`
	CustomerID        int64      `gorm:"column:customer_id;not null;comment:Customer ID" json:"customer_id"`
	OpportunityName   string     `gorm:"column:opportunity_name;not null;comment:Opportunity Name" json:"opportunity_name"`
	Stage             string     `gorm:"column:stage;not null;comment:Opportunity Stage" json:"stage"`
	Amount            float64    `gorm:"column:amount;type:decimal(18,2);not null;comment:Estimated Amount" json:"amount"`
	ExpectedCloseDate *time.Time `gorm:"column:expected_close_date;type:date;comment:Expected Close Date" json:"expected_close_date"`
	OwnerUserID       int64      `gorm:"column:owner_user_id;not null;comment:Owner User ID" json:"owner_user_id"`
	OwnerUserName     string     `gorm:"column:owner_user_name;not null;comment:Owner User Name" json:"owner_user_name"`
	Status            string     `gorm:"column:status;not null;comment:Status" json:"status"`
	Remark            string     `gorm:"column:remark;type:text;comment:Remark" json:"remark"`
	CreatedBy         int64      `gorm:"column:created_by;not null;comment:Created By" json:"created_by"`
	UpdatedBy         int64      `gorm:"column:updated_by;not null;comment:Updated By" json:"updated_by"`
	CreatedAt         int64      `gorm:"column:created_at;not null;autoCreateTime:milli;comment:Create Time in Milliseconds" json:"created_at"`
	UpdatedAt         int64      `gorm:"column:updated_at;not null;autoUpdateTime:milli;comment:Update Time in Milliseconds" json:"updated_at"`
	IsDeleted         bool       `gorm:"column:is_deleted;not null;comment:Soft Delete Flag" json:"is_deleted"`
}

func (*crmOpportunityModel) TableName() string {
	return tableNameCRMOpportunity
}

type crmFollowRecordModel struct {
	FollowRecordID int64      `gorm:"column:id;primaryKey;comment:Follow Record ID" json:"id"`
	TenantID       int64      `gorm:"column:tenant_id;not null;comment:Tenant ID" json:"tenant_id"`
	SpaceID        int64      `gorm:"column:space_id;not null;comment:Space ID" json:"space_id"`
	CustomerID     int64      `gorm:"column:customer_id;not null;comment:Customer ID" json:"customer_id"`
	ContactID      int64      `gorm:"column:contact_id;comment:Contact ID" json:"contact_id"`
	FollowType     string     `gorm:"column:follow_type;not null;comment:Follow Type" json:"follow_type"`
	Content        string     `gorm:"column:content;type:text;not null;comment:Follow Content" json:"content"`
	NextFollowTime *time.Time `gorm:"column:next_follow_time;type:datetime(3);comment:Next Follow Time" json:"next_follow_time"`
	OwnerUserID    int64      `gorm:"column:owner_user_id;not null;comment:Owner User ID" json:"owner_user_id"`
	OwnerUserName  string     `gorm:"column:owner_user_name;not null;comment:Owner User Name" json:"owner_user_name"`
	Status         string     `gorm:"column:status;not null;comment:Status" json:"status"`
	CreatedBy      int64      `gorm:"column:created_by;not null;comment:Created By" json:"created_by"`
	UpdatedBy      int64      `gorm:"column:updated_by;not null;comment:Updated By" json:"updated_by"`
	CreatedAt      int64      `gorm:"column:created_at;not null;autoCreateTime:milli;comment:Create Time in Milliseconds" json:"created_at"`
	UpdatedAt      int64      `gorm:"column:updated_at;not null;autoUpdateTime:milli;comment:Update Time in Milliseconds" json:"updated_at"`
	IsDeleted      bool       `gorm:"column:is_deleted;not null;comment:Soft Delete Flag" json:"is_deleted"`
}

func (*crmFollowRecordModel) TableName() string {
	return tableNameCRMFollowRecord
}

type crmProductModel struct {
	ProductID   int64   `gorm:"column:id;primaryKey;comment:Product ID" json:"id"`
	TenantID    int64   `gorm:"column:tenant_id;not null;comment:Tenant ID" json:"tenant_id"`
	SpaceID     int64   `gorm:"column:space_id;not null;comment:Space ID" json:"space_id"`
	ProductName string  `gorm:"column:product_name;not null;comment:Product Name" json:"product_name"`
	ProductCode string  `gorm:"column:product_code;comment:Product Code" json:"product_code"`
	Category    string  `gorm:"column:category;comment:Category" json:"category"`
	UnitPrice   float64 `gorm:"column:unit_price;type:decimal(18,2);not null;comment:Unit Price" json:"unit_price"`
	Status      string  `gorm:"column:status;not null;comment:Status" json:"status"`
	Remark      string  `gorm:"column:remark;type:text;comment:Remark" json:"remark"`
	CreatedBy   int64   `gorm:"column:created_by;not null;comment:Created By" json:"created_by"`
	UpdatedBy   int64   `gorm:"column:updated_by;not null;comment:Updated By" json:"updated_by"`
	CreatedAt   int64   `gorm:"column:created_at;not null;autoCreateTime:milli;comment:Create Time in Milliseconds" json:"created_at"`
	UpdatedAt   int64   `gorm:"column:updated_at;not null;autoUpdateTime:milli;comment:Update Time in Milliseconds" json:"updated_at"`
	IsDeleted   bool    `gorm:"column:is_deleted;not null;comment:Soft Delete Flag" json:"is_deleted"`
}

func (*crmProductModel) TableName() string {
	return tableNameCRMProduct
}

type crmSalesOrderModel struct {
	SalesOrderID  int64      `gorm:"column:id;primaryKey;comment:Sales Order ID" json:"id"`
	TenantID      int64      `gorm:"column:tenant_id;not null;comment:Tenant ID" json:"tenant_id"`
	SpaceID       int64      `gorm:"column:space_id;not null;comment:Space ID" json:"space_id"`
	CustomerID    int64      `gorm:"column:customer_id;not null;comment:Customer ID" json:"customer_id"`
	OpportunityID int64      `gorm:"column:opportunity_id;comment:Opportunity ID" json:"opportunity_id"`
	ProductID     int64      `gorm:"column:product_id;not null;comment:Product ID" json:"product_id"`
	ProductName   string     `gorm:"column:product_name;not null;comment:Product Snapshot Name" json:"product_name"`
	SalesUserID   int64      `gorm:"column:sales_user_id;not null;comment:Sales User ID" json:"sales_user_id"`
	SalesUserName string     `gorm:"column:sales_user_name;not null;comment:Sales User Name" json:"sales_user_name"`
	Quantity      float64    `gorm:"column:quantity;type:decimal(18,2);not null;comment:Sales Quantity" json:"quantity"`
	Amount        float64    `gorm:"column:amount;type:decimal(18,2);not null;comment:Sales Amount" json:"amount"`
	OrderDate     *time.Time `gorm:"column:order_date;type:date;comment:Order Date" json:"order_date"`
	Status        string     `gorm:"column:status;not null;comment:Status" json:"status"`
	Remark        string     `gorm:"column:remark;type:text;comment:Remark" json:"remark"`
	CreatedBy     int64      `gorm:"column:created_by;not null;comment:Created By" json:"created_by"`
	UpdatedBy     int64      `gorm:"column:updated_by;not null;comment:Updated By" json:"updated_by"`
	CreatedAt     int64      `gorm:"column:created_at;not null;autoCreateTime:milli;comment:Create Time in Milliseconds" json:"created_at"`
	UpdatedAt     int64      `gorm:"column:updated_at;not null;autoUpdateTime:milli;comment:Update Time in Milliseconds" json:"updated_at"`
	IsDeleted     bool       `gorm:"column:is_deleted;not null;comment:Soft Delete Flag" json:"is_deleted"`
}

func (*crmSalesOrderModel) TableName() string {
	return tableNameCRMSalesOrder
}
