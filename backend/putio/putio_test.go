// Test Put.io filesystem interface
package putio

import (
	"testing"

	"github.com/Evengard/rclone/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestPutio:",
		NilObject:  (*Object)(nil),
	})
}
