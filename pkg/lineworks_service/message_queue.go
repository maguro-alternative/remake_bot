package lineworks_service

import (
	"context"
	"log/slog"
	"time"
)

// MessageQueue handles asynchronous message sending
type MessageQueue struct {
	queue   chan *MessageRequest
	service *LineWorksService // LineWorksServiceへの参照に変更
	ctx     context.Context
	cancel  context.CancelFunc
	workers int
}

// NewMessageQueue creates a new message queue
func NewMessageQueue(service *LineWorksService, bufferSize, workers int) *MessageQueue {
	return &MessageQueue{
		queue:   make(chan *MessageRequest, bufferSize),
		service: service,
		workers: workers,
	}
}

// Start starts the message queue workers
func (mq *MessageQueue) Start(ctx context.Context) {
	mq.ctx, mq.cancel = context.WithCancel(ctx)
	
	// Start worker goroutines
	for i := 0; i < mq.workers; i++ {
		go mq.worker(i)
	}
	
	slog.InfoContext(ctx, "Message queue started", "workers", mq.workers, "buffer_size", cap(mq.queue))
}

// Stop stops the message queue
func (mq *MessageQueue) Stop() {
	if mq.cancel != nil {
		mq.cancel()
	}
	
	// Drain the queue
	go func() {
		for {
			select {
			case <-mq.queue:
				// Discard remaining messages
			default:
				return
			}
		}
	}()
}

// EnqueueMessage adds a message to the queue
func (mq *MessageQueue) EnqueueMessage(ctx context.Context, guildID, message string) error {
	req := &MessageRequest{
		GuildID: guildID,
		Message: message,
		Context: ctx,
	}
	
	select {
	case mq.queue <- req:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Queue is full, send directly
		return mq.sendMessageDirect(ctx, guildID, message)
	}
}

// worker processes messages from the queue
func (mq *MessageQueue) worker(id int) {
	slog.InfoContext(mq.ctx, "Message queue worker started", "worker_id", id)
	
	for {
		select {
		case req := <-mq.queue:
			startTime := time.Now()
			err := mq.sendMessageDirect(req.Context, req.GuildID, req.Message)
			duration := time.Since(startTime)
			
			if err != nil {
				slog.ErrorContext(req.Context, "error", "log", err.Error())
				slog.ErrorContext(req.Context, "Failed to process message from queue",
					"worker_id", id,
					"guildID", req.GuildID,
					"error", err,
					"duration", duration)
			} else {
				slog.DebugContext(req.Context, "Message processed successfully",
					"worker_id", id,
					"guildID", req.GuildID,
					"duration", duration)
			}
			
		case <-mq.ctx.Done():
			slog.InfoContext(mq.ctx, "Message queue worker stopped", "worker_id", id)
			return
		}
	}
}

// sendMessageDirect sends a message directly using the service
func (mq *MessageQueue) sendMessageDirect(ctx context.Context, guildID, message string) error {
	return mq.service.SendMessageDirect(ctx, guildID, message)
}

// GetQueueStatus returns queue status information
func (mq *MessageQueue) GetQueueStatus() QueueStatus {
	return QueueStatus{
		QueueLength: len(mq.queue),
		QueueCap:    cap(mq.queue),
		Workers:     mq.workers,
	}
}

// QueueStatus represents the status of the message queue
type QueueStatus struct {
	QueueLength int `json:"queue_length"`
	QueueCap    int `json:"queue_capacity"`
	Workers     int `json:"workers"`
}