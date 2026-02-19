# Files

## Foods

TBD

## Recipes

`recipes.md` in your home folder is filled with records that look like this:

```
# Jamaican Jerk Chicken
bowls | 150
- Chicken | 1
- Sugar | 2
- Potatoes | 4.2
- Flour | 0.2
```

In this example, the name of the recipe is Jamaican Jerk Chicken. The unit is bowls. 150 is the amount of "other calories" per unit, aka a shorthand for all of the stuff you don't want to get granular on or reference other foods for.

Then below that is a hyphen bullet list. Each item before the vertical bar is a string that keys to the `foods.md` file (see above); the number to the right of the bar is the number of units of that food per unit of this recipe. So in this case if chicken were in lbs, sugar were in tbsp, and potatoes were measured in cups, each "bowl" here would have 1 lb of chicken, 2 tbsp of sugar, and 4.2 cups of potatoes. The total caloric value of the bowl will be the calories of those ingredients, plus the 150 "other calories" item.
