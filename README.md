# pack

A Docker build system.

## Introduction

Pack is a build system for Docker images. It use three modules to perform the provisioning of your Docker image:

 * **install**
 * **compose**
 * **publish**

Each module is run as a sequential processing. If a module fails, the build process is aborted with a system exit.

### Install

The **install** module export a file _(a binary or an archive, for example)_ using a Docker container. This container _(called a packer)_ will run a command in order to create the required file for the **compose** module.

> **NOTE:** You can disable this module if you don't require a generated file for compose.

For further informations, please read the **Configuration** and **Packer** paragraph.

### Compose

The **compose** module create a Docker images using the current working directory as _"context"_ with its `Dockerfile`. In order to maintain a clean Docker environment, if an image exists with the same name and tag _(or repository)_, it will be removed.

Furthermore, a _packer_ is an image used to build your project with your development dependencies, in contrast with an image created with _compose_ that should only contains your - minimal - runtime dependencies.

### Publish

The **publish** module will push your image _(from the previous process: compose)_ on a docker registry.

> **NOTE:** It will use the same repository name defined in compose for the remote registry.

If your registry require an authentication, please use `docker login` in order to expose your credentials to the docker's daemon. For further informations, please visit the official [Docker documentation](https://docs.docker.com/engine/reference/commandline/login/)

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

### Compose

```yaml
compose:
    name: user/app
    no-cache: false
    pull: false
```

### Publish

```yaml
publish:
    hostname: docker.io
```

* **name**: Repository name (and optionally a tag) for the image.
* **no-cache**: Do not use cache when building the image.
* **pull**: Always attempt to pull a newer version of the image.

## Packer

> // TODO
