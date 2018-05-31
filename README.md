# echo[![GoDoc](https://godoc.org/github.com/Wang-Kai/echo?status.svg)](https://godoc.org/github.com/Wang-Kai/echo)

Echo is a configuration manager based etcd .

![](./images/echo_design.png)

### Main functions

- Read paramters under special directory from etcd easily
- Automatic update configuration data in memory, while add, update or delete value under special etcd directory in etcd

### Getting started
```shell
# run etcd on docker
docker pull etcd
docker run -d --name echo_etcd -p 2379:2379 etcd  

# set some example k/v paramters
docker exec -it echo_etcd /bin/sh

ETCDCTL_API=3 etcdctl put "hi/echo/name" kim
ETCDCTL_API=3 etcdctl put "hi/echo/age" 18

# download echo lib
$ go get -u github.com/Wang-Kai/echo

go test -v -run=GetConf
```

 ```golang
package misc

import (
	"log"

	"github.com/Wang-Kai/echo"
)

var AppConf echo.Config

const (
	etcdURI        = "http://127.0.0.1:2379"
	configParamDir = "my_config/"
)

// load config paramters
func init() {
	echoAgent, err := echo.New(etcdURI)
	if err != nil {
		log.Fatal(err)
	}

	config, err := echoAgent.GetConf(configParamDir)
	if err != nil {
		log.Fatal(err)
	}
	
	// export global configuration
	AppConf = config

	// use your config any where
	//  ...
}
  ```