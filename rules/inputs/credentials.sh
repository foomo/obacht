#!/bin/sh
home="$HOME"

get_mode() {
  if [ -f "$1" ]; then
    mode=$(stat -f '%Lp' "$1" 2>/dev/null || stat -c '%a' "$1" 2>/dev/null || echo "")
    echo "0$mode"
  else
    echo ""
  fi
}

# AWS credentials.
aws_creds="$home/.aws/credentials"
aws_exists=false
aws_mode=""
if [ -f "$aws_creds" ]; then
  aws_exists=true
  aws_mode=$(get_mode "$aws_creds")
fi

# .netrc file.
netrc="$home/.netrc"
netrc_exists=false
netrc_mode=""
if [ -f "$netrc" ]; then
  netrc_exists=true
  netrc_mode=$(get_mode "$netrc")
fi

# GCP application default credentials.
gcp_creds="$home/.config/gcloud/application_default_credentials.json"
gcp_exists=false
gcp_mode=""
if [ -f "$gcp_creds" ]; then
  gcp_exists=true
  gcp_mode=$(get_mode "$gcp_creds")
fi

# npmrc with auth token.
npmrc="$home/.npmrc"
npmrc_has_token=false
npmrc_mode=""
if [ -f "$npmrc" ]; then
  npmrc_mode=$(get_mode "$npmrc")
  if grep -q '_authToken' "$npmrc" 2>/dev/null; then
    npmrc_has_token=true
  fi
fi

printf '{"aws_exists": %s, "aws_mode": "%s", "netrc_exists": %s, "netrc_mode": "%s", "gcp_exists": %s, "gcp_mode": "%s", "npmrc_has_token": %s, "npmrc_mode": "%s"}' \
  "$aws_exists" "$aws_mode" "$netrc_exists" "$netrc_mode" "$gcp_exists" "$gcp_mode" "$npmrc_has_token" "$npmrc_mode"
