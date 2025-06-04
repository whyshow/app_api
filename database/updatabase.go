package database

import (
	"fmt"
	"gorm.io/gorm"
)

func UpgradeDatabase(models []interface{}) error {
	for _, model := range models {
		if err := processModel(model); err != nil {
			return fmt.Errorf("数据库升级失败: %v", err)
		}
	}
	return nil
}

func processModel(model interface{}) error {
	mig := DB.Migrator()

	// 检查表是否存在
	if !mig.HasTable(model) {
		if err := mig.CreateTable(model); err != nil {
			return err
		}
		return nil
	}

	// 获取模型结构信息
	stmt := &gorm.Statement{DB: DB}
	if err := stmt.Parse(model); err != nil {
		return err
	}

	// 获取现有字段
	columns, err := mig.ColumnTypes(model)
	if err != nil {
		return err
	}

	// 检查并添加缺失字段
	for _, field := range stmt.Schema.Fields {
		if field.DBName == "" || field.IgnoreMigration {
			continue
		}

		if !columnExists(columns, field.DBName) {
			if err := mig.AddColumn(model, field.DBName); err != nil {
				return fmt.Errorf("添加字段 %s 失败: %v", field.DBName, err)
			}
		}
	}

	return nil
}

func columnExists(columns []gorm.ColumnType, name string) bool {
	for _, col := range columns {
		if col.Name() == name {
			return true
		}
	}
	return false
}
