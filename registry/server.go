package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// 注册服务-服务端监听端口
const ServicePort = ":3000"

// 注册服务对外查询接口
const ServicesURL = "http://localhost" + ServicePort + "/services"

// 服务注册结构体
type registry struct {
	registrations []Registration
	// 保证 registrations 并发读写安全
	mutex *sync.RWMutex
}

type RegistrationService struct{}

// 全局变量 初始化
var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.RWMutex),
}

var once sync.Once

// 启动心跳检测
func SetupRegistryService() {
	// 只会在执行一次
	once.Do(func() {
		go reg.heartbeat(3 * time.Second)
	})
}

// 心跳检查
func (r *registry) heartbeat(freq time.Duration) {
	// 循环
	for {
		var wg sync.WaitGroup
		// 遍历所有的注册服务
		for _, regItem := range r.registrations {

			// 并发请求
			wg.Add(1)
			go func(reg Registration) {
				defer wg.Done()
				success := true
				// 请求3次重试
				for attemps := 0; attemps < 3; attemps++ {
					res, err := http.Get(reg.HeartbeatURL)
					if err != nil {
						log.Println(err)
					} else if res.StatusCode == http.StatusOK {
						log.Printf("Heartbeat check passed for %v", reg.ServiceName)
						// 检测成功
						if !success {
							// 进行添加
							r.add(reg)
						}
						break
					}

					// 检测失败
					log.Printf("Heartbeat check failed for %v", reg.ServiceName)
					if success {
						// 进行移除
						success = false
						r.remove(reg.ServiceURL)
					}
					time.Sleep(1 * time.Second)
				}
			}(regItem)
			// 等待所有 go routinue
			wg.Wait()

			// 按照一定频次进行 心跳检测
			time.Sleep(freq)
		}
	}
}

// 注册一个服务
func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()

	err := r.sendRequiredServices(reg)
	if err != nil {
		return err
	}

	// 服务出现的时候通知
	r.notify(patch{
		Added: []patchEntry{
			{
				Name: reg.ServiceName,
				URL:  reg.ServiceURL,
			},
		},
	})
	return nil
}

// 移除一个服务
func (r *registry) remove(url string) error {
	for i := range reg.registrations {
		if reg.registrations[i].ServiceURL == url {
			// 需要移除的服务通知
			r.notify(patch{
				Removed: []patchEntry{
					{
						Name: r.registrations[i].ServiceName,
						URL:  r.registrations[i].ServiceURL,
					},
				},
			})
			r.mutex.Lock()
			// 移除当前 item
			reg.registrations = append(reg.registrations[:i], reg.registrations[i+1:]...)
			r.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("Service at URL %s not found", url)
}

// 在服务中心发现服务
func (r registry) sendRequiredServices(reg Registration) error {
	// 读取当前服务注册中心的服务，用读锁即可
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var p patch
	for _, serviceReg := range r.registrations {
		// 再遍历依赖的
		for _, reqService := range reg.RequiredServices {
			if serviceReg.ServiceName == reqService {
				// 发现依赖的
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.ServiceName,
					URL:  serviceReg.ServiceURL,
				})
			}
		}

	}

	// 向服务发送
	err := r.sendPatch(p, reg.ServiceUpdateURL)
	if err != nil {
		return err
	}
	return nil
}

// 告知依赖项
func (r registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		return nil
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return err
	}
	return nil
}

// 通知依赖变更
func (r registry) notify(fullPatch patch) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 注册中心所有服务
	for _, reg := range r.registrations {
		// 并发通知
		go func(reg Registration) {
			// 每个注册的服务，对其依赖服务项进行处理
			for _, reqService := range reg.RequiredServices {
				// 初始化
				p := patch{Added: []patchEntry{}, Removed: []patchEntry{}}
				// 是否需要更新
				sendUpdate := false

				// 添加服务
				for _, added := range fullPatch.Added {
					// 如果添加的服务，正好是某个服务的依赖项
					if added.Name == reqService {
						p.Added = append(p.Added, added)
						sendUpdate = true
					}

				}
				// 移除服务
				for _, removed := range fullPatch.Removed {
					if removed.Name == reqService {
						p.Removed = append(p.Removed, removed)
						sendUpdate = true
					}
				}
				if sendUpdate {
					err := r.sendPatch(p, reg.ServiceUpdateURL)
					if err != nil {
						log.Println(err)
						return
					}
				}

			}
		}(reg)

	}
}

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
