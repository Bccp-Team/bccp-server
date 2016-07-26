Bccp
====

Bccp is a CI for high number of builds and repositories.


How to use
==========

Generate private key (.key)
--------------------------

```sh
# Key considerations for algorithm "RSA" ≥ 2048-bit
openssl genrsa -out server.key 2048

# Key considerations for algorithm "ECDSA" ≥ secp384r1
# List ECDSA the supported curves (openssl ecparam -list_curves)
openssl ecparam -genkey -name secp384r1 -out server.key
```

##### Generation of self-signed(x509) public key (PEM-encodings
`.pem`|`.crt`) based on the private (`.key`)

```sh
openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
```

Compile
-------

```
make
```

Run
---

```
bccp
```

API
===

See [API.md](https://github.com/Bccp-Team/bccp-server/blob/master/API.md)

Note: If you want to communicate with the API take care of the final '/' of
your URLs. None of the API's URLs have one.
