package storagenomad

import (
	"context"
	"os"
	"path"
	"testing"

	nomad "github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/assert"
)

var nomadClient *nomad.Client

const TestPrefix = "nomadtlstest"

// these tests needs a running Nomad server
func setupNomadEnv(t *testing.T) *NomadStorage {

	os.Setenv(EnvNamePrefix, TestPrefix)
	// os.Setenv("NOMAD_TOKEN", "2f9e03f8-714b-5e4d-65ea-c983d6b172c4")

	ns := New()

	// _, err := ns.NomadClient.SecureVariables().Delete(TestPrefix, &nomad.WriteOptions{})
	// assert.NoError(t, err)
	return ns
}

func TestNomadStorage_Store(t *testing.T) {
	ns := setupNomadEnv(t)

	varPath := path.Join("acme", "example.com", "sites", "example.com", "example.com.crt")
	ctx := context.Background()

	err := ns.Store(ctx, varPath, []byte("crt data"))
	assert.NoError(t, err)
}

// func TestNomadStorage_Exists(t *testing.T) {
// 	ns := setupNomadEnv(t)

// 	key := path.Join("acme", "example.com", "sites", "example.com", "example.com.crt")

// 	err := ns.Store(key, []byte("crt data"))
// 	assert.NoError(t, err)

// 	exists := ns.Exists(key)
// 	assert.True(t, exists)
// }

// func TestNomadStorage_Load(t *testing.T) {
// 	ns := setupNomadEnv(t)

// 	key := path.Join("acme", "example.com", "sites", "example.com", "example.com.crt")
// 	content := []byte("crt data")

// 	err := ns.Store(key, content)
// 	assert.NoError(t, err)

// 	contentLoded, err := ns.Load(key)
// 	assert.NoError(t, err)

// 	assert.Equal(t, content, contentLoded)
// }

// func TestNomadStorage_Delete(t *testing.T) {
// 	ns := setupNomadEnv(t)

// 	key := path.Join("acme", "example.com", "sites", "example.com", "example.com.crt")
// 	content := []byte("crt data")

// 	err := ns.Store(key, content)
// 	assert.NoError(t, err)

// 	err = ns.Delete(key)
// 	assert.NoError(t, err)

// 	exists := ns.Exists(key)
// 	assert.False(t, exists)

// 	contentLoaded, err := ns.Load(key)
// 	assert.Nil(t, contentLoaded)

// 	_, ok := err.(certmagic.ErrNotExist)
// 	assert.True(t, ok)
// }

// func TestNomadStorage_Stat(t *testing.T) {
// 	ns := setupNomadEnv(t)

// 	key := path.Join("acme", "example.com", "sites", "example.com", "example.com.crt")
// 	content := []byte("crt data")

// 	err := ns.Store(key, content)
// 	assert.NoError(t, err)

// 	info, err := ns.Stat(key)
// 	assert.NoError(t, err)

// 	assert.Equal(t, key, info.Key)
// }

// func TestNomadStorage_List(t *testing.T) {
// 	ns := setupNomadEnv(t)

// 	err := ns.Store(path.Join("acme", "example.com", "sites", "example.com", "example.com.crt"), []byte("crt"))
// 	assert.NoError(t, err)
// 	err = ns.Store(path.Join("acme", "example.com", "sites", "example.com", "example.com.key"), []byte("key"))
// 	assert.NoError(t, err)
// 	err = ns.Store(path.Join("acme", "example.com", "sites", "example.com", "example.com.json"), []byte("meta"))
// 	assert.NoError(t, err)

// 	keys, err := ns.List(path.Join("acme", "example.com", "sites", "example.com"), true)
// 	assert.NoError(t, err)
// 	assert.Len(t, keys, 3)
// 	assert.Contains(t, keys, path.Join("acme", "example.com", "sites", "example.com", "example.com.crt"))
// }

// func TestNomadStorage_ListNonRecursive(t *testing.T) {
// 	ns := setupNomadEnv(t)

// 	err := ns.Store(path.Join("acme", "example.com", "sites", "example.com", "example.com.crt"), []byte("crt"))
// 	assert.NoError(t, err)
// 	err = ns.Store(path.Join("acme", "example.com", "sites", "example.com", "example.com.key"), []byte("key"))
// 	assert.NoError(t, err)
// 	err = ns.Store(path.Join("acme", "example.com", "sites", "example.com", "example.com.json"), []byte("meta"))
// 	assert.NoError(t, err)

// 	keys, err := ns.List(path.Join("acme", "example.com", "sites"), false)
// 	assert.NoError(t, err)

// 	assert.Len(t, keys, 1)
// 	assert.Contains(t, keys, path.Join("acme", "example.com", "sites", "example.com"))
// }
