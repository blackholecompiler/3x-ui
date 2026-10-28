package xray

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mhsanaei/3x-ui/v3/internal/config"
)

const maxBinaryBytes = 200 << 20

// GetInstalledBinaryPath returns the persisted copy of the core last installed from the panel.
func GetInstalledBinaryPath() string {
	return filepath.Join(config.GetInstalledXrayFolderPath(), GetBinaryName())
}

// InstallBinary keeps core as the panel-installed copy, then puts it in place of the bundled one.
func InstallBinary(core io.Reader) error {
	if err := writeExecutableAtomic(GetInstalledBinaryPath(), core); err != nil {
		return err
	}
	return RestoreInstalledBinary()
}

// RestoreInstalledBinary copies the panel-installed core over the bundled one rather than
// running it in place: xray-core resolves geo files, certs and Lua beside its executable.
func RestoreInstalledBinary() error {
	core, err := os.Open(GetInstalledBinaryPath())
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer core.Close()
	return writeExecutableAtomic(GetBinaryPath(), core)
}

func writeExecutableAtomic(path string, content io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".xray-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	ok := false
	defer func() {
		_ = tmp.Close()
		if !ok {
			_ = os.Remove(tmpPath)
		}
	}()
	n, err := io.Copy(tmp, io.LimitReader(content, maxBinaryBytes+1))
	if err != nil {
		return err
	}
	if n > maxBinaryBytes {
		return fmt.Errorf("xray binary exceeds %d bytes", maxBinaryBytes)
	}
	if err := tmp.Chmod(0o755); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		_ = os.Remove(path)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}
	ok = true
	return nil
}
