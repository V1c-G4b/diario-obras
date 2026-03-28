# Diário de Obra

API REST para registro digital de obras de construção civil. O engenheiro registra entradas diárias com progresso, custos, condições climáticas e fotos.

Projeto de estudo de **Go + Docker + Kubernetes** com arquitetura Ports and Adapters.

## Stack

- **Go** (Gin + GORM) — API REST
- **PostgreSQL** — Banco de dados
- **MinIO** — Object storage S3-compatible (fotos)
- **Swagger** — Documentação interativa da API
- **Docker** — Multi-stage build com distroless (35MB)
- **Kubernetes** — Orquestração com Kind (local)

## Arquitetura

```mermaid
flowchart TB
    Client([Engenheiro na obra])
    Client --> Ingress

    subgraph Kubernetes Cluster
        Ingress[Ingress nginx :80]
        Ingress --> APIService[Service: api-service]
        APIService --> API1[Pod: API replica 1]
        APIService --> API2[Pod: API replica 2]

        API1 --> DBService[Service: db-service]
        API1 --> MinIOService[Service: minio-service]
        API2 --> DBService
        API2 --> MinIOService

        DBService --> PG[(PostgreSQL)]
        MinIOService --> MinIO[(MinIO S3)]
    end
```

## Arquitetura interna — Ports and Adapters

```mermaid
flowchart LR
    subgraph Adapters
        Handler[handler/\nGin HTTP]
        Repo[repository/\nGORM]
        Storage[storage/\nMinIO SDK]
    end

    subgraph Application
        ObraService[ObraService]
        EntradaService[EntradaService]
        FotoService[FotoService]
        ResponsavelService[ResponsavelService]
    end

    subgraph Domain
        Entity[entity/\nObra, Entrada, Foto]
        Port[port/\nInterfaces]
    end

    Handler --> ObraService
    Handler --> EntradaService
    Handler --> FotoService
    Handler --> ResponsavelService
    ObraService --> Port
    EntradaService --> Port
    FotoService --> Port
    ResponsavelService --> Port
    Repo -.->|implements| Port
    Storage -.->|implements| Port
    Port --> Entity
```

## Modelo de dados

```mermaid
erDiagram
    OBRA ||--o{ ENTRADA : tem
    ENTRADA ||--o{ FOTO : contem
    OBRA ||--o{ RESPONSAVEL : atribuida

    OBRA {
        uuid id PK
        string nome
        string endereco
        float progresso
        date data_inicio
        date data_estimada
        float gasto_total
    }

    ENTRADA {
        uuid id PK
        uuid obra_id FK
        uuid responsavel_id FK
        string observacoes
        string etapa
        float progresso_etapa
        float custo_dia
        int quantidade_trabalhadores
        int condicoes_climaticas
        bool paralisacao
        datetime criado_em
    }

    FOTO {
        uuid id PK
        uuid entrada_id FK
        string url_s3
        string descricao
        datetime criado_em
    }

    RESPONSAVEL {
        uuid id PK
        string nome
        string email
        string cargo
    }
```

## Endpoints

| Metodo | Rota | Descricao |
|--------|------|-----------|
| GET | `/ping` | Health check |
| GET | `/swagger/*any` | Documentacao Swagger |
| | | |
| **Obras** | | |
| POST | `/api/v1/obras` | Criar obra |
| GET | `/api/v1/obras` | Listar obras |
| GET | `/api/v1/obras/:id` | Buscar obra |
| PUT | `/api/v1/obras/:id` | Atualizar obra |
| DELETE | `/api/v1/obras/:id` | Deletar obra |
| | | |
| **Entradas** | | |
| POST | `/api/v1/obras/:id/entradas` | Criar entrada |
| GET | `/api/v1/obras/:id/entradas` | Listar entradas da obra |
| GET | `/api/v1/entradas/:id` | Buscar entrada |
| DELETE | `/api/v1/entradas/:id` | Deletar entrada |
| | | |
| **Fotos** | | |
| POST | `/api/v1/entradas/:id/fotos` | Upload foto (multipart) |
| GET | `/api/v1/entradas/:id/fotos` | Listar fotos da entrada |
| DELETE | `/api/v1/entradas/:id/fotos/:fotoId` | Deletar foto |
| | | |
| **Responsaveis** | | |
| POST | `/api/v1/responsaveis` | Criar responsavel |
| GET | `/api/v1/responsaveis` | Listar responsaveis |
| GET | `/api/v1/responsaveis/:id` | Buscar responsavel |
| PUT | `/api/v1/responsaveis/:id` | Atualizar responsavel |
| DELETE | `/api/v1/responsaveis/:id` | Deletar responsavel |

## Rodando com Docker Compose

```bash
cp .env.example .env
docker compose up -d
curl http://localhost:8080/ping
```

## Rodando com Kubernetes (Kind)

```bash
# Criar cluster
kind create cluster --name diario-obra --config kind-config.yaml

# Build e carregar imagem
docker build -t diario-obras-api:latest .
kind load docker-image diario-obras-api:latest --name diario-obra

# Aplicar manifests
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/minio.yaml
kubectl apply -f k8s/api.yaml
kubectl apply -f k8s/ingress.yaml
kubectl apply -f k8s/hpa.yaml

# Instalar Ingress controller
kubectl apply -f https://kind.sigs.k8s.io/examples/ingress/deploy-ingress-nginx.yaml

# Testar
curl http://localhost/ping
```

## Estrutura do projeto

```
├── cmd/api/main.go
├── internal/
│   ├── domain/
│   │   ├── entity/        # Obra, Entrada, Foto, Responsavel
│   │   └── port/          # Interfaces (repositories, storage)
│   ├── application/       # Services (use cases)
│   └── adapter/
│       ├── handler/       # Gin HTTP handlers + router
│       ├── repository/    # GORM implementations
│       └── storage/       # MinIO S3 implementation
├── docs/                  # Swagger (gerado pelo swag)
├── k8s/                   # Kubernetes manifests
├── Dockerfile             # Multi-stage distroless (35MB)
├── docker-compose.yml
└── kind-config.yaml
```

## Docker — progressão de tamanho

| Imagem | Tamanho |
|--------|---------|
| golang:alpine (runtime) | ~600MB |
| Alpine multi-stage | 75MB |
| Distroless multi-stage | **35MB** |

## Roadmap

- [x] Fase 1 — Docker multi-stage + distroless
- [x] Fase 2 — MinIO upload de fotos
- [x] Fase 3 — Kubernetes com Kind
- [ ] Fase 4 — Helm chart + OpenShift + CI/CD

## Autor

Victor Gabriel — [@v1c-g4b](https://github.com/v1c-g4b)

Acompanhe a jornada no [LinkedIn](https://linkedin.com/in/v1ctorg4briel) — #buildinpublic
