package database

import "gorm.io/gorm"

type PageQuery struct {
	PageNum  int `json:"pageNum" form:"pageNum" binding:"required,min=1" default:"1"`             // 页码
	PageSize int `json:"pageSize" form:"pageSize" binding:"required,min=1,max=500" default:"100"` // 分页大小
}

// Page 分页模型
type Page struct {
	PageNum   int         `json:"pageNum"`   // 页码
	PageSize  int         `json:"pageSize"`  // 分页大小
	Total     int         `json:"total"`     // 数据总量
	PageCount int         `json:"pageCount"` // 分页数量
	Result    interface{} `json:"result"`    // 分页大小
}

// GetOffset 获取分页偏移量
func (p *Page) GetOffset() int {
	if p.PageNum < 1 {
		p.PageNum = 1
	}
	return (p.PageNum - 1) * p.PageSize
}

func (p *Page) GetLimit() int {
	return p.PageSize
}

// done 执行完毕
func (p *Page) done() {
	if p.Total == 0 {
		return
	}
	p.PageCount = p.Total / p.PageSize
	if p.Total%p.PageSize != 0 {
		p.PageCount++
	}
}

// Execute 执行分页查询
func (p *Page) Execute(db *gorm.DB, result interface{}) error {
	if p.PageNum <= 0 {
		p.PageNum = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	// 查询具体明细
	err := db.Offset(p.GetOffset()).Limit(p.GetLimit()).Find(result).Error
	if err != nil {
		return err
	}
	p.Result = result

	// 查询总行数
	var total int64
	err = db.Offset(-1).Limit(-1).Count(&total).Error

	if err != nil {
		return err
	}
	p.Total = int(total)

	p.done()
	return nil
}
