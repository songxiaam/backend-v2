package gormx

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type GroupKeyFunc[T any] func(item *T) string

type QueryKeyFunc[T any] func(item *T) string

// SaveBatch 通用批量存储处理
func SaveBatch[T any](crmDB *gorm.DB, tableName string, items []*T, groupKey GroupKeyFunc[T], queryFields []string, queryKey QueryKeyFunc[T]) error {
	if len(items) == 0 {
		return nil
	}

	// 数据分组
	groupedItems := make(map[string][]*T)
	for _, item := range items {
		key := groupKey(item)
		groupedItems[key] = append(groupedItems[key], item)
	}

	// 遍历分组数据
	for key, group := range groupedItems {
		// 从分组键中解析查询条件
		conditions := strings.Split(key, ":")

		// 提取查询字段的值
		var queryValues []string
		for _, item := range group {
			queryValues = append(queryValues, queryKey(item))
		}

		// 查询已存在的数据
		var existing []map[string]interface{}
		if err := crmDB.Table(tableName).
			Select(queryFields).
			Where(fmt.Sprintf("%s IN ?", queryFields[0]), queryValues).
			Find(&existing).Error; err != nil {
			return err
		}

		// 构建存在的数据集合
		existingMap := make(map[string]bool)
		for _, e := range existing {
			key := fmt.Sprintf("%v:%v:%v", e[queryFields[0]], e[queryFields[1]], e[queryFields[2]])
			existingMap[key] = true
		}

		// 分出需要插入和更新的部分
		var toInsert, toUpdate []*T
		for _, item := range group {
			key := queryKey(item)
			if existingMap[key] {
				toUpdate = append(toUpdate, item)
			} else {
				toInsert = append(toInsert, item)
			}
		}

		// 批量插入（忽略冲突）
		if len(toInsert) > 0 {
			if err := crmDB.Clauses(clause.Insert{
				Modifier: "IGNORE",
			}).Create(&toInsert).Error; err != nil {
				return err
			}
		}

		// 更新
		if len(toUpdate) > 0 {
			for _, item := range toUpdate {
				updateQuery := crmDB.Table(tableName)
				for i, field := range queryFields {
					updateQuery = updateQuery.Where(fmt.Sprintf("%s = ?", field), conditions[i])
				}
				if err := updateQuery.Updates(item).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}
