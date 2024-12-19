package producer

import (
	"context"
	"sync"
	"time"
)

func (s *Service) ProduceMessages(ctx context.Context, wg *sync.WaitGroup) {
	s.ProducerOnStatus.Store(true)
	ticker := time.NewTicker(s.Config.Interval)
	defer ticker.Stop()

	offset := s.mustGetProducerOffset()
	poffset := &offset

	for {
		select {
		case <-ctx.Done():
			s.Logger.Println("Context canceled, saving final offset")
			s.mustSetProducerOffset(poffset)
			wg.Done()
			return
		case <-ticker.C:
			if !s.ProducerOnStatus.Load() {
				continue
			}
			limit := s.Config.BatchCount
			go s.PushMessagesToQ(ctx, limit, offset)
			(*poffset) += limit
		}
	}
}
