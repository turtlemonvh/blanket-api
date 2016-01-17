package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"time"
)

/*
The most basic task is a command line task

https://github.com/boltdb/bolt
- only 1 process can access bolt at a time, so everything has to go through a db access process that starts first

http://stackoverflow.com/questions/2886719/unix-sockets-in-go
- example of using named sockets

http://lua-users.org/lists/lua-l/2013-12/msg00047.html
- docker does REST over sockets

https://github.com/docker/docker-py/blob/master/docker/client.py#L60
https://github.com/docker/docker-py/blob/master/docker/unixconn/unixconn.py
- in python

https://github.com/docker/docker/blob/master/opts/hosts.go#L24
https://github.com/docker/docker/blob/master/opts/hosts_unix.go
- in golang
https://github.com/docker/docker/blob/master/opts/hosts_windows.go
- windows uses the default tcp host
https://github.com/docker/docker/blob/master/opts/hosts_test.go
- tests for various host types

https://github.com/docker/docker/blob/master/api/server/server_unix.go
- creating the server
https://github.com/docker/docker/tree/master/api
- read through how they structure their code

https://github.com/docker/docker/blob/master/docs/reference/api/docker_remote_api_v1.22.md
- api docs

https://github.com/docker/docker/blob/master/docs/extend/plugin_api.md
- how plugins work
- they talk over sockets
- exclusively http POST requests sent from the main process
- internally docker uses events for most subsystems, so easy to farm off to eternal processes written in any language
- the template system would be a good first thing to use this way, since the communication is always 1 way

https://github.com/docker/docker/blob/master/docs/userguide/basics.md#bind-docker-to-another-hostport-or-a-unix-socket
- docker networking security tip
https://github.com/docker/docker/blob/master/docs/installation/debian.md#giving-non-root-access
- more networking security on debian
https://github.com/docker/docker/blob/master/docs/articles/security.md#docker-daemon-attack-surface
- more security

http://stackoverflow.com/questions/9029174/af-unix-equivalent-for-windows
- handling things like unix domain sockets on windows

https://github.com/docker/go-connections
- helper library to work with network connections uses by docker
https://godoc.org/github.com/docker/go-connections/sockets#NewUnixSocket
- unix socket

http://stackoverflow.com/questions/2135595/creating-a-socket-restricted-to-localhost-connections-only
- can create localhost only sockets by binding to '127.0.0.1'
- this ensures that the user must be on that server
http://stackoverflow.com/questions/2135595/creating-a-socket-restricted-to-localhost-connections-only#comment2075905_2135752
- unix domain socket vs localhost only
- both work fine, 1 handles file system
http://stackoverflow.com/questions/2205073/how-to-create-java-socket-that-is-localhost-only
- java version

http://enterprisewebbook.com/ch8_websockets.html
- use websockets for UI parts of api, auto updates, etc.
- this will allow that part of the application to be easy switched for another websocket server (e.g. python)
- also nice for dashboard updates
https://github.com/docker/docker/blob/master/docs/reference/api/docker_remote_api_v1.22.md#attach-to-a-container-websocket
- e.g. usage for attaching to container

http://tldp.org/LDP/abs/html/tabexpansion.html
- basic command line completion
https://github.com/docker/docker/blob/master/contrib/completion/bash/docker
- docker's fancy version
https://github.com/kislyuk/argcomplete
- python can do it with some fancy packages

https://github.com/spf13/cobra/blob/master/bash_completions.md
- bash completion
http://unix.stackexchange.com/questions/149730/how-do-command-line-tools-have-their-own-autocomplete-list
- how it works

*/

// Arguments are typed per task
type TaskType struct {
	Id            string            `json:"id"`
	CreatedTs     int64             `json:"createdTs"`
	LastUpdatedTs int64             `json:"lastUpdatedTs"`
	Type          string            `json:"type"`
	DefaultEnv    map[string]string `json:"defaultEnv"`
	ConfigPath    string            `json:"configPath"` // path to config directory
}

func NewTaskType(typeName string, defaultEnv map[string]string) (TaskType, error) {
	id := uuid.NewV4().String()
	return TaskType{
		Id:            id,
		CreatedTs:     time.Now().Unix(),
		LastUpdatedTs: time.Now().Unix(),
		Type:          typeName,
		DefaultEnv:    defaultEnv,
		ConfigPath:    fmt.Sprintf("tasks/%s/", id),
	}, nil
}

func (t *TaskType) String() string {
	return fmt.Sprintf("%s (%s) [%d]", t.Type, t.Id, t.CreatedTs)
}

func (t *TaskType) ToJSON() (string, error) {
	bts, err := json.Marshal(t)
	return string(bts), err
}

// Tasks inherit all the environment params of a tasktype + more
func (t *TaskType) NewTask(envOverrides map[string]string) (Task, error) {
	return Task{
		CreatedTs:     time.Now().Unix(),
		LastUpdatedTs: time.Now().Unix(),
		TypeId:        t.Id,
		ExecEnv:       envOverrides,
	}, nil
}

type Task struct {
	CreatedTs     int64             `json:"createdTs"`
	LastUpdatedTs int64             `json:"lastUpdatedTs"`
	TypeId        string            `json:"type"`
	ExecEnv       map[string]string `json:"defaultEnv"`
}

func (t *Task) String() string {
	return fmt.Sprintf("%d [%d]", t.TypeId, t.CreatedTs)
}

func (t *Task) ToJSON() (string, error) {
	bts, err := json.Marshal(t)
	return string(bts), err
}

func (t *Task) Execute() error {
	return nil
}