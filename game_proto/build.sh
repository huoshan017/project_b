#!/bin/bash

protoc --go_out=./ --proto_path=./ game.proto common.proto error.proto