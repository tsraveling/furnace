package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type RecipePart struct {
	Food     string
	Quantity float64
}

type Recipe struct {
	Name          string
	Parts         []RecipePart
	OtherCalories int
	Units         string
	startLine     int
	endLine       int
}

type RecipeDB struct {
	filePath string
	items    []Recipe
	byName   map[string]*Recipe
}

func (db *RecipeDB) reload() error {
	data, err := os.ReadFile(db.filePath)
	if err != nil {
		return err
	}

	clear(db.byName)
	db.items = db.items[:0]

	lines := strings.Split(string(data), "\n")
	i := 0
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])

		// Look for recipe header: # Name
		if !strings.HasPrefix(line, "# ") {
			i++
			continue
		}

		startLine := i + 1 // 1-based
		name := strings.TrimSpace(line[2:])
		i++

		// Next non-empty line: units | otherCalories
		for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
			i++
		}
		if i >= len(lines) {
			break
		}

		metaParts := strings.Split(lines[i], "|")
		if len(metaParts) != 2 {
			i++
			continue
		}

		units := strings.TrimSpace(metaParts[0])
		otherCal, err := strconv.Atoi(strings.TrimSpace(metaParts[1]))
		if err != nil {
			i++
			continue
		}
		i++

		// Parse ingredient lines: - Food | quantity
		var parts []RecipePart
		for i < len(lines) {
			line = strings.TrimSpace(lines[i])
			if !strings.HasPrefix(line, "- ") {
				break
			}

			entry := strings.TrimPrefix(line, "- ")
			fields := strings.Split(entry, "|")
			if len(fields) != 2 {
				i++
				continue
			}

			qty, err := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
			if err != nil {
				i++
				continue
			}

			parts = append(parts, RecipePart{
				Food:     strings.TrimSpace(fields[0]),
				Quantity: qty,
			})
			i++
		}

		recipe := Recipe{
			Name:          name,
			Parts:         parts,
			OtherCalories: otherCal,
			Units:         units,
			startLine:     startLine,
			endLine:       i, // 1-based exclusive
		}
		db.items = append(db.items, recipe)
		db.byName[recipe.Name] = &db.items[len(db.items)-1]
	}
	return nil
}

func LoadRecipeDB(path string) (*RecipeDB, error) {
	// Create file if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte(""), 0644); err != nil {
			return nil, err
		}
	}

	db := &RecipeDB{
		byName:   make(map[string]*Recipe),
		filePath: path,
	}

	if err := db.reload(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *RecipeDB) Get(name string) *Recipe {
	return db.byName[name]
}

func (db *RecipeDB) All() []Recipe {
	return db.items
}

func (db *RecipeDB) Add(r Recipe) error {
	f, err := os.OpenFile(db.filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n", r.Name)
	fmt.Fprintf(&b, "%s | %d\n", r.Units, r.OtherCalories)
	for _, p := range r.Parts {
		fmt.Fprintf(&b, "- %s | %g\n", p.Food, p.Quantity)
	}
	b.WriteString("\n")

	if _, err := f.WriteString(b.String()); err != nil {
		return err
	}

	db.reload()
	return nil
}

func (db *RecipeDB) Delete(r Recipe) error {
	data, err := os.ReadFile(db.filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	// Remove lines from startLine-1 to endLine-1 (0-based)
	start := r.startLine - 1
	end := r.endLine - 1
	// Also remove a trailing blank line if present
	if end < len(lines) && strings.TrimSpace(lines[end]) == "" {
		end++
	}
	lines = append(lines[:start], lines[end:]...)

	if err := os.WriteFile(db.filePath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return err
	}

	db.reload()
	return nil
}
