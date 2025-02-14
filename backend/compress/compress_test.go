// Test Crypt filesystem interface
package compress

import (
	"os"
	"path/filepath"
	"testing"

	_ "github.com/Evengard/rclone/backend/drive"
	_ "github.com/Evengard/rclone/backend/local"
	_ "github.com/Evengard/rclone/backend/s3"
	_ "github.com/Evengard/rclone/backend/swift"
	"github.com/Evengard/rclone/fstest"
	"github.com/Evengard/rclone/fstest/fstests"
)

var defaultOpt = fstests.Opt{
	RemoteName: "TestCompress:",
	NilObject:  (*Object)(nil),
	UnimplementableFsMethods: []string{
		"OpenWriterAt",
		"OpenChunkWriter",
		"MergeDirs",
		"DirCacheFlush",
		"PutUnchecked",
		"PutStream",
		"UserInfo",
		"Disconnect",
	},
	TiersToTest:                  []string{"STANDARD", "STANDARD_IA"},
	UnimplementableObjectMethods: []string{},
}

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &defaultOpt)
}

// TestRemoteGzip tests GZIP compression
func TestRemoteGzip(t *testing.T) {
	if *fstest.RemoteName != "" {
		t.Skip("Skipping as -remote set")
	}
	tempdir := filepath.Join(os.TempDir(), "rclone-compress-test-gzip")
	name := "TestCompressGzip"
	opt := defaultOpt
	opt.RemoteName = name + ":"
	opt.ExtraConfig = []fstests.ExtraConfigItem{
		{Name: name, Key: "type", Value: "compress"},
		{Name: name, Key: "remote", Value: tempdir},
		{Name: name, Key: "compression_mode", Value: "gzip"},
	}
	opt.QuickTestOK = true
	fstests.Run(t, &opt)
}
