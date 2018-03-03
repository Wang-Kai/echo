package echo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type Config map[string]string

type Echo struct {
	Configs map[string]Config
}

var (
	cli = new(clientv3.Client)
)

/*
	New make link to etcd server, and create Echo instance

	@endponters: the url of etcd server
*/
func New(endponters ...string) (*Echo, error) {
	var err error
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   endponters,
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	// init echo instance' Configs attribute
	echo := &Echo{}
	echo.Configs = make(map[string]Config)
	return echo, nil
}

/*
	GetConf get all k/v from etcd prefix with special etcdDir, and save k/v in map

	@etcdDir: the dir of config keys prefix with
*/
func (e *Echo) GetConf(etcdDir string) (map[string]string, error) {
	resp, err := cli.Get(context.Background(), etcdDir, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var config = make(Config)

	for _, ev := range resp.Kvs {
		key := fmt.Sprintf("%s", ev.Key)
		val := fmt.Sprintf("%s", ev.Value)
		key = removeDir(key, 1)
		config[key] = val
	}

	e.Configs[etcdDir] = config
	go e.watchConfDir(etcdDir)

	return config, nil
}

/*
	watchConfDir watch special etcd's dir and update map value in memory

	@etcdDir the dir that will be watched
*/
func (e *Echo) watchConfDir(etcdDir string) {
	defer cli.Close()

	config := e.Configs[etcdDir]

	rch := cli.Watch(context.Background(), etcdDir, clientv3.WithPrefix())
	for wresp := range rch {

		for _, ev := range wresp.Events {
			key := fmt.Sprintf("%s", ev.Kv.Key)
			key = removeDir(key, 1)
			val := fmt.Sprintf("%s", ev.Kv.Value)

			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			if ev.Type == clientv3.EventTypePut {
				// Put operate
				config[key] = val
			} else if ev.Type == clientv3.EventTypeDelete {
				// Delete operate
				delete(config, key)
			}
		}
	}
}

/*
	removeDir remove special level dir

	@key the key will be operated
	@level how many level dir will be removed
*/
func removeDir(key string, level int) string {
	keySnippet := strings.Split(key, "/")
	return strings.Join(keySnippet[level:], "/")
}
