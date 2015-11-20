package install

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/crowley-io/pack/configuration"
)

const (
	// UserEnv is the environment variable name used to injected user id.
	UserEnv = "CROWLEY_PACK_USER"
	// GroupEnv is the environment variable name used to injected user's group id.
	GroupEnv = "CROWLEY_PACK_GROUP"
	// DirectoryEnv is the environment variable name used to define the working
	// directory output inside the container.
	DirectoryEnv = "CROWLEY_PACK_DIRECTORY"
	// OutputEnv is the environment variable name used to define the output
	// file's path inside the container.
	OutputEnv = "CROWLEY_PACK_OUTPUT"
)

// GetEnv return the required environment variables for the container.
func GetEnv(configuration *configuration.Configuration) ([]string, error) {

	e := configuration.Install.Environment
	p := path.Clean(configuration.Install.Path)
	o := path.Clean(configuration.Output)

	e = addUserEnv(e)
	e = addPathEnv(e, p, path.Join(p, o))

	return e, nil
}

// Return user's uid and gid
func whoami() (uid, gid string) {
	if u, err := user.Current(); err == nil {
		uid = u.Uid
		gid = u.Gid
	}
	return
}

// Return user's home directory
func home() (h string) {
	if u, err := user.Current(); err == nil {
		h = u.HomeDir
	}
	return
}

// Return user's working directory
func pwd() string {
	p, _ := os.Getwd()
	return p
}

func addUserEnv(env []string) []string {

	uid, gid := whoami()
	env = append(env, fmt.Sprintf("%s=%s", UserEnv, uid))
	env = append(env, fmt.Sprintf("%s=%s", GroupEnv, gid))

	return env
}

func addPathEnv(env []string, directory, output string) []string {

	env = append(env, fmt.Sprintf("%s=%s", DirectoryEnv, directory))
	env = append(env, fmt.Sprintf("%s=%s", OutputEnv, output))

	return env
}
