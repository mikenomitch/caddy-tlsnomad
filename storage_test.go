package storagenomad

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	nomad "github.com/hashicorp/nomad/api"
)

// var nomadClient *nomad.Client

const TestPrefix = "nomadtlstest"

// these tests needs a running Nomad server
func setupNomadEnv(t *testing.T) *NomadStorage {

	os.Setenv(EnvNamePrefix, TestPrefix)
	// os.Setenv("NOMAD_TOKEN", "2f9e03f8-714b-5e4d-65ea-c983d6b172c4")

	ns := New()
	ns.createNomadClient()

	_, err := ns.NomadClient.Variables().Delete(TestPrefix, &nomad.WriteOptions{})
	assert.NoError(t, err)
	return ns
}

func TestNomadStorage_Store(t *testing.T) {
	ns := setupNomadEnv(t)

	varPath := path.Join("acme", "examplecom", "sites", "examplecom", "examplecomcrt")
	ctx := context.Background()

	err := ns.Store(ctx, varPath, []byte("crt data"))
	assert.NoError(t, err)
}

func TestNomadStorage_Exists(t *testing.T) {
	ns := setupNomadEnv(t)
	ctx := context.Background()

	key := path.Join("acme", "examplecom", "sites", "examplecom", "examplecomcrt")

	err := ns.Store(ctx, key, []byte("crt data"))
	assert.NoError(t, err)

	exists := ns.Exists(ctx, key)
	assert.True(t, exists)
}

func TestNomadStorage_Load(t *testing.T) {
	ns := setupNomadEnv(t)
	ctx := context.Background()

	key := path.Join("acme", "examplecom", "sites", "examplecom", "examplecomcrt")
	content := []byte("crt data")

	err := ns.Store(ctx, key, content)
	assert.NoError(t, err)

	contentLoded, err := ns.Load(ctx, key)
	assert.NoError(t, err)

	assert.Equal(t, content, contentLoded)
}

func TestNomadStorage_Delete(t *testing.T) {
	ns := setupNomadEnv(t)
	ctx := context.Background()

	key := path.Join("acme", "examplecom", "sites", "examplecom", "examplecomcrt")
	content := []byte("crt data")

	err := ns.Store(ctx, key, content)
	assert.NoError(t, err)

	err = ns.Delete(ctx, key)
	assert.NoError(t, err)

	exists := ns.Exists(ctx, key)
	assert.False(t, exists)

	contentLoaded, err := ns.Load(ctx, key)
	assert.Nil(t, contentLoaded)

	assert.Equal(t, "variable not found", err.Error())
}

func TestNomadStorage_Stat(t *testing.T) {
	ns := setupNomadEnv(t)
	ctx := context.Background()

	key := path.Join("acme", "examplecom", "sites", "examplecom", "examplecomcrt")
	content := []byte("crt data")

	err := ns.Store(ctx, key, content)
	assert.NoError(t, err)

	info, err := ns.Stat(ctx, key)
	assert.NoError(t, err)

	assert.Equal(t, key, info.Key)
}

// func TestNomadStorage_List(t *testing.T) {
// 	ns := setupNomadEnv(t)

// 	err := ns.Store(path.Join("acme", "examplecom", "sites", "examplecom", "examplecomcrt"), []byte("crt"))
// 	assert.NoError(t, err)
// 	err = ns.Store(path.Join("acme", "examplecom", "sites", "examplecom", "examplecom.key"), []byte("key"))
// 	assert.NoError(t, err)
// 	err = ns.Store(path.Join("acme", "examplecom", "sites", "examplecom", "examplecom.json"), []byte("meta"))
// 	assert.NoError(t, err)

// 	keys, err := ns.List(path.Join("acme", "examplecom", "sites", "examplecom"), true)
// 	assert.NoError(t, err)
// 	assert.Len(t, keys, 3)
// 	assert.Contains(t, keys, path.Join("acme", "examplecom", "sites", "examplecom", "examplecomcrt"))
// }

// func TestNomadStorage_ListNonRecursive(t *testing.T) {
// 	ns := setupNomadEnv(t)

// 	err := ns.Store(path.Join("acme", "examplecom", "sites", "examplecom", "examplecomcrt"), []byte("crt"))
// 	assert.NoError(t, err)
// 	err = ns.Store(path.Join("acme", "examplecom", "sites", "examplecom", "examplecom.key"), []byte("key"))
// 	assert.NoError(t, err)
// 	err = ns.Store(path.Join("acme", "examplecom", "sites", "examplecom", "examplecom.json"), []byte("meta"))
// 	assert.NoError(t, err)

// 	keys, err := ns.List(path.Join("acme", "examplecom", "sites"), false)
// 	assert.NoError(t, err)

// 	assert.Len(t, keys, 1)
// 	assert.Contains(t, keys, path.Join("acme", "examplecom", "sites", "examplecom"))
// }
