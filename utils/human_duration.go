package utils

import (
	"strconv"
	"strings"
	"time"
)

// ParseDuration 将字符串类型的时间解析为时间戳
func ParseDuration(d string) (time.Duration, error) {
	d = strings.TrimSpace(d)
	dr, err := time.ParseDuration(d)
	if err == nil {
		return dr, nil
	}
	if strings.Contains(d, "d") {
		index := strings.Index(d, "d")

		hour, _ := strconv.Atoi(d[:index])
		dr = time.Hour * 24 * time.Duration(hour)   //转换为小时为单位
		ndr, err := time.ParseDuration(d[index+1:]) //查看后面是否还有小时制
		if err != nil {
			return dr, nil
		}
		return dr + ndr, nil
	}

	dv, err := strconv.ParseInt(d, 10, 64)
	return time.Duration(dv), err
}
