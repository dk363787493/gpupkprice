package test

import (
	"fmt"
	"gpuprice/middleware/mysql/model"
	"os"
	"testing"
)

func __TestQuery(t *testing.T) {
	os.Setenv("CONFIG_PATH", "/Users/zhangjinge/GolandProjects/gpuprice/config")
	err, items := model.GetItemTitleMemoryPage(1, 10)
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}

	fmt.Println("items:", len(items))
}
