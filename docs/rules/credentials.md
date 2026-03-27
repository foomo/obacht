# Credentials Rules

These rules check that credential files used by developer tools have appropriate file permissions. Overly permissive credentials allow other users or processes on the system to steal authentication tokens and access keys.

## CRD001: AWS credentials file has weak permissions

**Severity:** high

The AWS credentials file (`~/.aws/credentials`) contains access keys that grant access to cloud resources. It should only be readable by the owner.

**What it checks:**
- Whether `~/.aws/credentials` exists
- Whether its permissions are `0600`

**Remediation:**
```bash
chmod 600 ~/.aws/credentials
```

## CRD002: .netrc file has weak permissions

**Severity:** high

The `.netrc` file stores login credentials for remote machines in plaintext. It should only be readable by the owner.

**What it checks:**
- Whether `~/.netrc` exists
- Whether its permissions are `0600`

**Remediation:**
```bash
chmod 600 ~/.netrc
```

## CRD003: GCP credentials file has weak permissions

**Severity:** high

The GCP application default credentials file contains tokens that grant access to Google Cloud resources. It should only be readable by the owner.

**What it checks:**
- Whether `~/.config/gcloud/application_default_credentials.json` exists
- Whether its permissions are `0600`

**Remediation:**
```bash
chmod 600 ~/.config/gcloud/application_default_credentials.json
```

## CRD004: .npmrc with auth token has weak permissions

**Severity:** high

The `.npmrc` file may contain npm authentication tokens that grant publish access to packages. When an auth token is present, the file should only be readable by the owner.

**What it checks:**
- Whether `~/.npmrc` contains an `_authToken` entry
- Whether its permissions are `0600`

**Remediation:**
```bash
chmod 600 ~/.npmrc
```
