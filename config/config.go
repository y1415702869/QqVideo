package config

import "gopkg.in/ini.v1"

var (
	QqVideoCookie   string
	VuSessionCookie string
	EmailHost       string
	EmailPort       int
	EmailUsername   string
	EmailPass       string
	NotifyEmail     string
)

func InitConfig() error {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		return err
	}
	loadEmailConfig(file)
	loadCookieConfig(file)

	return nil
}

func loadEmailConfig(file *ini.File) { //加载email服务配置
	EmailHost = file.Section("email").Key("Host").MustString("")
	EmailUsername = file.Section("email").Key("Username").MustString("")
	EmailPort = file.Section("email").Key("Port").MustInt(465)
	EmailPass = file.Section("email").Key("Pass").MustString("")
	NotifyEmail = file.Section("email").Key("NotifyEmail").MustString("")
}

func loadCookieConfig(file *ini.File) { //加载qq video cookie
	QqVideoCookie = file.Section("cookie").Key("QqVideoCookie").MustString("")
	VuSessionCookie = file.Section("cookie").Key("VuSessionCookie").MustString("")
}
