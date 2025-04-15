package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/blog"
	"github.com/flipped-aurora/gin-vue-admin/server/model/class"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(class.Student{}, blog.Article{})
	if err != nil {
		return err
	}
	return nil
}
