module 11-chat-bbolt-begin

go 1.24.0

require github.com/parakeet-nest/parakeet v0.2.7

require (
	go.etcd.io/bbolt v1.3.10 // indirect
	golang.org/x/sys v0.4.0 // indirect
)

replace github.com/parakeet-nest/parakeet => ../../..
