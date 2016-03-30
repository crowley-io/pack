package install

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/crowley-io/pack/configuration"
)

// GetVolumes return the required volumes path for the container.
func GetVolumes(configuration *configuration.Configuration) ([]string, error) {

	path := configuration.Install.Path
	volumes := configuration.Install.Volumes

	var paths []string
	var p string
	var err error

	if p, err = parseVolumePath(fmt.Sprintf(".:%s:rw", path)); err != nil {
		return nil, err
	}

	paths = append(paths, p)

	for _, v := range volumes {

		if p, err = parseVolumePath(v); err != nil {
			return nil, err
		}

		paths = append(paths, p)

	}

	return paths, nil
}

func parseVolumePath(path string) (string, error) {

	p := strings.Split(path, ":")

	if len(p) > 3 {
		return "", fmt.Errorf("Volume %s has incorrect format, should be external:internal[:mode]", path)
	}

	external := ""
	internal := ""
	mode := ""

	if len(p) == 1 {
		external, internal = resolveVolumePath("", p[0])
	} else {
		external, internal = resolveVolumePath(p[0], p[1])
	}

	if len(p) == 3 {
		mode = p[2]
	}

	if err := isAbsolutePath("internal", internal); err != nil {
		return "", err
	}

	if err := isAbsolutePath("external", external); err != nil && external != "" {
		return "", err
	}

	return formatVolumePath(external, internal, mode), nil
}

func resolveVolumePath(external, internal string) (string, string) {

	if external == "" {
		return "", path.Clean(internal)
	}

	if len(external) >= 2 && external[:2] == "~/" {
		external = strings.Replace(filepath.Clean(external), "~", home(), 1)
	}

	i := path.Clean(internal)
	e, _ := filepath.Abs(filepath.Clean(external))

	return e, i
}

func formatVolumePath(external, internal, mode string) string {

	if external != "" {
		external = fmt.Sprintf("%s:", external)
	}

	if mode != "" {
		mode = fmt.Sprintf(":%s", mode)
	}

	return external + internal + mode
}

func isAbsolutePath(t, p string) error {
	if !path.IsAbs(p) {
		return fmt.Errorf("%s path '%s' is not an absolute path", t, p)
	}

	return nil
}
