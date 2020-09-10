package dao

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"qq/model"
	"sync"
)

var CommonDao Dao = new(jsonFileDao)

type Dao interface {
	Write(path string, obj interface{}) error
	Read(path string, obj interface{}) (interface{}, error)
}

type gobFileDao struct {
	sync.RWMutex
}

func (f *gobFileDao) Write(path string, obj interface{}) error {
	f.Lock()
	defer f.Unlock()
	file, err := os.Create(path)
	if err == nil {
		gob.NewEncoder(file).Encode(obj)
	}
	file.Close()
	return err
}
func (f *gobFileDao) Read(path string, obj interface{}) (interface{}, error) {
	f.RLock()
	defer f.RUnlock()
	file, err := os.Open(path)
	defer file.Close()
	if err == nil {
		gob.NewDecoder(file).Decode(obj)
	}
	return obj, err
}

type jsonFileDao struct {
	sync.RWMutex
}

func (f *jsonFileDao) Write(path string, obj interface{}) error {
	f.Lock()
	defer f.Unlock()
	file, err := os.Create(path)
	defer file.Close()
	if err == nil {
		err = json.NewEncoder(file).Encode(obj)
	}
	return err
}
func (f *jsonFileDao) Read(path string, obj interface{}) (interface{}, error) {
	f.RLock()
	defer f.RUnlock()
	file, err := os.Open(path)
	defer file.Close()
	if err == nil {
		json.NewDecoder(file).Decode(obj)
	}
	return obj, err
}

func Test() {
	user := model.User{}
	CommonDao.Write("./config/users.json", user)
	userr := new(model.User)
	CommonDao.Read("./config/users.json", userr)
	fmt.Println(*userr)

}
