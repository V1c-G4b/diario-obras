# Diário de Obra

API REST para registro digital de obras de construção civil. O engenheiro registra entradas diárias com progresso, custos, condições climáticas e fotos.

Projeto de estudo de **Go + Docker + Kubernetes + IaC + CI/CD** com arquitetura Ports and Adapters.

## Stack

- **Go** (Gin + GORM) — API REST
- **PostgreSQL** — Banco de dados
- **MinIO** — Object storage S3-compatible (fotos)
- **Swagger** — Documentação interativa da API
- **Prometheus + Grafana** — Métricas e dashboards (via kube-prometheus-stack)
- **Docker** — Multi-stage build com distroless (35MB)
- **Kubernetes** — Orquestração com Kind (local) e OpenShift
- **Helm** — Chart para deploy em Kubernetes e OpenShift
- **Terraform** — Provisionamento declarativo do cluster (namespaces, Helm releases, Tekton)
- **Tekton** — Pipeline CI/CD nativa em Kubernetes (build com Kaniko + deploy com Helm)

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

## Infraestrutura completa

```mermaid
flowchart TB
    Dev([Developer])
    TF[Terraform<br/>kind-local context]
    Dev -->|terraform apply| TF

    subgraph Host["Host (Docker)"]
        subgraph Kind["Kind Cluster — diario-obra"]
            direction TB

            TF -.->|provisiona| NS_APP
            TF -.->|helm_release monitoring| NS_MON
            TF -.->|helm_release diario-obra| NS_APP
            TF -.->|kubectl apply| NS_TEK

            subgraph NS_TEK["ns: tekton-pipelines"]
                TektonCtrl[Tekton Controller]
                RBAC[ClusterRole + Binding<br/>tekton-helm-role]
                Pipeline[Pipeline: deploy]
                TaskBuild[Task: build-image<br/>git-clone → Kaniko → ctr import]
                TaskDeploy[Task: deploy<br/>helm upgrade]
                Pipeline --> TaskBuild
                Pipeline --> TaskDeploy
                TaskBuild -->|ctr k8s.io| Containerd[(containerd do nó)]
            end

            subgraph NS_MON["ns: monitoring"]
                Prom[Prometheus Operator<br/>+ Prometheus]
                Graf[Grafana<br/>Ingress exposto]
                Prom --> Graf
            end

            subgraph NS_APP["ns: diario-obra"]
                Ing[Ingress nginx]
                APISvc[Service: api-service]
                API1[Pod: API replica 1]
                API2[Pod: API replica 2]
                HPA[HPA]
                SM[ServiceMonitor]
                PG[(PostgreSQL<br/>StatefulSet)]
                Minio[(MinIO<br/>Deployment)]
                DBSvc[Service: db-service]
                MinioSvc[Service: minio-service]

                Ing --> APISvc
                APISvc --> API1
                APISvc --> API2
                HPA -.->|scale| API1
                HPA -.->|scale| API2
                API1 --> DBSvc --> PG
                API2 --> DBSvc
                API1 --> MinioSvc --> Minio
                API2 --> MinioSvc
            end

            TaskDeploy -->|helm upgrade| NS_APP
            SM -.->|scrape| Prom
            Prom -.->|coleta /metrics| APISvc
        end
    end

    Client([Engenheiro na obra]) -->|HTTP :80| Ing
    Admin([Admin]) -->|HTTP| Graf
```

## Pipeline CI/CD (Tekton)

```mermaid
flowchart LR
    Trigger([TaskRun / PipelineRun])
    Trigger --> Build

    subgraph Pipeline: deploy
        Build[Task: build-image<br/>git-clone + Kaniko + ctr import]
        Deploy[Task: deploy<br/>helm upgrade]
        Build --> Deploy
    end

    Deploy --> Cluster[(Cluster Kind / OpenShift)]
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
| GET | `/api/v1/metrics` | Metricas Prometheus |
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

## Rodando com Kubernetes (Kind) — manifests

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

## Rodando com Helm

```bash
# Criar secrets (editar com credenciais reais)
cp chart/values-secrets.yaml chart/my-secrets.yaml
# editar chart/my-secrets.yaml

# Deploy local (Kind)
helm install diario-obra ./chart -n diario-obra --create-namespace -f chart/my-secrets.yaml

# Deploy OpenShift
helm install diario-obra ./chart -n diario-obra --create-namespace \
  -f chart/values-openshift.yaml \
  -f chart/my-secrets.yaml
```

## Provisionamento com Terraform

O Terraform provisiona, em um cluster Kind já existente (contexto `kind-local`):

- Namespace `diario-obra`
- Helm release do kube-prometheus-stack (Prometheus + Grafana) no namespace `monitoring`
- Helm release da aplicação usando o chart em `./chart`
- Instalação das Tekton Pipelines (via `kubectl apply`)

```bash
cd terraform

# Criar terraform.tfvars com as credenciais
cat > terraform.tfvars <<EOF
database_user          = "diario"
database_password      = "..."
storage_user           = "minio"
storage_password       = "..."
grafana_admin_password = "..."
EOF

terraform init
terraform apply
```

## Pipeline Tekton

A pipeline `deploy` roda dois Tasks sequenciais em um workspace compartilhado:

1. **build-image** — clona o repo, builda a imagem com Kaniko (`--no-push --tar-path`) e importa o tar no containerd do nó via `ctr` (namespace `k8s.io`).
2. **deploy** — roda `helm upgrade` do chart com as credenciais injetadas por parâmetro.

```bash
# RBAC para a ServiceAccount default do namespace tekton-pipelines
kubectl apply -f tekton/rbac.yaml

# Tasks e Pipeline
kubectl apply -f tekton/task-build.yaml
kubectl apply -f tekton/task-deploy.yaml
kubectl apply -f tekton/pipeline.yaml

# Executar apenas o build
kubectl apply -f tekton/taskrun-build.yaml
```

## Observabilidade

O middleware Prometheus coleta automaticamente:

- `http_requests_total` — contador por metodo, rota e status
- `http_request_duration_seconds` — histograma de latencia

Endpoint de métricas: `GET /api/v1/metrics`

No Kubernetes, o `ServiceMonitor` faz scraping automático a cada 15s para o Prometheus Operator (instalado via kube-prometheus-stack pelo Terraform). O Grafana fica exposto pelo `k8s/grafana-ingress.yaml`.

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
├── chart/                 # Helm chart (Kind + OpenShift)
├── terraform/             # IaC: namespace, Helm releases, Tekton
├── tekton/                # Pipeline CI/CD (build Kaniko + deploy Helm)
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
- [x] Fase 4 — Helm chart + métricas Prometheus
- [x] Fase 5 — Suporte a OpenShift
- [x] Fase 6 — IaC com Terraform (cluster, monitoring, Tekton)
- [x] Fase 7 — CI/CD com Tekton (build Kaniko + deploy Helm)
- [ ] Fase 8 — Triggers + webhook GitHub

## Autor

Victor Gabriel — [@v1c-g4b](https://github.com/v1c-g4b)

Acompanhe a jornada no [LinkedIn](https://linkedin.com/in/v1ctorg4briel) — #buildinpublic
