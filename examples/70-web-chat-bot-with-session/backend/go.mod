module 70-web-chat-bot-with-session

go 1.24.0

require (
	github.com/google/uuid v1.6.0
	github.com/parakeet-nest/parakeet v0.2.8
)

require (
	go.etcd.io/bbolt v1.3.11 // indirect
	golang.org/x/sys v0.28.0 // indirect
)

replace github.com/parakeet-nest/parakeet => ../../..
