# set the IMAGE_NAME variable
# It should consist of repo name, image name and version
# use hyphens as separator
IMAGE_NAME = gcr.io/ticker-beats/ticker-signals
# set the IMAGE_VERSION variable
# It should consist of the version of the image
# use hyphens as separator
IMAGE_VERSION = 17
# set the IMAGE_TAG variable
# It should consist of the IMAGE_NAME and IMAGE_VERSION
# use colon as separator
IMAGE_TAG = $(IMAGE_NAME):$(IMAGE_VERSION)
# set the IMAGE_LATEST variable
# It should consist of the IMAGE_NAME and latest
# use colon as separator
IMAGE_LATEST = $(IMAGE_NAME):latest


.PHONY: generate
generate: ## Traverses project recursively, running go generate commands
	$(GO) generate ./...

.PHONY: build
build: ## builds the project using docker build
	docker build -t $(IMAGE_TAG) .
	docker tag $(IMAGE_TAG) $(IMAGE_LATEST)
	docker push $(IMAGE_TAG)
	docker push $(IMAGE_LATEST)