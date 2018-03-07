BUILD_CONTAINER := joemiller/go-jail-build

create-build-container:
	@docker build -t $(BUILD_CONTAINER) -f Dockerfile.build .

devshell: create-build-container
	#@docker build -t joemiller/go-jail-devshell ./test/
	@docker run \
		--rm \
		-it \
		--privileged \
		-v"$$(pwd)":/go/src/github.com/joemiller/go-jail \
		$(BUILD_CONTAINER) \
		/bin/bash

deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

build: create-build-container
	@docker run \
		--rm \
		-it \
		"-v$$(pwd):/go/src/github.com/joemiller/go-jail" \
		$(BUILD_CONTAINER) \
		/bin/sh -c 'make deps && CGO_ENABLED=1 go build -tags "static_build" -installsuffix cgo -ldflags "-w -extldflags -static"'
	@ls -l ./go-jail

test: create-build-container
	@docker run \
		--rm \
		-it \
		--privileged \
		"-v$$(pwd):/go/src/github.com/joemiller/go-jail" \
		$(BUILD_CONTAINER) \
		/bin/sh -c 'bats ./test'

release:
	go get -u github.com/tcnksm/ghr
	ghr -t "$$GITHUB_TOKEN" -u "$$CIRCLE_PROJECT_USERNAME" -r "$$CIRCLE_PROJECT_REPONAME" \
		--replace "$$CIRCLE_BUILD_NUM" \
		./go-jail

.PHONY: devshell deps test build
