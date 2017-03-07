package contrib

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v5"
)

var redisCli *redis.Client

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
	cmd := redisCli.SAdd(domain, member)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func RemoveMember(domain string, member string) error {
	cmd := redisCli.SRem(domain, member)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}
