package engine

import (
	"QqVideo/config"
	"github.com/robfig/cron/v3"
	"log"
	"sync"
	"time"
)

var (
	wg                sync.WaitGroup
	cstSh, _          = time.LoadLocation("Asia/Shanghai")
	SignTimeRule      = "0 10 10 * * *" //每天10:10:00执行 每日签到
	Minutes60TimeRule = "0 10 20 * * *" //每天20:10:00执行 观看60分钟
)

// GoTask run task
func GoTask() {
	e := &Engine{}
	time.Local = cstSh
	c := cron.New(cron.WithSeconds(), cron.WithLocation(cstSh))

	//every day sign
	go runTask(c, e, &Params{
		Cookie:        config.QqVideoCookie,
		ReqUrl:        SignUrl,
		EmailSubject:  "每日签到",
		NotifyMsg:     "获得V力值:%d",
		WithResErrMsg: "每日签到失败",
	}, SignTimeRule)

	//Minutes60
	go runTask(c, e, &Params{
		Cookie:        config.QqVideoCookie,
		ReqUrl:        Minutes60Url,
		EmailSubject:  "观看视频满60分钟",
		NotifyMsg:     "获得V力值:%d",
		WithResErrMsg: "看60分钟V力值获取失败",
	}, Minutes60TimeRule)

	log.Println("服务已启动...")

	wg.Add(1)
	wg.Wait()
}

func runTask(c *cron.Cron, e *Engine, params *Params, timeRule string) {
	_, err := c.AddFunc(timeRule, func() { //开启定时任务
		e.Run(params) //任务处理
	})
	if err != nil {
		log.Fatalln("定时器启动失败,err:", err.Error())
	}
	c.Start()
}
