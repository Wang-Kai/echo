package echo

import (
	"testing"
	"time"

	"github.com/goushuyun/log"
)

func TestGetConf(t *testing.T) {
	echo, err := New("http://localhost:2379")
	if err != nil {
		t.Fatal(err)
	}

	config, err := echo.GetConf("uservice_app/")
	if err != nil {
		t.Fatal(err)
	}

	log.JSON(config)

	ticker := time.NewTicker(time.Second * 5)

	for range ticker.C {
		log.JSONIndent(config)
	}
}

func TestRemoveDir(t *testing.T) {
	key := removeDir("echo/klQn3lLIpjekYZI4/age", 1)
	t.Log(key)
}
