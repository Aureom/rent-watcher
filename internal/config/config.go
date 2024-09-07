package config

import (
	"encoding/json"
	"net/url"
	"os"
)

type Config struct {
	DatabaseURL    string        `json:"database_url"`
	DiscordToken   string        `json:"discord_token"`
	DiscordChannel string        `json:"discord_channel"`
	ArantesConfig  ArantesConfig `json:"arantes_config"`
}

type ArantesConfig struct {
	BaseURL    string        `json:"base_url"`
	MaxPages   int           `json:"max_pages"`
	UserAgent  string        `json:"user_agent"`
	BaseParams ArantesParams `json:"base_params"`
}

type ArantesParams struct {
	Cidade           string `json:"cidade"`
	Bairro           string `json:"bairro"`
	CategoriaImovel  string `json:"categoria_imovel"`
	Tipo             string `json:"tipo"`
	PrecoMin         string `json:"precoMin"`
	PrecoMax         string `json:"precoMax"`
	Quartos          string `json:"quartos"`
	Banheiros        string `json:"banheiros"`
	TipoOperacao     string `json:"tipoOperacao"`
	IDOnlyIntegrador string `json:"id_only_integrador"`
	IDIntegrador     string `json:"id_integrador"`
	OrderBy          string `json:"order_by"`
}

type URLValues url.Values

func Load() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
