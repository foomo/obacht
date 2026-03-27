package schema

// Facts is the type-safe contract between collectors and policies.
type Facts struct {
	SchemaVersion string      `json:"schema_version"`
	OS            OSFacts     `json:"os"`
	SSH           SSHFacts    `json:"ssh"`
	Git           GitFacts    `json:"git"`
	Docker        DockerFacts `json:"docker"`
	Kube          KubeFacts   `json:"kube"`
	Env           EnvFacts    `json:"env"`
	Shell         ShellFacts  `json:"shell"`
	Tools         ToolsFacts  `json:"tools"`
	Path          PathFacts   `json:"path"`
}

// OSFacts contains operating system information.
type OSFacts struct {
	OS                        string `json:"os"`
	Arch                      string `json:"arch"`
	Hostname                  string `json:"hostname"`
	SIPEnabled                bool   `json:"sip_enabled"`
	FileVaultEnabled          bool   `json:"filevault_enabled"`
	FirewallEnabled           bool   `json:"firewall_enabled"`
	StealthModeEnabled        bool   `json:"stealth_mode_enabled"`
	GatekeeperEnabled         bool   `json:"gatekeeper_enabled"`
	AutoLoginDisabled         bool   `json:"auto_login_disabled"`
	GuestAccountDisabled      bool   `json:"guest_account_disabled"`
	ScreenLockTimeoutSecs     int    `json:"screen_lock_timeout_seconds"`
	OSAutoUpdateEnabled       bool   `json:"os_auto_update_enabled"`
	AppAutoUpdateEnabled      bool   `json:"app_auto_update_enabled"`
	RSREnabled                bool   `json:"rsr_enabled"`
	ScreenSharingDisabled     bool   `json:"screen_sharing_disabled"`
	InternetSharingDisabled   bool   `json:"internet_sharing_disabled"`
	PrinterSharingDisabled    bool   `json:"printer_sharing_disabled"`
	RemoteAppleEventsDisabled bool   `json:"remote_apple_events_disabled"`
	AirdropSetting            string `json:"airdrop_setting"`
	RosettaInstalled          bool   `json:"rosetta_installed"`
	EDRDeployed               bool   `json:"edr_deployed"`
	LegacyKextsBlocked        bool   `json:"legacy_kexts_blocked"`
	MDMEnrolled               bool   `json:"mdm_enrolled"`
}

// SSHFacts contains SSH configuration and key information.
type SSHFacts struct {
	DirectoryExists bool     `json:"directory_exists"`
	DirectoryMode   string   `json:"directory_mode"`
	Keys            []SSHKey `json:"keys"`
	ConfigExists    bool     `json:"config_exists"`
}

// SSHKey describes a single SSH key file.
type SSHKey struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
}

// GitFacts contains Git configuration information.
type GitFacts struct {
	Installed        bool   `json:"installed"`
	Version          string `json:"version"`
	CredentialHelper string `json:"credential_helper"`
	SigningEnabled   bool   `json:"signing_enabled"`
	SigningFormat    string `json:"signing_format"`
}

// DockerFacts contains Docker configuration information.
type DockerFacts struct {
	Installed    bool   `json:"installed"`
	SocketExists bool   `json:"socket_exists"`
	SocketMode   string `json:"socket_mode"`
	UserInGroup  bool   `json:"user_in_group"`
}

// KubeFacts contains Kubernetes configuration information.
type KubeFacts struct {
	ConfigExists bool          `json:"config_exists"`
	ConfigMode   string        `json:"config_mode"`
	Contexts     []KubeContext `json:"contexts"`
}

// KubeContext describes a single Kubernetes context.
type KubeContext struct {
	Name    string `json:"name"`
	Cluster string `json:"cluster"`
}

// EnvFacts contains information about suspicious environment variables.
type EnvFacts struct {
	SuspiciousVars []SuspiciousVar `json:"suspicious_vars"`
}

// SuspiciousVar records that an environment variable matched a suspicious pattern.
// Only the variable name and matched pattern are stored, NEVER the value.
type SuspiciousVar struct {
	Name    string `json:"name"`
	Pattern string `json:"pattern"`
}

// ShellFacts contains shell configuration information.
type ShellFacts struct {
	Shell           string `json:"shell"`
	HistoryFile     string `json:"history_file"`
	HistoryFileMode string `json:"history_file_mode"`
	HistControl     string `json:"histcontrol"`
}

// ToolsFacts contains information about installed developer tools.
type ToolsFacts struct {
	Tools []ToolInfo `json:"tools"`
}

// ToolInfo describes a single tool.
type ToolInfo struct {
	Name      string `json:"name"`
	Installed bool   `json:"installed"`
	Version   string `json:"version"`
	Path      string `json:"path"`
}

// PathFacts contains information about PATH directories.
type PathFacts struct {
	Dirs []PathDir `json:"dirs"`
}

// PathDir describes a single directory in PATH.
type PathDir struct {
	Path       string `json:"path"`
	Exists     bool   `json:"exists"`
	Writable   bool   `json:"writable"`
	IsRelative bool   `json:"is_relative"`
}

// NewFacts returns a Facts with SchemaVersion set to "1.0".
func NewFacts() Facts {
	return Facts{
		SchemaVersion: "1.0",
	}
}
