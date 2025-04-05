
(readme created by AI) 
# 📝 test-task1:

---

**Project Description**  
A simple REST API written in Go with JWT-based authentication and PostgreSQL database support.

---

**Features**
- User registration with hashed password (bcrypt)
- Login with JWT token generation
- Protected endpoints using JWT
- CRUD operations on users
- Input validation (email, password length)
- PostgreSQL as storage
- Uses context and graceful shutdown

---

**Tech Stack**
- Go (net/http)
- PostgreSQL
- JWT (`github.com/golang-jwt/jwt/v5`)
- Bcrypt (`golang.org/x/crypto/bcrypt`)
- Validator (`github.com/go-playground/validator/v10`)
- SQL driver (`github.com/lib/pq`)

---

**Environment Variables (`.env`)**
```
CONFIG_PATH=./config/local.yaml
DB_PASSWORD=postgres
HASH_COST=10
JWT_SECRET=somesecret
```

---

**How to Run**
1. Create a PostgreSQL database (`test-task1`)
2. Set up your `.env` file
3. Run database migrations
4. Run the app using:
   ```bash
   go run cmd/app/main.go
   ```

---

**Available Endpoints**

| Method | Endpoint       | Auth | Description                   |
|--------|----------------|------|-------------------------------|
| POST   | `/login`       | ❌   | Login and get JWT             |
| POST   | `/users`       | ❌   | Register a new user           |
| GET    | `/users`       | ✅   | Get all users                 |
| GET    | `/users/{id}`  | ✅   | Get user by ID                |
| PUT    | `/users/{id}`  | ✅   | Update user (name/email)      |
| DELETE | `/users/{id}`  | ✅   | Delete user by ID             |

---

**Assignment Requirements**
- ✅ All endpoints implemented
- ✅ JWT authorization
- ✅ Passwords hashed with bcrypt
- ✅ Unique email constraint
- ✅ PostgreSQL storage
- ✅ Proper request validation
- ✅ Unauthorized access returns 401

---

## 📡 API Endpoints

All endpoints (except `POST /users` and `POST /login`) **require a valid JWT** in the `Authorization: Bearer <token>` header.

---

### 🔐 `POST /login`

**Description:** Authenticates the user and returns a JWT token.  
**Auth:** ❌ No  
**Body:**
```json
{
  "email": "john.doe@example.com",
  "password": "P@ssw0rd123"
}
```

**Response:**
```json
{
  "token": "<jwt_token>"
}
```

---

### ➕ `POST /users`

**Description:** Registers a new user (sign up).  
**Auth:** ❌ No  
**Body:**
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "password": "P@ssw0rd123"
}
```

**Response:**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john.doe@example.com"
}
```

---

### 📥 `GET /users`

**Description:** Returns a list of all users.  
**Auth:** ✅ Yes  
**Response:**
```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john.doe@example.com"
  },
  "..."
]
```

---

### 🔍 `GET /users/{id}`

**Description:** Returns a single user by ID.  
**Auth:** ✅ Yes  
**Response:**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john.doe@example.com"
}
```

---

### ✏️ `PUT /users/{id}`

**Description:** Updates a user’s name and email.  
**Auth:** ✅ Yes  
**Body:**
```json
{
  "name": "Updated Name",
  "email": "updated@example.com"
}
```

**Response:**
```json
{
  "name": "Updated Name",
  "email": "updated@example.com"
}
```

---

### ❌ `DELETE /users/{id}`

**Description:** Deletes a user by ID.  
**Auth:** ✅ Yes  
**Response:**  
Status `204 No Content` with no body.

