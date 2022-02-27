package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	xlog "ipv6_share/log"
	"net"
	"os"
	"strings"
)

type appObj struct {
	maxFileNum  int
	maxDepth    int
	shareFiles  []*shareFileInfo
	currentPath string
	version     string
	publicIpv6s []net.IP
}

var app *appObj

func main() {
	app = newAppObj()

	pathStr := flag.String("p", app.currentPath, "set share paths, use ; to separate, e.g. -p='/Users/src/my_photos;/Users/src/my_videos'")
	flag.Parse()
	if pathStr == nil {
		flag.Usage()
		return
	}
	fmt.Println(*pathStr)

	app.setFileList(strings.Split(*pathStr, ";"))
	xlog.Warn("app info", app)

	//if len(app.publicIpv6s) == 0 {
	//	xlog.Warn("no global unicast ipv6 address!!!")
	//	return
	//}
	//port := 1228
	//addr := fmt.Sprintf("[%s]:%d", app.publicIpv6s[0], port)

	addr := "127.0.0.1:1228"
	engine := gin.Default()
	router(engine)
	err := engine.Run(addr)
	if err != nil {
		xlog.Error("listen failed", addr, err)
		return
	}
}

func newAppObj() *appObj {
	newApp := &appObj{
		maxFileNum:  1000,
		maxDepth:    5,
		currentPath: getPwd(),
		version:     Version,
		publicIpv6s: getPublicIpv6(),
	}
	return newApp
}

func getPublicIpv6() []net.IP {
	ifaces, err := net.Interfaces()
	if err != nil {
		xlog.Error("Interfaces", err)
		return nil
	}

	ipv6Addrs := make([]net.IP, 0, 10)
	for _, i := range ifaces {
		addrs, tmpErr := i.Addrs()
		if tmpErr != nil {
			xlog.Error("Addrs", tmpErr)
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				xlog.Info("IPNet", v)
			case *net.IPAddr:
				ip = v.IP
				xlog.Info("IPAddr", v)
			}
			if ip == nil || ip.To4() != nil || !ip.IsGlobalUnicast() {
				continue
			}
			ipv6Addrs = append(ipv6Addrs, ip)
		}
	}

	return ipv6Addrs
}

func getPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		xlog.Error("Getwd", err)
		return ""
	}
	return dir
}

func (a *appObj) setFileList(paths []string) {
	a.shareFiles = make([]*shareFileInfo, 0, a.maxFileNum)
	for _, path := range paths {
		traverseDir(path, a.maxDepth, a.maxFileNum, a)
	}
	for i, file := range a.shareFiles {
		xlog.Info("file", i, file.File.Name())
	}
}

func traverseDir(path string, depth, maxFileNum int, app *appObj) {
	if depth < 0 {
		return
	}
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		xlog.Error("ReadDir", path, err)
		return
	}

	xlog.Info("find fileInfos", path)
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			xlog.Info("is dir", fileInfo.Name())
			traverseDir(path+string(os.PathSeparator)+fileInfo.Name(), depth-1, maxFileNum, app)
			continue
		}
		if len(app.shareFiles) >= maxFileNum {
			return
		}

		f, err := os.Open(path + string(os.PathSeparator) + fileInfo.Name())
		if err != nil {
			xlog.Error("open fileInfo err", err)
			continue
		}
		fi, err := newFileInfo(f, fileInfo)
		if err != nil {
			xlog.Error("newFileInfo err", err)
			continue
		}
		app.shareFiles = append(app.shareFiles, fi)
	}
}
