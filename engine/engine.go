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
	"strings"
	"time"
)

const (
	LoginUrl = "https://access.video.qq.com/user/auth_refresh?vappid=11059694&vsecret=" +
		"fdf61a6be0aad57132bc5cdf78ac30145b6cd2c1470b0cfe&type=qq&g_tk=&g_vstk=649650180&g_actk=1299537010&callback=" +
		"jQuery19102408648260124442_1648355235495&_=1648355235496"
	ReqUrl = "https://vip.video.qq.com/fcgi-bin/comm_cgi?name=hierarchical_task_system&cmd=2"
)

var (
	JsonReg                = regexp.MustCompile(`QZOutputJson=\((.*)\);`)
	FindCookieVuSessionReg = regexp.MustCompile(`vqq_vusession=([^;]*){1}`)
)

type Engine struct {
	videoCookie     string //login cookie
	vuSessionCookie string //签到cookie
}

// Run 启动签到
func (e *Engine) Run(cookie string) {
	e.videoCookie = cookie
	if err := e.getVuSessionCookie(); err != nil {
		log.Println(err)
		email.SendEmail(config.NotifyEmail, "cookie解析出错", err.Error())
		return
	}
	bytes, err := e.httpRequest(ReqUrl, e.vuSessionCookie, "https://m.v.qq.com", false)
	if err != nil {
		log.Println(err)
		email.SendEmail(config.NotifyEmail, "获取V力值请求失败", err.Error())
		return
	}
	score, err := e.withRes(&bytes)
	if err != nil {
		log.Println(err)
		email.SendEmail(config.NotifyEmail, "签到请求结果处理失败", err.Error())
		return
	}
	email.SendEmail(config.NotifyEmail,
		"签到成功",
		fmt.Sprintf("获得积分:%d<br/>签到时间:%s",
			score, time.Now().Local().Format("2006-01-02 15:04:05")))
	log.Println("签到成功,获得积分:", score)
}

//getVuSessionCookie get VuSession
func (e *Engine) getVuSessionCookie() error {
	findVuSession := FindCookieVuSessionReg.FindStringSubmatch(e.videoCookie)
	if len(findVuSession) != 2 {
		return errors.New("cookie错误")
	}
	//替换成sprint string
	vuSessionCookieSprint := strings.Replace(e.videoCookie, findVuSession[1], "%s", 1)
	vuSession, err := e.httpRequest(LoginUrl, e.videoCookie, "https://v.qq.com", true)
	if err != nil {
		return err
	}
	//获取v力值cookie获取成功
	e.vuSessionCookie = fmt.Sprintf(vuSessionCookieSprint, string(vuSession))
	return nil
}

//处理V力值签到结果
func (e *Engine) withRes(res *[]byte) (int, error) {
	RegSub := JsonReg.FindStringSubmatch(string(*res))
	if len(RegSub) != 2 {
		return 0, errors.New(fmt.Sprintf("返回数据解析失败,返回数据:%s", string(*res)))
	}

	resJson := e.formatJson([]byte(RegSub[1])) //转换数据格式
	if resJson.Ret != 0 {                      //error
		return 0, errors.New(fmt.Sprintf("积分签到失败,errCode:%d,errMsg:%s", resJson.Ret, resJson.Msg))
	}
	return resJson.CheckinScore, nil
}

//httpRequest network request
func (e *Engine) httpRequest(url, cookieStr, referer string, isGetVuSession bool) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("cookie", cookieStr)
	req.Header.Set("referer", referer)
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.83 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if isGetVuSession { //get vqq_vusession
		for _, cookie := range res.Cookies() {
			if cookie.Name == "vqq_vusession" {
				return []byte(cookie.Value), nil
			}
		}
		return nil, errors.New("登陆cookie失效")
	}
	//V力值签到
	return ioutil.ReadAll(res.Body)
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
