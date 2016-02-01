all: pack

setup:
	go get -d -t -v ./...

test: setup
	go test ./...

style:
	gofmt -w .

lint:
	golint ./...

pack: setup
	go build

clean:
	rm -rf pack

install: pack
	install -o root -g root -m 0755 pack /usr/local/bin/crowley-pack

artifacts:
	gox -osarch="linux/amd64" -output="crowley-pack_{{.OS}}-{{.Arch}}"

.PHONY: clean artifacts install test style lint pack
