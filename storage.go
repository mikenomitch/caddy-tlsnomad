package storagenomad

import (
	"context"
	"fmt"
	"io/fs"
	"path"
	"time"

	"github.com/caddyserver/certmagic"
	nomad "github.com/hashicorp/nomad/api"
	"go.uber.org/zap"
)

// NomadStorage allows to store certificates and other TLS resources
// in a shared cluster environment using Nomad Variables.
type NomadStorage struct {
	certmagic.Storage
	NomadClient *nomad.Client
	logger      *zap.SugaredLogger

	Address     string `json:"address"`
	Token       string `json:"token"`
	Timeout     int    `json:"timeout"`
	Prefix      string `json:"prefix"`
	ValuePrefix string `json:"value_prefix"`
	TlsEnabled  bool   `json:"tls_enabled"`
	TlsInsecure bool   `json:"tls_insecure"`
}

// New connects to Nomad and returns a NomadStorage
func New() *NomadStorage {
	// create NomadStorage and pre-set values
	s := NomadStorage{
		ValuePrefix: DefaultValuePrefix,
		Prefix:      DefaultPrefix,
		Timeout:     DefaultTimeout,
	}

	return &s
}

func (ns *NomadStorage) prefixKey(key string) string {
	return path.Join(ns.Prefix, key)
}

// Store saves data value as a variable in Nomad
func (ns NomadStorage) Store(ctx context.Context, key string, value []byte) error {
	items := &nomad.VariableItems{
		"Value":    string(value),
		"Modified": time.Now().String(),
	}

	sv := &nomad.Variable{
		Path:  ns.prefixKey(key),
		Items: *items,
	}

	opts := NomaWriteDefaults(ctx)

	if _, _, err := ns.NomadClient.Variables().Create(sv, opts); err != nil {
		msg := fmt.Sprintf("unable to store data for %s", ns.prefixKey(key))
		return wrapError(err, msg)
	}

	return nil
}

// Load retrieves the value for a key from Nomad KV
func (ns NomadStorage) Load(ctx context.Context, key string) ([]byte, error) {
	path := ns.prefixKey(key)
	opts := NomadQueryDefaults(ctx)
	items, _, err := ns.NomadClient.Variables().GetItems(path, opts)
	if err != nil {
		return nil, err
	}

	i := *items
	val := i["Value"]

	if val == "" {
		return nil, fs.ErrNotExist
	}

	return []byte(val), nil
}

// Delete a key from Nomad KV
func (ns NomadStorage) Delete(ctx context.Context, key string) error {
	path := ns.prefixKey(key)
	opts := NomaWriteDefaults(ctx)

	if _, err := ns.NomadClient.Variables().Delete(path, opts); err != nil {
		msg := fmt.Sprintf("unable to delete data for %s", ns.prefixKey(key))
		return wrapError(err, msg)
	}

	return nil
}

// Exists checks if a key exists
func (ns NomadStorage) Exists(ctx context.Context, key string) bool {
	path := ns.prefixKey(key)
	opts := NomadQueryDefaults(ctx)
	items, _, err := ns.NomadClient.Variables().GetItems(path, opts)
	if err != nil {
		// TODO: make sure this interface is okay
		return false
	}

	i := *items

	fmt.Println("i")
	fmt.Println(i)

	// val := i["Value"]
	if _, ok := i["Value"]; ok {
		return true
	}

	return false
}

// List returns a list with all keys under a given prefix
func (ns NomadStorage) List(ctx context.Context, prefix string, recursive bool) ([]string, error) {
	// TODO: LIST KEYS UNDER THE PREFIX

	// var keysFound []string

	// // get a list of all keys at prefix
	// keys, _, err := ns.NomadClient.KV().Keys(ns.prefixKey(prefix), "", NomadQueryDefaults(ctx))
	// if err != nil {
	// 	return keysFound, err
	// }

	// if len(keys) == 0 {
	// 	return keysFound, fs.ErrNotExist
	// }

	// // remove default prefix from keys
	// for _, key := range keys {
	// 	if strings.HasPrefix(key, ns.prefixKey(prefix)) {
	// 		key = strings.TrimPrefix(key, ns.Prefix+"/")
	// 		keysFound = append(keysFound, key)
	// 	}
	// }

	// // if recursive wanted, just return all keys
	// if recursive {
	// 	return keysFound, nil
	// }

	// // for non-recursive split path and look for unique keys just under given prefix
	// keysMap := make(map[string]bool)
	// for _, key := range keysFound {
	// 	dir := strings.Split(strings.TrimPrefix(key, prefix+"/"), "/")
	// 	keysMap[dir[0]] = true
	// }
	// keysFound = make([]string, 0)
	// for key := range keysMap {
	// 	keysFound = append(keysFound, path.Join(prefix, key))
	// }

	// return keysFound, nil

	keysFound := make([]string, 0)
	return keysFound, nil
}

// Stat returns statistic data of a key
func (ns NomadStorage) Stat(ctx context.Context, key string) (certmagic.KeyInfo, error) {
	// kv, _, err := ns.NomadClient.KV().Get(ns.prefixKey(key), NomadQueryDefaults(ctx))
	// if err != nil {
	// 	return certmagic.KeyInfo{}, fmt.Errorf("unable to obtain data for %s: %w", ns.prefixKey(key), fs.ErrNotExist)
	// }
	// if kv == nil {
	// 	return certmagic.KeyInfo{}, fs.ErrNotExist
	// }

	// return certmagic.KeyInfo{
	// 	Key:        key,
	// 	Modified:   kv.Mofified,
	// 	Size:       int64(len(kv.Value)),
	// 	IsTerminal: false,
	// }, nil

	return certmagic.KeyInfo{
		Key:        "wat",
		Size:       int64(len("wat")),
		IsTerminal: false,
	}, nil
}

func (ns *NomadStorage) createNomadClient() error {
	// get the default config
	nomadCfg := nomad.DefaultConfig()
	if ns.Address != "" {
		nomadCfg.Address = ns.Address
	}
	if ns.Token != "" {
		nomadCfg.SecretID = ns.Token
	}
	// if ns.TlsEnabled {
	// 	nomadCfg.Scheme = "https"
	// }

	nomadCfg.TLSConfig.Insecure = ns.TlsInsecure

	// set a dial context to prevent default keepalive
	// nomadCfg.Transport.DialContext = (&net.Dialer{
	// 	Timeout:   time.Duration(ns.Timeout) * time.Second,
	// 	KeepAlive: time.Duration(ns.Timeout) * time.Second,
	// }).DialContext

	// create the Nomad API client
	nomadClient, err := nomad.NewClient(nomadCfg)
	if err != nil {
		return wrapError(err, "unable to create Nomad client")
	}

	if _, err := nomadClient.Agent().NodeName(); err != nil {
		return wrapError(err, "unable to ping Nomad")
	}

	ns.NomadClient = nomadClient
	return nil
}

func NomadQueryDefaults(ctx context.Context) *nomad.QueryOptions {
	// TODO: Set some of these
	opts := &nomad.QueryOptions{}
	return opts.WithContext(ctx)
}

func NomaWriteDefaults(ctx context.Context) *nomad.WriteOptions {
	// TODO: Set some of these
	opts := &nomad.WriteOptions{}
	return opts.WithContext(ctx)
}
