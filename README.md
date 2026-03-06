# Blob API em Go

Este projeto implementa as rotas de manipulação de arquivos blob em Go, conforme especificação:

- POST   /blob                   - Requer autenticação (Admin Token)
- POST   /blob/multipart         - Requer autenticação (Admin Token)
- PUT    /blob/:id/part          - Requer autenticação (Admin Token)
- POST   /blob/:id/complete      - Requer autenticação (Admin Token)
- GET    /blob                   - Requer autenticação (Admin Token)
- GET    /blob/:id               - Público (sem autenticação) ou via assinatura
- HEAD   /blob/:id               - Público (sem autenticação) ou via assinatura
- DELETE /blob/:id               - Requer autenticação (Admin Token)
- GET    /health                 - Público (sem autenticação)

## Como rodar

1. Instale o Go (versão 1.20+)
2. Instale as dependências:
   ```bash
   go get github.com/gorilla/mux
   ```
3. Defina a variável de ambiente `ADMIN_TOKEN` para autenticação admin.
4. Rode o servidor:
   ```bash
   go run main.go
   ```

## Estrutura
- Todas as rotas e handlers estão em `main.go`.
- Middlewares e handlers podem ser separados em arquivos conforme necessidade.

## Observação
Os handlers estão prontos para serem implementados com lógica real de upload, download, etc.

---
Qualquer dúvida ou ajuste, só pedir!
