// Test Sia filesystem interface
package sia_test

import (
	"testing"

	"github.com/Evengard/rclone/backend/sia"

	"github.com/Evengard/rclone/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestSia:",
		NilObject:  (*sia.Object)(nil),
	})
}
