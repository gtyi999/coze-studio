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

type FollowRecord struct {
	FollowRecordID int64
	TenantID       int64
	SpaceID        int64

	CustomerID    int64
	ContactID     int64
	OpportunityID int64

	FollowType string
	Content    string

	// FollowContent is kept for compatibility with the current upper layer.
	FollowContent string

	NextFollowTime int64
	OwnerUserID    int64
	OwnerUserName  string
	Status         string

	AuditInfo
}

type FollowRecordFilter struct {
	Scope
	PageOption

	CustomerID          *int64
	ContactID           *int64
	OpportunityID       *int64
	OwnerUserID         *int64
	Keyword             string
	Status              *string
	CreatedAtStart      *int64
	CreatedAtEnd        *int64
	NextFollowTimeStart *int64
	NextFollowTimeEnd   *int64
}

func (r *FollowRecord) Normalize() {
	if r == nil {
		return
	}
	r.AuditInfo.Normalize()
	if r.Content == "" {
		r.Content = r.FollowContent
	}
	if r.FollowContent == "" {
		r.FollowContent = r.Content
	}
}
