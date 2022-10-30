package messages

const (
	startCommand       = "/start"
	addCommand         = "/add"
	getExpensesCommand = "/get"
	helpCommand        = "/help"
)

const (
	invalidCommand   = "Невалидная команда"
	invalidMsg       = "Невалидное сообщение. Посмотрите в /help"
	invalidTimestamp = "Невалидный формат даты. Используйте формат rfc3339"
	invalidRange     = "Невалидный диапазон"
	expensesNotFound = "Трат не найдено"
	help             = "Для добавления траты введите сообщение вида:\n📎, {размер траты}, {категория}\nДля просмотра отчета за день введите:\n📤 1\nОтчет за месяц: 📤 2\nОтчет за год: 📤 3"
	addMessage       = "Добавьте трату"
)
