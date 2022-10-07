package storagenomad

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	nomad "github.com/hashicorp/nomad/api"
)

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

	varPath := path.Join("acme", "example.com", "sites", "example.com", "example.com.crt")
	ctx := context.Background()

	err := ns.Store(ctx, varPath, []byte("crt data"))
	assert.NoError(t, err)
}

func TestNomadStorage_Exists(t *testing.T) {
	ns := setupNomadEnv(t)
	ctx := context.Background()

	key := path.Join("acme", "example.com", "sites", "example.com", "example.com.crt")

	err := ns.Store(ctx, key, []byte("crt data"))
	assert.NoError(t, err)

	exists := ns.Exists(ctx, key)
	assert.True(t, exists)
}

func TestNomadStorage_Load(t *testing.T) {
	ns := setupNomadEnv(t)
	ctx := context.Background()

	key := path.Join("acme", "example.com", "sites", "example.com", "example.com.crt")
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

	key := path.Join("acme", "example.com", "sites", "example.com", "example.com.crt")
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

	key := path.Join("acme", "example.com", "sites", "example.com", "example.com.crt")
	content := []byte("crt data")

	err := ns.Store(ctx, key, content)
	assert.NoError(t, err)

	info, err := ns.Stat(ctx, key)
	assert.NoError(t, err)

	assert.Equal(t, key, info.Key)
}

func TestNomadStorage_List(t *testing.T) {
	ns := setupNomadEnv(t)
	ctx := context.Background()

	err := ns.Store(ctx, path.Join("acme", "example.com", "sites", "example.com", "example.com.crt"), []byte("crt"))
	assert.NoError(t, err)
	err = ns.Store(ctx, path.Join("acme", "example.com", "sites", "first-level"), []byte("first"))
	assert.NoError(t, err)

	keys, err := ns.List(ctx, path.Join("acme", "example_com", "sites"), true)
	assert.NoError(t, err)
	assert.Len(t, keys, 2)
	assert.Contains(t, keys, path.Join("acme", "example_com", "sites", "example_com", "example_com_crt"))
}

func TestNomadStorage_ListNonRecursive(t *testing.T) {
	ns := setupNomadEnv(t)
	ctx := context.Background()

	err := ns.Store(ctx, path.Join("acme", "example.com", "sites", "example.com", "example.com.crt"), []byte("crt"))
	assert.NoError(t, err)
	err = ns.Store(ctx, path.Join("acme", "example.com", "sites", "first-level"), []byte("first"))
	assert.NoError(t, err)

	keys, err := ns.List(ctx, path.Join("acme", "example_com", "sites"), false)
	assert.NoError(t, err)
	assert.Len(t, keys, 1)
	assert.Contains(t, keys, path.Join("acme", "example_com", "sites", "first-level"))
}
