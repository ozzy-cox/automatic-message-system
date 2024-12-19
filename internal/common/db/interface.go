package db

import "iter"

type MessageRepository interface {
	GetUnsentMessagesFromDb(limit, offset int) iter.Seq2[*Message, error]
	GetSentMessagesFromDb(limit, offset int) iter.Seq2[*Message, error]
	MarkMessageAsSent(messageId int) error
}
