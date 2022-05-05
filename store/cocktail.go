package store

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Category struct {
	CategoryName string
	CategoryCode string
}

func NewCategory(name string) Category {
	return Category{
		CategoryName: name,
		CategoryCode: fmt.Sprintf("1%s", name),
	}
}

type Cocktail struct {
	CocktailName string
	CocktailCode string
	Category
}

func NewCocktail(name, categoryName string) Cocktail {
	return Cocktail{
		CocktailName: name,
		CocktailCode: fmt.Sprintf("2%s", name),
		Category:     NewCategory(categoryName),
	}
}

type MenuKeyboards map[string]tgbotapi.InlineKeyboardMarkup

var EmptyInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{},
)

func newDrinksKeyboard(buttons []tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup {
	const buttonsPerRow = 2
	nButtons := len(buttons)
	var rows [][]tgbotapi.InlineKeyboardButton
	for i := 0; i < nButtons; i += buttonsPerRow {
		limit := i + buttonsPerRow
		if limit > nButtons {
			limit = nButtons
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(buttons[i:limit]...))
	}
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Cancel", "0Cancel"),
		))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func newCategoriesKeyboard(categories []Category) tgbotapi.InlineKeyboardMarkup {
	buttons := make([]tgbotapi.InlineKeyboardButton, len(categories))
	for i, d := range categories {
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData(d.CategoryName, d.CategoryCode)
	}
	return newDrinksKeyboard(buttons)
}

func newCocktailsKeyboard(cocktails []Cocktail) tgbotapi.InlineKeyboardMarkup {
	buttons := make([]tgbotapi.InlineKeyboardButton, len(cocktails))
	for i, d := range cocktails {
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData(d.CocktailName, d.CocktailCode+d.CategoryCode)
	}
	return newDrinksKeyboard(buttons)
}
