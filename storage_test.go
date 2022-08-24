//go:build consul
// +build consul

package storagenomad

import (
	"github.com/caddyserver/certmagic"
	consul "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
	"time"
)

var consulClient *consul.Client

const TestPrefix = "consultlstest"

// these tests needs a running Nomad server
func setupNomadEnv(t *testing.T) *NomadStorage {

	os.Setenv(EnvNamePrefix, TestPrefix)
	os.Setenv(consul.HTTPTokenEnvName, "2f9e03f8-714b-5e4d-65ea-c983d6b172c4")

	cs, err := New()
	assert.NoError(t, err)

	_, err = cs.NomadClient.KV().DeleteTree(TestPrefix, nil)
	assert.NoError(t, err)
	return cs
}

func TestNomadStorage_Store(t *testing.T) {
	cs := setupNomadEnv(t)

	err := cs.Store(path.Join("acme", "example.com", "sites", "example.com", "example.com.crt"), []byte("crt data"))
	assert.NoError(t, err)
}

func TestNomadStorage_Exists(t *testing.T) {
	cs := setupNomadEnv(t)

	key := path.Join("acme", "example.com", "sites", "example.com", "example.com.crt")

	err := cs.Store(key, []byte("crt data"))
	assert.NoError(t, err)

	exists := cs.Exists(key)
	assert.True(t, exists)
}

func TestNomadStorage_Load(t *testing.T) {
	cs := setupNomadEnv(t)

	key := path.Join("acme", "example.com", "sites", "example.com", "example.com.crt")
	content := []byte("crt data")

	err := cs.Store(key, content)
	assert.NoError(t, err)

	contentLoded, err := cs.Load(key)
	assert.NoError(t, err)

	assert.Equal(t, content, contentLoded)
}

func TestNomadStorage_Delete(t *testing.T) {
	cs := setupNomadEnv(t)

	key := path.Join("acme", "example.com", "sites", "example.com", "example.com.crt")
	content := []byte("crt data")

	err := cs.Store(key, content)
	assert.NoError(t, err)

	err = cs.Delete(key)
	assert.NoError(t, err)

	exists := cs.Exists(key)
	assert.False(t, exists)

	contentLoaded, err := cs.Load(key)
	assert.Nil(t, contentLoaded)

	_, ok := err.(certmagic.ErrNotExist)
	assert.True(t, ok)
}

func TestNomadStorage_Stat(t *testing.T) {
	cs := setupNomadEnv(t)

	key := path.Join("acme", "example.com", "sites", "example.com", "example.com.crt")
	content := []byte("crt data")

	err := cs.Store(key, content)
	assert.NoError(t, err)

	info, err := cs.Stat(key)
	assert.NoError(t, err)

	assert.Equal(t, key, info.Key)
}

func TestNomadStorage_List(t *testing.T) {
	cs := setupNomadEnv(t)

	err := cs.Store(path.Join("acme", "example.com", "sites", "example.com", "example.com.crt"), []byte("crt"))
	assert.NoError(t, err)
	err = cs.Store(path.Join("acme", "example.com", "sites", "example.com", "example.com.key"), []byte("key"))
	assert.NoError(t, err)
	err = cs.Store(path.Join("acme", "example.com", "sites", "example.com", "example.com.json"), []byte("meta"))
	assert.NoError(t, err)

	keys, err := cs.List(path.Join("acme", "example.com", "sites", "example.com"), true)
	assert.NoError(t, err)
	assert.Len(t, keys, 3)
	assert.Contains(t, keys, path.Join("acme", "example.com", "sites", "example.com", "example.com.crt"))
}

func TestNomadStorage_ListNonRecursive(t *testing.T) {
	cs := setupNomadEnv(t)

	err := cs.Store(path.Join("acme", "example.com", "sites", "example.com", "example.com.crt"), []byte("crt"))
	assert.NoError(t, err)
	err = cs.Store(path.Join("acme", "example.com", "sites", "example.com", "example.com.key"), []byte("key"))
	assert.NoError(t, err)
	err = cs.Store(path.Join("acme", "example.com", "sites", "example.com", "example.com.json"), []byte("meta"))
	assert.NoError(t, err)

	keys, err := cs.List(path.Join("acme", "example.com", "sites"), false)
	assert.NoError(t, err)

	assert.Len(t, keys, 1)
	assert.Contains(t, keys, path.Join("acme", "example.com", "sites", "example.com"))
}

func TestNomadStorage_LockUnlock(t *testing.T) {
	cs := setupNomadEnv(t)
	lockKey := path.Join("acme", "example.com", "sites", "example.com", "lock")

	err := cs.Lock(lockKey)
	assert.NoError(t, err)

	err = cs.Unlock(lockKey)
	assert.NoError(t, err)
}

func TestNomadStorage_TwoLocks(t *testing.T) {
	cs := setupNomadEnv(t)
	cs2 := setupNomadEnv(t)
	lockKey := path.Join("acme", "example.com", "sites", "example.com", "lock")

	err := cs.Lock(lockKey)
	assert.NoError(t, err)

	go time.AfterFunc(5*time.Second, func() {
		err = cs.Unlock(lockKey)
		assert.NoError(t, err)
	})

	err = cs2.Lock(lockKey)
	assert.NoError(t, err)

	err = cs2.Unlock(lockKey)
	assert.NoError(t, err)
}
