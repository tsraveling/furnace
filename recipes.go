package main

type RecipePart struct {
	Food     string
	quantity float64
}

type Recipe struct {
	Parts         []RecipePart
	OtherCalories int
	Units         string
}
