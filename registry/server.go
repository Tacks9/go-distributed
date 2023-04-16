package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// 移除一个服务
func (r *registry) remove(url string) error {
	for i := range reg.registrations {
		if reg.registrations[i].ServiceURL == url {
			r.mutex.Lock()
			// 移除当前 item
			reg.registrations = append(reg.registrations[:i], reg.registrations[i+1:]...)
			r.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("Service at URL %s not found", url)
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
	case http.MethodDelete:
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(payload)
		log.Printf("Removing service at URL:%s \n", url)

		// 移除服务
		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
