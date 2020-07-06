#!/bin/bash
protoc -I protocol/ protocol/zbprotocol.proto --js_out=import_style=commonjs:jsweb/ --grpc-web_out=import_style=commonjs,mode=grpcwebtext:jsweb/
