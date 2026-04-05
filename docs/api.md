# API Outline

Base URL: `/api/v1`

Auth:

- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/refresh`
- `POST /auth/logout`
- `GET /me`

Houses:

- `GET /houses`
- `POST /houses/join-by-code`

Categories:

- `GET /houses/{houseId}/categories`
- `POST /houses/{houseId}/categories`
- `PATCH /houses/{houseId}/categories/{categoryId}`
- `DELETE /houses/{houseId}/categories/{categoryId}`

Posts:

- `GET /houses/{houseId}/posts?page=1&pageSize=10&categoryId=2`
- `GET /houses/{houseId}/posts/{postId}`
- `POST /houses/{houseId}/posts`
- `PATCH /houses/{houseId}/posts/{postId}`
- `DELETE /houses/{houseId}/posts/{postId}`

Comments:

- `GET /houses/{houseId}/posts/{postId}/comments`
- `POST /houses/{houseId}/posts/{postId}/comments`

Invite codes:

- `GET /houses/{houseId}/invite-codes`
- `POST /houses/{houseId}/invite-codes`
- `PATCH /houses/{houseId}/invite-codes/{inviteCodeId}/deactivate`
