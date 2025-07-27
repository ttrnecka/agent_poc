module github.com/ttrnecka/agent_poc/collector

go 1.23.0

require (
	github.com/fsnotify/fsnotify v1.9.0
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/gorilla/websocket v1.5.3
	github.com/ttrnecka/agent_poc/webapi v0.0.0-00010101000000-000000000000
)

replace github.com/ttrnecka/agent_poc/webapi => ../webapi

require (
	github.com/google/uuid v1.6.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
)
