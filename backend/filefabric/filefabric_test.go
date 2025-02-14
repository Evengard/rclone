// Test filefabric filesystem interface
package filefabric_test

import (
	"testing"

	"github.com/Evengard/rclone/backend/filefabric"
	"github.com/Evengard/rclone/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestFileFabric:",
		NilObject:  (*filefabric.Object)(nil),
	})
}
