package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	storege "gpuprice/middleware/mysql"
	"log/slog"
	"strings"
	"time"
)

var gpuItemsTable = "gpu_items"

type GpuItem struct {
	ID        uint      `gorm:"primarykey" json:"-"`
	ItemID    string    `gorm:"column:item_id" json:"-"`
	ItemName  string    `gorm:"column:item_name" json:"-"`
	BrandID   string    `gorm:"column:brand_id" json:"-"`
	BrandName string    `gorm:"column:brand_name" json:"brand_name"`
	URL       string    `gorm:"column:url" json:"url"`
	PicURL    string    `gorm:"column:pic_url" json:"pic_url"`
	Price     float64   `gorm:"column:price" json:"price"`
	Title     string    `gorm:"column:title" json:"title"`
	Sales     int       `gorm:"column:sales" json:"sales"`
	Star      uint8     `gorm:"column:star" json:"-"`
	Memory    uint8     `gorm:"column:memory" json:"memory"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// TableName overrides the table name used by UploadImg to `upload_img`
func (GpuItem) TableName() string {
	return gpuItemsTable
}

func FindGpuByItemId(itemId string) (*GpuItem, error) {
	var gpuItem GpuItem
	result := storege.GetMysqlDB().Where("item_id = ?", itemId).First(&gpuItem)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 没有找到记录不一定是一个错误
		}
		return nil, result.Error
	}
	return &gpuItem, nil
}

func InsertItem(item GpuItem) error {

	result := storege.GetMysqlDB().Create(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateItemMemory(itemId string, memory uint8) error {
	sql := fmt.Sprintf("update %s set memory=? where item_id=? ", gpuItemsTable)
	result := storege.GetMysqlDB().Exec(sql, memory, itemId)
	if result.Error != nil {
		// 错误处理
		fmt.Printf("Error occurred while updating: %v\n", result.Error)
		return result.Error
	}
	return nil
}

func GetItemTitleMemoryPage(page, size int) (error, []*GpuItem) {
	if page < 0 {
		return fmt.Errorf("page can not be less than 0"), nil
	}
	sql := fmt.Sprintf("select item_id,title,memory from %s limit ?,?", gpuItemsTable)
	tx := storege.GetMysqlDB().Raw(sql, (page-1)*size, size)
	rows, err := tx.Rows()
	if err != nil {
		slog.Error("err:" + err.Error())
		return err, nil
	}
	defer rows.Close()
	result := make([]*GpuItem, 0)
	for rows.Next() {
		var r = new(GpuItem)
		if err = rows.Scan(&r.ItemID, &r.Title, &r.Memory); err != nil {
			slog.Error(err.Error())
			break
		}
		result = append(result, r)
	}
	return nil, result
}

func GetItemsWithCondition(whereParam map[string]any, page, size int, sort string) (error, []GpuItem) {
	if page <= 0 || size <= 0 {
		return fmt.Errorf("page or size can not be less than 0, size: %d,page: %d", size, page), nil
	}
	wheresql, vs := GenerateWheresql(whereParam)
	slog.Info("wheresql:", wheresql)
	vs = append(vs, (page-1)*size)
	vs = append(vs, size)
	sortSql := ""
	if sort != "" {
		sortSql = fmt.Sprintf("order by price %s", sort)
	}
	sql := fmt.Sprintf("select item_id,sales,title,brand_name,url,pic_url,price,memory from %s %s %s limit ?,?", gpuItemsTable, wheresql, sortSql)
	fmt.Println("sql:", sql)
	tx := storege.GetMysqlDB().Raw(sql, vs...)
	rows, err := tx.Rows()
	if err != nil {
		slog.Error("err:" + err.Error())
		return err, nil
	}
	defer rows.Close()
	result := make([]GpuItem, 0)
	for rows.Next() {
		var r GpuItem
		if err = rows.Scan(&r.ItemID, &r.Sales, &r.Title, &r.BrandName, &r.URL, &r.PicURL, &r.Price, &r.Memory); err != nil {
			slog.Error(err.Error())
			break
		}
		result = append(result, r)
	}
	return nil, result
}

func GenerateWheresql(param map[string]any) (string, []interface{}) {
	var whereSql strings.Builder
	whereSql.WriteString("where 1=1 ")
	values := make([]interface{}, len(param))
	i := 0
	for k, v := range param {
		whereSql.WriteString(" and ")
		whereSql.WriteString(k)
		whereSql.WriteString("?")
		values[i] = v
		i++
	}
	return whereSql.String(), values
}

func GetItemsCntWithCondition(whereSql string) (error, int) {
	sql := fmt.Sprintf("select count(1) from %s %s ", gpuItemsTable, whereSql)
	tx := storege.GetMysqlDB().Raw(sql)
	row := tx.Row()

	var cnt int
	err := row.Scan(&cnt)
	if err != nil {
		slog.Error("err:" + err.Error())
		return err, 0
	}
	return nil, cnt
}
