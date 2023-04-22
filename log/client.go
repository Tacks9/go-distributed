package log

import (
	"bytes"
	"fmt"
	stlog "log"
	"net/http"

	"github.com/Tacks9/go-distributed/registry"
)

type clientLogger struct {
	url string
}

func SetClientLogger(ServiceURL string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	stlog.SetFlags(0)
	// 设置日志的输出
	stlog.SetOutput(&clientLogger{url: ServiceURL})
}

// 实现 IO Write 接口
func (cl clientLogger) Write(data []byte) (int, error) {
	b := bytes.NewBuffer([]byte(data))
	res, err := http.Post(cl.url+"/log", "text/plain", b)
	if err != nil {
		return 0, err
	}

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to send log message. Service responded with %d - %s", res.StatusCode, res.Status)
	}

	return len(data), nil
}
