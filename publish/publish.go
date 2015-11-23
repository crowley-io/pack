package publish

import (
	"fmt"

	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"
)

// Publish push the docker image into the docker registy.
func Publish(client docker.Docker, configuration *configuration.Configuration) error {

	if err := configuration.Validate(); err != nil {
		return err
	}

	name := configuration.Compose.Name
	hostname := configuration.Publish.Hostname

	repository, tag := parseRepositoryTag(fmt.Sprintf("%s/%s", hostname, name))

	option := docker.TagOptions{
		Name:       name,
		Repository: repository,
		Tag:        tag,
	}

	err := client.Tag(option)

	if err != nil {
		return err
	}

	// TODO Push docker image
	// nov. 23 17:33:02 xps-dash docker[404]: time="2015-11-23T17:33:02.643242192+01:00" level=info msg="POST /v1.21/images/localhost:5000/captaindash/influxdb/push?tag="
	// nov. 23 17:33:30 xps-dash docker[404]: time="2015-11-23T17:33:30.567899236+01:00" level=info msg="Signed manifest for captaindash/influxdb:latest using daemon's key: OKBX:SUSA:HQBV:VTSK:U6TO:6PGC:K4ZT:JAW2:IHR7:JLYS:Z743:RDD7"
	// nov. 23 17:33:46 xps-dash docker[404]: time="2015-11-23T17:33:46.369134172+01:00" level=info msg="POST /v1.21/build?buildargs=%7B%7D&cgroupparent=&cpuperiod=0&cpuquota=0&cpusetcpus=&cpusetmems=&cpushares=0&dockerfile=Dockerfile&memory=0&memswap=0&nocache=1&rm=1&t=captaindash%2Fcommon&ulimits=null"

	return nil
}
