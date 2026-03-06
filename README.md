# Blob (Binary Large OBject)

Minimalist file storage API with upload, download, authentication, rate limiting, and Redis/Postgres integration.

## API Routes

| Method  |           Route             | Private |  Description  |
|---------|-----------------------------|---------|---------------|
| `PUT`   | `/blob/:bucket`             | `true`  | Upload Blob   |
| `GET`   | `/blob`                     | `true`  | List blobs    |
| `GET`   | `/blob/:bucket/:id`         | `false` | Download blob |
| `HEAD`  | `/blob/:bucket/:id`         | `false` | Blob metadata |
| `DELETE`| `/blob/:bucket/:bucket/:id` | `true`  | Delete blob   |
| `GET`   | `/health`                   | `false` | Healthcheck   |
| `GET`   | `/`                         | `false` | Hello, World  |
