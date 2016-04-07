VERSION=0.1.4
NAME=pack
ORGANIZATION=crowley-io
PACKAGE=github.com/${ORGANIZATION}/${NAME}

GITHUB_USER=${ORGANIZATION}
GITHUB_REPO=${NAME}

ARTIFACTS = \
	crowley-${NAME}_linux-amd64

UPLOAD_CMD = github-release upload --user ${GITHUB_USER} --repo ${GITHUB_REPO} --tag "v${VERSION}" \
	--name ${FILE} --file ${FILE}

LDFLAGS="-X main.Version=v${VERSION}"

all: ${NAME}

setup:
	go get -u github.com/alecthomas/gometalinter
	go get -u github.com/aktau/github-release
	gometalinter --install --update
	git submodule update --init --recursive

test:
	@echo " -> $@"
	@script/test ${PACKAGE}

style:
	@echo " -> $@"
	@script/style ${PACKAGE}

lint:
	@echo " -> $@"
	@script/lint ${PACKAGE}

coverage:
	@echo " -> $@"
	@script/coverage ${PACKAGE}

${NAME}:
	go build -ldflags ${LDFLAGS} -o ${NAME}

clean:
	rm -rf ${NAME}
	rm -rf crowley-${NAME}_linux-amd64

install: ${NAME}
	install -o root -g root -m 0755 ${NAME} /usr/local/bin/crowley-${NAME}

release: artifacts
	git tag "v${VERSION}" && git push --tags
	github-release release --user ${GITHUB_USER} --repo ${GITHUB_REPO} --tag "v${VERSION}" \
		--name ${VERSION} --pre-release
	$(foreach FILE,$(ARTIFACTS),$(UPLOAD_CMD);)

crowley-${NAME}_linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags ${LDFLAGS} -o "crowley-${NAME}_linux-amd64"

artifacts: crowley-${NAME}_linux-amd64

.PHONY: clean ${NAME} install artifacts test style lint coverage release
