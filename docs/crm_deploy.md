# CRM Deploy Guide

## 1. Deployment stance

CRM phase 1 is deployed inside the existing `coze-server` service.

- No standalone CRM container is required.
- No extra gateway, queue, or storage component is required for CRM itself.
- CRM reuses the current login session, workspace membership, tenant resolution, and frontend `coze-web` entry.

## 2. CRM module dependencies

CRM depends on the same runtime that `coze-studio` already uses:

- `coze-server`: serves all CRM APIs under `/api/crm/...`
- `coze-web`: serves the CRM page under `/space/:space_id/crm`
- MySQL: stores CRM customers, contacts, opportunities, follow records, products, sales orders, and audit logs
- Existing session and workspace model: used for permission checks and tenant isolation

CRM does not require:

- a separate container
- a separate database instance
- ClickHouse, Elasticsearch, or other analytical storage
- new runtime environment variables

## 3. Configuration change checklist

### 3.1 Docker / Compose

No new Compose service is needed.

CRM is enabled by the existing startup chain as long as these files already include the CRM schema:

- `docker/volumes/mysql/schema.sql`
- `docker/atlas/opencoze_latest_schema.hcl`

Current Compose files already mount and apply the shared MySQL schema for `coze-server`:

- `docker/docker-compose.yml`
- `docker/docker-compose-debug.yml`

### 3.2 Helm

No new Helm deployment, service, or container is needed.

CRM must be included in the existing MySQL init job used by the chart:

- `helm/charts/opencoze/files/mysql/schema.sql`
- `helm/charts/opencoze/templates/mysql-init-configmap.yaml`
- `helm/charts/opencoze/templates/mysql-init-job.yaml`

This guide syncs the Helm schema with the Docker schema so Helm install and Helm upgrade both create the CRM tables.

### 3.3 Environment variables

No CRM-specific runtime environment variables are required.

CRM continues to use the existing server and database configuration:

- `LISTEN_ADDR`
- `SERVER_HOST`
- `MYSQL_HOST`
- `MYSQL_PORT`
- `MYSQL_USER`
- `MYSQL_PASSWORD`
- `MYSQL_DATABASE`
- `MYSQL_DSN`

If your current `coze-studio` deployment can connect to MySQL and serve workspace APIs, CRM does not need extra env wiring.

## 4. MySQL initialization and migration

### 4.1 Fresh Docker / Compose install

For a brand-new MySQL data directory, the existing Compose startup path is enough:

1. MySQL boots with `docker/volumes/mysql/schema.sql`
2. Atlas applies `docker/atlas/opencoze_latest_schema.hcl`
3. `coze-server` starts after MySQL is healthy

Recommended command:

```bash
docker compose -f docker/docker-compose.yml up -d
```

### 4.2 Existing local MySQL volume

If your local MySQL data directory was created before CRM was added, the mounted `init.sql` will not rerun automatically.

Apply the CRM schema explicitly:

```bash
./scripts/setup/apply_crm_schema.sh
```

Optional demo data import:

```bash
CRM_SPACE_ID=1 CRM_TENANT_ID=1 ./scripts/setup/init_crm_demo.sh
```

### 4.3 Existing server deployment

For an already-running environment, use one of these migration paths:

1. Apply the latest shared schema through Atlas.
2. Or run the CRM DDL directly with `backend/types/ddl/crm_phase1.sql`.

Atlas example:

```bash
export ATLAS_URL="mysql://coze:coze123@mysql:3306/opencoze?charset=utf8mb4&parseTime=True"
atlas schema apply -u "$ATLAS_URL" --to "file://docker/atlas/opencoze_latest_schema.hcl" --exclude "atlas_schema_revisions,table_*" --auto-approve
```

Manual DDL example:

```bash
mysql -h <mysql-host> -u <user> -p <database> < backend/types/ddl/crm_phase1.sql
```

### 4.4 Helm upgrade path

The Helm chart runs a post-install and post-upgrade MySQL init job.

After syncing the chart schema, upgrade the existing release:

```bash
helm upgrade <release-name> helm/charts/opencoze -n <namespace>
```

Then verify the MySQL init job completed successfully before checking `coze-server`.

## 5. Local development startup

### 5.1 Start shared dependencies

```bash
docker compose -f docker/docker-compose.yml up -d mysql redis elasticsearch minio etcd milvus nsqlookupd nsqd nsqadmin coze-server coze-web
```

Or start the full stack:

```bash
docker compose -f docker/docker-compose.yml up -d
```

### 5.2 Initialize CRM tables if needed

```bash
./scripts/setup/apply_crm_schema.sh
```

### 5.3 Import demo data for local verification

```bash
CRM_SPACE_ID=1 CRM_TENANT_ID=1 ./scripts/setup/init_crm_demo.sh
```

### 5.4 Open the CRM page

```text
http://localhost:8888/space/1/crm
```

## 6. Docker / Compose startup checkpoints

Use the checkpoints below after startup:

1. Check container status:

```bash
docker compose -f docker/docker-compose.yml ps
```

2. Confirm MySQL schema apply completed:

```bash
docker logs coze-mysql
```

Expected signal:

- `Atlas migrations completed successfully`

3. Confirm CRM tables exist:

```bash
docker exec -i coze-mysql mysql -uroot -proot opencoze -e "SHOW TABLES LIKE 'crm_%';"
```

Expected tables:

- `crm_customer`
- `crm_contact`
- `crm_opportunity`
- `crm_follow_record`
- `crm_product`
- `crm_sales_order`
- `crm_audit_log`

4. Confirm backend route is reachable:

```bash
curl -I http://localhost:8888/api/crm/customer/list
```

You should get an application response instead of a proxy or upstream failure.

## 7. Helm startup checkpoints

1. Install or upgrade the current release:

```bash
helm upgrade --install <release-name> helm/charts/opencoze -n <namespace> --create-namespace
```

2. Check the MySQL init hook:

```bash
kubectl get jobs -n <namespace>
kubectl logs job/<release-name>-mysql-init -n <namespace>
```

3. Check server and web rollout:

```bash
kubectl get pods -n <namespace>
kubectl rollout status deployment/<release-name>-server -n <namespace>
kubectl rollout status deployment/<release-name>-web -n <namespace>
```

4. Check CRM tables inside MySQL:

```bash
kubectl exec -it <mysql-pod> -n <namespace> -- mysql -u root -p<root-password> <database> -e "SHOW TABLES LIKE 'crm_%';"
```

## 8. Troubleshooting

### 8.1 CRM page opens but data is empty

Possible causes:

- no CRM rows exist for the current `tenant_id + space_id`
- you imported demo data into a different `space_id`
- the current user is not in the target workspace

Recommended checks:

- open `/space/<space_id>/crm`
- inspect `/api/crm/...` requests in browser DevTools
- import demo data again with the correct `CRM_SPACE_ID`

### 8.2 CRM API returns not found or route missing

Possible causes:

- backend image is older than the CRM code
- the current deployment did not include the updated `coze-server`
- route registration was not published together with the backend binary

Recommended checks:

- verify the backend image tag
- verify `coze-server` rollout completed
- inspect server logs for CRM route registration and request handling

### 8.3 CRM tables are missing after upgrade

Possible causes:

- local MySQL data directory was initialized before CRM landed
- Helm release still uses an older `files/mysql/schema.sql`
- the MySQL init job failed during `helm upgrade`

Recommended fixes:

- rerun `./scripts/setup/apply_crm_schema.sh` for local Docker
- rerun `atlas schema apply` against the target database
- inspect `kubectl logs job/<release-name>-mysql-init -n <namespace>`

### 8.4 Permission denied in CRM APIs

This is usually expected behavior rather than a deployment failure.

CRM enforces:

- logged-in user required
- workspace membership required
- tenant and space scope match required

Check whether the current account belongs to the requested `space_id`.

### 8.5 Frontend route exists but shows old UI

Possible causes:

- `coze-web` image was not updated
- browser cache still serves an older bundle

Recommended fixes:

- redeploy `coze-web`
- hard refresh the browser
- confirm the current frontend image tag matches the backend release

## 9. Release recommendation

For a normal CRM release inside the current `coze-studio` system:

1. Publish the updated `coze-server` image.
2. Publish the updated `coze-web` image.
3. Apply the shared MySQL schema update through Compose or Helm.
4. Verify `crm_%` tables exist.
5. Verify `/space/<space_id>/crm` loads and `/api/crm/dashboard/overview` returns data.

The CRM module should always ship as part of the existing `coze-server` and `coze-web` release, not as a standalone service.
