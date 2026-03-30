#!/bin/sh
# Check if git is installed.
if ! command -v git >/dev/null 2>&1; then
  printf '{"installed": false, "version": "", "credential_helper": "", "signing_enabled": false, "signing_format": "", "safe_directory_wildcard": false, "gitignore_missing_macos": [], "gitignore_missing_keys": [], "gitignore_missing_env": []}'
  exit 0
fi

version=$(git --version 2>/dev/null | tr -d '\n')
config=$(git config --global --list 2>/dev/null || true)

credential_helper=""
signing_enabled=false
signing_format=""
safe_directory_wildcard=false
excludes_file=""

IFS='
'
for line in $config; do
  key=$(printf '%s' "$line" | cut -d= -f1 | tr '[:upper:]' '[:lower:]')
  value=$(printf '%s' "$line" | cut -d= -f2-)
  case "$key" in
    credential.helper) credential_helper="$value" ;;
    commit.gpgsign)
      case "$(printf '%s' "$value" | tr '[:upper:]' '[:lower:]')" in
        true) signing_enabled=true ;;
      esac ;;
    gpg.format) signing_format="$value" ;;
    safe.directory)
      if [ "$value" = "*" ]; then safe_directory_wildcard=true; fi ;;
    core.excludesfile) excludes_file="$value" ;;
  esac
done

# Build a JSON array of patterns missing from a file.
# Usage: missing_patterns "pat1 pat2 pat3" /path/to/file
# Outputs: ["pat1","pat3"] or [] if all present.
missing_patterns() {
  patterns="$1"
  file="$2"
  result="["
  first=true
  oIFS="$IFS"; IFS=' '
  for p in $patterns; do
    if [ -z "$file" ] || ! grep -Fxq "$p" "$file" 2>/dev/null; then
      if [ "$first" = true ]; then first=false; else result="$result,"; fi
      result="$result\"$p\""
    fi
  done
  IFS="$oIFS"
  printf '%s]' "$result"
}

macos_patterns='.DS_Store .AppleDouble .LSOverride Icon? ._* .DocumentRevisions-V100 .fseventsd .Spotlight-V100 .TemporaryItems .Trashes .VolumeIcon.icns .com.apple.timemachine.donotpresent'
keys_patterns='id_rsa id_ed25519 id_ecdsa *.pem *.key *.p12 *.pfx *.cer *.crt *.jks *.p8 *.der *.jce *.keystore *.truststore *.p7b *.p7s *.p10 *.csr *.req *.p7c *.p7m *.p7r'
env_patterns='.env .env.local .env.*.local'

gitignore_file=""
if [ -n "$excludes_file" ]; then
  # Expand ~ in path.
  excludes_file=$(eval echo "$excludes_file")
  if [ -f "$excludes_file" ]; then
    gitignore_file="$excludes_file"
  fi
fi

# Empty gitignore_file means all patterns are missing.
missing_macos=$(missing_patterns "$macos_patterns" "$gitignore_file")
missing_keys=$(missing_patterns "$keys_patterns" "$gitignore_file")
missing_env=$(missing_patterns "$env_patterns" "$gitignore_file")

printf '{"installed": true, "version": "%s", "credential_helper": "%s", "signing_enabled": %s, "signing_format": "%s", "safe_directory_wildcard": %s, "gitignore_missing_macos": %s, "gitignore_missing_keys": %s, "gitignore_missing_env": %s}' \
  "$version" "$credential_helper" "$signing_enabled" "$signing_format" "$safe_directory_wildcard" "$missing_macos" "$missing_keys" "$missing_env"
