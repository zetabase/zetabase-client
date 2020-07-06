#!/bin/bash
protoc -I zbprotocol/ zbprotocol/zbprotocol.proto --js_out=import_style=commonjs:jsweb/ --grpc-web_out=import_style=commonjs,mode=grpcwebtext:jsweb/
