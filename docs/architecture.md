# HappyHouse Architecture

HappyHouse разделен на два независимых приложения:

- `backend/` содержит API на Go, реализованное в стиле clean architecture.
- `frontend/` содержит SPA на React + TypeScript.

Слои backend:

- `internal/domain` описывает доменные сущности и контракты репозиториев.
- `internal/usecase` содержит бизнес-правила.
- `internal/repository/mssql` реализует доступ к MS SQL Server.
- `internal/transport/http` содержит роутинг, middleware и DTO.

Поток данных:

1. HTTP handler принимает запрос и валидирует транспортный формат.
2. Use case проверяет членство в доме, роли и прикладные правила.
3. Repository читает или изменяет данные в MS SQL Server.
4. Handler возвращает единообразный JSON-ответ.

Frontend структура:

- `src/app` для роутера, layout и providers.
- `src/pages` для экранов.
- `src/features` для сценариев и форм.
- `src/entities` для типизации доменных сущностей.
- `src/shared` для API-клиента, UI и утилит.
