package logger

import (
	"io"
	"log"
	"os"
)

var Log *log.Logger

func init() {
	file, err := os.OpenFile("controller.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Unable to create log file:", err)
	}
	//defer file.Close()

	// 使用io.MultiWriter将日志输出同时写入控制台和文件
	logWriter := io.MultiWriter(os.Stdout, file)

	flag := log.Ldate | log.Ltime | log.Lmicroseconds // log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile
	Log = log.New(logWriter, "", flag)
}
