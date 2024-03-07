//go:build linux
// +build linux

package mount2

import (
	"testing"

	"github.com/Evengard/rclone/vfs/vfscommon"
	"github.com/Evengard/rclone/vfs/vfstest"
)

func TestMount(t *testing.T) {
	vfstest.RunTests(t, false, vfscommon.CacheModeOff, true, mount)
}
