package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const url = "https://api.followin.io/feed/list/recommended/news"

type FirstRequestBody struct {
	OnlyImportant bool `json:"only_important"`
	Count         int  `json:"count"`
	Page          int  `json:"page"`
}

type RequestBody struct {
	OnlyImportant bool   `json:"only_important"`
	Count         int    `json:"count"`
	Page          int    `json:"page"`
	LastCursor    string `json:"last_cursor"`
	LastSource    string `json:"last_source"`
}

type ResponseBody struct {
	Code int          `json:"code"`
	Data ResponseData `json:"data"`
	Msg  string       `json:"msg"`
}

type ResponseData struct {
	List         []Feed `json:"list"`
	RecRequestID string `json:"rec_request_id"`
	HasMore      bool   `json:"has_more"`
	LastCursor   string `json:"last_cursor"`
	LastSource   string `json:"last_source"`
	Source       string `json:"source"`
}

type Feed struct {
	ID                int64  `json:"id"`
	Title             string `json:"title"`
	Content           string `json:"content"`
	TranslatedTitle   string `json:"translated_title"`
	TranslatedContent string `json:"translated_content"`
	PublishTime       int64  `json:"publish_time"`
	Important         bool   `json:"important"`
	SourceURL         string `json:"source_url"`
	RelatedFeeds      int    `json:"related_feeds"`
	Nickname          string `json:"nickname"`
}

func BaseReqBody() RequestBody {
	return RequestBody{
		OnlyImportant: true,
		Count:         20,
		Page:          2,
		LastSource:    "algo",
	}
}

func FirstReqBody() FirstRequestBody {
	return FirstRequestBody{
		OnlyImportant: true,
		Count:         20,
		Page:          1,
	}
}

func (r RequestBody) Update(lc string) RequestBody {
	r.Page = r.Page + 1
	r.LastCursor = lc
	return r
}

func Post(requestBody RequestBody) ([]Feed, string) {
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Failed to marshal request body:", err)
		return nil, ""
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		fmt.Println("Failed to send POST request:", err)
		return nil, ""
	}
	defer resp.Body.Close()

	// 讀取回應內容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return nil, ""
	}

	resBody := ResponseBody{}
	err = json.Unmarshal(body, &resBody)
	if err != nil {
		fmt.Println(err)
		return nil, ""
	}

	if resBody.Code != 2000 {
		fmt.Println(resBody.Msg)
	}

	return resBody.Data.List, resBody.Data.LastCursor
}

func First() ([]Feed, string) {
	first := FirstReqBody()
	fmt.Println("第", 1, "次请求")
	requestBodyBytes, err := json.Marshal(first)
	if err != nil {
		fmt.Println("Failed to marshal request body:", err)
		return nil, ""
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		fmt.Println("Failed to send POST request:", err)
		return nil, ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return nil, ""
	}

	resBody := ResponseBody{}
	err = json.Unmarshal(body, &resBody)
	if err != nil {
		fmt.Println(err)
		return nil, ""
	}

	if resBody.Code != 2000 {
		fmt.Println(resBody.Msg)
	}

	list := resBody.Data.List
	l := len(list)
	last := list[l-1]
	fmt.Println("当前结尾标题", last.Title)
	fmt.Println("当前结尾值时间", last.PublishTime, JSUnixTimeToTime(int(last.PublishTime)))

	return resBody.Data.List, resBody.Data.LastCursor
}

func Start() []Feed {
	base := BaseReqBody()
	feedList, lc := First()
	base.LastCursor = lc
	for {
		fmt.Println("第", base.Page, "次请求")
		list, lc := Post(base)
		l := len(list)
		last := list[l-1]
		fmt.Println("当前结尾标题", last.Title)
		fmt.Println("当前结尾值时间", last.PublishTime, JSUnixTimeToTime(int(last.PublishTime)))
		feedList = append(feedList, list...)
		base = base.Update(lc)
		if TimeCaseCmp(JSUnixTimeToTime(int(last.PublishTime)), GetCurrentDay()) == false && TimeCaseCmp(JSUnixTimeToTime(int(last.PublishTime)), GetPreDay()) == false {
			break
		}
	}
	return feedList
}
