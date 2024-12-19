package queue

import "context"

type WriterClient interface {
	WriteMessages(c context.Context, m ...MessagePayload) error
	Close() error
}
type ReaderClient interface {
	ReadMessage(c context.Context) (MessagePayload, error)
	Close() error
}
