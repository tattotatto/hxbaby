# 后端开发工程师 — 任务手册

> 基于 [AI资料库系统设计](../specs/2026-06-12-children-health-ai-design.md) 和 [小程序工厂设计](../specs/2026-06-12-miniapp-factory-design.md)

---

## 职责范围

| 服务 | 语言 | 职责 |
|------|------|------|
| AI Service | Python | RAG管道、LLM适配、知识库、SSE对话API |
| Biz Service | Go | 业务API、认证鉴权、多租户、BaaS运行时 |
| CodeGen Service | Node.js | 代码生成引擎、模块编排、ZIP打包 |
| 基础设施 | - | PostgreSQL、Milvus、Redis、Docker、K8s |

---

## 一、AI资料库系统 (Phase 1, 4-6周)

### 1.1 Python AI 服务

#### 项目骨架
- **文件**: `ai-service/main.py`, `ai-service/config.py`, `ai-service/requirements.txt`
- 启动 FastAPI 应用，配置 CORS，挂载路由
- 环境变量管理 (`.env` → `Settings`)

#### 知识库管道
| 模块 | 文件 | 说明 |
|------|------|------|
| 文档加载器 | `ai-service/knowledge/loader.py` | 支持 PDF/DOCX/TXT/MD/CSV/XLSX |
| 分块器 | `ai-service/knowledge/chunker.py` | 递归语义分块，重叠控制 |
| Embedding | `ai-service/knowledge/embedder.py` | bge-large-zh-v1.5，单条+批量 |
| 向量存储 | `ai-service/knowledge/vector_store.py` | Milvus CRUD，索引创建，过滤检索 |

#### AI核心
| 模块 | 文件 | 说明 |
|------|------|------|
| LLM适配器 | `ai-service/core/llm_adapter.py` | 统一接口，支持通义千问/DeepSeek切换 |
| RAG管道 | `ai-service/core/rag_pipeline.py` | 检索→重排→上下文构建→产品匹配 |
| 症状分析 | `ai-service/core/symptom_analyzer.py` | LLM分析症状，输出风险等级+建议 |
| 对话API | `ai-service/api/chat.py` | `POST /ai/chat` SSE流式，完整RAG链路 |

#### 数据模型 (Pydantic)
- `ai-service/models/chat.py`: ChatRequest, ChatResponse, ProductRecommendation, SourceInfo

#### API 接口清单 (对内)
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/ai/health` | 健康检查 |
| POST | `/ai/chat` | 核心对话(SSE流式) |
| POST | `/ai/generate` | 文本生成(供AI Bridge调用) |
| POST | `/ai/admin/knowledge/upload` | 知识库文档上传 |
| POST | `/ai/admin/knowledge/reindex` | 重建索引 |

---

### 1.2 Go 业务服务

#### 项目骨架
- `biz-service/cmd/main.go`, `biz-service/internal/config/config.go`
- `biz-service/internal/router/router.go`

#### 数据模型 (GORM)
| 模型 | 文件 |
|------|------|
| Tenant | `internal/model/tenant.go` |
| User | `internal/model/user.go` |
| Child (含JSON生长记录) | `internal/model/child.go` |
| Conversation | `internal/model/conversation.go` |
| Message | `internal/model/message.go` |
| Product + ProductCategory | `internal/model/product.go` |
| Assessment + Recommendation | `internal/model/assessment.go` |

#### 业务层 (Handler → Service → Repository)
| 模块 | 文件 | 核心接口 |
|------|------|----------|
| 认证 | `handler/auth.go`, `service/auth.go`, `repository/user.go` | login, wx-login, register |
| 租户 | `handler/tenant.go` | 租户配置读取 |
| 儿童 | `handler/child.go` | CRUD + 生长记录 |
| 对话 | `handler/conversation.go` | 对话CRUD |
| Chat代理 | `handler/chat.go`, `service/ai_client.go` | SSE透传 → Python AI |
| 产品 | `handler/product.go` | 产品CRUD + 症状匹配 |

#### 中间件
| 中间件 | 文件 | 功能 |
|--------|------|------|
| AuthRequired | `middleware/auth.go` | JWT Bearer Token 校验 |
| TenantIsolation | `middleware/tenant.go` | 租户数据隔离 |
| RateLimit | `middleware/ratelimit.go` | 基于Redis的限流 |
| Logger | `middleware/logger.go` | 请求日志 |

#### API 接口清单 (对外)
| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| POST | `/api/v1/auth/login` | 无 | 账号密码登录 |
| POST | `/api/v1/auth/wx-login` | 无 | 微信code登录 |
| GET | `/api/v1/tenant/profile` | JWT | 租户信息 |
| POST | `/api/v1/children` | JWT | 创建儿童档案 |
| GET | `/api/v1/children` | JWT | 儿童列表 |
| GET | `/api/v1/children/:id` | JWT | 儿童详情 |
| POST | `/api/v1/children/:id/growth` | JWT | 添加生长记录 |
| POST | `/api/v1/conversations` | JWT | 创建对话 |
| **POST** | **`/api/v1/conversations/:id/chat`** | JWT | **核心：AI对话(SSE)** |
| GET | `/api/v1/conversations/:id/messages` | JWT | 对话历史 |
| POST | `/api/v1/conversations/:id/feedback` | JWT | 回复反馈 |
| GET | `/api/v1/products` | JWT | 产品列表 |
| GET | `/api/v1/products/match` | JWT | 症状匹配产品 |
| GET | `/api/v1/assessments/:id` | JWT | 评估详情 |
| GET | `/api/v1/assessments/child/:child_id` | JWT | 儿童评估历史 |
| POST | `/api/v1/admin/knowledge/upload` | JWT+Admin | 上传知识文档 |
| POST | `/api/v1/admin/knowledge/reindex` | JWT+Admin | 重建索引 |

---

## 二、小程序工厂平台 (Phase 2, 6-8周)

### 2.1 Go 扩展 — 工厂管理

#### 新增数据模型
| 模型 | 文件 |
|------|------|
| Customer | `internal/model/customer.go` |
| MiniappProject | `internal/model/miniapp_project.go` |
| BuildTask | `internal/model/build_task.go` |

#### 新增业务层
| 模块 | 核心接口 |
|------|----------|
| CustomerHandler | `POST /api/v1/auth/register`, `POST /api/v1/auth/login` (客户版) |
| ProjectHandler | CRUD + `POST /api/v1/projects/:id/build` (触发构建) |

#### API Key 工具
- `pkg/apikey/apikey.go`: 生成/校验 API Key + Secret

---

### 2.2 Go 扩展 — BaaS 运行时

#### 新增数据模型
| 模型 | 文件 | 用途 |
|------|------|------|
| CmsArticle | `internal/model/cms_article.go` | 文章(含AI生成标记) |
| BActivity | `internal/model/b_activity.go` | 活动(含AI文案/海报) |
| BProduct | `internal/model/b_product.go` | 商品(含AI卖点) |
| BOrder | `internal/model/b_order.go` | 订单 |
| BUser | `internal/model/b_user.go` | 终端用户 |
| BBooking | `internal/model/b_booking.go` | 预约 |
| BMember | `internal/model/b_member.go` | 会员 |

#### 新增中间件
| 中间件 | 功能 |
|--------|------|
| BaaSAuth | X-API-Key 校验 + 项目识别 |
| ProjectIsolate | project_id 强制注入，数据隔离 |

#### BaaS API 接口清单（供小程序调用）
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/baas/v1/auth/wx-login` | 终端用户微信登录 |
| GET | `/baas/v1/articles` | 文章列表 |
| GET | `/baas/v1/articles/:id` | 文章详情 |
| POST | `/baas/v1/ai/chat` | AI对话(→AI系统) |
| POST | `/baas/v1/ai/assess` | 症状评估(→AI系统) |
| GET | `/baas/v1/products` | 商品列表 |
| POST | `/baas/v1/orders` | 创建订单 |
| POST | `/baas/v1/payment/wechat` | 微信支付 |
| GET | `/baas/v1/activities` | 活动列表 |
| POST | `/baas/v1/activities/:id/signup` | 活动报名 |
| POST | `/baas/v1/activities/:id/checkin` | 签到 |
| GET | `/baas/v1/user/profile` | 用户信息 |
| PUT | `/baas/v1/user/profile` | 更新用户信息 |
| POST | `/baas/v1/user/children` | 添加儿童 |
| GET | `/baas/v1/bookings/slots` | 可预约时段 |
| POST | `/baas/v1/bookings` | 创建预约 |

#### AI Bridge (统一AI调用)
- `internal/service/ai_bridge.go`: 封装对Python AI服务的所有调用
- `internal/handler/ai_enhance.go`: 提供AI增强API(写文章/摘要/活动文案/卖点/复盘)

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/ai/generate-article` | AI写文章 |
| POST | `/api/v1/ai/generate-summary` | AI写摘要 |
| POST | `/api/v1/ai/generate-activity-copy` | AI活动文案 |
| POST | `/api/v1/ai/generate-poster` | AI海报生成 |
| POST | `/api/v1/ai/generate-selling-points` | AI卖点提炼 |
| POST | `/api/v1/ai/activity-report` | AI活动复盘 |
| POST | `/api/v1/ai/match-images` | AI配图推荐 |
| POST | `/api/v1/ai/suggest-tags` | AI标签推荐 |

---

### 2.3 Node.js 代码生成引擎

#### 项目骨架
- `codegen-service/package.json`, `codegen-service/server.js`
- `codegen-service/scripts/`

#### 核心模块
| 模块 | 文件 | 功能 |
|------|------|------|
| 依赖解析 | `scripts/resolver.js` | 模块依赖图+拓扑排序+冲突检测 |
| 配置注入 | `scripts/injector.js` | 模板变量替换(品牌/API端点/租户ID) |
| 编排引擎 | `scripts/compose.js` | 复制骨架→合并路由→拼装模块→生成README |
| 打包器 | `scripts/packager.js` | ZIP压缩+MD5校验 |
| 校验器 | `scripts/validator.js` | 生成后合法性检查 |

#### API 接口
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/health` | 健康检查 |
| GET | `/api/modules` | 获取可用模块列表 |
| POST | `/api/build` | 触发代码生成(接收project对象，返回zip包路径) |

#### 模块声明规范 (module.json)
```json
{
  "name": "cms",
  "displayName": "内容管理",
  "version": "1.0.0",
  "dependencies": ["base"],
  "pages": [{ "path": "pages/article/list", "style": {} }],
  "tabbar": { "pagePath": "pages/article/list", "text": "资讯", "iconPath": "...", "selectedIconPath": "..." },
  "api": ["content.js"],
  "components": [],
  "configSchema": {}
}
```

---

## 三、前后端接口契约

### Go Biz Service → Python AI Service (内部)

```
POST {AI_SERVICE_URL}/ai/chat
Content-Type: application/json
{
  "message": "string (用户问题)",
  "conversation_id": "string (可选)",
  "child_id": "string (可选)",
  "child_age_months": "int (可选)",
  "tenant_id": "string",
  "product_chunks": [{"product_id":"...", "name":"...", "embedding":[...]}],
  "history": [{"role":"user/assistant", "content":"..."}]
}

→ SSE Stream:
event: thinking | data: {"stage":"retrieving/analyzing", "message":"..."}
event: token    | data: {"content":"..."}
event: done     | data: {message_id, content, assessment, recommendations, sources}
```

### Go Biz Service → Node.js CodeGen Service (内部)

```
POST {CODEGEN_URL}/api/build
Content-Type: application/json
{
  "project": {
    "id": 123,
    "name": "宝宝健康助手",
    "modules": ["base","cms","ai-advisor","activity"],
    "brand_config": {"appName":"...","primaryColor":"#4caf50","logo":"..."},
    "api_key": "ak_xxx",
    "wx_app_id": "wxXXX"
  }
}

→ Response:
{
  "task_id": "uuid",
  "status": "done",
  "output_dir": "/path/to/output",
  "zip_path": "/path/to/output.zip",
  "md5": "abc123",
  "size_bytes": 2048000,
  "warnings": ["电商模块需要配置微信支付..."]
}
```

### 小程序 → BaaS API (外网)

```
Header: X-API-Key: ak_xxx
所有请求携带 X-API-Key，Go中间件自动识别project_id并隔离数据
```

---

## 四、开发顺序

```
Week 1-2: Python AI服务 → 知识库管道(loader→chunker→embedder→Milvus)
Week 2-3: Python AI服务 → LLM适配器 → RAG管道 → SSE Chat API
Week 3-4: Go业务服务 → 数据模型 → JWT认证 → 租户中间件 → Chat代理
Week 4-5: Go扩展 → 产品推荐 → 知识库管理API
Week 5-6: Go扩展 → 工厂管理(Customer/Project/Build) → BaaS认证 → 内容/活动API
Week 6-7: Node.js → 代码生成引擎(resolver→injector→compose→packager)
Week 7-8: Go扩展 → BaaS电商/支付/用户API → AI Bridge → AI增强API
Week 8-10: 联调测试 → 集成测试 → Docker部署 → 上线
```

## 五、环境依赖

```
开发环境: Docker Compose (PostgreSQL + Milvus + Redis + MinIO)
Python: 3.11 + FastAPI + LangChain + sentence-transformers + pymilvus
Go: 1.22 + Gin + GORM + golang-jwt
Node.js: 20 + Express + Handlebars + archiver
生产环境: 阿里云ACK(K8s) + SLB + SLS
```
