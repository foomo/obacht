#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# GCP application default credentials permissions.
gcp_creds="$HOME/.config/gcloud/application_default_credentials.json"
gcp_exists=false
gcp_mode=""
if [ -f "$gcp_creds" ]; then
  gcp_exists=true
  gcp_mode=$(stat -f '%Lp' "$gcp_creds" 2>/dev/null || stat -c '%a' "$gcp_creds" 2>/dev/null || echo "")
  gcp_mode="0$gcp_mode"
fi

printf '{"gcp_exists": %s, "gcp_mode": "%s"}' "$gcp_exists" "$gcp_mode" | emit_ok CRD003
