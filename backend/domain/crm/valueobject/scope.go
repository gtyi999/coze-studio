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

package valueobject

import (
	"fmt"

	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

type TenantScope struct {
	TenantID int64
	SpaceID  int64
}

func (s TenantScope) Validate() error {
	if s.SpaceID <= 0 {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "space_id is required"))
	}
	if s.TenantID <= 0 {
		return errorx.New(errno.ErrCRMInvalidParamCode, errorx.KV("msg", "tenant_id is required"))
	}
	return nil
}

func (s TenantScope) EnsureSame(tenantID int64, spaceID int64, resource string) error {
	if s.TenantID == tenantID && s.SpaceID == spaceID {
		return nil
	}

	return errorx.New(
		errno.ErrCRMPermissionCode,
		errorx.KV("msg", fmt.Sprintf("%s tenant scope mismatch", resource)),
	)
}
