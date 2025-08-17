package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/imnzr/sim-service-project/config"
	"github.com/imnzr/sim-service-project/internal/repository"
	"github.com/imnzr/sim-service-project/models"
	"github.com/imnzr/sim-service-project/utils"
)

type OrderService interface {
	BuyNumberFromService(ctx context.Context, service, country, operator string) (*models.ResponsOrderFromService, error)
	CheckSimOrderStatus(ctx context.Context, orderId int) (*models.ResponsOrderFromService, error)
	// CheckOrderStatus(ctx context.Context, orderId int) (*models.ResponsOrderFromService, error)
}

type OrderServiceImplementation struct {
	simOrderRepo repository.SimOrderRepository
	DB           *sql.DB
	Config       config.AppConfig
}

func NewOrderService(simOrderRepo repository.SimOrderRepository, db *sql.DB) OrderService {
	return &OrderServiceImplementation{
		simOrderRepo: simOrderRepo,
		DB:           db,
	}
}

// CheckOrderStatus implements OrderService.
// func (o *OrderServiceImplementation) CheckOrderStatus(ctx context.Context, orderId int) (*models.ResponsOrderFromService, error) {
// 	order, err := o.simOrderRepo.GetById(ctx, orderId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if order.SimOrderServiceId == 0 {
// 		return nil, fmt.Errorf("order not found or invalid order ID")
// 	}

// 	return o.CheckSimOrderStatus(ctx, order.SimOrderServiceId)
// }

// CheckSimOrderStatus implements OrderService.
func (o *OrderServiceImplementation) CheckSimOrderStatus(ctx context.Context, orderId int) (*models.ResponsOrderFromService, error) {
	client := http.Client{}
	url := fmt.Sprintf("%s/user/check/%d", os.Getenv("SIM_SERVICE_URL"), orderId)

	req, err := utils.NewRequestSIM("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.Config.SimServiceAPIKey))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("service error: %s", body)
	}

	var result models.ResponsOrderFromService

	errResult := json.NewDecoder(resp.Body).Decode(&result)
	if errResult != nil {
		return nil, fmt.Errorf("failed to decode response: %w", errResult)
	}
	if result.Id == 0 {
		return nil, fmt.Errorf("order not found or invalid order ID")
	}
	return &result, nil
}

func (o *OrderServiceImplementation) BuyNumberFromService(ctx context.Context, service, country, operator string) (*models.ResponsOrderFromService, error) {
	client := http.Client{}
	url := fmt.Sprintf("%s/user/buy/activation/%s/%s/%s", os.Getenv("SERVICE_API_URL"), country, operator, service)

	req, err := utils.NewRequestSIM("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	bodyStr := strings.TrimSpace(string(bodyBytes))
	log.Println("📦 Response body:", bodyStr)

	// ✅ Jika bukan JSON (biasanya plain text = error)
	if resp.StatusCode != http.StatusOK || !strings.HasPrefix(bodyStr, "{") {
		return nil, fmt.Errorf("service error: %s", bodyStr)
	}

	var result models.ResponsOrderFromService
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return &result, nil
}
