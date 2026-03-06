# Blob (Binary Large OBject)

🗂️ Minimalist file storage API with upload, download, authentication, rate limiting, and Redis/Postgres integration.

## Routes

| Method  |           Route               | Private | Description                           |
|---------|-------------------------------|---------|---------------------------------------|
| `PUT`   | `/blob`                       | `true`  | Upload Blob (full)                    |
| `POST`  | `/blob/:id/stream`            | `true`  | Upload Blob in streaming mode         |
| `GET`   | `/blob`                       | `true`  | List blobs                            |
| `GET`   | `/blob/:id`                   | `false` | Download Blob (full)                  |
| `GET`   | `/blob/:id/stream`            | `false` | Download Blob in streaming mode       |
| `HEAD`  | `/blob/:id`                   | `false` | Blob metadata                          |
| `POST`  | `/blob/:id/metadata`          | `true`  | Edit metadata                          |
| `DELETE`| `/blob/:id`                   | `true`  | Delete blob                            |
| `GET`   | `/health`                     | `false` | Healthcheck                            |
| `GET`   | `/`                           | `false` | Hello, World                            |

## Database Schema

| Column           | Type        | Nullable | Description                                      |
|------------------|------------ |--------- |-------------------------------------------------|
| `id`             | UUID        | No       | Unique identifier for each blob                 |
| `bucket`         | TEXT        | No       | Logical grouping for files                      |
| `filename`       | TEXT        | No       | Original file name                              |
| `mime`           | TEXT        | No       | File MIME type                                  |
| `size`           | BIGINT      | No       | File size in bytes                               |
| `hash`           | TEXT        | No       | Unique hash for integrity/deduplication         |
| `path`           | TEXT        | No       | Storage path or reference                       |
| `public`         | BOOLEAN     | Yes      | Whether blob is publicly accessible             |
| `download_count` | INT         | Yes      | Number of times the blob has been downloaded   |
| `created_at`     | TIMESTAMPTZ | No       | Timestamp when blob was created                 |
| `updated_at`     | TIMESTAMPTZ | No       | Timestamp when blob was last updated            |
| `expires_at`     | TIMESTAMPTZ | Yes      | Optional expiration date for automatic cleanup |
| `deleted_at`     | TIMESTAMPTZ | Yes      | Timestamp when blob was deleted (soft delete)  |