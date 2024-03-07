//go:build noselfupdate
// +build noselfupdate

package selfupdate

import (
	"github.com/Evengard/rclone/lib/buildinfo"
)

func init() {
	buildinfo.Tags = append(buildinfo.Tags, "noselfupdate")
}
