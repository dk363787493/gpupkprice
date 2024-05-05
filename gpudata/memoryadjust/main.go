package main

import (
	"fmt"
	"gpuprice/middleware/mysql/model"
	"log/slog"
	"strings"
)

func main() {
	memoryStrs := []uint8{
		4,
		6,
		8,
		12,
		16,
	}
	notMatch := make([]string, 0)
	for i := 1; ; i++ {
		err, items := model.GetItemTitleMemoryPage(i, 20)
		if err != nil || len(items) == 0 {
			break
		}
		for _, item := range items {
			f := false
			if strings.Contains(item.Title, fmt.Sprintf("%d GB", item.Memory)) ||
				strings.Contains(item.Title, fmt.Sprintf("%dGB", item.Memory)) ||
				strings.Contains(item.Title, fmt.Sprintf("%d gb", item.Memory)) ||
				strings.Contains(item.Title, fmt.Sprintf("%dgb", item.Memory)) ||
				strings.Contains(item.Title, fmt.Sprintf("%d G", item.Memory)) ||
				strings.Contains(item.Title, fmt.Sprintf("%dG", item.Memory)) ||
				strings.Contains(item.Title, fmt.Sprintf("%d g", item.Memory)) ||
				strings.Contains(item.Title, fmt.Sprintf("%dg", item.Memory)) {
				fmt.Println("db,matched,item id:" + item.ItemID)
				f = true
				continue
			}
			for _, mInt := range memoryStrs {
				if strings.Contains(item.Title, fmt.Sprintf("%d GB", mInt)) ||
					strings.Contains(item.Title, fmt.Sprintf("%dGB", mInt)) ||
					strings.Contains(item.Title, fmt.Sprintf("%d gb", mInt)) ||
					strings.Contains(item.Title, fmt.Sprintf("%dgb", mInt)) ||
					strings.Contains(item.Title, fmt.Sprintf("%d G", mInt)) ||
					strings.Contains(item.Title, fmt.Sprintf("%dG", mInt)) ||
					strings.Contains(item.Title, fmt.Sprintf("%d g", mInt)) ||
					strings.Contains(item.Title, fmt.Sprintf("%dg", mInt)) {
					_ = model.UpdateItemMemory(item.ItemID, mInt)
					f = true
					fmt.Println("manual,matched,item id:" + item.ItemID)
					break
				}

			}
			if !f {
				slog.Info("can not match memory of item:" + item.ItemID + ",tile:" + item.Title)
				notMatch = append(notMatch, item.ItemID)
			}

		}
	}

	for _, itemId := range notMatch {
		fmt.Println(itemId)
	}
}
