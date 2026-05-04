# Subscription Service

REST‑сервис для учёта онлайн‑подписок пользователей.
Тестовое задание для Effective Mobile (Junior Golang Developer).

[![GitHub Repository](https://img.shields.io/badge/GitHub-Repository-blue?logo=github)](https://github.com/bladsvytro/effective-mobile_test)

## Запуск

### Docker Compose
```bash
docker-compose up --build
```
Сервис будет доступен на `http://localhost:8080`.

### Локально
1. Установите Go 1.26+ и PostgreSQL.
2. Создайте БД:
   ```bash
   createdb -h localhost -U postgres subscriptions
   ```
3. При необходимости отредактируйте `config.yaml`.
4. Запустите:
   ```bash
   go run ./cmd/api
   ```

## Эндпоинты
- `POST   /api/v1/subscriptions` – создать подписку
- `GET    /api/v1/subscriptions` – список подписок (пагинация)
- `GET    /api/v1/subscriptions/:id` – получить подписку по ID
- `PUT    /api/v1/subscriptions/:id` – обновить подписку
- `DELETE /api/v1/subscriptions/:id` – удалить подписку (мягкое удаление)
- `GET    /api/v1/analytics/total` – суммарная стоимость за период

Подробная документация: `http://localhost:8080/swagger/index.html`

## Миграции

Используются версионные SQL‑миграции (папка `migrations/`). Применяются автоматически при старте.

## Логирование

Логи пишутся в stdout в JSON‑формате. Логируются HTTP‑запросы, ошибки, ключевые операции.

## Пример запроса

```bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

## Контакты

Владислав Наумов
[GitHub](https://github.com/bladsvytro)
blads.vytro@gmail.com

