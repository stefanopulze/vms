VERSION ?= 1.0.0

build:
	podman build -t harbor.pulze.cloud/voltronic/vms-core:$(VERSION) --arch=arm64 --build-arg VERSION=$(VERSION) .
	podman push harbor.pulze.cloud/voltronic/vms-core:$(VERSION)