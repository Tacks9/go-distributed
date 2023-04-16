package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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
