package bouncer.docker_test

import data.bouncer.docker

test_doc001_permissive_socket {
    findings := docker.findings with input as {
        "docker": {
            "installed": true,
            "socket_exists": true,
            "socket_mode": "0666",
            "user_in_group": false
        }
    }
    count(findings) == 1
    some f in findings
    f.rule_id == "DOC001"
}

test_doc001_correct_socket_perms {
    findings := docker.findings with input as {
        "docker": {
            "installed": true,
            "socket_exists": true,
            "socket_mode": "0660",
            "user_in_group": false
        }
    }
    count(findings) == 0
}

test_doc001_no_socket {
    findings := docker.findings with input as {
        "docker": {
            "installed": true,
            "socket_exists": false,
            "socket_mode": "",
            "user_in_group": false
        }
    }
    count(findings) == 0
}

test_doc002_user_in_group {
    findings := docker.findings with input as {
        "docker": {
            "installed": true,
            "socket_exists": false,
            "socket_mode": "",
            "user_in_group": true
        }
    }
    count(findings) == 1
    some f in findings
    f.rule_id == "DOC002"
}

test_doc002_user_not_in_group {
    findings := docker.findings with input as {
        "docker": {
            "installed": true,
            "socket_exists": false,
            "socket_mode": "",
            "user_in_group": false
        }
    }
    count(findings) == 0
}
