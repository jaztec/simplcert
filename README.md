[![Go Report Card](https://goreportcard.com/badge/github.com/jaztec/simplcert)](https://goreportcard.com/report/github.com/jaztec/simplcert)
![Status](https://github.com/jaztec/simplcert/actions/workflows/run-tests.yml/badge.svg)

# SimplCERT

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
$ simplcert verify
```

`root-crt` will display the root certificate as PEM encoded string to the terminal.
```bash
$ simplcert root-crt
```

`create` will create a certificate. You can use CLI flags or just run create and fill in 
the prompts. It is important to know the `--host` flag needs to be set to the domain name 
where the service will be reached. Or, if Docker is used, the `--host` flag should be set 
to the name of the docker container. If you need to support multiple hosts or add IP 
addresses this is supported. Just use comma's to separate the values like 
`--host "127.0.0.1,::1,localhost,hostname.tld"`
The command will give some additional prompts to clarify any values that are not provided.
```bash
$ simplcert create \
  --root-cert-path /path/to/root-ca \
  --host hostname.tld \
  --name "My server" \
  --days-valid 30 \
  --ecdsa \
  --is-server
```

### Usage examples

See [examples folder](examples) for some examples:

- [Go server/client](examples/go-server-client)
- [Rust server/client](examples/rust-server-client)
- [Rust server with Go client](examples/rust-server-go-client)

### Roadmap

| Status  | Target                        | Description                                                            |
|---------|-------------------------------|------------------------------------------------------------------------|
| &check; | Add Go server/client examples | Have an example folder displaying a working setup                      |
| &check; | Add more signing options      | RSA, DSA etc. (now only ecdsa is supported                             |
| &check; | Write output to file          | Add additional flag to write output to file instead of stdout          | 
| &check; | Add Rust example              | Add an example on how to use the certs in a Rust gRPC application      |
| &check; | Add cross language example    | Add an example on using a Rust server with a Go client and mTLS        |
| Open    | Add leaf certs                | Cross sign with dedicated leaf certs for server and client validations |
