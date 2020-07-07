package email

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/GGP1/palo/internal/cfg"
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
	defer p.Unlock()

	p.UserList[user.Email] = token

	jsonMap, err := json.Marshal(p.UserList)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	err = ioutil.WriteFile(cfg.PendingJSONPath, jsonMap, 0644)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}

// Read pending emails list
func (p *PendingList) Read() map[string]string {
	p.RLock()
	defer p.RUnlock()

	jsonFile, err := ioutil.ReadFile(cfg.PendingJSONPath)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	err = json.Unmarshal(jsonFile, &p.UserList)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	return p.UserList
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
