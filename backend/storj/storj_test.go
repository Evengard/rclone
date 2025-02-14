//go:build !plan9
// +build !plan9

// Test Storj filesystem interface
package storj_test

import (
	"testing"

	"github.com/Evengard/rclone/backend/storj"
	"github.com/Evengard/rclone/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName: "TestStorj:",
		NilObject:  (*storj.Object)(nil),
	})
}
