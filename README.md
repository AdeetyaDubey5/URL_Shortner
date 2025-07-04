# URL Shortener in Go

This is a simple yet powerful URL shortener service built in Go. It allows users to shorten long URLs into compact codes and redirect them to the original addresses. The service includes intelligent hashing, in-memory storage, HATEOAS-style JSON responses, and a minimal HTTP API.

---

## Features

- **Shorten URLs** using `MD5` hashing to generate 8-character short codes.
- **In-Memory Storage** using a Go map to store original and shortened URLs.
- **Redirection Support** with HTTP 302 (Found) response.
- **HATEOAS JSON Response** with useful links for clients.
- **Error Handling** for invalid inputs or unknown codes.
- **Modular Handlers** for `/shorten`, `/redirect/{code}`, and root route.

---

## How It Works

- `POST /shorten`: Submit a long URL and get a short code with links.
- `GET /redirect/{short_code}`: Automatically redirects to the original URL.

---

## Sample Request & Response

### POST `/shorten`

**Request Body:**
```json
{
  "url": "https://github.com"
}
```

**Response:**
```json
{
    "short_code": "3097fca9",
    "_links": [
        {
            "href": "http://localhost:3000/redirect/3097fca9",
            "rel": "redirect_to_original",
            "type": "text/html"
        },
        {
            "href": "http://localhost:3000/shorten",
            "rel": "self",
            "type": "application/json"
        }
    ]
}
```

---

## How to Run

### Prerequisites

- Go 1.17 or higher installed

### Steps

1. **Clone the repository:**
```bash
git clone https://github.com/your-username/go-url-shortener
cd go-url-shortener
```

2. **Run the application:**
```bash
go run main.go
```

3. The server starts on:
```
http://localhost:3000
```

---

## API Endpoints

| Method | Endpoint             | Description                        |
|--------|----------------------|------------------------------------|
| POST   | `/shorten`           | Accepts a long URL and returns a short code |
| GET    | `/redirect/{code}`   | Redirects to the original URL using the code |
| GET    | `/`                  | Basic landing handler              |

---

## Technologies Used

- Go `net/http`
- `crypto/md5` and `encoding/hex` for generating URL hashes
- JSON encoding with `encoding/json`
- Modular handler functions

---

## Example Output

Terminal output when server starts:
```
Starting URL Shortner
Original URL:  https://gemini.google.com/app/0f81ef95d7aacd96
Shortend URL:  1a2b3c4d
Starting serve on Port:3000
```

---

## Notes

- This version stores data in-memory using a map (`urlDB`). Data will be lost when the server restarts.
- For production use, integrate with a persistent data store (e.g., Redis, PostgreSQL).
- The hash function (MD5) is used only for educational/demonstration purposes and is not cryptographically secure.

---


