package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

//GetMD5Hash --
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

//GetIPFromNetWork --
func GetIPFromNetWork(name string) string {
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		if i.Name == name {
			addrs, _ := i.Addrs()
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						os.Stdout.WriteString(ipnet.IP.String() + "\n")
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return ""
}

func getFullPath(path string) string {
	absolutePath, _ := filepath.Abs(path)
	return absolutePath
}

//PrintFilesName --
func PrintFilesName(path, ip string) {
	if path == "" || ip == "" {
		return
	}
	fullPath := getFullPath(path)
	filepath.Walk(fullPath, func(path string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		name := fi.Name()
		if strings.Contains(name, "json") {
			b, e := ioutil.ReadFile(path)
			if e != nil {
				fmt.Println("read file error")
			}
			res := string(b)
			if strings.Contains(string(b), "${BIND_IP}") {
				os.Remove(path)
				res = strings.Replace(res, "${BIND_IP}", ip, -1)
				f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
				if err != nil {
					log.Error(err.Error())
					log.Error("consul.d replace error")
				}
				defer f.Close()
				_, e := io.WriteString(f, res)
				if e != nil {
					log.Error(e.Error())
				}
			}
		}
		return nil
	})
}
