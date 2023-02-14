[![Go Reference](https://pkg.go.dev/badge/github.com/jaztec/simplcert.svg)](https://pkg.go.dev/github.com/jaztec/simplcert)

# Certificate manager

The certificate manager is meant to easily generate keys pairs for mTLS purposes. It will generate
it's own root certificate and use that to sign server and client certificates that can later be used
to secure gRPC over mTLS for instance.

### Usage

Exporting the root certificate path is handy so you don't have to provide it to every call
```bash
$ export SCM_ROOT_CERT_PATH=/path/to/directory/for/root/cert
````

`verify` will check if the `root` certificate exists and if not, will create one 
```bash
# verify will check if the certs exist and generate a root cert if necessary 
$ simplcert verify
```

`root-crt` will display the root certificate as PEM encoded string to the terminal.
```bash
# root-crt will show the root certificate. This can be used as trusted root inside a gRPC client
$ simplcert root-crt
```

`create` will create a certificate. You can use CLI flags or just run create and fill in 
the prompts. It is important to know the `--host` flag needs to be set to the domain name 
where the service will be reached. Or, if Docker is used, the `--host` flag should be set 
to the name of the docker container.
```bash
# create will generate a certificate signed by the root certificate. Use the flags to generate a 
# certificate for a specific purpose.
$ simplcert create \
  --root-cert-path /path/to/root-ca \
  --host hostname.tld \
  --name "My server" \
  --is-server
```

### Usage examples

See [examples folder](examples) for some examples:

- [Go server/client](examples/go-server-client)
- [Rust server/client](examples/rust-server-client)

### Roadmap

| Status  | Target                   | Description                                                            |
|---------|--------------------------|------------------------------------------------------------------------|
| &check; | Add examples             | Have an example folder displaying a working setup                      |
| &check; | Add more signing options | RSA, DSA etc. (now only ecdsa is supported                             |
| &check; | Write output to file     | Add additional flag to write output to file instead of stdout          | 
| &check; | Add Rust example         | Add an example on how to use the certs in a Rust gRPC application      |
| Open    | Add leaf certs           | Cross sign with dedicated leaf certs for server and client validations |
