package install

import (
	"fmt"
	"path"
	"testing"

	"github.com/crowley-io/pack/configuration"
	"github.com/stretchr/testify/assert"
)

func TestArchivePath(t *testing.T) {

	p := "/media//archive.tar.gz"
	e := "/media/archive.tar.gz"
	c := path.Clean(p)

	assert.Equal(t, e, c)

}

func TestOutputPath(t *testing.T) {

	p := "/media/"
	e := "/media"

	c := path.Clean(p)

	assert.Equal(t, e, c)

}

func TestArchiveDirectory(t *testing.T) {

	p := "/media/archive.tar.gz"
	e := "/media"

	d := path.Dir(path.Clean(p))

	assert.Equal(t, e, d)

}

func TestParseVolumeWithoutMode(t *testing.T) {

	p := "/home/user/.npm:/root/.npm"

	v, err := parseVolumePath(p)

	assert.Nil(t, err)
	assert.Equal(t, p, v)

}

func TestParseVolumeWithMode(t *testing.T) {

	p := "/home/user/.npm:/root/.npm:rw"

	v, err := parseVolumePath(p)

	assert.Nil(t, err)
	assert.Equal(t, p, v)

}

func TestParseVolumeWithoutExternal(t *testing.T) {

	p := "/root/.npm"

	v, err := parseVolumePath(p)

	assert.Nil(t, err)
	assert.Equal(t, p, v)

}

func TestParseVolumeWithSyntaxError(t *testing.T) {

	p := "/home/user/.npm:/root/.npm:rw:3"

	v, err := parseVolumePath(p)

	assert.NotNil(t, err)
	assert.Empty(t, v)

}

func TestParseVolumeWithInternalRelativeError(t *testing.T) {

	p := "/home/user/.npm:./foo:rw"

	v, err := parseVolumePath(p)

	assert.NotNil(t, err)
	assert.Empty(t, v)

}

func TestParseVolumeWithHomeResolve(t *testing.T) {

	p := "~/.npm:/root/.npm"
	e := fmt.Sprintf("%s/.npm:/root/.npm", home())

	v, err := parseVolumePath(p)

	assert.Nil(t, err)
	assert.Equal(t, e, v)

}

func TestParseVolumeWithCleanHomeResolve(t *testing.T) {

	p := "~/:/root/"
	e := fmt.Sprintf("%s:/root", home())

	v, err := parseVolumePath(p)

	assert.Nil(t, err)
	assert.Equal(t, e, v)

}

func TestParseVolumeWithRelativeResolve(t *testing.T) {

	p := "./bin/:/usr/local/bin/:ro"
	e := fmt.Sprintf("%s/bin:/usr/local/bin:ro", pwd())

	v, err := parseVolumePath(p)

	assert.Nil(t, err)
	assert.Equal(t, e, v)

}

func TestParseVolumeWithRelativeResolve2(t *testing.T) {

	p := "./:/var/www"
	e := fmt.Sprintf("%s:/var/www", pwd())

	v, err := parseVolumePath(p)

	assert.Nil(t, err)
	assert.Equal(t, e, v)

}

func TestParseVolumeWithRelativeResolve3(t *testing.T) {

	p := "bin/app:/var/www:rw"
	e := fmt.Sprintf("%s/bin/app:/var/www:rw", pwd())

	v, err := parseVolumePath(p)

	assert.Nil(t, err)
	assert.Equal(t, e, v)

}

func TestParseVolumeWithRelativeResolve4(t *testing.T) {

	p := "~bin/app:/var/www:rw"
	e := fmt.Sprintf("%s/~bin/app:/var/www:rw", pwd())

	v, err := parseVolumePath(p)

	assert.Nil(t, err)
	assert.Equal(t, e, v)

}

func TestParseVolumeWithCleanRelativeResolve(t *testing.T) {

	p := "../bin/../app:/var/"
	e := fmt.Sprintf("%s/app:/var", path.Dir(pwd()))

	v, err := parseVolumePath(p)

	assert.Nil(t, err)
	assert.Equal(t, e, v)

}

func TestParseVolumeWithParentResolve(t *testing.T) {

	p := "../../:/var/"
	e := fmt.Sprintf("%s:/var", path.Dir(path.Dir(pwd())))

	v, err := parseVolumePath(p)

	assert.Nil(t, err)
	assert.Equal(t, e, v)

}

func TestGetVolumes(t *testing.T) {

	p := "/usr/local/bin"
	v1 := "/home/user/.ivy2:/root/.ivy2:rw"
	v2 := "/home/user/.sbt:/root/.sbt:rw"
	vp := fmt.Sprintf("%s:%s:rw", pwd(), p)
	v := []string{v1, v2}

	i, err := GetVolumes(getConfiguration(p, v))

	assert.Nil(t, err)
	assert.Contains(t, i, v1)
	assert.Contains(t, i, v2)
	assert.Contains(t, i, vp)

}

func TestGetVolumesEmptyPath(t *testing.T) {
	v, err := GetVolumes(getConfiguration("", nil))
	assert.NotNil(t, err)
	assert.Empty(t, v)
}

func TestGetVolumesSyntaxError(t *testing.T) {

	p := "/usr/local/bin"
	v := []string{"/home/user/.npm:/root/.npm:ro:x"}
	i, err := GetVolumes(getConfiguration(p, v))

	assert.NotNil(t, err)
	assert.Nil(t, i)

}

func getConfiguration(path string, volumes []string) *configuration.Configuration {
	return &configuration.Configuration{
		Install: configuration.Install{
			Path:    path,
			Volumes: volumes,
		},
	}
}
