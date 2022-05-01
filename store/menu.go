package store

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	FreshEasyMenu   = "AgACAgIAAxkBAAIDw2JBrqb52-4vUPAn3OovZHZ_ZnVyAAJrtzEbEpM5STL2aSKwlnLmAQADAgADcwADIwQ"
	ConnoisseurMenu = "AgACAgIAAxkBAAIDzmJBsV_nyDgddGRDu65yjV9S9djNAAJstzEbEpM5SR2KRHwgCPQaAQADAgADcwADIwQ"
)

type Menu map[string][]string

var menu = Menu{
	"Entree": {
		"Caipiroska Fragola",
		"Caipiroska Mango",
		"Mojito",
		"Virgin Mojito",
		"Spritz Aperol",
		"Spritz Campari",
	},
	// "Fresh&Easy": {
	// 	"Americano",
	// 	"Bay Breeze",
	// 	"Caipiroska Fragola",
	// 	"Caipiroska Mango",
	// 	"Cosmopolitan",
	// 	"Cuba Libre",
	// 	"Dark & Stormy",
	// 	"Gin Tonic",
	// 	"Margarita",
	// 	"Virgin Margarita",
	// 	"Mojito",
	// 	"Virgin Mojito",
	// 	"Moscow Mule",
	// 	"Spritz Aperol",
	// 	"Spritz Campari",
	// 	"Tequila Sunrise",
	// },
	// "Connoisseur": {
	// 	"Daiquiri",
	// 	"Dry Martini",
	// 	"God Father",
	// 	"Long Island Ice Tea",
	// 	"Manhattan",
	// 	"Negroni",
	// 	"Old Fashioned",
	// 	"Sex on the Beach",
	// },
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

var AllCocktails map[string]string = NewMenuList(menu)

func NewCocktailKeyboards(m Menu) (
	categoriesKeyboard tgbotapi.InlineKeyboardMarkup, cocktailkeyboards MenuKeyboards,
) {
	var categories []Category
	cocktailkeyboards = make(MenuKeyboards, len(m))
	for categoryName, cocktailNames := range m {
		category := NewCategory(categoryName)
		cocktails := make([]Cocktail, len(cocktailNames))
		for i, name := range cocktailNames {
			cocktails[i] = NewCocktail(name, categoryName)
		}
		categories = append(categories, category)
		cocktailkeyboards[category.CategoryCode] = newCocktailsKeyboard(cocktails)
	}
	categoriesKeyboard = newCategoriesKeyboard(categories)
	return categoriesKeyboard, cocktailkeyboards
}

var CategoriesKeyboard, CocktailKeyboards = NewCocktailKeyboards(menu)
