#!/bin/sh
set -e

# Generate config.yaml with runtime values
cat > /build/src/config/config.yaml << YAML
vault:
  path: "${VAULT_PATH:-/vault}"

backup:
  path: "${BACKUP_PATH:-/backup}"

template:
  path: "${TEMPLATE_PATH:-/templates}"

plugins:
  rest-api:
    url: ""
    token: ""
YAML

# Run the MCP server
exec /app/obsimcp
