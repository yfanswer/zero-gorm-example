package main

import (
	"database/sql"
	"flag"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"log"
	ycache "zero-gorm-example/model/cache"
)

var configFile = flag.String("f", "etc/practice.yaml", "the config file")

func main() {
	flag.Parse()

	if err := LoadConfig(*configFile); err != nil {
		log.Fatalf("%+v", err)
	}
	log.Println("configuration:", *Configuration())

	if err := NewGDB(Configuration().DSN); err != nil {
		log.Fatalf("%+v", err)
	}
	um := ycache.NewUserModel(GDB(), cache.CacheConf{
		{
			RedisConf: redis.RedisConf{
				Host: "127.0.0.1:6379",
				//Type: "node",
				//Pass: "",
				//Tls: false,
			},
			Weight: 100,
		}})
	err := um.Insert(&ycache.User{
		User: "aaa2",
		Name: sql.NullString{
			String: "yfaaa2",
			Valid:  true,
		},
		Password: "123456",
		Mobile:   "13243204942",
		Gender:   "男",
		Nickname: "yf",
		Tp:       2,
	})
	if err != nil {
		log.Printf("%+v", err)
		return
	}
	u, err := um.FindOne(1)
	if err != nil {
		log.Printf("%+v", err)
		return
	}
	log.Printf("u:%#v", *u)
}
