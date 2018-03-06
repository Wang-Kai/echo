# echo

It is a tool to manage your configuration parameters based [etcd](https://github.com/coreos/etcd) k/v store.


![](./images/echo_design.png)

### Main functions

- Read all configuration paramters while app init , save them in memory
- Automatic update configuration data in memory , while add or delete value under special etcd directory


### Getting started
```shell
$ go get -u github.com/Wang-Kai/echo
```

 ```golang
package misc

import (
	"log"

	"github.com/Wang-Kai/echo"
)

var (
	AppConf   echo.Config
)

const (
	etcdURI        = "http://127.0.0.1:2379"
	configParamDir = "my_app/"
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
}
 
 ```