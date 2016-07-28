package storage

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPersistenceService(t *testing.T) {
	assert := assert.New(t)

	s := New("data", "test.json")
	err := s.InitStorage()
	assert.Nil(err)

	r, _ := s.Get()
	r.Version++
	err = s.SaveAndCommit(r)
	assert.Nil(err)

	r2, err := s.Get()
	assert.Nil(err)

	assert.Equal(r2.Version, r.Version)
}