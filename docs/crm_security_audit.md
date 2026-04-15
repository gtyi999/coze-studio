# CRM Security And Audit

## Permission Boundary

CRM reuses the current repository session and workspace model.

- Authentication comes from the existing request session in context.
- Authorization is workspace-scoped: a user must belong to the requested `space_id`.
- CRM does not trust a client-supplied `tenant_id`.
- The server resolves `tenant_id` from the authorized workspace before entering the domain and repository layers.

This keeps CRM aligned with the current open-source model, where `space` is the stable boundary already exposed by the user and permission subsystems.

## Tenant Isolation Strategy

CRM persists both `tenant_id` and `space_id` on every business table and enforces them in every query scope.

- List queries always filter by `tenant_id + space_id + is_deleted = false`.
- Detail queries always filter by `tenant_id + space_id + id + is_deleted = false`.
- Update and delete queries always filter by `tenant_id + space_id + id + is_deleted = false`.
- Cross-entity writes are validated in the domain layer before persistence.

Examples:

- Creating a contact loads the target customer by ID and rejects the write if the customer's tenant scope does not match the contact scope.
- Creating an opportunity validates the referenced customer and optional contact in the same tenant scope.
- Creating a sales order validates the referenced customer, product, and optional opportunity in the same tenant scope.

In the current repository model there is no separate tenant aggregate exposed to CRM yet, so the authorized workspace ID is used as the resolved tenant ID. This keeps the schema future-proof while preserving current compatibility.

## Audit Field Flow

All CRM write paths maintain the standard fields below:

- `created_by`
- `updated_by`
- `created_at`
- `updated_at`

Flow rules:

1. The application layer resolves the current user from session context.
2. On create:
   - `tenant_id` is injected from the authorized workspace scope.
   - `created_by` and `updated_by` are both set to the current user.
3. On update:
   - `tenant_id` is re-injected from the authorized workspace scope.
   - `updated_by` is set to the current user.
   - `created_by` and `created_at` are inherited from the current record.
4. On delete:
   - soft delete is used
   - `updated_by` is set from the current operator context
   - `updated_at` is refreshed

## Audit Log Scope

CRM adds a minimal MySQL-backed audit table: `crm_audit_log`.

Recorded operations:

- create customer
- update customer
- delete customer
- create opportunity
- create sales order

Each audit record stores:

- `tenant_id`
- `space_id`
- `resource_type`
- `resource_id`
- `action`
- `operator_id`
- `before_snapshot`
- `after_snapshot`
- `operation_at`

The repository writes the business row and the audit log in the same database transaction for the audited operations above.

## Operational Notes

- The API surface remains workspace-oriented.
- Tenant isolation is enforced server-side, not delegated to frontend parameters.
- Audit logs are stored in MySQL only and do not require Elasticsearch, ClickHouse, or any external event pipeline.
