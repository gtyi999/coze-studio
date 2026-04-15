# CRM Design Overview

## 1. Scan Baseline

This draft is based on the current `coze-studio` repository structure and current worktree state.

### 1.1 Backend structure observed

- `backend/api`
  - `handler`
  - `middleware`
  - `model`
  - `router`
- `backend/application`
  - existing bounded/application modules such as `prompt`, `shortcutcmd`, `im`, `memory`, `workflow`
  - current worktree already has `backend/application/crm`
- `backend/domain`
  - domain packages commonly use `entity`, `repository`, `service`
  - some older modules also carry `internal/dal`
  - current worktree already has `backend/domain/crm/{entity,repository,service}`
- `backend/infra`
  - infrastructure is separated by concern, for example `infra/im`, `infra/storage`, `infra/cache`
  - current worktree already has `backend/infra/crm/impl/mysql`
- `idl`
  - existing namespaces such as `conversation`, `workflow`, `im`, `crm`
  - current worktree already has `idl/crm/crm.thrift`

### 1.2 Frontend structure observed

- `frontend/apps/coze-studio`
  - app-level route registration and page wrappers
- `frontend/packages`
  - `foundation`: layout, space/workspace shell, menu adaptation
  - `studio`: workspace pages and adapters
  - `arch`: API schema and infrastructure packages
- current worktree already has:
  - `frontend/apps/coze-studio/src/pages/crm.tsx`
  - `frontend/packages/studio/workspace/entry-adapter/src/pages/crm/`
  - `frontend/packages/arch/api-schema/src/idl/crm/crm.ts`

## 2. Why CRM should be a new bounded context

CRM is a typical business domain with its own master data, lifecycle, permission semantics, and persistence model. It does not fit naturally into existing bounded contexts such as:

- `memory`: this is resource/data capability infrastructure, not line-of-business customer lifecycle management
- `conversation` / `workflow` / `agent`: these are AI runtime and orchestration domains
- `plugin` / `connector` / `prompt`: these are development resources, not operational business entities

Creating CRM as a dedicated bounded context is the cleanest way to:

- keep customer, contact, opportunity, follow record, product, and sales order models cohesive
- avoid mixing business CRUD logic into AI-resource domains
- preserve current DDD boundaries in `application`, `domain`, and `infra`
- leave room for future sales workflow, reporting, approval, and external ERP/OMS integration without polluting existing modules

Recommended bounded context name: `crm`

## 3. Best existing references in the repo

### 3.1 Backend reference

No single existing module matches CRM end to end, so the best reference should be assembled from two existing patterns:

- Domain/Application CRUD layering reference:
  - `backend/application/prompt`
  - `backend/domain/prompt/{entity,repository,service}`
  - `backend/application/shortcutcmd`
  - `backend/domain/shortcutcmd/{entity,repository,service}`
- Manual Hertz API and route registration reference:
  - `backend/api/router/coze/im_manual.go`
  - `backend/api/handler/coze/im_admin_service.go`

Reason:

- `prompt` and `shortcutcmd` show the repo's common `application -> domain service -> repository` orchestration style
- `im` shows how to add a non-generated CRUD-like API entry with `router/handler/model` while still coexisting with the generated router

### 3.2 Frontend reference

The closest existing workspace CRUD page reference is IM:

- route wrapper: `frontend/apps/coze-studio/src/pages/im.tsx`
- route registration: `frontend/apps/coze-studio/src/routes/index.tsx`
- lazy load registration: `frontend/apps/coze-studio/src/routes/async-components.tsx`
- workspace submenu integration:
  - `frontend/packages/foundation/space-ui-adapter/src/const.ts`
  - `frontend/packages/foundation/space-ui-adapter/src/components/workspace-sub-menu/index.tsx`
- page implementation:
  - `frontend/packages/studio/workspace/entry-adapter/src/pages/im/index.tsx`
  - `frontend/packages/studio/workspace/entry-adapter/src/pages/im/components/*`
- frontend API schema pattern:
  - `frontend/packages/arch/api-schema/src/idl/im/im.ts`

Reason:

- IM already uses workspace menu + route + list page + side sheet forms + manual API schema
- CRM as a business workspace submodule should follow the same integration seam instead of opening a separate app shell

## 4. CRM phase 1 scope and non-goals

### 4.1 In scope

CRM phase 1 should focus on business master data and basic operational records:

- `customer`
- `contact`
- `opportunity`
- `follow_record`
- `product`
- `sales_order`

Recommended delivery order inside phase 1:

1. `customer`
2. `contact`
3. `product`
4. `opportunity`
5. `follow_record`
6. `sales_order`

This keeps the implementation aligned with entity dependency order.

### 4.2 Out of scope for phase 1

- inventory, warehouse, delivery, receivable/payable
- contract approval engine
- reporting center / BI dashboards
- generic metadata engine or low-code schema editor
- external ERP / OMS / CRM sync
- fine-grained field-level permission system
- complex workflow automation beyond basic CRUD and list/search

## 5. Backend layering recommendation

The backend should remain inside the current DDD layering and should not create a parallel architecture.

### 5.1 Recommended landing points

- API layer
  - `backend/api/model/crm`
  - `backend/api/handler/coze/crm_service.go`
  - `backend/api/router/coze/crm_manual.go`
- Application layer
  - `backend/application/crm/init.go`
  - `backend/application/crm/*.go`
- Domain layer
  - `backend/domain/crm/entity`
  - `backend/domain/crm/repository`
  - `backend/domain/crm/service`
- Infra layer
  - `backend/infra/crm/impl/mysql`
- Error codes
  - `backend/types/errno/crm.go`
- App initialization
  - `backend/application/application.go`
- IDL
  - `idl/crm/crm.thrift`

### 5.2 Suggested domain split

- `entity`
  - aggregate/data structures for `Customer`, `Contact`, `Opportunity`, `FollowRecord`, `Product`, `SalesOrder`
  - common scope/value objects such as `Scope`, `PageOption`, audit info
- `repository`
  - persistence interfaces by aggregate or by subdomain
- `service`
  - domain validation
  - relation checks
  - delete guard rules
  - status normalization

### 5.3 Suggested application responsibilities

`backend/application/crm` should own:

- session/user extraction from context
- workspace access check
- tenant resolution strategy
- DTO to domain entity translation
- orchestration across multiple domain objects when needed

It should not own SQL details.

### 5.4 Suggested API style

Based on the current repo, CRM should follow the manual IM route style:

- route path prefix: `/api/crm/...`
- handler style: Hertz handlers under `backend/api/handler/coze`
- request/response structs under `backend/api/model/crm`
- generated router remains untouched except one explicit registration call in `backend/api/router/register.go`

Recommended naming pattern:

- `/api/crm/customer/list`
- `/api/crm/customer/get`
- `/api/crm/customer/create`
- `/api/crm/customer/update`
- `/api/crm/customer/delete`

Repeat the same pattern for other aggregates.

## 6. Frontend menu, route, and page recommendation

CRM should be added as a workspace submodule, not a new standalone frontend application.

### 6.1 Recommended landing points

- workspace submenu enum
  - `frontend/packages/foundation/space-ui-adapter/src/const.ts`
- workspace submenu rendering
  - `frontend/packages/foundation/space-ui-adapter/src/components/workspace-sub-menu/index.tsx`
- app route registration
  - `frontend/apps/coze-studio/src/routes/index.tsx`
  - `frontend/apps/coze-studio/src/routes/async-components.tsx`
- app page wrapper
  - `frontend/apps/coze-studio/src/pages/crm.tsx`
- workspace page implementation
  - `frontend/packages/studio/workspace/entry-adapter/src/pages/crm/`
- workspace adapter package export
  - `frontend/packages/studio/workspace/entry-adapter/package.json`
- frontend API schema
  - `frontend/packages/arch/api-schema/src/idl/crm/crm.ts`
  - `frontend/packages/arch/api-schema/src/index.ts`
  - `frontend/packages/arch/api-schema/package.json`

### 6.2 Recommended page organization

Recommended page organization under `frontend/packages/studio/workspace/entry-adapter/src/pages/crm/`:

- `index.tsx`
- `page.tsx`
- `types.ts`
- `constants.tsx`
- `utils.ts`
- `components/`

This is intentionally aligned to the existing IM page structure.

### 6.3 Recommended UX pattern

Follow the IM page pattern rather than inventing a new UI system:

- tabbed page or sectioned list page
- search + filter toolbar
- table list
- create/edit side sheet
- delete confirmation modal
- workspace route shape: `/space/:space_id/crm`

## 7. MySQL storage recommendation

CRM master data should use MySQL as the primary store.

### 7.1 Why MySQL is the right fit

- CRM entities are strongly relational
- CRUD, pagination, filtering, and transactional consistency are primary needs
- phase 1 does not require high-throughput event storage or document-style flexible schema
- current repo already initializes MySQL through Docker and Atlas

### 7.2 Schema maintenance points in current repo

Because of the current startup path in `docker/docker-compose.yml`, schema changes should be reflected in both:

- `docker/volumes/mysql/schema.sql`
- `docker/atlas/opencoze_latest_schema.hcl`

Otherwise local initialization and Atlas apply can drift.

### 7.3 Recommended table set

- `crm_customer`
- `crm_contact`
- `crm_opportunity`
- `crm_follow_record`
- `crm_product`
- `crm_sales_order`

### 7.4 Recommended common columns

Each table should include at least:

- primary key `id`
- `tenant_id`
- `space_id`
- `creator_id`
- `updater_id`
- `created_at`
- `updated_at`
- `is_deleted`

Plus business status columns such as `status`, and relation keys such as `customer_id`, `contact_id`, or `opportunity_id` where applicable.

## 8. Multi-tenant and soft delete recommendation

### 8.1 Tenant strategy

The repo currently centers most workspace-facing features on `space_id`, and permission checks are workspace-oriented. A practical CRM phase 1 strategy is:

- persist both `tenant_id` and `space_id`
- treat `space_id` as the immediate authorization boundary
- if the open-source deployment does not yet expose an independent tenant domain, allow `tenant_id` to be derived from `space_id` in phase 1

This keeps the schema future-proof without forcing a new tenancy subsystem right now.

### 8.2 Authorization strategy

Recommended application-layer check sequence:

1. extract user/session from context
2. verify workspace access for `space_id`
3. resolve `tenant_id`
4. pass `tenant_id + space_id` into domain/repository scope

### 8.3 Soft delete strategy

All CRM tables should use soft delete via `is_deleted`.

Recommended rules:

- read/list queries always filter active rows implicitly
- delete API performs soft delete only
- unique keys should be designed carefully around nullable code fields
- cross-entity delete should add guard rules to avoid obvious orphaning

## 9. Suggested phase breakdown

### Phase 0: design and scaffold

- finalize bounded context boundaries
- freeze table set and naming
- settle route naming and request/response style
- wire menu and route entry

### Phase 1: master data and minimal sales chain

- customer CRUD
- contact CRUD
- product CRUD
- opportunity CRUD
- follow record CRUD
- sales order CRUD

### Phase 2: business constraints

- relation-aware delete rules
- opportunity stage transitions
- sales order amount aggregation
- customer/contact association validations

### Phase 3: collaboration and integration

- richer permission rules
- workflow/event integration
- search/indexing integration
- external system sync

## 10. Recommended new directories

Recommended CRM directories, aligned to the current repo and current worktree:

- `backend/api/model/crm`
- `backend/api/handler/coze/crm_service.go`
- `backend/api/router/coze/crm_manual.go`
- `backend/application/crm`
- `backend/domain/crm/entity`
- `backend/domain/crm/repository`
- `backend/domain/crm/service`
- `backend/infra/crm/impl/mysql`
- `idl/crm`
- `frontend/apps/coze-studio/src/pages/crm.tsx`
- `frontend/packages/studio/workspace/entry-adapter/src/pages/crm`
- `frontend/packages/arch/api-schema/src/idl/crm`
- `docs/crm_design_overview.md`

## 11. Recommended existing modules to reuse

### Backend

- `backend/application/prompt`
  - good reference for application service orchestration
- `backend/domain/prompt`
  - good reference for `entity/repository/service` organization
- `backend/application/shortcutcmd`
  - good reference for permission + request-to-domain transformation
- `backend/domain/shortcutcmd`
  - good reference for simple CRUD domain service/repository interface design
- `backend/api/router/coze/im_manual.go`
  - best manual route registration reference
- `backend/api/handler/coze/im_admin_service.go`
  - best manual Hertz CRUD-style handler reference

### Frontend

- `frontend/apps/coze-studio/src/pages/im.tsx`
  - app-level workspace wrapper reference
- `frontend/apps/coze-studio/src/routes/index.tsx`
  - route insertion point reference
- `frontend/apps/coze-studio/src/routes/async-components.tsx`
  - lazy load registration reference
- `frontend/packages/foundation/space-ui-adapter/src/components/workspace-sub-menu/index.tsx`
  - workspace menu insertion point reference
- `frontend/packages/studio/workspace/entry-adapter/src/pages/im`
  - best page-level CRUD/list/side-sheet reference
- `frontend/packages/arch/api-schema/src/idl/im/im.ts`
  - best manual frontend API schema reference

## 12. Key development risks

- Schema drift risk
  - current repo startup path uses both `schema.sql` and Atlas HCL, so CRM table changes must update both
- Reference pattern mismatch risk
  - some older domains keep DAL under `domain/internal`, while CRM is expected to live under `infra/impl`; the design needs to stay consistent inside CRM even if older modules differ
- Permission model risk
  - current repo is workspace-oriented; if CRM later requires tenant-level or org-level permission, phase 1 should not overdesign this
- Delete consistency risk
  - customer/contact/opportunity/order relations can easily produce orphan records if delete guards are omitted
- IDL drift risk
  - adding manual handler/model code without keeping `idl/crm` and frontend api schema aligned will create request/response drift
- Frontend integration risk
  - CRM should be inserted into the existing workspace menu and route chain, not built as a parallel page system
- Scope creep risk
  - CRM often expands quickly into approval, BI, inventory, and automation; these should remain explicitly out of phase 1

## 13. Draft conclusion

The safest CRM secondary-development path for the current `coze-studio` repo is:

- create `crm` as an independent bounded context
- keep backend inside existing DDD layers
- use MySQL as the source of truth
- use IM's manual API and workspace-page pattern for integration
- keep phase 1 focused on business master data and minimal sales chain CRUD

This approach is closest to the repository's existing architecture and has the lowest chance of creating architectural debt.

## 14. Security And Audit

The concrete CRM permission boundary, tenant isolation strategy, audit field flow, and MySQL audit-log design are documented in:

- `docs/crm_security_audit.md`
- `docs/crm_quickstart.md`
- `docs/crm_deploy.md`
