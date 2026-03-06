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

```bash
docker compose up --build -d
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

Exemplo com `curl`:

```bash
curl -X POST "http://127.0.0.1:3000/blob/upload" \
	-H "x-admin-token: <TOKEN_SECRET>" \
	-F "file=@./arquivo.pdf" \
	-F "bucket=docs" \
	-F "key=contratos/arquivo.pdf" \
	-F "public=false" \
	-F 'metadata={"origem":"manual"}'
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

## Seguranca

- assinatura HMAC-SHA512 com `nonce` e expiracao
- validacao de `bucket`/`key` e bloqueio de path traversal
- rate limit por IP nas rotas privadas
- headers `x-ratelimit-remaining` e `x-ratelimit-reset`

## Tester HTML

Existe um tester manual em:

- `api-tester.html`

Abra no navegador e configure a base URL (`http://127.0.0.1:3000`).

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
