// Package lineworks_service provides background service for LINE Works SDK integration
package lineworks_service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	sdklineworks "github.com/maguro-alternative/line-works-sdk-go/pkg/lineworks"
)

// LineWorksService manages LINE Works SDK clients for multiple guilds
type LineWorksService struct {
	configMap     map[string]*LineWorksConfig // for backward compatibility
	unifiedConfig *LineWorksConfig            // unified config for all guilds
	clientPool    *ClientPool
	messageQueue  *MessageQueue
	ctx           context.Context
	cancel        context.CancelFunc
}

// LineWorksConfig contains LINE Works configuration for a guild
type LineWorksConfig struct {
	WorksID   string
	Password  string
	ChannelNo int
}

// NewService creates a new LINE Works service with unified configuration
func NewService(worksID, password string, channelNo int) (*LineWorksService, error) {
	if worksID == "" || password == "" || channelNo == 0 {
		// Return disabled service
		return &LineWorksService{
			configMap:     make(map[string]*LineWorksConfig),
			clientPool:    NewClientPool(),
			unifiedConfig: nil,
		}, nil
	}

	// Create unified config for all guilds
	unifiedConfig := &LineWorksConfig{
		WorksID:   worksID,
		Password:  password,
		ChannelNo: channelNo,
	}

	service := &LineWorksService{
		configMap:     map[string]*LineWorksConfig{"unified": unifiedConfig},
		clientPool:    NewClientPool(),
		unifiedConfig: unifiedConfig,
	}
	
	// Create message queue with reference to the service
	service.messageQueue = NewMessageQueue(service, 100, 3)

	return service, nil
}

// Start starts the LINE Works service
func (s *LineWorksService) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)

	// If no unified config, service is disabled
	if s.unifiedConfig == nil {
		slog.InfoContext(ctx, "LINE Works service started in disabled mode (no configuration)")
		return nil
	}

	// Start client pool with unified configuration
	if err := s.clientPool.Start(ctx, s.configMap); err != nil {
		return fmt.Errorf("failed to start client pool: %w", err)
	}

	// Start message queue
	s.messageQueue.Start(ctx)

	// Pre-warm client immediately on startup for fast response
	if err := s.prewarmClientBlocking(ctx); err != nil {
		return fmt.Errorf("failed to prewarm LINE Works client: %w", err)
	}

	slog.InfoContext(ctx, "LINE Works service started with unified configuration", 
		"worksID", s.unifiedConfig.WorksID, "channelNo", s.unifiedConfig.ChannelNo)
	return nil
}

// Stop stops the LINE Works service
func (s *LineWorksService) Stop() {
	if s.cancel != nil {
		s.cancel()
	}

	// Stop message queue
	if s.messageQueue != nil {
		s.messageQueue.Stop()
	}

	// Stop client pool
	s.clientPool.Stop()

	// Wait for graceful shutdown
	time.Sleep(1 * time.Second)

	slog.Info("LINE Works service stopped")
}

// SendMessage sends a message to LINE Works using unified configuration
func (s *LineWorksService) SendMessage(ctx context.Context, guildID, message string) error {
	// Check if service is enabled
	if s.unifiedConfig == nil {
		return nil // Service is disabled, silently skip
	}
	
	// Try to send via queue first
	return s.messageQueue.EnqueueMessage(ctx, guildID, message)
}

// SendMessageDirect sends message directly to LINE Works (bypassing queue)
func (s *LineWorksService) SendMessageDirect(ctx context.Context, guildID, message string) error {
	// Check if service is enabled
	if s.unifiedConfig == nil {
		return nil // Service is disabled, silently skip
	}

	// Use unified config for any guild
	client, err := s.clientPool.GetClient(ctx, "unified", s.unifiedConfig)
	if err != nil {
		return fmt.Errorf("failed to get unified LINE Works client: %w", err)
	}

	_, err = client.SendTextMessage(ctx, s.unifiedConfig.ChannelNo, message, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to send LINE Works message",
			"guildID", guildID, "error", err)
		return err
	}

	slog.DebugContext(ctx, "LINE Works message sent successfully",
		"guildID", guildID, "channelNo", s.unifiedConfig.ChannelNo, "messageLength", len(message))

	return nil
}

// prewarmClientBlocking preloads the unified client during startup (blocking)
func (s *LineWorksService) prewarmClientBlocking(ctx context.Context) error {
	if s.unifiedConfig == nil {
		return nil
	}

	slog.InfoContext(ctx, "Pre-warming LINE Works client...", "worksID", s.unifiedConfig.WorksID)

	client, err := s.clientPool.GetClient(ctx, "unified", s.unifiedConfig)
	if err != nil {
		return fmt.Errorf("failed to create and login LINE Works client: %w", err)
	}

	// Test the client by getting user info
	myInfo := client.GetMyInfo()
	if myInfo == nil {
		return fmt.Errorf("LINE Works client test failed - unable to get user info")
	}

	slog.InfoContext(ctx, "LINE Works client pre-warmed and logged in successfully",
		"worksID", s.unifiedConfig.WorksID,
		"userName", myInfo.GetDisplayName())

	return nil
}

// prewarmClient preloads the unified client to improve first message latency (async)
func (s *LineWorksService) prewarmClient(ctx context.Context) {
	if err := s.prewarmClientBlocking(ctx); err != nil {
		slog.WarnContext(ctx, "Failed to prewarm unified client", "error", err)
	}
}

// GetServiceStatus returns comprehensive status information
func (s *LineWorksService) GetServiceStatus() ServiceStatus {
	status := ServiceStatus{
		Clients: s.clientPool.GetAllClients(),
	}
	
	if s.messageQueue != nil {
		status.Queue = s.messageQueue.GetQueueStatus()
	}
	
	return status
}

// RemoveClient removes a client from the service (for administrative purposes)
func (s *LineWorksService) RemoveClient(clientID string) {
	s.clientPool.RemoveClient(clientID)
}

// IsEnabled returns true if the service is enabled (has configuration)
func (s *LineWorksService) IsEnabled() bool {
	return s.unifiedConfig != nil
}

// ClientInfo represents LINE Works client information for a guild
type ClientInfo struct {
	Client   *sdklineworks.Client
	Config   *LineWorksConfig
	IsActive bool
	LastUsed time.Time
	mu       sync.RWMutex
}

// MessageRequest represents a message to be sent
type MessageRequest struct {
	GuildID string
	Message string
	Context context.Context
}

// ClientStatus represents the status of a LINE Works client
type ClientStatus struct {
	IsActive  bool      `json:"is_active"`
	LastUsed  time.Time `json:"last_used"`
	WorksID   string    `json:"works_id"`
	ChannelNo int       `json:"channel_no"`
}

// ServiceStatus represents the overall status of the LINE Works service
type ServiceStatus struct {
	Clients map[string]ClientStatus `json:"clients"`
	Queue   QueueStatus             `json:"queue"`
}

type LineWorksServiceMock struct{
	StartFunc func(ctx context.Context) error
	StopFunc  func()
	SendMessageFunc func(ctx context.Context, guildID, message string) error
	SendMessageDirectFunc func(ctx context.Context, guildID, message string) error
	GetServiceStatusFunc func() ServiceStatus
	RemoveClientFunc func(guildID string)
}

func (s *LineWorksServiceMock) Start(ctx context.Context) error {
	return s.StartFunc(ctx)
}

func (s *LineWorksServiceMock) Stop() {
	s.StopFunc()
}

func (s *LineWorksServiceMock) SendMessage(ctx context.Context, guildID, message string) error {
	return s.SendMessageFunc(ctx, guildID, message)
}

func (s *LineWorksServiceMock) SendMessageDirect(ctx context.Context, guildID, message string) error {
	return s.SendMessageDirectFunc(ctx, guildID, message)
}

func (s *LineWorksServiceMock) GetServiceStatus() ServiceStatus {
	return s.GetServiceStatusFunc()
}

func (s *LineWorksServiceMock) RemoveClient(guildID string) {
	s.RemoveClientFunc(guildID)
}
	
var (
	_ LineWorksServiceInterface = (*LineWorksService)(nil)
	_ LineWorksServiceInterface = (*LineWorksServiceMock)(nil)
)

// LineWorksServiceInterface defines the interface for LINE Works service
type LineWorksServiceInterface interface {
	Start(ctx context.Context) error
	Stop()
	SendMessage(ctx context.Context, guildID, message string) error
	SendMessageDirect(ctx context.Context, guildID, message string) error
	GetServiceStatus() ServiceStatus
	RemoveClient(guildID string)
}
