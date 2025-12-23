# furnace

A terminal application for tracking what you eat, as frictionlessly as possible. Add food quickly, look up the nutrition later, and set easy autocomplete meals and snacks for easy logging in the future.

## Installation and Setup

1. While this is in alpha, install locally by running `go install` from the repo root.
2. Copy this into `.config/furnace/config.ini` (feel free to change the path to whereever you'd like your food and log files to go):

```
[general]
homeFolder = "~/.config/furnace/"
```

That's it!

## Usage

You can see a **summary view** by simply typing `furnace`. This will show your logs for today and a calorie total. Follow the instructions in the help text to page through days, add items, etc.

You can log a new item by typing `furnace log`. This will drop you into the picker flow. Follow the instructions in the help text to log food, create new food items, and so on.

### Editing Logs and Food Items

For now, you cannot edit logs or food items from inside Furnace (although this is coming soon!). So instead simply use the text editor of your choice to edit either `logs.md` or `food.md` in the home folder you set above.

- Logs has the format `date | item | quantity` (e.g. of servings)
- Food has the format `item | units | calories`

Then your total daily caloric intake is outputted as `items today * quantity in units * calories per unit`. Easy!

## Development

### Prerequisites

- Go 1.25+
- [cocogitto](https://github.com/cocogitto/cocogitto) for commit linting (`brew install cocogitto`)

### Commit conventions

This project uses [Conventional Commits](https://www.conventionalcommits.org/). The git hooks will validate your commit messages.

Run this to use the comitted githooks directory:

```bash
git config core.hooksPath .githooks
```

Use conventional commit syntax when making commits or PRs:

```
feat: add new feature
fix: fix a bug
docs: update documentation
chore: maintenance tasks
```

## License

[GNU GPL v3](LICENSE)



