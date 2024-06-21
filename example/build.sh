gcc -shared -fPIC -o libvalue.so libvalue.c
go build -o main.run main.go
ldd main.run
go build -o single.run single.go
ldd single.run