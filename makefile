build:
	podman build -t harbor.pulze.cloud/voltronic/vms-core:dev --arch=arm64 .
	podman push harbor.pulze.cloud/voltronic/vms-core:dev

build-os:
	container build -t harbor.pulze.cloud/voltronic/vms-core:dev --arch=arm64 .
	container image push harbor.pulze.cloud/voltronic/vms-core:dev

update:
	ssh pi@192.168.1.42 "cd /opt/vms && docker compose pull && docker compose up -d"