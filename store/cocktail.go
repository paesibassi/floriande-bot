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
		name,
		fmt.Sprintf("1%s", name),
	}
}

type Cocktail struct {
	CocktailName string
	CocktailCode string
	Category
}

func NewCocktail(name string) Cocktail {
	return Cocktail{
		CocktailName: name,
		CocktailCode: fmt.Sprintf("2%s", name),
		Category:     Category{}, // TODO fill category
	}
}

type MenuKeyboards map[string]tgbotapi.InlineKeyboardMarkup

var EmptyInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{},
)

func newDrinksKeyboard(buttons []tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttons...),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Cancel", "0Cancel"),
		),
	)
	return keyboard
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
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData(d.CocktailName, d.CocktailCode)
	}
	return newDrinksKeyboard(buttons)
}
