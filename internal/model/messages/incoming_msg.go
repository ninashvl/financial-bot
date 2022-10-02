package messages

import (
	datahandler "gitlab.ozon.dev/lukyantsev-pa/lukyantsev-pavel/internal/model/handler"
)

type MessageSender interface {
	SendMessage(text string, userID int64) error
}

type Model struct {
	tgClient MessageSender
}

func New(tgClient MessageSender) *Model {
	return &Model{
		tgClient: tgClient,
	}
}

type Message struct {
	Text   string
	UserID int64
}

type ReportChek struct {
	check bool
	code  int
}

var addSpendChek bool = false
var reportChek ReportChek

func (s *Model) IncomingMessage(msg Message) error {
	if msg.Text == comandStart {
		return s.tgClient.SendMessage(textHello+textCommand+textComndAdd+textComndAddDesc+textComndReport+textComndReportDesc+textComndHelp, msg.UserID)
	} else if msg.Text == comandHelp {
		return s.tgClient.SendMessage(textCommand+textComndAdd+textComndAddDesc+textComndReport+textComndReportDesc+textComndHelp, msg.UserID)
	} else if msg.Text == comandAdd {
		addSpendChek = true
		return s.tgClient.SendMessage(comandAddDo, msg.UserID)
	} else if msg.Text == comandReport {
		reportChek.check = true
		switch msg.Text {
		case "/week":
			reportChek.code = 1
		case "/month":
			reportChek.code = 2
		case "year":
			reportChek.code = 3
		}
		return s.tgClient.SendMessage(comandReportDo, msg.UserID)
	}

	if addSpendChek {
		addSpendChek = false
		err := datahandler.AddSpend(msg.Text)
		if err == nil {
			return s.tgClient.SendMessage("", msg.UserID)
		} else {
			return s.tgClient.SendMessage(errBadSpendFormat, msg.UserID)
		}
	}

	if reportChek.check {
		reportChek.check = false
		result := datahandler.Report(reportChek.code)
		return s.tgClient.SendMessage(result, msg.UserID)
	}

	return s.tgClient.SendMessage(textErrorComnd, msg.UserID)
}
