#!/bin/bash

go build -mod=mod -o ./bin/client project_b/client
go build -mod=mod -o ./bin/game_server project_b/server/game_server