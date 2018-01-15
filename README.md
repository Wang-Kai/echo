# echo

This is an library which can help you to complete auto configurations. Echo library based [etcdv3](https://github.com/coreos/etcd), you can modify configurations in etcd, then echo can help you update or delete data in memory. **Don't need to reload app**.

The reason i write echo library is that my project need many static parameters, and those parameters are changeful. Besides, the workload of deploy a big project is much heavy, and it is easy to make mistakes if engineer is careless. Echo can solve this problem.

## Installation
```
$ go get -u -v github.com/Wang-Kai/echo
```

## Usage
```go
/*
	1) create your setting map (it should be map[string]string)
*/
var setting = map[string]string{
	"_name": "Wang",
	"age":   "26",
}

// 2) init echo agent
EchoAgent, err := echo.Init([]string{"localhost:2379"})
if err != nil {
	log.Fatal(err)
}

// 3) detory map defore app down
defer EchoAgent.Destroy()

// 4) take your configurations map to echo
err = EchoAgent.Trusteeship(SettingMap)
if err != nil {
	log.Fatal(err)
}
```

## Features

In order to meet a variety of application scenarios. For example you want to share same paramter in many apps, or you need use a paramter in many app with different value. Echo provide two kind paramters, global and private paramter.

* **Global paramter**: add "_" before paramter name, eg: '\_apiToken'
* **Private parameter**: **don't** prefix with '_', eg: "name"

I hope it can solve your problem, and enjoy it . 🤣
