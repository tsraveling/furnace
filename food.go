package main

import (
	"os"
	"strconv"
	"strings"
)

type FoodItem struct {
	Name     string
	Units    string
	Calories int
}

type FoodDB struct {
	items  []FoodItem
	byName map[string]*FoodItem
}

func LoadFoodDB(path string) (*FoodDB, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	db := &FoodDB{
		byName: make(map[string]*FoodItem),
	}

	for line := range strings.SplitSeq(string(data), "\n") {
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
		}
		db.items = append(db.items, item)
		db.byName[item.Name] = &db.items[len(db.items)-1]
	}

	return db, nil
}

func (db *FoodDB) Get(name string) *FoodItem {
	return db.byName[name]
}

func (db *FoodDB) All() []FoodItem {
	return db.items
}
