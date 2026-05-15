package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.os_auto_download_enabled
	f := {"rule_id": "OS035", "evidence": "Automatic download of OS updates is disabled"}
}
