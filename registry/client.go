package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// 客户端调用服务注册服务

// 向注册中心注册服务
func RegistryService(r Registration) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(r)
	if err != nil {
		return err
	}

	res, err := http.Post(ServicesURL, "application/json", buf)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to register service. Registry Service Response Code:%d", res.StatusCode)
	}

	return nil
}

// 向注册中心移除服务
func ShutdownService(url string) error {
	// 创建一个 HTTP 请求
	req, err := http.NewRequest(http.MethodDelete, ServicesURL, bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to deregister service. Registry Service Response Code:%d", res.StatusCode)
	}

	return nil
}
