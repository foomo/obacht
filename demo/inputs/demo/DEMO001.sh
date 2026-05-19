#!/bin/sh
# shellcheck shell=sh
# Fixture input for the README demo recording (DEMO001).
# Emits a static envelope so the scan output is deterministic.
set -eu

printf '%s\n' '{"rule_id":"DEMO001","status":"ok","data":{"path":"/var/run/docker.sock","mode":"0666","world_writable":true}}'
