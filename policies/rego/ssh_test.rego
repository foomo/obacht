package bouncer.ssh_test

import data.bouncer.ssh

test_ssh001_weak_key_perms {
    findings := ssh.findings with input as {
        "ssh": {
            "directory_exists": true,
            "directory_mode": "0700",
            "keys": [{"path": "/home/user/.ssh/id_rsa", "mode": "0644", "type": "rsa"}],
            "config_exists": false
        }
    }
    count(findings) == 1
    some f in findings
    f.rule_id == "SSH001"
}

test_ssh001_correct_key_perms {
    findings := ssh.findings with input as {
        "ssh": {
            "directory_exists": true,
            "directory_mode": "0700",
            "keys": [{"path": "/home/user/.ssh/id_rsa", "mode": "0600", "type": "rsa"}],
            "config_exists": false
        }
    }
    count(findings) == 0
}

test_ssh002_weak_dir_perms {
    findings := ssh.findings with input as {
        "ssh": {
            "directory_exists": true,
            "directory_mode": "0755",
            "keys": [],
            "config_exists": false
        }
    }
    count(findings) == 1
    some f in findings
    f.rule_id == "SSH002"
}
