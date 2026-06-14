# 多用户小程序工厂平台 — 设计文档

**日期**: 2026-06-12
**状态**: 已确认
**版本**: v1.0
**关联**: [AI资料库系统设计文档](2026-06-12-children-health-ai-design.md)

---

## 1. 项目概述

多用户小程序工厂平台是一个SaaS系统，允许合作机构（母婴机构等）在平台上快速创建和发布自己的微信小程序。客户在Web管理后台勾选功能模块、填写品牌配置，系统自动生成完整的UniApp跨端源码包（微信小程序 + H5），并提供全托管BaaS后端服务。

### 1.1 核心价值

```
客户注册 → 勾选模块 → 配置品牌 → 一键生成源码 → 下载 → 微信开发者工具打开 → 发布上线
```

客户只需聚焦内容和运营，后端零运维，AI能力开箱即用。

### 1.2 关键决策

| 维度 | 决策 |
|------|------|
| 目标用户 | 现有合作机构 + 未来开放注册客户（SaaS多租户） |
| 模块策略 | 全模块可插拔，客户按需自由组合 |
| 代码输出 | UniApp 跨端源码（编译到微信小程序 + H5） |
| 后端模式 | 全托管BaaS，客户下载前端源码，零后端运维 |
| 管理后台 | PC Web (React) + 移动端配套工具 |
| 技术架构 | Go扩展BaaS + Node.js代码生成引擎 |
| AI集成 | 统一AI Bridge → 复用AI资料库系统 |

---

## 2. 系统架构

### 2.1 架构模式

**Go扩展方案 (Node.js + Go)**：在AI系统Go业务服务基础上扩展BaaS能力，新增Node.js微服务处理UniApp代码生成。

### 2.2 架构全景

```
┌─────────────────────────────────────────────────┐
│  客户管理端: PC Web (React)  +  移动端工具          │
└────────────────────┬────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────┐
│  Go · Gin 统一入口                                 │
│  JWT · 租户识别 · 限流 · 统一日志                    │
├──────────────────┬──────────────────────────────┤
│  工厂管理服务      │  BaaS 运行时服务               │
│  (Go 新增)        │  (Go 扩展 + 复用AI系统)         │
│                  │                              │
│  · 客户注册/登录   │  📝 内容(AI增强)              │
│  · 套餐/计费      │  🎪 活动(AI增强)              │
│  · 项目CRUD      │  🛒 电商(AI增强)              │
│  · 模块选择配置    │  👤 用户/微信登录              │
│  · 域名/证书      │  🤖 AI Bridge → AI系统        │
│  · 源码下载       │  📅 预约 · 👑 会员 · 📊 数据   │
└──────┬───────────┴──────────────┬───────────────┘
       │                          │
       ▼                          ▼
┌──────────────┐    ┌─────────────────────────────┐
│ Node.js 微服务│    │  Python AI 服务 (复用)         │
│ 代码生成引擎   │    │  LLM · RAG · Embedding       │
│              │    │  文案/海报/分析/问答            │
│ · 模板库      │    └─────────────────────────────┘
│ · 编排引擎     │
│ · 配置注入     │
│ · ZIP打包     │
└──────┬───────┘
       │
       ▼
┌─────────────────────────────────────────────────┐
│  PostgreSQL (复用) · Redis (复用) · OSS (复用)      │
└─────────────────────────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────────────────┐
│  客户下载UniApp源码 → 微信开发者工具 → 发布上线       │
│  运行时调用 BaaS API (/baas/v1/*) → AI能力         │
└─────────────────────────────────────────────────┘
```

---

## 3. 模块可插拔架构

### 3.1 模块清单

| 模块 | 功能 | AI增强 | 依赖 |
|------|------|--------|------|
| 🏠 基础（必选） | 首页 · 用户中心 · 微信登录 · Tabbar | - | 无 |
| 📝 内容 | 文章 · 公告 · 轮播图 · 分类 | 写文章 · 配图 · 摘要 · 标签 · 选题推荐 | 基础 |
| 🤖 AI顾问 | AI对话 · 症状评估 · 生长记录 · 健康报告 | 核心AI服务（全部） | 基础 |
| 🎪 活动 | 活动列表 · 报名 · 签到 · 分享裂变 | 文案 · 海报 · 人群定向 · 复盘报告 | 基础 + 电商(可选) |
| 🛒 电商 | 商品 · 购物车 · 订单 · 微信支付 | 卖点 · 分类 · 关联推荐 · 商品图 | 基础 + 微信支付 |
| 📅 预约 | 顾问列表 · 时段 · 确认 · 提醒 | 排期建议 · 提醒文案 | 基础 + 微信支付(可选) |
| 👑 会员 | 等级 · 积分 · 优惠券 · 权益 | 个性化推送 · 权益推荐 | 基础 + 电商(可选) |
| 📊 数据 | 生长曲线 · 健康趋势 · 消费记录 | AI报告生成 | AI顾问 + 电商(可选) |

### 3.2 依赖解析规则

- 代码生成引擎在组合模块时自动进行依赖检测
- 选中模块但没有选其依赖模块时，自动补齐并提示用户
- 模块冲突时（如选了电商但没选微信支付），给出警告

---

## 4. 代码生成引擎

### 4.1 生成流程

```
客户触发 → 模块解析(依赖检测+冲突检查) → 模板拼装(合并pages.json/manifest/App.vue)
→ 配置注入(品牌/API端点/租户ID) → npm依赖合并 → ZIP打包 → OSS上传 → 下载链接
```

### 4.2 模板组织结构

```
templates/
├── base/              # 基础骨架（始终包含）
│   ├── App.vue.hbs
│   ├── main.js.hbs
│   ├── manifest.json.hbs
│   ├── pages.json.hbs       # 路由由编排引擎动态拼接
│   ├── config.js.hbs        # BaaS配置注入点
│   └── components/          # 通用组件(tabbar等)
│
├── modules/           # 功能模块模板（每个独立目录）
│   ├── cms/
│   ├── ai-advisor/
│   ├── shop/
│   ├── activity/            # 含AI文案/海报组件
│   ├── booking/
│   ├── member/
│   └── analytics/
│       └── module.json      # 模块声明: 名称/依赖/路由定义/API依赖
│
├── scripts/           # 生成引擎脚本
│   ├── compose.js           # 编排主入口
│   ├── resolver.js          # 依赖解析+拓扑排序
│   ├── injector.js          # 配置注入器
│   ├── packager.js          # ZIP打包器
│   └── validator.js         # 合法性校验
│
└── output/            # 生成产物临时目录（按task_id隔离）
```

### 4.3 技术选型

| 组件 | 选型 |
|------|------|
| Web框架 | Express / Fastify |
| 模板引擎 | Handlebars (.hbs) |
| 打包工具 | archiver (npm) |
| 运行时 | Node.js 20+ |

---

## 5. 数据模型

### 5.1 核心实体

**工厂管理层**:
- `customers`: 客户/商户（id, name, phone, email, company_name, plan, max_projects, status）
- `miniapp_projects`: 小程序项目（id, customer_id, name, wx_app_id, wx_mch_id, modules(JSONB), brand_config(JSONB), api_key, status, domain）
- `build_tasks`: 构建任务（id, project_id, modules_snapshot, config_snapshot, status, output_zip_url, output_md5, duration_ms）

**BaaS数据层**:
- `cms_articles`: 文章（id, project_id, title, content, summary, cover_image, ai_generated, ai_prompt, ai_images(JSONB)）
- `b_activities`: 活动（id, project_id, title, description, ai_copy(JSONB), ai_poster_url, ai_report(JSONB), signup_start, signup_end, checkin_enabled）
- `b_products`: 商品（id, project_id, name, description, ai_selling_points, ai_tags(JSONB), price, stock）
- `b_orders`: 订单（id, project_id, user_id, items(JSONB), total_amount, wx_transaction_id, status）
- `b_users`: 终端用户（id, project_id, wx_openid, nickname, children(JSONB)）
- `b_bookings`: 预约（id, project_id, user_id, consultant_id, slot_time, status）
- `b_members`: 会员（id, project_id, user_id, level, points, coupons(JSONB)）
- `b_activity_signups`: 活动报名（id, activity_id, user_id, checkin_at）
- `b_api_logs`: API调用日志（id, project_id, endpoint, method, status, duration_ms）

---

## 6. API设计

### 6.1 管理端API (`/api/v1/`)

| 分组 | 核心接口 |
|------|----------|
| 客户认证 | `POST /auth/register`, `POST /auth/login`, `GET /customer/profile`, `PUT /customer/profile` |
| 项目管理 | `GET /projects`, `POST /projects`, `GET /projects/:id`, `PUT /projects/:id`, `POST /projects/:id/build` |
| 代码生成 | `GET /projects/:id/builds`, `GET /builds/:id/status`, `GET /builds/:id/download`, `POST /builds/:id/retry` |
| 模板预览 | `GET /templates/modules`, `POST /templates/preview` |
| AI增强 | `POST /ai/generate-article`, `POST /ai/generate-summary`, `POST /ai/generate-activity-copy`, `POST /ai/generate-poster`, `POST /ai/generate-selling-points`, `POST /ai/activity-report` |

### 6.2 BaaS运行时API (`/baas/v1/`)

供生成的小程序调用，通过API Key鉴权和项目隔离：

| 分组 | 核心接口 |
|------|----------|
| 认证 | `POST /auth/wx-login` |
| 内容 | `GET /articles`, `GET /articles/:id` |
| AI对话 | `POST /ai/chat`, `POST /ai/assess`, `POST /ai/report` |
| 电商 | `GET /products`, `GET /products/:id`, `POST /orders`, `POST /payment/wechat` |
| 活动 | `GET /activities`, `GET /activities/:id`, `POST /activities/:id/signup`, `POST /activities/:id/checkin` |
| 用户 | `GET /user/profile`, `PUT /user/profile`, `POST /user/children` |
| 预约 | `GET /bookings/slots`, `POST /bookings` |

### 6.3 UniApp SDK

生成的小程序源码中内置 `baas-sdk.js`：

```javascript
// 初始化（自动从config.js读取）
import { BaaS } from '@/common/baas-sdk'
const baas = new BaaS({ apiKey, projectId, baseURL })

// 使用
const articles = await baas.getArticles({ category: 'health' })
const reply = await baas.aiChat({ message: '宝宝发烧怎么办', childId: 'xxx', stream: true })
const order = await baas.createOrder({ productId, quantity })
```

---

## 7. AI增强服务

### 7.1 统一AI调用架构

```
各业务服务 → AI Bridge (Go) → Python AI服务 → LLM/图像模型/RAG
```

### 7.2 AI能力矩阵

| 服务 | AI能力 | 调用模型 |
|------|--------|----------|
| 📝 内容 | 写文章、配图推荐、摘要、标签、选题推荐 | LLM文本 + 图像匹配 |
| 🎪 活动 | 文案生成、海报生成、人群定向、复盘报告 | LLM文本 + 图像生成 + 数据分析 |
| 🛒 电商 | 卖点提炼、智能分类、关联推荐、商品图优化 | LLM文本 + RAG匹配 + 图像处理 |
| 🤖 AI顾问 | 健康问答、症状评估、生长分析、产品推荐 | LLM + RAG知识检索（核心） |
| 📅 预约 | 智能排期建议、自动提醒文案 | LLM文本 |
| 👑 会员 | 个性化推送、权益推荐 | 数据分析 + LLM文案 |

---

## 8. 计费模型

| 版本 | 价格 | 项目数 | 核心限制 |
|------|------|--------|----------|
| 🆓 免费版 | ¥0 | 1个 | CMS基础模块，AI对话100次/月，源码下载3次/月 |
| 💚 基础版 | ¥299/月 | 3个 | 内容+AI顾问+活动，AI对话2000次/月，AI文案100次/月 |
| 💙 专业版 | ¥999/月 | 10个 | 全模块含电商/预约/会员，AI对话10000次/月，AI文案500次/月 |
| ❤️ 企业版 | 定制 | 无限 | 全部功能+定制模块+独立资源池+私有化部署可选 |

---

## 9. 技术栈总览

| 层级 | 组件 | 选型 |
|------|------|------|
| Go扩展 | 工厂管理+BaaS API+AI Bridge | Go · Gin (扩展现有服务) |
| Node.js | 代码生成引擎 | Express + Handlebars + archiver |
| 管理前端 | PC Web | React 18 + Ant Design 5 |
| 移动端工具 | 移动端配套 | UniApp |
| 数据库 | 业务数据 | PostgreSQL (扩展现有实例) |
| 缓存 | 高性能缓存 | Redis (复用) |
| 文件 | 源码包/素材/海报 | 阿里云OSS (复用) |
| 部署 | 容器编排 | Docker + K8s (复用) |

---

## 10. 两个系统的关系

```
┌──────────────────────────────────────────────────┐
│              小程序工厂平台 (本文档)                  │
│                                                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────┐  │
│  │ 管理后台  │  │ BaaS API │  │ AI Bridge ───────┼──┼──→ 内部调用
│  │ (React)  │  │ (Go扩展) │  │ (Go)             │  │
│  └──────────┘  └──────────┘  └──────────────────┘  │
│                                                      │
│  ┌──────────────────────────────────────────────┐   │
│  │  代码生成引擎 (Node.js)                         │   │
│  └──────────────────────────────────────────────┘   │
└──────────────────────────┬───────────────────────┘
                           │ 内部API
                           ▼
┌──────────────────────────────────────────────────┐
│              AI资料库系统 (关联文档)                 │
│                                                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────┐  │
│  │ Go 业务   │  │ Python   │  │ 知识库管道         │  │
│  │ (复用)    │  │ AI服务    │  │ (共享+租户产品)    │  │
│  └──────────┘  └──────────┘  └──────────────────┘  │
└──────────────────────────────────────────────────┘
```

**两个系统共享**：Go框架、PostgreSQL、Redis、OSS、Milvus、K8s集群

---

## 11. MVP与路线图

### 11.1 MVP (Phase 1) — 6-8周

| 周次 | 重点 |
|------|------|
| Week 1-2 | Go扩展：客户管理+项目管理API；管理后台：注册登录+项目CRUD |
| Week 2-3 | Node.js代码生成引擎MVP（基础+内容+AI顾问三模块）；模板库搭建 |
| Week 3-4 | BaaS运行时API（内容、用户、AI桥接）；UniApp SDK封装 |
| Week 4-5 | 微信支付集成；活动+电商模块；AI增强API（文案、海报） |
| Week 5-6 | 管理后台完善（配置向导、接入指引、下载管理） |
| Week 6-7 | 联调测试：管理后台→代码生成→BaaS→AI系统全链路 |
| Week 7-8 | 部署上线、灰度发布、文档完善 |

### 11.2 后续阶段

- **Phase 2** (+4周): 预约+会员+数据模块、AI海报质量优化、移动端管理工具
- **Phase 3** (+4周): 分析面板、A/B测试、知识库客户自主贡献、模板市场
- **Phase 4+** (长期): 多平台（支付宝/抖音小程序）、第三方模板开发者生态

---

## 12. 关键风险与对策

| 风险 | 对策 |
|------|------|
| 代码生成UniApp兼容性问题 | 建立自动化测试套件，每个模块在微信开发者工具中验证 |
| BaaS多租户数据泄露 | project_id强制注入，API Key+签名双重校验，渗透测试 |
| AI调用成本失控 | 免费版严格限流，AI Bridge统一计量+预算告警 |
| 微信支付接入门槛 | 提供标准化接入文档+视频教程，专人协助首单 |
| 模板维护成本随模块增多而增长 | 模块独立版本管理，自动化CI检测模板编译 |
