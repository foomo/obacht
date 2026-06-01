#!/bin/sh
# shellcheck shell=sh
# include: _lib/json.sh
# include: _lib/bumblebee.sh

bumblebee_skip_if_unavailable BUM008
bumblebee_scan_ecosystem BUM008 browser-extension
