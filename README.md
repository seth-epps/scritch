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

### Custom Templates
You can define your own custom scratch templates and pass that as a source to the scratch command.

_All files need to have extension `.tpl`_

Eg, if you store a template `go_http` at `~/.scratch_templates`
```sh
❯ find -d  ~/.scratch_templates/go_http_nested
/Users/{user}/.scratch_templates/go_http/main.go.tpl
/Users/{user}/.scratch_templates/go_http/go.mod.tpl
/Users/{user}/.scratch_templates/go_http/pkg/models.go.tpl
/Users/{user}/.scratch_templates/go_http/pkg
/Users/{user}/.scratch_templates/go_http

❯ scritch scratch --source=~/.scratch_templates/go_http
2023/08/25 13:35:10 Created scratch at /Users/seth.epps/.scritch/scratch/go_http/708ebf03-6538-47a8-b9e3-0e36e8fb033b
```

### Options

#### `scratch-path` (default `~/.scritch/scratch`)
The generated scratch destination.

#### `custom-sources` (default [[~/.scritch/templates]])
A list of paths to look for source templates when you don't provide an absolute path. Eg,
```sh
❯ scritch scratch --source=go_http --custom-sources ~/.scratch_templates
2023/08/25 13:41:16 Created scratch at /Users/seth.epps/.scritch/scratch/go_http/645d0805-0023-4c2d-9053-f9406c65d5f7
```


#### `editor-command`
A command to execute on the resulting scratch path as `<command> <resulting-scratch-path>`.
Eg, if you have vs-code installed and you specify `code`, vs-code will open with the workspace
set to the newly created scratch.

#### config (default `~/.scritch/config.yaml`)
All options can be configured through a config yaml. Eg,

```yaml
custom-sources:
  - ~/.scritch/templates
  - ~/.scratch_templates
scratch-path: ~/dev/scratch/files
editor-command: code
```

## TODO
 - [] Support for template variables
 - [] ~Automatic starting of a new terminal at the scratch location~ Probably not achievable
 - [] TUI Client (maybe charm?; eg, UI for `go fmt`,`go run main.go`, etc.)
 - [] ...Major refactor...
 - [] Unit tests... 😬