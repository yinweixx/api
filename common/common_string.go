package common

import (
	"errors"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	applications = []string{"api", "file"}
)

//CheckURL 检查URL地址
func CheckURL(url string) (api string, version string, application string, err error) {
	urls := strings.Split(url, "/")
	n := 0
	for _, s := range urls {
		if strings.Contains(s, "v") {
			_s := strings.Replace(s, "v", "", 1)
			_, err := strconv.Atoi(_s)
			if err == nil {
				n++
				version = s
			} else {
				api += "/" + s
			}
		} else if s != "" {
			api += "/" + s
		}
	}
	for _, a := range applications {
		if strings.Contains(url, a) {
			application = a
		}
	}
	if n == 0 || n > 1 || application == "" {
		log.WithFields(log.Fields{
			"prefix": "common.checkurl",
		}).Error("illegal param")
		api = ""
		version = ""
		err = errors.New("illegar params")
	}
	log.WithFields(log.Fields{
		"prefix":  "common.checkurl",
		"api":     api,
		"version": version,
	}).Info("success to get url")
	return
}

// func NewUtil() *Util {
// 	return &Util{
// 		Version: "v1.2",
// 	}
// }

// type Util struct {
// 	Version string
// }

// func (_this *Util) method1() {
// 	if _this.Version == "v1.0.0.1" {

// 	} else {

// 	}
// }
