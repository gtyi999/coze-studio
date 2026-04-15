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

type Product struct {
	ProductID int64
	TenantID  int64
	SpaceID   int64

	ProductName string
	ProductCode string
	Category    string
	Unit        string
	UnitPrice   float64
	Status      string

	Remark      string
	Description string

	AuditInfo
}

type ProductFilter struct {
	Scope
	PageOption

	Keyword        string
	Status         *string
	CreatedAtStart *int64
	CreatedAtEnd   *int64
}

func (p *Product) Normalize() {
	if p == nil {
		return
	}
	p.AuditInfo.Normalize()
	if p.Remark == "" {
		p.Remark = p.Description
	}
	if p.Description == "" {
		p.Description = p.Remark
	}
}
