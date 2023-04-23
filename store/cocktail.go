package store

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	CocktailName  string
	CocktailCode  string
	CocktailImage string
	Category
}

func NewCocktail(name, image, categoryName string) Cocktail {
	return Cocktail{
		CocktailName:  name,
		CocktailCode:  fmt.Sprintf("2%s", name),
		CocktailImage: image,
		Category:      NewCategory(categoryName),
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

func CocktailImageFileID(client *firestore.Client, cocktailName string) (*string, error) {
	cocktailPtr := new(Cocktail)
	dsnap, err := client.
		Collection("menu").
		Doc(cocktailName).
		Get(context.Background())
	if err != nil {
		// if document is not found
		if status.Code(err) == codes.NotFound {
			err = fmt.Errorf("cocktail %q not found: %s", cocktailName, err)
			return nil, err
		}
		// fallback in case of a different error
		err = fmt.Errorf("could not fetch cocktail details: %s", err)
		return nil, err
	}

	err = dsnap.DataTo(cocktailPtr)
	if err != nil {
		err = fmt.Errorf("could not unmarshal cocktail details: %s", err)
		return nil, err
	}

	fileId := cocktailPtr.CocktailImage
	if fileId == "" {
		return nil, fmt.Errorf("cocktail %q has no image defined", cocktailName)
	}
	return &fileId, nil
}
