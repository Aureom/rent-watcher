package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"net/url"
	"rent-watcher/internal/config"
	"rent-watcher/internal/models"
	"rent-watcher/internal/notifier"
	"rent-watcher/internal/storage"
	"strconv"
	"sync"
	"time"
)

type ArantesScraper struct {
	BaseScraper
	Config ArantesConfig
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
}

type ArantesConfig struct {
	BaseURL        string
	MaxPages       int
	UserAgent      string
	BaseParams     config.ArantesParams
	DestinationLat float64
	DestinationLng float64
}

func NewArantesScraper(config ArantesConfig, storage storage.Storage, notifier notifier.Notifier, geoProvider GeolocationProvider) *ArantesScraper {
	ctx, cancel := context.WithCancel(context.Background())
	return &ArantesScraper{
		BaseScraper: BaseScraper{
			Storage:             storage,
			Notifier:            notifier,
			GeolocationProvider: geoProvider,
			DestinationLat:      config.DestinationLat,
			DestinationLng:      config.DestinationLng,
		},
		Config: config,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (as *ArantesScraper) Scrape(ctx context.Context) error {
	as.mu.Lock()
	as.ctx, as.cancel = context.WithCancel(ctx)
	as.mu.Unlock()

	c := as.initCollector()
	as.setupCallbacks(c)
	return as.scrapePagination(c)
}

func (as *ArantesScraper) initCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.UserAgent(as.Config.UserAgent),
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*arantesimoveis.com*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})
	return c
}

func (as *ArantesScraper) setupCallbacks(c *colly.Collector) {
	c.OnHTML(".card-imovel", func(e *colly.HTMLElement) {
		as.processPropertyCard(e)
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL: %s failed with response: %v\nError: %v\n",
			r.Request.URL, r, err)
	})
}

func (as *ArantesScraper) processPropertyCard(e *colly.HTMLElement) {
	property, rawData := as.extractPropertyData(e)
	if err := as.ProcessProperty(as.ctx, property, rawData); err != nil {
		fmt.Printf("Error processing property: %v\n", err)
	}
}

func (as *ArantesScraper) extractPropertyData(e *colly.HTMLElement) (*models.Property, string) {
	var property models.Property
	jsonData := e.ChildAttr("input.json_imovel", "value")
	if err := json.Unmarshal([]byte(jsonData), &property); err != nil {
		log.Printf("Error unmarshalling JSON: %v\n", err)
		return &models.Property{}, jsonData
	}

	property.Quartos = getValueOrDefault(e.ChildText(".fa-bed + span"), property.Quartos)
	property.Banheiros = getValueOrDefault(e.ChildText(".fa-bath + span"), property.Banheiros)
	property.Metragem = getValueOrDefault(e.ChildText(".area span"), property.Metragem)
	property.Garagens = getValueOrDefault(e.ChildText(".fa-car + span"), property.Garagens)
	property.Price = getValueOrDefault(e.ChildText(".money"), property.Price)

	// add "https://www.arantesimoveis.com" to first photo URL if it's not a full URL already
	if property.FirstPhoto != "" && property.FirstPhoto[0] == '/' {
		property.FirstPhoto = fmt.Sprintf("https://www.arantesimoveis.com%s", property.FirstPhoto)
	}

	if detailLink := e.ChildAttr(".span-card-titulo a", "href"); detailLink == "" && property.ID != "" {
		property.ID = fmt.Sprintf("/detalhes/%s", property.ID)
	}

	return &property, jsonData
}

func getValueOrDefault(value, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}

func (as *ArantesScraper) scrapePagination(c *colly.Collector) error {
	baseParams := url.Values{
		"cidade":             {as.Config.BaseParams.Cidade},
		"bairro":             {as.Config.BaseParams.Bairro},
		"categoria_imovel":   {as.Config.BaseParams.CategoriaImovel},
		"tipo":               {as.Config.BaseParams.Tipo},
		"precoMin":           {as.Config.BaseParams.PrecoMin},
		"precoMax":           {as.Config.BaseParams.PrecoMax},
		"quartos":            {as.Config.BaseParams.Quartos},
		"banheiros":          {as.Config.BaseParams.Banheiros},
		"tipoOperacao":       {as.Config.BaseParams.TipoOperacao},
		"id_only_integrador": {as.Config.BaseParams.IDOnlyIntegrador},
		"id_integrador":      {as.Config.BaseParams.IDIntegrador},
		"order_by":           {as.Config.BaseParams.OrderBy},
	}

	for page := 1; page <= as.Config.MaxPages; page++ {
		select {
		case <-as.ctx.Done():
			return as.ctx.Err()
		default:
			params := baseParams
			params.Set("page", strconv.Itoa(page))
			pageURL := as.Config.BaseURL + "?" + params.Encode()

			log.Printf("Visiting page %d: %s\n", page, pageURL)
			err := c.Visit(pageURL)
			if err != nil {
				log.Printf("Failed to visit page %d: %v\n", page, err)
			}
		}
	}

	return nil
}
