package storage

import (
	"github.com/sebdehne/accountingserver/domain"
	"io/ioutil"
	"encoding/json"
	"os"
	"os/exec"
	"github.com/iris-contrib/errors"
	"strconv"
	"sync"
)

func New(storageDir, storageFilename string) Storage {
	return Storage{storageDir:storageDir, storageFilename:storageFilename}
}

type Storage struct {
	root            *domain.Root
	storageDir      string
	storageFilename string
	lock            sync.Mutex
}

func (s *Storage) Save(r domain.Root) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	defer os.Chdir(pwd)

	existing, _ := s.readFromFile()
	if existing.Version + 1 != r.Version {
		err = errors.New("Expected version " + strconv.Itoa(existing.Version + 1))
		return
	}

	os.Chdir(s.storageDir)

	json, err := json.Marshal(r)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(s.storageFilename, json, 0644)
	if err != nil {
		return
	}
	s.root = &r

	_, err = exec.Command("git", "add", s.storageFilename).Output()
	if err != nil {
		return
	}
	_, err = exec.Command("git", "commit", "-m", "version:" + strconv.Itoa(r.Version)).Output()

	return
}

func (s *Storage) Get() (r domain.Root, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.root == nil {
		r, err = s.readFromFile()
		if err == nil {
			s.root = &r
		}
	} else {
		r = *s.root
	}

	return
}

func (s *Storage) readFromFile() (r domain.Root, err error) {
	data, err := ioutil.ReadFile(s.storageDir + "/" + s.storageFilename)
	if err != nil {
		return
	}
	r = domain.Root{Version:1}
	err = json.Unmarshal(data, &r)
	return
}

func (s *Storage) InitStorage() (err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	defer os.Chdir(pwd)

	fileInfo, err := os.Stat(s.storageDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(s.storageDir, 0755)
		if err != nil {
			return
		}
		fileInfo, err = os.Stat(s.storageDir)
	}

	if err != nil {
		return
	}

	if !fileInfo.IsDir() {
		err = errors.New(s.storageDir + " is not a directory")
		return
	}

	// now we know it exists and it is a directory
	// next, ensure contains only json files and one .git directory
	files, err := ioutil.ReadDir(fileInfo.Name())
	if err != nil {
		return
	}

	hasGitRepo := false
	hasJsonFile := false
	for _, fi := range files {
		if fi.Name() == ".git" && fi.IsDir() {
			hasGitRepo = true
		} else if fi.Name() == s.storageFilename {
			hasJsonFile = true
		}
	}

	if len(files) == 1 && !hasGitRepo && !hasJsonFile {
		err = errors.New("Storage directory contains unknown files")
	} else if len(files) == 2 && !(hasGitRepo && hasJsonFile) {
		err = errors.New("Storage directory contains unknown files")
	} else if len(files) > 2 {
		err = errors.New("Storage directory contains unknown files")
	} else {
		os.Chdir(s.storageDir)
		if !hasGitRepo {
			_, err = exec.Command("git", "init").Output()
			if err != nil {
				return
			}
		}

		if !hasJsonFile {
			err = s.Save(domain.New())
		}
	}

	return
}
