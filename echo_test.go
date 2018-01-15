package echo

import (
	"testing"
	"time"

	"github.com/goushuyun/log"
)

func TestInit(t *testing.T) {
	Init([]string{"localhost:2379"})
}

func TestEcho(t *testing.T) {
	var setting = map[string]string{
		"_name": "Wang",
		"age":   "26",
	}

	echo, err := Init([]string{"localhost:2379"})
	if err != nil {
		t.Fatal(err)
	}

	err = echo.Trusteeship(setting)
	if err != nil {
		t.Fatal(err)
	}

	// test Destory
	go func(e *Echo) {
		timer := time.NewTimer(time.Second * 30)
		<-timer.C
		e.Destroy()
	}(echo)

	ticker := time.NewTicker(3 * time.Second)
	for {
		log.JSON(setting)
		<-ticker.C
	}
}
