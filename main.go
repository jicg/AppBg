package main

import (
	"github.com/gin-gonic/gin"

	"github.com/jicg/AppBg/router"
	"net/http"
	"github.com/lestrrat/go-file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"github.com/rifflock/lfshook"
	"github.com/jicg/AppBg/middleware/wxpay"
	"time"
	"path"
	_ "github.com/jicg/AppBg/conf"
	"github.com/jicg/AppBg/conf"
)

func main() {
	if conf.GetConf() == nil {
		return
	}
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	ConfigLocalFilesystemLogger(conf.GetConf().Logpath, conf.GetConf().Logpath, 6000*time.Second, 6000*time.Second)
	r.Static("/page", "page")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/page")
	})
	router.Router(r)
	http.HandleFunc("/wxpay", wxpay.WxpayCallback)
	http.Handle("/", r)
	http.ListenAndServe(":"+conf.GetConf().Port, nil)
}

func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	baseLogPaht := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPaht+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPaht),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", err)
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	})
	log.AddHook(lfHook)
}
