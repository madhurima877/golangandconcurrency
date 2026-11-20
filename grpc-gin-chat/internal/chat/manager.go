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

func (m *Manager) Register() {

}
