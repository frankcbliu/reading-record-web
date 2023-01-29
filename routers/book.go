package routers

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"reading-record-web/tools"

	"github.com/crawlab-team/crawlab-go-sdk/entity"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

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

func SearchHandler(c *gin.Context) {
	book_name := c.Query("book_name")

	var data []BookInfo
	BuildCollectorAndRun(book_name, &data)
	response := gin.H{"code": 0, "msg": "success", "data": data}
	c.IndentedJSON(http.StatusOK, response)
}

func TestGo(c *gin.Context) {
	books := []string{
		"作者:凯利•麦格尼格尔/Kelly McGonigal Ph.D.出版年: 2012-8-21装帧: PaperbackISBN: 9787993417657",
		"作者:【美】凯利·麦格尼格尔出版社:北京联合出版公司出品方:磨铁图书副标题: 斯坦福大学掌控自我的心理学课程原作名: The Joy of Movement:  How exercise helps us find happiness, hope, connection, and courage译者:江兰, 张旭, 刘婉婷出版年: 2021-9-1页数: 280定价: 55装帧: 平装ISBN: 9787559649508",
		"作者:李萌出版社:成都地图出版社出版年: 2018-10-1定价: 29.8装帧: 平装ISBN: 9787555710363",
		"作者:（美）凯利·麦格尼格尔出版社: 黑天鹅图书·北京联合出版公司出品方:磨铁图书副标题: 斯坦福高效实用的25堂心理学课译者:金磊出版年: 2018-9页数: 257定价: 48装帧: 平装ISBN: 9787559620941",
		"作者:凯利•麦格尼格尔出版社:印刷工业出版社出品方:磨铁图书副标题: 瑜伽实操篇原作名: Yoga for Pain Relief: Simple Practices to Calm Your Mind and Heal Your Chronic Pain译者:王岑卉出版年: 2013-5页数: 214定价: 39.8装帧: 平装ISBN: 9787514207804",
		"作者:【美】凯利·麦格尼格尔出版社:北京联合出版公司出品方:磨铁图书副标题: 斯坦福大学掌控情绪的心理学课程原作名: The Upside of Stress: Why Stress Is Good for You, and How to Get Good at It译者:王鹏程出版年: 2021-9-1页数: 256定价: 55装帧: 平装ISBN: 9787559649430",
		"作者:[美] 凯利·麦格尼格尔出版社: 文化发展出版社(原印刷工业出版社)出品方:磨铁图书副标题: 斯坦福大学最受欢迎心理学课程原作名: The Willpower Instinct:How Self-control Works,Why it Matters,and What You Can do to Get More of It译者:王岑卉出版年: 2012-8页数: 263定价: 39.80元装帧: 平装ISBN: 9787514205039",
		"作者:【美】凯利·麦格尼格尔出版社:北京联合出版公司出品方:磨铁图书副标题: 斯坦福大学广受欢迎的心理学课程原作名: THE WILLPOWER INSTINCT译者:王岑卉出版年: 2021-4-8页数: 256定价: 55装帧: 平装ISBN: 9787559647276",
		"作者:[美]凯利·麦格尼格尔出版社:北京联合出版公司出品方:磨铁图书副标题: 斯坦福大学最实用的心理学课程原作名: The Upside of Stress: Why Stress Is Good for You, and How to Get Good at It译者:王鹏程出版年: 2016-3-1页数: 272定价: 39.80元装帧: 平装ISBN: 9787550267831",
	}

	var data []BookInfo
	for _, v := range books {
		data = append(data, FormatBookInfoV2(v))
	}
	response := gin.H{"code": 0, "msg": "success", "data": data}
	c.IndentedJSON(http.StatusOK, response)
}

// 构建爬虫并执行
func BuildCollectorAndRun(book_name string, data *[]BookInfo) {
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
			"title":     tools.TrimSpace(e.ChildText(".title > h3 > a")),
			"image_url": e.ChildAttr(".pic > .nbg > img", "src"),
			"type":      tools.TrimSpace(e.ChildText(".title > h3 > span")),
			"next_url":  e.ChildAttr(".title > h3 > a", "href"),
		}
		if !strings.Contains(item["type"].(string), "书籍") {
			fmt.Println("[title]: ", item["title"].(string))
			fmt.Println("[type]: ", item["type"].(string))
			fmt.Println("[next_url]: ", item["next_url"].(string))
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

	c2.OnHTML(".subject", func(e *colly.HTMLElement) {
		lines := strings.Split(e.ChildText("#info"), "\n")
		info_text := ""
		for i := 0; i < len(lines); i++ {
			text := strings.Trim(lines[i], " \n")
			if text != "" {
				info_text += text
			}
		}
		fmt.Println(info_text)
		book_info := FormatBookInfoV2(info_text)
		book_info.BookName = book_name
		book_info.ImageUrl = tools.TrimSpace(e.ChildAttr(".nbg > img", "src"))
		*data = append(*data, book_info)
	})

	c1.Visit("https://www.douban.com/search?q=" + url.QueryEscape(book_name))
	c1.Wait()
	c2.Wait()
}

func FormatBookInfo(info_text string) BookInfo {
	var book_info BookInfo
	reg := regexp.MustCompile(`作者:(.*?)出版社:(.*?)副标题:(.*?)原作名:(.*?)译者:(.*?)出版年:(.*?)页数:(.*?)定价:(.*?)装帧:(.*?)ISBN:(.*)`)
	result1 := reg.FindAllStringSubmatch(info_text, -1)
	if len(result1) > 0 {
		book_info.Author = tools.TrimSpace(result1[0][1])
		book_info.Publication = tools.TrimSpace(result1[0][2])
		book_info.SubTitle = tools.TrimSpace(result1[0][3])
		book_info.Year = tools.TrimSpace(result1[0][6])
		book_info.Pages = tools.TrimSpace(result1[0][7])
		book_info.ISBN = tools.TrimSpace(result1[0][10])
		for i, v := range result1[0] {
			fmt.Println(i, tools.TrimSpace(v))
		}
	}
	fmt.Println("[BookInfo] Parse error. info:")
	fmt.Println(info_text)
	fmt.Println("[BookInfo] info over.")
	return book_info
}

func FormatBookInfoV2(info_text string) BookInfo {
	var book_info BookInfo
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
	fmt.Println(indexList)

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
		fmt.Println(v)
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
	fmt.Println(i)
	return tools.TrimSpace(src[i+1:])
}
