package utools

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"net"
	"os"
	"time"
)

var LocalIp string
var stmp *gomail.Dialer
var Pwd string
var HostName string

// var TimeZone, _ = time.LoadLocation("Asia/Shanghai")
func init() {
	Pwd, _ = os.Getwd()
	stmp = gomail.NewDialer("smtp.exmail.qq.com", 465, os.Getenv("email.name"), os.Getenv("email.pwd"))
	stmp.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	HostName, _ = os.Hostname()
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					LocalIp = ipnet.IP.String()
				}
			}
		}
	}
}

// RegisterTimeZone 注册时区 "Asia/Shanghai"
func RegisterTimeZone(zone string) {
	time.Local, _ = time.LoadLocation(zone)
}
