package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type FoodItem struct {
	Name     string
	Units    string
	Calories int
	line     int
}

type FoodDB struct {
	filePath string
	items    []FoodItem
	byName   map[string]*FoodItem
}

func (db *FoodDB) reload() error {
	data, err := os.ReadFile(db.filePath)
	if err != nil {
		return err
	}

	clear(db.byName)
	clear(db.items)

	lineNum := 0
	for line := range strings.SplitSeq(string(data), "\n") {
		lineNum++
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "|")

		// TODO: Add error handling here
		if len(parts) != 3 {
			continue
		}

		cal, err := strconv.Atoi(strings.TrimSpace(parts[2]))
		if err != nil {
			continue
		}

		item := FoodItem{
			Name:     strings.TrimSpace(parts[0]),
			Units:    strings.TrimSpace(parts[1]),
			Calories: cal,
			line:     lineNum,
		}
		db.items = append(db.items, item)
		db.byName[item.Name] = &db.items[len(db.items)-1]
	}
	return nil
}

func LoadFoodDB(path string) (*FoodDB, error) {

	db := &FoodDB{
		byName:   make(map[string]*FoodItem),
		filePath: path,
	}

	// Load from file
	db.reload()

	return db, nil
}

func (db *FoodDB) Get(name string) *FoodItem {
	return db.byName[name]
}

func (db *FoodDB) All() []FoodItem {
	return db.items
}

func (db *FoodDB) Delete(i FoodItem) error {
	err := deleteLine(db.filePath, i.line)
	if err != nil {
		return err
	}
	db.reload()
	return nil
}

func (db *FoodDB) Add(i FoodItem) error {

	// Append to file
	f, err := os.OpenFile(db.filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	line := fmt.Sprintf("%s | %s | %d\n", i.Name, i.Units, i.Calories)
	if _, err := f.WriteString(line); err != nil {
		return err
	}

	// Reload
	db.reload()
	return nil
}
