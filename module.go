package storagenomad

import (
	"os"
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/certmagic"
)

func init() {
	caddy.RegisterModule(NomadStorage{})
}

func (NomadStorage) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "caddy.storage.nomad",
		New: func() caddy.Module {
			return New()
		},
	}
}

// Provision is called by Caddy to prepare the module
func (cs *NomadStorage) Provision(ctx caddy.Context) error {
	cs.logger = ctx.Logger(cs).Sugar()
	cs.logger.Infof("TLS storage is using Nomad at %s", cs.Address)

	if prefix := os.Getenv(EnvNamePrefix); prefix != "" {
		cs.Prefix = prefix
	}

	if valueprefix := os.Getenv(EnvValuePrefix); valueprefix != "" {
		cs.ValuePrefix = valueprefix
	}

	return cs.createNomadClient()
}

func (cs *NomadStorage) CertMagicStorage() (certmagic.Storage, error) {
	return cs, nil
}

// UnmarshalCaddyfile parses plugin settings from Caddyfile
// storage nomad {
//     address      "127.0.0.1:8500"
//     token        "nomad-access-token"
//     timeout      10
//     prefix       "caddytls"
//     value_prefix "myprefix"
//     tls_enabled  "false"
//     tls_insecure "true"
// }
func (cs *NomadStorage) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		key := d.Val()
		var value string

		if !d.Args(&value) {
			continue
		}

		switch key {
		case "address":
			if value != "" {
				parsedAddress, err := caddy.ParseNetworkAddress(value)
				if err == nil {
					cs.Address = parsedAddress.JoinHostPort(0)
				}
			}
		case "token":
			if value != "" {
				cs.Token = value
			}
		case "timeout":
			if value != "" {
				timeParse, err := strconv.Atoi(value)
				if err == nil {
					cs.Timeout = timeParse
				}
			}
		case "prefix":
			if value != "" {
				cs.Prefix = value
			}
		case "value_prefix":
			if value != "" {
				cs.ValuePrefix = value
			}
		case "tls_enabled":
			if value != "" {
				tlsParse, err := strconv.ParseBool(value)
				if err == nil {
					cs.TlsEnabled = tlsParse
				}
			}
		case "tls_insecure":
			if value != "" {
				tlsInsecureParse, err := strconv.ParseBool(value)
				if err == nil {
					cs.TlsInsecure = tlsInsecureParse
				}
			}
		}
	}
	return nil
}
