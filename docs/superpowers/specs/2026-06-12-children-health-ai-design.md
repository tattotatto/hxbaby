# 儿童大健康AI资料库系统 — 设计文档

**日期**: 2026-06-12
**状态**: 已确认
**版本**: v1.0

---

## 1. 项目概述

与母婴机构合作，构建儿童大健康AI系统。系统以API形式对外提供对话式AI健康服务，第三方应用（如微信小程序）通过API接入，为其用户提供专业的儿童健康咨询、症状评估和产品推荐。

### 1.1 核心链路

```
用户描述症状 → AI检索医学知识库 → 分析评估 → 生成建议 + 推荐租户产品/就医
```

### 1.2 关键决策

| 维度 | 决策 |
|------|------|
| 系统形态 | API优先的后端服务，对接微信小程序等第三方应用 |
| 用户端 | 家长端（数据采集+基础建议）+ 顾问端（深度分析） |
| 数据范围 | 全面健康档案 + 机构健康产品库 |
| 交互方式 | 对话式（Chat），SSE流式返回 |
| 部署环境 | 国内云部署，国内大模型 |
| 资料库 | 结构化文档 + 非结构化资料混合 |
| 知识库策略 | 共享知识库（医学知识）+ 租户产品库（产品推荐） |

---

## 2. 系统架构

### 2.1 架构模式

**Python AI + Go 业务 混合架构** — Python负责AI大脑（RAG、LLM），Go负责业务身体（API网关、用户管理、产品管理）。

### 2.2 架构全景

```
┌─────────────────────────────────────────────────┐
│  客户端层: 微信小程序 · H5 · 机构后台 · 第三方App    │
└────────────────────┬────────────────────────────┘
                     │ HTTPS / REST API
                     ▼
┌─────────────────────────────────────────────────┐
│  API 网关层 (Go - Gin)                            │
│  认证鉴权 · 限流 · 路由 · 日志 · 监控               │
└──────────┬──────────────────┬───────────────────┘
           │                  │
           ▼                  ▼
┌──────────────────┐  ┌──────────────────┐
│  业务服务 (Go)    │  │  AI 服务 (Python) │
│  · 用户管理       │  │  · RAG 检索增强   │
│  · 对话管理       │  │  · LLM 对话生成   │
│  · 产品管理       │  │  · 症状分析评估   │
│  · 推荐逻辑       │  │  · 知识库检索     │
│  · 客户租户隔离   │  │  · 建议生成       │
└──────┬───────────┘  └──────┬───────────┘
       │                     │
       └─────────┬───────────┘
                 │
    ┌────────────┼────────────┐
    ▼            ▼            ▼
┌───────┐  ┌────────┐  ┌────────┐  ┌───────┐
│PostgreSQL│ │Milvus  │  │ Redis  │  │  OSS  │
│业务数据  │ │向量检索 │  │缓存会话│  │文件存储│
└───────┘  └────────┘  └────────┘  └───────┘
                 │
                 ▼
┌─────────────────────────────────────────────────┐
│  LLM 模型层: 通义千问 · DeepSeek · 文心一言        │
│  (适配器模式，可随时切换)                           │
└─────────────────────────────────────────────────┘
```

### 2.3 核心设计原则

- **服务解耦**: Go业务服务和Python AI服务通过内部API通信，各自独立部署、独立扩缩容
- **多租户架构**: 每个合作机构是一个租户，拥有独立的产品库、用户池，数据完全隔离
- **模型可替换**: 通过适配器模式封装LLM调用，支持随时切换底层模型
- **可观测性**: 全链路日志追踪、Token用量监控、检索质量评估、用户反馈闭环

---

## 3. 知识库与RAG设计

### 3.1 数据隔离策略

| 层级 | 范围 | 内容 |
|------|------|------|
| 🌐 共享知识库 | 所有租户共用 | 儿童医学知识、生长标准、症状判断、处理指南、营养建议、发育里程碑 |
| 🏪 租户产品库 | 每个租户独立 | 机构健康产品、价格、适用症状、使用说明 |
| 👤 租户业务数据 | 完全隔离 | 用户档案、对话记录、评估历史、运营数据 |

### 3.2 RAG知识库管道

```
文档导入 (PDF/Word/Excel/Markdown/图片OCR)
    → 预处理清洗 (去噪/格式化/脱敏)
    → 智能分块 (语义分块, 512-1024 tokens, 10-20%重叠)
    → 向量化 (bge-large-zh-v1.5, 1024维) + 元数据标注 (分类/年龄段/权威等级)
    → 质量审核 (自动校验 + 人工审核 + 版本管理)
    → 入库 (Milvus + PostgreSQL元数据索引)
```

### 3.3 检索增强生成流程（问答时）

1. **查询重写**: 用户原始问题改写为多条检索查询
2. **多路召回**: 向量检索(Milvus) + 关键词检索 + 元数据过滤(年龄段/分类)
3. **重排序**: Cross-encoder模型精排，取Top-5最相关片段
4. **知识融合**: 共享知识库结果 + 租户产品库匹配 → 组合上下文
5. **LLM生成**: System Prompt注入专业身份 + 检索上下文 → 生成专业回答 + 产品推荐
6. **引用溯源**: 回答附带知识来源引用，顾问端可追溯原始资料

---

## 4. 数据模型

### 4.1 核心实体

**共享知识层**:
- `knowledge_docs`: 知识文档（id, title, category, source, authority_level, version, status）
- `knowledge_chunks`: 知识块（id, doc_id(FK), content, embedding(vector), chunk_index, tags, age_range, keywords）

**租户隔离层**:
- `tenants`: 租户/机构（id, name, config, status）
- `users`: 用户（id, tenant_id(FK), role(parent/consultant/admin), phone, wx_openid, name）
- `children`: 儿童档案（id, parent_id(FK), name, gender, birth_date, growth_records(JSON)）

**业务交互层**:
- `conversations`: 对话（id, tenant_id, user_id, child_id, status）
- `messages`: 消息（id, conversation_id(FK), role(user/assistant), content, retrieved_chunks(JSON), tokens_used）
- `assessments`: 评估记录（id, conversation_id(FK), child_id, symptoms(JSON), ai_analysis, risk_level, suggestion_type）

**租户产品层**:
- `product_categories`: 产品分类（id, tenant_id(FK), name, parent_id）
- `products`: 产品（id, tenant_id(FK), category_id, name, description, symptoms_tags, age_range, price, status）
- `recommendations`: 推荐记录（id, message_id(FK), product_id, score, reason, user_feedback）

---

## 5. API设计

### 5.1 API结构

| 分组 | 核心接口 |
|------|----------|
| 认证 & 租户 | `POST /api/v1/auth/login`, `POST /api/v1/auth/wx-login`, `GET /api/v1/tenant/profile` |
| 儿童档案 | `POST /api/v1/children`, `GET /api/v1/children`, `GET /api/v1/children/:id`, `POST /api/v1/children/:id/growth` |
| 对话 & AI（核心） | `POST /api/v1/conversations`, **`POST /api/v1/conversations/:id/chat`**, `GET /api/v1/conversations/:id/messages`, `POST /api/v1/conversations/:id/feedback` |
| 产品 & 推荐 | `GET /api/v1/products`, `GET /api/v1/products/:id`, `GET /api/v1/products/match?symptoms=xxx` |
| 评估报告 | `GET /api/v1/assessments/:id`, `GET /api/v1/assessments/child/:child_id` |
| 知识库管理 | `POST /api/v1/admin/knowledge/upload`, `GET /api/v1/admin/knowledge/docs`, `POST /api/v1/admin/knowledge/reindex` |

### 5.2 核心API：对话接口

```
POST /api/v1/conversations/:id/chat

Request:
{
  "message": "宝宝3个月，最近两天拉便便有奶瓣，怎么办？",
  "child_id": "c_123",
  "stream": true
}

Response (SSE Stream):
event: thinking → {"stage":"retrieving","message":"正在检索..."}
event: token    → {"content":"根据您描述的情况..."}
event: done     → {
  "message_id": "msg_456",
  "assessment": {"symptoms":[...], "analysis":"...", "risk_level":"low"},
  "recommendations": [{"product_id":"...", "name":"...", "reason":"...", "score":0.92}],
  "sources": [{"doc_title":"...", "chunk_id":"..."}]
}
```

---

## 6. 技术栈

### 6.1 AI服务 (Python)

| 组件 | 选型 | 理由 |
|------|------|------|
| Web框架 | FastAPI | 原生async，SSE流式支持，自动OpenAPI文档 |
| RAG框架 | LangChain + LangGraph | 最成熟RAG生态，Agent工作流支持 |
| LLM主模型 | 通义千问 Qwen-Plus | 中文医疗理解强，API稳定，约¥1.5/百万token |
| LLM备选 | DeepSeek-V3 | 性价比极高（¥1/百万token），中文能力强 |
| Embedding | bge-large-zh-v1.5 | 中文SOTA，1024维，本地部署零成本 |
| Reranker | bge-reranker-v2-m3 | 多语言Cross-encoder，显著提升检索精度 |
| 文档解析 | Unstructured + MinerU | 多格式支持，MinerU专注中文PDF表格识别 |

### 6.2 业务服务 (Go)

| 组件 | 选型 | 理由 |
|------|------|------|
| Web框架 | Gin | 国内Go后端首选，高性能，中间件丰富 |
| ORM | GORM | Go最流行ORM，多租户scope支持 |
| API文档 | Swagger/OpenAPI | Gin-swagger自动生成 |
| 认证 | JWT + 微信OAuth | 无状态认证，适合API分发 |

### 6.3 数据 & 基础设施

| 组件 | 选型 |
|------|------|
| 业务数据库 | PostgreSQL 15+ |
| 向量数据库 | Milvus 2.4+ |
| 缓存层 | Redis 7+ |
| 文件存储 | 阿里云OSS / MinIO(dev) |
| 消息队列 | Redis Streams (初期) / RabbitMQ |
| 容器化 | Docker Compose (dev) → 阿里云ACK K8s (prod) |
| 监控 | Prometheus + Grafana + 阿里云SLS |

---

## 7. 项目目录结构

```
hxbaby/
├── ai-service/                    # Python AI 服务
│   ├── api/                       # FastAPI 路由
│   │   ├── chat.py                # 核心对话接口(SSE)
│   │   ├── knowledge.py           # 知识库管理接口
│   │   └── health.py              # 健康评估接口
│   ├── core/                      # 核心逻辑
│   │   ├── rag_pipeline.py        # RAG检索管道
│   │   ├── llm_adapter.py         # LLM适配器(模型切换)
│   │   ├── symptom_analyzer.py    # 症状分析引擎
│   │   └── product_matcher.py     # 产品匹配推荐
│   ├── knowledge/                 # 知识库管道
│   │   ├── loader.py              # 文档加载器
│   │   ├── chunker.py             # 智能分块
│   │   ├── embedder.py            # 向量化
│   │   └── vector_store.py        # Milvus操作
│   ├── models/                    # Pydantic数据模型
│   ├── Dockerfile
│   └── requirements.txt
│
├── biz-service/                   # Go 业务服务
│   ├── cmd/main.go                # 入口
│   ├── internal/
│   │   ├── handler/               # HTTP handlers
│   │   ├── service/               # 业务逻辑
│   │   ├── repository/            # 数据访问(GORM)
│   │   ├── middleware/            # 认证/日志/租户
│   │   └── model/                 # 数据模型
│   ├── pkg/                       # 公共工具库
│   ├── Dockerfile
│   └── go.mod
│
├── docker-compose.yml             # 本地开发环境
├── docs/                          # 设计文档
│   └── superpowers/specs/
└── scripts/                       # 工具脚本
    └── init_knowledge.py          # 知识库初始化
```

---

## 8. MVP与路线图

### 8.1 MVP (Phase 1) — 4-6周

- 知识库管道：文档导入→分块→向量化→入库
- 基础对话：单轮问答 + RAG检索
- 症状分析：基于描述的初步判断
- 产品推荐：症状→产品匹配
- 多租户基础隔离
- 核心API：对话、儿童档案、产品查询
- 微信小程序基本对接

### 8.2 开发路径

| 周次 | 重点 |
|------|------|
| Week 1-2 | 知识库管道搭建 (文档解析→分块→向量化→入库) |
| Week 2-3 | AI对话核心 (RAG检索管道 + Chat API + SSE流式) |
| Week 3-4 | Go业务服务 (API网关 + 多租户 + 用户/儿童档案) |
| Week 4-5 | 产品推荐引擎对接 (症状匹配 + 产品库检索 + 推荐排序) |
| Week 5-6 | 联调测试 + 部署上线 |

### 8.3 后续阶段

- **Phase 2** (+4周): 多轮对话+记忆、生长曲线分析、顾问端面板、个性化建议、知识库持续更新、反馈闭环
- **Phase 3** (+4周): 智能预警、多模态识别、家长-顾问桥接、运营分析、知识库自动扩写
- **Phase 4+** (长期): Fine-tuning专属模型、多机构协同、视频问诊辅助、智能穿戴数据接入

---

## 9. 关键风险与对策

| 风险 | 对策 |
|------|------|
| 知识库质量不足导致AI判断不准 | 建立专家审核流程，A/B测试检索效果，持续迭代 |
| 国内模型在儿科医学领域知识有限 | RAG补偿模型知识盲区，长期考虑Fine-tuning |
| 多租户下的性能隔离 | 租户级限流，关键路径缓存，独立资源池 |
| 医疗合规风险（AI建议 ≠ 诊断） | 强制免责声明，设置风险分级（咨询/建议就医/紧急就医） |
| 资料格式多样解析困难 | MinerU处理中文PDF，人工兜底审核 |

---

## 10. 待决事项

以下事项在进入实现前需要确认：

- [ ] 机构提供的资料库首批内容和量级确认
- [ ] 通义千问API账号和预算审批
- [ ] 目标部署云平台确认（阿里云/华为云/腾讯云）
- [ ] 首个合作机构接入时间节点
- [ ] 微信小程序接入方技术对接人确认
