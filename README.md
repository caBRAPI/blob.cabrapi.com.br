# Tarefas do Projeto

## Concluídas
- [x] Healthcheck (`/health`)
- [x] Hello, World (`/`)
- [x] Logger colorido e padronizado
- [x] Rate limiter global por IP usando Redis
- [x] Integração com Redis (cache, rate limit)
- [x] Integração com Postgres (persistência)
- [x] Carregamento de variáveis de ambiente via `.env`
- [x] Estrutura inicial de rotas e middlewares
- [x] Configuração de storage local para blobs
- [x] Configuração de limites de upload e tipos permitidos

## Em andamento / A Fazer
- [ ] Implementar todas as rotas de blob (upload simples, multipart, download, delete, listagem, metadados)
- [ ] Implementar autenticação JWT nas rotas privadas
- [ ] Implementar CORS configurável
- [ ] Implementar expiração de blobs (públicos/privados)
- [ ] Implementar controle de tamanho máximo de upload
- [ ] Implementar validação de MIME types
- [ ] Implementar testes automatizados (unitários e integração)
- [ ] Implementar tratamento global de erros e respostas padronizadas
- [ ] Implementar logs de acesso e erros detalhados
- [ ] Implementar documentação automática (Swagger/OpenAPI)
- [ ] Adicionar exemplos de uso no README
- [ ] Adicionar CI/CD (lint, test, build)
- [ ] Melhorar mensagens de erro e feedback para o usuário
- [ ] Adicionar exemplos de configuração de ambiente
- [ ] Adicionar suporte a múltiplos backends de storage (futuro)

# API Blob

| Método | Rota | Privado | Descrição |
|--------|------|---------|-----------|
| `POST` | `/blob` |  `true` | Upload simples |
| `POST` | `/blob/multipart` |  `true` | Inicia upload multipart |
| `PUT` | `/blob/:id/part` |  `true` | Upload de parte |
| `POST` | `/blob/:id/complete` |  `true` | Finaliza multipart |
| `GET` | `/blob` |  `true` | Listar blobs |
| `GET` | `/blob/:id` |  `false` | Download |
| `HEAD` | `/blob/:id` |  `false` | Metadados |
| `DELETE` | `/blob/:id` |  `true` | Deletar |
| `GET` | `/health` |  `false` | Healthcheck |
| `GET` | `/` |  `false` | Hello, World |