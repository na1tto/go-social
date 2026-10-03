```text
  ____          ____             _       _
 / ___| ___    / ___|  ___   ___(_) __ _| |
| |  _ / _ \   \___ \ / _ \ / __| |/ _` | |
| |_| | (_) |   ___) | (_) | (__| | (_| | |
 \____|\___/   |____/ \___/ \___|_|\__,_|_|
```

<div align="center">

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)
![Prometheus](https://img.shields.io/badge/Prometheus-metrics-E6522C?logo=prometheus&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-API-85EA2D?logo=swagger&logoColor=black)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black)
![TypeScript](https://img.shields.io/badge/TypeScript-web-3178C6?logo=typescript&logoColor=white)
![Vite](https://img.shields.io/badge/Vite-web-646CFF?logo=vite&logoColor=white)

**[Português](README.md) · [English](README.en.md)**

</div>

---

# Go Social

API de uma rede social em Go, construída para evoluir a prática de desenvolvimento backend: autenticação, persistência, cache e observabilidade aplicados a um mesmo projeto.

O núcleo já permite cadastrar usuários, publicar, comentar, seguir pessoas e consultar um feed. O foco atual é tornar esses fluxos mais confiáveis e observar o comportamento da aplicação sob carga.

> **Em desenvolvimento.** Há um frontend inicial em `web/`, com confirmação de conta e uma página principal provisória. Login, feed e publicação pela interface ainda fazem parte da evolução do projeto.

<details>
<summary><kbd>Funcionalidades implementadas · clique para expandir</kbd></summary>

- Cadastro com bcrypt, ativação por e-mail e autenticação JWT.
- Permissões por autoria e papéis: `user`, `moderator` e `admin`.
- Criação, leitura, edição e exclusão de posts; criação de comentários.
- Seguidores e feed com paginação, ordenação, busca textual e filtro por tags.
- Cache opcional de usuários autenticados no Redis.
- Rate limiting em memória, logs estruturados e métricas Prometheus.
- Migrations SQL, testes com mocks, CI e ambiente de desenvolvimento em container.

</details>

## Conteúdo

- [01 · Início rápido](#01--início-rápido)
- [02 · Tecnologias e arquitetura](#02--tecnologias-e-arquitetura)
- [03 · Configuração](#03--configuração)
- [04 · Uso da API](#04--uso-da-api)
- [05 · Observabilidade e decisões](#05--observabilidade-e-decisões)
- [06 · Limitações e próximos passos](#06--limitações-e-próximos-passos)
- [07 · Desenvolvimento](#07--desenvolvimento)

---

## 01 · Início rápido

O caminho recomendado é o **Dev Container**. Requer Git, Docker com Compose e um editor compatível, como VS Code com Dev Containers ou Zed. No Windows, use Docker Desktop com WSL 2.

Na raiz do projeto:

```sh
cp .env.example .env
```

No PowerShell, use `Copy-Item .env.example .env`. Abra o projeto no Dev Container e execute no terminal dele:

```sh
make migrate-up
make dev
```

A API fica em **http://localhost:8080**. Verifique com `curl http://localhost:8080/v1/health`; essa rota é pública atualmente. O container fornece Go, Air, Swag e Migrate; as migrations são executadas manualmente.

Para testar cadastro e ativação, configure o Mailtrap conforme a [seção de configuração](#03--configuração). O guia [DEVCONTAINER.md](DEVCONTAINER.md) detalha abertura no editor, ferramentas e solução de problemas.

<details>
<summary><kbd>Executar diretamente na máquina · alternativa</kbd></summary>

Instale Go compatível com [go.mod](go.mod), Docker Compose e o CLI Migrate com suporte a PostgreSQL. Make e Air são necessários para os atalhos de desenvolvimento.

No `.env`, troque os endereços de rede do Docker pelos publicados na máquina:

```dotenv
DB_ADDR=postgres://admin:adminpassword@localhost:15432/gosocial?sslmode=disable
REDIS_ADDR=localhost:16379
```

Depois:

```sh
docker compose up -d
make migrate-up
go run ./cmd/api
```

A API carrega `.env` automaticamente; o Makefile também lê esse arquivo. O Compose da raiz sobe a infraestrutura, mas não inicia a API.

</details>

---

## 02 · Tecnologias e arquitetura

| Camada | Ferramentas e responsabilidade |
|---|---|
| API | Go · Chi · go-playground/validator |
| Persistência | PostgreSQL 16 · SQL com `database/sql` e `lib/pq` · golang-migrate |
| Autenticação | JWT · bcrypt · papéis e verificação de autoria |
| Cache e limite de requisições | Redis 7 para usuários · janela fixa em memória para rate limiting |
| E-mail | Mailtrap SMTP sandbox ativo · implementação SendGrid disponível, sem uso no bootstrap |
| Observabilidade | Zap · Prometheus · expvar |
| Interface inicial | React 19 · TypeScript · Vite · React Router |
| Desenvolvimento | Docker Compose · Dev Containers · Air · Swag · GitHub Actions |

A API é um único serviço Go organizado em camadas. Middlewares tratam autenticação e aspectos comuns das requisições; handlers validam a entrada e coordenam operações; `internal/store` concentra o acesso ao banco. Ainda não há uma camada de serviços de negócio separada.

```text
Cliente HTTP / web (React)
          |
          v
  Chi + middlewares
          |
          v
  Handlers (cmd/api) ------> Mailtrap (ativação)
          |
          +---------------> Redis (cache de usuários)
          |
          v
  Store (internal/store)
          |
          v
      PostgreSQL

  API ----> logs Zap
  Prometheus ----> API /metrics
```

O modelo relaciona usuários e papéis, convites de ativação, posts, comentários e seguidores. Posts têm versão para controle de atualização concorrente; o banco usa `CITEXT`, `pg_trgm` e índices para busca e tags.

| Diretório | Onde procurar |
|---|---|
| `cmd/api/` | Rotas, handlers, middlewares e respostas HTTP |
| `internal/store/` | Consultas SQL, interfaces e cache |
| `internal/auth/`, `internal/mailer/` | JWT e envio de e-mail |
| `internal/observability/`, `internal/ratelimiter/` | Métricas e controle de requisições |
| `cmd/migrate/`, `internal/db/` | Migrations, conexão e seed |
| `web/` | Interface inicial de confirmação de conta |
| `docs/` | Swagger gerado e decisões em `docs/adr/` |
| `.devcontainer/`, `.github/workflows/` | Ambiente de desenvolvimento e CI |

---

## 03 · Configuração

Use [.env.example](.env.example) como ponto de partida. Ele foi preparado para a rede do Dev Container; `.env` é ignorado pelo Git.

| Variável | Uso atual |
|---|---|
| `ADDR` / `ENV` | Endereço HTTP (`:8080`) e ambiente |
| `DB_ADDR` | PostgreSQL: `db:5432` no container ou `localhost:15432` na máquina |
| `DB_USER`, `DB_PASSWORD`, `POSTGRES_DB` | Credenciais e banco criados pelo Compose |
| `REDIS_ENABLED` / `REDIS_ADDR` | Cache: `redis:6379` no container ou `localhost:16379` na máquina |
| `AUTH_TOKEN_SECRET` | Segredo JWT; adicione um valor próprio ao `.env` |
| `AUTH_BASIC_USER` / `AUTH_BASIC_PASS` | Protegem `/v1/debug/vars`; padrão local `admin:admin` |
| `FROM_EMAIL` / `MAILTRAP_*` | Remetente e configuração de e-mail |
| `FRONTEND_URL` | Base dos links de ativação; padrão `http://localhost:5173` |
| `RATE_LIMITER_ENABLED` / `RATELIMITER_REQUESTS_COUNT` | Limite local; padrão de 20 requisições por janela de 3 segundos |
| `EXTERNAL_URL` | Alimenta o host do Swagger; integração ainda requer revisão |

**E-mail:** preencha `FROM_EMAIL`, `MAILTRAP_API_KEY`, `MAILTRAP_USERNAME` e `MAILTRAP_PASSWORD`. O construtor exige a API key, mas o envio atual usa usuário e senha do SMTP sandbox. Se o envio falhar, o cadastro tenta desfazer a criação do usuário.

**CORS:** `CORS_ALLOWED_ORIGIN` está no exemplo, mas o código lê `FRONTEND_ORIGIN` e atualmente aceita qualquer origem via callback. A lista de métodos também não inclui `PATCH`. A correção está no roadmap.

**Serviços locais:** API `:8080` · PostgreSQL `:15432` · Redis `:16379` · Redis Commander `:8082` · Prometheus `:9090`.

<details>
<summary><kbd>Interface de ativação · opcional</kbd></summary>

Com Node compatível com [web/.nvmrc](web/.nvmrc) instalado, execute em outro terminal:

```sh
cd web
npm ci
npm run dev
```

A interface chama `http://localhost:8080/v1` por padrão. Para trocar a API, configure `VITE_API_URL` no ambiente do Vite. Mantenha `FRONTEND_URL` alinhada à porta exibida pelo frontend; o link do e-mail abre `/confirm/:token`.

</details>

---

## 04 · Uso da API

Base local: `http://localhost:8080/v1`. O fluxo é **cadastro → ativação → emissão de JWT → rotas autenticadas**.

| Método | Rota relativa a `/v1` | Acesso / finalidade |
|---|---|---|
| `POST` | `/authentication/user` | Público · recebe `username`, `email` e `password` |
| `PUT` | `/users/activate/{token}` | Público · ativa a conta |
| `POST` | `/authentication/token` | Público · recebe `email` e `password`, retorna JWT |
| `GET` | `/users/{userId}/` | JWT · consulta usuário |
| `PUT` | `/users/{userId}/follow` ou `/users/{userId}/unfollow` | JWT · segue ou deixa de seguir |
| `GET` | `/users/feed` | JWT · feed paginado |
| `POST` | `/posts/` | JWT · cria post com `title`, `content` e `tags` |
| `GET` / `POST` | `/posts/{postId}/` | JWT · lê post com comentários / comenta com `content` |
| `PATCH` / `DELETE` | `/posts/{postId}/` | JWT + autoria ou papel · edita / exclui |

Envie `Authorization: Bearer <token>` nas rotas protegidas. O autor pode editar e excluir seus posts; moderadores podem editar posts de outros usuários e administradores também podem excluí-los.

O feed aceita `limit` (1–20), `offset` (≥ 0), `sort` (`asc`/`desc`), `tags` (até 3, separadas por vírgula) e `search` (até 100 caracteres). `since` e `until` são interpretados, mas ainda não filtram a consulta SQL.

Respostas JSON usam `{"data": ...}` em sucesso e `{"error": "..."}` em falhas; exclusões podem retornar `204` sem corpo.

**Referência:** [Swagger YAML](docs/swagger.yaml) e [Swagger JSON](docs/swagger.json). A UI fica em `/v1/swagger/index.html`, mas a URL do schema montada pela aplicação ainda precisa de correção; consulte os arquivos se ela não carregar.

---

## 05 · Observabilidade e decisões

Os logs incluem `request_id`, método, rota, status, duração e estado do contexto. As métricas das rotas `/v1` registram volume, duração e requisições em andamento.

| Endpoint | Acesso atual |
|---|---|
| `GET /v1/health` | Público · estado, ambiente e versão |
| `GET /metrics` | Público · coleta Prometheus |
| `GET /v1/debug/vars` | Basic Auth · runtime e pool de conexões |

O [Prometheus local](internal/observability/prometheus/prometheus.yml) consulta `host.docker.internal:8080` a cada 15 segundos. Seu painel fica em **http://localhost:9090**.

### ADR 0001 · Cancelamentos HTTP

O [ADR 0001 — Tratamento de cancelamentos de requisições HTTP](docs/adr/0001-http-request-cancellation.md), aceito em **03/10/2026**, registra a decisão motivada pelos testes de carga: cancelamentos reconhecidos não devem virar erros internos nem sucessos artificiais.

A decisão prevê log em nível Info sem resposta JSON, correlação por `request_id` e preservação de status já escrito. Sem status escrito e com contexto cancelado, os logs usam `0` e as métricas usam `status="canceled"`. Esse `0` existe apenas na observabilidade; não é um código HTTP enviado ao cliente. Expiração de prazo fica fora do escopo do ADR.

**Alinhamento pendente:** em `errors.go`, a condição atual verifica apenas `errors.Is(err, context.Canceled)`; falta exigir também o cancelamento do contexto HTTP, conforme a decisão documentada.

---

## 06 · Limitações e próximos passos

O projeto já tem testes com mocks, CI, rate limiting e Dockerfile de runtime. A próxima etapa é ampliar e consolidar esses recursos. A tabela organiza pendências verificadas no código e propostas de evolução, sem compromisso de prazo.

| Limitação atual | Próximo passo proposto |
|---|---|
| Cadastro depende do Mailtrap sandbox síncrono; erro de inicialização do mailer é ignorado | Validar o bootstrap e evoluir envio, retentativas e recuperação do cadastro |
| Token de ativação também é devolvido no JSON; não há refresh token | Revisar o contrato de ativação e evoluir o ciclo de autenticação |
| CORS permissivo, configuração divergente e ausência de `PATCH` | Unificar configuração, restringir origens e cobrir os métodos usados |
| Filtros temporais do feed não chegam ao SQL | Implementar `since`/`until` e validar com testes de integração |
| Cache habilitado propaga falhas do Redis | Definir comportamento de contingência e invalidação |
| Rate limiter vive em memória por processo e não remove entradas antigas | Definir limpeza e estratégia para múltiplas instâncias |
| Swagger tem problemas na URL do schema e nos comandos de geração | Alinhar rotas, metadados e geração entre ambientes |
| Tratamento de cancelamentos ainda diverge do ADR | Completar a condição e ampliar testes de cancelamento e timeout |
| Frontend cobre apenas a confirmação de conta | Evoluir cadastro, login, feed e publicação |
| Dockerfile usa `cmd/api/*.go`, incluindo arquivos de teste no build | Corrigir e validar a imagem; ampliar testes de integração e validação de deploy |

Antes de publicar uma instância, também é necessário substituir os segredos de exemplo e definir a exposição de `/metrics` e das rotas de diagnóstico. A configuração atual é voltada ao desenvolvimento.

---

## 07 · Desenvolvimento

No Dev Container, com `.env` criado:

```sh
make test                            # testes Go
go vet ./...                         # análise estática
make migrate-up                      # aplica migrations
make migration name=nome_da_mudanca   # cria migration SQL
make dev                             # API com live reload
```

O [workflow de auditoria](.github/workflows/audit.yaml) já executa verificação de dependências, build, `go vet`, Staticcheck e testes com `-race`. O [CHANGELOG](CHANGELOG.md) registra as releases; os [ADRs](docs/adr/) explicam decisões de arquitetura.

<details>
<summary><kbd>Banco de demonstração e geração de documentação</kbd></summary>

`make seed` gera 100 usuários, 200 posts e 500 comentários, com senha de demonstração `123123`. O executável de seed lê `DB_ADDR` do ambiente: fora do Dev Container, exporte essa variável antes de executá-lo, pois o seed não carrega `.env` por conta própria.

`make migrate-down` reverte uma migration. `make reset-db` apaga o schema, recria e popula o banco; use somente em um banco descartável.

O Makefile oferece `make gen-docs` e `make gen-docs-win`. O primeiro ainda combina diretório de busca e caminho do arquivo principal de forma inconsistente. Uma alternativa direta, na raiz, é:

```sh
swag init -g ./api/main.go -d cmd,internal --parseDependency --parseInternal
swag fmt
```

Revise os arquivos gerados antes de incluí-los em um commit.

</details>

---

**Documentação revisada:** 03/10/2026 · **Idiomas:** [Português](README.md) / [English](README.en.md)
