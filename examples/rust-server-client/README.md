# Example application

This is an example application on how to use the generated certificates in use.

The certificates in the `certs` folder are pre-generated with the following commands.
The root certs path is set to the cert path inside this example so it will use the 
special root certificate for this example only.

For the server
```bash
$ simplcert create \
  --root-cert-path $PWD/certs \
  --name Server \
  --host server \
  --is-server \
  --output-path $PWD/certs
```

For the client
```bash
$ simplcert create \
  --root-cert-path $PWD/certs \
  --name Client \
  --host client \
  --output-path $PWD/certs
```

Run the example from this folder with docker-compose:
```bash
$ docker compose up
```

It will request a greeting from the server every 10 seconds over a with mTLS secured
connection to the server.
