package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gpuprice/middleware/mysql/model"
	"log/slog"
	"net/http"
	"strconv"
)

type SearchItemsRequest struct {
	BrandName string  `json:"brand_name"`
	Price     float64 `json:"price"`
	Sales     int     `json:"sales"`
	Start     uint8   `json:"start"`
}

type SearchItemsRespone struct {
	Data []model.GpuItem `json:"data"`
}

func SearchItemsHandler(c *gin.Context) {
	//var searchItemsRequest SearchItemsRequest
	//if err := c.BindJSON(&searchItemsRequest); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	param := make(map[string]any)
	// 获取URL查询参数
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	size, _ := strconv.ParseInt(c.DefaultQuery("size", "10"), 10, 64)
	brandName := c.Query("brand_name")
	minSales := c.Query("min_sales")
	maxSales := c.Query("max_sales")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	memory := c.Query("memory")
	star := c.Query("star")
	priceSort := c.Query("price_sort")

	if brandName != "" {
		param["brand_name="] = brandName
	} else {
		param["brand_name<>"] = ""
	}
	if minSales != "" {
		m, _ := strconv.ParseInt(minSales, 10, 64)
		k := fmt.Sprintf("sales>=")
		param[k] = m
	}
	if maxSales != "" {
		m, _ := strconv.ParseInt(maxSales, 10, 64)
		k := fmt.Sprintf("sales<=")
		param[k] = m
	}
	if minPrice != "" {
		m, _ := strconv.ParseFloat(minPrice, 64)
		k := fmt.Sprintf("price>=")
		param[k] = m
	}
	if maxPrice != "" {
		m, _ := strconv.ParseFloat(maxPrice, 64)
		k := fmt.Sprintf("price<=")
		param[k] = m
	}
	if star != "" {
		m, _ := strconv.ParseInt(star, 10, 64)
		k := fmt.Sprintf("star=")
		param[k] = m
	}
	if memory != "" {
		m, _ := strconv.ParseInt(memory, 10, 64)
		k := fmt.Sprintf("memory=")
		param[k] = m
		param["memory"] = m
	}

	err, items := model.GetItemsWithCondition(param, int(page), int(size), priceSort)
	if err != nil {
		slog.Error("err:", err.Error())
	}

	fmt.Println("len of items:", len(items))
	items = AppendUrlTag(items)
	c.JSON(http.StatusOK, SearchItemsRespone{items})
}

func AppendUrlTag(items []model.GpuItem) []model.GpuItem {
	for i := 0; i < len(items); i++ {
		itemPoiter := &(items[i])
		itemPoiter.URL = fmt.Sprintf("%s&tag=%s", itemPoiter.URL, "217308-20")
	}
	return items
}
