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

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusOpen     = "open"
	StatusDraft    = "draft"

	StageInitial = "initial"
)

type Scope struct {
	TenantID int64
	SpaceID  int64
}

type PageOption struct {
	Page     int
	PageSize int
}

type AuditInfo struct {
	CreatedBy int64
	UpdatedBy int64

	// CreatorID/UpdaterID are kept for compatibility with the current
	// application/api layer and are normalized from CreatedBy/UpdatedBy.
	CreatorID int64
	UpdaterID int64

	CreatedAt int64
	UpdatedAt int64
	IsDeleted bool
}

func (a *AuditInfo) Normalize() {
	if a == nil {
		return
	}
	if a.CreatedBy == 0 {
		a.CreatedBy = a.CreatorID
	}
	if a.CreatorID == 0 {
		a.CreatorID = a.CreatedBy
	}
	if a.UpdatedBy == 0 {
		a.UpdatedBy = a.UpdaterID
	}
	if a.UpdaterID == 0 {
		a.UpdaterID = a.UpdatedBy
	}
}
