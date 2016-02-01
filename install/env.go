package install

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"

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

	e := expandEnv(configuration.Install.Environment)
	p := path.Clean(configuration.Install.Path)
	o := path.Clean(configuration.Install.Output)

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

func expandEnv(list []string) []string {

	var env []string

	for _, e := range list {

		p := strings.SplitN(e, "=", 2)

		if len(p) == 2 {

			k := p[0]
			v := p[1]

			// Ignore environment variables replacement if value start and finish by
			// single quote.
			if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {

				v = v[1 : len(v)-1]
				e = fmt.Sprintf("%s=%s", k, v)

			} else {
				e = os.ExpandEnv(e)
			}

		}

		env = append(env, e)
	}

	return env
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
