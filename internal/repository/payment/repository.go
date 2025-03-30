package payment

import (
	"context"
	"cult/internal/domain"
	"errors"
	"sync"
	"time"
)

type Repository struct {
	payments map[string]*domain.Payment
	mu       sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{
		payments: make(map[string]*domain.Payment),
	}
}

func (r *Repository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.payments[payment.ID]; exists {
		return errors.New("payment already exists")
	}

	r.payments[payment.ID] = payment
	return nil
}

func (r *Repository) GetPayment(ctx context.Context, paymentID string) (*domain.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	payment, exists := r.payments[paymentID]
	if !exists {
		return nil, errors.New("payment not found")
	}

	return payment, nil
}

func (r *Repository) UpdatePayment(ctx context.Context, payment *domain.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.payments[payment.ID]; !exists {
		return errors.New("payment not found")
	}

	r.payments[payment.ID] = payment
	return nil
}

func (r *Repository) ListPayments(ctx context.Context, userID string, from, to time.Time, limit, offset int) ([]*domain.Payment, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Payment
	for _, p := range r.payments {
		if p.UserID == userID {
			if !from.IsZero() && p.CreatedAt.Before(from) {
				continue
			}
			if !to.IsZero() && p.CreatedAt.After(to) {
				continue
			}
			result = append(result, p)
		}
	}

	total := len(result)

	if offset > len(result) {
		return nil, total, nil
	}
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}

	return result[offset:end], total, nil
}
