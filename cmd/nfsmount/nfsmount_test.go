//go:build darwin && !cmount
// +build darwin,!cmount

package nfsmount

import (
	"testing"

	"github.com/Evengard/rclone/vfs/vfscommon"
	"github.com/Evengard/rclone/vfs/vfstest"
)

func TestMount(t *testing.T) {
	vfstest.RunTests(t, false, vfscommon.CacheModeMinimal, false, mount)
}
