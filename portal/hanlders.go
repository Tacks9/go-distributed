package portal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Tacks9/go-distributed/grades"
	"github.com/Tacks9/go-distributed/registry"
)

func RegisterHandlers() {
	// 重定向到 /students
	http.Handle("/", http.RedirectHandler("/students", http.StatusPermanentRedirect))

	h := new(studentsHandler)
	http.Handle("/students", h)
	http.Handle("/students/", h)
}

type studentsHandler struct{}

// 提供 HTTP 服务
func (sh studentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")
	switch len(pathSegments) {
	case 2: // /students
		sh.renderStudents(w, r)
	case 3: // /students/{:id}
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.renderStudent(w, r, id)
	case 4: // /students/{:id}/grades
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if strings.ToLower(pathSegments[3]) != "grades" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.renderGrades(w, r, id)

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// 渲染模版 学生列表
func (studentsHandler) renderStudents(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error retrieving students: ", err)
		}
	}()

	serviceURL, err := registry.GetProvider(registry.GradingService)
	if err != nil {
		return
	}

	res, err := http.Get(serviceURL + "/students")
	if err != nil {
		return
	}

	var s grades.Students
	err = json.NewDecoder(res.Body).Decode(&s)
	if err != nil {
		return
	}

	rootTemplate.Lookup("students.html").Execute(w, s)
}

// 渲染模版 学生页面
func (studentsHandler) renderStudent(w http.ResponseWriter, r *http.Request, id int) {

	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error retrieving students: ", err)
			return
		}
	}()

	// 获取服务地址
	serviceURL, err := registry.GetProvider(registry.GradingService)
	if err != nil {
		return
	}

	// 模拟请求具体的学生
	res, err := http.Get(fmt.Sprintf("%v/students/%v", serviceURL, id))
	if err != nil {
		return
	}

	// 解析到 s
	var s grades.Student
	err = json.NewDecoder(res.Body).Decode(&s)
	if err != nil {
		return
	}

	// 输出模版视图
	rootTemplate.Lookup("student.html").Execute(w, s)
}

// 添加学生成绩
func (studentsHandler) renderGrades(w http.ResponseWriter, r *http.Request, id int) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer func() {
		w.Header().Add("location", fmt.Sprintf("/students/%v", id))
		w.WriteHeader(http.StatusTemporaryRedirect)
	}()
	title := r.FormValue("Title")
	gradeType := r.FormValue("Type")
	score, err := strconv.ParseFloat(r.FormValue("Score"), 32)
	if err != nil {
		log.Println("Failed to parse score: ", err)
		return
	}
	// 赋值成绩
	g := grades.Grade{
		Title: title,
		Type:  grades.GradeType(gradeType),
		Score: float32(score),
	}
	data, err := json.Marshal(g)
	if err != nil {
		log.Println("Failed to convert grade to JSON: ", g, err)
	}

	serviceURL, err := registry.GetProvider(registry.GradingService)
	if err != nil {
		log.Println("Failed to retrieve instance of Grading Service", err)
		return
	}
	// 保存学生成绩
	res, err := http.Post(fmt.Sprintf("%v/students/%v/grades", serviceURL, id), "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("Failed to save grade to Grading Service", err)
		return
	}
	if res.StatusCode != http.StatusCreated {
		log.Println("Failed to save grade to Grading Service. Status: ", res.StatusCode)
		return
	}
}
