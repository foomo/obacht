#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# shellcheck source=../_lib/json.sh

# AWS credentials file permissions.
aws_creds="$HOME/.aws/credentials"
aws_exists=false
aws_mode=""
if [ -f "$aws_creds" ]; then
  aws_exists=true
  aws_mode=$(stat -f '%Lp' "$aws_creds" 2>/dev/null || stat -c '%a' "$aws_creds" 2>/dev/null || echo "")
  aws_mode="0$aws_mode"
fi

printf '{"aws_exists": %s, "aws_mode": "%s"}' "$aws_exists" "$aws_mode" | emit_ok CRD001
