package registry

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

// 注册服务-服务端监听端口
const ServicePort = ":3000"

// 注册服务对外查询接口
const ServicesURL = "http://localhost" + ServicePort + "/services"

// 服务注册结构体
type registry struct {
	registrations []Registration
	// 保证 registrations 并发读写安全
	mutex *sync.Mutex
}

// 注册一个服务
func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()

	return nil
}

// 全局变量 初始化
var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.Mutex),
}

type RegistrationService struct{}

// 注册服务-服务端
func (s RegistrationService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request receviced")

	switch r.Method {
	case http.MethodPost:
		var regItem Registration

		// json 解码
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&regItem)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service:%v with URL:%s \n", regItem.ServiceName,
			regItem.ServiceURL)

		// 注册服务
		err = reg.add(regItem)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
