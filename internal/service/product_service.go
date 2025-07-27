package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/imnzr/sim-service-project/internal/repository"
	"github.com/imnzr/sim-service-project/models"
	"github.com/imnzr/sim-service-project/utils"
)

type ProductInformation struct {
	Category string  `json:"category"`
	Qty      int     `json:"qty"`
	Price    float64 `json:"price"`
}

type ProductService interface {
	GetProductAvailable(service, country, operator string) (map[string]ProductInformation, error)
	SyncFromSimServices(ctx context.Context) error
}

type ProductServiceImplementation struct {
	Repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &ProductServiceImplementation{
		Repo: repo,
	}
}

// GetProductAvailable implements ProductService.
func (*ProductServiceImplementation) GetProductAvailable(service string, country string, operator string) (map[string]ProductInformation, error) {
	url := fmt.Sprintf("%sguest/products/%s/%s/%s", os.Getenv("SIM_API_URL_SERVICE"), service, country, operator)

	req, err := utils.NewRequestSIM("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyBytes, _ := io.ReadAll(res.Body)
	bodyText := string(bodyBytes)

	if !strings.HasPrefix(bodyText, "{") {
		return nil, fmt.Errorf("sim service error: %s", bodyText)
	}

	var result map[string]ProductInformation

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to decode product JSON: %w", err)
	}

	return result, nil
}

// SyncFromSimServices implements ProductService.
func (service *ProductServiceImplementation) SyncFromSimServices(ctx context.Context) error {
	url := os.Getenv("SIM_API_URL_SERVICE")

	format := fmt.Sprintf("%s/guest/prices", url)
	resp, err := http.Get(format)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var raw map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return err
	}

	kursRubel := 208.0
	markup := 2000.0

	for country, opsRaw := range raw {
		opsMap, ok := opsRaw.(map[string]interface{})
		if !ok {
			continue
		}
		for operator, servicesRaw := range opsMap {
			servicesMap, ok := servicesRaw.(map[string]interface{})
			if !ok {
				fmt.Println("gagal konversi ke map[string]interface{}: ", servicesRaw)
				continue
			}
			for serviceName, priceRaw := range servicesMap {
				serviceInfo, ok := priceRaw.(map[string]interface{})
				if !ok {
					log.Printf("gagal konversi service menjadi float64")
					continue
				}

				// ambil nilai cost dari map
				priceRubelRaw, ok := serviceInfo["cost"]
				if !ok {
					log.Printf("cost tidak ditemukan di %s - %s - %s", country, operator, serviceName)
					continue
				}
				priceRubel, ok := priceRubelRaw.(float64)
				if !ok {
					log.Printf("gagal konversi cost ke float64: %v", priceRubelRaw)
				}

				priceIDR := priceRubel * kursRubel
				priceSell := priceIDR + markup

				product := models.SimProduct{
					Service:      serviceName,
					Country:      country,
					Operator:     operator,
					PriceDefault: priceRubel,
					PriceSell:    priceSell,
				}
				if err := service.Repo.Upsert(ctx, &product); err != nil {
					log.Printf("gagal menyimpan product")
					fmt.Printf("gagal menyimpan product %+v: %v\n", product, err)
				}
			}
		}
	}
	return nil
}
