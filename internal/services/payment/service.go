package payment

import (
	"context"
	"cult/internal/domain"
	desc "cult/pkg"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	repo      Repository
	processor PaymentProcessor
	mu        sync.RWMutex
	desc.UnimplementedPaymentAPIServer
}

type Repository interface {
	CreatePayment(ctx context.Context, payment *domain.Payment) error
	GetPayment(ctx context.Context, paymentID string) (*domain.Payment, error)
	UpdatePayment(ctx context.Context, payment *domain.Payment) error
	ListPayments(ctx context.Context, userID string, from, to time.Time, limit, offset int) ([]*domain.Payment, int, error)
}

type PaymentProcessor interface {
	CreatePayment(payment *domain.Payment) (paymentURL string, err error)
	ProcessCallback(data []byte, signature string) (*domain.Payment, error)
	RefundPayment(paymentID string, amount int64) (string, error)
}

func NewService(repo Repository, processor PaymentProcessor) *Service {
	return &Service{
		repo:      repo,
		processor: processor,
	}
}

// CreatePayment реализует gRPC метод создания платежа
func (s *Service) CreatePayment(ctx context.Context, req *desc.CreatePaymentRequest) (*desc.CreatePaymentResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	payment := &domain.Payment{
		ID:        uuid.New().String(),
		UserID:    req.UserId,
		Status:    domain.PaymentStatus_PENDING,
		CreatedAt: time.Now(),
		Amount:    calculateAmount(req), // функция расчета суммы
		Currency:  "RUB",
	}

	switch {
	case req.GetBookingId() != "":
		payment.BookingID = req.GetBookingId()
		payment.PaymentType = domain.PaymentType_BOOKING
	case req.GetRentalId() != "":
		payment.RentalID = req.GetRentalId()
		payment.PaymentType = domain.PaymentType_RENTAL
	default:
		return nil, status.Error(codes.InvalidArgument, "either booking_id or rental_id must be provided")
	}

	paymentURL, err := s.processor.CreatePayment(payment)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create payment: %v", err)
	}

	if err := s.repo.CreatePayment(ctx, payment); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save payment: %v", err)
	}

	return &desc.CreatePaymentResponse{
		PaymentId:  payment.ID,
		PaymentUrl: paymentURL,
		Status:     payment.Status,
	}, nil
}

// GetPaymentStatus реализует gRPC метод проверки статуса платежа
func (s *Service) GetPaymentStatus(ctx context.Context, req *desc.GetPaymentStatusRequest) (*desc.GetPaymentStatusResponse, error) {
	if req.PaymentId == "" {
		return nil, status.Error(codes.InvalidArgument, "payment_id is required")
	}

	payment, err := s.repo.GetPayment(ctx, req.PaymentId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "payment not found: %v", err)
	}

	return &desc.GetPaymentStatusResponse{
		Status:      payment.Status,
		PaymentDate: timestamppb.New(payment.PaidAt),
		Amount:      payment.Amount,
		Currency:    payment.Currency,
	}, nil
}

// PaymentCallback обработчик callback от платежной системы
func (s *Service) PaymentCallback(ctx context.Context, req *desc.PaymentCallbackRequest) (*desc.PaymentCallbackResponse, error) {
	payment, err := s.processor.ProcessCallback([]byte(req.PaymentId), req.Signature)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid callback: %v", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	existing, err := s.repo.GetPayment(ctx, payment.ID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "payment not found: %v", err)
	}

	existing.Status = payment.Status
	existing.PaidAt = payment.PaidAt

	if err := s.repo.UpdatePayment(ctx, existing); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update payment: %v", err)
	}

	// Здесь можно добавить уведомление других сервисов об изменении статуса платежа

	return &desc.PaymentCallbackResponse{Success: true}, nil
}

// GetPaymentHistory реализует gRPC метод получения истории платежей
func (s *Service) GetPaymentHistory(ctx context.Context, req *desc.GetPaymentHistoryRequest) (*desc.GetPaymentHistoryResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	var from, to time.Time
	if req.From != nil {
		from = req.From.AsTime()
	}
	if req.To != nil {
		to = req.To.AsTime()
	}

	payments, total, err := s.repo.ListPayments(ctx, req.UserId, from, to, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get payments: %v", err)
	}

	resp := &desc.GetPaymentHistoryResponse{
		Total: int32(total),
	}

	for _, p := range payments {
		item := &desc.PaymentHistoryItem{
			PaymentId:     p.ID,
			Amount:        p.Amount,
			Currency:      p.Currency,
			Status:        p.Status,
			CreatedAt:     timestamppb.New(p.CreatedAt),
			PaymentMethod: p.Method,
		}

		switch p.PaymentType {
		case domain.PaymentType_BOOKING:
			item.PaymentType = &desc.PaymentHistoryItem_BookingId{BookingId: p.BookingID}
		case domain.PaymentType_RENTAL:
			item.PaymentType = &desc.PaymentHistoryItem_RentalId{RentalId: p.RentalID}
		}

		resp.Payments = append(resp.Payments, item)
	}

	return resp, nil
}

// RefundPayment реализует возврат платежа
func (s *Service) RefundPayment(ctx context.Context, req *desc.RefundPaymentRequest) (*desc.RefundPaymentResponse, error) {
	if req.PaymentId == "" {
		return nil, status.Error(codes.InvalidArgument, "payment_id is required")
	}

	payment, err := s.repo.GetPayment(ctx, req.PaymentId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "payment not found: %v", err)
	}

	if payment.Status != desc.PaymentStatus_COMPLETED {
		return nil, status.Errorf(codes.FailedPrecondition, "only completed payments can be refunded")
	}

	refundID, err := s.processor.RefundPayment(payment.ID, payment.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process refund: %v", err)
	}

	payment.Status = desc.PaymentStatus_REFUNDED
	if err := s.repo.UpdatePayment(ctx, payment); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update payment status: %v", err)
	}

	return &desc.RefundPaymentResponse{
		RefundId: refundID,
		Status:   payment.Status,
	}, nil
}

func calculateAmount(req *desc.CreatePaymentRequest) int64 {
	// Здесь должна быть логика расчета суммы платежа
	// В зависимости от типа (бронирование или аренда)
	// и других параметров
	return 1000 // Пример фиксированной суммы
}
