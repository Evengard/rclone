// Test Cache filesystem interface

//go:build !plan9 && !js && !race
// +build !plan9,!js,!race

package cache_test

import (
	"testing"

	"github.com/Evengard/rclone/backend/cache"
	_ "github.com/Evengard/rclone/backend/local"
	"github.com/Evengard/rclone/fstest/fstests"
)

// TestIntegration runs integration tests against the remote
func TestIntegration(t *testing.T) {
	fstests.Run(t, &fstests.Opt{
		RemoteName:                      "TestCache:",
		NilObject:                       (*cache.Object)(nil),
		UnimplementableFsMethods:        []string{"PublicLink", "OpenWriterAt", "OpenChunkWriter", "DirSetModTime", "MkdirMetadata"},
		UnimplementableObjectMethods:    []string{"MimeType", "ID", "GetTier", "SetTier", "Metadata"},
		UnimplementableDirectoryMethods: []string{"Metadata", "SetMetadata", "SetModTime"},
		SkipInvalidUTF8:                 true, // invalid UTF-8 confuses the cache
	})
}
