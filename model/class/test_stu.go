// Package class 自动生成模板Student
package class

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/datatypes"
)

// Student 学生 结构体  Student
type Student struct {
	global.GVA_MODEL
	Name              *string        `json:"name" form:"name" gorm:"column:name;comment:;" binding:"required"`                                                 //姓名
	Age               *int           `json:"age" form:"age" gorm:"column:age;comment:;" binding:"required"`                                                    //年龄
	Sex               *bool          `json:"sex" form:"sex" gorm:"column:sex;comment:;"`                                                                       //性别
	BiographicalNotes datatypes.JSON `json:"biographicalNotes" form:"biographicalNotes" gorm:"column:biographical_notes;comment:;" swaggertype:"array,object"` //简历
	Pic               string         `json:"pic" form:"pic" gorm:"column:pic;comment:;"`                                                                       //头像
	Group             *string        `json:"group" form:"group" gorm:"column:group;comment:;"`                                                                 //组别
}

// TableName 学生 Student自定义表名 student
func (Student) TableName() string {
	return "student"
}
