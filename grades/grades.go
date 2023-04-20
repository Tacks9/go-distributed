package grades

import (
	"fmt"
	"sync"
)

// 学生类型
type Student struct {
	ID        int
	FirstName string
	LastName  string
	Grades    []Grade
}

type Students []Student

// 全局变量
var students Students

// 互斥锁
var studentsMutex sync.Mutex

// 根据ID寻找学生信息
func (ss Students) GetByID(id int) (*Student, error) {
	for i := range ss {
		if ss[i].ID == id {
			return &ss[i], nil
		}
	}
	return nil, fmt.Errorf("Student with ID:%d not found", id)
}

// 成绩类型
type GradeType string

const (
	GradeQuiz = GradeType("Quiz")
	GradeTest = GradeType("Test")
	GradeExam = GradeType("Exam")
)

// 成绩
type Grade struct {
	Title string
	Type  GradeType
	Score float32
}

// 计算平均分
func (s Student) Average() float32 {
	var result float32
	for _, grade := range s.Grades {
		result += grade.Score
	}
	return result / float32(len(s.Grades))
}
