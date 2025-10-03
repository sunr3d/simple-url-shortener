# Simple URL Shortener

**Simple URL Shortener** — мини-сервис для сокращения ссылок с аналитикой переходов. Создает короткие ссылки, отслеживает переходы и получает детальную статистику по дням, месяцам или User Agent.

## Быстрый старт

### Запуск в Docker (рекомендуется)

Для быстрого запуска выполните:

```bash
make up
```

Эта команда запустит все необходимые сервисы в Docker-контейнерах.

### Остановка

```bash
make down
```

## Установка и настройка

### Требования

- Go 1.24 или выше
- Docker и Docker Compose
- Make

### Установка

1. Клон репозитория:

```bash
git clone https://github.com/sunr3d/simple-url-shortener.git
cd simple-url-shortener
```

2. Запуск сервиса:

```bash
make up
```

## Использование

После запуска сервис будет доступен по адресу `http://localhost:8080`.

### API Endpoints

#### Создание короткой ссылки
```bash
curl -X POST http://localhost:8080/shorten \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com/very/long/url"}'
```

Ответ:
```json
{
  "code": "abc12",
  "short_url": "http://localhost:8080/s/abc12"
}
```

#### Переход по короткой ссылке
```bash
curl -I http://localhost:8080/s/abc12
# Возвращает 302 редирект на оригинальный URL
```

#### Аналитика переходов
```bash
# Общая статистика
curl "http://localhost:8080/analytics/abc12"

# По дням за период
curl "http://localhost:8080/analytics/abc12?group=day&from=2025-10-01&to=2025-10-31"

# По месяцам
curl "http://localhost:8080/analytics/abc12?group=month"

# По User-Agent
curl "http://localhost:8080/analytics/abc12?group=ua"
```

### Примеры ответов аналитики

```json
{
  "total": 15,
  "by_day": [
    {"date": "2025-10-02", "count": 10},
    {"date": "2025-10-03", "count": 5}
  ]
}
```

```json
{
  "total": 15,
  "by_month": [
    {"year": 2025, "month": 10, "count": 15}
  ]
}
```

```json
{
  "total": 15,
  "by_user_agent": [
    {"ua": "Mozilla/5.0 (Chrome)", "count": 10},
    {"ua": "curl/8.5.0", "count": 5}
  ]
}