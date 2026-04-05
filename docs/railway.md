# Deploy to Railway

## Target setup

HappyHouse is now prepared for a Railway-native deployment:

1. `backend` service from `/backend`
2. `frontend` service from `/frontend`
3. `PostgreSQL` database provisioned directly in Railway

This is the simplest and cheapest production-like setup for the current project.

## Railway services

### 1. PostgreSQL

Create a PostgreSQL database inside Railway first.

Railway will generate a connection string automatically. Use that value as `DATABASE_URL` in the backend service.

### 2. Backend

Source folder: `/backend`

Required variables:

- `DATABASE_URL`
- `JWT_SECRET`
- `FRONTEND_ORIGIN`
- `UPLOAD_DIR=uploads`
- `ACCESS_TOKEN_TTL_MINUTES=30`
- `REFRESH_TOKEN_TTL_HOURS=168`
- `HTTP_READ_TIMEOUT_SECONDS=10`
- `HTTP_WRITE_TIMEOUT_SECONDS=15`
- `HTTP_SHUTDOWN_TIMEOUT_SECONDS=10`

Public health endpoint:

- `/health`

### 3. Frontend

Source folder: `/frontend`

Required build variable:

- `VITE_API_BASE_URL=https://<your-backend-domain>/api/v1`

Important:

- `VITE_API_BASE_URL` is read at build time
- set it before the frontend deploy starts

## Recommended deploy order

1. Create one Railway project.
2. Add PostgreSQL.
3. Add service `backend` from `/backend`.
4. Add service `frontend` from `/frontend`.
5. Put PostgreSQL `DATABASE_URL` into backend variables.
6. Deploy backend.
7. Generate backend public domain.
8. Put `VITE_API_BASE_URL=https://<backend-domain>/api/v1` into frontend variables.
9. Deploy frontend.
10. Generate frontend public domain.
11. Update backend variable `FRONTEND_ORIGIN=https://<frontend-domain>`.
12. Redeploy backend.

## Database bootstrap

After PostgreSQL is created, run the SQL files in this order:

1. `db/migrations/001_create_schema.sql`
2. `db/migrations/002_add_post_image.sql`
3. `db/migrations/003_create_chat_messages.sql`
4. optional demo data: `db/seed/001_demo_seed.sql`
5. optional showcase data for diploma demo: `db/seed/003_penza_showcase_seed.sql`

If you use a fresh database, `001_create_schema.sql` already contains the full current schema. Files `002` and `003` are left as safe incremental migrations.

You can apply these files without `psql` by using the local Go helper:

```powershell
cd backend
$env:DATABASE_URL="<paste Railway DATABASE_PUBLIC_URL here>"
go run ./cmd/dbexec ..\db\migrations\001_create_schema.sql ..\db\migrations\002_add_post_image.sql ..\db\migrations\003_create_chat_messages.sql ..\db\seed\003_penza_showcase_seed.sql
```

## What is already ready in the repo

- backend Dockerfile: [backend/Dockerfile](/C:/Users/User/go-services/DIPLOM/backend/Dockerfile)
- frontend Dockerfile: [frontend/Dockerfile](/C:/Users/User/go-services/DIPLOM/frontend/Dockerfile)
- backend env example: [backend/.env.example](/C:/Users/User/go-services/DIPLOM/backend/.env.example)
- frontend env example: [frontend/.env.example](/C:/Users/User/go-services/DIPLOM/frontend/.env.example)

## Local PostgreSQL example

For local development, a typical connection string looks like this:

```env
DATABASE_URL=postgres://happyhouse:happyhouse@localhost:5432/happyhouse?sslmode=disable
```

## Conclusion

For Railway this architecture is now straightforward:

- backend on Railway
- frontend on Railway
- PostgreSQL on Railway

No external MSSQL hosting is required anymore.
