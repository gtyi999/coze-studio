#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
ENV_FILE="$ROOT_DIR/docker/.env"
DDL_FILE="$ROOT_DIR/backend/types/ddl/crm_phase1.sql"

MYSQL_CONTAINER="${MYSQL_CONTAINER:-coze-mysql}"
MYSQL_DATABASE="${MYSQL_DATABASE:-opencoze}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-}"

if [ -f "$ENV_FILE" ]; then
    set -a
    source "$ENV_FILE"
    set +a
fi

if [ ! -f "$DDL_FILE" ]; then
    echo "CRM DDL file not found: $DDL_FILE"
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

docker exec -i "$MYSQL_CONTAINER" mysql "-u${MYSQL_USER}" "-p${MYSQL_PASSWORD}" "${MYSQL_DATABASE}" < "$DDL_FILE"

echo "CRM phase 1 schema applied successfully."
echo "DDL source: $DDL_FILE"
