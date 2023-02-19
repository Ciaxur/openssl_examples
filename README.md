# nginx reverse proxy with mTLS
This is an example of how certs are generated and what options are expected BY
nginx to allow client & server mTLS to work, as server must have the server option and likewise
with the client.


# Certificate Generation
NOTE: For no password, pass in `-nodes`.

REALLY good resource: https://www.golinuxcloud.com/revoke-certificate-generate-crl-openssl/

## Generating DH Params file
```sh
openssl dhparam -out dh.pem 4096
```

## Generating keys
```sh
# ED25519
> openssl genpkey -algorithm ed25519 -out secret.key

# RSA
> openssl genrsa -out secret.key 4096
```


## Generate CSR from config
```sh
> openssl req -new -in secret.key -config csr.config -out secret.csr
```

## Signing CSR (gen Certificate)
```sh
> openssl x509 -req \
  -days 365 \
  -CA trusted_ca.crt \
  -CAkey trusted_ca.key \
  -CAcreateserial \
  -extfile csr_file.config \
  -extensions usr_crt \
  -in secret.csr \
  -sha256 \
  -out secret.crt
```

## Self-Signing Certificate
```sh
> openssl req \
  -new \
  -x509 \
  -days 3650 \
  -config ./opensslv3.conf \
  -key ./secret.key \
  -sha256 \
  -out ./secret.crt
```

## Certificate Extentions
https://www.golinuxcloud.com/add-x509-extensions-to-certificate-openssl/

These are used for OpenSSLv3 options, which allow setting "purpose" of the certificate.
See `man openssl-verification-options` for more detail.

See https://www.openssl.org/docs/manmaster/man5/x509v3_config.html for KeyUsage options.
