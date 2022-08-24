package storagenomad

import (
	"context"
	"path"

	"github.com/caddyserver/certmagic"
	nomad "github.com/hashicorp/nomad/api"
	"github.com/pteich/errors"
	"go.uber.org/zap"
)

// NomadStorage allows to store certificates and other TLS resources
// in a shared cluster environment using Nomad Secure Variables.
// It uses distributed locks to ensure consistency.
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

func (cs *NomadStorage) prefixKey(key string) string {
	return path.Join(cs.Prefix, key)
}

// Store saves data value as a secure variable in Nomad
func (cs NomadStorage) Store(ctx context.Context, key string, value []byte) error {
	// kv := &nomad.KVPair{Key: cs.prefixKey(key)}

	// // prepare the stored data
	// nomadData := &StorageData{
	// 	Value:    value,
	// 	Modified: time.Now(),
	// }

	// kv.Value = nomadData

	// opts := nomad.WriteOptions{}
	// if _, err = cs.NomadClient.KV().Put(kv, opts.WithContext(ctx)); err != nil {
	// 	return errors.Wrapf(err, "unable to store data for %s", cs.prefixKey(key))
	// }

	// TODO: Create a Nomad SV with the Data

	return nil
}

// Load retrieves the value for a key from Nomad KV
func (cs NomadStorage) Load(ctx context.Context, key string) ([]byte, error) {
	cs.logger.Debugf("loading data from Nomad for %s", key)

	// TODO: Load a Nomad SV

	// kv, _, err := cs.NomadClient.KV().Get(cs.prefixKey(key), NomadQueryDefaults(ctx))
	// if err != nil {
	// 	return nil, err
	// }

	// if kv == nil {
	// 	return nil, fs.ErrNotExist
	// }

	// return kv.Value, nil

	return []byte("foo"), nil
}

// Delete a key from Nomad KV
func (cs NomadStorage) Delete(ctx context.Context, key string) error {
	cs.logger.Infof("deleting key %s from Nomad", key)

	// // first obtain existing keypair
	// kv, _, err := cs.NomadClient.KV().Get(cs.prefixKey(key), NomadQueryDefaults(ctx))
	// if err != nil {
	// 	return fmt.Errorf("%s: %w", err, fs.ErrNotExist)
	// }

	// if kv == nil {
	// 	return fs.ErrNotExist
	// }

	// // now do a Check-And-Set operation to verify we really deleted the key
	// if success, _, err := cs.NomadClient.KV().DeleteCAS(kv, nil); err != nil {
	// 	return errors.Wrapf(err, "unable to delete data for %s", cs.prefixKey(key))
	// } else if !success {
	// 	return errors.Errorf("failed to lock data delete for %s", cs.prefixKey(key))
	// }

	// TODO: DELETE a Nomad SV

	return nil
}

// Exists checks if a key exists
func (cs NomadStorage) Exists(ctx context.Context, key string) bool {
	// TODO: READ A NOMAD SV

	// kv, _, err := cs.NomadClient.KV().Get(cs.prefixKey(key), NomadQueryDefaults(ctx))
	// if kv != nil && err == nil {
	// 	return true
	// }
	// return false

	return true
}

// List returns a list with all keys under a given prefix
func (cs NomadStorage) List(ctx context.Context, prefix string, recursive bool) ([]string, error) {
	// TODO: LIST KEYS UNDER THE PREFIX

	// var keysFound []string

	// // get a list of all keys at prefix
	// keys, _, err := cs.NomadClient.KV().Keys(cs.prefixKey(prefix), "", NomadQueryDefaults(ctx))
	// if err != nil {
	// 	return keysFound, err
	// }

	// if len(keys) == 0 {
	// 	return keysFound, fs.ErrNotExist
	// }

	// // remove default prefix from keys
	// for _, key := range keys {
	// 	if strings.HasPrefix(key, cs.prefixKey(prefix)) {
	// 		key = strings.TrimPrefix(key, cs.Prefix+"/")
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
func (cs NomadStorage) Stat(ctx context.Context, key string) (certmagic.KeyInfo, error) {
	// kv, _, err := cs.NomadClient.KV().Get(cs.prefixKey(key), NomadQueryDefaults(ctx))
	// if err != nil {
	// 	return certmagic.KeyInfo{}, fmt.Errorf("unable to obtain data for %s: %w", cs.prefixKey(key), fs.ErrNotExist)
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

func (cs *NomadStorage) createNomadClient() error {
	// get the default config
	nomadCfg := nomad.DefaultConfig()
	if cs.Address != "" {
		nomadCfg.Address = cs.Address
	}
	if cs.Token != "" {
		nomadCfg.SecretID = cs.Token
	}
	// if cs.TlsEnabled {
	// 	nomadCfg.Scheme = "https"
	// }

	nomadCfg.TLSConfig.Insecure = cs.TlsInsecure

	// set a dial context to prevent default keepalive
	// nomadCfg.Transport.DialContext = (&net.Dialer{
	// 	Timeout:   time.Duration(cs.Timeout) * time.Second,
	// 	KeepAlive: time.Duration(cs.Timeout) * time.Second,
	// }).DialContext

	// create the Nomad API client
	nomadClient, err := nomad.NewClient(nomadCfg)
	if err != nil {
		return errors.Wrap(err, "unable to create Nomad client")
	}

	if _, err := nomadClient.Agent().NodeName(); err != nil {
		return errors.Wrap(err, "unable to ping Nomad")
	}

	cs.NomadClient = nomadClient
	return nil
}

func NomadQueryDefaults(ctx context.Context) *nomad.QueryOptions {
	// TODO: Set some of these
	opts := &nomad.QueryOptions{}
	return opts.WithContext(ctx)
}
