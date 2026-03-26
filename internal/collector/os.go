package collector

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/franklinkim/bouncer/pkg/schema"
)

// OSCollector gathers facts about the operating system.
type OSCollector struct{}

// NewOSCollector returns a new OSCollector.
func NewOSCollector() *OSCollector {
	return &OSCollector{}
}

// Name returns the collector name.
func (c *OSCollector) Name() string {
	return "os"
}

// Collect populates facts.OS with operating system information.
func (c *OSCollector) Collect(_ context.Context, facts *schema.Facts) Result {
	facts.OS.OS = runtime.GOOS
	facts.OS.Arch = runtime.GOARCH

	hostname, _ := os.Hostname()
	facts.OS.Hostname = hostname

	if runtime.GOOS == "darwin" {
		c.collectDarwin(facts)
	}

	return Result{Name: c.Name(), Status: StatusOK}
}

// collectDarwin gathers macOS-specific security facts.
func (c *OSCollector) collectDarwin(facts *schema.Facts) {
	facts.OS.SIPEnabled = cmdContains("csrutil", []string{"status"}, "enabled")
	facts.OS.FileVaultEnabled = cmdContains("fdesetup", []string{"status"}, "On")
	facts.OS.FirewallEnabled = cmdContains("/usr/libexec/ApplicationFirewall/socketfilterfw", []string{"--getglobalstate"}, "enabled")
	facts.OS.StealthModeEnabled = cmdContains("/usr/libexec/ApplicationFirewall/socketfilterfw", []string{"--getstealthmode"}, "enabled")
	facts.OS.GatekeeperEnabled = cmdContains("spctl", []string{"--status"}, "enabled")
	facts.OS.AutoLoginDisabled = !cmdContains("defaults", []string{"read", "/Library/Preferences/com.apple.loginwindow", "autoLoginUser"}, "")
	facts.OS.GuestAccountDisabled = cmdContains("defaults", []string{"read", "/Library/Preferences/com.apple.loginwindow", "GuestEnabled"}, "0")
	facts.OS.ScreenLockTimeoutSecs = readScreenLockTimeout()
	facts.OS.OSAutoUpdateEnabled = cmdContains("defaults", []string{"read", "/Library/Preferences/com.apple.SoftwareUpdate", "AutomaticallyInstallMacOSUpdates"}, "1")
	facts.OS.AppAutoUpdateEnabled = cmdContains("defaults", []string{"read", "/Library/Preferences/com.apple.commerce", "AutoUpdate"}, "1")
	facts.OS.RSREnabled = cmdContains("defaults", []string{"read", "/Library/Preferences/com.apple.SoftwareUpdate", "ConfigDataInstall"}, "1")
	facts.OS.UserIsStandardAccount = !cmdContains("groups", nil, "admin")
	facts.OS.ScreenSharingDisabled = !cmdContains("launchctl", []string{"list", "com.apple.screensharing"}, "")
	facts.OS.InternetSharingDisabled = !cmdContains("defaults", []string{"read", "/Library/Preferences/SystemConfiguration/com.apple.nat", "Enabled"}, "1")
	facts.OS.PrinterSharingDisabled = !cmdContains("cupsctl", nil, "share_printers=0")
	facts.OS.RemoteAppleEventsDisabled = !cmdContains("launchctl", []string{"list", "com.apple.AEServer"}, "")
	facts.OS.AirdropSetting = readAirdropSetting()
	facts.OS.RosettaInstalled = cmdSucceeds("pgrep", []string{"-q", "oahd"})
	facts.OS.EDRDeployed = detectEDR()
	facts.OS.LegacyKextsBlocked = !cmdContains("kmutil", []string{"showloaded", "--list-only"}, "com.apple")
	facts.OS.MDMEnrolled = cmdContains("profiles", []string{"status", "-type", "enrollment"}, "MDM enrollment: Yes")
}

// cmdContains runs a command and checks if stdout contains a substring.
func cmdContains(name string, args []string, substr string) bool {
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		return false
	}
	if substr == "" {
		return true
	}
	return strings.Contains(string(out), substr)
}

// cmdSucceeds returns true if the command exits with status 0.
func cmdSucceeds(name string, args []string) bool {
	return exec.Command(name, args...).Run() == nil
}

func readScreenLockTimeout() int {
	out, err := exec.Command("defaults", "-currentHost", "read", "com.apple.screensaver", "idleTime").CombinedOutput()
	if err != nil {
		return 0
	}
	val, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0
	}
	return val
}

func readAirdropSetting() string {
	out, err := exec.Command("defaults", "read", "com.apple.sharingd", "DiscoverableMode").CombinedOutput()
	if err != nil {
		return "off"
	}
	switch strings.TrimSpace(string(out)) {
	case "Everyone":
		return "everyone"
	case "Contacts Only":
		return "contacts_only"
	default:
		return "off"
	}
}

func detectEDR() bool {
	knownAgents := []string{
		"com.crowdstrike.falcon",
		"com.sentinelone",
		"com.carbon.black",
		"com.microsoft.wdav",
	}
	for _, agent := range knownAgents {
		if cmdSucceeds("launchctl", []string{"list", agent}) {
			return true
		}
	}
	return false
}
