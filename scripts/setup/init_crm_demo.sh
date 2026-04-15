#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

"$SCRIPT_DIR/apply_crm_schema.sh"
"$SCRIPT_DIR/apply_crm_demo_data.sh"

echo "CRM local demo initialization finished."
echo "You can now open: /space/${CRM_SPACE_ID:-1}/crm"
