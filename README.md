# TO DOs

- Figure out how this gets swapped in
- Remove errors library
- Remove any unnecessary libs/functionality
- Get a test running
- Test out locally
- Add acknowledgements

# Caddy 2 cluster / Certmagic TLS cluster support for Nomad Secure Variables

[Nomad Secure Variable](https://github.com/hashicorp/nomad) Storage for [Caddy](https://github.com/caddyserver/caddy) TLS data.

This cluster plugin enables Caddy 2 to store TLS data like keys and certificates as Nomad Secure Variables so you don't have to rely on a shared filesystem.
This allows you to use Caddy 2 in distributed environment and use a centralized storage for auto-generated certificates that is
shared between all Caddy instances.

With this plugin it is possible to use multiple Caddy instances with the same HTTPS domain for instance with DNS round-robin.

The version of this plugin in the master branch supports Caddy 2.0.0+ using CertMagic's [Storage Interface](https://pkg.go.dev/github.com/caddyserver/certmagic?tab=doc#Storage)

## Older versions

This will only work with Caddy 2.

## Docker Image

You can use my Docker image `mnomitch/caddy-tlsnomad` (https://hub.docker.com/r/mnomitch/caddy-tlsnomad) to get an image of Caddy including this Nomad TLS plugin build in. It is based on the official caddy:2 image and will be build automatically on changes in this repo.

## Configuration

### Caddy configuration

You need to specify `nomad` as the storage module in Caddy's configuration. This can be done in the config file of using the [admin API](https://caddyserver.com/docs/api).

JSON ([reference](https://caddyserver.com/docs/json/))

```
{
  "admin": {
    "listen": "0.0.0.0:2019"
  },
  "storage": {
    "module": "nomad",
    "address": "localhost:4646",
    "prefix": "caddytls",
    "token": "nomad-access-token",
  }
}
```

Caddyfile ([reference](https://caddyserver.com/docs/caddyfile/options))

```
{
    storage nomad {
           address      "127.0.0.1:4646"
           token        "nomad-access-token"
           timeout      10
           prefix       "caddytls"
           value_prefix "myprefix"
           aes_key      "nomadtls-1234567890-caddytls-32"
           tls_enabled  "false"
           tls_insecure "true"
    }
}

:443 {
}
```

### Nomad configuration

Because this plugin uses the official Nomad API client you can use all ENV variables like `nomad_HTTP_ADDR` or `nomad_HTTP_TOKEN`
to define your Nomad address and token. For more information see https://github.com/hashicorp/nomad/blob/master/api/api.go

Without any further configuration a running Nomad on 127.0.0.1:4646 is assumed.

There are additional ENV variables for this plugin:

- `CADDY_CLUSTERING_nomad_PREFIX` defines the prefix for the keys in the Secure Variable. Default is `caddytls`

### Nomad ACL Policy

To access Nomad you need a token with a valid ACL policy. Assuming you configured `cadytls` as your Secure Variable path prefix you can use the following settings:

```
namespace "default" {
  secure_variables {
    path "cadytls/*" {
      capabilities = ["write", "read", "destroy"]
    }
  }
}
```
