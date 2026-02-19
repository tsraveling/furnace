package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type log struct {
	date     time.Time
	name     string
	units    string
	quantity float64
	calories int
	isRecipe bool
	line     int
}

func loadLogs() []log {
	path := cfg.getPath("logs.md")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var logs []log

	lineNum := 0
	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)
		lineNum++

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 3 && len(parts) != 4 {
			continue
		}

		date, err := time.Parse("2006-01-02", strings.TrimSpace(parts[0]))
		if err != nil {
			continue
		}

		itemName := strings.TrimSpace(parts[1])
		quantity, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
		if err != nil {
			continue
		}

		isRecipe := len(parts) == 4 && strings.TrimSpace(parts[3]) == "recipe"

		var calories int
		var units string

		if isRecipe {
			recipe := cfg.recipeDB.Get(itemName)
			if recipe == nil {
				continue
			}
			calories = int(float64(recipe.CaloriesPerUnit()) * quantity)
			units = recipe.Units
		} else {
			food := cfg.foodDB.Get(itemName)
			if food == nil {
				continue
			}
			calories = int(float64(food.Calories) * quantity)
			units = food.Units
		}

		logs = append(logs, log{
			date:     date,
			name:     itemName,
			units:    units,
			quantity: quantity,
			calories: calories,
			isRecipe: isRecipe,
			line:     lineNum,
		})
	}

	// Sort these by date so we can easily filter them
	slices.SortFunc(logs, func(a, b log) int {
		return a.date.Compare(b.date)
	})

	return logs
}

func writeLog(itemName string, quantity float64, date time.Time) error {
	path := cfg.getPath("logs.md")

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	line := fmt.Sprintf("%s | %s | %.2f\n", date.Format("2006-01-02"), itemName, quantity)
	_, err = f.WriteString(line)
	return err
}

func writeRecipeLog(recipeName string, quantity float64, date time.Time) error {
	path := cfg.getPath("logs.md")

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	line := fmt.Sprintf("%s | %s | %.2f | recipe\n", date.Format("2006-01-02"), recipeName, quantity)
	_, err = f.WriteString(line)
	return err
}
