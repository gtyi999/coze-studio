# CRM Quickstart

## 1. Local prerequisites

- Start MySQL in Docker first:

```bash
docker compose -f docker/docker-compose.yml up -d mysql
```

- CRM local scripts default to:
  - `MYSQL_CONTAINER=coze-mysql`
  - `MYSQL_DATABASE=opencoze`
  - `CRM_SPACE_ID=1`
  - `CRM_TENANT_ID=$CRM_SPACE_ID`

If your local workspace is not `1`, export `CRM_SPACE_ID` before running the seed script.

## 2. Initialize CRM tables

Apply CRM schema only:

```bash
./scripts/setup/apply_crm_schema.sh
```

## 3. Import demo data

Import demo data only:

```bash
CRM_SPACE_ID=1 CRM_TENANT_ID=1 ./scripts/setup/apply_crm_demo_data.sh
```

One-click initialize schema plus demo data:

```bash
CRM_SPACE_ID=1 CRM_TENANT_ID=1 ./scripts/setup/init_crm_demo.sh
```

The demo import is idempotent inside the selected tenant scope. The script clears existing CRM rows for the same `tenant_id + space_id` before inserting fresh data.

## 4. Demo data shape

The seed script inserts a minimal but usable CRM dataset:

- 10 customers
- 20 contacts
  - every customer has 2 contacts
- 10 opportunities
- 20 follow records
- 10 products
- 20 sales orders

All demo rows are written into the same `tenant_id + space_id` scope so the CRM dashboard and list APIs can query them directly. The timestamps are generated relative to the current date, so the dashboard trend and recent updates remain meaningful after import.

## 5. Access the CRM page

After frontend and backend are running, open:

```text
/space/<space_id>/crm
```

Example:

```text
http://localhost:8888/space/1/crm
```

Use the same `space_id` that you passed to `CRM_SPACE_ID` when importing demo data.

## 6. Debug CRM APIs

Recommended local debugging flow:

1. Log in through the browser.
2. Open the CRM page.
3. Use DevTools Network panel to inspect requests under `/api/crm/...`.

Common endpoints:

- `GET /api/crm/dashboard/overview?space_id=<space_id>`
- `GET /api/crm/customer/list?space_id=<space_id>&page=1&page_size=10`
- `GET /api/crm/opportunity/list?space_id=<space_id>&page=1&page_size=10`
- `GET /api/crm/sales_order/list?space_id=<space_id>&page=1&page_size=10`
- `POST /api/crm/customer/create`

If you want to replay requests with `curl` or Postman, copy the browser cookie or auth header from a successful request and keep the same `space_id`.

## 7. Test commands

Backend:

```bash
cd backend
GOFLAGS='-ldflags=-checklinkname=0' go test ./domain/crm/... ./application/crm ./infra/crm/... ./api/router/coze -count=1
```

Frontend:

```bash
cd frontend/packages/studio/workspace/entry-adapter
corepack pnpm install
corepack pnpm test -- crm-page.test.tsx
```
