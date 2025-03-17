package class

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/class"
	classReq "github.com/flipped-aurora/gin-vue-admin/server/model/class/request"
)

type StudentService struct{}

// CreateStudent 创建学生记录
// Author [yourname](https://github.com/yourname)
func (stuService *StudentService) CreateStudent(stu *class.Student) (err error) {
	err = global.GVA_DB.Create(stu).Error
	return err
}

// DeleteStudent 删除学生记录
// Author [yourname](https://github.com/yourname)
func (stuService *StudentService) DeleteStudent(ID string) (err error) {
	err = global.GVA_DB.Delete(&class.Student{}, "id = ?", ID).Error
	return err
}

// DeleteStudentByIds 批量删除学生记录
// Author [yourname](https://github.com/yourname)
func (stuService *StudentService) DeleteStudentByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]class.Student{}, "id in ?", IDs).Error
	return err
}

// UpdateStudent 更新学生记录
// Author [yourname](https://github.com/yourname)
func (stuService *StudentService) UpdateStudent(stu class.Student) (err error) {
	err = global.GVA_DB.Model(&class.Student{}).Where("id = ?", stu.ID).Updates(&stu).Error
	return err
}

// GetStudent 根据ID获取学生记录
// Author [yourname](https://github.com/yourname)
func (stuService *StudentService) GetStudent(ID string) (stu class.Student, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&stu).Error
	return
}

// GetStudentInfoList 分页获取学生记录
// Author [yourname](https://github.com/yourname)
func (stuService *StudentService) GetStudentInfoList(info classReq.StudentSearch) (list []class.Student, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&class.Student{})
	var stus []class.Student
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.StartAge != nil && info.EndAge != nil {
		db = db.Where("age BETWEEN ? AND ? ", info.StartAge, info.EndAge)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&stus).Error
	return stus, total, err
}
func (stuService *StudentService) GetStudentPublic() {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}

// NewStudentMethod 测试学生新方法
// Author [yourname](https://github.com/yourname)
func (stuService *StudentService) NewStudentMethod() (err error) {
	// 请在这里实现自己的业务逻辑
	db := global.GVA_DB.Model(&class.Student{})
	return db.Error
}
