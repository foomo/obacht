package bouncer.kube_test

import data.bouncer.kube

test_kub001_weak_permissions {
    findings := kube.findings with input as {
        "kube": {
            "config_exists": true,
            "config_mode": "0644",
            "contexts": []
        }
    }
    count(findings) == 1
    some f in findings
    f.rule_id == "KUB001"
}

test_kub001_correct_permissions {
    findings := kube.findings with input as {
        "kube": {
            "config_exists": true,
            "config_mode": "0600",
            "contexts": []
        }
    }
    count(findings) == 0
}

test_kub001_no_config {
    findings := kube.findings with input as {
        "kube": {
            "config_exists": false,
            "config_mode": "",
            "contexts": []
        }
    }
    count(findings) == 0
}

test_kub002_prod_context {
    findings := kube.findings with input as {
        "kube": {
            "config_exists": true,
            "config_mode": "0600",
            "contexts": [
                {"name": "prod-context", "cluster": "prod-cluster"},
                {"name": "dev-context", "cluster": "dev-cluster"}
            ]
        }
    }
    count(findings) == 1
    some f in findings
    f.rule_id == "KUB002"
}

test_kub002_no_prod_context {
    findings := kube.findings with input as {
        "kube": {
            "config_exists": true,
            "config_mode": "0600",
            "contexts": [
                {"name": "dev-context", "cluster": "dev-cluster"},
                {"name": "staging-context", "cluster": "staging-cluster"}
            ]
        }
    }
    count(findings) == 0
}
