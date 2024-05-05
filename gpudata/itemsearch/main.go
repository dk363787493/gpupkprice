package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"gpuprice/middleware/mysql/model"
	"io/ioutil"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

var apiKey = "t3760372173"
var secret = "20240417"
var memory = 16

func main() {
	apis := []string{
		"https://api-1.onebound.cn/amazon",
		"https://api-2.onebound.cn/amazon",
		"https://api-3.onebound.cn/amazon",
	}
	apiSign := 1
	var page int64 = 1

	keyWord := fmt.Sprintf("GPU+%dGB", memory)
	var param = []any{
		apis[0],
		apiKey,
		secret,
		keyWord,
		page,
	}
	//url := fmt.Sprintf("https://api-gw.onebound.cn/amazon/item_search/?key=<您自己的apiKey>&secret=<您自己的apiSecret>&q=鞋子&start_price=&end_price=&page=&cat=&discount_only=&sort=&page_size=&seller_info=&nick=&ppath=", params)
	for {
		err, result := doRequest(param)
		if err != nil {
			apiSign++
			api := apis[apiSign%3]
			param[0] = api
			fmt.Println("api:", api)
			time.Sleep(time.Second)
			continue
		}
		pagecount := result.Get("items.pagecount").Int()
		slog.Info("items.page:", page)
		slog.Info("\nitems.pagecount:", pagecount)
		err = CreateItem(result)
		if err != nil {
			apiSign++
			api := apis[apiSign%3]
			param[0] = api
			fmt.Println("api:", api)
			time.Sleep(time.Second)
			continue
		}
		if page >= pagecount {
			break
		}
		page++
		param[4] = page
		fmt.Println("result:", result.String())
	}

}

func CreateItem(r *gjson.Result) error {

	items := r.Get("items.item").Array()
	if len(items) == 0 {
		slog.Warn("len of items is 0")
		slog.Warn("result:", r.String())
		return fmt.Errorf("len of items is 0")
	}
	for _, item := range items {
		title := item.Get("title").String()
		picUrl := item.Get("pic_url").String()
		priceStr := item.Get("price").String()
		numIid := item.Get("num_iid").String()
		sales := item.Get("sales").Int()
		detailUrl := item.Get("detail_url").String()
		price, _ := strconv.ParseFloat(priceStr, 64)
		data := model.GpuItem{
			ItemID: numIid,
			Title:  title,
			PicURL: picUrl,
			Sales:  int(sales),
			URL:    detailUrl,
			Price:  price,
			Memory: uint8(memory),
		}
		err := model.InsertItem(data)
		if err != nil {
			slog.Error("err:", err)
		}
	}

	return nil
}

func doRequest(param []any) (error, *gjson.Result) {
	url := fmt.Sprintf("%s/item_search/?key=%s&secret=%s&q=%s&page=%d", param...)
	slog.Info("request url:" + url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("err:", err.Error())
		return err, nil
	}
	req.Header.Set("Authorization", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("err:", err.Error())
		return err, nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		slog.Error("err:", err.Error())
		return err, nil
	}
	r := gjson.Parse(string(body))
	return nil, &r
}
