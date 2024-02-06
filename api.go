package qishutaLib

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"regexp"
	"strings"
)

var getDownURLMatches = regexp.MustCompile(`get_down_url\('.*?','(.*?)','.*?'\);`)

type APP struct {
	Request *resty.Request
	Client  *Client
}

func (app *APP) newReq(path string, query map[string]string) (ResponseInterface, error) {
	if query == nil {
		query = make(map[string]string)
	}
	response, err := app.Request.SetQueryParams(query).Get(path)
	if err != nil {
		return nil, err
	} else if response.StatusCode() != 200 {
		return nil, fmt.Errorf("status code is %v", response.StatusCode())
	} else {
		return &Response{response}, nil
	}
}
func (app *APP) GetCover(coverPath string) ([]byte, error) {
	cover, err := app.Request.Get(coverPath)
	if err != nil {
		return nil, err
	}
	if cover.StatusCode() != 200 {
		return nil, fmt.Errorf("status code is %v", cover.StatusCode())
	}
	if cover.Body() == nil {
		return nil, fmt.Errorf("cover body is nil")
	}
	return cover.Body(), nil
}
func Replace(selection *goquery.Selection, old string) string {
	return strings.ReplaceAll(selection.Text(), old, "")
}
func (app *APP) GetBookInfo(bookId string) (*BookInfoModel, error) {
	response, err := app.newReq(fmt.Sprintf("/Shtml%v.html", bookId), nil)
	if err != nil {
		return nil, err
	}
	document := response.Document()
	bookName := document.Find(".detail_right h1").Text()
	if bookName == "" {
		return nil, fmt.Errorf("book not found,bookId is %v", bookId)
	}
	// 《雷武》全集 ,取《》中的内容
	bookName = regexp.MustCompile(`《(.*)》`).FindAllStringSubmatch(bookName, -1)[0][1]
	bookInfo := &BookInfoModel{
		BookId:      bookId,
		BookName:    bookName,
		Author:      Replace(document.Find(".detail .small").Eq(5), "书籍作者："),
		Cover:       document.Find(".detail_pic img").AttrOr("src", ""),
		UpdateDate:  Replace(document.Find(".detail .small").Eq(3), "更新日期："),
		Status:      Replace(document.Find(".detail .small").Eq(4), "连载状态："),
		ClickInfo:   Replace(document.Find(".detail .small").Eq(0), "点击次数："),
		FileSize:    Replace(document.Find(".detail .small").Eq(1), "文件大小："),
		Description: document.Find(".showInfo p").First().Text(),
	}

	if matches := getDownURLMatches.FindStringSubmatch(response.String()); len(matches) > 1 {
		bookInfo.Download = matches[1]
	} else {
		bookInfo.Download = "download url not found"
	}
	return bookInfo, nil
}

func (app *APP) GetCatalogs(bookId string) ([]CatalogModel, error) {
	var chapters []CatalogModel
	response, err := app.newReq(fmt.Sprintf("/du/%v/%v/", bookId[0:1], bookId), nil)
	if err != nil {
		return nil, err
	}
	response.Document().Find("#info .pc_list").Eq(1).Find("ul li").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Find("a").Attr("href")
		chapters = append(chapters, CatalogModel{
			BookID:       bookId,
			ChapterIndex: fmt.Sprintf("%v", i+1),
			ChapterTitle: s.Find("a").Text(),
			ChapterId:    strings.ReplaceAll(url, ".html", ""),
		})
	})
	return chapters, nil

}

func (app *APP) GetContent(bookId string, chapterId string) (*ContentModel, error) {
	response, err := app.newReq(fmt.Sprintf("/du/%v/%v/%v.html", bookId[:2], bookId, chapterId), nil)
	if err != nil {
		return nil, err
	}
	var content string

	document := response.Document()
	document.Find("#content1").Contents().Each(func(i int, s *goquery.Selection) {
		if !s.Is("p") {
			content += s.Text()
		}
	})
	if strings.TrimSpace(content) == "" {
		return nil, fmt.Errorf("content is empty,bookId is %v,chapterId is %v", bookId, chapterId)
	}
	return &ContentModel{
		ID:           chapterId,
		BookId:       bookId,
		ChapterId:    chapterId,
		ChapterTitle: document.Find("h1").First().Text(),
		Content:      content,
		ChapterWord:  len([]rune(content)),
	}, nil

}
func (app *APP) GetSearch(keyword string) ([]SearchModel, error) {
	response, err := app.newReq("/search.html", map[string]string{"searchkey": keyword})
	if err != nil {
		return nil, err
	}
	var books []SearchModel
	response.Document().Find("table.grid tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		bookId, _ := s.Find("td:nth-child(2) a").Attr("href")
		books = append(books, SearchModel{
			Index:         fmt.Sprintf("%v", i+1),
			BookId:        regexp.MustCompile(`\d+`).FindAllString(bookId, -1)[0],
			BookName:      s.Find("td:nth-child(2) a").Text(),
			BookAuthor:    s.Find("td:nth-child(3)").Text(),
			LatestChapter: s.Find("td:nth-child(4) a").Text(),
			Update:        s.Find("td:nth-child(5)").Text(),
		})
	})
	return books, nil
}

func (app *APP) GetBookTypeList(typeNumber int, page int) ([]TypeListBookInfoModel, error) {
	var books []TypeListBookInfoModel
	response, err := app.newReq(fmt.Sprintf("soft/sort0%v/index_%v.html", typeNumber, page), nil)
	if err != nil {
		return nil, err
	}
	response.Document().Find(".list ul li").Each(func(i int, s *goquery.Selection) {
		books = append(books, TypeListBookInfoModel{
			Index:       fmt.Sprintf("%v", i+1),
			BookName:    s.Find("a").First().Text(),
			BookId:      regexp.MustCompile(`\d+`).FindAllString(s.Find("a").AttrOr("href", ""), -1)[0],
			Cover:       s.Find("img").AttrOr("src", ""),
			Description: s.Find(".u").Text(),
		})
	})

	return books, nil

}
func (app *APP) GetBookshelf() ([]BookshelfModel, error) {
	response, err := app.newReq("/bookcase.php", nil)
	if err != nil {
		return nil, err
	}
	var books []BookshelfModel
	response.Document().Find(".grid tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		bookId, _ := s.Find("td:nth-child(2) a").Attr("href")
		books = append(books, BookshelfModel{
			Index:         fmt.Sprintf("%v", i+1),
			BookId:        regexp.MustCompile(`\d+`).FindAllString(bookId, -1)[0],
			BookName:      s.Find("td:nth-child(2) a").Text(),
			LatestChapter: s.Find("td:nth-child(3) a").Text(),
			UpdateDate:    s.Find("td:nth-child(5)").Text(),
		})
	})
	return books, nil
}
func (app *APP) GetCookie(username string, password string) (string, error) {
	response, err := app.Request.SetFormData(map[string]string{
		"LoginForm[username]": username,
		"LoginForm[password]": password,
		"submit":              "%E7%99%BB%C2%A0%C2%A0%E5%BD%95",
	}).Post("/login.php")
	if err != nil {
		return "", err
	} else if response.Cookies() == nil {
		return "", fmt.Errorf("cookies is nil,username or password is wrong")
	} else {
		var cookie string
		for _, c := range response.Cookies() {
			cookie += fmt.Sprintf("%v=%v; ", c.Name, c.Value)
		}
		return cookie, nil
	}
}
