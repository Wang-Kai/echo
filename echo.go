package echo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/dchest/uniuri"
)

var cli *clientv3.Client

func Init(endpoints []string) (*Echo, error) {
	var err error

	// connect etcd
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	// gen an nonce data
	var nonce = uniuri.New()
	return &Echo{Nonce: nonce}, nil
}

type Echo struct {
	Nonce string
	Map   map[string]string
}

// export an API to accept a map as dynamic config
func (e *Echo) Trusteeship(kvMap map[string]string) error {
	// put k/v into etcd
	for key, val := range kvMap {
		if strings.HasPrefix(key, "_") {
			// global params
			key = fmt.Sprintf("echo/%s", key)
		} else {
			// private params
			key = fmt.Sprintf("echo/%s/%s", e.Nonce, key)
		}

		_, err := cli.Put(context.Background(), key, val)
		if err != nil {
			return err
		}
	}

	e.Map = kvMap
	// start to watch
	go e.watch()
	return nil
}

// Destroy Trusteeship , remove all k/v in map from etcd
func (e *Echo) Destroy() error {
	defer cli.Close()

	// delete all private params
	var key = fmt.Sprintf("echo/%s", e.Nonce)

	log.Printf("Delete ====>  %s", key)
	_, err := cli.Delete(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for k := range e.Map {
		delete(e.Map, k)
	}
	return nil
}

func (e *Echo) watch() {
	rch := cli.Watch(context.Background(), "echo", clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			var (
				key = fmt.Sprintf("%s", ev.Kv.Key)
				val = fmt.Sprintf("%s", ev.Kv.Value)
			)
			// check if this action is concern me
			keys := strings.Split(key, "/")
			if len(keys) == 3 && keys[1] != e.Nonce {
				// this private params is not concern me
				continue
			}

			mapKey := keys[len(keys)-1]

			switch ev.Type {
			case mvccpb.PUT:
				e.Map[mapKey] = val
			case mvccpb.DELETE:
				delete(e.Map, mapKey)
			}
		}
	}
}
