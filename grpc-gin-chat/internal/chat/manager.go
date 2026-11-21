package chat

import (
	"sync"

	pb "grpc-gin-chat/gen/proto"

	"google.golang.org/grpc"
)

type Manager struct {
	mu      sync.Mutex
	clients map[string]grpc.BidiStreamingServer[pb.ChatMessage, pb.ChatMessage]
}

func (m *Manager) Register(userID string, stream grpc.BidiStreamingServer[pb.ChatMessage, pb.ChatMessage]) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[userID] = stream
}
