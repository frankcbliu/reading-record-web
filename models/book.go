package models

import (
	"fmt"

	"gorm.io/gorm"
)

// 记录从豆瓣搜索得到的书籍信息
type Book struct {
	gorm.Model
	BookName    string
	Author      string
	Publication string
	SubTitle    string
	Year        string
	Pages       string
	ISBN        string
	ImageUrl    string
}

// 用于 routers 使用
type BookInfo struct {
	BookName    string `json:"book_name"`   // 书名
	Author      string `json:"author"`      // 作者
	Publication string `json:"publication"` // 出版社
	SubTitle    string `json:"sub_title"`   // 副标题
	Year        string `json:"year"`        // 出版年份
	Pages       string `json:"pages"`       // 页数
	ISBN        string `json:"ISBN"`        // ISBN
	ImageUrl    string `json:"image_url"`   // 图片地址
}

// 根据书名查询结果
func GetBookInfo(book_name string, book_infos *[]BookInfo) bool {
	var books []Book

	db := InitModel()
	db.Find(&books, "book_name = ?", book_name)

	if len(books) <= 0 {
		fmt.Println("[models.book] 查找不到: ", len(books))
		return false
	}
	fmt.Println("[models.book] 查询结果: ", len(books))

	for _, v := range books {
		var book_info BookInfo
		book_info.BookName = v.BookName
		book_info.Author = v.Author
		book_info.Publication = v.Publication
		book_info.SubTitle = v.SubTitle
		book_info.Year = v.Year
		book_info.Pages = v.Pages
		book_info.ISBN = v.ISBN
		book_info.ImageUrl = v.ImageUrl
		*book_infos = append(*book_infos, book_info)
	}

	fmt.Println("[models.book] Use Cache data!")
	return true
}

// 缓存查询后的书籍内容
func CacheBookInfo(book_info BookInfo) bool {
	var book Book
	if book_info.ISBN == "" {
		fmt.Println("[models.book][Cache Error] book: ", book_info, " 's isbn is empty.")
		return false
	}

	if QueryBookModel(&book, book_info.ISBN) {
		fmt.Println("[models.book][Cache Error] book: ", book_info.ISBN, " is exist.")
		return false
	}

	db := InitModel()
	db.Create(&Book{
		BookName: book_info.BookName, Author: book_info.Author, Publication: book_info.Publication,
		SubTitle: book_info.SubTitle, Year: book_info.Year, Pages: book_info.Pages,
		ISBN: book_info.ISBN, ImageUrl: book_info.ImageUrl})

	fmt.Println("[models.book] Cache data success!")
	return true
}

// 查询 book 是否存在
func QueryBookModel(book *Book, isbn string) bool {
	db := InitModel()
	db.First(book, "isbn = ?", isbn)
	if book.ID > 0 {
		fmt.Println("[models.book]", book.ISBN, "exist")
		return true
	}
	fmt.Println("[models.book]", book.ISBN, "not exist")
	return false
}
