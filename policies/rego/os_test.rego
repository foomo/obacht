package bouncer.os_test

import data.bouncer.os

test_os001_always_passes {
    findings := os.findings with input as {
        "os": {
            "os": "linux",
            "arch": "amd64",
            "hostname": "myhost"
        }
    }
    count(findings) == 0
}
