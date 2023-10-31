#!/bin/bash

KEY_DIR="$(dirname "$0")/../dev/keys"
mkdir -p "$KEY_DIR"

openssl ecparam -name prime256v1 -genkey -noout -out "$KEY_DIR/private.ec.key"
openssl ec -in "$KEY_DIR/private.ec.key" -pubout -out "$KEY_DIR/public.pem"