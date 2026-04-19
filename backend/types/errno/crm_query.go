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

const (
	ErrCRMQueryInvalidParamCode   = 113000010
	ErrCRMQueryIntentUnclearCode  = 113000011
	ErrCRMQuerySecurityRejectCode = 113000012
	ErrCRMQueryTimeoutCode        = 113000013
	ErrCRMQueryNoDataCode         = 113000014
	ErrCRMQueryFeaturePendingCode = 113000015
)

func init() {
	code.Register(
		ErrCRMQueryInvalidParamCode,
		"invalid crm query parameter: {msg}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrCRMQueryIntentUnclearCode,
		"crm query intent is unclear: {msg}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrCRMQuerySecurityRejectCode,
		"crm query rejected by security policy: {msg}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrCRMQueryTimeoutCode,
		"crm query timeout: {msg}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrCRMQueryNoDataCode,
		"crm query returned no data: {msg}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrCRMQueryFeaturePendingCode,
		"crm query feature is pending: {msg}",
		code.WithAffectStability(false),
	)
}
