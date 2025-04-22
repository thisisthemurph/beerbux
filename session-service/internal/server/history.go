package server

import (
	"context"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/history"
	"github.com/thisisthemurph/beerbux/session-service/protos/historypb"
)

type HistoryServer struct {
	historypb.UnimplementedHistoryServer
	historyRepository history.HistoryRepository
}

func NewHistoryServer(historyRepository history.HistoryRepository) *HistoryServer {
	return &HistoryServer{
		historyRepository: historyRepository,
	}
}

func (s *HistoryServer) GetBySessionID(ctx context.Context, r *historypb.GetBySessionIDRequest) (*historypb.SessionHistoryResponse, error) {
	sessionHistory, err := s.historyRepository.GetBySessionID(ctx, r.SessionId)
	if err != nil {
		return nil, err
	}

	return sessionHistory, nil
}
