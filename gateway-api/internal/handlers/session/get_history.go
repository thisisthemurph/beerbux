package session

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/shared/send"
	"github.com/thisisthemurph/beerbux/session-service/protos/historypb"
	"github.com/thisisthemurph/beerbux/transaction-service/pkg/fn"
	"google.golang.org/protobuf/types/known/anypb"
)

var ErrUnknownHistoryEventType = errors.New("unknown event type")

type GetSessionHistoryHandler struct {
	logger        *slog.Logger
	historyClient historypb.HistoryClient
}

func NewGetSessionHistoryHandler(logger *slog.Logger, historyClient historypb.HistoryClient) http.Handler {
	return &GetSessionHistoryHandler{
		logger:        logger,
		historyClient: historyClient,
	}
}

func (h *GetSessionHistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sessionId, err := uuid.Parse(r.PathValue("sessionId"))
	if err != nil {
		send.BadRequest(w, "Invalid sessionId")
		return
	}

	history, err := h.historyClient.GetBySessionID(r.Context(), &historypb.GetBySessionIDRequest{
		SessionId: sessionId.String(),
	})
	if err != nil {
		h.logger.Error("Error fetching session history", "error", err)
		send.InternalServerError(w, "Error fetching session history")
		return
	}

	response := h.parseHistory(history)
	send.JSON(w, response, http.StatusOK)
}

const EventTypeTransactionCreated = "transaction_created"

type HistoryResponse struct {
	SessionID string         `json:"sessionId"`
	Events    []HistoryEvent `json:"events"`
}

type HistoryEvent struct {
	ID        int64       `json:"id"`
	MemberID  string      `json:"memberId"`
	EventType string      `json:"eventType"`
	EventData interface{} `json:"eventData"`
	CreatedAt string      `json:"createdAt"`
}

type TransactionCreatedEventData struct {
	TransactionID string            `json:"transactionId"`
	Lines         []TransactionLine `json:"lines"`
}

type TransactionLine struct {
	MemberID string  `json:"memberId"`
	Amount   float64 `json:"amount"`
}

func (h *GetSessionHistoryHandler) parseHistory(history *historypb.SessionHistoryResponse) HistoryResponse {
	response := HistoryResponse{
		SessionID: history.SessionId,
	}

	for _, e := range history.Events {
		eventData, err := parseEventData(e.EventType, e.EventData)
		if err != nil {
			h.logger.Error("Failed to parse event data", "eventType", e.EventData, "error", err)
			continue
		}

		response.Events = append(response.Events, HistoryEvent{
			ID:        e.Id,
			MemberID:  e.MemberId,
			EventType: e.EventType,
			EventData: eventData,
			CreatedAt: e.CreatedAt,
		})
	}

	return response
}

func parseEventData(eventType string, d *anypb.Any) (interface{}, error) {
	switch eventType {
	case EventTypeTransactionCreated:
		return parseTransactionCreatedEventData(d)
	default:
		return nil, ErrUnknownHistoryEventType
	}
}

func parseTransactionCreatedEventData(d *anypb.Any) (interface{}, error) {
	var msg historypb.TransactionCreatedEventData
	if err := d.UnmarshalTo(&msg); err != nil {
		return nil, err
	}

	return TransactionCreatedEventData{
		TransactionID: msg.TransactionId,
		Lines: fn.Map(msg.Lines, func(line *historypb.TransactionLine) TransactionLine {
			return TransactionLine{
				MemberID: line.MemberId,
				Amount:   line.Amount,
			}
		}),
	}, nil
}
