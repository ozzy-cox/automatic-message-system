package producer

import (
	"context"
	"sync"
	"time"
)

func (s *Service) ProduceMessages(ctx context.Context, wg *sync.WaitGroup, ticker <-chan time.Time) {
	offset := s.mustGetProducerOffset()
	poffset := &offset

	for {
		if !s.ProducerOnStatus.Load() {
			continue
		}
		limit := s.Config.BatchCount
		count := s.PushMessagesToQ(ctx, limit, offset)
		(*poffset) += count
		select {
		case <-ticker:
			continue
		case <-ctx.Done():
			s.Logger.Println("Context canceled, saving final offset")
			s.mustSetProducerOffset(poffset)
			wg.Done()
			return
		}
	}
}
