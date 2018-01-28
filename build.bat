del .\pkg\* /s/q
go build -o ./bin/ngadmin.exe ./src/ngengine/apps/ngadmin
go build -o ./bin/test.exe ./src/ngengine/apps/test
go build -o ./bin/client.exe ./src/ngengine/apps/client