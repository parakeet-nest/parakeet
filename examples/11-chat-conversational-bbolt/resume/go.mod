module 11-chat-bbolt-resume

go 1.22.1

require github.com/parakeet-nest/parakeet v0.0.4

require (
	go.etcd.io/bbolt v1.3.10 // indirect
	golang.org/x/sys v0.4.0 // indirect
)

replace github.com/parakeet-nest/parakeet => ../../..
