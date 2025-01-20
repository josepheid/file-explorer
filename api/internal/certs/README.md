# Generating certs

Follow your OS specific instructions to install `mkcert` and `nss` (if using firefox):
https://github.com/FiloSottile/mkcert?tab=readme-ov-file#installation

Then run:

```
mkcert -install
mkcert localhost
```

This should generate `localhost.pem` and `localhost-key.pem` in the `certs` directory.
