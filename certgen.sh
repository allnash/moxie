#!/bin/bash
# Generate a Sample certificate
echo "Please type HOSTNAME example: *.hostname"
openssl req -newkey rsa:4096 \
            -x509 \
            -sha256 \
            -days 3650 \
            -nodes \
            -out server.crt \
            -keyout server.key