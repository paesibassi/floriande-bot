package controller

type language int

const (
	EN language = iota
	IT
)

func userLanguage(code string) language {
	switch code {
	case "it":
		return IT
	default:
		return EN
	}
}

type multiLanguageMessage int

const (
	dontUnderstand multiLanguageMessage = iota
	bookConfirmation
	orderConfirmation
	youChose
	youCanceled
	chooseDrinkInCategory
	orderReady
	welcomeText
	helpText
	chooseEvent
	hereTheMenu
	chooseCategory
	yourOrders
	dontKnowCommand
	reservedForBarman
)

type messagesMap map[multiLanguageMessage][2]string

var mss messagesMap = messagesMap{
	dontUnderstand: {
		`I don't understand: '%s'\.
Please try the command ` + "`/drink`" + ` to get guided help while ordering\.
Otherwise, make sure to type the exact name of the drink \(including capitals\)\.
Get in touch with [us](tg://user?id=%d) if anything doesn't work\!`,
		`Non capisco: '%s'\.
Prova il comando ` + "`/drink`" + ` per ricevere aiuto guidato per fare un ordine\.
Altrimenti, presta attenzione ad inserire correttamente il nome del drink \(incluse maiuscole\)\.
[Contattaci](tg://user?id=%d) direttamente se qualcosa non funziona\!`,
	},
	bookConfirmation: {
		"Great! We reserved a spot for you %s!",
		"Ottimo! Abbiamo riservato un posto per te %s!",
	},
	orderConfirmation: {
		"A %s is coming soon! %v",
		"Stiamo preparando un %s per te! %v",
	},
	youChose: {
		"You selected %s",
		"Hai scelto %s",
	},
	youCanceled: {
		"You canceled your order",
		"Hai cancellato il tuo ordine",
	},
	chooseDrinkInCategory: {
		"2. good choice! Now choose your favorite drink in this category",
		"2. ottima scelta! Adesso scegli il tuo drink preferito in questa categoria",
	},
	orderReady: {
		"Your order %s is ready! Enjoy!",
		"Il tuo ordine per un %s è pronto! Enjoy!",
	},
	welcomeText: {
		`Welcome by Floriande Lounge bar %v\.

Have you already reserved a spot for an upcoming event? You can do so with the ` +
			"*`/book`*" + ` command\.
Please use the ` + "*`/menu`*" + ` command to download our latest drink selection\. %v
You can order a drink from here using the ` + "`/drink`" + ` command, or you can type ` +
			`the name of the cocktail if you know it already\. Make sure you spell it correctly\!
Then, you can check if you have any order waiting to be prepared and served` +
			`with the ` + "`/orders`" + ` command\. %v
Please let [us](tg://user?id=%d) know if you have suggestions for improvement\.
We hope you enjoy you stay\! %v

What would you like to drink today? %v`,
		`Benvenuti al Floriande Lounge bar %v\.

Hai già prenotato un posto per un evento in programma? Puoi farlo con il comando` +
			"*`/book`*" + ` \.
Usa il comando ` + "*`/menu`*" + ` per scaricare il menu delle nostre selezioni di drink\. %v
Puoi ordinare un drink usando il comando ` + "`/drink`" + ` , oppure puoi digitare direttamente ` +
			`il nome del cocktail se lo conosci già\. Fai attenzione a scriverlo correttamente\!
Quindi, puoi controllare se hai qualche ordine in attesa di essere preparato e servito` +
			`con il comando ` + "`/orders`" + ` \. %v
[Contattaci](tg://user?id=%d) direttamente se hai qualche suggerimento\.
Buon divertimento\! %v

Di cosa hai voglia oggi? %v`,
	},
	helpText: {
		"You can use the *`/book`* command to reserve for an event, " +
			"*`/menu`* to download the digital version of our cocktail menu, " +
			"*`/drink`* to order a cocktail" + ` \(you will be guided through the process\), ` +
			"*`/orders`*" + ` to see the cocktail\(s\) you have ordered and are being mixed\.
The commands ` + "*`/list`* and *`/serve`*," + ` are reserved for the barman\.`,
		"Puoi usare il comando *`/book`* per prenotare un posto per un evento, " +
			"*`/menu`* per scaricare la versione digitale del nostro menu di cocktails, " +
			"*`/drink`* per ordinare un cocktail" + ` \(sarai guidato nella scelta\), ` +
			"*`/orders`*" + ` per controllare gli ordini in attesa di essere preparati\.
I comandi ` + "*`/list`* e *`/serve`*" + ` sono riservati per il barman\.`,
	},
	chooseEvent: {
		"We are happy you want to reserve your place with us. Which event would you like to join?",
		"Siamo lieti che tu voglia prenotare un posto con noi. A quale evento vorresti partecipare?",
	},
	hereTheMenu: {
		"Here you go, our digital menu. What would you like to drink?",
		"Ecco, il nostro menu digitale. Cosa ti piacerebbe ordinare?",
	},
	chooseCategory: {
		"1. let's start with choosing a category %v; what kind of cocktail would you like?",
		"1. cominciamo con lo scegliere una categoria %v; che tipo di cocktail ti piace?",
	},
	yourOrders: {
		"Your current order(s):\n%s",
		"I/l tuo ordine/i attuale/i:\n%s",
	},
	dontKnowCommand: {
		"I don't know that command",
		"Non conosco quel comando",
	},
	reservedForBarman: {
		"The `%v` command is reserved for the barman",
		"Il comando `%v` è riservato per il barman",
	},
}
