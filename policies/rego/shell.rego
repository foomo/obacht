package bouncer.shell

findings[f] {
    input.shell.history_file_mode != ""
    input.shell.history_file_mode != "0600"
    f := {
        "rule_id": "SHL001",
        "evidence": sprintf("History file %s has mode %s (expected 0600)", [input.shell.history_file, input.shell.history_file_mode])
    }
}
