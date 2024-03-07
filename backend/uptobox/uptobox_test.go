// Test Uptobox filesystem interface
package uptobox_test

import (
	"testing"

	"github.com/Evengard/rclone/backend/uptobox"
	"github.com/Evengard/rclone/fstest"
	"github.com/Evengard/rclone/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	if *fstest.RemoteName == "" {
		*fstest.RemoteName = "TestUptobox:"
	}
	fstests.Run(t, &fstests.Opt{
		RemoteName: *fstest.RemoteName,
		NilObject:  (*uptobox.Object)(nil),
	})
}
