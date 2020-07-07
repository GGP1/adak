package email

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/GGP1/palo/internal/cfg"
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
	defer v.Unlock()

	v.UserList[email] = token

	jsonMap, err := json.Marshal(v.UserList)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	err = ioutil.WriteFile(cfg.ValidatedJSONPath, jsonMap, 0644)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}

// Read validated emails list
func (v *ValidatedList) Read() map[string]string {
	v.RLock()
	defer v.RUnlock()

	jsonFile, err := ioutil.ReadFile(cfg.ValidatedJSONPath)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	err = json.Unmarshal(jsonFile, &v.UserList)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	return v.UserList
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
