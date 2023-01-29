package models

import "gorm.io/gorm"

// 计时单元
type RecordItem struct {
	gorm.Model
	RecordType string // 计时类型 book/(预留)

	// 书本信息
	BookName      string // 书名
	Author        string // 作者
	Publication   string // 出版社
	PublicateDate string // 出版日期
	// 进度信息
	AllPage     uint32 // 总页数
	CurrentPage uint32 // 当前页数

	// 书籍管理
	ReadingStatus string // 阅读状态: 正在读/想读/已读/弃读/闲置
	BookSelfId    uint32 // 书架ID
	Label         string // 标签: 用,分隔
}
