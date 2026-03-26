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
func (c *OSCollector) Collect(ctx context.Context, facts *schema.Facts) Result {
	facts.OS.OS = runtime.GOOS
	facts.OS.Arch = runtime.GOARCH

	hostname, _ := os.Hostname()
	facts.OS.Hostname = hostname

	if runtime.GOOS == "darwin" {
		c.collectDarwin(ctx, facts)
	}

	return Result{Name: c.Name(), Status: StatusOK}
}

// collectDarwin gathers macOS-specific security facts.
func (c *OSCollector) collectDarwin(ctx context.Context, facts *schema.Facts) {
	facts.OS.SIPEnabled = cmdContains(ctx, "csrutil", []string{"status"}, "enabled")
	facts.OS.FileVaultEnabled = cmdContains(ctx, "fdesetup", []string{"status"}, "On")
	facts.OS.FirewallEnabled = cmdContains(ctx, "/usr/libexec/ApplicationFirewall/socketfilterfw", []string{"--getglobalstate"}, "enabled")
	facts.OS.StealthModeEnabled = cmdContains(ctx, "/usr/libexec/ApplicationFirewall/socketfilterfw", []string{"--getstealthmode"}, "enabled")
	facts.OS.GatekeeperEnabled = cmdContains(ctx, "spctl", []string{"--status"}, "enabled")
	facts.OS.AutoLoginDisabled = !cmdContains(ctx, "defaults", []string{"read", "/Library/Preferences/com.apple.loginwindow", "autoLoginUser"}, "")
	facts.OS.GuestAccountDisabled = cmdContains(ctx, "defaults", []string{"read", "/Library/Preferences/com.apple.loginwindow", "GuestEnabled"}, "0")
	facts.OS.ScreenLockTimeoutSecs = readScreenLockTimeout(ctx)
	facts.OS.OSAutoUpdateEnabled = cmdContains(ctx, "defaults", []string{"read", "/Library/Preferences/com.apple.SoftwareUpdate", "AutomaticallyInstallMacOSUpdates"}, "1")
	facts.OS.AppAutoUpdateEnabled = cmdContains(ctx, "defaults", []string{"read", "/Library/Preferences/com.apple.commerce", "AutoUpdate"}, "1")
	facts.OS.RSREnabled = cmdContains(ctx, "defaults", []string{"read", "/Library/Preferences/com.apple.SoftwareUpdate", "ConfigDataInstall"}, "1")
	facts.OS.UserIsStandardAccount = !cmdContains(ctx, "groups", nil, "admin")
	facts.OS.ScreenSharingDisabled = !cmdContains(ctx, "launchctl", []string{"list", "com.apple.screensharing"}, "")
	facts.OS.InternetSharingDisabled = !cmdContains(ctx, "defaults", []string{"read", "/Library/Preferences/SystemConfiguration/com.apple.nat", "Enabled"}, "1")
	facts.OS.PrinterSharingDisabled = cmdContains(ctx, "cupsctl", nil, "_share_printers=0")
	facts.OS.RemoteAppleEventsDisabled = !cmdContains(ctx, "launchctl", []string{"list", "com.apple.AEServer"}, "")
	facts.OS.AirdropSetting = readAirdropSetting(ctx)
	facts.OS.RosettaInstalled = cmdSucceeds(ctx, "pgrep", []string{"-q", "oahd"})
	facts.OS.EDRDeployed = detectEDR(ctx)
	facts.OS.LegacyKextsBlocked = !cmdContains(ctx, "kmutil", []string{"showloaded", "--list-only"}, "com.apple")
	facts.OS.MDMEnrolled = cmdContains(ctx, "profiles", []string{"status", "-type", "enrollment"}, "MDM enrollment: Yes")
}

// cmdContains runs a command and checks if stdout contains a substring.
func cmdContains(ctx context.Context, name string, args []string, substr string) bool {
	out, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	if err != nil {
		return false
	}

	if substr == "" {
		return true
	}

	return strings.Contains(string(out), substr)
}

// cmdSucceeds returns true if the command exits with status 0.
func cmdSucceeds(ctx context.Context, name string, args []string) bool {
	return exec.CommandContext(ctx, name, args...).Run() == nil
}

func readScreenLockTimeout(ctx context.Context) int {
	out, err := exec.CommandContext(ctx, "defaults", "-currentHost", "read", "com.apple.screensaver", "idleTime").CombinedOutput()
	if err != nil {
		return 0
	}

	val, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0
	}

	return val
}

func readAirdropSetting(ctx context.Context) string {
	out, err := exec.CommandContext(ctx, "defaults", "read", "com.apple.sharingd", "DiscoverableMode").CombinedOutput()
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

func detectEDR(ctx context.Context) bool {
	knownAgents := []string{
		"com.crowdstrike.falcon",
		"com.sentinelone",
		"com.carbon.black",
		"com.microsoft.wdav",
	}
	for _, agent := range knownAgents {
		if cmdSucceeds(ctx, "launchctl", []string{"list", agent}) {
			return true
		}
	}

	return false
}
