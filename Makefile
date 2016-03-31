VERSION=0.1.3
NAME=pack
ORGANIZATION=crowley-io
PACKAGE=github.com/${ORGANIZATION}/${NAME}

GITHUB_USER=${ORGANIZATION}
GITHUB_REPO=${NAME}

ARTIFACTS = \
	crowley-${NAME}_linux-amd64

UPLOAD_CMD = github-release upload --user ${GITHUB_USER} --repo ${GITHUB_REPO} --tag "v${VERSION}" \
	--name ${FILE} --file ${FILE}

all: ${NAME}

setup:
	go get -u github.com/mitchellh/gox
	go get -u github.com/alecthomas/gometalinter

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
	go build -o ${NAME}

clean:
	rm -rf ${NAME}

install: ${NAME}
	install -o root -g root -m 0755 ${NAME} /usr/local/bin/crowley-${NAME}

release: artifacts
	git tag "v${VERSION}" && git push --tags
	github-release release --user ${GITHUB_USER} --repo ${GITHUB_REPO} --tag "v${VERSION}" \
		--name ${VERSION} --pre-release
	$(foreach FILE,$(ARTIFACTS),$(UPLOAD_CMD);)

artifacts:
	gox -osarch="linux/amd64" -output="crowley-${NAME}_{{.OS}}-{{.Arch}}"

.PHONY: clean ${NAME} install artifacts test style lint coverage release
