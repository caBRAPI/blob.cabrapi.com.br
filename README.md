# Blob (Binary Large OBject)

Minimalist file storage API with upload, download, authentication, rate limiting, and Redis/Postgres integration.

## API Routes

| Method  |           Route             | Private |  Description  |
|---------|-----------------------------|---------|---------------|
| `PUT`   | `/blob`                     | `true`  | Upload Blob   |
| `GET`   | `/blob`                     | `true`  | List blobs    |
| `POST`  | `/blob/:uuid`               | `true`  | Edit metadata |
| `GET`   | `/blob/:uuid`               | `false` | Download blob |
| `HEAD`  | `/blob/:uuid`               | `false` | Blob metadata |
| `DELETE`| `/blob/:uuid`               | `true`  | Delete blob   |
| `GET`   | `/health`                   | `false` | Healthcheck   |
| `GET`   | `/`                         | `false` | Hello, World  |

## Database Schema
