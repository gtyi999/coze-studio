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

package crm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	crossuser "github.com/coze-dev/coze-studio/backend/crossdomain/user"
	"github.com/coze-dev/coze-studio/backend/domain/crm/auditctx"
	"github.com/coze-dev/coze-studio/backend/domain/crm/entity"
	domainservice "github.com/coze-dev/coze-studio/backend/domain/crm/service"
	userentity "github.com/coze-dev/coze-studio/backend/domain/user/entity"
	mockcrm "github.com/coze-dev/coze-studio/backend/internal/mock/domain/crm"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

type fakeCrossUserService struct {
	spaces []*userentity.Space
}

func (f *fakeCrossUserService) GetUserSpaceList(_ context.Context, _ int64) ([]*userentity.Space, error) {
	return f.spaces, nil
}

func (f *fakeCrossUserService) GetUserSpaceBySpaceID(_ context.Context, _ []int64) ([]*userentity.Space, error) {
	return f.spaces, nil
}

func TestCRMApplicationServiceCreateCustomerInjectsScopeAndAudit(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mockcrm.NewMockRepository(ctrl)

	svc := &CRMApplicationService{
		DomainSVC: domainservice.NewService(&domainservice.Components{
			Repository: repo,
		}),
	}
	ctx := newCRMApplicationTestContext(t)

	repo.EXPECT().CreateCustomer(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, customer *entity.Customer) (*entity.Customer, error) {
		assert.Equal(t, int64(1), customer.TenantID)
		assert.Equal(t, int64(1), customer.SpaceID)
		assert.Equal(t, int64(123), customer.CreatedBy)
		assert.Equal(t, int64(123), customer.UpdatedBy)
		assert.Equal(t, int64(123), customer.CreatorID)
		assert.Equal(t, int64(123), customer.UpdaterID)
		assert.Equal(t, int64(123), customer.OwnerUserID)
		customer.CustomerID = 9001
		return customer, nil
	})

	customer, err := svc.CreateCustomer(ctx, &entity.Customer{
		SpaceID:      1,
		CustomerName: "Acme Health",
		AuditInfo:    entity.AuditInfo{},
	})

	require.NoError(t, err)
	require.NotNil(t, customer)
	assert.Equal(t, int64(9001), customer.CustomerID)
}

func TestCRMApplicationServiceDeleteCustomerInjectsActorContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mockcrm.NewMockRepository(ctrl)

	svc := &CRMApplicationService{
		DomainSVC: domainservice.NewService(&domainservice.Components{
			Repository: repo,
		}),
	}
	ctx := newCRMApplicationTestContext(t)

	repo.EXPECT().GetCustomerByID(gomock.Any(), int64(2001)).Return(&entity.Customer{
		CustomerID: 2001,
		TenantID:   1,
		SpaceID:    1,
	}, nil)
	repo.EXPECT().CountActiveContactsByCustomer(gomock.Any(), &entity.Scope{TenantID: 1, SpaceID: 1}, int64(2001)).Return(int64(0), nil)
	repo.EXPECT().CountActiveOpportunitiesByCustomer(gomock.Any(), &entity.Scope{TenantID: 1, SpaceID: 1}, int64(2001)).Return(int64(0), nil)
	repo.EXPECT().DeleteCustomer(gomock.Any(), &entity.Scope{TenantID: 1, SpaceID: 1}, int64(2001)).DoAndReturn(func(ctx context.Context, scope *entity.Scope, customerID int64) error {
		actor, ok := auditctx.ActorFromContext(ctx)
		require.True(t, ok)
		assert.Equal(t, int64(2001), customerID)
		assert.Equal(t, int64(1), scope.TenantID)
		assert.Equal(t, int64(1), scope.SpaceID)
		assert.Equal(t, int64(123), actor.UserID)
		assert.Equal(t, int64(1), actor.TenantID)
		assert.Equal(t, int64(1), actor.SpaceID)
		return nil
	})

	require.NoError(t, svc.DeleteCustomer(ctx, 1, 2001))
}

func newCRMApplicationTestContext(t *testing.T) context.Context {
	t.Helper()

	oldCrossUser := crossuser.DefaultSVC()
	crossuser.SetDefaultSVC(&fakeCrossUserService{
		spaces: []*userentity.Space{
			{ID: 1, Name: "default"},
		},
	})
	t.Cleanup(func() {
		crossuser.SetDefaultSVC(oldCrossUser)
	})

	ctx := ctxcache.Init(context.Background())
	ctxcache.Store(ctx, consts.SessionDataKeyInCtx, &userentity.Session{
		UserID: 123,
	})
	return ctx
}
