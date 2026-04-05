# Run Guide

## 1. Create the database

1. Создайте базу `HappyHouse` в MS SQL Server.
2. Выполните `db/migrations/001_create_schema.sql`.
3. При необходимости загрузите демо-данные из `db/seed/001_demo_seed.sql`.

## 2. Configure backend

1. Скопируйте `backend/.env.example` в `backend/.env` или задайте переменные окружения вручную.
2. Укажите актуальный `DATABASE_URL`.

## 3. Start backend

```powershell
cd backend
go mod tidy
go run ./cmd/api
```

## 4. Start frontend

```powershell
cd frontend
npm install
npm run dev
```

## 5. Demo flow

1. Зарегистрировать пользователя.
2. Войти по логину и паролю.
3. Вступить в дом по invite code.
4. Создать пост.
5. Отфильтровать посты по категории.
6. Открыть пост и добавить комментарий.
7. Под admin-ролью создать категорию или деактивировать invite code.

Демо-учетки из seed:

- `admin_demo` / `demo1234`
- `resident_demo` / `demo1234`
- invite code: `HAPPY2026`
