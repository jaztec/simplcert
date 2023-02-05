[![Go Reference](https://pkg.go.dev/badge/github.com/jaztec/simplcert.svg)](https://pkg.go.dev/github.com/jaztec/simplcert)

# Certificate manager

The certificate manager is meant to easily generate keys pairs for mTLS purposes. It will generate
it's own root certificate and use that to sign server and client certificates that can later be used
to secure gRPC over mTLS for instance.

### Usage

```bash
$ export SCM_CERTS_PATH=/path/to/directory/for/root/cert
# verify will check if the certs exist and generate a root cert if necessary 
$ simplcert verify
# root-crt will show the root certificate. This can be used as trusted root inside a gRPC client
$ simplcert root-crt
# create will generate a certificate signed by the root certificate. Use the flags to generate a 
# certificate for a specific purpose.
$ simplcert create --host hostname.tld --name "My server" --is_server
```

### Roadmap

|       | Target         | Description                                                            |
|-------|-------------------------|----------------------------------------------------------------|
| Open  | Add more signing options | DSA, ECDSA etc.                                              |
| Open  | Add leaf certs | Cross sign with dedicated leaf certs for server and client validations |
| Open  | Write output to file | Add additional flag to write output to file instead of stdout | 
