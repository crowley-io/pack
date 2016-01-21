package install

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/crowley-io/pack/configuration"
	"github.com/stretchr/testify/assert"
)

func TestWhoami(t *testing.T) {

	uid, gid := whoami()

	assert.NotEmpty(t, uid)
	assert.NotEmpty(t, gid)

	u, err := strconv.ParseUint(uid, 10, 64)

	assert.Nil(t, err)
	assert.True(t, u >= 0)

	g, err := strconv.ParseUint(gid, 10, 64)

	assert.Nil(t, err)
	assert.True(t, g >= 0)

}

func TestHome(t *testing.T) {

	h := home()

	assert.NotEmpty(t, h)
	assert.True(t, pathExist(h), "home path doesn't exist")

}

func TestPwd(t *testing.T) {

	p := pwd()

	assert.NotEmpty(t, p)
	assert.True(t, pathExist(p), "pwd path doesn't exist")

}

func TestGetEnv(t *testing.T) {

	uid, gid := whoami()
	path := os.Getenv("PATH")

	o := "libshaped.so"
	p := "/media/app"
	u := "DB_URI=mongodb://user:password@host:27017/db"
	s := "SECRET_KEY='3ztP7$Xqoef=VUdPa'"
	f := "GOPATH=$PATH:/usr/local/go"

	c := &configuration.Configuration{
		Install: configuration.Install{
			Path:        p,
			Output:      o,
			Environment: []string{u, f, s},
		},
	}

	e, err := GetEnv(c)

	assert.Nil(t, err)
	assert.NotEmpty(t, e)
	assert.Contains(t, e, u)
	assert.Contains(t, e, fmt.Sprintf("CROWLEY_PACK_USER=%s", uid))
	assert.Contains(t, e, fmt.Sprintf("CROWLEY_PACK_GROUP=%s", gid))
	assert.Contains(t, e, fmt.Sprintf("CROWLEY_PACK_DIRECTORY=%s", p))
	assert.Contains(t, e, fmt.Sprintf("CROWLEY_PACK_OUTPUT=%s/%s", p, o))
	assert.Contains(t, e, fmt.Sprintf("GOPATH=%s:/usr/local/go", path))
	assert.Contains(t, e, "SECRET_KEY=3ztP7$Xqoef=VUdPa")

}
