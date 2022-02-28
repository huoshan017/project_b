go build -o client_core_import/clientcore.a -buildmode=c-archive project_b/client_core_import
go build -o client_core_import/clientcore.dll -buildmode=c-shared project_b/client_core_import