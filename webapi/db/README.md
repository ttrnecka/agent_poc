to test

systemctl --user start podman.socket

podman info --format '{{.Host.RemoteSocket.Path}}'

export DOCKER_HOST=unix://<your_podman_socket_location>

go test ./...