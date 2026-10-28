package xray

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

func fakeCoreScript(version string) string {
	return "#!/bin/sh\necho \"Xray " + version + " (fixture)\"\n"
}

func writeFakeCore(t *testing.T, path, version string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(fakeCoreScript(version)), 0o755); err != nil {
		t.Fatalf("write fake core: %v", err)
	}
}

func bundledCoreVersion(t *testing.T) string {
	t.Helper()
	p := newProcess(&Config{})
	p.refreshVersion()
	return p.GetXrayVersion()
}

func useInstallFolders(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("shell fixture is Unix-only")
	}
	t.Setenv("XUI_BIN_FOLDER", t.TempDir())
	t.Setenv("XUI_DB_FOLDER", t.TempDir())
}

// A recreated container brings the image's bundled core back into the bin
// folder; the core picked in the panel must win over it on the next start.
func TestInstalledCoreSurvivesBundledCoreReplacement(t *testing.T) {
	useInstallFolders(t)
	writeFakeCore(t, GetBinaryPath(), "26.10.21")

	if err := InstallBinary(strings.NewReader(fakeCoreScript("25.12.8"))); err != nil {
		t.Fatalf("InstallBinary: %v", err)
	}
	if got := bundledCoreVersion(t); got != "25.12.8" {
		t.Fatalf("version after install = %q, want 25.12.8", got)
	}

	writeFakeCore(t, GetBinaryPath(), "26.10.21")
	if err := RestoreInstalledBinary(); err != nil {
		t.Fatalf("RestoreInstalledBinary: %v", err)
	}
	if got := bundledCoreVersion(t); got != "25.12.8" {
		t.Fatalf("version after image replaced the bin folder = %q, want 25.12.8", got)
	}
}

func TestRestoreInstalledBinaryKeepsBundledCoreWhenNothingInstalled(t *testing.T) {
	useInstallFolders(t)
	writeFakeCore(t, GetBinaryPath(), "26.10.21")

	if err := RestoreInstalledBinary(); err != nil {
		t.Fatalf("RestoreInstalledBinary: %v", err)
	}
	if got := bundledCoreVersion(t); got != "26.10.21" {
		t.Fatalf("version = %q, want bundled 26.10.21", got)
	}
}
