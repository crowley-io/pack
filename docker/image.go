package docker

import (
	api "github.com/fsouza/go-dockerclient"
)

// Image options
const (
	rmiForce   = false
	rmiNoPrune = false
)

func (d docker) ImageID(name string) string {

	i, err := d.client.InspectImage(name)

	if err != nil {
		return ""
	}

	return i.ID
}

func (d docker) RemoveImage(name string) error {
	if i, _ := d.client.InspectImage(name); i != nil {
		if err := d.client.RemoveImageExtended(name, removeImageOptions()); err != nil {
			return err
		}
	}
	return nil
}

func removeImageOptions() api.RemoveImageOptions {
	return api.RemoveImageOptions{
		Force:   rmiForce,
		NoPrune: rmiNoPrune,
	}
}
