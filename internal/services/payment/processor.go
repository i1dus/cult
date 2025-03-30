package payment

import (
	"bytes"
	"context"
	"crypto/sha256"
	"cult/internal/domain"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// SberProcessor реализует интеграцию с API Сбербанка
type SberProcessor struct {
	baseURL      string
	userName     string
	password     string
	token        string
	tokenExpires time.Time
	httpClient   *http.Client
}

func NewSberProcessor(baseURL, userName, password string) *SberProcessor {
	return &SberProcessor{
		baseURL:    baseURL,
		userName:   userName,
		password:   password,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// CreatePayment создает платеж в Сбербанке
func (p *SberProcessor) CreatePayment(payment *domain.Payment) (string, error) {
	// Обновляем токен если нужно
	if err := p.ensureAuth(context.Background()); err != nil {
		return "", fmt.Errorf("auth failed: %w", err)
	}

	// Формируем запрос
	reqData := map[string]interface{}{
		"orderNumber": payment.ID,
		"amount":      payment.Amount * 100, // Сбербанк ожидает сумму в копейках
		"currency":    "643",                // RUB по ISO 4217
		"returnUrl":   "https://yourdomain.com/payment/success",
		"failUrl":     "https://yourdomain.com/payment/fail",
		"description": fmt.Sprintf("Payment for %s", payment.PaymentType),
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		p.baseURL+"/payment/rest/register.do",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		OrderID string `json:"orderId"`
		FormURL string `json:"formUrl"`
		Error   struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Error.Code != "" {
		return "", fmt.Errorf("sberbank error: %s (code %s)", result.Error.Message, result.Error.Code)
	}

	return result.FormURL, nil
}

// ProcessCallback обрабатывает callback от Сбербанка
func (p *SberProcessor) ProcessCallback(data []byte, signature string) (*domain.Payment, error) {
	var callback struct {
		OrderNumber string `json:"mdOrder"`
		Status      int    `json:"status"`
		Signature   string `json:"checksum"`
	}

	if err := json.Unmarshal(data, &callback); err != nil {
		return nil, fmt.Errorf("failed to parse callback: %w", err)
	}

	// Проверяем подпись
	if !p.verifySignature(data, callback.Signature) {
		return nil, fmt.Errorf("invalid signature")
	}

	payment := &domain.Payment{
		ID:     callback.OrderNumber,
		Status: p.mapStatus(callback.Status),
		PaidAt: time.Now(),
	}

	return payment, nil
}

// RefundPayment выполняет возврат платежа
func (p *SberProcessor) RefundPayment(paymentID string, amount int64) (string, error) {
	if err := p.ensureAuth(context.Background()); err != nil {
		return "", fmt.Errorf("auth failed: %w", err)
	}

	reqData := map[string]interface{}{
		"orderId": paymentID,
		"amount":  strconv.FormatInt(amount*100, 10), // В копейках
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		p.baseURL+"/payment/rest/refund.do",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		RefundID string `json:"refundId"`
		Error    struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Error.Code != "" {
		return "", fmt.Errorf("sberbank error: %s (code %s)", result.Error.Message, result.Error.Code)
	}

	return result.RefundID, nil
}

// ensureAuth проверяет и обновляет токен авторизации
func (p *SberProcessor) ensureAuth(ctx context.Context) error {
	if p.token != "" && time.Now().Before(p.tokenExpires) {
		return nil
	}

	req, err := http.NewRequest(
		"GET",
		p.baseURL+"/payment/rest/getToken.do",
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create auth request: %w", err)
	}

	req.SetBasicAuth(p.userName, p.password)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	var authResp struct {
		Token   string `json:"token"`
		Expires int64  `json:"expires_in"`
		Error   struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("failed to decode auth response: %w", err)
	}

	if authResp.Error.Code != "" {
		return fmt.Errorf("auth failed: %s (code %s)", authResp.Error.Message, authResp.Error.Code)
	}

	p.token = authResp.Token
	p.tokenExpires = time.Now().Add(time.Duration(authResp.Expires) * time.Second)

	return nil
}

// verifySignature проверяет подпись callback
func (p *SberProcessor) verifySignature(data []byte, signature string) bool {
	// В реальной реализации нужно использовать секретный ключ для проверки подписи
	// Это упрощенный пример
	hasher := sha256.New()
	hasher.Write(data)
	computedSig := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	return computedSig == signature
}

// mapStatus преобразует статус Сбербанка в наш PaymentStatus
func (p *SberProcessor) mapStatus(sberStatus int) domain.PaymentStatus {
	switch sberStatus {
	case 0: // Заказ зарегистрирован, но не оплачен
		return domain.PaymentStatus_PENDING
	case 1: // Предавторизованная сумма захолдирована
		return domain.PaymentStatus_PROCESSING
	case 2: // Полная авторизация суммы
		return domain.PaymentStatus_COMPLETED
	case 3: // Авторизация отменена
		return domain.PaymentStatus_CANCELLED
	case 4: // По транзакции была проведена операция возврата
		return domain.PaymentStatus_REFUNDED
	case 5: // Инициирована авторизация через ACS банка-эмитента
		return domain.PaymentStatus_PROCESSING
	case 6: // Авторизация отклонена
		return domain.PaymentStatus_FAILED
	default:
		return domain.PaymentStatus_UNDEFINED_PAYMENT_STATUS
	}
}
