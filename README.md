
# vmmigrator-micro (Makefile Automated)

Microservice-based VM migration system built in Go, using **gRPC** for internal service communication
and **Gin** as the HTTP gateway.

This version includes **Makefile automation** for generating Protobuf and gRPC code.

---

# 🚀 Architecture Overview

```
Client (HTTP/JSON)
        ↓
API Gateway (Gin HTTP server)
        ↓ gRPC
Orchestrator Service
        ↓ gRPC
Provider Services (vmware, openstack, ...)
```

Each provider is a separate microservice running its own gRPC server.

---

# 📁 Project Structure

```
vmmigrator-micro/
├── Makefile
├── go.mod
├── proto/
│   └── vmmigrator.proto
├── services/
│   ├── api-gateway/
│   ├── orchestrator/
│   ├── provider-vmware/
│   └── provider-openstack/
└── infra/
    └── docker-compose.yml
```

---

# 🛠 Requirements

Install:

- `go` ≥ 1.20
- `protoc` ≥ 3.19
- `protoc-gen-go`
- `protoc-gen-go-grpc`

Install plugins:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Ensure `$GOPATH/bin` is in your `$PATH`.

---

# ⚙️ Proto Generation (Fully Automated)

Run:

```bash
make generate
```

This runs:

```
protoc --go_out=paths=source_relative:proto        --go-grpc_out=paths=source_relative:proto        proto/vmmigrator.proto
```

---

# ▶️ Running All Services

Follow this sequence:

### 1️⃣ Generate protobuf code

```bash
make generate
```

### 2️⃣ Run VMware provider service

```bash
go run ./services/provider-vmware/cmd
```

### 3️⃣ Run OpenStack provider service

```bash
go run ./services/provider-openstack/cmd
```

### 4️⃣ Run Orchestrator service

```bash
go run ./services/orchestrator/cmd
```

### 5️⃣ Run API Gateway

```bash
go run ./services/api-gateway/cmd
```


cd ~/vmmigrator

# ensure proto is generated
make generate

# rebuild images from scratch
docker compose -f infra/docker-compose.yml build --no-cache

# run
docker compose -f infra/docker-compose.yml up

---

# 📡 Test Migration API

Trigger migration:

```bash
curl -X POST http://localhost:8080/migrations   -H "Content-Type: application/json"   -d '{
    "source_provider": "vmware",
    "dest_provider":   "openstack",
    "vm_ids": ["vm-101"]
  }'
```

Check migration status:

```bash
curl http://localhost:8080/migrations/<job_id>
```

---

# 🧱 Notes

- Provider implementations are **stubbed**.
- Orchestrator uses **in-memory** job state (upgradeable to Postgres).
- gRPC internal communication is fully wired using generated protobufs.

---

# 📌 Next Steps I Can Build For You

✔ Add PostgreSQL job store  
✔ Add Redis/RabbitMQ worker queue  
✔ Implement real VMware + OpenStack API logic  
✔ Add Kubernetes deployment configs  
✔ Add TLS/mTLS between services  

Just tell me what you'd like next! 😊
