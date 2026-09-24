# คู่มือระบบออกแบบระบบ IoT Flow Management Module

## พร้อม Kafka, InfluxDB, Elasticsearch, Socket.IO, Logging, LLM, AI Embedding, Blockchain, Vector Database และระบบ Logging Management

---

## สารบัญ

1. [ภาพรวมระบบ](#1-ภาพรวมระบบ)
2. [โครงสร้าง Module](#2-โครงสร้าง-module)
3. [Domain Layer](#3-domain-layer)
4. [Application Layer](#4-application-layer)
5. [Infrastructure Layer](#5-infrastructure-layer)
6. [Interface Layer](#6-interface-layer)
7. [Background Workers / Consumers](#7-background-workers--consumers)
8. [Database Migrations](#8-database-migrations)
9. [Workflow Diagrams](#9-workflow-diagrams)
10. [การติดตั้งและใช้งาน](#10-การติดตั้งและใช้งาน)
11. [Frontend Integration Guide](#11-frontend-integration-guide)
12. [สรุปและแนวทางพัฒนาเพิ่มเติม](#12-สรุปและแนวทางพัฒนาเพิ่มเติม)

---

## 1. ภาพรวมระบบ

ระบบ IoT Flow Management Module เป็นแพลตฟอร์มประมวลผลข้อมูลแบบเรียลไทม์ สำหรับอุปกรณ์ IoT ด้วยสถาปัตยกรรม Event‑Driven และ Clean Architecture โดยมีเป้าหมายหลัก:

- รองรับข้อมูลปริมาณสูง (High Throughput)
- แยกส่วนการประมวลผล (Decoupling) ด้วย Kafka
- จัดเก็บข้อมูลอนุกรมเวลา (Time‑series) ด้วย InfluxDB
- ค้นหาและจัดเก็บ Logs ด้วย Elasticsearch
- แสดงผลแบบ Real‑time ผ่าน Socket.IO
- เสริมความอัจฉริยะด้วย LLM (สำหรับสรุป, ตอบคำถาม, รายงาน)
- ใช้ AI Embedding เพื่อค้นหาเชิงความหมาย (Semantic Search) และ RAG
- ใช้ Blockchain เพื่อบันทึกหลักฐานที่ไม่สามารถแก้ไขได้ (Immutable Audit Trail)
- ใช้ Vector Database เฉพาะทาง (Milvus/Qdrant) สำหรับการค้นหาเวกเตอร์ประสิทธิภาพสูง
- ระบบ Logging Management ที่เก็บผลเฉลี่ยรายปีและลบข้อมูลเก่าโดยอัตโนมัติ

### 1.1 สถาปัตยกรรมหลัก

```
External Input (MQTT/HTTP)
       │
       ▼
┌──────────────┐
│  Flow Engine │  (Goroutines, Channels)
└──────────────┘
       │
       ▼
┌──────────────┐
│Node Executors│  (Filter, Function, LLM, Embedding, Blockchain, etc.)
└──────────────┘
       │
       ▼
   Kafka Topics
  ┌────────────────────────────────────────────────────────────────┐
  │ flow.data     │──► Consumer Group ──► InfluxDB, ES, Socket.IO │
  │ flow.alarm    │──► Consumer Group ──► ES, Socket.IO, LLM      │
  │ flow.llm.request    │──► LLM Consumer ──► LLM Service         │
  │ flow.embedding.request │──► Embedding Consumer ──► Vector DB  │
  │ flow.blockchain.request │──► Blockchain Consumer ──► Network  │
  └────────────────────────────────────────────────────────────────┘
       │
       ▼
  Admin API (HTTP) + Socket.IO Dashboard
```

### 1.2 ส่วนประกอบหลัก

| ส่วนประกอบ | บทบาท |
|------------|--------|
| **Kafka** | Message Bus หลัก แยก Producer และ Consumer |
| **InfluxDB** | เก็บ Time‑series data (ค่า sensor, metrics) |
| **Elasticsearch** | เก็บ Logs, Audit, และ Vector (dense_vector) |
| **Vector Database** (Milvus/Qdrant) | เก็บ Embedding vectors สำหรับ similarity search |
| **Socket.IO** | Real‑time Dashboard |
| **LLM** | สรุป, ตอบคำถาม, สร้างรายงานอัตโนมัติ |
| **AI Embedding** | สร้างเวกเตอร์จากข้อมูล sensor/ข้อความ |
| **Blockchain** | เก็บ hash ของข้อมูลสำคัญใน Smart Contract |
| **Aggregator** | คำนวณสถิติรายปี และลบ raw data เก่า |
| **Cron Scheduler** | เรียกใช้งาน Aggregator และ Cleanup เป็นระยะ |

---

## 2. โครงสร้าง Module

โครงสร้างโค้ดภายใน `internal/modules/iotflow/` (Clean Architecture):

```
internal/modules/iotflow/
├── domain/
│   ├── entity/
│   │   ├── flow.go
│   │   ├── node.go
│   │   └── wire.go
│   ├── value_object/
│   │   ├── node_type.go          // รวม NodeTypeLLM, NodeTypeEmbedding, NodeTypeBlockchain
│   │   ├── node_status.go
│   │   ├── flow_status.go
│   │   ├── data_payload.go
│   │   └── port.go
│   ├── repository/
│   │   ├── flow_repository.go
│   │   └── flow_cache_repository.go
│   ├── service/
│   │   ├── flow_engine.go
│   │   ├── node_executor.go
│   │   └── flow_validator.go
│   └── errors/
│       └── errors.go
├── application/
│   ├── create_flow.go
│   ├── deploy_flow.go
│   ├── start_flow.go
│   ├── stop_flow.go
│   ├── get_flow.go
│   ├── list_flows.go
│   ├── update_flow.go
│   ├── delete_flow.go
│   ├── add_node.go
│   ├── add_wire.go
│   ├── remove_node.go
│   ├── remove_wire.go
│   ├── generate_aggregates.go
│   └── dto.go
├── infrastructure/
│   ├── persistence/
│   │   └── postgres/
│   │       ├── models.go
│   │       └── flow_repo_impl.go
│   ├── cache/
│   │   └── redis/
│   │       ├── redis_client.go
│   │       ├── flow_cache_repo_impl.go
│   │       ├── rate_limiter.go
│   │       └── flow_instance_store.go
│   ├── messaging/
│   │   ├── kafka_client.go
│   │   ├── producer.go
│   │   └── consumer.go
│   ├── timeseries/
│   │   ├── influx_client.go
│   │   └── writer.go
│   ├── search/
│   │   ├── es_client.go
│   │   ├── indexer.go
│   │   └── vector_store.go
│   ├── vectordb/
│   │   ├── interface.go
│   │   ├── milvus_client.go
│   │   ├── consumer.go
│   │   └── index_manager.go
│   ├── llm/
│   │   ├── client.go
│   │   └── consumer.go
│   ├── embedding/
│   │   ├── client.go
│   │   └── consumer.go
│   ├── blockchain/
│   │   ├── client.go
│   │   └── consumer.go
│   ├── aggregator/
│   │   ├── aggregator.go
│   │   ├── scheduler.go
│   │   └── cleanup.go
│   ├── engine/
│   │   ├── runtime.go
│   │   ├── node_runner.go
│   │   └── builtin_nodes/
│   │       ├── mqtt_input.go
│   │       ├── filter.go
│   │       ├── function.go
│   │       ├── http_output.go
│   │       ├── debug.go
│   │       ├── delay.go
│   │       ├── schedule_trigger.go
│   │       ├── kafka_output.go
│   │       ├── influx_output.go
│   │       ├── es_output.go
│   │       ├── llm_query.go
│   │       ├── embedding.go
│   │       └── blockchain_record.go
│   ├── mqtt/
│   │   └── mqtt_adapter.go
│   ├── logging/
│   │   └── logger.go
│   └── helpers/
│       ├── timezone.go
│       ├── alarm_validator.go
│       └── utils.go
└── interfaces/
    ├── http/
    │   ├── flow_handler.go
    │   └── routes.go
    ├── middleware/
    │   └── rate_limit.go
    ├── socketio/
    │   ├── server.go
    │   └── handlers.go
    └── websocket/
        └── flow_monitor.go
```

---

## 3. Domain Layer

### 3.1 Value Objects

```go
// domain/value_object/node_type.go
const (
    NodeTypeMQTTInput       NodeType = "mqtt_input"
    NodeTypeFilter          NodeType = "filter"
    NodeTypeFunction        NodeType = "function"
    NodeTypeHTTPOutput      NodeType = "http_output"
    NodeTypeDebug           NodeType = "debug"
    NodeTypeDelay           NodeType = "delay"
    NodeTypeScheduleTrigger NodeType = "schedule_trigger"
    NodeTypeKafkaOutput     NodeType = "kafka_output"
    NodeTypeInfluxOutput    NodeType = "influx_output"
    NodeTypeESOutput        NodeType = "es_output"
    NodeTypeLLM             NodeType = "llm"
    NodeTypeEmbedding       NodeType = "embedding"
    NodeTypeBlockchain      NodeType = "blockchain"
)
```

### 3.2 Entity: Flow

```go
// domain/entity/flow.go
type Flow struct {
    ID          int64
    Name        string
    Status      value_object.FlowStatus
    UserID      string
    Nodes       []*Node
    Wires       []*Wire
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### 3.3 Entity: Node

```go
// domain/entity/node.go
type Node struct {
    ID          int64
    FlowID      int64
    Name        string
    Type        value_object.NodeType
    Config      map[string]interface{}
    Status      value_object.NodeStatus
    CreatedAt   time.Time
}
```

---

## 4. Application Layer

### 4.1 Use Cases

#### 4.1.1 StartFlowUseCase (ปรับปรุง)

```go
// application/start_flow.go
func (uc *StartFlowUseCase) Execute(ctx context.Context, flowID int64) error {
    // 1. โหลด Flow จาก DB
    flow, err := uc.repo.FindByID(ctx, flowID)
    if err != nil { return err }

    // 2. ตรวจสอบสถานะ
    if flow.Status != value_object.FlowStatusDeployed {
        return errors.ErrFlowNotDeployed
    }

    // 3. เรียก Engine เพื่อเริ่มทำงาน
    if err := uc.engine.StartFlow(ctx, flow); err != nil {
        return err
    }

    // 4. อัปเดตสถานะใน DB และ Cache
    flow.Status = value_object.FlowStatusRunning
    uc.repo.Update(ctx, flow)
    uc.cacheRepo.SetFlowStatus(ctx, flowID, flow.Status)

    // 5. ส่ง Event ไป Kafka
    event := map[string]interface{}{
        "type":      "flow.started",
        "flow_id":   flowID,
        "user_id":   flow.UserID,
        "timestamp": time.Now(),
    }
    _ = uc.kafkaProducer.Publish(ctx, "flow.events", event)
    _ = uc.kafkaProducer.Publish(ctx, "flow.blockchain.request", event)

    return nil
}
```

#### 4.1.2 GenerateAggregatesUseCase

```go
// application/generate_aggregates.go
type GenerateAggregatesUseCase struct {
    influxClient *timeseries.InfluxClient
    esClient     *search.ESClient
    db           *gorm.DB
}

func (uc *GenerateAggregatesUseCase) Execute(ctx context.Context) error {
    // 1. คำนวณ aggregate จาก InfluxDB และ ES
    aggregates := uc.computeAggregates(ctx)
    // 2. บันทึกลง PostgreSQL
    uc.saveAggregates(ctx, aggregates)
    // 3. ลบ raw data ที่เก่ากว่า 1 ปี
    uc.cleanupOldData(ctx)
    return nil
}
```

### 4.2 DTOs

```go
// application/dto.go

type LLMRequest struct {
    FlowID      int64                  `json:"flow_id"`
    NodeID      int64                  `json:"node_id"`
    Prompt      string                 `json:"prompt"`
    Model       string                 `json:"model"`
    Temperature float32                `json:"temperature"`
    Payload     map[string]interface{} `json:"payload"`
}

type EmbeddingRequest struct {
    FlowID      int64                  `json:"flow_id"`
    NodeID      int64                  `json:"node_id"`
    InputFields []string               `json:"input_fields"`
    Payload     map[string]interface{} `json:"payload"`
}

type BlockchainRequest struct {
    FlowID   int64                  `json:"flow_id"`
    NodeID   int64                  `json:"node_id"`
    DataHash string                 `json:"data_hash"`
    Payload  map[string]interface{} `json:"payload"`
}

type VectorIndexRequest struct {
    FlowID     int64                  `json:"flow_id"`
    NodeID     int64                  `json:"node_id"`
    Collection string                 `json:"collection"`
    Vector     []float64              `json:"vector"`
    Metadata   map[string]interface{} `json:"metadata"`
}

type SensorDataDTO struct {
    ID        string                 `json:"id"`
    DeviceID  string                 `json:"deviceId"`
    Value     float64                `json:"value"`
    Timestamp time.Time              `json:"timestamp"`
    FlowID    string                 `json:"flowId"`
    Metadata  map[string]interface{} `json:"metadata"`
}

type AlertDTO struct {
    ID        string    `json:"id"`
    Message   string    `json:"message"`
    Severity  string    `json:"severity"` // info, warning, error, success
    Source    string    `json:"source"`
    Timestamp time.Time `json:"timestamp"`
}
```

---

## 5. Infrastructure Layer

### 5.1 Persistence (PostgreSQL)

```sql
-- ตารางเพิ่มเติมสำหรับ Aggregates
CREATE TABLE sensor_aggregates (
    id          BIGSERIAL PRIMARY KEY,
    flow_id     BIGINT NOT NULL,
    device_id   TEXT NOT NULL,
    metric      TEXT NOT NULL,
    period      TEXT NOT NULL,   -- 'year', 'month', 'day'
    year        INTEGER NOT NULL,
    avg_value   DOUBLE PRECISION,
    max_value   DOUBLE PRECISION,
    min_value   DOUBLE PRECISION,
    count       BIGINT,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 5.2 Redis Integration

ใช้สำหรับ Cache, Rate Limit, และ Distributed State ของ Flow Instances

### 5.3 Kafka Integration

#### Producer

```go
// infrastructure/messaging/producer.go
func PublishLLMRequest(ctx context.Context, kc *KafkaClient, req *dto.LLMRequest) error {
    return kc.Publish(ctx, "flow.llm.request", req)
}

func PublishEmbeddingRequest(ctx context.Context, kc *KafkaClient, req *dto.EmbeddingRequest) error {
    return kc.Publish(ctx, "flow.embedding.request", req)
}

func PublishBlockchainRequest(ctx context.Context, kc *KafkaClient, req *dto.BlockchainRequest) error {
    return kc.Publish(ctx, "flow.blockchain.request", req)
}
```

#### Consumer Group

```go
// infrastructure/messaging/consumer.go
type FlowConsumer struct {
    Influx *timeseries.InfluxClient
    ES     *search.ESClient
    Socket *socketio.Server
}

func (c *FlowConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        var payload value_object.DataPayload
        json.Unmarshal(msg.Value, &payload)
        
        // Write to InfluxDB
        c.Influx.WritePoint(payload)
        
        // Index to Elasticsearch
        c.ES.IndexDocument(payload)
        
        // Broadcast via Socket.IO
        c.Socket.BroadcastToRoom("flow:"+payload.FlowID, "data", payload)
        
        sess.MarkMessage(msg, "")
    }
    return nil
}
```

### 5.4 InfluxDB Integration

```go
// infrastructure/timeseries/influx_client.go
type InfluxClient struct {
    client influxdb2.Client
    bucket string
    org    string
}

func (c *InfluxClient) WritePoint(payload value_object.DataPayload) error {
    point := influxdb2.NewPoint(
        "sensor",
        map[string]string{
            "device_id": payload.DeviceID,
            "flow_id":   payload.FlowID,
        },
        map[string]interface{}{
            "value": payload.Value,
        },
        payload.Timestamp,
    )
    return c.client.WritePoint(c.bucket, c.org, point)
}
```

### 5.5 Elasticsearch Integration

```go
// infrastructure/search/es_client.go
type ESClient struct {
    client *elasticsearch.Client
}

func (e *ESClient) CreateVectorIndex(index string, dims int) error {
    mapping := fmt.Sprintf(`{
        "mappings": {
            "properties": {
                "embedding": {
                    "type": "dense_vector",
                    "dims": %d,
                    "index": true,
                    "similarity": "cosine"
                },
                "text": { "type": "text" },
                "flow_id": { "type": "long" },
                "timestamp": { "type": "date" }
            }
        }
    }`, dims)
    req := bytes.NewReader([]byte(mapping))
    _, err := e.client.Indices.Create(index, e.client.Indices.Create.WithBody(req))
    return err
}

func (e *ESClient) VectorSearch(index string, queryVector []float64, size int) ([]map[string]interface{}, error) {
    // ใช้ ES 8.x knn query
    var result []map[string]interface{}
    // ... implementation
    return result, nil
}
```

### 5.6 Vector Database (Milvus/Qdrant)

#### Interface

```go
// infrastructure/vectordb/interface.go
type VectorStore interface {
    Insert(ctx context.Context, collection string, vector []float64, metadata map[string]interface{}) error
    Search(ctx context.Context, collection string, queryVector []float64, topK int) ([]SearchResult, error)
    Delete(ctx context.Context, collection string, ids ...string) error
    CreateCollection(ctx context.Context, name string, dims int) error
}

type SearchResult struct {
    ID       string                 `json:"id"`
    Score    float64                `json:"score"`
    Metadata map[string]interface{} `json:"metadata"`
}
```

#### Milvus Implementation

```go
// infrastructure/vectordb/milvus_client.go
import "github.com/milvus-io/milvus-sdk-go/v2/client"

type MilvusClient struct {
    client client.Client
}

func (m *MilvusClient) Insert(ctx context.Context, collection string, vector []float64, metadata map[string]interface{}) error {
    // แปลง metadata เป็น fields ที่ Milvus รองรับ
    // ...
}
```

### 5.7 LLM Integration

#### Client

```go
// infrastructure/llm/client.go
type LLMClient interface {
    Generate(ctx context.Context, prompt string, model string, temp float32) (string, error)
}

type OpenAIClient struct {
    client *openai.Client
}

func (c *OpenAIClient) Generate(ctx context.Context, prompt string, model string, temp float32) (string, error) {
    resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model:       model,
        Messages:    []openai.ChatCompletionMessage{{Role: "user", Content: prompt}},
        Temperature: temp,
    })
    if err != nil {
        return "", err
    }
    return resp.Choices[0].Message.Content, nil
}
```

#### Consumer

```go
// infrastructure/llm/consumer.go
type LLMConsumer struct {
    client   LLMClient
    producer *messaging.KafkaClient
}

func (c *LLMConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        var req dto.LLMRequest
        json.Unmarshal(msg.Value, &req)
        result, err := c.client.Generate(context.Background(), req.Prompt, req.Model, req.Temperature)
        if err != nil {
            slog.Error("LLM error", "err", err)
            continue
        }
        response := map[string]interface{}{
            "flow_id":  req.FlowID,
            "node_id":  req.NodeID,
            "result":   result,
            "original": req.Payload,
        }
        _ = c.producer.Publish(context.Background(), "flow.llm.response", response)
        sess.MarkMessage(msg, "")
    }
    return nil
}
```

### 5.8 AI Embedding Integration

#### Client

```go
// infrastructure/embedding/client.go
type EmbeddingClient interface {
    Embed(ctx context.Context, text string) ([]float64, error)
}

type RESTEmbeddingClient struct {
    url        string
    httpClient *http.Client
}

func (c *RESTEmbeddingClient) Embed(ctx context.Context, text string) ([]float64, error) {
    // POST ไปยัง embedding service
    reqBody := map[string]string{"text": text}
    jsonData, _ := json.Marshal(reqBody)
    req, _ := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    var result struct {
        Vector []float64 `json:"vector"`
    }
    json.NewDecoder(resp.Body).Decode(&result)
    return result.Vector, nil
}
```

### 5.9 Blockchain Integration

#### Client

```go
// infrastructure/blockchain/client.go
type BlockchainClient struct {
    client   *ethclient.Client
    contract *AuditContract
    auth     *bind.TransactOpts
}

func (b *BlockchainClient) StoreDataHash(ctx context.Context, dataHash [32]byte) (string, error) {
    tx, err := b.contract.StoreData(b.auth, dataHash)
    if err != nil {
        return "", err
    }
    return tx.Hash().Hex(), nil
}
```

### 5.10 Built-in Node Executors

#### LLM Executor

```go
// infrastructure/engine/builtin_nodes/llm_query.go
type LLMExecutor struct {
    kafkaProducer *messaging.KafkaClient
}

func (e *LLMExecutor) Execute(ctx context.Context, node *entity.Node, input <-chan value_object.DataPayload) (<-chan value_object.DataPayload, error) {
    output := make(chan value_object.DataPayload)
    go func() {
        defer close(output)
        for payload := range input {
            promptTemplate, _ := node.Config["prompt_template"].(string)
            model, _ := node.Config["model"].(string)
            temp, _ := node.Config["temperature"].(float64)
            prompt := renderTemplate(promptTemplate, payload.Payload)
            req := &dto.LLMRequest{
                FlowID:      payload.FlowID,
                NodeID:      node.ID,
                Prompt:      prompt,
                Model:       model,
                Temperature: float32(temp),
                Payload:     payload.Payload,
            }
            _ = e.kafkaProducer.Publish(ctx, "flow.llm.request", req)
            output <- payload
        }
    }()
    return output, nil
}
```

#### Embedding Executor

```go
// infrastructure/engine/builtin_nodes/embedding.go
type EmbeddingExecutor struct {
    kafkaProducer *messaging.KafkaClient
}

func (e *EmbeddingExecutor) Execute(ctx context.Context, node *entity.Node, input <-chan value_object.DataPayload) (<-chan value_object.DataPayload, error) {
    output := make(chan value_object.DataPayload)
    go func() {
        defer close(output)
        for payload := range input {
            inputFields, _ := node.Config["input_fields"].([]interface{})
            req := &dto.EmbeddingRequest{
                FlowID:      payload.FlowID,
                NodeID:      node.ID,
                InputFields: toStringSlice(inputFields),
                Payload:     payload.Payload,
            }
            _ = e.kafkaProducer.Publish(ctx, "flow.embedding.request", req)
            output <- payload
        }
    }()
    return output, nil
}
```

#### Blockchain Executor

```go
// infrastructure/engine/builtin_nodes/blockchain_record.go
type BlockchainExecutor struct {
    kafkaProducer *messaging.KafkaClient
}

func (e *BlockchainExecutor) Execute(ctx context.Context, node *entity.Node, input <-chan value_object.DataPayload) (<-chan value_object.DataPayload, error) {
    output := make(chan value_object.DataPayload)
    go func() {
        defer close(output)
        for payload := range input {
            fieldsToHash, _ := node.Config["fields_to_hash"].([]interface{})
            hash := computeHashFromFields(payload.Payload, fieldsToHash)
            req := &dto.BlockchainRequest{
                FlowID:   payload.FlowID,
                NodeID:   node.ID,
                DataHash: hash,
                Payload:  payload.Payload,
            }
            _ = e.kafkaProducer.Publish(ctx, "flow.blockchain.request", req)
            output <- payload
        }
    }()
    return output, nil
}
```

### 5.11 Aggregator & Cleanup

```go
// infrastructure/aggregator/aggregator.go
type Aggregator struct {
    influxClient *timeseries.InfluxClient
    esClient     *search.ESClient
    db           *gorm.DB
}

func (a *Aggregator) ComputeYearlyAggregates(ctx context.Context) ([]AggregateResult, error) {
    // 1. Query InfluxDB เพื่อหาค่าเฉลี่ย, สูงสุด, ต่ำสุด, count ของแต่ละ metric ในรอบ 1 ปี
    query := `from(bucket:"iotflow")
              |> range(start: -1y)
              |> filter(fn: (r) => r._measurement == "sensor")
              |> aggregateWindow(every: 1y, fn: mean)`
    // ... ดำเนินการ
    // 2. บันทึกผลลัพธ์ลง PostgreSQL
    return results, nil
}

func (a *Aggregator) CleanupOldData(ctx context.Context) error {
    // ลบ raw data ที่เก่ากว่า 1 ปีออกจาก InfluxDB
    // ลบ logs ที่เก่ากว่า 1 ปีออกจาก Elasticsearch
    _, err := a.esClient.DeleteByQuery(ctx, "logs-*", `{"query": {"range": {"timestamp": {"lt": "now-365d"}}}}`)
    return err
}
```

#### Scheduler

```go
// infrastructure/aggregator/scheduler.go
import "github.com/robfig/cron/v3"

func StartAggregationScheduler(influx *timeseries.InfluxClient, es *search.ESClient, db *gorm.DB) {
    agg := &Aggregator{influxClient: influx, esClient: es, db: db}
    c := cron.New()
    // ทำงานทุกวัน เวลา 02:00
    c.AddFunc("0 2 * * *", func() {
        ctx := context.Background()
        if err := agg.ComputeYearlyAggregates(ctx); err != nil {
            slog.Error("Aggregation failed", "err", err)
        }
        if err := agg.CleanupOldData(ctx); err != nil {
            slog.Error("Cleanup failed", "err", err)
        }
    })
    c.Start()
}
```

---

## 6. Interface Layer

### 6.1 HTTP Handlers

```go
// interfaces/http/flow_handler.go
type FlowHandler struct {
    createFlowUC  *application.CreateFlowUseCase
    startFlowUC   *application.StartFlowUseCase
    stopFlowUC    *application.StopFlowUseCase
    getFlowUC     *application.GetFlowUseCase
    listFlowsUC   *application.ListFlowsUseCase
}

func (h *FlowHandler) CreateFlow(c *gin.Context) {
    var req CreateFlowRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    flow, err := h.createFlowUC.Execute(c.Request.Context(), req.Name, req.UserID)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(201, flow)
}

func (h *FlowHandler) StartFlow(c *gin.Context) {
    flowID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    if err := h.startFlowUC.Execute(c.Request.Context(), flowID); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, gin.H{"status": "started"})
}
```

### 6.2 Routes

```go
// interfaces/http/routes.go
func SetupRoutes(r *gin.Engine, handler *FlowHandler) {
    api := r.Group("/api/v1/iot")
    {
        api.GET("/flows", handler.ListFlows)
        api.POST("/flows", handler.CreateFlow)
        api.GET("/flows/:id", handler.GetFlow)
        api.POST("/flows/:id/start", handler.StartFlow)
        api.POST("/flows/:id/stop", handler.StopFlow)
        api.POST("/flows/:id/deploy", handler.DeployFlow)
        api.DELETE("/flows/:id", handler.DeleteFlow)
        
        // Sensor Data API
        api.GET("/sensor-data/history", handler.GetSensorHistory)
        api.GET("/sensor-data/latest", handler.GetLatestSensorData)
        
        // Alert API
        api.GET("/alerts", handler.GetAlerts)
    }
}
```

### 6.3 Socket.IO Handler

```go
// interfaces/socketio/server.go
import "github.com/googollee/go-socket.io"

func NewSocketIOServer() (*socketio.Server, error) {
    server := socketio.NewServer(nil)
    server.OnConnect("/", func(s socketio.Conn) error {
        s.SetContext("")
        return nil
    })
    server.OnEvent("/", "subscribe", func(s socketio.Conn, flowID string) {
        s.Join("flow:" + flowID)
    })
    server.OnEvent("/", "unsubscribe", func(s socketio.Conn, flowID string) {
        s.Leave("flow:" + flowID)
    })
    server.OnDisconnect("/", func(s socketio.Conn, reason string) {
        // cleanup
    })
    return server, nil
}
```

### 6.4 Rate Limit Middleware

```go
// interfaces/middleware/rate_limit.go
func RateLimitMiddleware(redisClient *redis.Client) gin.HandlerFunc {
    limiter := redis.NewRateLimiter(redisClient)
    return func(c *gin.Context) {
        key := c.ClientIP()
        allowed, err := limiter.Allow(c.Request.Context(), key, 100, time.Minute)
        if err != nil || !allowed {
            c.JSON(429, gin.H{"error": "Too many requests"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

---

## 7. Background Workers / Consumers

### 7.1 Kafka Consumer Group (หลัก)

```go
// infrastructure/messaging/consumer.go
func (c *FlowConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        switch msg.Topic {
        case "flow.data":
            c.handleFlowData(msg)
        case "flow.alarm":
            c.handleFlowAlarm(msg)
        case "flow.llm.response":
            c.handleLLMResponse(msg)
        case "flow.blockchain.response":
            c.handleBlockchainResponse(msg)
        }
        sess.MarkMessage(msg, "")
    }
    return nil
}
```

### 7.2 InfluxDB Batch Writer

```go
// infrastructure/timeseries/writer.go
type BatchWriter struct {
    client   *InfluxClient
    buffer   []*influxdb2.Point
    maxSize  int
    flushInterval time.Duration
}

func (w *BatchWriter) Write(point *influxdb2.Point) {
    w.buffer = append(w.buffer, point)
    if len(w.buffer) >= w.maxSize {
        w.Flush()
    }
}

func (w *BatchWriter) Flush() {
    if len(w.buffer) == 0 {
        return
    }
    w.client.client.WriteBatch(w.client.bucket, w.client.org, w.buffer...)
    w.buffer = nil
}
```

### 7.3 Elasticsearch Bulk Indexer

```go
// infrastructure/search/indexer.go
type BulkIndexer struct {
    client *ESClient
    buffer []map[string]interface{}
    maxSize int
}

func (i *BulkIndexer) Index(index string, doc map[string]interface{}) {
    i.buffer = append(i.buffer, doc)
    if len(i.buffer) >= i.maxSize {
        i.Flush(index)
    }
}

func (i *BulkIndexer) Flush(index string) {
    if len(i.buffer) == 0 {
        return
    }
    // Bulk index to Elasticsearch
    i.client.BulkIndex(index, i.buffer)
    i.buffer = nil
}
```

### 7.4 WebSocket Broadcaster

```go
// interfaces/websocket/flow_monitor.go
type FlowMonitor struct {
    server *socketio.Server
}

func (m *FlowMonitor) BroadcastData(flowID string, data interface{}) {
    m.server.BroadcastToRoom("/", "flow:"+flowID, "flow:data", data)
}

func (m *FlowMonitor) BroadcastAlert(flowID string, alert AlertDTO) {
    m.server.BroadcastToRoom("/", "flow:"+flowID, "flow:alert", alert)
}
```

### 7.5 LLM Consumer

```go
// infrastructure/llm/consumer.go
type LLMConsumer struct {
    client   LLMClient
    producer *messaging.KafkaClient
}

func (c *LLMConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        var req dto.LLMRequest
        json.Unmarshal(msg.Value, &req)
        result, err := c.client.Generate(context.Background(), req.Prompt, req.Model, req.Temperature)
        if err != nil {
            slog.Error("LLM error", "err", err)
            continue
        }
        response := map[string]interface{}{
            "flow_id":  req.FlowID,
            "node_id":  req.NodeID,
            "result":   result,
            "original": req.Payload,
        }
        _ = c.producer.Publish(context.Background(), "flow.llm.response", response)
        sess.MarkMessage(msg, "")
    }
    return nil
}
```

### 7.6 Embedding Consumer

```go
// infrastructure/embedding/consumer.go
type EmbeddingConsumer struct {
    EmbedClient EmbeddingClient
    VectorStore  vectordb.VectorStore
    Producer    *messaging.KafkaClient
}

func (c *EmbeddingConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        var req dto.EmbeddingRequest
        json.Unmarshal(msg.Value, &req)
        text := buildTextFromFields(req.Payload, req.InputFields)
        vec, err := c.EmbedClient.Embed(context.Background(), text)
        if err != nil {
            slog.Error("Embedding error", "err", err)
            continue
        }
        err = c.VectorStore.Insert(context.Background(), "sensor_vectors", vec, req.Payload)
        if err != nil {
            slog.Error("Vector insert error", "err", err)
            continue
        }
        sess.MarkMessage(msg, "")
    }
    return nil
}
```

### 7.7 Blockchain Consumer

```go
// infrastructure/blockchain/consumer.go
type BlockchainConsumer struct {
    Client   *BlockchainClient
    Producer *messaging.KafkaClient
}

func (c *BlockchainConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        var req dto.BlockchainRequest
        json.Unmarshal(msg.Value, &req)
        hashBytes, _ := hex.DecodeString(req.DataHash)
        var hashArr [32]byte
        copy(hashArr[:], hashBytes)
        txHash, err := c.Client.StoreDataHash(context.Background(), hashArr)
        if err != nil {
            slog.Error("Blockchain error", "err", err)
            continue
        }
        response := map[string]interface{}{
            "flow_id":  req.FlowID,
            "node_id":  req.NodeID,
            "tx_hash":  txHash,
            "original": req.Payload,
        }
        _ = c.Producer.Publish(context.Background(), "flow.blockchain.response", response)
        sess.MarkMessage(msg, "")
    }
    return nil
}
```

---

## 8. Database Migrations

```sql
-- 001_create_flows.up.sql
CREATE TABLE flows (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    status TEXT NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

-- 002_create_nodes.up.sql
CREATE TABLE nodes (
    id BIGSERIAL PRIMARY KEY,
    flow_id BIGINT REFERENCES flows(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    config JSONB,
    status TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 003_create_wires.up.sql
CREATE TABLE wires (
    id BIGSERIAL PRIMARY KEY,
    flow_id BIGINT REFERENCES flows(id) ON DELETE CASCADE,
    source_node_id BIGINT REFERENCES nodes(id),
    target_node_id BIGINT REFERENCES nodes(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- 004_create_aggregates.up.sql
CREATE TABLE sensor_aggregates (
    id BIGSERIAL PRIMARY KEY,
    flow_id BIGINT NOT NULL,
    device_id TEXT NOT NULL,
    metric TEXT NOT NULL,
    period TEXT NOT NULL,
    year INTEGER NOT NULL,
    avg_value DOUBLE PRECISION,
    max_value DOUBLE PRECISION,
    min_value DOUBLE PRECISION,
    count BIGINT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## 9. Workflow Diagrams

### 9.1 สถาปัตยกรรมโดยรวม

```mermaid
graph TD
    subgraph External
        MQTT[MQTT Broker]
        HTTP[HTTP API]
    end
    subgraph Core
        Engine[Flow Engine]
        Nodes[Node Executors]
    end
    subgraph MessageBus[Message Bus - Kafka]
        Kafka[Kafka Topics]
    end
    subgraph Storage
        PG[(PostgreSQL)]
        Redis[(Redis)]
        Influx[(InfluxDB)]
        ES[(Elasticsearch + Vector)]
        VDB[(Vector DB - Milvus/Qdrant)]
    end
    subgraph Realtime
        Socket[Socket.IO Server]
        Dashboard[Dashboard UI]
    end
    subgraph AI
        LLM[LLM Service]
        Embed[Embedding Service]
    end
    subgraph Blockchain
        BC[Blockchain Network]
    end
    subgraph Workers
        Consumer[Kafka Consumer Group]
        LLMConsumer[LLM Consumer]
        EmbedConsumer[Embedding Consumer]
        BCWorker[Blockchain Consumer]
        VectorWorker[Vector Consumer]
        Aggregator[Aggregator Cron]
    end

    MQTT --> Nodes
    HTTP --> Engine
    Engine --> Nodes
    Nodes -->|Output| Kafka
    Kafka --> Consumer
    Consumer --> Influx
    Consumer --> ES
    Consumer --> Socket
    Socket --> Dashboard

    Nodes -->|LLM Request| Kafka --> LLMConsumer --> LLM --> Kafka --> Consumer
    Nodes -->|Embedding Request| Kafka --> EmbedConsumer --> Embed --> VectorWorker
    VectorWorker --> VDB
    VectorWorker --> ES
    Nodes -->|Blockchain Request| Kafka --> BCWorker --> BC --> Kafka --> Consumer

    Aggregator --> Influx
    Aggregator --> ES
    Aggregator --> PG

    Engine --> PG
    Engine --> Redis
    Engine --> Kafka
```

### 9.2 Sequence Diagram

```mermaid
sequenceDiagram
    participant Device
    participant MQTT
    participant InputNode
    participant Engine
    participant Kafka
    participant LLMConsumer
    participant EmbedConsumer
    participant BCWorker
    participant VectorConsumer
    participant InfluxDB
    participant ES
    participant VDB
    participant SocketIO
    participant Dashboard
    participant Aggregator

    Device->>MQTT: publish sensor data
    MQTT->>InputNode: on message
    InputNode->>Engine: forward payload
    Engine->>Engine: process via nodes
    Engine->>Kafka: publish flow.data
    Engine->>Kafka: publish flow.llm.request (if LLM node)
    Engine->>Kafka: publish flow.embedding.request (if Embedding node)
    Engine->>Kafka: publish flow.blockchain.request (if Blockchain node)

    Kafka->>LLMConsumer: consume request
    LLMConsumer->>LLMService: call API
    LLMService-->>LLMConsumer: response
    LLMConsumer->>Kafka: publish flow.llm.response

    Kafka->>EmbedConsumer: consume request
    EmbedConsumer->>EmbedService: generate vector
    EmbedService-->>EmbedConsumer: vector
    EmbedConsumer->>Kafka: publish flow.vector.index
    Kafka->>VectorConsumer: consume vector
    VectorConsumer->>VDB: insert vector
    VectorConsumer->>ES: index vector doc

    Kafka->>BCWorker: consume request
    BCWorker->>Blockchain: send transaction
    Blockchain-->>BCWorker: tx hash
    BCWorker->>Kafka: publish flow.blockchain.response

    Kafka->>Consumer: consume all responses + flow.data
    Consumer->>InfluxDB: write point
    Consumer->>ES: index log
    Consumer->>SocketIO: emit event
    SocketIO->>Dashboard: real-time update

    Aggregator->>InfluxDB: query aggregates (daily)
    Aggregator->>ES: query logs
    Aggregator->>PG: store aggregates
    Aggregator->>InfluxDB: delete old data
    Aggregator->>ES: delete old logs
```

---

## 10. การติดตั้งและใช้งาน

### 10.1 Dependencies

```bash
go get github.com/IBM/sarama
go get github.com/influxdata/influxdb-client-go/v2
go get github.com/elastic/go-elasticsearch/v8
go get github.com/googollee/go-socket.io
go get github.com/sashabaranov/go-openai
go get github.com/ethereum/go-ethereum
go get github.com/milvus-io/milvus-sdk-go/v2
go get github.com/robfig/cron/v3
go get github.com/redis/go-redis/v9
go get github.com/lib/pq
go get github.com/gin-gonic/gin
```

### 10.2 Environment Variables

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=iotflow

# Redis
REDIS_HOST=localhost:6379
REDIS_PASSWORD=

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_GROUP_ID=iotflow-consumer

# InfluxDB
INFLUX_URL=http://localhost:8086
INFLUX_TOKEN=my-token
INFLUX_ORG=my-org
INFLUX_BUCKET=iotflow

# Elasticsearch
ES_ADDRESSES=http://localhost:9200
ES_USERNAME=elastic
ES_PASSWORD=changeme

# Socket.IO
SOCKETIO_PORT=3000

# LLM (OpenAI)
OPENAI_API_KEY=sk-...
LLM_MODEL=gpt-4o-mini
LLM_PROVIDER=openai

# Embedding Service
EMBEDDING_SERVICE_URL=http://embedding-service:5000/embed

# Blockchain
ETH_RPC_URL=https://mainnet.infura.io/v3/...
ETH_PRIVATE_KEY=0x...
BLOCKCHAIN_CONTRACT_ADDRESS=0x...

# Vector DB (Milvus)
MILVUS_HOST=localhost
MILVUS_PORT=19530
VECTOR_COLLECTION=sensor_vectors
VECTOR_DIM=384

# Aggregation
AGGREGATION_CRON="0 2 * * *"
RETENTION_DAYS=365
```

### 10.3 การเริ่มต้นระบบ (Wire Up ใน `main.go`)

```go
package main

import (
    "context"
    "log"
    "os"
    "strings"

    "github.com/yourproject/internal/modules/iotflow/infrastructure/messaging"
    "github.com/yourproject/internal/modules/iotflow/infrastructure/timeseries"
    "github.com/yourproject/internal/modules/iotflow/infrastructure/search"
    "github.com/yourproject/internal/modules/iotflow/infrastructure/vectordb"
    "github.com/yourproject/internal/modules/iotflow/infrastructure/llm"
    "github.com/yourproject/internal/modules/iotflow/infrastructure/embedding"
    "github.com/yourproject/internal/modules/iotflow/infrastructure/blockchain"
    "github.com/yourproject/internal/modules/iotflow/infrastructure/aggregator"
    "github.com/yourproject/internal/modules/iotflow/infrastructure/cache/redis"
    "github.com/yourproject/internal/modules/iotflow/infrastructure/engine"
    "github.com/yourproject/internal/modules/iotflow/infrastructure/engine/builtin_nodes"
    "github.com/yourproject/internal/modules/iotflow/domain/value_object"
    "github.com/yourproject/internal/modules/iotflow/interfaces/http"
    "github.com/yourproject/internal/modules/iotflow/interfaces/socketio"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    ctx := context.Background()

    // 1. PostgreSQL
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // 2. Redis
    redisClient := redis.NewClient(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PASSWORD"))

    // 3. Kafka
    kafkaClient, err := messaging.NewKafkaClient(
        strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
        os.Getenv("KAFKA_GROUP_ID"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 4. InfluxDB
    influxClient := timeseries.NewInfluxClient(
        os.Getenv("INFLUX_URL"),
        os.Getenv("INFLUX_TOKEN"),
        os.Getenv("INFLUX_ORG"),
        os.Getenv("INFLUX_BUCKET"),
    )

    // 5. Elasticsearch
    esClient, err := search.NewESClient(
        strings.Split(os.Getenv("ES_ADDRESSES"), ","),
        os.Getenv("ES_USERNAME"),
        os.Getenv("ES_PASSWORD"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 6. Vector DB (Milvus)
    milvusClient, err := vectordb.NewMilvusClient(
        os.Getenv("MILVUS_HOST"),
        os.Getenv("MILVUS_PORT"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 7. Socket.IO
    socketServer, err := socketio.NewSocketIOServer()
    if err != nil {
        log.Fatal(err)
    }
    go socketServer.Serve()

    // 8. LLM Client
    llmClient := llm.NewOpenAIClient(os.Getenv("OPENAI_API_KEY"))

    // 9. Embedding Client
    embedClient := embedding.NewRESTEmbeddingClient(os.Getenv("EMBEDDING_SERVICE_URL"))

    // 10. Blockchain Client
    bcClient, err := blockchain.NewBlockchainClient(
        os.Getenv("ETH_RPC_URL"),
        os.Getenv("ETH_PRIVATE_KEY"),
        os.Getenv("BLOCKCHAIN_CONTRACT_ADDRESS"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 11. Engine
    instanceStore := redis.NewFlowInstanceStore(redisClient)
    engine := engine.NewRuntime(instanceStore, "instance-1", kafkaClient)

    // 12. Register Executors
    engine.RegisterExecutor(value_object.NodeTypeKafkaOutput, builtin_nodes.NewKafkaOutputExecutor(kafkaClient))
    engine.RegisterExecutor(value_object.NodeTypeInfluxOutput, builtin_nodes.NewInfluxOutputExecutor(influxClient))
    engine.RegisterExecutor(value_object.NodeTypeESOutput, builtin_nodes.NewESOutputExecutor(esClient))
    engine.RegisterExecutor(value_object.NodeTypeLLM, builtin_nodes.NewLLMExecutor(kafkaClient))
    engine.RegisterExecutor(value_object.NodeTypeEmbedding, builtin_nodes.NewEmbeddingExecutor(kafkaClient))
    engine.RegisterExecutor(value_object.NodeTypeBlockchain, builtin_nodes.NewBlockchainExecutor(kafkaClient))

    // 13. Start Consumers
    // 13.1 Main Consumer (Influx, ES, Socket)
    mainConsumer := &messaging.FlowConsumer{
        Influx: influxClient,
        ES:     esClient,
        Socket: socketServer,
    }
    go kafkaClient.Consume(ctx, []string{"flow.data", "flow.alarm", "flow.llm.response", "flow.blockchain.response"}, mainConsumer)

    // 13.2 LLM Consumer
    llmConsumer := &llm.LLMConsumer{
        Client:   llmClient,
        Producer: kafkaClient,
    }
    go kafkaClient.Consume(ctx, []string{"flow.llm.request"}, llmConsumer)

    // 13.3 Embedding Consumer
    embedConsumer := &embedding.EmbeddingConsumer{
        EmbedClient: embedClient,
        VectorStore: milvusClient,
        Producer:    kafkaClient,
    }
    go kafkaClient.Consume(ctx, []string{"flow.embedding.request"}, embedConsumer)

    // 13.4 Blockchain Consumer
    bcConsumer := &blockchain.BlockchainConsumer{
        Client:   bcClient,
        Producer: kafkaClient,
    }
    go kafkaClient.Consume(ctx, []string{"flow.blockchain.request"}, bcConsumer)

    // 14. Start Aggregation Scheduler
    go aggregator.StartAggregationScheduler(influxClient, esClient, db)

    // 15. Setup HTTP Server
    r := gin.Default()
    flowHandler := http.NewFlowHandler(/* dependencies */)
    http.SetupRoutes(r, flowHandler)
    r.Run(":8080")
}
```

---

## 11. Frontend Integration Guide

### 11.1 Angular Module Structure

สร้าง Angular feature module ภายใต้ `src/app/features/iot/`:

```
src/app/features/iot/
├── domain/
│   ├── entities/
│   │   ├── sensor-data.entity.ts
│   │   └── alert.entity.ts
│   ├── repositories/
│   │   └── iot.repository.ts
│   └── use-cases/
│       ├── get-sensor-history.use-case.ts
│       ├── get-latest-sensor.use-case.ts
│       ├── connect-realtime.use-case.ts
│       └── get-alerts.use-case.ts
├── data/
│   ├── dtos/
│   │   ├── sensor-data.dto.ts
│   │   └── alert.dto.ts
│   ├── datasources/
│   │   └── iot.api.datasource.ts
│   └── repositories/
│       └── iot.repository.impl.ts
└── presentation/
    └── pages/
        ├── dashboard/
        │   ├── dashboard.component.ts
        │   ├── dashboard.component.html
        │   └── dashboard.component.scss
        └── components/
            ├── sensor-chart/
            ├── alert-list/
            └── toast-container/
```

### 11.2 Domain Layer (Angular)

#### Entity: SensorData

```typescript
// domain/entities/sensor-data.entity.ts
export class SensorData {
  private constructor(
    public readonly id: string,
    public readonly deviceId: string,
    public readonly value: number,
    public readonly timestamp: Date,
    public readonly flowId: string,
    public readonly metadata?: Record<string, unknown>
  ) {}

  static create(props: SensorDataProps): SensorData {
    this.validate(props);
    return new SensorData(
      props.id,
      props.deviceId,
      props.value,
      props.timestamp,
      props.flowId,
      props.metadata
    );
  }

  private static validate(props: SensorDataProps): void {
    if (!props.deviceId) throw new Error('Device ID is required');
    if (props.value === undefined || props.value === null) {
      throw new Error('Value is required');
    }
  }

  isAbnormal(threshold: number): boolean {
    return this.value > threshold;
  }

  toJSON(): SensorDataDTO {
    return {
      id: this.id,
      deviceId: this.deviceId,
      value: this.value,
      timestamp: this.timestamp.toISOString(),
      flowId: this.flowId,
      metadata: this.metadata,
    };
  }
}

export interface SensorDataProps {
  id: string;
  deviceId: string;
  value: number;
  timestamp: Date;
  flowId: string;
  metadata?: Record<string, unknown>;
}
```

#### Value Object: Alert

```typescript
// domain/value-objects/alert.value-object.ts
export class Alert {
  private constructor(
    public readonly message: string,
    public readonly severity: AlertSeverity,
    public readonly timestamp: Date,
    public readonly source: string
  ) {}

  static create(message: string, severity: AlertSeverity, source: string): Alert {
    if (!message || message.trim().length === 0) {
      throw new Error('Alert message is required');
    }
    return new Alert(message, severity, new Date(), source);
  }
}

export type AlertSeverity = 'info' | 'warning' | 'error' | 'success';
```

#### Repository Interface

```typescript
// domain/repositories/iot.repository.ts
import { Observable } from 'rxjs';
import { SensorData } from '../entities/sensor-data.entity';
import { Alert } from '../value-objects/alert.value-object';

export interface IIoTRepository {
  getHistory(flowId: string, from: Date, to: Date): Observable<SensorData[]>;
  getLatest(flowId: string): Observable<SensorData | null>;
  getAlerts(flowId: string): Observable<Alert[]>;
  connect(flowId: string): Observable<SensorData>;
  disconnect(): void;
}

export const IOT_REPOSITORY = new InjectionToken<IIoTRepository>('IIoTRepository');
```

#### Use Cases

```typescript
// domain/use-cases/get-sensor-history.use-case.ts
@Injectable({ providedIn: 'root' })
export class GetSensorHistoryUseCase {
  constructor(@Inject(IOT_REPOSITORY) private repository: IIoTRepository) {}

  execute(flowId: string, from: Date, to: Date): Observable<SensorData[]> {
    return this.repository.getHistory(flowId, from, to);
  }
}

// domain/use-cases/get-latest-sensor.use-case.ts
@Injectable({ providedIn: 'root' })
export class GetLatestSensorUseCase {
  constructor(@Inject(IOT_REPOSITORY) private repository: IIoTRepository) {}

  execute(flowId: string): Observable<SensorData | null> {
    return this.repository.getLatest(flowId);
  }
}

// domain/use-cases/get-alerts.use-case.ts
@Injectable({ providedIn: 'root' })
export class GetAlertsUseCase {
  constructor(@Inject(IOT_REPOSITORY) private repository: IIoTRepository) {}

  execute(flowId: string): Observable<Alert[]> {
    return this.repository.getAlerts(flowId);
  }
}

// domain/use-cases/connect-realtime.use-case.ts
@Injectable({ providedIn: 'root' })
export class ConnectRealtimeUseCase {
  constructor(@Inject(IOT_REPOSITORY) private repository: IIoTRepository) {}

  execute(flowId: string): Observable<SensorData> {
    return this.repository.connect(flowId);
  }
}
```

### 11.3 Infrastructure Layer (Angular)

#### DTOs

```typescript
// data/dtos/sensor-data.dto.ts
export interface SensorDataDTO {
  id: string;
  deviceId: string;
  value: number;
  timestamp: string;
  flowId: string;
  metadata?: Record<string, unknown>;
}

// data/dtos/alert.dto.ts
export interface AlertDTO {
  id: string;
  message: string;
  severity: string;
  source: string;
  timestamp: string;
}
```

#### REST API Datasource

```typescript
// data/datasources/iot.api.datasource.ts
@Injectable({ providedIn: 'root' })
export class IoTApiDatasource {
  constructor(
    private http: HttpClient,
    private fallbackService: ApiFallbackService
  ) {}

  private endpoint(path: string): string {
    return `${this.fallbackService.getActiveBaseUrl()}${path}`;
  }

  getSensorHistory(flowId: string, from: Date, to: Date): Observable<SensorDataDTO[]> {
    const params = new HttpParams()
      .set('flowId', flowId)
      .set('from', from.toISOString())
      .set('to', to.toISOString());
    return this.http.get<SensorDataDTO[]>(this.endpoint('/api/v1/iot/sensor-data/history'), { params });
  }

  getLatestSensor(flowId: string): Observable<SensorDataDTO> {
    return this.http.get<SensorDataDTO>(this.endpoint('/api/v1/iot/sensor-data/latest'), {
      params: { flowId }
    });
  }

  getAlerts(flowId: string): Observable<AlertDTO[]> {
    return this.http.get<AlertDTO[]>(this.endpoint('/api/v1/iot/alerts'), {
      params: { flowId }
    });
  }
}
```

#### WebSocket Service

```typescript
// core/infrastructure/websocket/socket.service.ts
@Injectable({ providedIn: 'root' })
export class SocketService {
  private socket: Socket | null = null;
  private dataSubject = new Subject<SensorData>();
  private alertSubject = new Subject<Alert>();

  connect(flowId: string): void {
    if (this.socket?.connected) return;

    this.socket = io(environment.socketUrl, {
      query: { flowId },
      transports: ['websocket'],
    });

    this.socket.on('flow:data', (payload: SensorDataDTO) => {
      const sensorData = SensorData.create({
        id: payload.id,
        deviceId: payload.deviceId,
        value: payload.value,
        timestamp: new Date(payload.timestamp),
        flowId: payload.flowId,
        metadata: payload.metadata,
      });
      this.dataSubject.next(sensorData);
    });

    this.socket.on('flow:alert', (payload: AlertDTO) => {
      const alert = Alert.create(
        payload.message,
        payload.severity as AlertSeverity,
        payload.source
      );
      this.alertSubject.next(alert);
    });

    this.socket.on('connect', () => {
      console.log('Socket.IO connected');
    });
  }

  disconnect(): void {
    this.socket?.disconnect();
    this.socket = null;
  }

  get onData(): Observable<SensorData> {
    return this.dataSubject.asObservable();
  }

  get onAlert(): Observable<Alert> {
    return this.alertSubject.asObservable();
  }

  send(event: string, data: unknown): void {
    this.socket?.emit(event, data);
  }
}
```

#### Repository Implementation

```typescript
// data/repositories/iot.repository.impl.ts
@Injectable({ providedIn: 'root' })
export class IoTRepositoryImpl implements IIoTRepository {
  constructor(
    private api: IoTApiDatasource,
    private socket: SocketService
  ) {}

  getHistory(flowId: string, from: Date, to: Date): Observable<SensorData[]> {
    return this.api.getSensorHistory(flowId, from, to).pipe(
      map(dtos => dtos.map(dto => SensorData.create({
        id: dto.id,
        deviceId: dto.deviceId,
        value: dto.value,
        timestamp: new Date(dto.timestamp),
        flowId: dto.flowId,
        metadata: dto.metadata,
      })))
    );
  }

  getLatest(flowId: string): Observable<SensorData | null> {
    return this.api.getLatestSensor(flowId).pipe(
      map(dto => dto ? SensorData.create({
        id: dto.id,
        deviceId: dto.deviceId,
        value: dto.value,
        timestamp: new Date(dto.timestamp),
        flowId: dto.flowId,
        metadata: dto.metadata,
      }) : null)
    );
  }

  getAlerts(flowId: string): Observable<Alert[]> {
    return this.api.getAlerts(flowId).pipe(
      map(dtos => dtos.map(dto => Alert.create(
        dto.message,
        dto.severity as AlertSeverity,
        dto.source
      )))
    );
  }

  connect(flowId: string): Observable<SensorData> {
    this.socket.connect(flowId);
    return this.socket.onData;
  }

  disconnect(): void {
    this.socket.disconnect();
  }
}
```

### 11.4 Application Layer (Angular)

```typescript
// application/services/dashboard.service.ts
export interface DashboardState {
  history: SensorData[];
  latest: SensorData | null;
  alerts: Alert[];
  isConnected: boolean;
}

@Injectable({ providedIn: 'root' })
export class DashboardService {
  private getHistoryUC = inject(GetSensorHistoryUseCase);
  private getLatestUC = inject(GetLatestSensorUseCase);
  private getAlertsUC = inject(GetAlertsUseCase);
  private connectRealtimeUC = inject(ConnectRealtimeUseCase);
  private socketService = inject(SocketService);

  private stateSubject = new BehaviorSubject<DashboardState>({
    history: [],
    latest: null,
    alerts: [],
    isConnected: false,
  });

  state$ = this.stateSubject.asObservable();
  history$ = this.state$.pipe(map(state => state.history));
  latest$ = this.state$.pipe(map(state => state.latest));
  alerts$ = this.state$.pipe(map(state => state.alerts));

  private flowId: string | null = null;
  private destroy$ = new Subject<void>();

  initialize(flowId: string): void {
    this.flowId = flowId;
    const now = new Date();
    const oneHourAgo = new Date(now.getTime() - 60 * 60 * 1000);

    // 1. โหลด History ผ่าน REST API
    combineLatest([
      this.getHistoryUC.execute(flowId, oneHourAgo, now),
      this.getLatestUC.execute(flowId),
      this.getAlertsUC.execute(flowId),
    ]).pipe(takeUntil(this.destroy$)).subscribe(([history, latest, alerts]) => {
      this.stateSubject.next({
        ...this.stateSubject.value,
        history,
        latest,
        alerts,
      });
    });

    // 2. เชื่อมต่อ WebSocket สำหรับ Real-time
    this.connectRealtimeUC.execute(flowId).pipe(
      takeUntil(this.destroy$),
      scan((acc: SensorData[], data: SensorData) => {
        const newHistory = [data, ...acc].slice(0, 1000);
        return newHistory;
      }, [])
    ).subscribe((history) => {
      this.stateSubject.next({
        ...this.stateSubject.value,
        history,
        latest: history[0] || null,
        isConnected: true,
      });
    });

    // 3. รับ Alert
    this.socketService.onAlert.pipe(
      takeUntil(this.destroy$)
    ).subscribe((alert) => {
      const currentAlerts = this.stateSubject.value.alerts;
      this.stateSubject.next({
        ...this.stateSubject.value,
        alerts: [alert, ...currentAlerts].slice(0, 100),
      });
    });
  }

  disconnect(): void {
    this.socketService.disconnect();
    this.destroy$.next();
    this.destroy$.complete();
    this.stateSubject.next({
      ...this.stateSubject.value,
      isConnected: false,
    });
  }

  clearAlerts(): void {
    this.stateSubject.next({
      ...this.stateSubject.value,
      alerts: [],
    });
  }
}
```

### 11.5 Presentation Layer (Angular)

#### Dashboard Component

```typescript
// presentation/pages/dashboard/dashboard.component.ts
@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, TranslatePipe],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})
export class DashboardComponent implements OnInit, OnDestroy {
  private dashboardService = inject(DashboardService);
  private toastService = inject(ToastService);
  private destroy$ = new Subject<void>();

  history$ = this.dashboardService.history$;
  latest$ = this.dashboardService.latest$;
  alerts$ = this.dashboardService.alerts$;
  isConnected$ = this.dashboardService.state$.pipe(
    map(state => state.isConnected)
  );

  latestData: SensorData | null = null;
  alertCount = 0;

  ngOnInit(): void {
    this.dashboardService.initialize('flow-123');

    this.dashboardService.alerts$
      .pipe(takeUntil(this.destroy$))
      .subscribe((alerts) => {
        this.alertCount = alerts.length;
        if (alerts.length > 0) {
          const latestAlert = alerts[0];
          this.toastService.show({
            message: latestAlert.message,
            severity: latestAlert.severity,
            duration: 5000,
          });
        }
      });

    this.dashboardService.latest$
      .pipe(takeUntil(this.destroy$))
      .subscribe((data) => {
        this.latestData = data;
      });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
    this.dashboardService.disconnect();
  }

  clearAlerts(): void {
    this.dashboardService.clearAlerts();
  }

  trackById(index: number, item: SensorData): string {
    return item.id;
  }
}
```

#### Dashboard Template

```html
<!-- presentation/pages/dashboard/dashboard.component.html -->
<div class="dashboard-container">
  <header class="dashboard-header">
    <h1>{{ 'iot.title' | translate }}</h1>
    <div class="connection-status">
      <span [class.connected]="(isConnected$ | async)">
        {{ (isConnected$ | async) ? '🟢 Online' : '🔴 Offline' }}
      </span>
    </div>
  </header>

  <div class="latest-data" *ngIf="latestData">
    <h2>{{ 'iot.latestReading' | translate }}</h2>
    <div class="sensor-card">
      <div class="sensor-value">{{ latestData.value }}</div>
      <div class="sensor-meta">
        <span>{{ 'iot.device' | translate }}: {{ latestData.deviceId }}</span>
        <span>{{ 'iot.time' | translate }}: {{ latestData.timestamp | date:'HH:mm:ss' }}</span>
      </div>
    </div>
  </div>

  <div class="alert-banner" *ngIf="alertCount > 0" (click)="clearAlerts()">
    ⚠️ {{ alertCount }} {{ 'iot.alerts' | translate }} - {{ 'iot.clickToClear' | translate }}
  </div>

  <div class="history-section">
    <h2>{{ 'iot.history' | translate }}</h2>
    <div class="history-list">
      <div *ngFor="let item of (history$ | async) | slice:0:50; trackBy: trackById"
           class="history-item">
        <span class="time">{{ item.timestamp | date:'HH:mm:ss' }}</span>
        <span class="value">{{ item.value }}</span>
        <span class="device">{{ item.deviceId }}</span>
      </div>
    </div>
  </div>
</div>
```

#### Toast Service

```typescript
// presentation/services/toast.service.ts
export interface ToastMessage {
  id?: string;
  message: string;
  severity: 'info' | 'warning' | 'error' | 'success';
  duration?: number;
  title?: string;
}

@Injectable({ providedIn: 'root' })
export class ToastService {
  private toastSubject = new Subject<ToastMessage>();
  toasts$ = this.toastSubject.asObservable();

  show(toast: ToastMessage): void {
    const id = crypto.randomUUID();
    this.toastSubject.next({ ...toast, id });
  }

  success(message: string, title?: string): void {
    this.show({ message, severity: 'success', title });
  }

  error(message: string, title?: string): void {
    this.show({ message, severity: 'error', title });
  }

  warning(message: string, title?: string): void {
    this.show({ message, severity: 'warning', title });
  }

  info(message: string, title?: string): void {
    this.show({ message, severity: 'info', title });
  }
}
```

### 11.6 Environment Configuration

```typescript
// environments/environment.ts
export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080',
  socketUrl: 'http://localhost:3000',
};
```

### 11.7 Provider Configuration

```typescript
// app.module.ts or feature module
@NgModule({
  providers: [
    {
      provide: IOT_REPOSITORY,
      useClass: IoTRepositoryImpl,
    },
  ],
})
export class IotModule {}
```

### 11.8 API Endpoints Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/iot/flows` | List all flows |
| POST | `/api/v1/iot/flows` | Create new flow |
| GET | `/api/v1/iot/flows/:id` | Get flow details |
| POST | `/api/v1/iot/flows/:id/start` | Start flow |
| POST | `/api/v1/iot/flows/:id/stop` | Stop flow |
| GET | `/api/v1/iot/sensor-data/history` | Get sensor history |
| GET | `/api/v1/iot/sensor-data/latest` | Get latest sensor data |
| GET | `/api/v1/iot/alerts` | Get alerts |

### 11.9 WebSocket Events

| Event | Direction | Payload | Description |
|-------|-----------|---------|-------------|
| `flow:data` | Server → Client | `SensorDataDTO` | Real-time sensor data |
| `flow:alert` | Server → Client | `AlertDTO` | Real-time alert notification |
| `subscribe` | Client → Server | `{ flowId: string }` | Subscribe to flow updates |
| `unsubscribe` | Client → Server | `{ flowId: string }` | Unsubscribe from flow updates |

---

## 12. สรุปและแนวทางพัฒนาเพิ่มเติม

### 12.1 สรุป

เราได้สร้างระบบ IoT Flow Management Module ที่สมบูรณ์พร้อมความสามารถ:

- **Event‑Driven Architecture** ด้วย Kafka ช่วยให้แยกส่วนและขยายขนาดได้
- **Time‑series Storage** ด้วย InfluxDB สำหรับการวิเคราะห์และแสดงผล
- **Search & Logging** ด้วย Elasticsearch พร้อมความสามารถ Vector Search
- **Vector Database** เฉพาะทาง (Milvus) สำหรับ Semantic Search ประสิทธิภาพสูง
- **Real‑time Dashboard** ผ่าน Socket.IO
- **AI‑Powered** ด้วย LLM และ Embedding สำหรับการสรุป, ตอบคำถาม, และค้นหาเชิงความหมาย
- **Immutable Audit** ด้วย Blockchain สำหรับบันทึกหลักฐานที่ไม่สามารถแก้ไขได้
- **Logging Management** ที่เก็บสถิติรายปีและลบข้อมูลเก่าโดยอัตโนมัติ

### 12.2 แนวทางพัฒนาเพิ่มเติม

1. **Dashboard UI** – สร้างหน้าแสดงผลแบบ Real‑time โดยใช้ Socket.IO เพื่อแสดงกราฟ, ผลลัพธ์ LLM, และการค้นหาเวกเตอร์
2. **Kibana Integration** – ใช้ Kibana สำหรับ visualize logs และทำ vector search
3. **Batch Blockchain Transactions** – รวมหลาย ๆ requests เข้าด้วยกันเพื่อลดค่า Gas
4. **Caching Embeddings** – ใช้ Redis แคชเวกเตอร์ที่สร้างแล้วเพื่อลดภาระการคำนวณ
5. **Circuit Breaker & Retry** – สำหรับการเรียก LLM/Embedding API เพื่อเพิ่มความเสถียร
6. **Auto‑scaling** – ปรับจำนวน Kafka Consumers และ Workers ตามโหลด
7. **Multi‑tenant** – รองรับหลายองค์กร/ผู้ใช้ ด้วยการแยกข้อมูลผ่าน `flow_id` และ `user_id`
8. **Integration Tests** – เขียนทดสอบการทำงานร่วมกันของทุก component
9. **Monitoring** – ใช้ Prometheus + Grafana เพื่อตรวจสอบสถานะของระบบและ Kafka Lag
10. **Data Retention Policy** – ปรับแต่ง Retention ของ InfluxDB และ ILM ของ Elasticsearch
11. **Mobile App** – พัฒนาแอปพลิเคชันมือถือสำหรับดูข้อมูลแบบเรียลไทม์
12. **Webhook Support** – รองรับการส่งข้อมูลไปยังเว็บhook ของ
# ระบบตั้งค่าใช้ Database PostgreSQL สำหรับ IoT Flow Management Module

## สารบัญ

1. [ภาพรวมการตั้งค่า PostgreSQL](#1-ภาพรวมการตั้งค่า-postgresql)
2. [การติดตั้ง PostgreSQL](#2-การติดตั้ง-postgresql)
3. [การสร้าง Database และ User](#3-การสร้าง-database-และ-user)
4. [โครงสร้างตาราง (Schema)](#4-โครงสร้างตาราง-schema)
5. [การเชื่อมต่อ PostgreSQL ใน Go](#5-การเชื่อมต่อ-postgresql-ใน-go)
6. [การทำ Database Migration](#6-การทำ-database-migration)
7. [การตั้งค่า Environment Variables](#7-การตั้งค่า-environment-variables)
8. [PostgreSQL Configuration Optimization](#8-postgresql-configuration-optimization)
9. [การ Backup และ Restore](#9-การ-backup-และ-restore)
10. [Monitoring และ Troubleshooting](#10-monitoring-และ-troubleshooting)
11. [Docker Setup](#11-docker-setup)

---

## 1. ภาพรวมการตั้งค่า PostgreSQL

PostgreSQL เป็นฐานข้อมูลหลักของระบบ IoT Flow Management Module ใช้สำหรับเก็บ:

| ตาราง | วัตถุประสงค์ |
|-------|-------------|
| `flows` | เก็บข้อมูล Flow (ชื่อ, สถานะ, ผู้สร้าง) |
| `nodes` | เก็บข้อมูล Node ในแต่ละ Flow |
| `wires` | เก็บความสัมพันธ์ระหว่าง Nodes (การเชื่อมต่อ) |
| `sensor_aggregates` | เก็บข้อมูลสถิติสรุปประจำปี |
| `users` | เก็บข้อมูลผู้ใช้ (ถ้าใช้ระบบ Auth) |
| `audit_logs` | เก็บประวัติการทำงานของระบบ |

---

## 2. การติดตั้ง PostgreSQL

### 2.1 ติดตั้งบน Ubuntu/Debian

```bash
# Update package list
sudo apt update

# Install PostgreSQL
sudo apt install postgresql postgresql-contrib -y

# Check status
sudo systemctl status postgresql

# Start PostgreSQL
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

### 2.2 ติดตั้งบน macOS (Homebrew)

```bash
# Install PostgreSQL
brew install postgresql@15

# Start PostgreSQL
brew services start postgresql@15

# Add to PATH
echo 'export PATH="/opt/homebrew/opt/postgresql@15/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### 2.3 ติดตั้งบน Windows

1. ดาวน์โหลด installer จาก [https://www.postgresql.org/download/windows/](https://www.postgresql.org/download/windows/)
2. เรียกใช้ installer และเลือก components:
   - PostgreSQL Server
   - pgAdmin (optional)
   - Command Line Tools
3. กำหนดรหัสผ่านสำหรับ `postgres` user
4. กำหนด port (default: 5432)

---

## 3. การสร้าง Database และ User

### 3.1 เข้าสู่ระบบ PostgreSQL

```bash
# Login as postgres user
sudo -u postgres psql

# หรือใช้ password
psql -U postgres -h localhost
```

### 3.2 สร้าง Database User

```sql
-- สร้าง user สำหรับแอปพลิเคชัน
CREATE USER iotflow_user WITH PASSWORD 'your_secure_password';

-- ให้สิทธิ์ในการสร้าง database
ALTER USER iotflow_user CREATEDB;

-- ให้สิทธิ์ทั้งหมดบน database
GRANT ALL PRIVILEGES ON DATABASE iotflow TO iotflow_user;
```

### 3.3 สร้าง Database

```sql
-- สร้าง database
CREATE DATABASE iotflow 
    OWNER iotflow_user 
    ENCODING 'UTF8' 
    LC_COLLATE 'en_US.UTF-8' 
    LC_CTYPE 'en_US.UTF-8' 
    TEMPLATE template0;

-- เชื่อมต่อกับ database ที่สร้าง
\c iotflow;

-- ให้สิทธิ์ทั้งหมด
GRANT ALL PRIVILEGES ON DATABASE iotflow TO iotflow_user;
```

### 3.4 ทดสอบการเชื่อมต่อ

```bash
# ทดสอบเชื่อมต่อด้วย user ที่สร้าง
psql -U iotflow_user -d iotflow -h localhost

# ดูตารางทั้งหมด
\dt

# ออกจากระบบ
\q
```

---

## 4. โครงสร้างตาราง (Schema)

### 4.1 Migration Files

สร้างโฟลเดอร์ `migrations` ใน project root:

```
migrations/
├── 001_create_flows_table.up.sql
├── 001_create_flows_table.down.sql
├── 002_create_nodes_table.up.sql
├── 002_create_nodes_table.down.sql
├── 003_create_wires_table.up.sql
├── 003_create_wires_table.down.sql
├── 004_create_sensor_aggregates_table.up.sql
├── 004_create_sensor_aggregates_table.down.sql
├── 005_create_users_table.up.sql
├── 005_create_users_table.down.sql
├── 006_create_audit_logs_table.up.sql
├── 006_create_audit_logs_table.down.sql
└── 007_add_indexes.up.sql
```

### 4.2 001_create_flows_table.up.sql

```sql
-- 001_create_flows_table.up.sql
CREATE TABLE IF NOT EXISTS flows (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    user_id UUID NOT NULL,
    config JSONB DEFAULT '{}',
    version INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Indexes
CREATE INDEX idx_flows_user_id ON flows(user_id);
CREATE INDEX idx_flows_status ON flows(status);
CREATE INDEX idx_flows_created_at ON flows(created_at);
CREATE INDEX idx_flows_deleted_at ON flows(deleted_at) WHERE deleted_at IS NOT NULL;
```

### 4.3 001_create_flows_table.down.sql

```sql
DROP TABLE IF EXISTS flows;
```

### 4.4 002_create_nodes_table.up.sql

```sql
-- 002_create_nodes_table.up.sql
CREATE TABLE IF NOT EXISTS nodes (
    id BIGSERIAL PRIMARY KEY,
    flow_id BIGINT NOT NULL REFERENCES flows(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    config JSONB DEFAULT '{}',
    position_x INTEGER DEFAULT 0,
    position_y INTEGER DEFAULT 0,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_nodes_flow_id ON nodes(flow_id);
CREATE INDEX idx_nodes_type ON nodes(type);
CREATE INDEX idx_nodes_status ON nodes(status);
```

### 4.5 003_create_wires_table.up.sql

```sql
-- 003_create_wires_table.up.sql
CREATE TABLE IF NOT EXISTS wires (
    id BIGSERIAL PRIMARY KEY,
    flow_id BIGINT NOT NULL REFERENCES flows(id) ON DELETE CASCADE,
    source_node_id BIGINT NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    target_node_id BIGINT NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    config JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(flow_id, source_node_id, target_node_id)
);

-- Indexes
CREATE INDEX idx_wires_flow_id ON wires(flow_id);
CREATE INDEX idx_wires_source_node_id ON wires(source_node_id);
CREATE INDEX idx_wires_target_node_id ON wires(target_node_id);
```

### 4.6 004_create_sensor_aggregates_table.up.sql

```sql
-- 004_create_sensor_aggregates_table.up.sql
CREATE TABLE IF NOT EXISTS sensor_aggregates (
    id BIGSERIAL PRIMARY KEY,
    flow_id BIGINT NOT NULL,
    device_id VARCHAR(255) NOT NULL,
    metric VARCHAR(100) NOT NULL,
    period VARCHAR(20) NOT NULL, -- 'hour', 'day', 'month', 'year'
    year INTEGER NOT NULL,
    month INTEGER,
    day INTEGER,
    hour INTEGER,
    avg_value DOUBLE PRECISION,
    max_value DOUBLE PRECISION,
    min_value DOUBLE PRECISION,
    sum_value DOUBLE PRECISION,
    count BIGINT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_sensor_aggregates_flow_id ON sensor_aggregates(flow_id);
CREATE INDEX idx_sensor_aggregates_device_id ON sensor_aggregates(device_id);
CREATE INDEX idx_sensor_aggregates_period_year ON sensor_aggregates(period, year);
CREATE INDEX idx_sensor_aggregates_metric ON sensor_aggregates(metric);
CREATE UNIQUE INDEX idx_sensor_aggregates_unique ON sensor_aggregates(flow_id, device_id, metric, period, year, month, day, hour);
```

### 4.7 005_create_users_table.up.sql

```sql
-- 005_create_users_table.up.sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    role VARCHAR(50) DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role);
```

### 4.8 006_create_audit_logs_table.up.sql

```sql
-- 006_create_audit_logs_table.up.sql
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id VARCHAR(255),
    details JSONB DEFAULT '{}',
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_resource_type ON audit_logs(resource_type);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
```

### 4.9 007_add_indexes.up.sql

```sql
-- 007_add_indexes.up.sql

-- Partial index สำหรับ flows ที่ยังไม่ถูกลบ
CREATE INDEX idx_flows_active ON flows(id) WHERE deleted_at IS NULL;

-- Composite index สำหรับค้นหา flows
CREATE INDEX idx_flows_user_status ON flows(user_id, status) WHERE deleted_at IS NULL;

-- Index สำหรับ JSONB queries
CREATE INDEX idx_flows_config ON flows USING GIN (config);

-- Index สำหรับ nodes config
CREATE INDEX idx_nodes_config ON nodes USING GIN (config);

-- Index สำหรับ search
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_flows_name_trgm ON flows USING GIN (name gin_trgm_ops);
CREATE INDEX idx_nodes_name_trgm ON nodes USING GIN (name gin_trgm_ops);
```

---

## 5. การเชื่อมต่อ PostgreSQL ใน Go

### 5.1 ติดตั้ง Dependencies

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/golang-migrate/migrate/v4
go get -u github.com/golang-migrate/migrate/v4/database/postgres
go get -u github.com/golang-migrate/migrate/v4/source/file
```

### 5.2 Database Connection

```go
// infrastructure/persistence/postgres/db.go
package postgres

import (
    "fmt"
    "log"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type DBConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
    Timezone string
}

func NewDBConnection(config DBConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
        config.Host,
        config.User,
        config.Password,
        config.DBName,
        config.Port,
        config.SSLMode,
        config.Timezone,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
        SkipDefaultTransaction: true,
        PrepareStmt:            true,
    })

    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    // Get underlying SQL DB
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }

    // Connection Pool Settings
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    sqlDB.SetConnMaxIdleTime(10 * time.Minute)

    return db, nil
}

// TestConnection ทดสอบการเชื่อมต่อ
func TestConnection(db *gorm.DB) error {
    sqlDB, err := db.DB()
    if err != nil {
        return err
    }
    return sqlDB.Ping()
}
```

### 5.3 Repository Implementation

```go
// infrastructure/persistence/postgres/flow_repo_impl.go
package postgres

import (
    "context"
    "time"

    "gorm.io/gorm"
    "github.com/yourproject/internal/modules/iotflow/domain/entity"
    "github.com/yourproject/internal/modules/iotflow/domain/repository"
)

type FlowRepositoryImpl struct {
    db *gorm.DB
}

func NewFlowRepository(db *gorm.DB) repository.FlowRepository {
    return &FlowRepositoryImpl{db: db}
}

func (r *FlowRepositoryImpl) Create(ctx context.Context, flow *entity.Flow) error {
    model := toFlowModel(flow)
    return r.db.WithContext(ctx).Create(model).Error
}

func (r *FlowRepositoryImpl) FindByID(ctx context.Context, id int64) (*entity.Flow, error) {
    var model FlowModel
    err := r.db.WithContext(ctx).
        Preload("Nodes").
        Preload("Wires").
        Where("id = ? AND deleted_at IS NULL", id).
        First(&model).Error
    if err != nil {
        return nil, err
    }
    return toFlowEntity(&model), nil
}

func (r *FlowRepositoryImpl) Update(ctx context.Context, flow *entity.Flow) error {
    model := toFlowModel(flow)
    model.UpdatedAt = time.Now()
    return r.db.WithContext(ctx).Save(model).Error
}

func (r *FlowRepositoryImpl) Delete(ctx context.Context, id int64) error {
    return r.db.WithContext(ctx).
        Model(&FlowModel{}).
        Where("id = ?", id).
        Update("deleted_at", time.Now()).
        Error
}

func (r *FlowRepositoryImpl) List(ctx context.Context, userID string, limit, offset int) ([]*entity.Flow, int64, error) {
    var models []FlowModel
    var total int64

    query := r.db.WithContext(ctx).
        Model(&FlowModel{}).
        Where("user_id = ? AND deleted_at IS NULL", userID)

    // Count total
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Get paginated results
    err := query.
        Preload("Nodes").
        Preload("Wires").
        Limit(limit).
        Offset(offset).
        Order("created_at DESC").
        Find(&models).Error
    if err != nil {
        return nil, 0, err
    }

    flows := make([]*entity.Flow, len(models))
    for i, model := range models {
        flows[i] = toFlowEntity(&model)
    }

    return flows, total, nil
}
```

### 5.4 Models

```go
// infrastructure/persistence/postgres/models.go
package postgres

import (
    "database/sql/driver"
    "encoding/json"
    "time"

    "gorm.io/gorm"
)

// FlowModel GORM model สำหรับตาราง flows
type FlowModel struct {
    ID          int64          `gorm:"primaryKey"`
    Name        string         `gorm:"column:name;type:varchar(255);not null"`
    Description string         `gorm:"column:description;type:text"`
    Status      string         `gorm:"column:status;type:varchar(50);not null;default:'draft'"`
    UserID      string         `gorm:"column:user_id;type:uuid;not null"`
    Config      JSONMap        `gorm:"column:config;type:jsonb;default:'{}'"`
    Version     int            `gorm:"column:version;default:1"`
    CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
    DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
    Nodes       []NodeModel    `gorm:"foreignKey:FlowID;constraint:OnDelete:CASCADE"`
    Wires       []WireModel    `gorm:"foreignKey:FlowID;constraint:OnDelete:CASCADE"`
}

func (FlowModel) TableName() string {
    return "flows"
}

// NodeModel GORM model สำหรับตาราง nodes
type NodeModel struct {
    ID        int64     `gorm:"primaryKey"`
    FlowID    int64     `gorm:"column:flow_id;not null;index"`
    Name      string    `gorm:"column:name;type:varchar(255);not null"`
    Type      string    `gorm:"column:type;type:varchar(50);not null"`
    Config    JSONMap   `gorm:"column:config;type:jsonb;default:'{}'"`
    PositionX int       `gorm:"column:position_x;default:0"`
    PositionY int       `gorm:"column:position_y;default:0"`
    Status    string    `gorm:"column:status;type:varchar(50);default:'active'"`
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (NodeModel) TableName() string {
    return "nodes"
}

// WireModel GORM model สำหรับตาราง wires
type WireModel struct {
    ID           int64     `gorm:"primaryKey"`
    FlowID       int64     `gorm:"column:flow_id;not null;index"`
    SourceNodeID int64     `gorm:"column:source_node_id;not null"`
    TargetNodeID int64     `gorm:"column:target_node_id;not null"`
    Config       JSONMap   `gorm:"column:config;type:jsonb;default:'{}'"`
    CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (WireModel) TableName() string {
    return "wires"
}

// SensorAggregateModel GORM model สำหรับตาราง sensor_aggregates
type SensorAggregateModel struct {
    ID        int64     `gorm:"primaryKey"`
    FlowID    int64     `gorm:"column:flow_id;not null;index"`
    DeviceID  string    `gorm:"column:device_id;type:varchar(255);not null"`
    Metric    string    `gorm:"column:metric;type:varchar(100);not null"`
    Period    string    `gorm:"column:period;type:varchar(20);not null"`
    Year      int       `gorm:"column:year;not null"`
    Month     int       `gorm:"column:month"`
    Day       int       `gorm:"column:day"`
    Hour      int       `gorm:"column:hour"`
    AvgValue  *float64  `gorm:"column:avg_value;type:double precision"`
    MaxValue  *float64  `gorm:"column:max_value;type:double precision"`
    MinValue  *float64  `gorm:"column:min_value;type:double precision"`
    SumValue  *float64  `gorm:"column:sum_value;type:double precision"`
    Count     *int64    `gorm:"column:count;type:bigint"`
    Metadata  JSONMap   `gorm:"column:metadata;type:jsonb;default:'{}'"`
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (SensorAggregateModel) TableName() string {
    return "sensor_aggregates"
}

// JSONMap สำหรับ JSONB fields
type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
    if j == nil {
        return []byte("{}"), nil
    }
    return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
    if value == nil {
        *j = make(JSONMap)
        return nil
    }

    var bytes []byte
    switch v := value.(type) {
    case []byte:
        bytes = v
    case string:
        bytes = []byte(v)
    default:
        return nil
    }

    if len(bytes) == 0 {
        *j = make(JSONMap)
        return nil
    }

    return json.Unmarshal(bytes, j)
}
```

### 5.5 Entity to Model Mappers

```go
// infrastructure/persistence/postgres/mapper.go
package postgres

import (
    "github.com/yourproject/internal/modules/iotflow/domain/entity"
    "github.com/yourproject/internal/modules/iotflow/domain/value_object"
)

func toFlowModel(flow *entity.Flow) *FlowModel {
    if flow == nil {
        return nil
    }

    model := &FlowModel{
        ID:          flow.ID,
        Name:        flow.Name,
        Status:      string(flow.Status),
        UserID:      flow.UserID,
        Config:      JSONMap(flow.Config),
        Version:     flow.Version,
        CreatedAt:   flow.CreatedAt,
        UpdatedAt:   flow.UpdatedAt,
    }

    // Convert Nodes
    if len(flow.Nodes) > 0 {
        model.Nodes = make([]NodeModel, len(flow.Nodes))
        for i, node := range flow.Nodes {
            model.Nodes[i] = *toNodeModel(node)
        }
    }

    // Convert Wires
    if len(flow.Wires) > 0 {
        model.Wires = make([]WireModel, len(flow.Wires))
        for i, wire := range flow.Wires {
            model.Wires[i] = *toWireModel(wire)
        }
    }

    return model
}

func toFlowEntity(model *FlowModel) *entity.Flow {
    if model == nil {
        return nil
    }

    flow := &entity.Flow{
        ID:        model.ID,
        Name:      model.Name,
        Status:    value_object.FlowStatus(model.Status),
        UserID:    model.UserID,
        Config:    map[string]interface{}(model.Config),
        Version:   model.Version,
        CreatedAt: model.CreatedAt,
        UpdatedAt: model.UpdatedAt,
    }

    // Convert Nodes
    if len(model.Nodes) > 0 {
        flow.Nodes = make([]*entity.Node, len(model.Nodes))
        for i, nodeModel := range model.Nodes {
            flow.Nodes[i] = toNodeEntity(&nodeModel)
        }
    }

    // Convert Wires
    if len(model.Wires) > 0 {
        flow.Wires = make([]*entity.Wire, len(model.Wires))
        for i, wireModel := range model.Wires {
            flow.Wires[i] = toWireEntity(&wireModel)
        }
    }

    return flow
}

func toNodeModel(node *entity.Node) *NodeModel {
    if node == nil {
        return nil
    }
    return &NodeModel{
        ID:        node.ID,
        FlowID:    node.FlowID,
        Name:      node.Name,
        Type:      string(node.Type),
        Config:    JSONMap(node.Config),
        Status:    string(node.Status),
        CreatedAt: node.CreatedAt,
        UpdatedAt: node.UpdatedAt,
    }
}

func toNodeEntity(model *NodeModel) *entity.Node {
    if model == nil {
        return nil
    }
    return &entity.Node{
        ID:        model.ID,
        FlowID:    model.FlowID,
        Name:      model.Name,
        Type:      value_object.NodeType(model.Type),
        Config:    map[string]interface{}(model.Config),
        Status:    value_object.NodeStatus(model.Status),
        CreatedAt: model.CreatedAt,
        UpdatedAt: model.UpdatedAt,
    }
}

func toWireModel(wire *entity.Wire) *WireModel {
    if wire == nil {
        return nil
    }
    return &WireModel{
        ID:           wire.ID,
        FlowID:       wire.FlowID,
        SourceNodeID: wire.SourceNodeID,
        TargetNodeID: wire.TargetNodeID,
        Config:       JSONMap(wire.Config),
        CreatedAt:    wire.CreatedAt,
        UpdatedAt:    wire.UpdatedAt,
    }
}

func toWireEntity(model *WireModel) *entity.Wire {
    if model == nil {
        return nil
    }
    return &entity.Wire{
        ID:           model.ID,
        FlowID:       model.FlowID,
        SourceNodeID: model.SourceNodeID,
        TargetNodeID: model.TargetNodeID,
        Config:       map[string]interface{}(model.Config),
        CreatedAt:    model.CreatedAt,
        UpdatedAt:    model.UpdatedAt,
    }
}
```

---

## 6. การทำ Database Migration

### 6.1 Migration Manager

```go
// infrastructure/persistence/postgres/migration.go
package postgres

import (
    "fmt"
    "log"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "gorm.io/gorm"
)

type MigrationManager struct {
    db         *gorm.DB
    migrationsPath string
}

func NewMigrationManager(db *gorm.DB, migrationsPath string) *MigrationManager {
    return &MigrationManager{
        db:              db,
        migrationsPath: migrationsPath,
    }
}

func (m *MigrationManager) Up() error {
    sqlDB, err := m.db.DB()
    if err != nil {
        return err
    }

    driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
    if err != nil {
        return err
    }

    migrator, err := migrate.NewWithDatabaseInstance(
        fmt.Sprintf("file://%s", m.migrationsPath),
        "postgres", driver,
    )
    if err != nil {
        return err
    }

    err = migrator.Up()
    if err != nil && err != migrate.ErrNoChange {
        return err
    }

    log.Println("Migrations applied successfully")
    return nil
}

func (m *MigrationManager) Down() error {
    sqlDB, err := m.db.DB()
    if err != nil {
        return err
    }

    driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
    if err != nil {
        return err
    }

    migrator, err := migrate.NewWithDatabaseInstance(
        fmt.Sprintf("file://%s", m.migrationsPath),
        "postgres", driver,
    )
    if err != nil {
        return err
    }

    err = migrator.Down()
    if err != nil && err != migrate.ErrNoChange {
        return err
    }

    log.Println("Migrations rolled back successfully")
    return nil
}

func (m *MigrationManager) Version() (uint, bool, error) {
    sqlDB, err := m.db.DB()
    if err != nil {
        return 0, false, err
    }

    driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
    if err != nil {
        return 0, false, err
    }

    migrator, err := migrate.NewWithDatabaseInstance(
        fmt.Sprintf("file://%s", m.migrationsPath),
        "postgres", driver,
    )
    if err != nil {
        return 0, false, err
    }

    return migrator.Version()
}
```

### 6.2 การใช้งาน Migration ใน main.go

```go
// main.go
func main() {
    // ... โค้ดอื่นๆ

    // 1. เชื่อมต่อ PostgreSQL
    dbConfig := postgres.DBConfig{
        Host:     os.Getenv("DB_HOST"),
        Port:     os.Getenv("DB_PORT"),
        User:     os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASSWORD"),
        DBName:   os.Getenv("DB_NAME"),
        SSLMode:  "disable",
        Timezone: "Asia/Bangkok",
    }

    db, err := postgres.NewDBConnection(dbConfig)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // 2. ทดสอบการเชื่อมต่อ
    if err := postgres.TestConnection(db); err != nil {
        log.Fatal("Database connection test failed:", err)
    }

    // 3. รัน Migration
    migrationManager := postgres.NewMigrationManager(db, "./migrations")
    if err := migrationManager.Up(); err != nil {
        log.Fatal("Migration failed:", err)
    }

    // 4. ตรวจสอบเวอร์ชัน
    version, dirty, err := migrationManager.Version()
    if err == nil {
        log.Printf("Current migration version: %d, dirty: %v", version, dirty)
    }

    // ... โค้ดอื่นๆ
}
```

### 6.3 Run Migration ด้วย CLI

```bash
# ติดตั้ง migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# สร้าง migration file
migrate create -ext sql -dir migrations -seq create_flows_table

# รัน migration
migrate -database "postgresql://iotflow_user:password@localhost:5432/iotflow?sslmode=disable" -path migrations up

# ตรวจสอบสถานะ
migrate -database "postgresql://iotflow_user:password@localhost:5432/iotflow?sslmode=disable" -path migrations version

# Rollback 1 step
migrate -database "postgresql://iotflow_user:password@localhost:5432/iotflow?sslmode=disable" -path migrations down 1

# Rollback ทั้งหมด
migrate -database "postgresql://iotflow_user:password@localhost:5432/iotflow?sslmode=disable" -path migrations down
```

---

## 7. การตั้งค่า Environment Variables

### 7.1 .env file

```env
# .env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=iotflow_user
DB_PASSWORD=your_secure_password
DB_NAME=iotflow
DB_SSL_MODE=disable
DB_TIMEZONE=Asia/Bangkok

# Connection Pool
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=3600
DB_CONN_MAX_IDLE_TIME=600

# Migration
MIGRATIONS_PATH=./migrations
```

### 7.2 โหลด Environment Variables

```go
// config/config.go
package config

import (
    "os"
    "strconv"
    "time"

    "github.com/joho/godotenv"
)

type Config struct {
    Database DatabaseConfig
}

type DatabaseConfig struct {
    Host            string
    Port            string
    User            string
    Password        string
    DBName          string
    SSLMode         string
    Timezone        string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    ConnMaxIdleTime time.Duration
}

func LoadConfig() (*Config, error) {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        // ไม่มี .env file ให้ใช้ environment variables แทน
    }

    cfg := &Config{
        Database: DatabaseConfig{
            Host:            getEnv("DB_HOST", "localhost"),
            Port:            getEnv("DB_PORT", "5432"),
            User:            getEnv("DB_USER", "postgres"),
            Password:        getEnv("DB_PASSWORD", ""),
            DBName:          getEnv("DB_NAME", "iotflow"),
            SSLMode:         getEnv("DB_SSL_MODE", "disable"),
            Timezone:        getEnv("DB_TIMEZONE", "Asia/Bangkok"),
            MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
            MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
            ConnMaxLifetime: time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME", 3600)) * time.Second,
            ConnMaxIdleTime: time.Duration(getEnvAsInt("DB_CONN_MAX_IDLE_TIME", 600)) * time.Second,
        },
    }

    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}
```

---

## 8. PostgreSQL Configuration Optimization

### 8.1 postgresql.conf

```ini
# /etc/postgresql/15/main/postgresql.conf (Ubuntu)

# Connection Settings
max_connections = 200
superuser_reserved_connections = 3

# Memory Settings
shared_buffers = 1GB                    # 25% of system RAM
work_mem = 64MB                         # Per connection
maintenance_work_mem = 256MB
effective_cache_size = 3GB              # 75% of system RAM
shared_preload_libraries = 'pg_stat_statements'

# Write Settings
wal_buffers = 16MB
checkpoint_completion_target = 0.9
wal_writer_delay = 200ms
synchronous_commit = off               # สำหรับ dev, on สำหรับ prod

# Query Settings
random_page_cost = 1.1
effective_io_concurrency = 200

# Logging
log_destination = 'jsonlog'
logging_collector = on
log_directory = '/var/log/postgresql'
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_rotation_age = 1d
log_rotation_size = 100MB
log_min_duration_statement = 5000      # Log queries > 5 seconds
log_checkpoints = on
log_connections = on
log_disconnections = on
log_lock_waits = on
log_temp_files = 0

# Monitoringshared_preload_libraries = 'pg_stat_statements'
pg_stat_statements.track = all
pg_stat_statements.max = 10000

# Performance
enable_seqscan = on
enable_indexscan = on
enable_bitmapscan = on
```

### 8.2 Connection Pooling ด้วย PgBouncer

```ini
# pgbouncer.ini
[databases]
iotflow = host=localhost port=5432 dbname=iotflow

[pgbouncer]
listen_addr = *
listen_port = 6432
auth_type = md5
auth_file = /etc/pgbouncer/userlist.txt
admin_users = postgres

pool_mode = transaction
default_pool_size = 20
max_client_conn = 1000
reserve_pool_size = 5
reserve_pool_timeout = 3
```

---

## 9. การ Backup และ Restore

### 9.1 Backup Script

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR="/backups/postgresql"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="iotflow"
DB_USER="iotflow_user"

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup
pg_dump -U $DB_USER -h localhost -d $DB_NAME \
    --format=custom \
    --compress=9 \
    --file="$BACKUP_DIR/${DB_NAME}_${DATE}.backup"

# Cleanup old backups (keep last 7 days)
find $BACKUP_DIR -name "*.backup" -mtime +7 -delete

echo "Backup completed: ${BACKUP_DIR}/${DB_NAME}_${DATE}.backup"
```

### 9.2 Restore Script

```bash
#!/bin/bash
# restore.sh

BACKUP_FILE=$1
DB_NAME="iotflow"
DB_USER="iotflow_user"

if [ -z "$BACKUP_FILE" ]; then
    echo "Usage: ./restore.sh <backup_file>"
    exit 1
fi

# Restore
pg_restore -U $DB_USER -h localhost -d $DB_NAME \
    --clean \
    --if-exists \
    $BACKUP_FILE

echo "Restore completed: $BACKUP_FILE"
```

### 9.3 Automated Backup (Cron)

```bash
# Add to crontab
# Daily backup at 2:00 AM
0 2 * * * /path/to/backup.sh
```

---

## 10. Monitoring และ Troubleshooting

### 10.1 Query Performance Monitoring

```sql
-- Check slow queries
SELECT 
    query,
    calls,
    total_exec_time,
    mean_exec_time,
    stddev_exec_time,
    rows
FROM pg_stat_statements
ORDER BY total_exec_time DESC
LIMIT 20;

-- Check active connections
SELECT 
    pid,
    usename,
    application_name,
    client_addr,
    state,
    query,
    now() - query_start AS duration
FROM pg_stat_activity
WHERE state = 'active'
ORDER BY duration DESC;

-- Check table sizes
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size('public.' || tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size('public.' || tablename) DESC;

-- Check index usage
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;
```

### 10.2 Troubleshooting Common Issues

```sql
-- Kill a long-running query
SELECT pid, query, state, now() - query_start AS duration
FROM pg_stat_activity 
WHERE state = 'active' AND now() - query_start > interval '5 minutes';

-- Kill specific query
SELECT pg_terminate_backend(pid);

-- Check locks
SELECT 
    locked.relation::regclass AS table_name,
    locked.mode,
    locked.granted,
    locker.usename AS locker_user,
    locker.query AS locker_query,
    blocked.usename AS blocked_user,
    blocked.query AS blocked_query
FROM pg_locks locked
JOIN pg_stat_activity locker ON locker.pid = locked.pid
JOIN pg_locks blocked_lock ON blocked_lock.relation = locked.relation AND blocked_lock.pid != locked.pid
JOIN pg_stat_activity blocked ON blocked.pid = blocked_lock.pid
WHERE NOT locked.granted;

-- VACUUM and ANALYZE
VACUUM ANALYZE;
VACUUM FULL;  -- Use with caution
```

---

## 11. Docker Setup

### 11.1 Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: iotflow-postgres
    environment:
      POSTGRES_USER: iotflow_user
      POSTGRES_PASSWORD: your_secure_password
      POSTGRES_DB: iotflow
      POSTGRES_INITDB_ARGS: "--encoding=UTF8 --locale=en_US.UTF-8"
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d
      - ./postgresql.conf:/etc/postgresql/postgresql.conf
    command: postgres -c config_file=/etc/postgresql/postgresql.conf
    networks:
      - iotflow-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U iotflow_user -d iotflow"]
      interval: 10s
      timeout: 5s
      retries: 5

  pgbouncer:
    image: edoburu/pgbouncer:latest
    container_name: iotflow-pgbouncer
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_NAME: iotflow
      DB_USER: iotflow_user
      DB_PASSWORD: your_secure_password
      POOL_MODE: transaction
      DEFAULT_POOL_SIZE: 20
      RESERVE_POOL_SIZE: 5
      MAX_CLIENT_CONN: 1000
    ports:
      - "6432:5432"
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - iotflow-network
    restart: unless-stopped

  pgadmin:
    image: dpage/pgadmin4:latest
    container_name: iotflow-pgadmin
    environment:
      PGADMIN_DEFAULT_EMAIL: admin@iotflow.com
      PGADMIN_DEFAULT_PASSWORD: admin_password
    ports:
      - "5050:80"
    depends_on:
      - postgres
    networks:
      - iotflow-network
    restart: unless-stopped

volumes:
  postgres_data:

networks:
  iotflow-network:
    driver: bridge
```

### 11.2 .env สำหรับ Docker

```env
# .env.docker
POSTGRES_USER=iotflow_user
POSTGRES_PASSWORD=your_secure_password
POSTGRES_DB=iotflow
POSTGRES_PORT=5432

PGADMIN_DEFAULT_EMAIL=admin@iotflow.com
PGADMIN_DEFAULT_PASSWORD=admin_password
```

### 11.3 Docker Commands

```bash
# Start all services
docker-compose up -d

# Start only PostgreSQL
docker-compose up -d postgres

# Check logs
docker-compose logs -f postgres

# Connect to PostgreSQL inside container
docker exec -it iotflow-postgres psql -U iotflow_user -d iotflow

# Run migration
docker-compose exec postgres psql -U iotflow_user -d iotflow -f /docker-entrypoint-initdb.d/001_create_flows_table.up.sql

# Backup database
docker exec -t iotflow-postgres pg_dump -U iotflow_user iotflow > backup.sql

# Restore database
cat backup.sql | docker exec -i iotflow-postgres psql -U iotflow_user iotflow

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

---

## สรุป

| Component | การตั้งค่า |
|-----------|-----------|
| **Database** | PostgreSQL 15+ |
| **Connection** | GORM + pgx driver |
| **Migration** | golang-migrate |
| **Connection Pool** | MaxOpenConns=100, MaxIdleConns=10 |
| **Connection String** | `host=localhost user=iotflow_user password=*** dbname=iotflow port=5432 sslmode=disable` |
| **Indexes** | Primary keys, Foreign keys, Composite indexes, GIN indexes for JSONB |
| **Backup** | pg_dump daily (cron) |
| **Monitoring** | pg_stat_statements, pg_stat_activity |
| **Docker** | Docker Compose with PostgreSQL, PgBouncer, pgAdmin |

# สรุปและแนวทางพัฒนาเพิ่มเติม (ฉบับสมบูรณ์)

## 12. สรุปและแนวทางพัฒนาเพิ่มเติม

### 12.1 สรุประบบ

เราได้สร้างระบบ IoT Flow Management Module ที่สมบูรณ์พร้อมความสามารถ:

| หมวดหมู่ | ความสามารถ | รายละเอียด |
|----------|------------|------------|
| **สถาปัตยกรรม** | Event‑Driven + Clean Architecture | แยกส่วนการทำงาน ช่วยให้ขยายขนาดและบำรุงรักษาง่าย |
| **Message Bus** | Apache Kafka | รองรับปริมาณข้อมูลสูง แยก Producer/Consumer |
| **Time‑series Storage** | InfluxDB | เก็บข้อมูลอนุกรมเวลาสำหรับวิเคราะห์และแสดงผล |
| **Search & Logging** | Elasticsearch | พร้อมความสามารถ Vector Search ด้วย dense_vector |
| **Vector Database** | Milvus/Qdrant | สำหรับ Semantic Search ประสิทธิภาพสูง |
| **Real‑time Dashboard** | Socket.IO | แสดงผลข้อมูลแบบเรียลไทม์ |
| **AI‑Powered** | LLM + Embedding | สรุปผล, ตอบคำถาม, ค้นหาเชิงความหมาย (RAG) |
| **Immutable Audit** | Blockchain | บันทึกหลักฐานที่ไม่สามารถแก้ไขได้ |
| **Logging Management** | Aggregator + Scheduler | เก็บสถิติรายปีและลบข้อมูลเก่าอัตโนมัติ |
| **Database** | PostgreSQL | เก็บข้อมูลหลัก (Flow, Node, Wire, Aggregate) |
| **Caching** | Redis | Cache, Rate Limit, Distributed State |

---

### 12.2 แนวทางพัฒนาเพิ่มเติม

#### 12.2.1 ด้าน Frontend / UI

| ลำดับ | งาน | รายละเอียด | ความสำคัญ |
|-------|-----|------------|-----------|
| 1 | Dashboard UI | สร้างหน้าแสดงผลแบบ Real‑time ด้วย Socket.IO แสดงกราฟ Sensor Data, ผลลัพธ์ LLM, การค้นหาเวกเตอร์ | สูง |
| 2 | Flow Builder | สร้าง Drag-and-Drop UI สำหรับออกแบบ Flow (Node.js + React/Angular) | สูง |
| 3 | Kibana Integration | ใช้ Kibana สำหรับ visualize logs และทำ vector search | ปานกลาง |
| 4 | Mobile App | พัฒนาแอปพลิเคชันมือถือ (React Native/Flutter) สำหรับดูข้อมูลแบบเรียลไทม์ | ต่ำ |
| 5 | Custom Alert UI | หน้าแสดงและจัดการ Alert พร้อมการกรองและแจ้งเตือน | ปานกลาง |

#### 12.2.2 ด้าน Backend / Infrastructure

| ลำดับ | งาน | รายละเอียด | ความสำคัญ |
|-------|-----|------------|-----------|
| 6 | Batch Blockchain Transactions | รวมหลาย requests เข้าด้วยกันเพื่อลดค่า Gas | ปานกลาง |
| 7 | Caching Embeddings | ใช้ Redis แคชเวกเตอร์ที่สร้างแล้วเพื่อลดภาระการคำนวณ | ปานกลาง |
| 8 | Circuit Breaker & Retry | สำหรับการเรียก LLM/Embedding API เพื่อเพิ่มความเสถียร | สูง |
| 9 | Auto‑scaling | ปรับจำนวน Kafka Consumers และ Workers ตามโหลด (KEDA) | ปานกลาง |
| 10 | Multi‑tenant | รองรับหลายองค์กร/ผู้ใช้ ด้วยการแยกข้อมูลผ่าน `flow_id` และ `user_id` | สูง |
| 11 | Integration Tests | เขียนทดสอบการทำงานร่วมกันของทุก component | สูง |
| 12 | Monitoring | ใช้ Prometheus + Grafana เพื่อตรวจสอบสถานะของระบบและ Kafka Lag | สูง |
| 13 | Data Retention Policy | ปรับแต่ง Retention ของ InfluxDB และ ILM ของ Elasticsearch | ปานกลาง |
| 14 | Webhook Support | รองรับการส่งข้อมูลไปยัง Webhook ของ第三方 | ต่ำ |
| 15 | API Gateway | เพิ่ม API Gateway (Kong/Traefik) สำหรับจัดการ API และ Rate Limit | ปานกลาง |
| 16 | Distributed Tracing | เพิ่ม OpenTelemetry + Jaeger สำหรับติดตาม request | ปานกลาง |
| 17 | Secret Management | ใช้ HashiCorp Vault สำหรับจัดการ Secrets | ต่ำ |
| 18 | Blue-Green Deployment | รองรับการ Deploy แบบไม่หยุดระบบ | ปานกลาง |

#### 12.2.3 ด้าน Data / Analytics

| ลำดับ | งาน | รายละเอียด | ความสำคัญ |
|-------|-----|------------|-----------|
| 19 | Real‑time Analytics | ใช้ Flink หรือ Spark Streaming สำหรับวิเคราะห์ข้อมูลแบบเรียลไทม์ | ปานกลาง |
| 20 | Anomaly Detection | ใช้ Machine Learning สำหรับตรวจจับความผิดปกติของ Sensor Data | สูง |
| 21 | Predictive Maintenance | ทำนายการเสียหายของอุปกรณ์จากข้อมูล Sensor | สูง |
| 22 | Data Lake Integration | เก็บข้อมูลดิบใน Data Lake (S3/MinIO) สำหรับการวิเคราะห์เชิงลึก | ปานกลาง |
| 23 | BI Dashboard | เชื่อมต่อกับ Tableau/Power BI สำหรับรายงานเชิงธุรกิจ | ต่ำ |

#### 12.2.4 ด้าน Security

| ลำดับ | งาน | รายละเอียด | ความสำคัญ |
|-------|-----|------------|-----------|
| 24 | JWT Authentication | เพิ่มระบบ Authentication และ Authorization | สูง |
| 25 | RBAC | ระบบควบคุมสิทธิ์ตามบทบาท (Admin, User, Viewer) | สูง |
| 26 | Audit Trail | บันทึกทุกการกระทำในระบบ (ใคร ทำอะไร เมื่อไหร่) | สูง |
| 27 | Data Encryption | เข้ารหัสข้อมูลที่ sensitive ใน Database | ปานกลาง |
| 28 | API Security | เพิ่ม CORS, CSRF Protection, Rate Limiting ที่เข้มงวด | สูง |
| 29 | OAuth2/SSO | รองรับการ Login ด้วย Google, Microsoft, etc. | ต่ำ |

#### 12.2.5 ด้าน DevOps / CI-CD

| ลำดับ | งาน | รายละเอียด | ความสำคัญ |
|-------|-----|------------|-----------|
| 30 | CI/CD Pipeline | ใช้ GitHub Actions/GitLab CI สำหรับ Build, Test, Deploy อัตโนมัติ | สูง |
| 31 | Kubernetes | จัดการ Container ด้วย K8s (Helm Charts) | สูง |
| 32 | Infrastructure as Code | ใช้ Terraform สำหรับจัดการ Infrastructure | ปานกลาง |
| 33 | Log Aggregation | รวม Logs ทั้งระบบไว้ที่ศูนย์กลาง (ELK Stack) | สูง |
| 34 | Alerting | ตั้ง Alert เมื่อระบบมีปัญหาผ่าน Slack/Email/PagerDuty | สูง |

---

### 12.3 ส่วนที่ขาดและควรเพิ่มเติม

#### 12.3.1 Domain Layer ที่ขาด

```go
// domain/entity/wire.go (เพิ่ม)
package entity

type Wire struct {
    ID           int64
    FlowID       int64
    SourceNodeID int64
    TargetNodeID int64
    Config       map[string]interface{}
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

// domain/value_object/flow_status.go (เพิ่ม)
type FlowStatus string

const (
    FlowStatusDraft    FlowStatus = "draft"
    FlowStatusDeployed FlowStatus = "deployed"
    FlowStatusRunning  FlowStatus = "running"
    FlowStatusStopped  FlowStatus = "stopped"
    FlowStatusError    FlowStatus = "error"
)

// domain/value_object/node_status.go (เพิ่ม)
type NodeStatus string

const (
    NodeStatusActive   NodeStatus = "active"
    NodeStatusInactive NodeStatus = "inactive"
    NodeStatusError    NodeStatus = "error"
)

// domain/errors/errors.go (เพิ่ม)
var (
    ErrFlowNotFound      = errors.New("flow not found")
    ErrFlowAlreadyExists = errors.New("flow already exists")
    ErrFlowNotDeployed   = errors.New("flow not deployed")
    ErrFlowAlreadyRunning = errors.New("flow already running")
    ErrNodeNotFound      = errors.New("node not found")
    ErrWireNotFound      = errors.New("wire not found")
    ErrInvalidNodeType   = errors.New("invalid node type")
)
```

#### 12.3.2 Application Layer ที่ขาด

```go
// application/deploy_flow.go (เพิ่ม)
type DeployFlowUseCase struct {
    repo     repository.FlowRepository
    validator *service.FlowValidator
}

func (uc *DeployFlowUseCase) Execute(ctx context.Context, flowID int64) error {
    flow, err := uc.repo.FindByID(ctx, flowID)
    if err != nil {
        return err
    }
    if flow.Status != value_object.FlowStatusDraft {
        return errors.ErrFlowNotDraft
    }
    if err := uc.validator.Validate(flow); err != nil {
        return err
    }
    flow.Status = value_object.FlowStatusDeployed
    return uc.repo.Update(ctx, flow)
}

// application/add_node.go (เพิ่ม)
type AddNodeUseCase struct {
    repo repository.FlowRepository
}

func (uc *AddNodeUseCase) Execute(ctx context.Context, flowID int64, name string, nodeType value_object.NodeType, config map[string]interface{}) (*entity.Node, error) {
    flow, err := uc.repo.FindByID(ctx, flowID)
    if err != nil {
        return nil, err
    }
    if flow.Status != value_object.FlowStatusDraft {
        return nil, errors.ErrFlowNotDraft
    }
    node := &entity.Node{
        FlowID: flowID,
        Name:   name,
        Type:   nodeType,
        Config: config,
        Status: value_object.NodeStatusActive,
    }
    // เพิ่ม node และบันทึก
    return node, nil
}

// application/remove_node.go (เพิ่ม)
// application/add_wire.go (เพิ่ม)
// application/remove_wire.go (เพิ่ม)
// application/update_flow.go (เพิ่ม)
// application/delete_flow.go (เพิ่ม)
// application/get_flow.go (เพิ่ม)
// application/list_flows.go (เพิ่ม)
// application/generate_aggregates.go (เพิ่ม)
```

#### 12.3.3 Infrastructure Layer ที่ขาด

```go
// infrastructure/engine/runtime.go (เพิ่ม)
package engine

type Runtime struct {
    instanceStore *redis.FlowInstanceStore
    instanceID    string
    executors     map[value_object.NodeType]NodeExecutor
    runningFlows  sync.Map
    kafkaClient   *messaging.KafkaClient
}

func NewRuntime(store *redis.FlowInstanceStore, instanceID string, kafka *messaging.KafkaClient) *Runtime {
    return &Runtime{
        instanceStore: store,
        instanceID:    instanceID,
        executors:     make(map[value_object.NodeType]NodeExecutor),
        runningFlows:  sync.Map{},
        kafkaClient:   kafka,
    }
}

func (r *Runtime) RegisterExecutor(nodeType value_object.NodeType, executor NodeExecutor) {
    r.executors[nodeType] = executor
}

func (r *Runtime) StartFlow(ctx context.Context, flow *entity.Flow) error {
    // ตรวจสอบว่า flow กำลังรันอยู่หรือไม่
    if _, ok := r.runningFlows.Load(flow.ID); ok {
        return errors.ErrFlowAlreadyRunning
    }
    
    // เริ่มต้นการทำงานของแต่ละ Node
    for _, node := range flow.Nodes {
        if err := r.startNode(ctx, flow, node); err != nil {
            return err
        }
    }
    
    r.runningFlows.Store(flow.ID, true)
    return r.instanceStore.SetRunning(ctx, flow.ID, r.instanceID)
}

func (r *Runtime) StopFlow(ctx context.Context, flowID int64) error {
    if _, ok := r.runningFlows.Load(flowID); !ok {
        return errors.ErrFlowNotRunning
    }
    r.runningFlows.Delete(flowID)
    return r.instanceStore.ClearRunning(ctx, flowID, r.instanceID)
}

// infrastructure/engine/node_runner.go (เพิ่ม)
type NodeRunner struct {
    node     *entity.Node
    executor NodeExecutor
    inputCh  <-chan value_object.DataPayload
    outputCh chan<- value_object.DataPayload
}

func (nr *NodeRunner) Run(ctx context.Context) error {
    output, err := nr.executor.Execute(ctx, nr.node, nr.inputCh)
    if err != nil {
        return err
    }
    go func() {
        for payload := range output {
            select {
            case nr.outputCh <- payload:
            case <-ctx.Done():
                return
            }
        }
    }()
    return nil
}

// NodeExecutor interface
type NodeExecutor interface {
    Execute(ctx context.Context, node *entity.Node, input <-chan value_object.DataPayload) (<-chan value_object.DataPayload, error)
}
```

#### 12.3.4 MQTT Adapter ที่ขาด

```go
// infrastructure/mqtt/mqtt_adapter.go (เพิ่ม)
package mqtt

import (
    "context"
    "encoding/json"
    "time"

    mqtt "github.com/eclipse/paho.mqtt.golang"
    "github.com/yourproject/internal/modules/iotflow/domain/value_object"
    "github.com/yourproject/internal/modules/iotflow/infrastructure/messaging"
)

type MQTTAdapter struct {
    client       mqtt.Client
    kafkaClient  *messaging.KafkaClient
    flowEngine   *engine.Runtime
}

func NewMQTTAdapter(broker string, clientID string, kafka *messaging.KafkaClient, engine *engine.Runtime) (*MQTTAdapter, error) {
    opts := mqtt.NewClientOptions()
    opts.AddBroker(broker)
    opts.SetClientID(clientID)
    opts.SetAutoReconnect(true)
    opts.SetConnectRetry(true)
    opts.SetConnectTimeout(30 * time.Second)
    
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        return nil, token.Error()
    }
    
    adapter := &MQTTAdapter{
        client:      client,
        kafkaClient: kafka,
        flowEngine:  engine,
    }
    
    return adapter, nil
}

func (a *MQTTAdapter) Subscribe(topic string, flowID string) error {
    token := a.client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
        var payload value_object.DataPayload
        if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
            return
        }
        payload.FlowID = flowID
        payload.Timestamp = time.Now()
        a.kafkaClient.Publish(context.Background(), "flow.data", payload)
    })
    return token.Error()
}

func (a *MQTTAdapter) Publish(topic string, payload interface{}) error {
    data, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    token := a.client.Publish(topic, 1, false, data)
    return token.Error()
}
```

#### 12.3.5 Logging System ที่ขาด

```go
// infrastructure/logging/logger.go (เพิ่ม)
package logging

import (
    "context"
    "encoding/json"
    "time"

    "github.com/yourproject/internal/modules/iotflow/infrastructure/search"
)

type LogEntry struct {
    Timestamp   time.Time              `json:"timestamp"`
    Level       string                 `json:"level"` // info, warn, error, debug
    FlowID      int64                  `json:"flow_id,omitempty"`
    NodeID      int64                  `json:"node_id,omitempty"`
    DeviceID    string                 `json:"device_id,omitempty"`
    Message     string                 `json:"message"`
    Error       string                 `json:"error,omitempty"`
    StackTrace  string                 `json:"stack_trace,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
    UserID      string                 `json:"user_id,omitempty"`
    IPAddress   string                 `json:"ip_address,omitempty"`
    Duration    int64                  `json:"duration_ms,omitempty"`
}

type Logger struct {
    esClient *search.ESClient
    index    string
}

func NewLogger(esClient *search.ESClient, index string) *Logger {
    return &Logger{
        esClient: esClient,
        index:    index,
    }
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...interface{}) {
    l.log(ctx, "info", msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...interface{}) {
    l.log(ctx, "warn", msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...interface{}) {
    l.log(ctx, "error", msg, fields...)
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...interface{}) {
    l.log(ctx, "debug", msg, fields...)
}

func (l *Logger) log(ctx context.Context, level string, msg string, fields ...interface{}) {
    entry := LogEntry{
        Timestamp: time.Now(),
        Level:     level,
        Message:   msg,
        Metadata:  make(map[string]interface{}),
    }
    
    for i := 0; i < len(fields); i += 2 {
        if i+1 < len(fields) {
            key, ok := fields[i].(string)
            if ok {
                entry.Metadata[key] = fields[i+1]
            }
        }
    }
    
    // Index to Elasticsearch (async)
    go func() {
        data, _ := json.Marshal(entry)
        l.esClient.IndexDocument(l.index, data)
    }()
}
```

---

### 12.4 Prompt สำหรับการพัฒนาในอนาคต

#### Prompt 1: เพิ่ม Feature ใหม่ (LLM Chat Interface)

```
ต้องการเพิ่มฟีเจอร์ Chat Interface สำหรับให้ผู้ใช้สอบถามข้อมูล Sensor แบบ Natural Language 
โดยใช้ LLM (RAG) และ Vector Database

รายละเอียด:
1. ผู้ใช้สามารถพิมพ์คำถาม เช่น "อุณหภูมิเฉลี่ยของเซ็นเซอร์ A ในเดือนที่แล้วคือเท่าไหร่?"
2. ระบบจะแปลงคำถามเป็น Vector Search Query
3. ค้นหาข้อมูลที่เกี่ยวข้องจาก Vector Database
4. ใช้ LLM สร้างคำตอบจากข้อมูลที่ค้นพบ
5. แสดงผลผ่าน Socket.IO แบบ Real-time

สิ่งที่ต้องเพิ่ม:
- HTTP Endpoint: POST /api/v1/iot/chat
- LLM Prompt Template สำหรับ RAG
- Vector Search Integration
- WebSocket Event: chat:response
- Rate Limiting สำหรับ Chat API
- การเก็บประวัติการสนทนา

ไฟล์ที่เกี่ยวข้อง:
- internal/modules/iotflow/interfaces/http/chat_handler.go (ใหม่)
- internal/modules/iotflow/application/chat_use_case.go (ใหม่)
- internal/modules/iotflow/infrastructure/llm/rag_client.go (ใหม่)
- frontend/src/app/features/chat/ (ใหม่)
```

#### Prompt 2: เพิ่ม Feature ใหม่ (Automated Alert Rules)

```
ต้องการเพิ่มระบบ Automated Alert Rules ที่ผู้ใช้สามารถกำหนดเงื่อนไขการแจ้งเตือนได้เอง

รายละเอียด:
1. ผู้ใช้สร้าง Rule ด้วย GUI (เช่น "แจ้งเตือนเมื่ออุณหภูมิ > 80°C ติดต่อกัน 3 ครั้ง")
2. ระบบประเมิน Rules กับข้อมูล Sensor แบบ Real-time (ใช้ Kafka Streams หรือ Flink)
3. เมื่อเงื่อนไขตรง จะส่ง Alert ผ่าน Socket.IO และ Email/Line/Slack
4. เก็บประวัติการแจ้งเตือนใน Elasticsearch
5. แสดง Dashboard สรุป Alert

สิ่งที่ต้องเพิ่ม:
- Alert Rule Entity และ CRUD
- Rule Engine (CEL หรือ JSONata)
- Kafka Streams Processor
- Notification Adapters (Email, Line, Slack)
- Alert Dashboard UI

ไฟล์ที่เกี่ยวข้อง:
- internal/modules/iotflow/domain/entity/alert_rule.go (ใหม่)
- internal/modules/iotflow/application/create_alert_rule.go (ใหม่)
- internal/modules/iotflow/infrastructure/streams/alert_processor.go (ใหม่)
- internal/modules/iotflow/infrastructure/notifications/ (ใหม่)
- frontend/src/app/features/alerts/ (ใหม่)
```

#### Prompt 3: ปรับปรุง Performance (Data Partitioning)

```
ต้องการปรับปรุงประสิทธิภาพการ query ข้อมูล Sensor Data ที่มีปริมาณมาก (>100 ล้านเรคคอร์ด)

รายละเอียด:
1. แบ่ง Partition ใน InfluxDB ตามช่วงเวลา (monthly partitioning)
2. ใช้ Time-series Indexing ใน PostgreSQL
3. เพิ่ม Materialized Views สำหรับข้อมูลที่ query บ่อย
4. ใช้ Read Replica สำหรับ Report Queries
5. ใช้ TimescaleDB extension แทน InfluxDB (option)

สิ่งที่ต้องเพิ่ม:
- Database Migration สำหรับ Partitioning
- Query Optimizer
- Read Replica Configuration
- TimescaleDB Integration (optional)

ไฟล์ที่เกี่ยวข้อง:
- migrations/008_partition_sensor_data.up.sql (ใหม่)
- infrastructure/persistence/postgres/partition_manager.go (ใหม่)
- infrastructure/timeseries/query_optimizer.go (ใหม่)
- config/database_replica.go (ใหม่)
```

#### Prompt 4: เพิ่ม Feature ใหม่ (Device Management)

```
ต้องการเพิ่มระบบจัดการอุปกรณ์ IoT (Device Management)

รายละเอียด:
1. ลงทะเบียนอุปกรณ์ (Device ID, Type, Location, Tags)
2. ดูสถานะอุปกรณ์ (Online/Offline, Last Seen, Battery Level)
3. OTA (Over-The-Air) Update Firmware
4. Device Group Management
5. Device Shadow (State Sync)
6. Device Lifecycle Management (Provision, Activate, Decommission)

สิ่งที่ต้องเพิ่ม:
- Device Entity และ CRUD
- Device Registry Service
- Device Shadow Store (Redis)
- OTA Update Service
- Device Group Service
- Device Management UI

ไฟล์ที่เกี่ยวข้อง:
- internal/modules/iotflow/domain/entity/device.go (ใหม่)
- internal/modules/iotflow/application/device/ (ใหม่)
- internal/modules/iotflow/infrastructure/device/ (ใหม่)
- frontend/src/app/features/devices/ (ใหม่)
```

#### Prompt 5: เพิ่ม Feature ใหม่ (Data Export & Report)

```
ต้องการเพิ่มระบบ Export และ Report ข้อมูล Sensor

รายละเอียด:
1. Export ข้อมูลเป็น CSV, Excel, JSON, PDF
2. สร้าง Report แบบ Scheduled (รายวัน/รายสัปดาห์/รายเดือน)
3. ส่ง Report ทาง Email หรือจัดเก็บใน Cloud Storage
4. Report Template ที่ปรับแต่งได้
5. Data Visualization ใน Report (Chart, Table)

สิ่งที่ต้องเพิ่ม:
- Report Entity และ CRUD
- Report Generator Service
- Export Service (CSV, Excel, PDF)
- Email Sender
- Report Scheduler

ไฟล์ที่เกี่ยวข้อง:
- internal/modules/iotflow/domain/entity/report.go (ใหม่)
- internal/modules/iotflow/application/report/ (ใหม่)
- internal/modules/iotflow/infrastructure/export/ (ใหม่)
- internal/modules/iotflow/infrastructure/report/ (ใหม่)
- frontend/src/app/features/reports/ (ใหม่)
```

#### Prompt 6: เพิ่ม Feature ใหม่ (Data Pipeline Monitoring)

```
ต้องการเพิ่มระบบ Monitoring สำหรับ Data Pipeline

รายละเอียด:
1. แสดง Kafka Consumer Lag
2. แสดงอัตราการรับ/ส่งข้อมูล (Throughput)
3. แสดงความล่าช้า (Latency) ของแต่ละ Stage
4. Alert เมื่อ Pipeline มีปัญหา
5. History ของ Pipeline Performance

สิ่งที่ต้องเพิ่ม:
- Kafka Admin Client
- Metrics Collector (Prometheus)
- Dashboard UI
- Alert Rules

ไฟล์ที่เกี่ยวข้อง:
- internal/modules/iotflow/infrastructure/monitoring/kafka_monitor.go (ใหม่)
- internal/modules/iotflow/infrastructure/monitoring/pipeline_metrics.go (ใหม่)
- internal/modules/iotflow/interfaces/http/monitoring_handler.go (ใหม่)
- frontend/src/app/features/monitoring/ (ใหม่)
```

#### Prompt 7: เพิ่ม Feature ใหม่ (Data Simulation)

```
ต้องการเพิ่มระบบจำลองข้อมูล Sensor สำหรับ Testing และ Demo

รายละเอียด:
1. สร้าง Simulation Scenario (Device, Pattern, Frequency)
2. กำหนดรูปแบบข้อมูล (Random, Sine Wave, Step, etc.)
3. กำหนดความถี่ในการส่งข้อมูล
4. เริ่ม/หยุด Simulation
5. ดูผลลัพธ์แบบ Real-time

สิ่งที่ต้องเพิ่ม:
- Simulation Entity และ CRUD
- Data Generator Service
- Simulation Scheduler
- Simulation Dashboard

ไฟล์ที่เกี่ยวข้อง:
- internal/modules/iotflow/domain/entity/simulation.go (ใหม่)
- internal/modules/iotflow/application/simulation/ (ใหม่)
- internal/modules/iotflow/infrastructure/simulator/ (ใหม่)
- frontend/src/app/features/simulation/ (ใหม่)
```

#### Prompt 8: ปรับปรุง Security (Zero Trust)

```
ต้องการปรับปรุงระบบ Security เป็น Zero Trust Architecture

รายละเอียด:
1. Mutual TLS (mTLS) สำหรับการสื่อสารระหว่าง Service
2. Service Mesh (Istio/Linkerd)
3. JWT with Short-lived Tokens + Refresh Token
4. API Key Management
5. Network Policy (Kubernetes Network Policies)
6. Secrets Encryption (Vault)
7. Security Audit Log

สิ่งที่ต้องเพิ่ม:
- mTLS Configuration
- Service Mesh Setup
- Token Service
- API Key Service
- Network Policies
- Vault Integration

ไฟล์ที่เกี่ยวข้อง:
- config/tls_config.go (ใหม่)
- infrastructure/security/ (ใหม่)
- kubernetes/network-policies/ (ใหม่)
- kubernetes/service-mesh/ (ใหม่)
```

---

### 12.5 ตารางสรุป Feature ที่ควรเพิ่ม

| ลำดับ | Feature | Priority | Complexity | Dependencies |
|-------|---------|----------|------------|--------------|
| 1 | Dashboard UI | ⭐⭐⭐ | กลาง | Socket.IO |
| 2 | JWT Authentication | ⭐⭐⭐ | กลาง | - |
| 3 | Flow Builder UI | ⭐⭐⭐ | สูง | - |
| 4 | Alert Rules | ⭐⭐⭐ | สูง | Kafka Streams |
| 5 | Device Management | ⭐⭐⭐ | กลาง | - |
| 6 | Data Export/Report | ⭐⭐ | กลาง | - |
| 7 | Anomaly Detection | ⭐⭐⭐ | สูง | ML Model |
| 8 | CI/CD Pipeline | ⭐⭐⭐ | กลาง | - |
| 9 | Kubernetes Deployment | ⭐⭐⭐ | สูง | Docker |
| 10 | Monitoring (Prometheus/Grafana) | ⭐⭐⭐ | กลาง | - |
| 11 | Chat Interface (RAG) | ⭐⭐ | สูง | LLM, Vector DB |
| 12 | Multi-tenant | ⭐⭐⭐ | สูง | - |
| 13 | Data Simulation | ⭐⭐ | กลาง | - |
| 14 | Predictive Maintenance | ⭐⭐ | สูง | ML Model |
| 15 | OTA Update | ⭐⭐ | สูง | Device Management |
| 16 | Zero Trust Security | ⭐⭐ | สูง | Service Mesh |
| 17 | Data Pipeline Monitoring | ⭐⭐ | กลาง | Kafka |
| 18 | Mobile App | ⭐ | สูง | - |
| 19 | Webhook Support | ⭐ | ต่ำ | - |
| 20 | OAuth2/SSO | ⭐ | กลาง | - |

---

### 12.6 บทสรุป

ระบบ IoT Flow Management Module ที่ออกแบบนี้เป็นรากฐานที่แข็งแกร่งสำหรับการจัดการข้อมูล IoT แบบเรียลไทม์ ด้วยสถาปัตยกรรม Event-Driven ที่แยกส่วนการทำงานอย่างชัดเจน ทำให้สามารถ:

1. **ขยายขนาดได้ง่าย** - แต่ละ component สามารถ scale ได้อิสระ
2. **บำรุงรักษาง่าย** - Clean Architecture ทำให้การเปลี่ยนแปลงไม่กระทบส่วนอื่น
3. **รองรับเทคโนโลยีใหม่** - สามารถเพิ่ม LLM, Embedding, Blockchain ได้โดยไม่ต้องเปลี่ยนโครงสร้างหลัก
4. **มีความยืดหยุ่นสูง** - รองรับการเชื่อมต่อกับระบบอื่นๆ ได้ง่าย

แนวทางพัฒนาเพิ่มเติมที่ระบุไว้จะช่วยให้ระบบเติบโตไปพร้อมกับความต้องการทางธุรกิจ และสามารถปรับตัวตามเทคโนโลยีที่เปลี่ยนแปลงไปได้อย่างมีประสิทธิภาพ

# ระบบ IoT Flow Management Module สำหรับ ERP - ดีไซน์ใหม่เพื่อการขาย

## สารบัญ

1. [Executive Summary](#1-executive-summary)
2. [Product Positioning](#2-product-positioning)
3. [ระบบสถาปัตยกรรมใหม่](#3-ระบบสถาปัตยกรรมใหม่)
4. [ERP Modules ที่พัฒนาเอง](#4-erp-modules-ที่พัฒนาเอง)
5. [Enterprise Features](#5-enterprise-features)
6. [Packaging & Pricing](#6-packaging-pricing)
7. [Implementation Roadmap](#7-implementation-roadmap)
8. [Go-to-Market Strategy](#8-go-to-market-strategy)
9. [Technical Architecture - v2.0](#9-technical-architecture---v20)
10. [Code Implementation](#10-code-implementation)

---

## 1. Executive Summary

### 1.1 Vision
**"เปลี่ยนโรงงานทุกแห่งให้เป็น Smart Factory ด้วย IoT และ AI ที่ใช้งานได้จริง"**

### 1.2 Value Proposition

| สำหรับลูกค้า | คุณค่าที่ได้รับ |
|-------------|---------------|
| **ผู้บริหาร** | ตัดสินใจด้วยข้อมูลจริง (Data-Driven Decision) ลดต้นทุน 15-30% |
| **ฝ่ายผลิต** | เพิ่มประสิทธิภาพการผลิต (OEE) ลด Downtime |
| **ฝ่ายซ่อมบำรุง** | Predictive Maintenance ลดการเสียหายฉุกเฉิน |
| **ฝ่ายคุณภาพ** | ตรวจจับความผิดปกติได้ทันที ลดของเสีย |
| **ฝ่ายไอที** | ระบบที่เสถียร ปลอดภัย ขยายได้ รองรับมาตรฐานสากล |

### 1.3 เป้าหมายทางการตลาด

| ระยะ | เป้าหมาย |
|------|----------|
| ปีที่ 1 | 30 โรงงาน, รายได้ 30-50 ล้านบาท |
| ปีที่ 2 | 100 โรงงาน, รายได้ 150-200 ล้านบาท |
| ปีที่ 3 | 300 โรงงาน, รายได้ 500+ ล้านบาท |

---

## 2. Product Positioning

### 2.1 ชื่อผลิตภัณฑ์
**"IQ-Factory"** - Intelligent Factory Platform

### 2.2 จุดแข็งทางการแข่งขัน

| ปัจจัย | IQ-Factory | Competitor A | Competitor B |
|--------|------------|--------------|--------------|
| ราคา | ★★★★★ | ★★ | ★★★ |
| ความง่ายในการติดตั้ง | ★★★★★ | ★★★ | ★★ |
| ครบวงจร (IoT + ERP) | ★★★★★ | ★★★ | ★★ |
| AI/ML ในตัว | ★★★★★ | ★★ | ★★★ |
| การปรับแต่งได้ | ★★★★★ | ★★★ | ★★★★ |
| ระบบรักษาความปลอดภัย | ★★★★★ | ★★★★ | ★★★ |
| ระบบ Blockchain | ★★★★★ | ★ | ★ |
| ภาษาไทย | ★★★★★ | ★ | ★★ |

### 2.3 กลุ่มลูกค้าเป้าหมาย

| กลุ่ม | ขนาด | ลักษณะ | การตัดสินใจซื้อ |
|------|------|--------|---------------|
| **โรงงานขนาดใหญ่** | 500+ | มี IT ใน-house, ต้องการระบบครบวงจร | 6-12 เดือน |
| **โรงงานขนาดกลาง** | 100-500 | ต้องการระบบสำเร็จรูป ปรับแต่งได้ | 3-6 เดือน |
| **โรงงานขนาดเล็ก** | <100 | ต้องการระบบราคาประหยัด ใช้งานง่าย | 1-3 เดือน |
| **นิคมอุตสาหกรรม** | กลุ่ม | ซื้อเป็นกลุ่มเพื่อบริหารจัดการรวม | 6-12 เดือน |

---

## 3. ระบบสถาปัตยกรรมใหม่

### 3.1 สถาปัตยกรรมรวม (v2.0 - Enterprise Edition)

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         IQ-Factory Enterprise Platform                          │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  ┌──────────────────────────────────────────────────────────────────────────┐   │
│  │                        UI Layer (Angular/React)                          │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐     │   │
│  │  │ Dashboard │ │  Flow    │ │  ERP     │ │  Report  │ │  Admin   │     │   │
│  │  │          │ │ Builder  │ │ Modules  │ │  Center  │ │  Console │     │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘     │   │
│  └──────────────────────────────────────────────────────────────────────────┘   │
│                                    │                                            │
│  ┌─────────────────────────────────▼───────────────────────────────────────┐    │
│  │                          API Gateway (Kong)                             │    │
│  │                    Authentication, Rate Limit, Logging                  │    │
│  └─────────────────────────────────┬───────────────────────────────────────┘    │
│                                    │                                            │
│  ┌──────────────┬─────────────────┼─────────────────┬──────────────────────┐    │
│  │              │                 │                 │                      │    │
│  ▼              ▼                 ▼                 ▼                      │    │
│┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐          │    │
││  IoT     │ │  Flow    │ │  ERP     │ │  AI/ML   │ │  Analytics│          │    │
││  Core    │ │  Engine  │ │  Modules │ │  Engine  │ │  Engine  │          │    │
│└──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘          │    │
│     │            │            │            │            │                 │    │
│     └────────────┴────────────┼────────────┴────────────┘                 │    │
│                               │                                            │    │
│                    ┌──────────▼──────────┐                                 │    │
│                    │   Event Bus (Kafka)  │                                 │    │
│                    └──────────┬──────────┘                                 │    │
│                               │                                            │    │
│  ┌────────────┬───────────────┼───────────────┬──────────────────────┐    │    │
│  │            │               │               │                      │    │    │
│  ▼            ▼               ▼               ▼                      │    │    │
│┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐    │    │    │
││PostgreSQL│ │ InfluxDB │ │ Elastic │ │  Redis   │ │  Vector  │    │    │    │
││ (Main)   │ │ (Time)   │ │ (Search)│ │ (Cache)  │ │   DB     │    │    │    │
│└──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘    │    │    │
│                                                                        │    │
│  ┌──────────────────────────────────────────────────────────────┐     │    │
│  │                  External Integrations                       │     │    │
│  │  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐          │     │    │
│  │  │MQTT  │ │OPC UA│ │Modbus│ │SAP  │ │Line  │ │Email │          │     │    │
│  │  └─────┘ └─────┘ └─────┘ └─────┘ └─────┘ └─────┘          │     │    │
│  └──────────────────────────────────────────────────────────────┘     │    │
│                                                                        │    │
└─────────────────────────────────────────────────────────────────────────┘    │
```

### 3.2 Core Modules

| Module | Responsibility | Technology |
|--------|---------------|------------|
| **IoT Core** | Device Management, Data Ingestion | Go, MQTT, Kafka |
| **Flow Engine** | Data Processing Pipeline | Go, Kafka Streams |
| **ERP Modules** | Business Logic | Go, PostgreSQL |
| **AI/ML Engine** | Predictions, Anomaly Detection | Python, TensorFlow |
| **Analytics Engine** | Reports, Dashboards | Go, Elasticsearch |
| **API Gateway** | Authentication, Routing | Kong, JWT |

---

## 4. ERP Modules ที่พัฒนาเอง

### 4.1 โมดูลหลัก

```
┌─────────────────────────────────────────────────────────────────┐
│                    IQ-Factory ERP Modules                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │  Production  │  │  Maintenance │  │  Quality     │        │
│  │  Management  │  │  Management  │  │  Management  │        │
│  ├──────────────┤  ├──────────────┤  ├──────────────┤        │
│  │ • MES        │  │ • CMMS       │  │ • SPC        │        │
│  │ • OEE        │  │ • PM Schedule│  │ • QC Check   │        │
│  │ • Production │  │ • Work Orders│  │ • Non-conform│        │
│  │   Planning   │  │ • Spare Parts│  │ • CAPA       │        │
│  │ • Downtime   │  │ • History    │  │ • Calibration│        │
│  └──────────────┘  └──────────────┘  └──────────────┘        │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │  Inventory   │  │  Procurement │  │  Sales &     │        │
│  │  Management  │  │  Management  │  │  Distribution│        │
│  ├──────────────┤  ├──────────────┤  ├──────────────┤        │
│  │ • Stock      │  │ • PO         │  │ • SO         │        │
│  │ • Warehouse  │  │ • Supplier   │  │ • Delivery   │        │
│  │ • FIFO/LIFO  │  │ • RFQ        │  │ • Invoice    │        │
│  │ • Stocktake  │  │ • Contract   │  │ • Customer   │        │
│  └──────────────┘  └──────────────┘  └──────────────┘        │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │  HR          │  │  Finance     │  │  Analytics   │        │
│  │  Management  │  │  & Accounting│  │  & Dashboard │        │
│  ├──────────────┤  ├──────────────┤  ├──────────────┤        │
│  │ • Employee   │  │ • GL         │  │ • KPI        │        │
│  │ • Attendance │  │ • AP/AR      │  │ • Reports    │        │
│  │ • Payroll    │  │ • Budget     │  │ • Forecast   │        │
│  │ • Training   │  │ • Asset      │  │ • Benchmark  │        │
│  └──────────────┘  └──────────────┘  └──────────────┘        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 รายละเอียดโมดูล

#### 4.2.1 Production Management (MES)

```go
// domain/entity/production.go
package entity

type ProductionOrder struct {
    ID            int64     `json:"id"`
    OrderNo       string    `json:"order_no"`
    ProductID     int64     `json:"product_id"`
    Quantity      float64   `json:"quantity"`
    StartDate     time.Time `json:"start_date"`
    EndDate       time.Time `json:"end_date"`
    Status        string    `json:"status"` // pending, running, completed, canceled
    Priority      int       `json:"priority"`
    MachineID     int64     `json:"machine_id"`
    OperatorID    int64     `json:"operator_id"`
    BatchNo       string    `json:"batch_no"`
    ActualQuantity float64  `json:"actual_quantity"`
    ScrapQuantity  float64  `json:"scrap_quantity"`
    DefectRate     float64  `json:"defect_rate"`
    OEE            float64  `json:"oee"` // Overall Equipment Effectiveness
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

type OEE struct {
    ID            int64     `json:"id"`
    MachineID     int64     `json:"machine_id"`
    Date          time.Time `json:"date"`
    Availability  float64   `json:"availability"`   // 0-1
    Performance   float64   `json:"performance"`    // 0-1
    Quality       float64   `json:"quality"`        // 0-1
    OverallOEE    float64   `json:"overall_oee"`    // 0-1
    PlannedTime   float64   `json:"planned_time"`   // minutes
    OperatingTime float64   `json:"operating_time"` // minutes
    IdealTime     float64   `json:"ideal_time"`     // minutes
    GoodCount     int64     `json:"good_count"`
    TotalCount    int64     `json:"total_count"`
    DowntimeReasons map[string]float64 `json:"downtime_reasons"`
    CreatedAt     time.Time `json:"created_at"`
}
```

#### 4.2.2 Maintenance Management (CMMS)

```go
// domain/entity/maintenance.go
package entity

type MaintenanceWorkOrder struct {
    ID            int64     `json:"id"`
    WorkOrderNo   string    `json:"work_order_no"`
    MachineID     int64     `json:"machine_id"`
    Type          string    `json:"type"` // preventive, corrective, emergency
    Priority      string    `json:"priority"` // low, medium, high, critical
    Description   string    `json:"description"`
    ReportedBy    int64     `json:"reported_by"`
    AssignedTo    int64     `json:"assigned_to"`
    StartDate     time.Time `json:"start_date"`
    EndDate       time.Time `json:"end_date"`
    Status        string    `json:"status"` // open, in_progress, completed, canceled
    Cost          float64   `json:"cost"`
    SpareParts    []SparePartUsage `json:"spare_parts"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type PreventiveMaintenanceSchedule struct {
    ID            int64     `json:"id"`
    MachineID     int64     `json:"machine_id"`
    ScheduleType  string    `json:"schedule_type"` // daily, weekly, monthly, yearly
    Interval      int       `json:"interval"`
    LastRun       time.Time `json:"last_run"`
    NextRun       time.Time `json:"next_run"`
    Tasks         []string  `json:"tasks"`
    EstimatedDuration int   `json:"estimated_duration"` // minutes
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type SparePart struct {
    ID            int64     `json:"id"`
    PartNo        string    `json:"part_no"`
    Name          string    `json:"name"`
    Category      string    `json:"category"`
    Quantity      int       `json:"quantity"`
    MinQuantity   int       `json:"min_quantity"`
    MaxQuantity   int       `json:"max_quantity"`
    Location      string    `json:"location"`
    SupplierID    int64     `json:"supplier_id"`
    UnitPrice     float64   `json:"unit_price"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}
```

#### 4.2.3 Quality Management

```go
// domain/entity/quality.go
package entity

type QualityCheckPoint struct {
    ID            int64     `json:"id"`
    ProductID     int64     `json:"product_id"`
    ProcessID     int64     `json:"process_id"`
    Parameter     string    `json:"parameter"`
    Specification string    `json:"specification"`
    UpperLimit    float64   `json:"upper_limit"`
    LowerLimit    float64   `json:"lower_limit"`
    Unit          string    `json:"unit"`
    CheckMethod   string    `json:"check_method"`
    Frequency     string    `json:"frequency"`
    CreatedAt     time.Time `json:"created_at"`
}

type QualityResult struct {
    ID            int64     `json:"id"`
    BatchNo       string    `json:"batch_no"`
    ProductID     int64     `json:"product_id"`
    CheckPointID  int64     `json:"check_point_id"`
    Value         float64   `json:"value"`
    Status        string    `json:"status"` // pass, fail, warning
    InspectorID   int64     `json:"inspector_id"`
    CheckedAt     time.Time `json:"checked_at"`
    CreatedAt     time.Time `json:"created_at"`
}

type NonConformance struct {
    ID            int64     `json:"id"`
    NCNo          string    `json:"nc_no"`
    BatchNo       string    `json:"batch_no"`
    ProductID     int64     `json:"product_id"`
    Severity      string    `json:"severity"` // minor, major, critical
    Description   string    `json:"description"`
    RootCause     string    `json:"root_cause"`
    CorrectiveAction string `json:"corrective_action"`
    Status        string    `json:"status"` // open, investigating, resolved, closed
    ReportedBy    int64     `json:"reported_by"`
    ResolvedBy    int64     `json:"resolved_by"`
    CreatedAt     time.Time `json:"created_at"`
    ResolvedAt    time.Time `json:"resolved_at"`
}
```

#### 4.2.4 Inventory Management

```go
// domain/entity/inventory.go
package entity

type InventoryItem struct {
    ID            int64     `json:"id"`
    Code          string    `json:"code"`
    Name          string    `json:"name"`
    Category      string    `json:"category"`
    Unit          string    `json:"unit"`
    Quantity      float64   `json:"quantity"`
    Reserved      float64   `json:"reserved"`
    MinStock      float64   `json:"min_stock"`
    MaxStock      float64   `json:"max_stock"`
    ReorderPoint  float64   `json:"reorder_point"`
    Location      string    `json:"location"`
    CostMethod    string    `json:"cost_method"` // fifo, lifo, average
    UnitCost      float64   `json:"unit_cost"`
    SupplierID    int64     `json:"supplier_id"`
    LastUpdated   time.Time `json:"last_updated"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type InventoryTransaction struct {
    ID            int64     `json:"id"`
    ItemID        int64     `json:"item_id"`
    TransactionType string  `json:"transaction_type"` // in, out, adjustment
    Quantity      float64   `json:"quantity"`
    UnitCost      float64   `json:"unit_cost"`
    TotalCost     float64   `json:"total_cost"`
    ReferenceType string    `json:"reference_type"` // po, so, adjustment
    ReferenceID   int64     `json:"reference_id"`
    Description   string    `json:"description"`
    CreatedBy     int64     `json:"created_by"`
    CreatedAt     time.Time `json:"created_at"`
}
```

#### 4.2.5 Procurement Management

```go
// domain/entity/procurement.go
package entity

type PurchaseOrder struct {
    ID            int64     `json:"id"`
    PONo          string    `json:"po_no"`
    SupplierID    int64     `json:"supplier_id"`
    OrderDate     time.Time `json:"order_date"`
    ExpectedDate  time.Time `json:"expected_date"`
    Status        string    `json:"status"` // draft, sent, confirmed, received, closed
    Items         []POItem  `json:"items"`
    SubTotal      float64   `json:"sub_total"`
    Tax           float64   `json:"tax"`
    Total         float64   `json:"total"`
    Currency      string    `json:"currency"`
    PaymentTerms  string    `json:"payment_terms"`
    ShippingTerms string    `json:"shipping_terms"`
    CreatedBy     int64     `json:"created_by"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type POItem struct {
    ID            int64     `json:"id"`
    POID          int64     `json:"po_id"`
    ItemID        int64     `json:"item_id"`
    Quantity      int       `json:"quantity"`
    Received      int       `json:"received"`
    UnitPrice     float64   `json:"unit_price"`
    TotalPrice    float64   `json:"total_price"`
    DeliveryDate  time.Time `json:"delivery_date"`
}

type Supplier struct {
    ID            int64     `json:"id"`
    Code          string    `json:"code"`
    Name          string    `json:"name"`
    TaxID         string    `json:"tax_id"`
    Address       string    `json:"address"`
    Phone         string    `json:"phone"`
    Email         string    `json:"email"`
    ContactPerson string    `json:"contact_person"`
    Rating        float64   `json:"rating"` // 0-5
    Status        string    `json:"status"` // active, inactive
    PaymentTerms  string    `json:"payment_terms"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}
```

#### 4.2.6 Financial Management

```go
// domain/entity/finance.go
package entity

type ChartOfAccount struct {
    ID            int64     `json:"id"`
    AccountNo     string    `json:"account_no"`
    AccountName   string    `json:"account_name"`
    AccountType   string    `json:"account_type"` // asset, liability, equity, revenue, expense
    ParentID      *int64    `json:"parent_id"`
    IsActive      bool      `json:"is_active"`
    CreatedAt     time.Time `json:"created_at"`
}

type JournalEntry struct {
    ID            int64     `json:"id"`
    EntryNo       string    `json:"entry_no"`
    EntryDate     time.Time `json:"entry_date"`
    Description   string    `json:"description"`
    Reference     string    `json:"reference"`
    Lines         []JournalLine `json:"lines"`
    CreatedBy     int64     `json:"created_by"`
    ApprovedBy    int64     `json:"approved_by"`
    Status        string    `json:"status"` // draft, posted, canceled
    CreatedAt     time.Time `json:"created_at"`
    PostedAt      *time.Time `json:"posted_at"`
}

type JournalLine struct {
    ID            int64     `json:"id"`
    EntryID       int64     `json:"entry_id"`
    AccountID     int64     `json:"account_id"`
    DebitAmount   float64   `json:"debit_amount"`
    CreditAmount  float64   `json:"credit_amount"`
    Description   string    `json:"description"`
}

type Budget struct {
    ID            int64     `json:"id"`
    Year          int       `json:"year"`
    AccountID     int64     `json:"account_id"`
    DepartmentID  int64     `json:"department_id"`
    Amount        float64   `json:"amount"`
    ActualAmount  float64   `json:"actual_amount"`
    Variance      float64   `json:"variance"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}
```

---

## 5. Enterprise Features

### 5.1 Multi-Tenant Architecture

```go
// domain/entity/tenant.go
package entity

type Tenant struct {
    ID            int64     `json:"id"`
    Code          string    `json:"code"`
    Name          string    `json:"name"`
    Database      string    `json:"database"` // สำหรับแยก DB
    Schema        string    `json:"schema"`   // หรือแยก Schema
    Status        string    `json:"status"`   // active, inactive, suspended
    Config        TenantConfig `json:"config"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type TenantConfig struct {
    MaxUsers      int       `json:"max_users"`
    MaxDevices    int       `json:"max_devices"`
    MaxFlows      int       `json:"max_flows"`
    Features      []string  `json:"features"` // which modules enabled
    DataRetention int       `json:"data_retention"` // days
    BackupSchedule string   `json:"backup_schedule"`
}
```

### 5.2 User Management & RBAC

```go
// domain/entity/user.go
package entity

type User struct {
    ID            int64     `json:"id"`
    TenantID      int64     `json:"tenant_id"`
    Email         string    `json:"email"`
    Username      string    `json:"username"`
    PasswordHash  string    `json:"-"`
    FirstName     string    `json:"first_name"`
    LastName      string    `json:"last_name"`
    Role          string    `json:"role"` // admin, manager, operator, viewer
    Department    string    `json:"department"`
    Phone         string    `json:"phone"`
    IsActive      bool      `json:"is_active"`
    LastLoginAt   *time.Time `json:"last_login_at"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type Permission struct {
    ID            int64     `json:"id"`
    Resource      string    `json:"resource"`
    Action        string    `json:"action"` // create, read, update, delete
    Description   string    `json:"description"`
}

type RolePermission struct {
    Role         string    `json:"role"`
    PermissionID int64     `json:"permission_id"`
}

// RBAC Matrix
/*
Role           | Resources                              | Actions
---------------|----------------------------------------|----------------
Admin          | All                                    | All
Manager        | Production, Maintenance, Quality, Inventory | CRUD
Manager        | Finance, HR, Procurement               | Read Only
Operator       | Production, Maintenance                | Create, Read
Operator       | Quality                               | Create, Read
Viewer         | All                                   | Read Only
*/
```

### 5.3 Audit Trail

```go
// domain/entity/audit.go
package entity

type AuditLog struct {
    ID            int64     `json:"id"`
    TenantID      int64     `json:"tenant_id"`
    UserID        int64     `json:"user_id"`
    Action        string    `json:"action"`
    ResourceType  string    `json:"resource_type"`
    ResourceID    string    `json:"resource_id"`
    OldValue      JSONMap   `json:"old_value"`
    NewValue      JSONMap   `json:"new_value"`
    IPAddress     string    `json:"ip_address"`
    UserAgent     string    `json:"user_agent"`
    CreatedAt     time.Time `json:"created_at"`
}
```

### 5.4 Notification System

```go
// domain/entity/notification.go
package entity

type Notification struct {
    ID            int64     `json:"id"`
    TenantID      int64     `json:"tenant_id"`
    UserID        int64     `json:"user_id"`
    Type          string    `json:"type"` // alert, reminder, report, system
    Severity      string    `json:"severity"` // info, warning, error, success
    Title         string    `json:"title"`
    Message       string    `json:"message"`
    Data          JSONMap   `json:"data"`
    IsRead        bool      `json:"is_read"`
    ReadAt        *time.Time `json:"read_at"`
    CreatedAt     time.Time `json:"created_at"`
}

type NotificationChannel struct {
    ID            int64     `json:"id"`
    TenantID      int64     `json:"tenant_id"`
    Type          string    `json:"type"` // email, line, slack, webhook
    Config        JSONMap   `json:"config"`
    IsEnabled     bool      `json:"is_enabled"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}
```

---

## 6. Packaging & Pricing

### 6.1 ระดับผลิตภัณฑ์

| Feature | Starter | Professional | Enterprise | Ultimate |
|---------|---------|--------------|------------|----------|
| **ราคา/เดือน** | ฿9,900 | ฿29,900 | ฿99,900 | ฿299,900 |
| **ราคา/ปี** | ฿99,000 | ฿299,000 | ฿999,000 | ฿2,999,000 |
| **จำนวนอุปกรณ์** | 10 | 50 | 500 | Unlimited |
| **จำนวน User** | 5 | 20 | 100 | Unlimited |
| **Data Retention** | 30 วัน | 90 วัน | 365 วัน | Unlimited |
| **IoT Core** | ✅ | ✅ | ✅ | ✅ |
| **Flow Engine** | ✅ | ✅ | ✅ | ✅ |
| **Dashboard** | ✅ | ✅ | ✅ | ✅ |
| **MES Module** | ❌ | ✅ | ✅ | ✅ |
| **CMMS Module** | ❌ | ✅ | ✅ | ✅ |
| **Quality Module** | ❌ | ✅ | ✅ | ✅ |
| **Inventory Module** | ❌ | ✅ | ✅ | ✅ |
| **Procurement Module** | ❌ | ❌ | ✅ | ✅ |
| **Finance Module** | ❌ | ❌ | ✅ | ✅ |
| **HR Module** | ❌ | ❌ | ✅ | ✅ |
| **AI/ML Engine** | ❌ | ❌ | ✅ | ✅ |
| **Blockchain** | ❌ | ❌ | ✅ | ✅ |
| **Multi-Tenant** | ❌ | ❌ | ❌ | ✅ |
| **Custom Integration** | ❌ | ❌ | ❌ | ✅ |
| **On-Premise** | ❌ | ❌ | ❌ | ✅ |
| **24/7 Support** | ❌ | ❌ | ✅ | ✅ |
| **SLA 99.9%** | ❌ | ❌ | ✅ | ✅ |

### 6.2 โมดูลเสริม (Add-ons)

| Module | Price/Year | Description |
|--------|-----------|-------------|
| **AI Predictive Maintenance** | ฿200,000 | ใช้ ML พยากรณ์การเสียหายของเครื่องจักร |
| **Computer Vision** | ฿300,000 | ตรวจจับคุณภาพด้วยกล้อง AI |
| **Energy Management** | ฿150,000 | วิเคราะห์และลดการใช้พลังงาน |
| **Supply Chain Integration** | ฿250,000 | เชื่อมต่อกับ Supplier/Customer |
| **Mobile App** | ฿100,000 | แอปพลิเคชันมือถือสำหรับผู้ปฏิบัติงาน |
| **Custom Report** | ฿50,000 | สร้างรายงานเฉพาะตามความต้องการ |
| **Training & Onboarding** | ฿100,000 | อบรมและติดตั้งระบบ |

### 6.3 บริการ Professional Services

| Service | Price | Description |
|---------|-------|-------------|
| **Implementation** | ฿300,000-1,000,000 | ติดตั้งและปรับแต่งระบบตามความต้องการ |
| **Integration** | ฿200,000-500,000 | เชื่อมต่อกับระบบเดิม (SAP, ERP, MES) |
| **Custom Development** | ฿100,000/เดือน | พัฒนาฟีเจอร์เพิ่มเติมตามความต้องการ |
| **Training** | ฿50,000-200,000 | อบรมการใช้งานระบบ |
| **Consulting** | ฿5,000/ชั่วโมง | ให้คำปรึกษาเกี่ยวกับ Smart Factory |

---

## 7. Implementation Roadmap

### 7.1 ระยะเวลาการพัฒนา

```
Phase 1: Foundation (2-3 เดือน)
├── Core Platform (IoT Core, Flow Engine, API Gateway)
├── PostgreSQL + InfluxDB + Elasticsearch Setup
├── Basic Dashboard
└── User Management + Authentication

Phase 2: ERP Modules (3-4 เดือน)
├── Production Management (MES)
├── Maintenance Management (CMMS)
├── Quality Management
└── Inventory Management

Phase 3: Advanced Features (2-3 เดือน)
├── Procurement Management
├── Finance & Accounting
├── HR Management
├── AI/ML Engine
└── Blockchain Integration

Phase 4: Enterprise Features (2-3 เดือน)
├── Multi-Tenant
├── Advanced Reporting
├── Notification System
├── Mobile App
└── Integration APIs

Phase 5: Production Ready (1-2 เดือน)
├── Performance Testing
├── Security Audit
├── Documentation
├── Training Materials
└── Deployment Automation
```

### 7.2 Sprint Plan

| Sprint | Duration | Deliverables |
|--------|----------|--------------|
| Sprint 1-2 | 4 weeks | Project Setup, Database Schema, Core Entities |
| Sprint 3-4 | 4 weeks | IoT Core, MQTT, Kafka Integration |
| Sprint 5-6 | 4 weeks | Flow Engine, Node Executors |
| Sprint 7-8 | 4 weeks | API Gateway, Authentication, Basic UI |
| Sprint 9-10 | 4 weeks | Production Module (MES) |
| Sprint 11-12 | 4 weeks | Maintenance Module (CMMS) |
| Sprint 13-14 | 4 weeks | Quality Module |
| Sprint 15-16 | 4 weeks | Inventory Module |
| Sprint 17-18 | 4 weeks | Procurement Module |
| Sprint 19-20 | 4 weeks | Finance Module |
| Sprint 21-22 | 4 weeks | AI/ML Engine |
| Sprint 23-24 | 4 weeks | Enterprise Features |
| Sprint 25-26 | 4 weeks | Testing, Documentation, Deployment |

---

## 8. Go-to-Market Strategy

### 8.1 ช่องทางการขาย

| Channel | Strategy | Timeline |
|---------|----------|----------|
| **Direct Sales** | ทีมขายของเราเอง เน้นโรงงานขนาดใหญ่ | Q1-Q4 |
| **Channel Partners** | Partner กับ SI, IT Consultant, System Integrator | Q2-Q4 |
| **Digital Marketing** | SEO, Google Ads, LinkedIn, Content Marketing | Ongoing |
| **Events** | งานสัมมนา, งานแสดงสินค้า (Manufacturing Expo) | Q2, Q4 |
| **Referral Program** | ส่วนลด 10% สำหรับลูกค้าที่แนะนำเพื่อน | Ongoing |

### 8.2 กลยุทธ์การตลาด

| ระยะ | กิจกรรม | เป้าหมาย |
|------|---------|----------|
| **Q1** | เปิดตัวผลิตภัณฑ์, งานสัมมนาเปิดตัว | สร้างการรับรู้ |
| **Q2** | โรดโชว์ 5 จังหวัด, งาน Manufacturing Expo | ได้ลูกค้ารายแรก 5 ราย |
| **Q3** | Case Study, Webinar, Partner Training | ได้ลูกค้า 20 ราย |
| **Q4** | Campaign สิ้นปี, งาน IoT Conference | ได้ลูกค้า 30 ราย |

### 8.3 เป้าหมายการขาย

| ปี | เป้าหมายลูกค้า | รายได้ | ทีมขาย |
|----|---------------|--------|--------|
| ปีที่ 1 | 30 | 30-50 ล้านบาท | 3-5 คน |
| ปีที่ 2 | 100 | 150-200 ล้านบาท | 10-15 คน |
| ปีที่ 3 | 300 | 500+ ล้านบาท | 20-30 คน |

---

## 9. Technical Architecture - v2.0

### 9.1 Technology Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| **Backend** | Go (Golang) | High performance, concurrent processing |
| **API Gateway** | Kong | Authentication, Rate Limit, Routing |
| **Message Bus** | Apache Kafka | Event-driven architecture |
| **Main Database** | PostgreSQL 15+ | Transactional data |
| **Time-series DB** | InfluxDB 2.x | Sensor data, metrics |
| **Search Engine** | Elasticsearch 8.x | Logs, search, vector search |
| **Vector DB** | Milvus/Qdrant | Embedding, semantic search |
| **Cache** | Redis 7.x | Caching, session, rate limit |
| **AI/ML** | Python + TensorFlow/PyTorch | Model training, predictions |
| **Frontend** | Angular 17 / React 18 | Web application |
| **Mobile** | React Native / Flutter | Mobile application |
| **Container** | Docker + Kubernetes | Deployment, scaling |
| **Monitoring** | Prometheus + Grafana | Metrics, monitoring |
| **Logging** | ELK Stack | Centralized logging |
| **CI/CD** | GitLab CI / GitHub Actions | Automation |

### 9.2 High-Level Architecture Diagram

```
┌──────────────────────────────────────────────────────────────────────────────────────┐
│                              Client Layer                                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │    Web UI   │  │  Mobile App │  │  Desktop App│  │  API Client │              │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘              │
└──────────────────────────────────────┬───────────────────────────────────────────────┘
                                       │
┌──────────────────────────────────────▼───────────────────────────────────────────────┐
│                              API Gateway Layer                                      │
│  ┌─────────────────────────────────────────────────────────────────────────────────┐│
│  │                              Kong Gateway                                      ││
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐           ││
│  │  │ Auth    │  │ Rate    │  │ Cache   │  │ Logging │  │ Circuit │           ││
│  │  │         │  │ Limit   │  │         │  │         │  │ Breaker │           ││
│  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘  └─────────┘           ││
│  └─────────────────────────────────────────────────────────────────────────────────┘│
└──────────────────────────────────────┬───────────────────────────────────────────────┘
                                       │
┌──────────────────────────────────────┼───────────────────────────────────────────────┐
│                                      │                                               │
│  ┌───────────────────────────────────▼───────────────────────────────────────────┐  │
│  │                         Service Layer (Microservices)                         │  │
│  │                                                                               │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐        │  │
│  │  │  Auth    │ │  Tenant  │ │   IoT    │ │   Flow   │ │   ERP    │        │  │
│  │  │ Service  │ │  Service │ │  Core    │ │  Engine  │ │  Service │        │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘        │  │
│  │                                                                               │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐        │  │
│  │  │   AI     │ │  Report  │ │  Notify  │ │  Audit   │ │  Export  │        │  │
│  │  │ Service  │ │  Service │ │  Service │ │  Service │ │  Service │        │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘        │  │
│  └───────────────────────────────────────────────────────────────────────────────┘  │
│                                       │                                              │
└───────────────────────────────────────┼──────────────────────────────────────────────┘
                                        │
┌───────────────────────────────────────▼──────────────────────────────────────────────┐
│                              Event Bus Layer                                        │
│  ┌─────────────────────────────────────────────────────────────────────────────────┐│
│  │                            Apache Kafka                                        ││
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐           ││
│  │  │ iot.data│  │ flow    │  │ erp     │  │ audit   │  │ notify  │           ││
│  │  │         │  │ events  │  │ events  │  │         │  │         │           ││
│  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘  └─────────┘           ││
│  └─────────────────────────────────────────────────────────────────────────────────┘│
└───────────────────────────────────────┬──────────────────────────────────────────────┘
                                        │
┌───────────────────────────────────────┼──────────────────────────────────────────────┐
│  ┌───────────────────────────────────▼───────────────────────────────────────────┐  │
│  │                        Data Layer                                              │  │
│  │                                                                               │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐        │  │
│  │  │PostgreSQL│ │ InfluxDB │ │Elastic   │ │  Redis   │ │  Vector  │        │  │
│  │  │ (Main)   │ │ (Time)   │ │ (Search) │ │ (Cache)  │ │   DB     │        │  │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘        │  │
│  └───────────────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────────────────┘
```

### 9.3 Data Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              Data Flow                                              │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│   Device ──► MQTT ──► IoT Core ──► Kafka ──► Flow Engine ──► Kafka ──► Consumer    │
│                              │                        │                            │
│                              ▼                        ▼                            │
│                         PostgreSQL              InfluxDB + ES                      │
│                         (Device)                (Sensor Data)                      │
│                              │                        │                            │
│                              ▼                        ▼                            │
│   User ──► API ──► Flow Engine ──► Process Data ──► Save to DB                      │
│                              │                                                        │
│                              ▼                                                        │
│                         AI/ML Engine                                                 │
│                         (Predictions)                                                │
│                              │                                                        │
│                              ▼                                                        │
│                         Notification Service                                        │
│                         (Alert, Report)                                              │
│                              │                                                        │
│                              ▼                                                        │
│                         User (Dashboard, Email, Line)                               │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 10. Code Implementation

### 10.1 Multi-Tenant Middleware

```go
// infrastructure/middleware/tenant.go
package middleware

import (
    "context"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/yourproject/internal/modules/iotflow/domain/entity"
)

type TenantKey string

const (
    TenantContextKey TenantKey = "tenant"
)

func TenantMiddleware(tenantService TenantService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get tenant from header or subdomain
        tenantCode := c.GetHeader("X-Tenant-Id")
        if tenantCode == "" {
            // Extract from subdomain
            host := c.Request.Host
            parts := strings.Split(host, ".")
            if len(parts) >= 3 {
                tenantCode = parts[0]
            }
        }

        if tenantCode == "" {
            c.JSON(400, gin.H{"error": "Tenant ID is required"})
            c.Abort()
            return
        }

        tenant, err := tenantService.GetByCode(c.Request.Context(), tenantCode)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid tenant"})
            c.Abort()
            return
        }

        if tenant.Status != "active" {
            c.JSON(403, gin.H{"error": "Tenant is not active"})
            c.Abort()
            return
        }

        ctx := context.WithValue(c.Request.Context(), TenantContextKey, tenant)
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}

func GetTenant(ctx context.Context) *entity.Tenant {
    if tenant, ok := ctx.Value(TenantContextKey).(*entity.Tenant); ok {
        return tenant
    }
    return nil
}
```

### 10.2 RBAC Middleware

```go
// infrastructure/middleware/auth.go
package middleware

import (
    "context"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
    jwtSecret []byte
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
    return &AuthMiddleware{jwtSecret: []byte(secret)}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Authorization header is required"})
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(401, gin.H{"error": "Invalid token format"})
            c.Abort()
            return
        }

        tokenString := parts[1]
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return m.jwtSecret, nil
        })

        if err != nil || !token.Valid {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(401, gin.H{"error": "Invalid claims"})
            c.Abort()
            return
        }

        c.Set("user_id", claims["user_id"])
        c.Set("tenant_id", claims["tenant_id"])
        c.Set("role", claims["role"])
        c.Next()
    }
}

func (m *AuthMiddleware) Authorize(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists {
            c.JSON(403, gin.H{"error": "Role not found"})
            c.Abort()
            return
        }

        // Hierarchical RBAC: admin > manager > operator > viewer
        roleHierarchy := map[string]int{
            "admin":    4,
            "manager":  3,
            "operator": 2,
            "viewer":   1,
        }

        roleVal, ok1 := roleHierarchy[role.(string)]
        requiredVal, ok2 := roleHierarchy[requiredRole]
        if !ok1 || !ok2 || roleVal < requiredVal {
            c.JSON(403, gin.H{"error": "Insufficient permissions"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### 10.3 Production Service Implementation

```go
// application/production/production_service.go
package production

import (
    "context"
    "time"

    "github.com/yourproject/internal/modules/erp/domain/entity"
    "github.com/yourproject/internal/modules/erp/domain/repository"
    "github.com/yourproject/internal/modules/iotflow/infrastructure/messaging"
)

type ProductionService struct {
    orderRepo    repository.ProductionOrderRepository
    oeeRepo      repository.OEERepository
    machineRepo  repository.MachineRepository
    kafkaClient  *messaging.KafkaClient
}

func NewProductionService(
    orderRepo repository.ProductionOrderRepository,
    oeeRepo repository.OEERepository,
    machineRepo repository.MachineRepository,
    kafka *messaging.KafkaClient,
) *ProductionService {
    return &ProductionService{
        orderRepo:   orderRepo,
        oeeRepo:     oeeRepo,
        machineRepo: machineRepo,
        kafkaClient: kafka,
    }
}

func (s *ProductionService) CreateProductionOrder(ctx context.Context, order *entity.ProductionOrder) error {
    // Generate order number
    order.OrderNo = generateOrderNo()
    order.Status = "pending"
    order.CreatedAt = time.Now()

    if err := s.orderRepo.Create(ctx, order); err != nil {
        return err
    }

    // Publish event
    event := map[string]interface{}{
        "event_type": "production_order_created",
        "order_id":   order.ID,
        "order_no":   order.OrderNo,
        "product_id": order.ProductID,
        "quantity":   order.Quantity,
        "timestamp":  time.Now(),
    }
    s.kafkaClient.Publish(ctx, "erp.production.events", event)

    return nil
}

func (s *ProductionService) CalculateOEE(ctx context.Context, machineID int64, date time.Time) (*entity.OEE, error) {
    // Get production data
    orders, err := s.orderRepo.GetByMachineAndDate(ctx, machineID, date)
    if err != nil {
        return nil, err
    }

    machine, err := s.machineRepo.GetByID(ctx, machineID)
    if err != nil {
        return nil, err
    }

    // Calculate OEE
    plannedTime := 8 * 60.0 // 8 hours in minutes
    operatingTime := plannedTime - machine.Downtime
    availability := operatingTime / plannedTime

    idealCycleTime := machine.IdealCycleTime // minutes per unit
    totalUnits := 0.0
    goodUnits := 0.0
    for _, order := range orders {
        totalUnits += order.ActualQuantity
        goodUnits += order.ActualQuantity - order.ScrapQuantity
    }

    performance := (totalUnits * idealCycleTime) / operatingTime
    quality := goodUnits / totalUnits

    oee := &entity.OEE{
        MachineID:     machineID,
        Date:          date,
        Availability:  availability,
        Performance:   performance,
        Quality:       quality,
        OverallOEE:    availability * performance * quality,
        PlannedTime:   plannedTime,
        OperatingTime: operatingTime,
        IdealTime:     totalUnits * idealCycleTime,
        GoodCount:     int64(goodUnits),
        TotalCount:    int64(totalUnits),
        CreatedAt:     time.Now(),
    }

    if err := s.oeeRepo.Save(ctx, oee); err != nil {
        return nil, err
    }

    // Publish OEE event
    event := map[string]interface{}{
        "event_type":    "oee_calculated",
        "machine_id":    machineID,
        "date":          date,
        "oee":           oee.OverallOEE,
        "availability":  oee.Availability,
        "performance":   oee.Performance,
        "quality":       oee.Quality,
        "timestamp":     time.Now(),
    }
    s.kafkaClient.Publish(ctx, "erp.production.events", event)

    return oee, nil
}

func generateOrderNo() string {
    return "PO-" + time.Now().Format("20060102") + "-" + randomString(6)
}

func randomString(n int) string {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}
```

### 10.4 Maintenance Service Implementation

```go
// application/maintenance/maintenance_service.go
package maintenance

import (
    "context"
    "time"

    "github.com/yourproject/internal/modules/erp/domain/entity"
    "github.com/yourproject/internal/modules/erp/domain/repository"
)

type MaintenanceService struct {
    workOrderRepo  repository.MaintenanceWorkOrderRepository
    scheduleRepo   repository.PreventiveMaintenanceRepository
    sparePartRepo  repository.SparePartRepository
    notificationSvc *notification.NotificationService
}

func NewMaintenanceService(
    workOrderRepo repository.MaintenanceWorkOrderRepository,
    scheduleRepo repository.PreventiveMaintenanceRepository,
    sparePartRepo repository.SparePartRepository,
    notificationSvc *notification.NotificationService,
) *MaintenanceService {
    return &MaintenanceService{
        workOrderRepo:  workOrderRepo,
        scheduleRepo:   scheduleRepo,
        sparePartRepo:  sparePartRepo,
        notificationSvc: notificationSvc,
    }
}

func (s *MaintenanceService) CreateWorkOrder(ctx context.Context, wo *entity.MaintenanceWorkOrder) error {
    wo.WorkOrderNo = generateWONo()
    wo.Status = "open"
    wo.CreatedAt = time.Now()

    if err := s.workOrderRepo.Create(ctx, wo); err != nil {
        return err
    }

    // Send notification
    s.notificationSvc.Send(ctx, &entity.Notification{
        UserID:   wo.AssignedTo,
        Type:     "maintenance",
        Severity: "warning",
        Title:    "งานซ่อมบำรุงใหม่",
        Message:  "มีงานซ่อมบำรุงหมายเลข " + wo.WorkOrderNo + " รอการดำเนินการ",
        Data: map[string]interface{}{
            "work_order_id": wo.ID,
            "machine_id":    wo.MachineID,
            "priority":      wo.Priority,
        },
    })

    return nil
}

func (s *MaintenanceService) GeneratePreventiveSchedule(ctx context.Context) error {
    schedules, err := s.scheduleRepo.GetDueSchedules(ctx)
    if err != nil {
        return err
    }

    for _, schedule := range schedules {
        wo := &entity.MaintenanceWorkOrder{
            MachineID:   schedule.MachineID,
            Type:        "preventive",
            Priority:    "medium",
            Description: "งานบำรุงรักษาเชิงป้องกันตามตาราง",
            StartDate:   time.Now(),
            Status:      "open",
        }
        if err := s.CreateWorkOrder(ctx, wo); err != nil {
            return err
        }

        // Update next run
        schedule.LastRun = time.Now()
        schedule.NextRun = calculateNextRun(schedule)
        if err := s.scheduleRepo.Update(ctx, schedule); err != nil {
            return err
        }
    }

    return nil
}

func calculateNextRun(schedule *entity.PreventiveMaintenanceSchedule) time.Time {
    now := time.Now()
    switch schedule.ScheduleType {
    case "daily":
        return now.AddDate(0, 0, schedule.Interval)
    case "weekly":
        return now.AddDate(0, 0, schedule.Interval*7)
    case "monthly":
        return now.AddDate(0, schedule.Interval, 0)
    case "yearly":
        return now.AddDate(schedule.Interval, 0, 0)
    default:
        return now.AddDate(0, 0, 1)
    }
}

func generateWONo() string {
    return "WO-" + time.Now().Format("20060102") + "-" + randomString(6)
}
```

### 10.5 Quality Service Implementation

```go
// application/quality/quality_service.go
package quality

import (
    "context"
    "time"

    "github.com/yourproject/internal/modules/erp/domain/entity"
    "github.com/yourproject/internal/modules/erp/domain/repository"
)

type QualityService struct {
    checkPointRepo repository.QualityCheckPointRepository
    resultRepo     repository.QualityResultRepository
    ncRepo         repository.NonConformanceRepository
}

func NewQualityService(
    checkPointRepo repository.QualityCheckPointRepository,
    resultRepo repository.QualityResultRepository,
    ncRepo repository.NonConformanceRepository,
) *QualityService {
    return &QualityService{
        checkPointRepo: checkPointRepo,
        resultRepo:     resultRepo,
        ncRepo:         ncRepo,
    }
}

func (s *QualityService) RecordQualityResult(ctx context.Context, result *entity.QualityResult) error {
    // Check if value is within specification
    checkPoint, err := s.checkPointRepo.GetByID(ctx, result.CheckPointID)
    if err != nil {
        return err
    }

    if result.Value < checkPoint.LowerLimit {
        result.Status = "fail"
    } else if result.Value > checkPoint.UpperLimit {
        result.Status = "fail"
    } else if result.Value < checkPoint.LowerLimit*1.1 || result.Value > checkPoint.UpperLimit*0.9 {
        result.Status = "warning"
    } else {
        result.Status = "pass"
    }

    result.CheckedAt = time.Now()
    result.CreatedAt = time.Now()

    if err := s.resultRepo.Create(ctx, result); err != nil {
        return err
    }

    // Create Non-Conformance if failed
    if result.Status == "fail" {
        nc := &entity.NonConformance{
            NCNo:        generateNCNo(),
            BatchNo:     result.BatchNo,
            ProductID:   result.ProductID,
            Severity:    "major",
            Description: "สินค้าไม่ผ่าน QC ที่จุดตรวจ " + checkPoint.Parameter,
            Status:      "open",
            ReportedBy:  result.InspectorID,
            CreatedAt:   time.Now(),
        }
        if err := s.ncRepo.Create(ctx, nc); err != nil {
            return err
        }
    }

    return nil
}

func (s *QualityService) GetSPCData(ctx context.Context, productID int64, parameter string, days int) ([]*entity.QualityResult, error) {
    return s.resultRepo.GetByProductAndParameter(ctx, productID, parameter, days)
}

func (s *QualityService) CalculateCPK(ctx context.Context, productID int64, parameter string) (float64, error) {
    // Get data
    results, err := s.resultRepo.GetByProductAndParameter(ctx, productID, parameter, 30)
    if err != nil {
        return 0, err
    }

    if len(results) == 0 {
        return 0, nil
    }

    // Calculate mean and standard deviation
    mean := 0.0
    for _, r := range results {
        mean += r.Value
    }
    mean /= float64(len(results))

    variance := 0.0
    for _, r := range results {
        variance += (r.Value - mean) * (r.Value - mean)
    }
    variance /= float64(len(results))
    stdDev := math.Sqrt(variance)

    // Get specification limits
    checkPoints, err := s.checkPointRepo.GetByProductAndParameter(ctx, productID, parameter)
    if err != nil || len(checkPoints) == 0 {
        return 0, nil
    }

    USL := checkPoints[0].UpperLimit
    LSL := checkPoints[0].LowerLimit

    cpu := (USL - mean) / (3 * stdDev)
    cpl := (mean - LSL) / (3 * stdDev)

    return math.Min(cpu, cpl), nil
}

func generateNCNo() string {
    return "NC-" + time.Now().Format("20060102") + "-" + randomString(6)
}
```

### 10.6 API Handlers

```go
// interfaces/http/erp_handler.go
package http

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/yourproject/internal/modules/erp/application/production"
    "github.com/yourproject/internal/modules/erp/application/maintenance"
    "github.com/yourproject/internal/modules/erp/application/quality"
)

type ERPHandler struct {
    productionSvc  *production.ProductionService
    maintenanceSvc *maintenance.MaintenanceService
    qualitySvc     *quality.QualityService
}

func NewERPHandler(
    productionSvc *production.ProductionService,
    maintenanceSvc *maintenance.MaintenanceService,
    qualitySvc *quality.QualityService,
) *ERPHandler {
    return &ERPHandler{
        productionSvc:  productionSvc,
        maintenanceSvc: maintenanceSvc,
        qualitySvc:     qualitySvc,
    }
}

// Production endpoints
func (h *ERPHandler) CreateProductionOrder(c *gin.Context) {
    var req struct {
        ProductID int64     `json:"product_id"`
        Quantity  float64   `json:"quantity"`
        StartDate time.Time `json:"start_date"`
        EndDate   time.Time `json:"end_date"`
        MachineID int64     `json:"machine_id"`
        OperatorID int64    `json:"operator_id"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    order := &entity.ProductionOrder{
        ProductID: req.ProductID,
        Quantity:  req.Quantity,
        StartDate: req.StartDate,
        EndDate:   req.EndDate,
        MachineID: req.MachineID,
        OperatorID: req.OperatorID,
    }

    if err := h.productionSvc.CreateProductionOrder(c.Request.Context(), order); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(201, order)
}

func (h *ERPHandler) GetProductionOrders(c *gin.Context) {
    tenant := middleware.GetTenant(c.Request.Context())
    orders, err := h.productionSvc.GetOrdersByTenant(c.Request.Context(), tenant.ID)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, orders)
}

func (h *ERPHandler) GetOEE(c *gin.Context) {
    machineID, _ := strconv.ParseInt(c.Query("machine_id"), 10, 64)
    dateStr := c.Query("date")

    date, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        date = time.Now()
    }

    oee, err := h.productionSvc.CalculateOEE(c.Request.Context(), machineID, date)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, oee)
}

// Maintenance endpoints
func (h *ERPHandler) CreateWorkOrder(c *gin.Context) {
    var req struct {
        MachineID   int64     `json:"machine_id"`
        Type        string    `json:"type"`
        Priority    string    `json:"priority"`
        Description string    `json:"description"`
        AssignedTo  int64     `json:"assigned_to"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    wo := &entity.MaintenanceWorkOrder{
        MachineID:   req.MachineID,
        Type:        req.Type,
        Priority:    req.Priority,
        Description: req.Description,
        AssignedTo:  req.AssignedTo,
    }

    if err := h.maintenanceSvc.CreateWorkOrder(c.Request.Context(), wo); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(201, wo)
}

func (h *ERPHandler) GetWorkOrders(c *gin.Context) {
    status := c.Query("status")
    machineID, _ := strconv.ParseInt(c.Query("machine_id"), 10, 64)

    orders, err := h.maintenanceSvc.GetWorkOrders(c.Request.Context(), status, machineID)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, orders)
}

// Quality endpoints
func (h *ERPHandler) RecordQualityResult(c *gin.Context) {
    var req struct {
        BatchNo      string  `json:"batch_no"`
        ProductID    int64   `json:"product_id"`
        CheckPointID int64   `json:"check_point_id"`
        Value        float64 `json:"value"`
        InspectorID  int64   `json:"inspector_id"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    result := &entity.QualityResult{
        BatchNo:      req.BatchNo,
        ProductID:    req.ProductID,
        CheckPointID: req.CheckPointID,
        Value:        req.Value,
        InspectorID:  req.InspectorID,
    }

    if err := h.qualitySvc.RecordQualityResult(c.Request.Context(), result); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(201, result)
}

func (h *ERPHandler) GetSPCData(c *gin.Context) {
    productID, _ := strconv.ParseInt(c.Query("product_id"), 10, 64)
    parameter := c.Query("parameter")
    days, _ := strconv.Atoi(c.Query("days"))
    if days == 0 {
        days = 30
    }

    data, err := h.qualitySvc.GetSPCData(c.Request.Context(), productID, parameter, days)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, data)
}

func (h *ERPHandler) CalculateCPK(c *gin.Context) {
    productID, _ := strconv.ParseInt(c.Query("product_id"), 10, 64)
    parameter := c.Query("parameter")

    cpk, err := h.qualitySvc.CalculateCPK(c.Request.Context(), productID, parameter)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, gin.H{"cpk": cpk})
}
```

### 10.7 Routes Setup

```go
// interfaces/http/routes.go
package http

import (
    "github.com/gin-gonic/gin"
    "github.com/yourproject/internal/infrastructure/middleware"
)

func SetupERPRoutes(r *gin.Engine, handler *ERPHandler, auth *middleware.AuthMiddleware) {
    api := r.Group("/api/v1/erp")
    api.Use(auth.Authenticate())
    api.Use(middleware.TenantMiddleware(tenantService))

    // Production routes
    production := api.Group("/production")
    {
        production.POST("/orders", auth.Authorize("operator"), handler.CreateProductionOrder)
        production.GET("/orders", auth.Authorize("viewer"), handler.GetProductionOrders)
        production.GET("/oee", auth.Authorize("viewer"), handler.GetOEE)
        production.POST("/oee/calculate", auth.Authorize("manager"), handler.CalculateOEE)
    }

    // Maintenance routes
    maintenance := api.Group("/maintenance")
    {
        maintenance.POST("/work-orders", auth.Authorize("operator"), handler.CreateWorkOrder)
        maintenance.GET("/work-orders", auth.Authorize("viewer"), handler.GetWorkOrders)
        maintenance.PUT("/work-orders/:id", auth.Authorize("manager"), handler.UpdateWorkOrder)
        maintenance.POST("/schedules/generate", auth.Authorize("manager"), handler.GenerateSchedules)
        maintenance.GET("/spare-parts", auth.Authorize("viewer"), handler.GetSpareParts)
    }

    // Quality routes
    quality := api.Group("/quality")
    {
        quality.POST("/results", auth.Authorize("operator"), handler.RecordQualityResult)
        quality.GET("/spc", auth.Authorize("viewer"), handler.GetSPCData)
        quality.GET("/cpk", auth.Authorize("viewer"), handler.CalculateCPK)
        quality.GET("/non-conformance", auth.Authorize("viewer"), handler.GetNonConformance)
        quality.PUT("/non-conformance/:id", auth.Authorize("manager"), handler.UpdateNonConformance)
    }

    // Inventory routes
    inventory := api.Group("/inventory")
    {
        inventory.GET("/items", auth.Authorize("viewer"), handler.GetInventoryItems)
        inventory.POST("/items", auth.Authorize("manager"), handler.CreateInventoryItem)
        inventory.PUT("/items/:id", auth.Authorize("manager"), handler.UpdateInventoryItem)
        inventory.GET("/transactions", auth.Authorize("viewer"), handler.GetInventoryTransactions)
        inventory.POST("/transactions", auth.Authorize("operator"), handler.CreateInventoryTransaction)
    }

    // Procurement routes
    procurement := api.Group("/procurement")
    {
        procurement.GET("/purchase-orders", auth.Authorize("viewer"), handler.GetPurchaseOrders)
        procurement.POST("/purchase-orders", auth.Authorize("manager"), handler.CreatePurchaseOrder)
        procurement.GET("/suppliers", auth.Authorize("viewer"), handler.GetSuppliers)
        procurement.POST("/suppliers", auth.Authorize("manager"), handler.CreateSupplier)
    }
}
```

---

## สรุป

ระบบ IQ-Factory Platform ที่ออกแบบใหม่นี้เป็นโซลูชันที่สมบูรณ์สำหรับการเปลี่ยนโรงงานให้เป็น Smart Factory ด้วย:

1. **IoT Platform** ที่ robust และ scalable
2. **ERP Modules** ที่ครอบคลุมทุกกระบวนการในโรงงาน
3. **AI/ML** ที่ช่วยเพิ่มประสิทธิภาพและลดต้นทุน
4. **Enterprise Features** ที่รองรับองค์กรทุกระดับ
5. **Business Model** ที่ชัดเจนและสามารถขายได้จริง