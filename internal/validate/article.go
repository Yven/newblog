package validate

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type List struct {
	Keyword  string `json:"keyword"`
	Tag      int64  `json:"tag"`
	Category int64  `json:"category"`
}

type Article struct {
	Slug    string `json:"slug" binding:"required,alphanumunicode"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	TagList string `json:"tag_list" binding:"required,idStringList"`
	Cid     int64  `json:"cid" binding:"required,number"`
}

var IdStringList validator.Func = func(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if ok {
		idList := strings.SplitSeq(data, ",")
		// 用于存储转换后的标签ID
		for id := range idList {
			// 去除空格
			id = strings.TrimSpace(id)
			if id == "" {
				continue
			}
			// 将字符串转换为int64
			tagId, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				return false
			}
			if tagId <= 0 {
				return false
			}
		}
	}

	return true
}
