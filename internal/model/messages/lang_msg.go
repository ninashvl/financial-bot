package messages

const (
	textHello             string = "Hello!\n"
	textCommand           string = "Команды:\n"
	textComndAdd          string = "Добавить трату: /spend summ/category/date\n"
	textComndAddDesc      string = "\t\t\t\t- сумма должна быть положительная\n\t\t\t\t- поле категории не должно быть пустым\n\t\t\t\t- дата в формате year-month-day\n"
	textComndReport       string = "Отобразить отчет: /report\n"
	textComndReportDesc   string = "\t\t\t\t- отчет за неделю: week\n\t\t\t\t- отчет за месяц: month\n\t\t\t\t- отчет за год: year\n"
	textComndHelp         string = "Справка: /help"
	textErrorComnd        string = "Не знаю эту команду"

	comandStart    string = "/start"
	comandAdd      string = "/spend"
	comandAddDo    string = "введите трату"
	comandReport   string = "/report"
	comandReportDo string = "введите период: /week, /month, /year"
	comandHelp     string = "/help"

	errBadSpendFormat string = "bad spend format"
)
