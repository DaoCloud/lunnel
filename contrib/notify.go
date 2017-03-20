package contrib

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v5"
)

var redisCli *redis.Client
var notifyKey string

func InitNotify(notifyUrl string, nk string) error {
	if notifyUrl == "" {
		log.Fatalln("notifyUrl is empty")
	}
	if nk == "" {
		log.Fatalln("notifyKey is empty")
	}
	notifyKey = nk
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
	cmd := redisCli.Del(notifyKey)
	if cmd.Err() != nil {
		log.WithFields(log.Fields{"err": cmd.Err()}).Errorln("init redis,delete key failed!")
		return cmd.Err()
	}
	return nil
}

func AddMember(domain string, member string) error {
	cmd := redisCli.SAdd(notifyKey, member)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	log.WithFields(log.Fields{"key": notifyKey, "member": member}).Debugln("redis add pub url success!")
	return nil
}

func RemoveMember(domain string, member string) error {
	cmd := redisCli.SRem(notifyKey, member)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	log.WithFields(log.Fields{"key": notifyKey, "member": member}).Debugln("redis remove pub url success!")
	return nil
}
