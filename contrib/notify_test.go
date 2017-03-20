package contrib

import (
	"fmt"
	"testing"
	"time"
)

func TestMemberAddAndRemove(t *testing.T) {
	err := InitNotify("192.168.1.26:6379")
	if err != nil {
		t.Errorf("redis init failed!err:=%v", err)
		return
	}
	err = AddMember("3.tunnel.daocloud.co", "http://1.2.tunnel.daocloud.co:54333")
	if err != nil {
		t.Errorf("redis add failed!err:=%v", err)
		return
	}
	time.Sleep(time.Second)
	str, err := redisCli.SMembers("ngrok.3.tunnel.daocloud.co").Result()
	if err != nil {
		t.Errorf("redis smember failed!err:=%v", err)
		return
	}
	fmt.Println(str)
	if len(str) == 0 {
		t.Errorf("add member to redis failed!smember return nil")
		return
	}
	if str[0] != "http://1.2.tunnel.daocloud.co:54333" {
		t.Errorf("add member to redis failed!str(%s) not match", str[0])
		return
	}
}
