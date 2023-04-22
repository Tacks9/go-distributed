package log

import (
	"io/ioutil"
	stlog "log"
	"net/http"
	"os"
)

// 标准库 log
var logger *stlog.Logger

// 自定义日志类型
type fileLog string

// 自定义日志类型
// 实现 io.Writer 接口中的 Write() 方法 将日志写入文件中
func (fl fileLog) Write(data []byte) (int, error) {
	fd, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	return fd.Write(data)
}

// 初始化日志记录器
// 将日志写入指定的位置，并包含日期和时间
func Run(destination string) {
	// 日志记录器
	logger = stlog.New(fileLog(destination), "[GO] - ", stlog.LstdFlags)
}

// 注册 日志服务
func RegisterHandlers() {
	// 处理 /log POST 请求
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			msg, err := ioutil.ReadAll(r.Body)
			if err != nil || len(msg) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// 记录日志
			write(string(msg))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

// 日志记录
func write(message string) {
	// 日志信息写入目标位置
	logger.Printf("%v \n", message)
}
