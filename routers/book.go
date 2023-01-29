package routers

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"reading-record-web/models"
	"reading-record-web/tools"

	"github.com/crawlab-team/crawlab-go-sdk/entity"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

// /book/search?book_name=?
func SearchHandler(c *gin.Context) {
	book_name := c.Query("book_name")
	fmt.Println("[routers.book] book_name: ", book_name)

	var data []models.BookInfo
	// 如果书名为空，直接跳过; 书名非空时，先查询数据库
	if book_name != "" && !models.GetBookInfo(book_name, &data) {
		fmt.Println("[routers.book] Find Cache error. ready to query from douban.")
		BuildCollectorAndRun(book_name, &data)
	}

	response := gin.H{"code": 0, "msg": "success", "data": data}
	c.IndentedJSON(http.StatusOK, response)
}

// 构建爬虫并执行
func BuildCollectorAndRun(book_name string, data *[]models.BookInfo) {
	c1 := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36"),
	)
	c2 := c1.Clone()

	// Find and visit all links
	c1.OnHTML(".result", func(e *colly.HTMLElement) {
		// 判断有无评分信息，没有就跳过
		if e.ChildText(".rating-info") == "" {
			return
		}
		// 解析信息
		item := entity.Item{
			"type":     tools.TrimSpace(e.ChildText(".title > h3 > span")),
			"next_url": e.ChildAttr(".title > h3 > a", "href"),
		}
		if !strings.Contains(item["type"].(string), "书籍") {
			return
		}
		// 解析具体跳转链接
		c2.Visit(e.Request.AbsoluteURL(item["next_url"].(string)))
	})

	c1.OnRequest(func(r *colly.Request) {
		fmt.Println("[C-1]Visiting", r.URL)
	})

	c2.OnRequest(func(r *colly.Request) {
		fmt.Println("[C-2]Visiting", r.URL)
	})

	c1.OnError(func(r *colly.Response, err error) {
		fmt.Println("[C-1]Visiting", r.Request.URL, "failed:", err)
	})

	c2.OnError(func(r *colly.Response, err error) {
		fmt.Println("[C-2]Visiting", r.Request.URL, "failed:", err)
	})

	c2.OnHTML(".subject", func(e *colly.HTMLElement) {
		lines := strings.Split(e.ChildText("#info"), "\n")
		info_text := ""
		for i := 0; i < len(lines); i++ {
			text := strings.Trim(lines[i], " \n")
			if text != "" {
				info_text += text
			}
		}
		book_info := FormatBookInfoV2(info_text)
		book_info.BookName = book_name
		book_info.ImageUrl = tools.TrimSpace(e.ChildAttr(".nbg > img", "src"))
		models.CacheBookInfo(book_info)
		*data = append(*data, book_info)
	})

	c1.Visit("https://www.douban.com/search?cat=1001&q=" + url.QueryEscape(book_name))
	c1.Wait()
	c2.Wait()
}

func FormatBookInfoV2(info_text string) models.BookInfo {
	var book_info models.BookInfo
	var indexList []int
	indexList = append(indexList, strings.Index(info_text, "作者"))
	indexList = append(indexList, strings.Index(info_text, "出版社"))
	indexList = append(indexList, strings.Index(info_text, "出品方"))
	indexList = append(indexList, strings.Index(info_text, "副标题"))
	indexList = append(indexList, strings.Index(info_text, "原作名"))
	indexList = append(indexList, strings.Index(info_text, "出版年"))
	indexList = append(indexList, strings.Index(info_text, "译者"))
	indexList = append(indexList, strings.Index(info_text, "页数"))
	indexList = append(indexList, strings.Index(info_text, "定价"))
	indexList = append(indexList, strings.Index(info_text, "装帧"))
	indexList = append(indexList, strings.Index(info_text, "ISBN"))
	sort.Ints(indexList)

	// 分割字符串
	var splitStrList []string
	for i := 0; i < len(indexList)-1; i++ {
		// 可以分割
		if indexList[i] > -1 {
			splitStrList = append(splitStrList, info_text[indexList[i]:indexList[i+1]])
		}
	}
	if indexList[len(indexList)-1] > -1 {
		splitStrList = append(splitStrList, info_text[indexList[len(indexList)-1]:])
	}

	// 处理分割后的字符串
	for _, v := range splitStrList {
		if strings.Contains(v, "作者") {
			book_info.Author = SplitAndGetStr(v)
		} else if strings.Contains(v, "出版社") {
			book_info.Publication = SplitAndGetStr(v)
		} else if strings.Contains(v, "副标题") {
			book_info.SubTitle = SplitAndGetStr(v)
		} else if strings.Contains(v, "出版年") {
			book_info.Year = SplitAndGetStr(v)
		} else if strings.Contains(v, "ISBN") {
			book_info.ISBN = SplitAndGetStr(v)
		} else if strings.Contains(v, "页数") {
			book_info.Pages = SplitAndGetStr(v)
		}
	}
	return book_info
}

// 用 : 分割并取出值
func SplitAndGetStr(src string) string {
	i := strings.Index(src, ":")
	return tools.TrimSpace(src[i+1:])
}
