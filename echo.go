package echo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
)

// Config a k/v string map to save setting info
type Config map[string]string

// Get return value in config for `k`
func (c Config) Get(k string) (val string, exist bool) {
	val, exist = c[k]

	return val, exist
}

// Echo the instance of echo lib
type Echo struct {
	Configs map[string]Config
}

var cli = new(clientv3.Client)

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
	echo := &Echo{
		Configs: make(map[string]Config),
	}
	return echo, nil
}

/*
	GetConf get all k/v from etcd prefix with special etcdDir, and save k/v in map

	@etcdDir: the dir of config keys prefix with
*/
func (e *Echo) GetConf(etcdDir string) (Config, error) {
	config, ok := e.Configs[etcdDir]
	if ok {
		return config, nil
	}

	resp, err := cli.Get(context.Background(), etcdDir, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	config = make(Config)

	for _, ev := range resp.Kvs {
		key, val := fmt.Sprintf("%s", ev.Key), fmt.Sprintf("%s", ev.Value)
		key = removeDirPrefix(key)

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
			key = removeDirPrefix(key)
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
	removeDir remove dir prefix, and return key name

	@key the key will be operated
	@level how many level dir will be removed
*/
func removeDirPrefix(key string) string {
	dirSplts := strings.Split(key, "/")

	return dirSplts[len(dirSplts)-1]
}
