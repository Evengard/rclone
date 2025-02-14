// Test Zoho filesystem interface
package zoho_test

import (
	"testing"

	"github.com/Evengard/rclone/backend/zoho"
	"github.com/Evengard/rclone/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestZoho:",
		NilObject:  (*zoho.Object)(nil),
	})
}
