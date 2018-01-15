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
	if you want to make global configuration paramters,
	add "_" before paramter name,
	otherwise means it is an private paramter
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

// 3) take your configurations map to echo
err = EchoAgent.Trusteeship(SettingMap)
if err != nil {
	log.Fatal(err)
}

// 4) detory map defore app down
defer EchoAgent.Destroy()
```
