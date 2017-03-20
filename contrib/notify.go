package contrib

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v5"
)

var redisCli *redis.Client
var keyPrefix = "ngrok"

func InitNotify(notifyUrl string) error {
	if notifyUrl == "" {
		log.Fatalln("redisAddr is empty")
	}
	redisCli = redis.NewClient(&redis.Options{
		Addr:     notifyUrl,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := redisCli.Ping().Result()
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Errorln("redis ping failed!")
		return err
	}
	log.WithFields(log.Fields{"pong": pong}).Infoln("init Notify,ping redis success!")
	return nil
}

func AddMember(domain string, member string) error {
	cmd := redisCli.SAdd(fmt.Sprintf("%s.%s", keyPrefix, domain), member)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	log.WithFields(log.Fields{"key": fmt.Sprintf("%s.%s", keyPrefix, domain), "member": member}).Debugln("redis add pub url success!")
	return nil
}

func RemoveMember(domain string, member string) error {
	cmd := redisCli.SRem(fmt.Sprintf("%s.%s", keyPrefix, domain), member)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	log.WithFields(log.Fields{"key": fmt.Sprintf("%s.%s", keyPrefix, domain), "member": member}).Debugln("redis remove pub url success!")
	return nil
}
