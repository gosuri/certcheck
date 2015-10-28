#!/bin/sh

echo "generating server cert"
openssl req -new -nodes -x509 -out server.crt -keyout server.key -days 3650 -subj "/C=US/ST=CA/L=Earth/O=CertCheck Org/OU=Tech/CN=www.example.com/email=$1"
echo "generating client cert"
openssl req -new -nodes -x509 -out client.crt -keyout client.key -days 3650 -subj "/C=US/ST=CA/L=Earth/O=CertCheck Org/OU=Tech/CN=www.example.com/email=$1"
