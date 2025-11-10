# API Documentation

## Обзор

Budget API предоставляет полноценный REST API для управления личным бюджетом с авторизацией через Telegram.

Базовый URL: `http://localhost:8080/api`

## Аутентификация

Все защищенные endpoints требуют JWT токен в заголовке:
```
Authorization: Bearer <your-jwt-token>
```

### Получение токена

Токен можно получить через Telegram авторизацию.

## Endpoints

### 1. Аутентификация

#### POST /api/auth/telegram
Авторизация/регистрация через Telegram Login Widget.

**Request:**
```json
{
  "id": 123456789,
  "first_name": "Иван",
  "last_name": "Иванов",
  "username": "ivan_ivanov",
  "photo_url": "https://t.me/i/userpic/320/...",
  "auth_date": 1735588800,
  "hash": "abcdef1234567890..."
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "telegram_id": 123456789,
    "telegram_username": "ivan_ivanov",
    "first_name": "Иван",
    "last_name": "Иванов",
    "name": "Иван Иванов",
    "created_at": "2024-01-10T10:00:00Z",
    "updated_at": "2024-01-10T10:00:00Z"
  }
}
```

#### GET /api/auth/me
Получить информацию о текущем пользователе (требуется авторизация).

**Response:**
```json
{
  "id": 1,
  "telegram_id": 123456789,
  "name": "Иван Иванов",
  ...
}
```

---

### 2. Счета (Accounts)

#### GET /api/accounts
Получить все счета текущего пользователя.

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "name": "Основная карта",
    "type": "checking",
    "currency": "RUB",
    "initial_balance": 10000.0,
    "current_balance": 12500.0,
    "color": "#FF5733",
    "icon": "💳",
    "is_active": true,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-10T15:30:00Z"
  }
]
```

#### POST /api/accounts
Создать новый счет.

**Типы счетов:**
- `checking` - Расчетный счет
- `savings` - Сберегательный счет
- `cash` - Наличные
- `credit` - Кредитная карта
- `investment` - Инвестиционный счет
- `other` - Другое

**Request:**
```json
{
  "name": "Основная карта",
  "type": "checking",
  "currency": "RUB",
  "initial_balance": 10000.0,
  "color": "#FF5733",
  "icon": "💳"
}
```

#### GET /api/accounts/{id}
Получить конкретный счет по ID.

#### PUT /api/accounts/{id}
Обновить счет.

**Request:**
```json
{
  "name": "Обновленное название",
  "type": "savings",
  "currency": "RUB",
  "color": "#4CAF50",
  "icon": "💰"
}
```

#### DELETE /api/accounts/{id}
Удалить счет.

---

### 3. Категории (Categories)

#### GET /api/categories
Получить все категории. Опциональный query параметр: `?type=income` или `?type=expense`.

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "name": "Зарплата",
    "type": "income",
    "color": "#4CAF50",
    "icon": "💰",
    "parent_id": null,
    "is_active": true,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  },
  {
    "id": 2,
    "user_id": 1,
    "name": "Продукты",
    "type": "expense",
    "color": "#FF5733",
    "icon": "🛒",
    "parent_id": null,
    "is_active": true,
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
]
```

#### POST /api/categories
Создать категорию.

**Request:**
```json
{
  "name": "Продукты",
  "type": "expense",
  "color": "#4CAF50",
  "icon": "🛒",
  "parent_id": null
}
```

#### GET /api/categories/{id}
Получить категорию по ID.

#### PUT /api/categories/{id}
Обновить категорию.

#### DELETE /api/categories/{id}
Удалить категорию.

---

### 4. Транзакции (Transactions)

#### GET /api/transactions
Получить транзакции с фильтрацией.

**Query Parameters:**
- `account_id` - ID счета
- `category_id` - ID категории
- `type` - Тип (income/expense/transfer)
- `start_date` - Начальная дата (YYYY-MM-DD)
- `end_date` - Конечная дата (YYYY-MM-DD)
- `limit` - Количество записей (по умолчанию: 50)
- `offset` - Смещение для пагинации

**Пример:** `GET /api/transactions?type=expense&limit=10&start_date=2024-01-01`

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "account_id": 1,
    "category_id": 2,
    "type": "expense",
    "amount": 500.0,
    "currency": "RUB",
    "description": "Покупка продуктов",
    "transaction_date": "2024-01-15T10:30:00Z",
    "notes": "Магазин у дома",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
]
```

#### POST /api/transactions
Создать транзакцию.

**Типы транзакций:**
- `income` - Доход (увеличивает баланс счета)
- `expense` - Расход (уменьшает баланс счета)
- `transfer` - Перевод (между двумя счетами)

**Пример расхода:**
```json
{
  "account_id": 1,
  "category_id": 2,
  "type": "expense",
  "amount": 500.0,
  "currency": "RUB",
  "description": "Покупка продуктов",
  "transaction_date": "2024-01-15T10:30:00Z",
  "notes": "Магазин у дома"
}
```

**Пример дохода:**
```json
{
  "account_id": 1,
  "category_id": 1,
  "type": "income",
  "amount": 50000.0,
  "currency": "RUB",
  "description": "Зарплата за январь",
  "transaction_date": "2024-01-31T00:00:00Z"
}
```

**Пример перевода:**
```json
{
  "account_id": 1,
  "type": "transfer",
  "amount": 1000.0,
  "currency": "RUB",
  "to_account_id": 2,
  "description": "Перевод на сберегательный счет",
  "transaction_date": "2024-01-15T10:30:00Z"
}
```

#### GET /api/transactions/{id}
Получить транзакцию по ID.

#### PUT /api/transactions/{id}
Обновить транзакцию (можно изменить только описание, категорию, дату и заметки).

**Request:**
```json
{
  "category_id": 3,
  "description": "Обновленное описание",
  "transaction_date": "2024-01-15T12:00:00Z",
  "notes": "Новые заметки"
}
```

#### DELETE /api/transactions/{id}
Удалить транзакцию. Баланс счетов автоматически корректируется.

---

### 5. Бюджеты (Budgets)

#### GET /api/budgets
Получить все бюджеты пользователя.

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "category_id": 2,
    "amount": 15000.0,
    "period": "monthly",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-12-31T23:59:59Z",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  }
]
```

#### POST /api/budgets
Создать бюджет для категории.

**Периоды:**
- `monthly` - Ежемесячный
- `yearly` - Ежегодный

**Request:**
```json
{
  "category_id": 2,
  "amount": 15000.0,
  "period": "monthly",
  "start_date": "2024-01-01",
  "end_date": "2024-12-31"
}
```

#### GET /api/budgets/{id}
Получить бюджет по ID.

#### GET /api/budgets/{id}/status
Получить статус бюджета (потрачено, осталось, процент выполнения).

**Response:**
```json
{
  "budget": {
    "id": 1,
    "category_id": 2,
    "amount": 15000.0,
    ...
  },
  "spent": 8500.0,
  "remaining": 6500.0,
  "percentage": 56.67,
  "is_exceeded": false
}
```

#### PUT /api/budgets/{id}
Обновить бюджет.

**Request:**
```json
{
  "amount": 20000.0,
  "period": "monthly",
  "start_date": "2024-01-01",
  "end_date": "2024-12-31"
}
```

#### DELETE /api/budgets/{id}
Удалить бюджет.

---

### 6. Статистика (Statistics)

#### GET /api/stats/category-summary
Получить сводку расходов/доходов по категориям.

**Query Parameters:**
- `start_date` - Начальная дата (YYYY-MM-DD)
- `end_date` - Конечная дата (YYYY-MM-DD)

**Пример:** `GET /api/stats/category-summary?start_date=2024-01-01&end_date=2024-01-31`

**Response:**
```json
[
  {
    "category_id": 2,
    "category_name": "Продукты",
    "type": "expense",
    "total_amount": 12500.0,
    "count": 45
  },
  {
    "category_id": 3,
    "category_name": "Транспорт",
    "type": "expense",
    "total_amount": 3500.0,
    "count": 20
  },
  {
    "category_id": 1,
    "category_name": "Зарплата",
    "type": "income",
    "total_amount": 50000.0,
    "count": 1
  }
]
```

#### GET /api/stats/monthly-balance
Получить помесячный баланс доходов и расходов.

**Query Parameters:**
- `start_date` - Начальная дата (YYYY-MM-DD)
- `end_date` - Конечная дата (YYYY-MM-DD)

**Пример:** `GET /api/stats/monthly-balance?start_date=2023-01-01&end_date=2024-12-31`

**Response:**
```json
[
  {
    "month": "2024-01",
    "income": 50000.0,
    "expense": 35000.0,
    "balance": 15000.0
  },
  {
    "month": "2024-02",
    "income": 50000.0,
    "expense": 32000.0,
    "balance": 18000.0
  }
]
```

---

## Коды ошибок

- `200 OK` - Успешный запрос
- `201 Created` - Ресурс успешно создан
- `204 No Content` - Успешное удаление
- `400 Bad Request` - Неверный запрос
- `401 Unauthorized` - Требуется авторизация
- `403 Forbidden` - Доступ запрещен
- `404 Not Found` - Ресурс не найден
- `500 Internal Server Error` - Внутренняя ошибка сервера

## Примеры использования

### Создание полного рабочего процесса

1. **Авторизация через Telegram**
```bash
curl -X POST http://localhost:8080/api/auth/telegram \
  -H "Content-Type: application/json" \
  -d '{"id": 123456789, "first_name": "Test", ...}'
```

2. **Создание счета**
```bash
curl -X POST http://localhost:8080/api/accounts \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Основная карта", "type": "checking", "currency": "RUB", "initial_balance": 10000}'
```

3. **Создание категорий**
```bash
# Категория дохода
curl -X POST http://localhost:8080/api/categories \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Зарплата", "type": "income", "icon": "💰"}'

# Категория расхода
curl -X POST http://localhost:8080/api/categories \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Продукты", "type": "expense", "icon": "🛒"}'
```

4. **Добавление транзакций**
```bash
# Доход
curl -X POST http://localhost:8080/api/transactions \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"account_id": 1, "category_id": 1, "type": "income", "amount": 50000, "description": "Зарплата"}'

# Расход
curl -X POST http://localhost:8080/api/transactions \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"account_id": 1, "category_id": 2, "type": "expense", "amount": 500, "description": "Продукты"}'
```

5. **Просмотр статистики**
```bash
curl -X GET "http://localhost:8080/api/stats/category-summary?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer <token>"
```
