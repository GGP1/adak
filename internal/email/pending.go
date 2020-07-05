package email

import (
	"fmt"
	"sync"

	"github.com/GGP1/palo/pkg/model"
)

// PendingList contains users that already verified their email
type PendingList struct {
	UserList map[string]string
	mutex    sync.RWMutex
}

// NewPendingList returns a new list with user's emails and tokens
func NewPendingList() *PendingList {
	return &PendingList{
		UserList: make(map[string]string),
	}
}

// Add verified user to the list
func (p *PendingList) Add(user model.User, token string) {
	p.mutex.Lock()
	p.UserList[user.Email] = token
	p.mutex.Unlock()
}

// Print verified users list
func (p *PendingList) Print() {
	p.mutex.RLock()
	for k, v := range p.UserList {
		fmt.Println("Pending list:")
		fmt.Printf("[%s] %s\n", k, v)
	}
	p.mutex.RUnlock()
}
