# API Documentation

## Auth Endpoints

### Register
- **Method:** POST
- **Path:** `/user/register`
- **Auth Required:** No
- **Request Body:**
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```

### Login
- **Method:** POST
- **Path:** `/user/login`
- **Auth Required:** No
- **Request Body:**
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```

### Refresh Token
- **Method:** POST
- **Path:** `/user/refresh`
- **Auth Required:** No
- **Request Body:**
  ```json
  {
    "refresh_token": "string"
  }
  ```

### Get User Info
- **Method:** GET
- **Path:** `/user/me`
- **Auth Required:** Yes
- **Headers:** `Authorization: Bearer {token}`

## Notes Endpoints

### Get Public Note
- **Method:** GET
- **Path:** `/public/{id}`
- **Auth Required:** No
- **URL Params:** `id=[string]`

### Create Note
- **Method:** POST
- **Path:** `/note/create`
- **Auth Required:** Yes
- **Headers:** `Authorization: Bearer {token}`
- **Request Body:**
  ```json
  {
    "title": "string"
  }
  ```

### Get Note
- **Method:** GET
- **Path:** `/note/{id}`
- **Auth Required:** Yes
- **Headers:** `Authorization: Bearer {token}`
- **URL Params:** `id=[string]`

### Get User Notes
- **Method:** GET
- **Path:** `/note/list`
- **Auth Required:** Yes
- **Headers:** `Authorization: Bearer {token}`

### Update Full Note
- **Method:** PUT
- **Path:** `/note/{id}`
- **Auth Required:** Yes
- **Headers:** `Authorization: Bearer {token}`
- **URL Params:** `id=[string]`
- **Request Body:**
  ```json
  {
    "id": 0,
    "user_id": "string",
    "title": "string",
    "nodes": [
      {
        "id": 0,
        "note_id": 0,
        "order": 0,
        "content_type": "text|image",
        "content": "string",
        "image": "string"
      }
    ],
    "created_at": "timestamp",
    "updated_at": "timestamp",
    "archived_at": "timestamp",
    "public_id": "string"
  }
  ```

### Update Note Title
- **Method:** PATCH
- **Path:** `/note/{id}`
- **Auth Required:** Yes
- **Headers:** `Authorization: Bearer {token}`
- **URL Params:** `id=[string]`
- **Request Body:**
  ```json
  {
    "title": "string"
  }
  ```

### Update Note Order
- **Method:** PATCH
- **Path:** `/note/{id}/order`
- **Auth Required:** Yes
- **Headers:** `Authorization: Bearer {token}`
- **URL Params:** `id=[string]`
- **Request Body:**
  ```json
  {
    "old_order": 0,
    "new_order": 0
  }
  ```

### Make Note Public
- **Method:** PATCH
- **Path:** `/note/{id}/public`
- **Auth Required:** Yes
- **Headers:** `Authorization: Bearer {token}`
- **URL Params:** `id=[string]`

### Make Note Private
- **Method:** PATCH
- **Path:** `/note/{id}/private`
- **Auth Required:** Yes
- **Headers:** `Authorization: Bearer {token}`
- **URL Params:** `id=[string]`

### Archive Note
- **Method:** PATCH
- **Path:** `/note/{id}/archive`
- **Auth Required:** Yes
- **Headers:** `Authorization: Bearer {token}`
- **URL Params:** `id=[string]`

### Unarchive Note
- **Method:** PATCH
- **Path:** `/note/{id}/unarchive`
- **Auth Required:** Yes
- **Headers:** `Authorization: Bearer {token}`
- **URL Params:** `id=[string]`

### Delete Note
- **Method:** DELETE
- **Path:** `/note/{id}`
- **Auth Required:** Yes
- **Headers:** `Authorization: Bearer {token}`
- **URL Params:** `id=[string]`

## Note Nodes Endpoints

### Add Node
- **Method:** POST
- **Path:** `/node`
- **Auth Required:** Yes
- **Headers:** `Authorization: Bearer {token}`
- **Request Body:**
  ```json
  {
    "note_id": 0,
    "content_type": "text|image",
    "content": "string"
  }
  ```

### Get Node Image
- **Method:** GET
- **Path:** `/node/{id}/image`
- **Auth Required:** Yes
- **Headers:** `Authorization: Bearer {token}`
- **URL Params:** `id=[string]`

### Update Node Content
- **Method:** PATCH
- **Path:** `/node/{id}`
- **Auth Required:** Yes
- **Headers:** `Authorization: Bearer {token}`
- **URL Params:** `id=[string]`
- **Request Body:**
  ```json
  {
    "content": "string"
  }
  ```

### Upload Node Image
- **Method:** PATCH
- **Path:** `/node/{id}/image`
- **Auth Required:** Yes
- **Headers:** `Authorization: Bearer {token}`
- **URL Params:** `id=[string]`
- **Content-Type:** `multipart/form-data`

### Delete Node
- **Method:** DELETE
- **Path:** `/node/{id}`
- **Auth Required:** Yes
- **Headers:** `Authorization: Bearer {token}`
- **URL Params:** `id=[string]` 