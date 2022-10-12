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
func (ns *NomadStorage) Provision(ctx caddy.Context) error {
	ns.logger = ctx.Logger(ns).Sugar()
	ns.logger.Infof("Version 0.11 - TLS storage is using Nomad at %s", ns.Address)

	if prefix := os.Getenv(EnvNamePrefix); prefix != "" {
		ns.Prefix = prefix
	}

	if valueprefix := os.Getenv(EnvValuePrefix); valueprefix != "" {
		ns.ValuePrefix = valueprefix
	}

	return ns.createNomadClient()
}

func (ns *NomadStorage) CertMagicStorage() (certmagic.Storage, error) {
	return ns, nil
}

// UnmarshalCaddyfile parses plugin settings from Caddyfile
//
//	storage nomad {
//	    address      "http://127.0.0.1:8500"
//	    token        "nomad-access-token"
//	    timeout      10
//	    prefix       "caddytls"
//	    value_prefix "myprefix"
//	    tls_enabled  "false"
//	    tls_insecure "true"
//	}
func (ns *NomadStorage) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		key := d.Val()
		var value string

		if !d.Args(&value) {
			continue
		}

		switch key {
		case "address":
			if value != "" {
				ns.Address = value
			}
		case "token":
			if value != "" {
				ns.Token = value
			}
		case "timeout":
			if value != "" {
				timeParse, err := strconv.Atoi(value)
				if err == nil {
					ns.Timeout = timeParse
				}
			}
		case "prefix":
			if value != "" {
				ns.Prefix = value
			}
		case "value_prefix":
			if value != "" {
				ns.ValuePrefix = value
			}
		case "tls_enabled":
			if value != "" {
				tlsParse, err := strconv.ParseBool(value)
				if err == nil {
					ns.TlsEnabled = tlsParse
				}
			}
		case "tls_insecure":
			if value != "" {
				tlsInsecureParse, err := strconv.ParseBool(value)
				if err == nil {
					ns.TlsInsecure = tlsInsecureParse
				}
			}
		}
	}
	return nil
}
