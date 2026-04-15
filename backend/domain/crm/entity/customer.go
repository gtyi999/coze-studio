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

type Customer struct {
	CustomerID int64
	TenantID   int64
	SpaceID    int64

	CustomerName string
	CustomerCode string
	Industry     string

	Level         string
	CustomerLevel string
	CustomerType  string

	OwnerUserID   int64
	OwnerUserName string

	Status  string
	Mobile  string
	Email   string
	Address string

	Remark      string
	Description string

	AuditInfo
}

type CustomerFilter struct {
	Scope
	PageOption

	Keyword        string
	OwnerUserID    *int64
	Status         *string
	CreatedAtStart *int64
	CreatedAtEnd   *int64
}

func (c *Customer) Normalize() {
	if c == nil {
		return
	}
	c.AuditInfo.Normalize()
	if c.Level == "" {
		c.Level = c.CustomerLevel
	}
	if c.CustomerLevel == "" {
		c.CustomerLevel = c.Level
	}
	if c.Remark == "" {
		c.Remark = c.Description
	}
	if c.Description == "" {
		c.Description = c.Remark
	}
}
