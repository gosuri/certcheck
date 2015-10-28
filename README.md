# certcheck

certcheck is simple utitily to check TLS authentication for a server and client

## Building

Install GoLang and ensure $GOPATH is configured. Get the source:

```sh
$ mkdir -p $GOPATH/src/github.com/gosuri/certcheck
$ git clone https://github.com/gosuri/certcheck.git $GOPATH/src/github.com/gosuri/certcheck
```

To build:

```sh
$ cd $GOPATH/src/github.com/gosuri/certcheck
$ go build
```

The above command will create a binary `certcheck` in the source directory

## Generate certificates

The below command will generate cert and key pair for `server` and `client`. Optionally, replace `foo@example.com` with an email.

```sh
./gencerts.sh foo@example.com
```

## Run 

To run, simply execute `certcheck`

```sh
./certcheck -h
```

## Usage

```sh
./certcheck -h

Usage of ./certcheck:
  -addr=":8080": Address to listen on
  -tls-client-cert="client.crt": Path for client TLS cert
  -tls-client-key="client.key": Path for client TLS key
  -tls-server-cert="server.crt": Path for server TLS cert
  -tls-server-key="server.key": Path for server TLS key
```
