package main

import (
	"fmt"
	"os"
	"time"
)

type log struct {
	item     *FoodItem
	quantity float64
	calories int
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
