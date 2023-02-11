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
$ simplcert create --host hostname.tld --name "My server" --is-server
```

See [examples folder](examples) for some examples

### Roadmap

| Status | Target                   | Description                                                            |
|--------|--------------------------|------------------------------------------------------------------------|
| Done   | Add examples             | Have an example folder displaying a working setup                      |
| Open   | Add more signing options | RSA, DSA etc. (now only ecdsa is supported                             |
| Open   | Add leaf certs           | Cross sign with dedicated leaf certs for server and client validations |
| Open   | Write output to file     | Add additional flag to write output to file instead of stdout          | 
