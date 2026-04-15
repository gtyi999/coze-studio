#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
ENV_FILE="$ROOT_DIR/docker/.env"
DEMO_SQL_FILE="$ROOT_DIR/backend/types/ddl/crm_demo_data.sql.tpl"

MYSQL_CONTAINER="${MYSQL_CONTAINER:-coze-mysql}"
MYSQL_DATABASE="${MYSQL_DATABASE:-opencoze}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-}"

CRM_SPACE_ID="${CRM_SPACE_ID:-1}"
CRM_TENANT_ID="${CRM_TENANT_ID:-$CRM_SPACE_ID}"

if [ -f "$ENV_FILE" ]; then
    set -a
    source "$ENV_FILE"
    set +a
fi

if [ ! -f "$DEMO_SQL_FILE" ]; then
    echo "CRM demo SQL template not found: $DEMO_SQL_FILE"
    exit 1
fi

if ! docker ps --format '{{.Names}}' | grep -q "^${MYSQL_CONTAINER}\$"; then
    echo "MySQL container '${MYSQL_CONTAINER}' is not running."
    echo "Please start MySQL first, for example:"
    echo "  docker compose -f docker/docker-compose.yml up -d mysql"
    exit 1
fi

if [ -z "$MYSQL_PASSWORD" ]; then
    if [ "${MYSQL_USER}" = "root" ]; then
        MYSQL_PASSWORD="${MYSQL_ROOT_PASSWORD:-root}"
    else
        MYSQL_PASSWORD="${MYSQL_PASSWORD:-coze123}"
    fi
fi

TMP_SQL="$(mktemp)"
trap 'rm -f "$TMP_SQL"' EXIT

sed \
    -e "s/__CRM_TENANT_ID__/${CRM_TENANT_ID}/g" \
    -e "s/__CRM_SPACE_ID__/${CRM_SPACE_ID}/g" \
    "$DEMO_SQL_FILE" > "$TMP_SQL"

docker exec -i "$MYSQL_CONTAINER" mysql "-u${MYSQL_USER}" "-p${MYSQL_PASSWORD}" "${MYSQL_DATABASE}" < "$TMP_SQL"

echo "CRM demo data imported successfully."
echo "Tenant scope: tenant_id=${CRM_TENANT_ID}, space_id=${CRM_SPACE_ID}"
echo "Seed counts: 10 customers, 20 contacts, 10 opportunities, 20 follow records, 10 products, 20 sales orders"
