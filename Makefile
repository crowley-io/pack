VERSION=0.1.2

GITHUB_USER="crowley-io"
GITHUB_REPO="pack"

ARTIFACTS = \
	crowley-pack_linux-amd64

UPLOAD_CMD = github-release upload --user ${GITHUB_USER} --repo ${GITHUB_REPO} --tag "v${VERSION}" \
	--name ${FILE} --file ${FILE}

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

release: artifacts
	git tag "v${VERSION}" && git push --tags
	github-release release --user ${GITHUB_USER} --repo ${GITHUB_REPO} --tag "v${VERSION}" \
		--name ${VERSION} --pre-release
	$(foreach FILE,$(ARTIFACTS),$(UPLOAD_CMD);)

artifacts:
	gox -osarch="linux/amd64" -output="crowley-pack_{{.OS}}-{{.Arch}}"

.PHONY: clean artifacts install test style lint pack release
