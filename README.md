# pack

A Docker build system.

## Introduction

Pack is a build system for Docker images. It use three modules to perform the provisioning of your Docker image:

 * **install**
 * **compose**
 * **publish**

Each module is run as a sequential processing. If a module fails, the build process is aborted with a system exit.

### Install

The **install** module export a file _(a binary or an archive, for example)_ using a Docker container. This container _(called a packer)_ will run a command in order to create the required file for the **compose** module _(and its Dockerfile)_.

The packer image must contains your development dependencies to execute your project's build and follow a guideline _(see **Packer** paragraph)_, in contrast with the image created with _compose_ which should only contains your minimal runtime dependencies.

> **NOTE:** You can disable this module if you don't require a generated file for compose.

For further informations, please read the **Configuration** and **Packer** paragraph.

### Compose

The **compose** module create a Docker images using the current working directory as _"context"_ with its `Dockerfile`. Also, its highly recommended to keep a lightweight image with only your runtime dependencies in a minimal way.

> **NOTE:** In order to maintain a clean Docker environment, if an image exists with the same name and tag _(or repository)_, it will be removed.

### Publish

The **publish** module will push your image _(from the previous process: compose)_ on a docker registry.

> **NOTE:** It will use the same repository name defined in compose for the remote registry.

If your registry require an authentication, please use `docker login` in order to expose your credentials to the docker's daemon. For further informations, please visit the official [Docker documentation](https://docs.docker.com/engine/reference/commandline/login/)

## Usage

```
Usage: crowley-pack [OPTIONS]

Docker build system.

Options:
  -f, --file="packer.yml"   Configuration file
```

**Example:**

`cd my-app && crowley-pack`

## Configuration

### Install

```yaml
install:
    disable: false
    image: user/app-packer
    path: /usr/local/app
    output: app.tar.gz
    volumes:
        - ~/.ssh:/home/user/.ssh
        - /media/hd0:/usr/local/data
    environment:
        - API_URI=http://api.example.com
        - API_ACCESS_KEY='3ztP7$Xqoef=VUdPa'
        - API_SECRET_KEY=$SECRET_KEY
```

* **disable**: Do not perform install process.
* **image**: Image used to perform the install process _(a packer)_.
* **path**: Path inside the container which will contains your working directory.
* **output**: Filename used as output when your project build is successful.
* **volumes**: Volumes to mount for the running container.
* **environment**: Environment variables to inject for the running container.

### Compose

```yaml
compose:
    name: user/app
    no-cache: false
    pull: false
```

* **name**: Repository name (and optionally a tag) for the image.
* **no-cache**: Do not use cache when building the image.
* **pull**: Always attempt to pull a newer version of the image.

### Publish

```yaml
publish:
    hostname: docker.io
```

* **hostname**: Registry hostname (with a optional port).

## Packer

A packer image should respect the following guidelines to ensure that the building process runs smoothly without any issue.

First of all, the container receives the following environment variables:

 * `CROWLEY_PACK_USER`
 * `CROWLEY_PACK_GROUP`
 * `CROWLEY_PACK_DIRECTORY`
 * `CROWLEY_PACK_OUTPUT`

These environment variables define a guideline configuration for the build process such as the user and the group of the output file, its path and also its working directory.

For example, you can define the following minimal `Dockerfile`:

```dockerfile
# Dockerfile for a generic packer
FROM debian:jessie

RUN apt-get update && apt-get install -y build-essential
ADD pack /usr/local/bin/
RUN chmod +x /usr/local/bin/pack
CMD pack
```

Using this script as `pack`:

```bash
#!/bin/bash

# Fail hard and fast
set -eo pipefail

# First, change the working directory
cd ${CROWLEY_PACK_DIRECTORY}

# Then, make will build our binary in $CROWLEY_PACK_OUTPUT
make

# Finally, update output owner since we may run as root user...
chown ${CROWLEY_PACK_USER}:${CROWLEY_PACK_GROUP} ${CROWLEY_PACK_OUTPUT}
```

This should provide you a minimal framework on how to build your packer image.
