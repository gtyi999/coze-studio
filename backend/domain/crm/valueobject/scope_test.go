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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

func TestTenantScopeValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		scope        TenantScope
		expectedCode int32
	}{
		{
			name:  "valid scope",
			scope: TenantScope{TenantID: 1, SpaceID: 2},
		},
		{
			name:         "missing tenant id",
			scope:        TenantScope{SpaceID: 2},
			expectedCode: errno.ErrCRMInvalidParamCode,
		},
		{
			name:         "missing space id",
			scope:        TenantScope{TenantID: 1},
			expectedCode: errno.ErrCRMInvalidParamCode,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.scope.Validate()
			if tt.expectedCode == 0 {
				require.NoError(t, err)
				return
			}

			requireErrorCode(t, err, tt.expectedCode)
		})
	}
}

func TestTenantScopeEnsureSame(t *testing.T) {
	t.Parallel()

	scope := TenantScope{TenantID: 1, SpaceID: 10}

	require.NoError(t, scope.EnsureSame(1, 10, "customer"))
	requireErrorCode(t, scope.EnsureSame(2, 10, "customer"), errno.ErrCRMPermissionCode)
}

func requireErrorCode(t *testing.T, err error, code int32) {
	t.Helper()

	require.Error(t, err)
	var statusErr errorx.StatusError
	require.True(t, errors.As(err, &statusErr))
	assert.Equal(t, code, statusErr.Code())
}
