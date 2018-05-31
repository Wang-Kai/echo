package echo

import (
	"testing"
	"time"

	"log"
)

func TestGetConf(t *testing.T) {
	echo, err := New("http://localhost:2379")
	if err != nil {
		t.Fatal(err)
	}

	config, err := echo.GetConf("/")
	if err != nil {
		t.Fatal(err)
	}

	ticker := time.NewTicker(time.Second * 3)
	for range ticker.C {
		log.Printf("%+v", config)

		log.Println(config.Get("name"))
	}
}

func TestRemoveDir(t *testing.T) {
	key := removeDirPrefix("echo/klQn3lLIpjekYZI4/age", 1)
	log.Println(key)
}
