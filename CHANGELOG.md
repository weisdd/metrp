# CHANGELOG

## 0.5.4

- Bump deps:
  - Alpine to `3.16.0`;
  - `caarlos0/env/v6` to `6.9.3`;
  - `go-yaml/yaml` to `3.0.1`;
- Bump go to `1.18.3`.

## 0.5.3

- Automatically set `GOMAXPROCS` to match Linux container CPU quota via [uber-go/automaxprocs](https://github.com/uber-go/automaxprocs). Enabled by default, can be turned off via `METRP_SET_GOMAXPROCS: false`.

## 0.5.2

- Add graceful shutdown;
- Bump alpine to 3.15.1;
- Migrate to go 1.18 (including `net/netip`).

## 0.5.1

- Simplify creation of insecure transport;
- Migrate to go 1.17;
- Bump alpine version;
- Add GitHub workflow to automatically publish images;
- Add dependabot alerts.

## 0.5.0

- Added support for forwarding bearer token (`/var/run/secrets/kubernetes.io/serviceaccount/token`). It's enabled by default and controlled via `METRP_USE_TOKEN` env;
- Bugfix: do not forward basic auth credentials to upstreams.

## 0.4.0

- Added support for `METRP_PREFERRED_IPV4` (useful for Downward API).

## 0.3.0

- Added custom transport to enforce `InsecureSkipVerify: true` (useful for kube-controller-manager).

## 0.2.1

- Code refactoring;
- Now, only GET method is allowed;
- URLs in metrp.yaml are now validated;
- Improved README.md.

## 0.2.0

- Added `METRP_PREFERRED_IPV4_CIDR`, so now it's possible to bind web-servers only to an interface that belongs to a specific prefix, e.g. `192.168.0.0/24`.

## 0.1.1

- Improve error handling.

## 0.1.0

- Initial release.
