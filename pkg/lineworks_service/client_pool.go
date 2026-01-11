package lineworks_service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	sdklineworks "github.com/maguro-alternative/line-works-sdk-go/pkg/lineworks"
)

// ClientPool manages multiple LINE Works clients
type ClientPool struct {
	clients map[string]*ClientInfo
	mu      sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewClientPool creates a new client pool
func NewClientPool() *ClientPool {
	return &ClientPool{
		clients: make(map[string]*ClientInfo),
	}
}

// Start starts the client pool
func (cp *ClientPool) Start(ctx context.Context, configs map[string]*LineWorksConfig) error {
	cp.ctx, cp.cancel = context.WithCancel(ctx)
	
	// Start cleanup worker
	go cp.cleanupWorker()
	
	return nil
}

// Stop stops the client pool
func (cp *ClientPool) Stop() {
	if cp.cancel != nil {
		cp.cancel()
	}
	
	cp.mu.Lock()
	defer cp.mu.Unlock()
	
	// Clear all clients
	cp.clients = make(map[string]*ClientInfo)
}

// GetClient gets or creates a LINE Works client for a guild
func (cp *ClientPool) GetClient(ctx context.Context, guildID string, config *LineWorksConfig) (*sdklineworks.Client, error) {
	cp.mu.RLock()
	clientInfo, exists := cp.clients[guildID]
	cp.mu.RUnlock()
	
	// If client exists and is active, return it
	if exists && clientInfo.IsActive {
		clientInfo.mu.Lock()
		clientInfo.LastUsed = time.Now()
		clientInfo.mu.Unlock()
		return clientInfo.Client, nil
	}
	
	// Create new client
	return cp.createClient(ctx, guildID, config)
}

// createClient creates a new LINE Works client for a guild
func (cp *ClientPool) createClient(ctx context.Context, guildID string, config *LineWorksConfig) (*sdklineworks.Client, error) {
	// Create client
	client := sdklineworks.NewClient(config.WorksID, config.Password)

	slog.InfoContext(ctx, "login line works")
	// Login
	if err := client.Login(ctx); err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}
	
	// Store client info
	cp.mu.Lock()
	cp.clients[guildID] = &ClientInfo{
		Client:   client,
		Config:   config,
		IsActive: true,
		LastUsed: time.Now(),
	}
	cp.mu.Unlock()
	
	// Start keep-alive for this client
	go cp.keepAlive(guildID)
	
	return client, nil
}

// RemoveClient removes a client from the pool
func (cp *ClientPool) RemoveClient(guildID string) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	
	if clientInfo, exists := cp.clients[guildID]; exists {
		clientInfo.IsActive = false
		delete(cp.clients, guildID)
	}
}

// keepAlive maintains client connection
func (cp *ClientPool) keepAlive(guildID string) {
	ticker := time.NewTicker(30 * time.Minute) // Keep alive every 30 minutes
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			cp.mu.RLock()
			clientInfo, exists := cp.clients[guildID]
			cp.mu.RUnlock()
			
			if !exists || !clientInfo.IsActive {
				return // Client removed or inactive
			}
			
			// Test client connectivity
			myInfo := clientInfo.Client.GetMyInfo()
			if myInfo == nil {
				// Mark client as inactive and remove
				cp.RemoveClient(guildID)
				return
			}
			
		case <-cp.ctx.Done():
			return
		}
	}
}

// cleanupWorker removes unused clients
func (cp *ClientPool) cleanupWorker() {
	ticker := time.NewTicker(15 * time.Minute) // Cleanup every 15 minutes
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			cp.cleanupUnusedClients()
		case <-cp.ctx.Done():
			return
		}
	}
}

// cleanupUnusedClients removes clients that haven't been used for a while
func (cp *ClientPool) cleanupUnusedClients() {
	cutoff := time.Now().Add(-1 * time.Hour) // Remove clients unused for 1 hour
	
	cp.mu.Lock()
	defer cp.mu.Unlock()
	
	for guildID, clientInfo := range cp.clients {
		clientInfo.mu.RLock()
		lastUsed := clientInfo.LastUsed
		clientInfo.mu.RUnlock()
		
		if lastUsed.Before(cutoff) {
			clientInfo.IsActive = false
			delete(cp.clients, guildID)
		}
	}
}

// GetAllClients returns information about all clients
func (cp *ClientPool) GetAllClients() map[string]ClientStatus {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	
	status := make(map[string]ClientStatus)
	for guildID, clientInfo := range cp.clients {
		clientInfo.mu.RLock()
		status[guildID] = ClientStatus{
			IsActive:  clientInfo.IsActive,
			LastUsed:  clientInfo.LastUsed,
			WorksID:   clientInfo.Config.WorksID,
			ChannelNo: clientInfo.Config.ChannelNo,
		}
		clientInfo.mu.RUnlock()
	}
	
	return status
}
