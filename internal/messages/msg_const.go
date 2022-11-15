package messages

const (
	startCommand          = "/start"
	helpCommand           = "/help"
	addCommand            = "/add"
	getExpensesCommand    = "/get"
	chooseCurrencyCommand = "/choose_currency"
	addLimit              = "/add_limit"
)

const (
	invalidCommand   = "Невалидная команда"
	invalidMsg       = "Невалидное сообщение. Посмотрите в /help"
	invalidTimestamp = "Невалидный формат даты. Используйте формат rfc3339"
	invalidRange     = "Невалидный диапазон"
	expensesNotFound = "Трат не найдено"
	help             = "Описание команд:\n\n/add\nДобавьте трату в таком формате: \n{размер траты}, {категория}. " +
		"Например, '200, магазин'.\nДля добавления траты на определенную дату введите данные в таком формате: \n{размер траты}, {категория}, {yy-mm-dd}. Например, '200, Магазин, 2021-06-02'. \n\n" +
		"/get\nПозволяет получить отчет по тратам за день, месяц или год\n\n/choose_currency\nПозволяет выбрать валюту трат. Изначально все траты отображаются в рублях."
	addMessage         = "Добавьте трату"
	savedMsg           = "Сохранено"
	currencySaved      = "Валюта трат сохранена"
	invalidCurrency    = "Невалидная валюта"
	limitSuccessfulSet = "Лимит установлен"
	inputLimit         = "Введите лимит на траты в рублях на месяц"
)
