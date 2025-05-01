module 11-chat-bbolt-begin

go 1.24.0

require github.com/parakeet-nest/parakeet v0.2.7

require (
	github.com/google/uuid v1.6.0 // indirect
	go.etcd.io/bbolt v1.4.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
)

replace github.com/parakeet-nest/parakeet => ../../..
