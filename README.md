# Casbin Authz

Casbin Authz is a middleware plugin for [Traefik](https://github.com/traefik/traefik) for authorization.

## Configuration

### Static Config for traefik

```

```toml
pilot:
  token: xxxxx

experimental:
    plugins:
        plugindemo:
            moduleName: "github.com/casbin/traefik-authz"
            version: "v0.0.1"
```