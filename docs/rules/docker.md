# Docker Rules

## DOC001: Docker socket has overly permissive access

**Severity:** high

The Docker socket (`/var/run/docker.sock`) provides root-equivalent access to the host system. If the socket is world-readable or world-writable, any user or process can control Docker and potentially escalate privileges.

**What it checks:**
- File permissions on `/var/run/docker.sock`
- Ensures the socket is not accessible by others

**Remediation:**
```bash
sudo chmod 660 /var/run/docker.sock
```

## DOC002: User is in the docker group

**Severity:** warn

Membership in the `docker` group grants root-equivalent access to the host via the Docker daemon. This is a known privilege escalation vector.

**What it checks:**
- Whether the current user is a member of the `docker` group

**Remediation:**

Consider using rootless Docker or requiring `sudo` for Docker commands instead of relying on group membership:

```bash
# Remove user from docker group
sudo gpasswd -d $USER docker

# Use rootless Docker instead
dockerd-rootless-setuptool.sh install
```
