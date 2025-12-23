package main

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

type config struct {

	// Config file
	homeFolder string

	// Loaded at beginning
	foodDB *FoodDB

	// Window width
	ww int
}

func (c *config) fullWidth() int {
	return c.ww - 8
}

// TODO: Make these configurable
func (c *config) updateWW(ww int) {
	c.ww = max(30, min(ww, 80))
}

var cfg config

func (c *config) getPath(filename string) string {
	return filepath.Join(c.homeFolder, filename)
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

func readConfig() config {
	// Get home
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	// Path to config
	configPath := filepath.Join(homeDir, ".config", "furnace", "config.ini")

	// Load the INI file
	cfg_file, err := ini.Load(configPath)
	if err != nil {
		panic(err)
	}

	// Read values
	ret := config{ww: 30}
	section := cfg_file.Section("general")
	ret.homeFolder = expandPath(section.Key("homeFolder").String())

	// Load food library
	db, err := LoadFoodDB(filepath.Join(ret.homeFolder, "food.md"))
	if err != nil {
		panic(err)
	}
	ret.foodDB = db

	return ret
}

func (c *config) saveNewFood(i FoodItem) {
	err := c.foodDB.Add(filepath.Join(c.homeFolder, "food.md"), i)
	if err != nil {
		panic(err)
	}
}
