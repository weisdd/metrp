# CHANGELOG

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
