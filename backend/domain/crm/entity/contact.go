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

type Contact struct {
	ContactID    int64
	TenantID     int64
	SpaceID      int64
	CustomerID   int64
	CustomerName string

	ContactName string
	Mobile      string
	Email       string

	Title    string
	Position string
	Gender   string

	IsPrimary bool
	Status    string

	Remark      string
	Description string

	AuditInfo
}

type ContactFilter struct {
	Scope
	PageOption

	CustomerID     *int64
	Keyword        string
	Status         *string
	CreatedAtStart *int64
	CreatedAtEnd   *int64
}

func (c *Contact) Normalize() {
	if c == nil {
		return
	}
	c.AuditInfo.Normalize()
	if c.Title == "" {
		c.Title = c.Position
	}
	if c.Position == "" {
		c.Position = c.Title
	}
	if c.Remark == "" {
		c.Remark = c.Description
	}
	if c.Description == "" {
		c.Description = c.Remark
	}
}
