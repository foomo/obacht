#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/bumblebee.sh

bumblebee_skip_if_unavailable BUM001
bumblebee_scan_ecosystem BUM001 npm
