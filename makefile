build:
	podman build -t harbor.pulze.cloud/voltronic/vms-core:dev --arch=arm64 .
	podman push harbor.pulze.cloud/voltronic/vms-core:dev

build-os:
	container build -t harbor.pulze.cloud/voltronic/vms-core:dev --arch=arm64 .
	container image push harbor.pulze.cloud/voltronic/vms-core:dev