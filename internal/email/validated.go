package email

import (
	"fmt"
	"sync"
)

// ValidatedList contains users that already verified their email
type ValidatedList struct {
	UserList map[string]string
	sync.RWMutex
}

// NewValidatedList returns a new list
func NewValidatedList() *ValidatedList {
	return &ValidatedList{
		UserList: make(map[string]string),
	}
}

// Add user to the list
func (v *ValidatedList) Add(email, token string) {
	v.Lock()
	v.UserList[email] = token
	v.Unlock()
}

// Print users list
func (v *ValidatedList) Print() {
	v.RLock()
	for k, val := range v.UserList {
		fmt.Println("Verified list:")
		fmt.Printf("[%s] %s\n", k, val)
	}
	v.RUnlock()
}

// Remove deletes a key from the map
func (v *ValidatedList) Remove(key string) {
	v.Lock()
	defer v.Unlock()

	for k := range v.UserList {
		if k == key {
			delete(v.UserList, k)
		}
	}
}

// Size returns the length of the user list
func (v *ValidatedList) Size() int {
	v.RLock()
	defer v.RUnlock()

	return len(v.UserList)
}
