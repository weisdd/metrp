# metrp

Metrics Reverse Proxy (metrp) is a trivial reverse-proxy that helps to publish metrics endpoints that are not directly reachable (e.g. endpoint is bound to localhost).

## Key features

* easy to configure;
* tiny CPU and memory footprint (e.g., with a few endpoints, it stays under a few `m` and `15Mi`);
* it's intentionally prohibited to use non-GET methods or pass any additional parameters to endpoints;
* optional basic authentication;
* can optionally forward its bearer token (k8s service account);
* you can specify a preferred prefix, so metrp will listen only on an interface with an IP that belongs to the prefix.

## Current limitations

* the preferred prefix feature currently support only IPv4.

## Similar projects

### pambrose/prometheus-proxy

[Link](https://github.com/pambrose/prometheus-proxy)

Pluses:

* In an environment, where you cannot afford to open a port on a Kubernetes node (=> let an agent to use host network), prometheus-proxy is a way to go.

Minuses:

* High memory footprint due to complex architecture and Java runtime.

### Caddy / nginx

The same thing could be done using a reverse-proxy such as `caddy`, `nginx`, or others, though they seem to be a bit heavy for the task and more difficult to configure.

## Installation

Docker images are published on [ghcr.io/weisdd/metrp](https://github.com/weisdd/metrp/pkgs/container/metrp).

## Configuration

### Environment variables

| Variable                    | Default Value | Description                                                  |
| --------------------------- | ------------- | ------------------------------------------------------------ |
| `METRP_PREFERRED_IPV4`      |               | Allows to bind web-server only to a specific IP, e.g. `192.168.0.1`. |
| `METRP_PREFERRED_IPV4_PREFIX` |               | Allows to bind web-server only to an interface that belongs to a specific prefix, e.g. `192.168.0.0/24`. |
| `METRP_SET_GOMAXPROCS` | `true` | Automatically set `GOMAXPROCS` to match Linux container CPU quota via [uber-go/automaxprocs](https://github.com/uber-go/automaxprocs). |
| `METRP_PORT`                | `8080`        | Port the web server will listen on.                          |
| `METRP_CONFIG_PATH`         | `metrp.yaml`  | Path to a file with a list of endpoints                      |
| `METRP_BASIC_AUTH`          | `false`       | Whether to enable basic authentication.                      |
| `METRP_USERNAME`            |               | Username to require if basic authentication is enabled.      |
| `METRP_PASSWORD`            |               | Password to require if basic authentication is enabled.      |
| `METRP_USE_TOKEN` | `true` | Whether to forward k8s token to upstream (useful for `kube-apiserver` and other upstreams that require authorization). Be careful though with untrusted upstreams! |
| `METRP_READ_TIMEOUT`        | `10s`         | `ReadTimeout` covers the time from when the connection is accepted to when the request body is fully read (if you do read the body, otherwise to the end of the headers). [More details](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/) |
| `METRP_SHUTDOWN_TIMEOUT` | `20s` | Maximum amount of time to wait for all connections to be closed. [More details](https://pkg.go.dev/net/http#Server.Shutdown) |
| `METRP_WRITE_TIMEOUT`       | `10s`           | `WriteTimeout` normally covers the time from the end of the request header read to the end of the response write (a.k.a. the lifetime of the ServeHTTP). [More details](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/) |

### metrp.yaml syntax

```yaml
apiserver:               "https://127.0.0.1:6443/metrics"
etcd:                    "http://127.0.0.1:2381/metrics"
kube-controller-manager: "https://127.0.0.1:10257/metrics"
kube-proxy:              "http://127.0.0.1:10249/metrics"
netdata:                 "http://127.0.0.1:19999/api/v1/allmetrics?format=prometheus"
```

With this configuration, you'll get the following list of URIs:

* `/metrics/apiserver`;
* `/metrics/etcd`;
* `/metrics/kube-controller-manager`;
* `/metrics/kube-proxy`;
* `/metrics/netdata`.
