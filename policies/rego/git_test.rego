package bouncer.git_test

import data.bouncer.git

test_git001_store_helper {
    findings := git.findings with input as {
        "git": {
            "installed": true,
            "credential_helper": "store",
            "signing_enabled": true
        }
    }
    count(findings) == 1
    some f in findings
    f.rule_id == "GIT001"
}

test_git001_safe_helper {
    findings := git.findings with input as {
        "git": {
            "installed": true,
            "credential_helper": "osxkeychain",
            "signing_enabled": true
        }
    }
    count(findings) == 0
}

test_git002_no_signing {
    findings := git.findings with input as {
        "git": {
            "installed": true,
            "credential_helper": "osxkeychain",
            "signing_enabled": false
        }
    }
    count(findings) == 1
    some f in findings
    f.rule_id == "GIT002"
}
