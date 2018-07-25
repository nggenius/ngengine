rem del .\pkg\* /s/q
go build -o ./bin/ngadmin.exe ./src/ngengine/apps/ngadmin
go build -o ./bin/test.exe ./src/ngengine/apps/test
go build -o ./bin/client.exe ./src/ngengine/apps/client
go build -o ./bin/login.exe ./src/ngengine/apps/ngadmin/services/login
go build -o ./bin/nest.exe ./src/ngengine/apps/ngadmin/services/nest
go build -o ./bin/region.exe ./src/ngengine/apps/ngadmin/services/region
go build -o ./bin/store.exe ./src/ngengine/apps/ngadmin/services/store
go build -o ./bin/world.exe ./src/ngengine/apps/ngadmin/services/world
