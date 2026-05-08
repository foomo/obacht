#!/bin/sh
# Fixture input for the README demo recording.
# Emits a static JSON document so `obacht scan --rules-dir demo` produces
# the same set of findings on every machine.
set -eu

cat <<'JSON'
{
  "docker_socket": {
    "path": "/var/run/docker.sock",
    "mode": "0666",
    "world_writable": true
  },
  "ssh_keys": [
    {
      "path": "~/.ssh/id_rsa",
      "algorithm": "RSA",
      "bits": 2048,
      "mode": "0644"
    }
  ],
  "shell_history": {
    "path": "~/.bash_history",
    "mode": "0644",
    "world_readable": true
  }
}
JSON
