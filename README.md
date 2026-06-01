# SilentDaily — Документация

## Что это

Асинхронный замена дейлика. Разработчики пишут боту в конце дня, LLM извлекает прогресс/планы/блокеры, утром в 08:00 тимлид получает дайджест 3-2-1, каждый разработчик — свой план на день.

---

## Архитектура

```
┌─────────────────────────────────────────────────────┐
│  Telegram Bot (polling)    REST API (Fiber :80)      │
│  /start, /help, текст  →   POST /api/v1/updates      │
│                            POST /api/v1/members      │
│                            POST /api/v1/teams        │
│                            GET  /api/v1/digest       │
└──────────────┬──────────────────────────────────────┘
               │ daily_updates (status: queued)
               ▼
┌─────────────────────────────────────────────────────┐
│  Async Dispatcher (горутина)                         │
│  FOR UPDATE SKIP LOCKED → ExtractWorker              │
│  OpenRouter LLM → structured_updates                 │
└──────────────────────────────────────────────────────┘
               │ cron 08:00
               ▼
┌─────────────────────────────────────────────────────┐
│  DigestWorker → LLM → digest → Telegram рассылка    │
└──────────────────────────────────────────────────────┘
```

### Трёхслойная архитектура

```
handler → service → repository
```

- `handler/` — HTTP и Telegram, только парсинг/ответ
- `service/` — бизнес-логика
- `repository/` — SQL через GORM
- `domain/` — чистые Go-структуры без тегов
- `model/` — GORM DAO с тегами, маппинг через generic `mapper.go`
- `provider/` — DI-контейнер на `reflect`

---

## Как получить Telegram Bot Token

1. Открой Telegram, найди **@BotFather**
2. Напиши `/newbot`
3. Придумай имя бота (например: `SilentDaily Bot`)
4. Придумай username (должен заканчиваться на `bot`, например: `silentdaily_bot`)
5. BotFather пришлёт токен вида: `7123456789:AAHxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`
6. Вставь его в `.env`:

```env
TELEGRAM_TOKEN=7123456789:AAHxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```


---
### Как получить OpenRouter API Key

1. Зайди на [openrouter.ai](https://openrouter.ai)
2. Зарегистрируйся → Keys → Create Key
3. Бесплатные модели: `google/gemma-3-12b-it:free`, `meta-llama/llama-3.1-8b-instruct:free`

---

## Запуск

### Локально через Docker

```bash

# 2. Запусти
docker-compose up --build

# 3. Проверь
curl http://localhost:8080/ping
# → pong
```

### Миграции

```bash
sh src/scripts/migrations.sh --up
```

---

## REST API

### GET /ping
Healthcheck.
```bash
curl http://localhost:8080/ping
# → pong
```

---

### POST /api/v1/teams
Создать команду. Делается один раз тимлидом.

```bash
curl -X POST http://localhost:8080/api/v1/teams \
  -H "Content-Type: application/json" \
  -d '{"name": "Backend Team"}'
```

Ответ:
```json
{"id": 1, "name": "Backend Team"}
```

---

### POST /api/v1/members
Зарегистрировать участника команды. `telegram_user_id` — это числовой ID пользователя в Telegram (не username).

```bash
curl -X POST http://localhost:8080/api/v1/members \
  -H "Content-Type: application/json" \
  -d '{
    "telegram_user_id": 123456789,
    "name": "Тим Малов",
    "team_id": 1,
    "is_lead": true
  }'
```

Ответ:
```json
{"status": "registered"}
```

> **Как узнать свой telegram_user_id:** напиши боту **@userinfobot** в Telegram — он пришлёт твой числовой ID.

---

### POST /api/v1/updates
Отправить апдейт напрямую через API (альтернатива боту).

```bash
curl -X POST http://localhost:8080/api/v1/updates \
  -H "Content-Type: application/json" \
  -d '{
    "telegram_user_id": 123456789,
    "raw_text": "Сегодня закончил авторизацию и написал тесты. Завтра займусь деплоем на staging. Блокеров нет."
  }'
```

Ответ:
```json
{"status": "queued"}
```

---

### GET /api/v1/digest
Получить дайджест команды за дату.

```bash
# Дайджест за вчера (по умолчанию)
curl "http://localhost:8080/api/v1/digest?team_id=1"

# Дайджест за конкретную дату
curl "http://localhost:8080/api/v1/digest?team_id=1&date=2026-05-31"
```

Ответ:
```json
{
  "team_id": 1,
  "date": "2026-05-31T00:00:00Z",
  "lead_digest": " ТОП-3 ПРОГРЕССА:\n1. ..."
}
```

---

## Telegram Bot

### Команды бота

| Команда | Описание |
|---------|----------|
| `/start` | Приветствие и инструкция |
| `/help` | Как пользоваться |
| любой текст | Принимается как дневной апдейт |

### Флоу для разработчика

1. Тимлид регистрирует всех через `POST /api/v1/members`
2. В конце дня каждый пишет боту в свободной форме:
   > «Сегодня закончил авторизацию, написал тесты для login endpoint. Завтра займусь деплоем. Блокеров нет.»
3. Бот отвечает: «✅ Апдейт принят!»
4. В 08:00 следующего дня:
   - Тимлид получает дайджест 3-2-1
   - Каждый разработчик получает свой план на день

### Что делает LLM

**ExtractWorker** — получает сырой текст, отправляет в OpenRouter с промптом:
```
Извлеки из текста:
PROGRESS: что сделано (через ;)
PLANS: что планируется (через ;)
BLOCKERS: что мешает (через ; или "нет")
```

**DigestWorker** — собирает все structured_updates команды за вчера, отправляет в LLM:
```
Сформируй дайджест 3-2-1:
🟢 ТОП-3 ПРОГРЕССА
📋 2 ГЛАВНЫХ ПЛАНА
🔴 1 КЛЮЧЕВОЙ БЛОКЕР
```

---

## Async-система

### Dispatcher
Горутина, которая каждые 2 секунды:
1. Берёт одну запись `daily_updates WHERE status = 'queued'` с `FOR UPDATE SKIP LOCKED`
2. Ставит статус `in_progress`, инкрементит `attempts`
3. Запускает `ExtractWorker` в отдельной горутине
4. Слушает канал результатов → ставит `done` или `failed`

`FOR UPDATE SKIP LOCKED` — гарантирует что два воркера не возьмут одну задачу, даже если запущено несколько инстансов.

### Semaphore
Ограничивает количество одновременных воркеров (по умолчанию 5).

### Retry
Если воркер упал — `attempts++`, статус обратно в `queued`. После 3 попыток — `failed`.

---