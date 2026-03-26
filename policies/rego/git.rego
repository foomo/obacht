package bouncer.git

findings[f] {
    input.git.credential_helper == "store"
    f := {
        "rule_id": "GIT001",
        "evidence": "Git credential helper is set to 'store' which saves passwords in plaintext"
    }
}

findings[f] {
    input.git.installed
    not input.git.signing_enabled
    f := {
        "rule_id": "GIT002",
        "evidence": "Git commit signing is not enabled"
    }
}
