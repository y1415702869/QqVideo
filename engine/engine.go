package engine

import (
	"QqVideo/config"
	"QqVideo/email"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

const (
	LoginUrl = "https://access.video.qq.com/user/auth_refresh?vappid=11059694&vsecret=" +
		"fdf61a6be0aad57132bc5cdf78ac30145b6cd2c1470b0cfe&type=qq&g_tk=&g_vstk=649650180&g_actk=1299537010&callback=" +
		"jQuery19102408648260124442_1648355235495&_=1648355235496"
	ReqUrl      = "https://vip.video.qq.com/fcgi-bin/comm_cgi?name=hierarchical_task_system&cmd=2"
	NotifyEmail = "1784605674@qq.com"
)

var JsonReg = regexp.MustCompile(`QZOutputJson=\((.*)\);`)

type Engine struct {
	videoCookie string
}

// Run 启动签到
func (e *Engine) Run(cookie string) {
	e.videoCookie = cookie
	vuSession, err := e.getVqqVuSession()
	if err != nil {
		log.Println(err)
		email.SendEmail(NotifyEmail, "登陆cookie失效", err.Error())
		return
	}
	bytes, err := e.httpReq(vuSession)
	if err != nil {
		log.Println(err)
		email.SendEmail(NotifyEmail, "签到请求失败", err.Error())
		return
	}
	score, err := e.withRes(&bytes)
	if err != nil {
		log.Println(err)
		email.SendEmail(NotifyEmail, "签到请求结果处理失败", err.Error())
		return
	}
	email.SendEmail(NotifyEmail,
		"签到成功",
		fmt.Sprintf("获得积分:%d<br/>签到时间:%s",
			score, time.Now().Format("2006-01-02 15:04:05")))
	log.Println("签到成功,获得积分:", score)
}

//处理结果
func (e *Engine) withRes(res *[]byte) (int, error) {
	var score int
	RegSub := JsonReg.FindStringSubmatch(string(*res))
	if len(RegSub) != 2 {
		return score, errors.New(fmt.Sprintf("返回数据解析失败,返回数据:%s", string(*res)))
	}

	resJson := e.formatJson([]byte(RegSub[1])) //转换数据格式
	if resJson.Ret != 0 {                      //error
		return score, errors.New(fmt.Sprintf("积分签到失败,errCode:%d,errMsg:%s", resJson.Ret, resJson.Msg))
	}
	score = resJson.CheckinScore
	return score, nil
}

//获取vqq vu session
func (e *Engine) getVqqVuSession() (string, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", LoginUrl, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("cookie", e.videoCookie)
	req.Header.Set("referer", "https://v.qq.com")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.83 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	for _, cookie := range res.Cookies() {
		if cookie.Name == "vqq_vusession" {
			return cookie.Value, nil
		}
	}
	return "", errors.New("登陆cookie失效")
}

//发起请求
func (e *Engine) httpReq(vuSession string) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", ReqUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("cookie", fmt.Sprintf(config.VuSessionCookie, vuSession))
	req.Header.Set("referer", "https://m.v.qq.com")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.83 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

type ResJson struct {
	Ret          int    `json:"ret"`
	CheckinScore int    `json:"checkin_score"`
	Msg          string `json:"msg"`
}

//转换返回json
func (e *Engine) formatJson(bytes []byte) *ResJson {
	var jsonStr ResJson
	json.Unmarshal(bytes, &jsonStr)
	return &jsonStr
}
