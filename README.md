# blob.cabrapi.com.br

Servico HTTP para upload, listagem, visualizacao e remocao de arquivos com:

- metadados em SQLite (`Prisma`)
- arquivos em disco local (`data/blob-storage`)
- URL assinada para blobs privados
- limite de taxa para rotas sensiveis

## Requisitos

- Node.js 20+
- npm
- Docker (opcional)

## Configuracao

Use o arquivo de exemplo:

```bash
cp .env.example .env
```

No Windows (PowerShell):

```powershell
Copy-Item .env.example .env
```

Ajuste obrigatoriamente:

- `TOKEN_SECRET`

Token administrativo:

- use `TOKEN_SECRET` em `x-admin-token` (ou `Authorization: Bearer ...`) para operacoes protegidas
- operacoes protegidas: upload, delete e listagem de privados

## Executar Local (Passo a Passo)

```bash
npm install
npm run db:prepare
npm run dev
```

API local:

- `http://127.0.0.1:3000`

Healthcheck rapido:

```bash
curl http://127.0.0.1:3000/health
```

## Docker (volume em `data/`)

Antes de subir o container, gere seu `.env`:

```bash
cp .env.example .env
```

```bash
docker compose up --build -d
```

Verificar status e healthcheck:

```bash
docker compose ps
docker compose logs -f blob
```

Persistencia:

- `./data/blob.db`
- `./data/blob-storage/...`

## Endpoints

Base URL: `http://127.0.0.1:3000`

- `GET /health`
- `POST /blob/upload`
- `GET /blob?page=1&pageSize=20&bucket=docs&public=true`
- `GET /blob/:id/sign?ttl=120`
- `GET /blob/:id?exp=<exp>&n=<nonce>&sig=<sig>`
- `DELETE /blob/:id`

### Upload (`POST /blob/upload`)

`multipart/form-data`:

- `file` (obrigatorio)
- `bucket` (opcional)
- `key` (opcional)
- `public` (opcional: `true|false`)
- `metadata` (opcional: JSON string)
- `expiresAt` (opcional: ISO 8601 ou timestamp)

Exemplo com `curl`:

```bash
curl -X POST "http://127.0.0.1:3000/blob/upload" \
	-H "x-admin-token: <TOKEN_SECRET>" \
	-F "file=@./arquivo.pdf" \
	-F "bucket=docs" \
	-F "key=contratos/arquivo.pdf" \
	-F "public=false" \
	-F 'metadata={"origem":"manual"}' \
	-F "expiresAt=2026-12-31T23:59:59Z"
```

### Visualizar Blob Publico

```bash
curl -L "http://127.0.0.1:3000/blob/<id>" --output arquivo.bin
```

### Download Privado (URL Assinada)

1. Chame `GET /blob/:id/sign?ttl=120`
2. Use a `url` retornada para baixar o arquivo

Exemplo:

```bash
curl "http://127.0.0.1:3000/blob/<id>/sign?ttl=120"
curl -L "http://127.0.0.1:3000/blob/<id>?exp=<exp>&n=<nonce>&sig=<sig>" --output arquivo.bin
```

### Listagem

Comportamento de seguranca:

- sem token administrativo, a API lista apenas blobs publicos
- para listar privados, envie `x-admin-token` com `TOKEN_SECRET`

```bash
curl "http://127.0.0.1:3000/blob?page=1&pageSize=10"
curl "http://127.0.0.1:3000/blob?page=1&pageSize=10&bucket=docs"
curl "http://127.0.0.1:3000/blob?page=1&pageSize=10&public=true"
curl -H "x-admin-token: <TOKEN_SECRET>" "http://127.0.0.1:3000/blob?page=1&pageSize=10&public=false"
```

### Remocao

```bash
curl -X DELETE -H "x-admin-token: <TOKEN_SECRET>" "http://127.0.0.1:3000/blob/<id>"
```

### Expiração de Blobs Privados

- Por padrão, arquivos privados não expiram automaticamente.
- Você pode definir uma data de expiração ao fazer upload usando o campo `expiresAt` (ISO 8601 ou timestamp).
- Após a expiração (`expiresAt`), nem mesmo o admin pode acessar o arquivo.
- URLs assinadas (GET /blob/:id/sign) sempre respeitam o TTL solicitado, mas nunca ultrapassam a data de expiração se definida.
- O admin pode baixar arquivos privados a qualquer momento (sem TTL/assinatura), exceto se expirados.

#### Exemplo de upload com expiração:

```bash
curl -X POST "http://127.0.0.1:3000/blob/upload" \
	-H "x-admin-token: <TOKEN_SECRET>" \
	-F "file=@./arquivo.pdf" \
	-F "expiresAt=2026-12-31T23:59:59Z"
```

#### Exemplo de download como admin (sem TTL):

```bash
curl -L -H "x-admin-token: <TOKEN_SECRET>" "http://127.0.0.1:3000/blob/<id>" --output arquivo.bin
```

#### Exemplo de download externo (URL assinada, respeita expiração):

```bash
curl "http://127.0.0.1:3000/blob/<id>/sign?ttl=3600"
curl -L "http://127.0.0.1:3000/blob/<id>?exp=<exp>&n=<nonce>&sig=<sig>" --output arquivo.bin
```

Se o arquivo estiver expirado, qualquer tentativa de acesso retorna erro 410 (Gone).

## Seguranca

- assinatura HMAC-SHA512 com `nonce` e expiracao
- validacao de `bucket`/`key` e bloqueio de path traversal
- rate limit por IP nas rotas privadas
- headers `x-ratelimit-remaining` e `x-ratelimit-reset`

## Scripts

```bash
npm run dev
npm start
npm run db:generate
npm run db:push
npm run db:prepare
npm run typecheck
npm run check
```

## Troubleshooting

- Erro de schema/DB:
	- execute `npm run db:prepare`.
- Erro de assinatura em blob privado:
	- gere novamente via `/blob/:id/sign` e use os parametros `exp`, `n`, `sig` sem alteracao.
- Upload rejeitado por MIME:
	- ajuste `ALLOWED_MIME_TYPES` no `.env` ou remova a variavel para aceitar qualquer MIME.
