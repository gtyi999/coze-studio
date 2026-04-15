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

package errno

import "github.com/coze-dev/coze-studio/backend/pkg/errorx/code"

// CRM: 113 000 000 ~ 113 999 999
const (
	ErrCRMInvalidParamCode   = 113000000
	ErrCRMPermissionCode     = 113000001
	ErrCRMRecordNotFoundCode = 113000002
	ErrCRMOperateCode        = 113000003
)

func init() {
	code.Register(
		ErrCRMInvalidParamCode,
		"invalid crm parameter: {msg}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrCRMPermissionCode,
		"crm permission denied: {msg}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrCRMRecordNotFoundCode,
		"crm record not found",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrCRMOperateCode,
		"crm operate failed: {msg}",
		code.WithAffectStability(true),
	)
}
