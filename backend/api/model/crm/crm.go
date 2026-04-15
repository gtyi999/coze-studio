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

package crm

type PageQuery struct {
	Page     int `query:"page" json:"page"`
	PageSize int `query:"page_size" json:"page_size"`
}

type TimeRangeQuery struct {
	CreatedAtStart string `query:"created_at_start"`
	CreatedAtEnd   string `query:"created_at_end"`
}

type CustomerListQuery struct {
	SpaceID     string `query:"space_id"`
	Keyword     string `query:"keyword"`
	Status      string `query:"status"`
	OwnerUserID string `query:"owner_user_id"`
	TimeRangeQuery
	PageQuery
}

type ContactListQuery struct {
	SpaceID    string `query:"space_id"`
	CustomerID string `query:"customer_id"`
	Keyword    string `query:"keyword"`
	Status     string `query:"status"`
	TimeRangeQuery
	PageQuery
}

type OpportunityListQuery struct {
	SpaceID                string `query:"space_id"`
	CustomerID             string `query:"customer_id"`
	OwnerUserID            string `query:"owner_user_id"`
	Keyword                string `query:"keyword"`
	Status                 string `query:"status"`
	ExpectedCloseDateStart string `query:"expected_close_date_start"`
	ExpectedCloseDateEnd   string `query:"expected_close_date_end"`
	TimeRangeQuery
	PageQuery
}

type FollowRecordListQuery struct {
	SpaceID             string `query:"space_id"`
	CustomerID          string `query:"customer_id"`
	ContactID           string `query:"contact_id"`
	OwnerUserID         string `query:"owner_user_id"`
	Keyword             string `query:"keyword"`
	Status              string `query:"status"`
	NextFollowTimeStart string `query:"next_follow_time_start"`
	NextFollowTimeEnd   string `query:"next_follow_time_end"`
	TimeRangeQuery
	PageQuery
}

type ProductListQuery struct {
	SpaceID string `query:"space_id"`
	Keyword string `query:"keyword"`
	Status  string `query:"status"`
	TimeRangeQuery
	PageQuery
}

type SalesOrderListQuery struct {
	SpaceID        string `query:"space_id"`
	CustomerID     string `query:"customer_id"`
	OpportunityID  string `query:"opportunity_id"`
	ProductID      string `query:"product_id"`
	SalesUserID    string `query:"sales_user_id"`
	Keyword        string `query:"keyword"`
	Status         string `query:"status"`
	OrderDateStart string `query:"order_date_start"`
	OrderDateEnd   string `query:"order_date_end"`
	TimeRangeQuery
	PageQuery
}

type DashboardOverviewQuery struct {
	SpaceID string `query:"space_id"`
}

type GetCustomerRequest struct {
	SpaceID    string `query:"space_id"`
	CustomerID string `query:"customer_id"`
}

type CreateCustomerRequest struct {
	SpaceID       string `json:"space_id"`
	CustomerName  string `json:"customer_name"`
	CustomerCode  string `json:"customer_code"`
	Industry      string `json:"industry"`
	Level         string `json:"level"`
	OwnerUserID   string `json:"owner_user_id"`
	OwnerUserName string `json:"owner_user_name"`
	Status        string `json:"status"`
	Mobile        string `json:"mobile"`
	Email         string `json:"email"`
	Address       string `json:"address"`
	Remark        string `json:"remark"`
}

type UpdateCustomerRequest struct {
	SpaceID       string `json:"space_id"`
	CustomerID    string `json:"customer_id"`
	CustomerName  string `json:"customer_name"`
	CustomerCode  string `json:"customer_code"`
	Industry      string `json:"industry"`
	Level         string `json:"level"`
	OwnerUserID   string `json:"owner_user_id"`
	OwnerUserName string `json:"owner_user_name"`
	Status        string `json:"status"`
	Mobile        string `json:"mobile"`
	Email         string `json:"email"`
	Address       string `json:"address"`
	Remark        string `json:"remark"`
}

type DeleteCustomerRequest struct {
	SpaceID    string `json:"space_id"`
	CustomerID string `json:"customer_id"`
}

type CustomerData struct {
	CustomerID    string `json:"customer_id,omitempty"`
	TenantID      string `json:"tenant_id,omitempty"`
	SpaceID       string `json:"space_id,omitempty"`
	CustomerName  string `json:"customer_name,omitempty"`
	CustomerCode  string `json:"customer_code,omitempty"`
	Industry      string `json:"industry,omitempty"`
	Level         string `json:"level,omitempty"`
	OwnerUserID   string `json:"owner_user_id,omitempty"`
	OwnerUserName string `json:"owner_user_name,omitempty"`
	Status        string `json:"status,omitempty"`
	Mobile        string `json:"mobile,omitempty"`
	Email         string `json:"email,omitempty"`
	Address       string `json:"address,omitempty"`
	Remark        string `json:"remark,omitempty"`
	CreatedBy     string `json:"created_by,omitempty"`
	UpdatedBy     string `json:"updated_by,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}

type CustomerListData struct {
	List  []*CustomerData `json:"list"`
	Total int64           `json:"total"`
}

type GetContactRequest struct {
	SpaceID   string `query:"space_id"`
	ContactID string `query:"contact_id"`
}

type CreateContactRequest struct {
	SpaceID     string `json:"space_id"`
	CustomerID  string `json:"customer_id"`
	ContactName string `json:"contact_name"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	Title       string `json:"title"`
	IsPrimary   bool   `json:"is_primary"`
	Status      string `json:"status"`
	Remark      string `json:"remark"`
}

type UpdateContactRequest struct {
	SpaceID     string `json:"space_id"`
	ContactID   string `json:"contact_id"`
	CustomerID  string `json:"customer_id"`
	ContactName string `json:"contact_name"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	Title       string `json:"title"`
	IsPrimary   bool   `json:"is_primary"`
	Status      string `json:"status"`
	Remark      string `json:"remark"`
}

type DeleteContactRequest struct {
	SpaceID   string `json:"space_id"`
	ContactID string `json:"contact_id"`
}

type ContactData struct {
	ContactID    string `json:"contact_id,omitempty"`
	TenantID     string `json:"tenant_id,omitempty"`
	SpaceID      string `json:"space_id,omitempty"`
	CustomerID   string `json:"customer_id,omitempty"`
	CustomerName string `json:"customer_name,omitempty"`
	ContactName  string `json:"contact_name,omitempty"`
	Mobile       string `json:"mobile,omitempty"`
	Email        string `json:"email,omitempty"`
	Title        string `json:"title,omitempty"`
	IsPrimary    bool   `json:"is_primary"`
	Status       string `json:"status,omitempty"`
	Remark       string `json:"remark,omitempty"`
	CreatedBy    string `json:"created_by,omitempty"`
	UpdatedBy    string `json:"updated_by,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

type ContactListData struct {
	List  []*ContactData `json:"list"`
	Total int64          `json:"total"`
}

type GetOpportunityRequest struct {
	SpaceID       string `query:"space_id"`
	OpportunityID string `query:"opportunity_id"`
}

type CreateOpportunityRequest struct {
	SpaceID           string `json:"space_id"`
	CustomerID        string `json:"customer_id"`
	OpportunityName   string `json:"opportunity_name"`
	Stage             string `json:"stage"`
	Amount            string `json:"amount"`
	ExpectedCloseDate string `json:"expected_close_date"`
	OwnerUserID       string `json:"owner_user_id"`
	OwnerUserName     string `json:"owner_user_name"`
	Status            string `json:"status"`
	Remark            string `json:"remark"`
}

type UpdateOpportunityRequest struct {
	SpaceID           string `json:"space_id"`
	OpportunityID     string `json:"opportunity_id"`
	CustomerID        string `json:"customer_id"`
	OpportunityName   string `json:"opportunity_name"`
	Stage             string `json:"stage"`
	Amount            string `json:"amount"`
	ExpectedCloseDate string `json:"expected_close_date"`
	OwnerUserID       string `json:"owner_user_id"`
	OwnerUserName     string `json:"owner_user_name"`
	Status            string `json:"status"`
	Remark            string `json:"remark"`
}

type DeleteOpportunityRequest struct {
	SpaceID       string `json:"space_id"`
	OpportunityID string `json:"opportunity_id"`
}

type OpportunityData struct {
	OpportunityID     string `json:"opportunity_id,omitempty"`
	TenantID          string `json:"tenant_id,omitempty"`
	SpaceID           string `json:"space_id,omitempty"`
	CustomerID        string `json:"customer_id,omitempty"`
	OpportunityName   string `json:"opportunity_name,omitempty"`
	Stage             string `json:"stage,omitempty"`
	Amount            string `json:"amount,omitempty"`
	ExpectedCloseDate string `json:"expected_close_date,omitempty"`
	OwnerUserID       string `json:"owner_user_id,omitempty"`
	OwnerUserName     string `json:"owner_user_name,omitempty"`
	Status            string `json:"status,omitempty"`
	Remark            string `json:"remark,omitempty"`
	CreatedBy         string `json:"created_by,omitempty"`
	UpdatedBy         string `json:"updated_by,omitempty"`
	CreatedAt         string `json:"created_at,omitempty"`
	UpdatedAt         string `json:"updated_at,omitempty"`
}

type OpportunityListData struct {
	List  []*OpportunityData `json:"list"`
	Total int64              `json:"total"`
}

type GetFollowRecordRequest struct {
	SpaceID        string `query:"space_id"`
	FollowRecordID string `query:"follow_record_id"`
}

type CreateFollowRecordRequest struct {
	SpaceID        string `json:"space_id"`
	CustomerID     string `json:"customer_id"`
	ContactID      string `json:"contact_id"`
	FollowType     string `json:"follow_type"`
	Content        string `json:"content"`
	NextFollowTime string `json:"next_follow_time"`
	OwnerUserID    string `json:"owner_user_id"`
	OwnerUserName  string `json:"owner_user_name"`
	Status         string `json:"status"`
}

type UpdateFollowRecordRequest struct {
	SpaceID        string `json:"space_id"`
	FollowRecordID string `json:"follow_record_id"`
	CustomerID     string `json:"customer_id"`
	ContactID      string `json:"contact_id"`
	FollowType     string `json:"follow_type"`
	Content        string `json:"content"`
	NextFollowTime string `json:"next_follow_time"`
	OwnerUserID    string `json:"owner_user_id"`
	OwnerUserName  string `json:"owner_user_name"`
	Status         string `json:"status"`
}

type DeleteFollowRecordRequest struct {
	SpaceID        string `json:"space_id"`
	FollowRecordID string `json:"follow_record_id"`
}

type FollowRecordData struct {
	FollowRecordID string `json:"follow_record_id,omitempty"`
	TenantID       string `json:"tenant_id,omitempty"`
	SpaceID        string `json:"space_id,omitempty"`
	CustomerID     string `json:"customer_id,omitempty"`
	ContactID      string `json:"contact_id,omitempty"`
	FollowType     string `json:"follow_type,omitempty"`
	Content        string `json:"content,omitempty"`
	NextFollowTime string `json:"next_follow_time,omitempty"`
	OwnerUserID    string `json:"owner_user_id,omitempty"`
	OwnerUserName  string `json:"owner_user_name,omitempty"`
	Status         string `json:"status,omitempty"`
	CreatedBy      string `json:"created_by,omitempty"`
	UpdatedBy      string `json:"updated_by,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

type FollowRecordListData struct {
	List  []*FollowRecordData `json:"list"`
	Total int64               `json:"total"`
}

type GetProductRequest struct {
	SpaceID   string `query:"space_id"`
	ProductID string `query:"product_id"`
}

type CreateProductRequest struct {
	SpaceID     string `json:"space_id"`
	ProductName string `json:"product_name"`
	ProductCode string `json:"product_code"`
	Category    string `json:"category"`
	UnitPrice   string `json:"unit_price"`
	Status      string `json:"status"`
	Remark      string `json:"remark"`
}

type UpdateProductRequest struct {
	SpaceID     string `json:"space_id"`
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	ProductCode string `json:"product_code"`
	Category    string `json:"category"`
	UnitPrice   string `json:"unit_price"`
	Status      string `json:"status"`
	Remark      string `json:"remark"`
}

type DeleteProductRequest struct {
	SpaceID   string `json:"space_id"`
	ProductID string `json:"product_id"`
}

type ProductData struct {
	ProductID   string `json:"product_id,omitempty"`
	TenantID    string `json:"tenant_id,omitempty"`
	SpaceID     string `json:"space_id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	ProductCode string `json:"product_code,omitempty"`
	Category    string `json:"category,omitempty"`
	UnitPrice   string `json:"unit_price,omitempty"`
	Status      string `json:"status,omitempty"`
	Remark      string `json:"remark,omitempty"`
	CreatedBy   string `json:"created_by,omitempty"`
	UpdatedBy   string `json:"updated_by,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type ProductListData struct {
	List  []*ProductData `json:"list"`
	Total int64          `json:"total"`
}

type GetSalesOrderRequest struct {
	SpaceID      string `query:"space_id"`
	SalesOrderID string `query:"sales_order_id"`
}

type CreateSalesOrderRequest struct {
	SpaceID       string `json:"space_id"`
	CustomerID    string `json:"customer_id"`
	OpportunityID string `json:"opportunity_id"`
	ProductID     string `json:"product_id"`
	ProductName   string `json:"product_name"`
	SalesUserID   string `json:"sales_user_id"`
	SalesUserName string `json:"sales_user_name"`
	Quantity      string `json:"quantity"`
	Amount        string `json:"amount"`
	OrderDate     string `json:"order_date"`
	Status        string `json:"status"`
	Remark        string `json:"remark"`
}

type UpdateSalesOrderRequest struct {
	SpaceID       string `json:"space_id"`
	SalesOrderID  string `json:"sales_order_id"`
	CustomerID    string `json:"customer_id"`
	OpportunityID string `json:"opportunity_id"`
	ProductID     string `json:"product_id"`
	ProductName   string `json:"product_name"`
	SalesUserID   string `json:"sales_user_id"`
	SalesUserName string `json:"sales_user_name"`
	Quantity      string `json:"quantity"`
	Amount        string `json:"amount"`
	OrderDate     string `json:"order_date"`
	Status        string `json:"status"`
	Remark        string `json:"remark"`
}

type DeleteSalesOrderRequest struct {
	SpaceID      string `json:"space_id"`
	SalesOrderID string `json:"sales_order_id"`
}

type SalesOrderData struct {
	SalesOrderID  string `json:"sales_order_id,omitempty"`
	TenantID      string `json:"tenant_id,omitempty"`
	SpaceID       string `json:"space_id,omitempty"`
	CustomerID    string `json:"customer_id,omitempty"`
	OpportunityID string `json:"opportunity_id,omitempty"`
	ProductID     string `json:"product_id,omitempty"`
	ProductName   string `json:"product_name,omitempty"`
	SalesUserID   string `json:"sales_user_id,omitempty"`
	SalesUserName string `json:"sales_user_name,omitempty"`
	Quantity      string `json:"quantity,omitempty"`
	Amount        string `json:"amount,omitempty"`
	OrderDate     string `json:"order_date,omitempty"`
	Status        string `json:"status,omitempty"`
	Remark        string `json:"remark,omitempty"`
	CreatedBy     string `json:"created_by,omitempty"`
	UpdatedBy     string `json:"updated_by,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}

type SalesOrderListData struct {
	List  []*SalesOrderData `json:"list"`
	Total int64             `json:"total"`
}

type DashboardOrderTrendData struct {
	Date        string `json:"date"`
	OrderCount  int64  `json:"order_count"`
	OrderAmount string `json:"order_amount"`
}

type DashboardOverviewData struct {
	CustomerTotal             int64                      `json:"customer_total"`
	NewCustomersThisMonth     int64                      `json:"new_customers_this_month"`
	OpportunityTotalAmount    string                     `json:"opportunity_total_amount"`
	NewOpportunitiesThisMonth int64                      `json:"new_opportunities_this_month"`
	SalesOrderTotalAmount     string                     `json:"sales_order_total_amount"`
	RecentOrderTrend          []*DashboardOrderTrendData `json:"recent_order_trend"`
}
