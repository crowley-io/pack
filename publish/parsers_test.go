package publish

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRepositoryTag(t *testing.T) {

	repo := "localhost.localdomain:5000/foo/bar:latest"
	name := "localhost.localdomain:5000/foo/bar"
	tag := "latest"

	n, i := parseRepositoryTag(repo)

	assert.Equal(t, name, n)
	assert.Equal(t, tag, i)

}

func TestParseRepositoryNoTag(t *testing.T) {

	repo := "localhost.localdomain:5000/foo/bar"
	name := "localhost.localdomain:5000/foo/bar"
	tag := ""

	n, i := parseRepositoryTag(repo)

	assert.Equal(t, name, n)
	assert.Equal(t, tag, i)

}

func TestParseRepositoryNoPort(t *testing.T) {

	repo := "localhost.localdomain/foo/bar"
	name := "localhost.localdomain/foo/bar"
	tag := ""

	n, i := parseRepositoryTag(repo)

	assert.Equal(t, name, n)
	assert.Equal(t, tag, i)

}

func TestParseRepositoryDigest(t *testing.T) {

	repo := "localhost:5000/foo/bar@sha256:bc8813ea7b3603864987522f02a76101c17ad122e1c46d790efc0fca78ca7bfb"
	name := "localhost:5000/foo/bar"
	tag := "sha256:bc8813ea7b3603864987522f02a76101c17ad122e1c46d790efc0fca78ca7bfb"

	n, i := parseRepositoryTag(repo)

	assert.Equal(t, name, n)
	assert.Equal(t, tag, i)

}
