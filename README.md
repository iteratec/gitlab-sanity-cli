# GitLab Sanity CLI

GitLab Sanity CLI is a command line tool for cleanup GitLab server.

## Motivation

Why this tool ?

This CLI was made to automate sanity jobs (like clean old unused GitLab Runners) which is currently not possible by GitLab server WebUI.

The main purpose was to save time for cleanup tousend of runners and group-runners.

## Features

The CLI is currently able to list, remove and archive the resources: users, projects, runners, grouprunners. 

### Parameter Matrix

|Parameter|Type|Default Value|Description|
|---|---|---|---|
|-u, --url|`string`|''|The Gitlab API URL|
|--insecure|`boolean`|false|Skip certificate Verfication for Gitlab API URL|
|-t, --token|`string`|''|The GitLab API Access Token|
|-o, --operation|`string`|''|Action to run (see below)|
|-r, --resource|`string`|''|GitLab Resource to interact with|
|-i, --identifier|`int`|-1|Specific Resource ID|
|-a, --age|`int`|36|Filter by last activity in months|
|-q, --query|`string`|''|Search by name|
|-s, --state|`string`|''|Filter list by state|
|-d, --dry-run|`boolean`|false|Dry run, does not change/delete any resources|
|-n, --num-concurrent-api-calls|`int`|100|Limit the amount of concurrent API Calls|


| Action | Resource | Query filter applicable | Age filter applicable | Status filter applicable | Example |
|---|---|---|---|---|---|
|list|user|YES|-|-| List user with name admin: <br> `gitlab-sanity-cli -o list -r user -q admin`|
|list|project|YES|YES|-| List projects older two years: <br> `gitlab-sanity-cli -o list -r project -a 24`|
|list|runner|YES|-|YES| List docker based runner: <br> `gitlab-sanity-cli -o list -r runner -q docker`|
|list|groupRunner|YES|-|YES| List online kubernetes based runner: <br> `gitlab-sanity-cli -o list -r groupRunner -q kubernetes -s online` |
|delete|user|-|-|-| <b>Delete is not capable on users</b> |
|delete|project|-|-|-| Remove project with ID 123: <br>  `gitlab-sanity-cli -o delete -r project -i 123`|
|delete|runner|-|-|-| Remove runner with ID 123: <br>  `gitlab-sanity-cli -o delete -r runner -i 123`|
|delete|groupRunner|-|-|-| Remove runner with ID 123: <br> `gitlab-sanity-cli -o delete -r groupRunner -i 123`|
|delete-all|user|-|-|-|<b>Delete-All is not capable on users</b> |
|delete-all|project|YES|YES|-| Remove all projects with name testing: <br>`gitlab-sanity-cli -o delete-all -r project -a 0 -q testing` <br><br> Remove all projects older than five years: <br> `gitlab-sanity-cli -o delete-all -r project -a 60`|
|delete-all|runner|YES|-|YES| Remove all offline runner: <br>`gitlab-sanity-cli -o delete-all -r runner -s offline`|
|delete-all|groupRunner|YES|-|YES| Remove all groupRunner (offline and online): <br>`gitlab-sanity-cli -o delete-all -r groupRunner` |
|archive|project|-|-|-|Archive project with ID 123:<br>`gitlab-sanity-cli -o archive -r project -i 123`|
|archive-all|project|YES|YES|-|Archive project with name testing:<br>`gitlab-sanity-cli -o archive-all -r project -q testing -a 0`|

## How to run

### Requirements

- [GitLab Access Token with api, repository access](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html)

### Download binary version

On Linux/MacOS/FreeBSD

```sh
export OS=`uname -s`
curl -L -O https://github.com/iteratec/gitlab-sanity-cli/releases/latest/download/gitlab-sanity-cli.${OS}
curl -L -O https://github.com/iteratec/gitlab-sanity-cli/releases/latest/download/gitlab-sanity-cli.${OS}.sha256
```

On Windows (open Powershell or Cmd and run follow commands)

```cmd
curl -L -O https://github.com/iteratec/gitlab-sanity-cli/releases/latest/download/gitlab-sanity-cli.windows
curl -L -O https://github.com/iteratec/gitlab-sanity-cli/releases/latest/download/gitlab-sanity-cli.windows.sha256
```

### Verify Download

On Linux

```sh
sha256sum gitlab-sanity-cli.${OS}
```

On MacOS/FreeBSD

```sh
shasum -a 256 gitlab-sanity-cli.${OS}
```

On Windows

```cmd
CertUtil -hashfile gitlab-sanity-cli.windows SHA256
```

```powershell
Get-FileHash gitlab-sanity-cli.windows -Algorithm SHA256
```

### Extract Binary

On Linux/MacOS/Windows

```sh
tar xvzf gitlab-sanity-cli.${OS}.tar.gz
chmod 0755 ./gitlab-sanity-cli
```

On Windows

```sh
powershell -command "Expand-Archive -Force 'gitlab-sanity-cli.windows.zip' '.'"
```


See [Parameter Matrix](#parameter-matrix) from above for examples

## How to run from source

### Requirements

- [go](https://golang.org)
- [GitLab Access Token with api, repository access](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html)

First Run

```sh
go get -d -v ./...
```

Run command without building binary

```sh
go run cmd/main.go -h
```

## How to build 

### Requirements

- [go](https://golang.org)

Use make to create the binaries

For Windows x86_64

```sh
make windows
```

For Linux x86_64

```sh
make linux
```

For MacOS x86_64

```sh
make darwin
```

For FreeBSD x86_64

```sh
make freebsd
```

For any other OS and Architecture:

See https://golang.org/doc/install/source#environment)

```sh
#
# MacOS (M1/arm64) Example
#
export target_os="darwin"
export target_arch="arm64"
env GOOS=${target_os} GOARCH=${target_arch} go build -ldflags "-extldflags '-static'" -o ./gitlab-sanity-cli.${target_os}.${target_arch} cmd/main.go
```


## Update go modules

```sh
# List all used modules
go list -m all

# List all available versions from module
go list -m -versions github.com/xanzy/go-gitlab

# Get specific version from module
go get github.com/xanzy/go-gitlab@v0.50.4
```

# Use the Code

[see architecture](architecture.md)

# Sources

- [GitLab API](https://docs.gitlab.com/ee/api)
- [go-gitlab](https://github.com/xanzy/go-gitlab)
