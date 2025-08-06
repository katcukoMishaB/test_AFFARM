Tracker-core - Микросервис для отслеживания криптовалют

Технологии

- Go- основной язык разработки
- Gin - фреймворк
- Bun DB - ORM для работы с базой данных
- PostgreSQL - база данных
- Docker - контейнеризация
- KuCoin API - источник данных о ценах

Запуск с Docker Compose

- docker-compose up --build

Локальный запуск

- go run tracker-core/cmd/service/main.go

API Endpoints

``` http POST /currency/add Content-Type: application/json ```


``` http DELETE /currency/remove Content-Type: application/json ```


``` http GET /currency/price Content-Type: application/json ```


``` http GET /currency/list ```

Table currencies
- `id`
- `symbol`
- `name`
- `is_active`
- `created_at` 
- `updated_at`

Table prices
- `id`
- `currency_id`
- `price`
- `timestamp`
- `created_at`
