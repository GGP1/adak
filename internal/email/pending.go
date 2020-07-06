package email

import (
	"fmt"
	"sync"

	"github.com/GGP1/palo/pkg/model"
)

// PendingList contains users that need to validate their email
type PendingList struct {
	UserList map[string]string
	sync.RWMutex
}

// NewPendingList returns a new pending list
func NewPendingList() *PendingList {
	return &PendingList{
		UserList: make(map[string]string),
	}
}

// Add pending user to the list
func (p *PendingList) Add(user model.User, token string) {
	p.Lock()
	p.UserList[user.Email] = token
	p.Unlock()
}

// Print pending users list
func (p *PendingList) Print() {
	p.RLock()
	for k, v := range p.UserList {
		fmt.Println("Pending list:")
		fmt.Printf("[%s] %s\n", k, v)
	}
	p.RUnlock()
}

// Remove deletes a key from the map
func (p *PendingList) Remove(key string) {
	p.Lock()
	defer p.Unlock()

	for k := range p.UserList {
		if k == key {
			delete(p.UserList, k)
		}
	}
}

// Size returns the length of the user list
func (p *PendingList) Size() int {
	p.RLock()
	defer p.RUnlock()

	return len(p.UserList)
}
