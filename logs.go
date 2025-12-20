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
	item     *FoodItem
	quantity float64
	calories int
}

func loadLogs() []log {
	path := cfg.getPath("logs.md")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var logs []log

	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 3 {
			continue
		}

		date, err := time.Parse("2006-01-02", strings.TrimSpace(parts[0]))
		if err != nil {
			continue
		}

		// Format: date | itemName | quantity
		itemName := strings.TrimSpace(parts[1])
		quantity, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
		if err != nil {
			continue
		}

		item := cfg.foodDB.Get(itemName)
		if item == nil {
			continue
		}

		calories := int(float64(item.Calories) * quantity)
		logs = append(logs, log{
			date:     date,
			item:     item,
			quantity: quantity,
			calories: calories,
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

func logsForDate(logs []log, day time.Time) []log {
	y, m, d := day.Date()
	var filtered []log
	for _, l := range logs {
		ly, lm, ld := l.date.Date()
		if ly == y && lm == m && ld == d {
			filtered = append(filtered, l)
		}
	}
	return filtered
}
