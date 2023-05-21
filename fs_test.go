package aurum

import (
	"path/filepath"
	"testing"

	"github.com/hansmi/aurum/internal/testutil"
)

func TestWritableDirFS(t *testing.T) {
	tmpdir := t.TempDir()

	d := newWritableDirFS(tmpdir)

	if err := d.WriteFile("test1", nil, 0o644); err != nil {
		t.Errorf("WriteFile() failed: %v", err)
	}

	testutil.MustLstat(t, filepath.Join(tmpdir, "test1"))
}
