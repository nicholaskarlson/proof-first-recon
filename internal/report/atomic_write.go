package report

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// writeFileAtomic writes data to path by writing to a temp file in the same
// directory and then renaming it into place. This prevents partial writes.
// On Windows, os.Rename cannot replace an existing file, so we remove the
// destination first as a best-effort overwrite-safe behavior.
func writeFileAtomic(path string, data []byte, perm fs.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	base := filepath.Base(path)
	tmp, err := os.CreateTemp(dir, "."+base+".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	cleanup := func() {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
	}
	defer cleanup()

	if err := tmp.Chmod(perm); err != nil {
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpName, path); err == nil {
		return nil
	} else {
		// Best-effort Windows overwrite support.
		if err2 := os.Remove(path); err2 != nil && !errors.Is(err2, os.ErrNotExist) {
			return err
		}
		return os.Rename(tmpName, path)
	}
}
