package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

type config struct {

	// Config file
	homeFolder  string
	dailyTarget int

	// Loaded at beginning
	foodDB   *FoodDB
	recipeDB *RecipeDB

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

const defaultHomeFolder = "~/furnace"
const defaultDailyTarget = 2000

var defaultConfigFile = fmt.Sprintf(`[general]
homeFolder = "%s"
dailyTarget = %d
`, defaultHomeFolder, defaultDailyTarget)

// offers to create a default config; returns false if the user declines
func promptCreateConfig(configPath string) bool {
	fmt.Printf("Config file not found at %s\n", configPath)
	fmt.Println("Create it with defaults?")
	fmt.Printf("  homeFolder:  %s\n", defaultHomeFolder)
	fmt.Printf("  dailyTarget: %d\n", defaultDailyTarget)
	fmt.Print("[y/N]: ")

	answer, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	answer = strings.ToLower(strings.TrimSpace(answer))
	if answer != "y" && answer != "yes" {
		return false
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(configPath, []byte(defaultConfigFile), 0644); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(expandPath(defaultHomeFolder), 0755); err != nil {
		panic(err)
	}
	fmt.Printf("Created %s\n", configPath)
	return true
}

func readConfig() config {
	// Get home
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	// Path to config
	configPath := filepath.Join(homeDir, ".config", "furnace", "config.ini")

	// Offer to create the config on first run
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if !promptCreateConfig(configPath) {
			fmt.Println("No config created. Exiting.")
			os.Exit(0)
		}
	}

	// Load the INI file
	cfg_file, err := ini.Load(configPath)
	if err != nil {
		panic(err)
	}

	// Read values
	ret := config{ww: 30}
	section := cfg_file.Section("general")
	ret.homeFolder = expandPath(section.Key("homeFolder").String())
	ret.dailyTarget = section.Key("dailyTarget").MustInt(0) // Default to 0 (no goal)

	// Load food library
	db, err := LoadFoodDB(filepath.Join(ret.homeFolder, "food.md"))
	if err != nil {
		panic(err)
	}
	ret.foodDB = db

	// Load recipe library
	rdb, err := LoadRecipeDB(filepath.Join(ret.homeFolder, "recipes.md"))
	if err != nil {
		panic(err)
	}
	ret.recipeDB = rdb

	return ret
}
