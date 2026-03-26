package schema_test

import (
	"encoding/json"
	"testing"

	"github.com/franklinkim/bouncer/pkg/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFacts_RoundTrip(t *testing.T) {
	facts := schema.NewFacts()
	facts.OS = schema.OSFacts{
		OS:       "linux",
		Arch:     "amd64",
		Hostname: "devbox",
	}
	facts.SSH = schema.SSHFacts{
		DirectoryExists: true,
		DirectoryMode:   "0700",
		Keys: []schema.SSHKey{
			{Path: "/home/user/.ssh/id_ed25519", Mode: "0600", Type: "ed25519"},
		},
		ConfigExists: true,
	}
	facts.Git = schema.GitFacts{
		Installed:        true,
		Version:          "2.43.0",
		CredentialHelper: "osxkeychain",
		SigningEnabled:   true,
		SigningFormat:    "ssh",
	}
	facts.Docker = schema.DockerFacts{
		Installed:    true,
		SocketExists: true,
		SocketMode:   "0660",
		UserInGroup:  true,
	}
	facts.Kube = schema.KubeFacts{
		ConfigExists: true,
		ConfigMode:   "0600",
		Contexts: []schema.KubeContext{
			{Name: "prod", Cluster: "prod-cluster"},
		},
	}
	facts.Env = schema.EnvFacts{
		SuspiciousVars: []schema.SuspiciousVar{
			{Name: "AWS_SECRET_ACCESS_KEY", Pattern: "AWS_SECRET"},
		},
	}
	facts.Shell = schema.ShellFacts{
		Shell:           "/bin/zsh",
		HistoryFile:     "/home/user/.zsh_history",
		HistoryFileMode: "0600",
		HistControl:     "ignoreboth",
	}
	facts.Tools = schema.ToolsFacts{
		Tools: []schema.ToolInfo{
			{Name: "gpg", Installed: true, Version: "2.4.0", Path: "/usr/bin/gpg"},
		},
	}
	facts.Path = schema.PathFacts{
		Dirs: []schema.PathDir{
			{Path: "/usr/local/bin", Exists: true, Writable: false, IsRelative: false},
			{Path: "bin", Exists: true, Writable: true, IsRelative: true},
		},
	}

	data, err := json.Marshal(facts)
	require.NoError(t, err)

	var decoded schema.Facts

	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, facts, decoded)
	assert.Equal(t, "1.0", decoded.SchemaVersion)
}

func TestScanResult_RoundTrip(t *testing.T) {
	results := []schema.CheckResult{
		{
			RuleID:      "SSH-001",
			Title:       "SSH directory permissions",
			Severity:    schema.SeverityCritical,
			Category:    "ssh",
			Status:      schema.StatusFail,
			Evidence:    "mode is 0755, expected 0700",
			Remediation: "chmod 700 ~/.ssh",
		},
		{
			RuleID:   "GIT-001",
			Title:    "Git is installed",
			Severity: schema.SeverityInfo,
			Category: "git",
			Status:   schema.StatusPass,
		},
	}

	sr := schema.NewScanResult(results)

	data, err := json.Marshal(sr)
	require.NoError(t, err)

	var decoded schema.ScanResult

	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, sr, decoded)
	assert.Equal(t, "1.0", decoded.SchemaVersion)
}

func TestNewScanResult_Summary(t *testing.T) {
	results := []schema.CheckResult{
		{RuleID: "R1", Severity: schema.SeverityCritical, Status: schema.StatusFail},
		{RuleID: "R2", Severity: schema.SeverityHigh, Status: schema.StatusFail},
		{RuleID: "R3", Severity: schema.SeverityWarn, Status: schema.StatusPass},
		{RuleID: "R4", Severity: schema.SeverityInfo, Status: schema.StatusPass},
		{RuleID: "R5", Severity: schema.SeverityInfo, Status: schema.StatusSkip},
		{RuleID: "R6", Severity: schema.SeverityCritical, Status: schema.StatusError},
	}

	sr := schema.NewScanResult(results)

	assert.Equal(t, 6, sr.Summary.Total)
	assert.Equal(t, 2, sr.Summary.Passed)
	assert.Equal(t, 2, sr.Summary.Failed)
	assert.Equal(t, 1, sr.Summary.Skipped)
	assert.Equal(t, 1, sr.Summary.Errors)
	assert.Equal(t, 2, sr.Summary.Critical)
	assert.Equal(t, 1, sr.Summary.High)
	assert.Equal(t, 1, sr.Summary.Warn)
	assert.Equal(t, 2, sr.Summary.Info)
}
