# Quick Start Guide

## Быстрый старт

### 1. Установка и настройка

```bash
# Клонировать репозиторий
git clone https://github.com/pinghoyk/budget-api.git
cd budget-api

# Установить зависимости
go mod download

# Создать .env файл
cp .env.example .env
```

### 2. Настройка Telegram бота

1. Откройте Telegram и найдите [@BotFather](https://t.me/BotFather)
2. Отправьте команду `/newbot` и следуйте инструкциям
3. Получите токен бота и добавьте в `.env`:
   ```
   TELEGRAM_BOT_TOKEN=123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11
   ```
4. Сгенерируйте секретный ключ для JWT:
   ```bash
   # Генерация случайного ключа (Linux/Mac)
   openssl rand -base64 32
   # или просто используйте любую длинную случайную строку
   ```
5. Добавьте JWT секрет в `.env`:
   ```
   JWT_SECRET=ваш-сгенерированный-ключ
   ```

### 3. Запуск сервера

```bash
go run cmd/api/main.go
```

Или собрать и запустить:
```bash
go build -o budget-api cmd/api/main.go
./budget-api
```

Сервер будет доступен по адресу: `http://localhost:8080`

### 4. Проверка работы

```bash
# Проверка health check
curl http://localhost:8080/api/health
# Ответ: OK
```

### 5. Интеграция Telegram Login Widget (для фронтенда)

Добавьте на вашу веб-страницу:

```html
<script async src="https://telegram.org/js/telegram-widget.js?22" 
        data-telegram-login="ваш_bot_username" 
        data-size="large" 
        data-onauth="onTelegramAuth(user)" 
        data-request-access="write">
</script>

<script type="text/javascript">
  function onTelegramAuth(user) {
    // Отправить данные на сервер
    fetch('http://localhost:8080/api/auth/telegram', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(user)
    })
    .then(response => response.json())
    .then(data => {
      // Сохранить токен
      localStorage.setItem('token', data.token);
      console.log('Авторизован:', data.user);
    });
  }
</script>
```

### 6. Использование API

После авторизации используйте полученный токен для запросов:

```bash
# Сохраните токен в переменную
TOKEN="ваш-jwt-токен"

# Получить информацию о пользователе
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/api/auth/me

# Создать счет
curl -X POST http://localhost:8080/api/accounts \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Основной счет",
       "type": "checking",
       "currency": "RUB",
       "initial_balance": 10000
     }'

# Создать категорию
curl -X POST http://localhost:8080/api/categories \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Продукты",
       "type": "expense",
       "icon": "🛒"
     }'

# Добавить транзакцию
curl -X POST http://localhost:8080/api/transactions \
     -H "Authorization: Bearer $TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "account_id": 1,
       "category_id": 1,
       "type": "expense",
       "amount": 500,
       "description": "Покупка продуктов",
       "transaction_date": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
     }'

# Получить статистику
curl -H "Authorization: Bearer $TOKEN" \
     "http://localhost:8080/api/stats/category-summary?start_date=2024-01-01&end_date=2024-12-31"
```

### 7. Типичный рабочий процесс

1. **Авторизация** → Получение JWT токена
2. **Создание счетов** → Банковские карты, наличные, и т.д.
3. **Создание категорий** → Доходы и расходы
4. **Добавление транзакций** → Учет доходов и расходов
5. **Установка бюджетов** → Планирование расходов
6. **Просмотр статистики** → Анализ финансов

### Дополнительная информация

- Полная документация API: [API_DOCS.md](API_DOCS.md)
- Changelog: [CHANGELOG.md](CHANGELOG.md)
- Основное README: [README.md](README.md)

### Troubleshooting

**Проблема:** Ошибка "JWT_SECRET не установлен"
- **Решение:** Убедитесь, что в файле `.env` указан `JWT_SECRET`

**Проблема:** Ошибка "Invalid Telegram authentication"
- **Решение:** Проверьте правильность `TELEGRAM_BOT_TOKEN` в `.env`

**Проблема:** Порт уже занят
- **Решение:** Измените `PORT` в `.env` или остановите процесс на порту 8080

**Проблема:** База данных не создается
- **Решение:** Проверьте права доступа к директории, указанной в `DB_PATH`
