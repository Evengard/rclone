package tree

import (
	"bytes"
	"context"
	"testing"

	_ "github.com/Evengard/rclone/backend/local"
	"github.com/Evengard/rclone/fs"
	"github.com/Evengard/rclone/fstest"
	"github.com/a8m/tree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTree(t *testing.T) {
	fstest.Initialise()

	buf := new(bytes.Buffer)

	f, err := fs.NewFs(context.Background(), "testfiles")
	require.NoError(t, err)
	err = Tree(f, buf, new(tree.Options))
	require.NoError(t, err)
	assert.Equal(t, `/
├── file1
├── file2
├── file3
└── subdir
    ├── file4
    └── file5

1 directories, 5 files
`, buf.String())
}
