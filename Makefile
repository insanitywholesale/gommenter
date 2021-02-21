.PHONY: buildwithvars

# the real command is the following
# go build -v -ldflags "-X main.commitHash=$(git rev-parse --short HEAD) -X main.commitDate=$(git log -1 --format=%ci | awk '{ print $1 }')"
buildwithvars:
	rm -rf ./gommenter; go build -v -ldflags "-X main.commitHash=$$(git rev-parse --short HEAD) -X main.commitDate=$$(git log -1 --format=%ci | awk '{ print $$1 }')"
