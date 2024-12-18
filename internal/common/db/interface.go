package db

import "iter"

type IMessageRepository interface {
	GetUnsentMessagesFromDb(limit, offset int) iter.Seq2[*Message, error]
	GetSentMessagesFromDb(limit, offset int) iter.Seq2[*Message, error]
	SetMessageSent(messageId int) error
}
