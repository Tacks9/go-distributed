package grades

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// 内部使用
type studentsHandler struct{}

func RegisterHandlers() {
	handler := new(studentsHandler)
	// 集合资源
	http.Handle("/students", handler)
	// 单个资源
	http.Handle("/students/", handler)
}

// 处理 HTTP 请求，支持下面三种场景
// /students
// /students/{id}
// /students/{id}/grades
func (sh studentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathInfo := strings.Split(r.URL.Path, "/")
	switch len(pathInfo) {
	case 2:
		// students
		sh.getAll(w, r)
	case 3:
		// students/id
		// 获取第二个参数
		id, err := strconv.Atoi(pathInfo[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.getOne(w, r, id)
	case 4:
		// students/id/grades
		id, err := strconv.Atoi(pathInfo[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// 新增
		sh.addGrade(w, r, id)
	default:
		w.WriteHeader(http.StatusNotFound)

	}
}

// 获取权量数据
func (sh studentsHandler) getAll(w http.ResponseWriter, r *http.Request) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	data, err := sh.toJSON(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

// 获取权量数据
func (sh studentsHandler) getOne(w http.ResponseWriter, r *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	// 获取单个数据
	student, err := students.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}
	data, err := sh.toJSON(student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to student:%q", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

// 录入成绩
func (sh studentsHandler) addGrade(w http.ResponseWriter, r *http.Request, id int) {
	// 加锁
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	// 读取数据
	student, err := students.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}
	// 解析 body
	var g Grade
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&g)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	// 新增对话窗口
	student.Grades = append(student.Grades, g)
	w.WriteHeader(http.StatusCreated)
	data, err := sh.toJSON(g)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "applicaiton/json")
	w.Write(data)
}

// 转化 JSON
func (sh studentsHandler) toJSON(obj interface{}) ([]byte, error) {
	fmt.Println(obj)
	var b bytes.Buffer
	encode := json.NewEncoder(&b)
	err := encode.Encode(obj)
	if err != nil {
		return b.Bytes(), fmt.Errorf("Failed to serialize students:%q", err)
	}
	return b.Bytes(), nil
}
