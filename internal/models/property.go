package models

type Property struct {
	ID             string `json:"id"`
	FirstPhoto     string `json:"first_foto"`
	Price          string `json:"preco"`
	Logradouro     string `json:"logradouro"`
	Bairro         string `json:"bairro"`
	Cidade         string `json:"cidade"`
	Metragem       string `json:"metragem"`
	Quartos        string `json:"quartos"`
	Banheiros      string `json:"banheiros"`
	Suites         string `json:"suites"`
	Garagens       string `json:"garagens"`
	TipoImovel     string `json:"tipo_imovel"`
	DistanceMeters int    `json:"distance_meters"`
}
