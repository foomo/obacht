package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.file_sharing_disabled
	f := {"rule_id": "OS028", "evidence": "File Sharing (SMB) is enabled"}
}
