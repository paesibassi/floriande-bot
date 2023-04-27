package store

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/api/iterator"
)

const (
	FreshEasyMenu   = "AgACAgIAAxkBAAIDw2JBrqb52-4vUPAn3OovZHZ_ZnVyAAJrtzEbEpM5STL2aSKwlnLmAQADAgADcwADIwQ"
	ConnoisseurMenu = "AgACAgIAAxkBAAIDzmJBsV_nyDgddGRDu65yjV9S9djNAAJstzEbEpM5SR2KRHwgCPQaAQADAgADcwADIwQ"
)

type Menu map[string][]string

func CocktailsMenu(client *firestore.Client) Menu {
	var cocktail Cocktail
	menu := make(map[string][]string)
	iter := client.
		Collection("menu").
		Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		err = doc.DataTo(&cocktail)
		if err != nil {
			continue
		}
		menu[cocktail.CategoryName] = append(menu[cocktail.CategoryName], cocktail.CocktailName)
	}
	return menu
}

func NewMenuList(m Menu) map[string]string {
	drinksMap := make(map[string]string)
	for category, drinks := range m {
		for _, d := range drinks {
			drinksMap[d] = category
		}
	}
	return drinksMap
}

func NewCocktailKeyboards(m Menu) (
	categoriesKeyboard tgbotapi.InlineKeyboardMarkup, cocktailkeyboards MenuKeyboards,
) {
	var categories []Category
	cocktailkeyboards = make(MenuKeyboards, len(m))
	for categoryName, cocktailNames := range m {
		category := NewCategory(categoryName)
		cocktails := make([]Cocktail, len(cocktailNames))
		for i, name := range cocktailNames {
			cocktails[i] = NewCocktail(name, "", categoryName) // FIXME image empty here?
		}
		categories = append(categories, category)
		cocktailkeyboards[category.CategoryCode] = newCocktailsKeyboard(cocktails)
	}
	categoriesKeyboard = newCategoriesKeyboard(categories)
	return categoriesKeyboard, cocktailkeyboards
}
