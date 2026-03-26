package bouncer.shell_test

import data.bouncer.shell

test_shl001_weak_history_perms {
    findings := shell.findings with input as {
        "shell": {
            "shell": "/bin/zsh",
            "history_file": "/home/user/.zsh_history",
            "history_file_mode": "0644",
            "histcontrol": ""
        }
    }
    count(findings) == 1
    some f in findings
    f.rule_id == "SHL001"
}

test_shl001_correct_history_perms {
    findings := shell.findings with input as {
        "shell": {
            "shell": "/bin/zsh",
            "history_file": "/home/user/.zsh_history",
            "history_file_mode": "0600",
            "histcontrol": ""
        }
    }
    count(findings) == 0
}

test_shl001_no_history_file {
    findings := shell.findings with input as {
        "shell": {
            "shell": "/bin/zsh",
            "history_file": "",
            "history_file_mode": "",
            "histcontrol": ""
        }
    }
    count(findings) == 0
}
