package email

import (
	"fmt"
	"sync"
)

// VerifiedList contains users that already verified their email
type VerifiedList struct {
	UserList map[string]string
	mutex    sync.RWMutex
}

// NewVerifiedList returns a new list with user's emails and tokens
func NewVerifiedList() *VerifiedList {
	return &VerifiedList{
		UserList: make(map[string]string),
	}
}

// Add verified user to the list
func (v *VerifiedList) Add(email, token string) {
	v.mutex.Lock()
	v.UserList[email] = token
	v.mutex.Unlock()
}

// Print verified users list
func (v *VerifiedList) Print() {
	v.mutex.RLock()
	for k, val := range v.UserList {
		fmt.Println("Verified list:")
		fmt.Printf("[%s] %s\n", k, val)
	}
	v.mutex.RUnlock()
}
