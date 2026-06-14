# 架构师执行计划 — 小块任务看板

> **执行模式**：架构师分配 → 开发者实现 → 架构师验证 → 下一块
> **规则**：每块1-2天完成，验证通过才能继续。后端和前端可并行时标注 🟢。

---

## Sprint 0: 基础设施搭建 (2天)

### Chunk 0.1 — 🔧 后端：Docker开发环境 [1天]

**分配**: 后端开发工程师

**任务**:
1. 编写 `docker-compose.yml`（PostgreSQL 16 + Milvus 2.4 + Redis 7 + MinIO + Etcd）
2. 编写 `.env.example` 和 `.gitignore`
3. `docker compose up -d` 启动所有服务
4. 验证：`docker compose ps` 全部 running

**验收标准**:
- [ ] `docker compose up -d` 一键启动所有依赖服务
- [ ] PostgreSQL 端口 5432 可连接
- [ ] Milvus 端口 19530 可连接
- [ ] Redis 端口 6379 可连接

---

### Chunk 0.2 — 🔧 后端：项目骨架 [1天]

**分配**: 后端开发工程师

**任务**:
1. Python AI 服务骨架：`ai-service/` 目录结构、`main.py`、`config.py`、`requirements.txt`、`Dockerfile`
2. Go 业务服务骨架：`biz-service/` 目录结构、`cmd/main.go`、`config.go`、`router.go`、`Dockerfile`
3. 两个服务各实现 `/health` 端点
4. 验证：`curl localhost:8001/ai/health` 和 `curl localhost:8080/health` 都返回 `{"status":"ok"}`

**验收标准**:
- [ ] Python 服务启动，`/ai/health` 返回 ok
- [ ] Go 服务启动，`/health` 返回 ok
- [ ] 两个服务都能在 Docker 中构建和运行

---

## Sprint 1: AI知识库管道 (Week 1-2)

### Chunk 1.1 — 🔧 后端：文档加载器 [1天]

**分配**: 后端开发工程师

**任务**: 实现 `ai-service/knowledge/loader.py`
- 支持 PDF / DOCX / TXT / Markdown / CSV / Excel
- 统一返回 `[{content, metadata}]`
- 不支持格式抛异常
- 文件不存在抛异常

**验收标准**:
- [ ] 加载 `.txt` 文件返回正确内容
- [ ] 加载 `.csv` 文件返回表格markdown
- [ ] 加载不存在的文件抛出 `DocumentLoadError`
- [ ] 加载 `.xyz` 抛出 "Unsupported format"
- [ ] 单元测试全部 PASS

---

### Chunk 1.2 — 🔧 后端：文档分块器 [1天]

**分配**: 后端开发工程师

**任务**: 实现 `ai-service/knowledge/chunker.py`
- 递归语义分块（按段落→句子→字符优先级分割）
- 支持 `chunk_size` 和 `overlap` 参数
- 返回 `[{content, chunk_index, metadata}]`
- 短文本（< chunk_size）返回单个chunk
- 每个chunk不超过 chunk_size

**验收标准**:
- [ ] `chunk_document("hello", 800, 120)` 返回 1 个chunk
- [ ] 2000字符文本 + chunk_size=500 → 生成 4+ 个chunk
- [ ] overlap 验证：前一个chunk末尾出现在后一个chunk开头
- [ ] metadata 正确传递到每个chunk
- [ ] 单元测试全部 PASS

---

### Chunk 1.3 — 🔧 后端：Embedding服务 [1天]

**分配**: 后端开发工程师

**任务**: 实现 `ai-service/knowledge/embedder.py`
- 加载 bge-large-zh-v1.5 模型（首次自动下载）
- `embed(text)` 返回 1024维向量
- `embed_batch(texts)` 批量向量化
- 单例模式，全局复用
- 相似文本余弦相似度 > 不相似文本

**验收标准**:
- [ ] `embed("测试")` 返回 1024 个 float
- [ ] `embed_batch(["a", "b", "c"])` 返回 3×1024
- [ ] "宝宝消化不好" vs "婴儿消化不良" 相似度 > "今天天气很好"
- [ ] 单元测试全部 PASS

---

### Chunk 1.4 — 🔧 后端：Milvus向量存储 [1天]

**分配**: 后端开发工程师

**任务**: 实现 `ai-service/knowledge/vector_store.py`
- 连接 Milvus，自动创建 Collection（含索引）
- `insert(chunks, embeddings)` 批量插入
- `search(query_vector, top_k, filter_expr)` 向量检索
- 支持元数据过滤（如 `age_range == "0-1岁"`）

**验收标准**:
- [ ] 插入 2 条数据后检索返回正确结果
- [ ] `filter_expr` 过滤生效
- [ ] Collection 支持 drop 和重建
- [ ] 单元测试全部 PASS（需要 Milvus 运行）

---

### Chunk 1.5 — 🔧 后端：知识库完整管道集成 [1天]

**分配**: 后端开发工程师

**任务**: 编写 `scripts/init_knowledge.py`
- 接收文档目录路径
- 遍历所有文件 → load → chunk → embed → insert into Milvus
- 打印进度和统计信息

**验收标准**:
- [ ] 将一个包含 3 个 PDF/TXT 文件的目录导入知识库
- [ ] Milvus 中可检索到对应内容
- [ ] 控制台输出处理统计（文件数、chunk数、耗时）

> 🏗️ **架构师验证点 #1**：知识库管道端到端可用。准备进入AI对话核心。

---

## Sprint 2: AI对话核心 (Week 2-3)

### Chunk 2.1 — 🔧 后端：LLM适配器 [1.5天]

**分配**: 后端开发工程师

**任务**: 实现 `ai-service/core/llm_adapter.py`
- 封装通义千问 API（通过 LangChain ChatOpenAI 兼容接口）
- `generate(messages)` 同步生成
- `generate_stream(messages)` 流式生成
- 支持模型切换（qwen-plus / deepseek-chat）
- `get_available_models()` 返回可用模型列表

**验收标准**:
- [ ] `generate([{"role":"user","content":"你好"}])` 返回非空字符串
- [ ] `generate_stream(...)` 逐 token yield
- [ ] 切换到 deepseek-chat 后正常工作
- [ ] 单元测试全部 PASS（mock LLM 调用）

---

### Chunk 2.2 — 🔧 后端：RAG检索管道 [1.5天]

**分配**: 后端开发工程师

**任务**: 实现 `ai-service/core/rag_pipeline.py`
- `rewrite_query(query)` 查询重写（LLM优化检索词）
- `retrieve(query, top_k, age_months, category)` 向量检索 + 元数据过滤
- `rerank(query, chunks, top_k)` LLM重排序
- `build_context(chunks)` 组合检索结果为LLM上下文
- `search_product_knowledge(query, product_chunks)` 产品库匹配

**验收标准**:
- [ ] `retrieve("宝宝发烧", age_months=6)` 返回相关结果且包含年龄过滤
- [ ] `rewrite_query("娃拉肚子")` 改写为更专业的检索词
- [ ] `rerank` 将最相关结果排到前面
- [ ] `build_context` 生成带来源标注的上下文
- [ ] 单元测试全部 PASS

---

### Chunk 2.3 — 🔧 后端：SSE对话API [2天]

**分配**: 后端开发工程师

**任务**: 实现 `ai-service/api/chat.py`
- `POST /ai/chat` 端点，接收 `ChatRequest`
- SSE流式响应：`event: thinking` → `event: token`... → `event: done`
- System Prompt 注入（儿科专业身份+安全规则）
- done 事件包含：message_id, content, assessment, recommendations, sources
- 实现 `symptom_analyzer.py`：症状分析→风险分级→建议类型
- 在 `main.py` 中注册路由

**验收标准**:
- [ ] `curl -N POST /ai/chat` 返回 SSE 流
- [ ] 流中包含 thinking → token → done 完整事件序列
- [ ] done 事件 JSON 包含 assessment（症状分析结果）
- [ ] done 事件 JSON 包含 sources（知识来源引用）
- [ ] 集成测试全部 PASS

> 🏗️ **架构师验证点 #2**：AI对话核心链路打通（query → retrieve → rerank → context → LLM → SSE）。准备进入Go业务层。

---

## Sprint 3: Go业务服务 + 多租户 (Week 3-4)

### Chunk 3.1 — 🔧 后端：数据库模型与迁移 [1天]

**分配**: 后端开发工程师

**任务**:
1. 实现 8 个 GORM Model（Tenant, User, Child, Conversation, Message, Product, ProductCategory, Assessment, Recommendation）
2. `config/database.go` 数据库连接 + AutoMigrate
3. `cmd/main.go` 启动时自动建表
4. 验证：数据库中出现对应表

**验收标准**:
- [ ] 启动 Go 服务后 PostgreSQL 中自动创建 8+ 张表
- [ ] Child.GrowthRecords JSON 字段可正常存取
- [ ] 表结构符合设计文档数据模型

---

### Chunk 3.2 — 🔧 后端：JWT认证 + 租户中间件 [1.5天]

**分配**: 后端开发工程师

**任务**:
1. `pkg/jwt/jwt.go`：生成/解析 JWT Token（含 user_id, tenant_id, role）
2. `middleware/auth.go`：Bearer Token 校验中间件
3. `middleware/tenant.go`：租户数据隔离中间件
4. `handler/auth.go` + `service/auth.go` + `repository/user.go`：登录/注册 API
5. 路由注册

**验收标准**:
- [ ] `POST /api/v1/auth/login` 返回 JWT Token
- [ ] 无 Token 访问受保护接口返回 401
- [ ] Token 中 tenant_id 注入到请求上下文
- [ ] 不同租户用户互相看不到对方数据

---

### Chunk 3.3 — 🔧 后端：Chat代理端点 [1.5天]

**分配**: 后端开发工程师

**任务**:
1. `service/ai_client.go`：HTTP Client 调用 Python AI 服务 `/ai/chat`
2. `handler/chat.go`：接收请求 → 组装payload（含历史+儿童月龄+产品chunks）→ 透传SSE到客户端
3. `POST /api/v1/conversations/:id/chat` 端点
4. Go → Python 内部调用（不经过公网）

**验收标准**:
- [ ] `curl -N POST /api/v1/conversations/1/chat` 返回 SSE 流
- [ ] 服务间调用延迟 < 50ms（本地）
- [ ] Go 正确透传 Python AI 服务的 SSE 事件

> 🏗️ **架构师验证点 #3**：Go → Python 全链路打通。AI对话从Go入口到Python AI再返回客户端。

---

### Chunk 3.4 — 🔧 后端：儿童档案 + 产品管理 CRUD [1.5天]

**分配**: 后端开发工程师

**任务**:
1. Child Handler+Service+Repository（CRUD + 生长记录管理）
2. Product Handler+Service+Repository（CRUD + 症状匹配查询）
3. Conversation Handler+Service+Repository（对话历史管理）

**验收标准**:
- [ ] `POST/GET /api/v1/children` 正常增删改查
- [ ] `POST /api/v1/children/:id/growth` 添加生长记录
- [ ] `GET /api/v1/products/match?symptoms=腹泻` 返回匹配产品
- [ ] 对话历史正确按 conversation_id 分组

> 🏗️ **架构师验证点 #4**：AI资料库系统 MVP 后端完成。可开始小程序工厂后端。

---

## Sprint 4: 工厂管理 + BaaS (Week 4-6)

### Chunk 4.1 — 🔧 后端：客户+项目+构建模型 [1天]

**分配**: 后端开发工程师

**任务**:
1. Customer, MiniappProject, BuildTask 三个Model
2. `pkg/apikey/apikey.go` API Key/Secret 生成工具
3. AutoMigrate 更新

**验收标准**:
- [ ] 数据库新建 3 张表
- [ ] `apikey.GenerateKey()` 返回 `ak_` 开头的 64 位随机字符串
- [ ] `apikey.GenerateSecret()` 返回 64 位 hex

---

### Chunk 4.2 — 🔧 后端：客户注册登录+项目管理API [1.5天]

**分配**: 后端开发工程师

**任务**:
1. Customer Handler+Service+Repository（注册/登录/套餐查询）
2. Project Handler+Service+Repository（CRUD + 项目数限制检查 + API Key自动生成）
3. 路由注册

**验收标准**:
- [ ] 客户注册 → 自动分配免费套餐（max_projects=1）
- [ ] 创建项目 → 自动生成 api_key + api_secret
- [ ] 第2个项目创建被拒绝（免费版限制）
- [ ] 创建项目时 modules 至少包含 "base"

> 🟢 **可并行**: Chunk 4.2 完成后，前端可开始 Chunk 5.1

---

### Chunk 4.3 — 🎨 前端：管理后台骨架+登录注册 [2天]

**分配**: 前端设计师/开发工程师

**任务**:
1. Vite + React 18 + Ant Design 5 项目初始化
2. `api/client.ts` axios 封装（Token 拦截器）
3. `pages/Login.tsx` 手机号+密码登录
4. `pages/Register.tsx` 注册页
5. `pages/Dashboard.tsx` 项目列表（空状态+项目卡片）
6. `components/Layout.tsx` 管理后台布局

**验收标准**:
- [ ] 登录页 UI 完整，输入校验正常
- [ ] 登录成功跳转 Dashboard
- [ ] Dashboard 展示项目列表（空/有数据两种状态）
- [ ] Token 过期自动跳转登录页
- [ ] 页面适配 PC 端（1280px+）

---

### Chunk 4.4 — 🎨 前端：创建向导（3步） [2.5天]

**分配**: 前端设计师/开发工程师

**任务**:
1. `pages/ProjectCreate.tsx` 基本信息（名称、描述、微信AppID）
2. `components/ModuleSelector.tsx` 模块勾选卡片（8模块、依赖提示、AI标签）
3. `pages/ProjectConfig.tsx` 整合模块选择器
4. `components/BrandConfig.tsx` 品牌配置（Logo上传、色板、手机预览）
5. `pages/ProjectBrand.tsx` 品牌配置页
6. 步骤导航（Steps 组件，支持前进/后退）

**验收标准**:
- [ ] 3步向导流程顺畅，数据在步骤间保留
- [ ] 模块卡片展示正确（图标、名称、AI标识、依赖说明）
- [ ] 选电商模块时提示"需要微信支付配置"
- [ ] 基础模块不可取消
- [ ] Logo 上传 + 预览正常
- [ ] 主题色选择实时反映到手机预览

---

### Chunk 4.5 — 🔧 后端：BaaS认证+内容API [1.5天]

**分配**: 后端开发工程师

**任务**:
1. `middleware/baas_auth.go` X-API-Key 校验中间件
2. `middleware/project_isolate.go` 项目数据隔离
3. CmsArticle Model + Repository
4. BaaS Content Handler（文章列表/详情，已发布过滤，阅读量统计）
5. BaaS 路由组注册

**验收标准**:
- [ ] 无 X-API-Key 返回 401
- [ ] 无效 X-API-Key 返回 401
- [ ] `GET /baas/v1/articles` 只返回当前项目的文章
- [ ] 文章详情请求自动 +1 阅读量

> 🟢 **可并行**: Chunk 4.5 完成后，前端可开始内容管理页面

---

### Chunk 4.6 — 🔧 后端：Node.js代码生成引擎 [3天] ⚡关键模块

**分配**: 后端开发工程师

**任务**:
1. `codegen-service/` 项目初始化（Express + Handlebars + archiver）
2. `scripts/resolver.js` 模块依赖解析 + 拓扑排序 + 冲突检测
3. `scripts/injector.js` 配置注入（变量替换）
4. `scripts/compose.js` 编排引擎（骨架+路由合并+模块拼装+npm配置+README）
5. `scripts/packager.js` ZIP打包 + MD5
6. `scripts/validator.js` 生成后合法性校验
7. `server.js` Express API：
   - `GET /api/modules` 获取可用模块列表
   - `POST /api/build` 触发代码生成

**验收标准**:
- [ ] `POST /api/build` 输入 project JSON → 返回 ZIP 包路径 + MD5
- [ ] 选择 ["base","cms","ai-advisor"] → ZIP 包含 3 个模块的页面文件
- [ ] pages.json 路由正确合并
- [ ] `{{appName}}` 等变量被正确替换
- [ ] 依赖自动补齐（选 shop → 自动加 base）
- [ ] 拓扑排序保证 base 在 shop 前面
- [ ] 生成的 README.md 包含接入指引
- [ ] ZIP 可正常解压，结构完整

> 🏗️ **架构师验证点 #5**：代码生成引擎核心完成。输入项目配置 → 输出可用的 UniApp 工程。

---

### Chunk 4.7 — 🔧 后端：Go触发构建+OSS上传 [1.5天]

**分配**: 后端开发工程师

**任务**:
1. BuildTask Repository
2. `POST /api/v1/projects/:id/build` Handler：
   - 创建 BuildTask（status: pending）
   - 调用 Node.js 代码生成 API
   - 轮询/等待结果
   - 上传 ZIP 到 OSS
   - 更新 BuildTask（status: done, zip_url, md5, duration_ms）
3. `GET /api/v1/builds/:id/status` 查询构建状态
4. `GET /api/v1/builds/:id/download` OSS 重定向下载

**验收标准**:
- [ ] `POST /projects/1/build` → 返回 build_id
- [ ] 构建完成后 status 变为 done
- [ ] `GET /builds/:id/download` 返回 ZIP 文件
- [ ] 下载的 ZIP MD5 与 BuildTask 记录一致

> 🟢 **可并行**: Chunk 4.7 完成后，前端可开始构建下载页面

---

### Chunk 4.8 — 🎨 前端：构建下载+接入指引 [1.5天]

**分配**: 前端设计师/开发工程师

**任务**:
1. `components/BuildProgress.tsx` 6阶段构建进度动画
2. `pages/ProjectBuild.tsx` 构建触发→进度→下载→历史
3. `pages/ProjectGuide.tsx` 4步接入指引（下载→安装→配置→发布）
4. API Key 复制、BaaS 文档链接

**验收标准**:
- [ ] 点击"构建" → 进度条动画展示 6 个阶段
- [ ] 构建完成 → 下载按钮可用
- [ ] 接入指引页面展示 4 步流程，内容完整
- [ ] API Key 一键复制功能正常

---

## Sprint 5: AI增强 + BaaS电商/活动 (Week 6-7)

### Chunk 5.1 — 🔧 后端：AI Bridge + AI增强API [2天]

**分配**: 后端开发工程师

**任务**:
1. `service/ai_bridge.go` 封装所有AI调用：
   - `GenerateArticle(topic, category)` AI写文章
   - `GenerateSummary(article)` AI写摘要
   - `GenerateActivityCopy(title, desc)` AI活动文案
   - `GenerateSellingPoints(name, desc)` AI卖点提炼
   - `GenerateActivityReport(name, stats)` AI复盘
2. `handler/ai_enhance.go` AI增强API端点
3. Python AI 服务添加 `/ai/generate` 通用文本生成端点

**验收标准**:
- [ ] `POST /api/v1/ai/generate-article` 返回 AI 生成的文章
- [ ] `POST /api/v1/ai/generate-activity-copy` 返回文案（含标题+详情+推文+朋友圈）
- [ ] `POST /api/v1/ai/generate-selling-points` 返回 3-5 个卖点
- [ ] AI Bridge 调用超时（30s）不阻塞，返回错误信息

> 🟢 **可并行**: Chunk 5.1 完成后，前端可开始 AI 文案面板开发

---

### Chunk 5.2 — 🎨 前端：AI文案面板 + 内容管理后台 [2天]

**分配**: 前端设计师/开发工程师

**任务**:
1. `components/AICopyPanel.tsx` 通用 AI 文案生成面板
   - 输入区（主题/关键词/分类）
   - 生成按钮 + Loading
   - 结果区（Tab：文章/摘要/推文/朋友圈）
   - 一键复制、一键应用到编辑器
2. `pages/ContentManager.tsx` 文章管理（表格+新建+编辑）
   - 新建时嵌入 AI 写作助手
   - 文章列表展示 AI 生成标记

**验收标准**:
- [ ] AI 文案面板输入主题 → 生成 → 展示结果
- [ ] 生成结果可切换 Tab 查看
- [ ] 一键复制到剪贴板
- [ ] 内容管理页可新建/编辑/发布文章

---

### Chunk 5.3 — 🎨 前端：活动管理后台 [1.5天]

**分配**: 前端设计师/开发工程师

**任务**:
1. `pages/ActivityManager.tsx` 活动列表+新建/编辑
   - 活动基本信息表单
   - 嵌入 AI 文案面板（生成活动标题+详情+推文）
   - AI 海报生成入口
2. AI 海报预览组件

**验收标准**:
- [ ] 新建活动时可一键生成文案
- [ ] 活动列表展示报名数、状态
- [ ] 海报生成入口可用

---

### Chunk 5.4 — 🔧 后端：BaaS电商+活动+用户API [2天]

**分配**: 后端开发工程师

**任务**:
1. BProduct + BOrder Model + Repository + Handler
2. BActivity + BActivitySignup Model + Repository + Handler
3. BUser Model + Repository + Handler（wx-login, profile）
4. BBooking Model + Repository + Handler
5. 微信支付集成（`service/wechat_pay.go`）：
   - 统一下单
   - 支付回调处理

**验收标准**:
- [ ] BaaS 电商接口完整（商品列表/详情/下单）
- [ ] BaaS 活动接口完整（列表/详情/报名/签到）
- [ ] BaaS 用户接口完整（微信登录/档案/儿童）
- [ ] 所有 BaaS 接口通过 X-API-Key 隔离数据

---

### Chunk 5.5 — 🎨 前端：UniApp模块模板开发（第一批） [3天] ⚡

**分配**: 前端设计师/开发工程师

**任务**:
1. **基础骨架** (`templates/base/`)：
   - `App.vue.hbs`, `main.js.hbs`, `manifest.json.hbs`, `pages.json.hbs`
   - `config.js.hbs`, `common/baas-sdk.js.hbs`, `common/request.js.hbs`
   - `components/tabbar.vue.hbs`
2. **📝 内容模块** (`templates/modules/cms/`)：
   - 文章列表页、文章详情页
   - `module.json` 声明文件
   - `api/content.js.hbs`
3. **🤖 AI顾问模块** (`templates/modules/ai-advisor/`)：
   - 对话页（聊天气泡+SSE流式+快捷入口+免责声明）
   - 症状评估页、生长记录页、健康报告页
   - `module.json`, `api/ai.js.hbs`

**验收标准**:
- [ ] 基础骨架生成后可在微信开发者工具中正常打开
- [ ] 内容模块页面渲染正常，API 调用 BaaS 接口
- [ ] AI 对话页支持 SSE 流式展示（打字效果）
- [ ] Tabbar 正确显示选中模块的菜单
- [ ] 所有 `{{变量}}` 占位符符合 injector.js 规范

> 🏗️ **架构师验证点 #6**：第一批模板生成→编译→运行可用。

---

### Chunk 5.6 — 🎨 前端：UniApp模块模板开发（第二批） [2.5天]

**分配**: 前端设计师/开发工程师

**任务**:
4. **🎪 活动模块** (`templates/modules/activity/`)：
   - 活动列表/详情/报名/签到页
   - AI 文案生成组件、AI 海报预览组件
5. **🛒 电商模块** (`templates/modules/shop/`)：
   - 商品列表/详情/购物车/结算/订单列表

**验收标准**:
- [ ] 活动模块完整（报名→签到全流程）
- [ ] 电商模块完整（浏览→加购→下单→支付）
- [ ] AI 文案组件在小程序中可用

---

## Sprint 6: 集成联调 + 部署上线 (Week 7-8)

### Chunk 6.1 — 🔧 后端+🎨 前端：全链路集成测试 [2天]

**分配**: 后端 + 前端 联合

**任务**:
1. **端到端测试用例**（按用户场景）：
   - 场景A：客户注册 → 创建项目 → 选3模块 → 构建 → 下载 → 编译 → 发布
   - 场景B：终端用户打开小程序 → 微信登录 → AI对话 → 查看活动 → 报名
   - 场景C：管理端写文章 → AI生成文案 → 发布 → 小程序查看文章
   - 场景D：小程序下单 → 微信支付 → 订单完成
2. Bug 修复
3. 性能基准测试

**验收标准**:
- [ ] 4 个核心场景全部跑通
- [ ] AI 对话端到端延迟 < 5 秒
- [ ] 代码生成 < 30 秒
- [ ] 无 P0/P1 级别 Bug

---

### Chunk 6.2 — 🔧 后端：Docker部署配置 [1天]

**分配**: 后端开发工程师

**任务**:
1. `docker-compose.full.yml`（全部服务：PG+Milvus+Redis+Python AI+Go Biz+Node.js CodeGen）
2. Nginx 反向代理配置
3. 环境变量生产配置
4. 健康检查脚本

**验收标准**:
- [ ] `docker compose -f docker-compose.full.yml up -d` 一键启动全部服务
- [ ] 所有服务健康检查通过
- [ ] Nginx 正确代理到各服务

---

### Chunk 6.3 — 🎨 前端：管理后台收尾 [1天]

**分配**: 前端设计师/开发工程师

**任务**:
1. `pages/Billing.tsx` 套餐对比页
2. 全局错误处理（网络异常、API 错误提示）
3. Loading 状态完善
4. 响应式适配收尾

---

### Chunk 6.4 — 🎨 前端：移动端管理工具 (UniApp) [2天]

**分配**: 前端设计师/开发工程师

**任务**:
1. 使用自己的 UniApp 框架生成一个管理工具小程序
2. 核心功能：项目列表查看、数据概览、构建状态推送

**验收标准**:
- [ ] 移动端可查看项目列表
- [ ] 移动端可查看基本数据统计

> 🏗️ **架构师最终验证点**：全系统就绪，可上线。

---

## 📊 总览

```
Sprint 0 (2天)    ██░░░░░░░░░░  基础设施
Sprint 1 (5天)    █████░░░░░░░  AI知识库管道
Sprint 2 (5天)    ██████░░░░░░  AI对话核心
Sprint 3 (5.5天)  ███████░░░░░  Go业务+多租户
Sprint 4 (12天)   ████████████  工厂管理+BaaS+代码生成引擎+管理后台 ⚡
Sprint 5 (10.5天) ██████████░░  AI增强+BaaS完整+UniApp模板 ⚡
Sprint 6 (6天)    ████████████  集成测试+部署上线
─────────────────────────────────────────────
合计: ~46天 (约9周)，含前后端并行优化
```

| 角色 | 总人天 | 占比 |
|------|--------|------|
| 🔧 后端 | ~28天 | 60% |
| 🎨 前端 | ~18天 | 40% |

---

## 🏗️ 架构师验证点（6个Gate）

| Gate | 位置 | 验证内容 |
|------|------|----------|
| #1 | Chunk 1.5 后 | 知识库管道端到端：文档→chunk→向量→Milvus检索 |
| #2 | Chunk 2.3 后 | AI对话链路：query→检索→重排→LLM→SSE |
| #3 | Chunk 3.3 后 | Go→Python全链路：Go入口→Python AI→SSE透传→客户端 |
| #4 | Chunk 3.4 后 | AI系统MVP完成：知识库+对话+儿童档案+产品推荐+多租户 |
| #5 | Chunk 4.6 后 | 代码生成引擎：项目配置→完整UniApp工程 ZIP |
| #6 | Chunk 5.5 后 | 第一批模板生成→微信开发者工具编译通过 |
| 🚀 | Chunk 6.4 后 | 最终验证：全系统就绪，可上线 |
