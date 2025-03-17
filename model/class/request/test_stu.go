package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type StudentSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	StartAge       *int       `json:"startAge" form:"startAge"`
	EndAge         *int       `json:"endAge" form:"endAge"`
	request.PageInfo
}
