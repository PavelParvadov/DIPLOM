# Run Guide

## 1. Create PostgreSQL database

You can use either:

1. local PostgreSQL
2. PostgreSQL in Docker
3. Railway PostgreSQL

Database name for local development:

- `happyhouse`

## 2. Apply schema

Run SQL files in this order:

1. `db/migrations/001_create_schema.sql`
2. `db/migrations/002_add_post_image.sql`
3. `db/migrations/003_create_chat_messages.sql`
4. `db/migrations/004_create_media_assets.sql`

Optional seed files:

- `db/seed/001_demo_seed.sql`
- `db/seed/003_penza_showcase_seed.sql`

## 3. Configure backend

1. Copy [backend/.env.example](/C:/Users/User/go-services/DIPLOM/backend/.env.example) to `backend/.env`
2. Set actual `DATABASE_URL`

Example:

```env
DATABASE_URL=postgres://happyhouse:happyhouse@localhost:5432/happyhouse?sslmode=disable
```

## 4. Start backend

```powershell
cd backend
go mod tidy
go run ./cmd/api
```

## 5. Start frontend

```powershell
cd frontend
npm install
npm run dev
```

## 6. Demo flow

1. Зарегистрировать пользователя.
2. Войти по логину и паролю.
3. Создать дом или вступить в дом по invite code.
4. Создать пост.
5. Отфильтровать посты по категории.
6. Открыть пост и добавить комментарий.
7. Открыть чат дома и отправить сообщение.
8. Под ролью admin создать категорию или деактивировать invite code.

## Demo accounts

Basic seed:

- `admin_demo` / `demo1234`
- `resident_demo` / `demo1234`
- invite code: `HAPPY2026`

Showcase seed:

- `alexey.morozov` / `demo1234`
- `irina.lapina` / `demo1234`
