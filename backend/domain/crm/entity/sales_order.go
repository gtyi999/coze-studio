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

package entity

type SalesOrder struct {
	SalesOrderID int64
	TenantID     int64
	SpaceID      int64

	CustomerID    int64
	OpportunityID int64
	ProductID     int64

	ProductName   string
	SalesUserID   int64
	SalesUserName string
	Quantity      float64
	Amount        float64
	OrderDate     string
	Status        string

	Remark      string
	Description string

	// Compatibility fields for the current upper layer implementation.
	ContactID      int64
	OrderNo        string
	ProductSummary string
	TotalAmount    float64
	SignTime       int64

	AuditInfo
}

type SalesOrderFilter struct {
	Scope
	PageOption

	CustomerID     *int64
	OpportunityID  *int64
	ProductID      *int64
	SalesUserID    *int64
	Keyword        string
	Status         *string
	CreatedAtStart *int64
	CreatedAtEnd   *int64
	OrderDateStart *string
	OrderDateEnd   *string
}

func (o *SalesOrder) Normalize() {
	if o == nil {
		return
	}
	o.AuditInfo.Normalize()
	if o.Amount == 0 && o.TotalAmount != 0 {
		o.Amount = o.TotalAmount
	}
	if o.TotalAmount == 0 && o.Amount != 0 {
		o.TotalAmount = o.Amount
	}
	if o.Remark == "" {
		o.Remark = o.Description
	}
	if o.Description == "" {
		o.Description = o.Remark
	}
	if o.ProductName == "" {
		o.ProductName = o.ProductSummary
	}
	if o.ProductSummary == "" {
		o.ProductSummary = o.ProductName
	}
}
