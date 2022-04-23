package debughelper

import (
	"fmt"
	"log"
	"time"
)

var Debug = false

func ShowError(err error) {
	if !Debug {
		return
	}
	log.Println(err)
}

func ShowInfo(info string) {
	if !Debug {
		return
	}
	fmt.Println(time.Now().GoString() + ":" + info)
}

func ShowBuf(buf []byte) {
	if !Debug {
		return
	}
	fmt.Println(time.Now().GoString() + ":" + string(buf))
}
