package db

import "iter"

type MessageRepository interface {
	GetMessages(limit, offset int) iter.Seq2[*Message, error]
	GetSentMessages(limit, offset int) iter.Seq2[*Message, error]
	MarkMessageAsSent(messageId int) error
}
