package main

import (
	"QqVideo/config"
	"QqVideo/engine"
	"github.com/robfig/cron/v3"
	"log"
	"sync"
	"time"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalln("配置文件加载失败,err:", err.Error())
	}
	e := &engine.Engine{}
	var wg sync.WaitGroup

	go func(e *engine.Engine) {
		var cstSh, _ = time.LoadLocation("Asia/Shanghai")
		c := cron.New(cron.WithSeconds(), cron.WithLocation(cstSh))
		_, err := c.AddFunc("0 10 10 * * *", func() { //每天10:10:00执行
			e.Run(config.QqVideoCookie) //开启签到处理
		})

		if err != nil {
			log.Fatalln("定时器启动失败,err:", err.Error())
		}

		c.Start()

		select {}
	}(e)
	log.Println("服务已启动...")
	wg.Add(1)

	wg.Wait()
}
