# grpcio
A basic gRPC server with middlewares, gRPC-Gateway, ReDoc, etc..

## Requirements
### protobuf
>  You can just install protoc.
> [protoc-3.17.3](https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-linux-x86_64.zip)

refer: [protobuf v3.17.3](https://github.com/protocolbuffers/protobuf/tree/v3.17.3)

install: 
1. Download [protobuf release v3.17.3](https://github.com/protocolbuffers/protobuf/releases/tag/v3.17.3), and install it by README.
2. Build from source.

### protoc plugins
> `protoc-gen-go`, `protoc-gen-go-grpc`, `protoc-gen-grpc-gateway`, and `protoc-gen-openapiv2`

install: `bash ./tools/install.sh`

## Test

## Compile
1. Compile protos: `bash ./scripts/gen_pb.sh`
2. Build **flyio**: `cd cmd/flyio && go build`
3. Start **flyio**: `./flyio -conf=./conf.yaml`

### gRPC-Gateway
Echo: `curl --data '{"msg": "hello"}' -H "Content-Type: application/json"   http://127.0.0.1:8080/apiv1/client_test_service/echo`

### ReDoc
Browser URL: http://127.0.0.1:8081/docs/

## Git
### Commits
A specification for adding human and machine readable meaning to commit messages.

refer: [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)
### .gitconfig
```
[user]
    name = yourname
    email = youremail
[push]
    default = simple
[diff]
    tool = vimdiff
[difftool]
	prompt = false
[credential]
    helper = store
[color]
    diff = auto
    status = auto
    branch = auto
    interactive = auto
    ui = true
    pager = true
[color "status"]
    added = green
    changed = red bold
    untracked = magenta bold
[color "branch"]
    remote = yellow
[merge]
    tool = vimdiff
[mergetool]
    prompt = false
[alias]
    d = difftool
    di = difftool
    co = checkout
    br = branch
    pr = pull --rebase
    dc = dcommit
    ci = commit
    st = status
    last = log -1 HEAD
    lg = log --color --graph --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit
    lg5 = lg -5
    lg10 = lg -10
    unstage = reset HEAD --
    s = "!f() { rev=${1-HEAD}; git difftool $rev^ $rev; }; f"
```

## References
- [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
- [redoc](https://github.com/Redocly/redoc)
- [go-microservice-demo](https://github.com/win5do/go-microservice-demo)
