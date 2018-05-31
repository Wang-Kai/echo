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

	config, err := echo.GetConf("hi")
	if err != nil {
		t.Fatal(err)
	}

	ticker := time.NewTicker(time.Second * 2)
	for range ticker.C {
		log.Printf("%+v", config)
	}
}

func TestRemoveDir(t *testing.T) {
	key := removeDirPrefix("echo/klQn3lL/IpjekYZI4/age/me")
	log.Println(key)
}
