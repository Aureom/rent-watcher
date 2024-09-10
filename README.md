# 🏠 Rent Watcher

Scrapping de imóveis que monitora imobiliárias e alerta sempre que um novo imóvel for listado

![image](https://github.com/user-attachments/assets/6b133f16-3975-4028-81c2-60716e2fabc5)

-----------------------

## 🔍 Imobiliárias Suportadas
- [Arantes Imóveis](https://arantesimoveis.com/)

-----------------------

## ⚙️ Configuração

Você precisará criar um arquivo `config.json` na raiz do projeto. Um exemplo de configuração é fornecido no arquivo `config-example.json`.

Exemplo de configuração:

```json
{
  "database_url": "file:./database.db",
  "discord_token": "<discord_bot_token>",
  "discord_channel": "<discord_channel_id>",
  "google_maps_api_key": "<api_key_google_maps>",
  "destination_lat": -10.8249467,
  "destination_lng": -42.7278008,
  "arantes_config": {
    "base_url": "https://www.arantesimoveis.com/listagem/",
    "max_pages": 5,
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
    "base_params": {
      "cidade": "1",
      "bairro": "142",
      "categoria_imovel": "1",
      "tipo": "2",
      "precoMin": "",
      "precoMax": "2.000,00",
      "quartos": "",
      "banheiros": "",
      "tipoOperacao": "2",
      "id_only_integrador": "",
      "id_integrador": "",
      "order_by": ""
    }
  }
}
```

- `discord_token`: Token do bot do Discord.
- `discord_channel`: O ID do canal onde as notificações serão enviadas.
- `google_maps_api_key`: Chave da API do Google Maps (opcional).
- `destination_lat` | `destination_lng`: Latitude e Longitude do local que você deseja calcular a distância a partir dos imóveis.