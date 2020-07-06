#!/bin/bash
protoc -I zbprotocol/ zbprotocol/zbprotocol.proto --go_out=plugins=grpc:.
