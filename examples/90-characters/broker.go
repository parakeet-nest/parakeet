package main
// This package is experimental
import (
	"sort"
	"sync"
	"time"
)

type Mail struct {
	Content     string
	SenderID    string
	RecipientID string
	Date        time.Time
	Read        bool
}

// Define a struct that will share a list of messages
type MailBroker struct {
	list []Mail
	mu   sync.Mutex
}

// Add a message to the shared list
func (s *MailBroker) Add(mail Mail) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.list = append(s.list, mail)
	// Sort the list from oldest to newest
	sort.Slice(s.list, func(i, j int) bool {
		return s.list[i].Date.Before(s.list[j].Date)
	})
}

// Get a copy of the shared list
func (s *MailBroker) Get() []Mail {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Return a copy of the list to prevent external modifications
	copiedList := make([]Mail, len(s.list))
	copy(copiedList, s.list)
	return copiedList
}

// Mark a message as read
func (s *MailBroker) MarkAsRead(index int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index >= 0 && index < len(s.list) {
		s.list[index].Read = true
	}
}
