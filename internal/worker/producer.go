package worker

import (
	"context"
	"sync"
	"time"

	"github.com/ozzy-cox/automatic-message-system/internal/db"
)

const limit = 2

func ProduceMessages(wg *sync.WaitGroup, ctx context.Context, chanMessage chan db.Message, ticker *time.Ticker) {
	offset := mustGetProducerOffset()
	poffset := &offset
	for {
		select {
		case <-ctx.Done():
			mustSetProducerOffset(poffset)
			wg.Done()
			return
		case <-ticker.C:
			if !ProducerOnStatus.Load() {
				continue
			}
			rows, err := db.DbConnection.Query("SELECT * FROM messages LIMIT $1 OFFSET $2", limit, offset)
			if err != nil {
			}

			for rows.Next() {
				var msg db.Message
				err := rows.Scan(
					&msg.ID,
					&msg.Content,
					&msg.To,
					&msg.Sent,
					&msg.SentAt,
					&msg.CreatedAt,
				)
				if err != nil {
					panic("Failed to scan messages")
				}
				chanMessage <- msg
				(*poffset)++
			}
		}
	}
}
