module 01-talk-to-llm

go 1.23.1

require github.com/parakeet-nest/parakeet v0.2.2

require (
	go.etcd.io/bbolt v1.3.10 // indirect
	golang.org/x/sys v0.20.0 // indirect
)

replace github.com/parakeet-nest/parakeet => ../../..
