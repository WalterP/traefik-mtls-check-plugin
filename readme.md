# Traefik mTLS Check Plugin




The `traefik_mtls_check_plugin` package is a middleware plugin for Traefik that performs mutual TLS (mTLS) certificate validation. It allows Traefik to verify the client certificates presented during TLS handshake against a given CA certificate. The main difference between this middleware plugin and TLSOptions, that the plugin allows you to specify the response code and message that should be returned to the client in case of a failed validation.
On other hand, TLSOptions with clientAuthType: RequireAndVerifyClientCert will return ssl error to the client.
The combination of:
- TLSOptions with clientAuthType: RequestClientCert
- mTlsCheck middleware
- Errors middleware with custom errorPage

will allow you to return custom errorPage to the client in case of failed validation and the unwanted visitors will not be able to recognize that you are using mTLS.

## Installation

To use the `traefik_mtls_check_plugin`, you need to have Traefik installed and configured. Follow the instructions below to install the plugin:

1. Install Traefik: [Traefik Installation Guide](https://doc.traefik.io/traefik/getting-started/install-traefik/)
2. [Configure Traefik to use the plugin](https://plugins.traefik.io/install): Add the following lines to your Traefik configuration file (`traefik.toml` or `traefik.yaml`):


```toml
[experimental.plugins.mTlsCheck]
modulename = "github.com/WalterP/traefik-mtls-check-plugin"
version = "v0.1.0"
```

```yaml
experimental:
  plugins:
    mTlsCheck:
      modulename: "github.com/WalterP/traefik-mtls-check-plugin"
      version: "v0.1.0"
```

```cli
- "--experimental.plugins.mTlsCheck.modulename=github.com/WalterP/traefik-mtls-check-plugin"
- "--experimental.plugins.mTlsCheck.version=v0.1.0"
```

4. Restart Traefik: `traefik restart`

## Configuration

The `traefik_mtls_check_plugin` supports the following configuration options in the Traefik middleware configuration file:


- `responseCode`: The HTTP response code to return when the certificate validation fails.
- `caCert`: The CA certificate content as a string. (Optional if `caCertPath` is provided). Mostly for kubernetes usage. See example below.
- `caCertPath`: The path to the CA certificate file. (Optional if `caCert` is provided). Mostly for docker usage.
- `message`: The message to return in the response body when the certificate validation fails.

## Usage

Once the `traefik_mtls_check_plugin` is installed and configured, you can use it by adding it as a middleware to your Traefik routes. Here's an example of how to configure the plugin in your Traefik routing:


Docker:

```labels
- "traefik.http.middlewares.test-redirect.plugin.mTlsCheck.message=Not - Found"
- "traefik.http.middlewares.test-redirect.plugin.mTlsCheck.ResponseCode=404"
- "traefik.http.middlewares.test-redirect.plugin.mTlsCheck.caCertPath=/certs/mtls/ca.crt"
- "traefik.http.routers.whoami.middlewares=test-redirect"

```

Kubernetes:

```yaml
apiVersion: traefik.io/v1alpha1
kind: Middleware
metadata:
  name: mtls-check
spec:
  plugin:
    mTlsCheck:
      message: "Not-found"
      responseCode: 404
      caCert: "urn:k8s:secret:client-mtls:ca.crt"
```

Example with combination of custom errorPage:

```dockerlabels
- "traefik.http.routers.whoami.tls.options=mtls@file"
- "traefik.http.middlewares.test-redirect.plugin.mTlsCheck.message=Not - Found"
- "traefik.http.middlewares.test-redirect.plugin.mTlsCheck.ResponseCode=404"
- "traefik.http.middlewares.test-redirect.plugin.mTlsCheck.caCertPath=/certs/mtls/ca.crt"
- "traefik.http.middlewares.test-errors.errors.status=400-499"
- "traefik.http.middlewares.test-errors.errors.service=errorServer@docker"
- "traefik.http.routers.whoami.middlewares=test-errors,test-redirect"
- "traefik.http.services.whoami.loadbalancer.server.port=8082"
```

mtlsOptions used to request client certificate:
```markdown
    mtls:
      sniStrict: false
      clientAuth:
        clientAuthType: RequestClientCert
        caFiles:
          - /etc/traefik/ca.crt
```
