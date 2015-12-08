# Examples

## Registry

For these examples, we will use a custom registry to publish our image.

From the official [Docker documentation](https://docs.docker.com/registry/), in order to have a quick prototype, please launch your registry inside a container as follow:

`docker run -d -p 5000:5000 --name registry registry:2`

 > **NOTE**: If you need to clean or remove your local registry, execute the following command:
 > `docker stop registry && docker rm -v registry`

## Testing

If you want to test the system, you can use `go-hello` example by executing these command:

`cd go-hello && crowley-pack`

However, if the crowley-pack's binary is not present in your user's `PATH`, you can use the following workaround:

```
cd pack/
make && cp pack /tmp/crowley-pack
chmod u+x /tmp/crowley-pack
cd examples/go-hello
/tmp/crowley-pack
```

> **NOTE**: Your user must have access to your Docker daemon.

Then, remove all reference to the image built locally:

```
docker rmi localhost:5000/user/go-hello
docker rmi user/go-hello
```

Finally use your local registry to pull back the image built previously and run it:

```
docker run --rm -it localhost:5000/user/go-hello
```
