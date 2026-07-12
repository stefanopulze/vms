VERSION ?= 1.2.0

build:
	podman build -t harbor.pulze.cloud/voltronic/vms-core:$(VERSION) --arch=arm64 --build-arg VERSION=$(VERSION) .
	podman push harbor.pulze.cloud/voltronic/vms-core:$(VERSION)

build-os:
	container build -t harbor.pulze.cloud/voltronic/vms-core:$(VERSION) --arch=arm64 --build-arg VERSION=$(VERSION) .
	container image push harbor.pulze.cloud/voltronic/vms-core:$(VERSION)

update:
	ssh pi@192.168.1.42 "cd /opt/vms && docker compose pull && docker compose up -d"