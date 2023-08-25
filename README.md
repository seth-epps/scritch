# Scritch
CLI to generate scratch files for different programming languages to address my frustration with the lack of features in the go playground (no autocomplete, parameter/quick info, manually adding dependencies, limited applicability, etc.)

Currently just a hack so a lot will probably change.

## Build
```sh
❯ make build
go generate ./...
go build -o dist/scritch
```

## Install
```sh
❯ go install .

❯ which scritch
/Users/seth.epps/go/bin/scritch

❯ scritch scratch go --variant=http
2023/08/22 16:08:34 Created scratch at /Users/seth.epps/.scritch/go/http/0e584181-8c34-4153-9a72-3c6dccbeb0bc
```

## Usage
### Default scratch files (Hello world)
```sh
❯ ./scritch scratch go
2023/08/22 15:45:55 Created scratch at /Users/seth.epps/.scritch/go/default/11f4467b-87f7-455c-a1a4-fdb91a4c7afd

❯ cd /Users/seth.epps/.scritch/go/default/11f4467b-87f7-455c-a1a4-fdb91a4c7afd

❯ go run main.go
Hello!
```

### Scratch Variants
Eg, the go http variant
```sh
❯ ./dist/scritch scratch go --variant=http
2023/08/22 16:03:52 Created scratch at /Users/seth.epps/.scritch/go/http/b086906a-8140-4925-814b-eaeed1130ffe

❯ cd /Users/seth.epps/.scritch/go/http/b086906a-8140-4925-814b-eaeed1130ffe

❯ go run main.go &
[1] 60461

❯ curl localhost:8080/
2023/08/22 16:05:21 Recieved req [&{GET / HTTP/1.1 1 1 map[Accept:[*/*] User-Agent:[curl/8.1.2]] {} <nil> 0 [] false localhost:8080 map[] map[] <nil> map[] 127.0.0.1:58186 / <nil> <nil> <nil> 0x140000e21e0}]
{"message":"Hello!"}

## ❯ kill %1
## [1]  + 60461 terminated  go run main.go
```

## TODO
 - [] Support for template variables
 - [] ~Automatic starting of a new terminal at the scratch location~ Probably not achievable
 - [] TUI Client (maybe charm?; eg, UI for `go fmt`,`go run main.go`, etc.)
 - [] ...Major refactor...
 - [] Unit tests... 😬