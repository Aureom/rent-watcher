# 🏠 Rent Watcher

Scrapping de imoveis em imobiliarias para alertar quando um novo imovel for listado

![image](https://github.com/user-attachments/assets/6b133f16-3975-4028-81c2-60716e2fabc5)

### Imobiliarias suportadas atualmente
- [Arantes](https://arantesimoveis.com/)

-----------------------

## Configuração

Para rodar o projeto, você deve criar um arquivo `config.json` na raiz do projeto, um arquivo de exemplo é fornecido em `config.example.json`.

Aqui está um exemplo de configuração:

```json
{
  "database_url": "file:./database.db",
  "discord_token": "",
  "discord_channel": "",
  "arantes_config": {
    "base_url": "https://www.arantesimoveis.com/listagem/",
    "max_pages": 5,
    "user_agent": "",
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