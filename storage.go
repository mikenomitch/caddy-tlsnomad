package storagenomad

import (
	"context"
	"fmt"
	"io/fs"
	"path"
	"strings"
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

func (ns *NomadStorage) readyKey(key string) string {
	return ns.escapeKey(ns.prefixKey(key))
}

func (ns *NomadStorage) prefixKey(key string) string {
	return path.Join(ns.Prefix, key)
}

func (ns *NomadStorage) escapeKey(key string) string {
	return strings.ReplaceAll(key, ".", "_")
}

// Store saves data value as a variable in Nomad
func (ns NomadStorage) Store(ctx context.Context, key string, value []byte) error {
	escapedAndPrefixedKey := ns.readyKey(key)

	items := &nomad.VariableItems{
		"Value":    string(value),
		"Modified": time.Now().Format(time.RFC3339),
	}

	sv := &nomad.Variable{
		Path:  escapedAndPrefixedKey,
		Items: *items,
	}

	opts := NomaWriteDefaults(ctx)

	if _, _, err := ns.NomadClient.Variables().Create(sv, opts); err != nil {
		msg := fmt.Sprintf("unable to store data for %s", escapedAndPrefixedKey)
		return wrapError(err, msg)
	}

	return nil
}

// Load retrieves the value for a key from Nomad KV
func (ns NomadStorage) Load(ctx context.Context, key string) ([]byte, error) {
	path := ns.readyKey(key)
	opts := NomadQueryDefaults(ctx)
	v, _, err := ns.NomadClient.Variables().Read(path, opts)
	if err != nil {
		msg := fmt.Sprintf("unable to read data for %s", ns.readyKey(key))
		return nil, wrapError(err, msg)
	}

	items := v.Items

	if val, ok := items["Value"]; ok {
		return []byte(val), nil
	}

	msg := fmt.Sprintf("unable to read data for %s", ns.readyKey(key))
	return nil, wrapError(fs.ErrNotExist, msg)
}

// Delete a key from Nomad KV
func (ns NomadStorage) Delete(ctx context.Context, key string) error {
	path := ns.readyKey(key)
	opts := NomaWriteDefaults(ctx)

	if _, err := ns.NomadClient.Variables().Delete(path, opts); err != nil {
		msg := fmt.Sprintf("unable to delete data for %s", ns.readyKey(key))
		return wrapError(err, msg)
	}

	return nil
}

// Exists checks if a key exists
func (ns NomadStorage) Exists(ctx context.Context, key string) bool {
	path := ns.readyKey(key)
	opts := NomadQueryDefaults(ctx)

	v, _, err := ns.NomadClient.Variables().Read(path, opts)
	if err != nil {
		return false
	}

	i := v.Items
	if _, ok := i["Value"]; ok {
		return true
	}

	return false
}

// List returns a list with all keys under a given prefix
func (ns NomadStorage) List(ctx context.Context, prefix string, recursive bool) ([]string, error) {
	var keysFound []string

	path := ns.prefixKey(prefix)
	opts := NomadQueryDefaults(ctx)
	keys, _, err := ns.NomadClient.Variables().PrefixList(path, opts)
	if err != nil {
		msg := fmt.Sprintf("unable to list data for %s", path)
		return nil, wrapError(err, msg)
	}

	for _, k := range keys {
		key := k.Path
		if strings.HasPrefix(key, path) {
			pf := path + "/"
			trimmedKey := strings.TrimPrefix(key, pf)
			isNested := strings.Contains(trimmedKey, "/")

			if recursive || !isNested {
				matchingPath := strings.TrimPrefix(key, ns.Prefix+"/")
				keysFound = append(keysFound, matchingPath)
			}
		}
	}

	if len(keys) == 0 {
		return keysFound, fs.ErrNotExist
	}

	return keysFound, nil
}

// Stat returns statistic data of a key
func (ns NomadStorage) Stat(ctx context.Context, key string) (certmagic.KeyInfo, error) {
	path := ns.readyKey(key)
	opts := NomadQueryDefaults(ctx)
	v, _, err := ns.NomadClient.Variables().Read(path, opts)
	if err != nil {
		msg := fmt.Sprintf("unable to read stats for %s", path)
		return certmagic.KeyInfo{}, wrapError(err, msg)
	}

	items := v.Items
	modified, mok := items["Modified"]
	val, vok := items["Value"]
	t, err := time.Parse(time.RFC3339, modified)

	if err != nil {
		msg := fmt.Sprintf("error parsing time when getting stats on %s", path)
		return certmagic.KeyInfo{}, wrapError(err, msg)
	}

	if mok && vok {
		return certmagic.KeyInfo{
			Key:        key,
			Modified:   t,
			Size:       int64(len(val)),
			IsTerminal: false,
		}, nil
	}

	msg := fmt.Sprintf("error reading value for stats %s", path)
	return certmagic.KeyInfo{}, fmt.Errorf(msg)
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
	opts := &nomad.QueryOptions{}
	return opts.WithContext(ctx)
}

func NomaWriteDefaults(ctx context.Context) *nomad.WriteOptions {
	opts := &nomad.WriteOptions{}
	return opts.WithContext(ctx)
}
