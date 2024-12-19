package producer_test

import (
	"context"
	"iter"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/ozzy-cox/automatic-message-system/internal/common/db"
	"github.com/ozzy-cox/automatic-message-system/internal/common/logger"
	"github.com/ozzy-cox/automatic-message-system/internal/common/queue"
	"github.com/ozzy-cox/automatic-message-system/internal/producer"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) GetUnsentMessages(limit, offset int) iter.Seq2[*db.Message, error] {
	args := m.Called(limit, offset)
	return args.Get(0).(iter.Seq2[*db.Message, error])
}

func (m *MockMessageRepository) GetSentMessages(limit, offset int) iter.Seq2[*db.Message, error] {
	args := m.Called(limit, offset)
	return args.Get(0).(iter.Seq2[*db.Message, error])
}

func (m *MockMessageRepository) MarkMessageAsSent(messageId int) error {
	args := m.Called(messageId)
	return args.Error(0)
}

type MockQueue struct {
	mock.Mock
	messagesChan chan queue.MessagePayload
}

func NewMockQueue() *MockQueue {
	return &MockQueue{
		messagesChan: make(chan queue.MessagePayload, 100),
	}
}

func (m *MockQueue) WriteMessages(ctx context.Context, msgs ...queue.MessagePayload) error {
	args := m.Called(ctx, msgs)
	for _, msg := range msgs {
		select {
		case m.messagesChan <- msg:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return args.Error(0)
}

func (m *MockQueue) Close() error {
	args := m.Called()
	close(m.messagesChan)
	return args.Error(0)
}

func setupTestService(t *testing.T) (*producer.Service, *MockMessageRepository, *MockQueue, *miniredis.Miniredis) {
	mockRepo := new(MockMessageRepository)
	mockQueue := NewMockQueue()
	mockRedis, err := miniredis.Run()
	mockRedis.FlushAll()
	require.NoError(t, err)

	cfg := &producer.ProducerConfig{
		Interval:   50 * time.Millisecond,
		BatchCount: 2,
	}

	loggerInst, _ := logger.NewLogger(logger.Config{LogToStdout: false})

	service := producer.NewProducerService(
		cfg,
		redis.NewClient(&redis.Options{
			Addr: mockRedis.Addr(),
		}),
		mockRepo,
		mockQueue,
		loggerInst,
	)

	return service, mockRepo, mockQueue, mockRedis
}

func sliceToIter[T any](sl []*T) iter.Seq2[*T, error] {
	var msq iter.Seq2[*T, error]
	msq = func(yield func(*T, error) bool) {
		for _, msg := range sl {
			if !yield(msg, nil) {
				return
			}
		}
	}
	return msq
}

func drainMessagesForCount(channel chan queue.MessagePayload, countInMs int) int {
	messagesCount := 0
	timeout := time.After(time.Duration(countInMs) * time.Millisecond)

	for {
		select {
		case <-channel:
			messagesCount++
		case <-timeout:
			return messagesCount
		}
	}
}

func TestProduceMessagesAsync(t *testing.T) {
	service, mockRepo, mockQueue, _ := setupTestService(t)

	now := time.Now()
	messages := []*db.Message{
		{ID: 1, Content: "Test1", To: "123", CreatedAt: &now},
		{ID: 2, Content: "Test2", To: "456", CreatedAt: &now},
	}

	mockRepo.On("GetUnsentMessages", mock.Anything, mock.Anything).Return(sliceToIter(messages))
	mockQueue.On("WriteMessages", mock.Anything, mock.Anything).Return(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	ticker := time.NewTicker(service.Config.Interval)
	defer ticker.Stop()
	service.ProducerOnStatus.Store(true)
	go service.ProduceMessages(ctx, &wg, ticker.C)

	receivedMessages := make([]queue.MessagePayload, 0)
	timeout := time.After(150 * time.Millisecond)

	for i := 0; i < len(messages); i++ {
		select {
		case msg := <-mockQueue.messagesChan:
			receivedMessages = append(receivedMessages, msg)
		case <-timeout:
			t.Fatal("Timeout waiting for messages")
		}
	}

	cancel()
	wg.Wait()

	assert.Len(t, receivedMessages, len(messages))
	assert.Equal(t, messages[0].ID, receivedMessages[0].ID)
	assert.Equal(t, messages[1].ID, receivedMessages[1].ID)
}

func TestProducerToggle(t *testing.T) {
	service, mockRepo, mockQueue, _ := setupTestService(t)

	now := time.Now()
	messages := []*db.Message{
		{ID: 1, Content: "Test1", To: "123", CreatedAt: &now},
		{ID: 2, Content: "Test2", To: "234", CreatedAt: &now},
	}

	mockRepo.On("GetUnsentMessages", mock.Anything, mock.Anything).Return(sliceToIter(messages))
	mockQueue.On("WriteMessages", mock.Anything, mock.Anything).Return(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(1)

	ticker := make(chan time.Time, 1)

	service.ProducerOnStatus.Store(true)
	go service.ProduceMessages(ctx, &wg, ticker)

	ticker <- time.Now()
	time.Sleep(100 * time.Millisecond)

	messagesCount := drainMessagesForCount(mockQueue.messagesChan, 200)
	assert.Greater(t, messagesCount, 0, "Should have received messages when producer was on")

	// Turn the producer off, simulate time passing by pushing ticker
	service.ProducerOnStatus.Store(false)
	ticker <- time.Now()
	messagesCount = drainMessagesForCount(mockQueue.messagesChan, 300)
	assert.Equal(t, messagesCount, 0, "No messages when producer is off")

	service.ProducerOnStatus.Store(true)
	time.Sleep(100 * time.Millisecond)

	ticker <- time.Now()
	time.Sleep(100 * time.Millisecond)

	messagesCount = drainMessagesForCount(mockQueue.messagesChan, 200)
	assert.Greater(t, messagesCount, 0, "Should have received messages when producer was on")

	cancel()
	wg.Wait()
}

func TestEachMessageProducedOnce(t *testing.T) {
	service, mockRepo, mockQueue, mockRedis := setupTestService(t)

	now := time.Now()
	messages := []*db.Message{
		{ID: 1, Content: "Test1", To: "123", CreatedAt: &now},
		{ID: 2, Content: "Test2", To: "234", CreatedAt: &now},
	}

	mockRepo.On("GetUnsentMessages", mock.Anything, mock.Anything).Return(sliceToIter(messages))
	mockQueue.On("WriteMessages", mock.Anything, mock.Anything).Return(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(1)

	ticker := make(chan time.Time, 1)
	ticker <- time.Now()

	service.ProducerOnStatus.Store(true)
	go service.ProduceMessages(ctx, &wg, ticker)
	time.Sleep(200 * time.Millisecond)
	mockRepo.AssertCalled(t, "GetUnsentMessages", 2, 0)

	cancel()
	wg.Wait()

	ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
	wg.Add(1)
	ticker <- time.Now()
	service.ProducerOnStatus.Store(true)
	go service.ProduceMessages(ctx, &wg, ticker)
	time.Sleep(200 * time.Millisecond)

	mockRepo.AssertCalled(t, "GetUnsentMessages", 2, 2)
	mockRepo.AssertNumberOfCalls(t, "GetUnsentMessages", 2)

	cancel()
	wg.Wait()
}
