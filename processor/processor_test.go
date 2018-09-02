package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtension(t *testing.T) {
	t.Parallel()

	result := ext("foobar.png")

	assert.Equal(t, "foobar.jpeg", result)
}

func TestFilenameProcessor(t *testing.T) {
	t.Parallel()

	file := filename("/foo/bar/images/abcdef.png", "/foo/bar/thumbnails")

	assert.Equal(t, "/foo/bar/thumbnails/abcdef.jpeg", file)
}
