package bean

import (
	"reflect"
)

var pageConf = pageConfig{
	MinPage:     1,
	MaxSize:     100,
	DefaultSize: 10,
}

type pageConfig struct {
	MinPage     int // 最小分页
	MaxSize     int // 最大条目
	DefaultSize int // 默认条目
}

// SetPageConfig 设置分页参数
func SetPageConfig(minPage, maxSize, defaultSize int) {
	pageConf = pageConfig{
		MinPage:     minPage,
		MaxSize:     maxSize,
		DefaultSize: defaultSize,
	}
}

// 入参
type PageParam struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

func (p *PageParam) GetOffset() int {
	if p.Page < pageConf.MinPage {
		p.Page = pageConf.MinPage
	}
	return (p.Page - 1) * p.Size
}

func (p *PageParam) GetLimit() int {
	if p.Size > pageConf.MaxSize {
		p.Size = pageConf.MaxSize
	}
	if p.Size == 0 {
		p.Size = pageConf.DefaultSize
	}
	return p.Size
}

func (p *PageParam) HasPage(total int64) bool {
	return total > int64(p.GetOffset())
}

type ResultPage struct {
	PageParam
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

func NewResultPage(param PageParam, list any, total int64) ResultPage {
	of := reflect.ValueOf(list)
	switch of.Kind() {
	case reflect.Slice:
		if of.Len() == 0 {
			list = make([]int, 0)
		}
	}
	return ResultPage{
		PageParam: param,
		List:      list,
		Total:     total,
	}
}
