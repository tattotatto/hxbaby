# 前端设计师/开发工程师 — 任务手册

> 基于 [小程序工厂设计](../specs/2026-06-12-miniapp-factory-design.md)

---

## 职责范围

| 项目 | 技术栈 | 职责 |
|------|--------|------|
| 管理后台 (PC Web) | React 18 + Ant Design 5 + Vite | 客户注册登录、小程序项目管理、模块配置、源码下载 |
| UniApp 模块模板 | Vue 3 + UniApp + Handlebars | 8大功能模块的页面模板开发 |
| BaaS SDK | JavaScript | 封装API调用，内置到生成的源码中 |
| 移动端工具 | UniApp | 客户移动端管理配套 |

> **注意**：后端API由后端工程师提供，前端按接口文档对接。

---

## 一、管理后台 (React + Ant Design 5)

### 1.1 项目初始化

```bash
npm create vite@latest factory-admin -- --template react-ts
cd factory-admin
npm install antd @ant-design/icons axios react-router-dom dayjs
npm install --save-dev @types/react-router-dom
```

**目录结构**:
```
factory-admin/src/
├── main.tsx              # 入口
├── App.tsx               # 路由配置
├── api/                  # API 调用封装
│   ├── client.ts         # axios 实例(拦截器: 自动附Token, 401跳转登录)
│   ├── auth.ts           # 登录注册API
│   ├── project.ts        # 项目管理API
│   ├── build.ts          # 构建下载API
│   └── ai.ts             # AI增强API
├── pages/                # 页面组件
│   ├── Login.tsx
│   ├── Register.tsx
│   ├── Dashboard.tsx     # 项目列表
│   ├── ProjectCreate.tsx # 创建向导(步骤1: 基本信息)
│   ├── ProjectConfig.tsx # 模块选择(步骤2)
│   ├── ProjectBrand.tsx  # 品牌配置(步骤3)
│   ├── ProjectBuild.tsx  # 构建下载(步骤4)
│   ├── ProjectGuide.tsx  # 接入指引
│   ├── ContentManager.tsx # 内容管理(文章CRUD)
│   ├── ActivityManager.tsx # 活动管理
│   ├── ShopManager.tsx   # 商品管理
│   └── Billing.tsx       # 套餐计费
├── components/           # 通用组件
│   ├── Layout.tsx        # 管理后台布局(Sider+Header+Content)
│   ├── ModuleSelector.tsx # 模块勾选卡片
│   ├── BrandConfig.tsx   # 品牌配置表单
│   ├── BuildProgress.tsx # 构建进度展示
│   └── AICopyPanel.tsx   # AI文案生成面板
└── hooks/
    └── useAuth.ts        # 认证状态Hook
```

### 1.2 页面清单与功能

#### 登录/注册页 (`Login.tsx`, `Register.tsx`)
- 手机号+密码登录
- 手机号+验证码登录(预留)
- 微信扫码登录(预留)
- 注册：手机号+密码+机构名称
- 登录成功后存储 token 到 localStorage，跳转 Dashboard

#### 项目列表 (`Dashboard.tsx`)
- **Header**: 欢迎语 + 套餐信息 + 升级按钮
- **统计卡片**: 项目数、本月API调用量、下载次数
- **项目卡片列表**:
  - 项目名称、状态(draft/active)
  - 已选模块标签展示
  - 操作按钮：配置、构建、下载、删除
- **空状态**: "还没有项目，点击创建第一个小程序"
- **创建按钮**: 跳转创建向导

#### 创建向导 (3步)
**Step 1 - 基本信息** (`ProjectCreate.tsx`):
- 小程序名称（必填）
- 一句话描述
- 微信AppID（选填，后续可在微信开放平台获取）

**Step 2 - 模块选择** (`ProjectConfig.tsx`):
- 8个模块以卡片形式展示，每张卡片包含：
  - 模块图标 + 名称
  - 功能简述
  - AI增强标签（带✨标识）
  - 依赖说明（如"需要基础模块"）
- 基础模块默认选中且不可取消
- 选中模块时自动检测依赖：如有缺失弹窗提示"将自动添加XX模块"
- 右侧实时显示已选模块摘要
- AI增强能力预览面板

**模块卡片参考设计**（每张卡片约 200×160px）:
```
┌─────────────────────┐
│ 🤖 AI健康顾问    ✨AI│
│ 智能问答·症状评估    │
│ 生长记录·健康报告    │
│                     │
│ 依赖: 基础模块       │
│ ☑️ 已选择            │
└─────────────────────┘
```

**Step 3 - 品牌配置** (`ProjectBrand.tsx`):
- 应用名称（默认取项目名）
- Logo上传（拖拽上传，预览）
- 主题色选择（预设色板 + 自定义色值）
- 辅助色选择
- 底部版权信息
- 实时预览：右侧模拟手机屏幕展示效果

#### 构建下载 (`ProjectBuild.tsx`)
- **构建按钮**: 一键触发，显示确认弹窗（模块清单 + 预计时间）
- **构建进度**: 6个阶段动画
  ```
  模块解析 → 模板拼装 → 配置注入 → 依赖合并 → 打包压缩 → 完成
  ```
- **完成状态**:
  - 下载按钮（下载ZIP包）
  - MD5校验值
  - 文件大小
  - 警告列表（如有）
  - 构建历史记录
  - "查看接入指引"跳转按钮

#### 接入指引 (`ProjectGuide.tsx`)
- **4步接入流程**（带步骤编号和插图）:
  1. 下载源码并解压
  2. 安装依赖：`npm install`
  3. 配置微信AppID（在 `src/manifest.json` 中修改）
  4. 打开微信开发者工具，导入项目目录
- **API配置信息卡片**（API Key + Secret，支持一键复制）
- **BaaS接口文档链接**
- **常见问题FAQ**

#### 内容管理 (`ContentManager.tsx`)
- 文章列表（表格：标题、分类、发布时间、阅读量、AI生成标记）
- 新建文章：富文本编辑器 + **AI写作助手面板**（输入主题→一键生成→编辑→发布）
- AI辅助：自动摘要、自动标签、配图推荐

#### 活动管理 (`ActivityManager.tsx`)
- 活动列表（表格：标题、时间、报名数、状态）
- 新建活动：基本信息表单 + **AI文案面板**（输入主题→生成标题+详情+推文）
- **AI海报生成**：选择模板 → AI生成海报预览 → 下载

#### AI文案面板组件 (`AICopyPanel.tsx`)
- 输入区：主题/关键词
- 生成按钮：带loading动画
- 结果展示：多Tab切换（文章/摘要/推文/朋友圈）
- 一键复制/一键应用到编辑器

---

## 二、UniApp 模块模板

### 2.1 模板技术规范

- **框架**: UniApp (Vue 3 Composition API)
- **模板引擎**: Handlebars (`.hbs` 文件，变量占位符 `{{key}}`)
- **样式**: 每个模块独立 `.vue` 单文件，使用 `<style scoped>`
- **API调用**: 通过 `baas-sdk.js` 统一封装

### 2.2 基础骨架模板 (`templates/base/`)

| 文件 | 说明 |
|------|------|
| `App.vue.hbs` | 应用入口，引入全局样式 |
| `main.js.hbs` | Vue实例化，挂载Pinia、全局组件 |
| `manifest.json.hbs` | UniApp配置（含 `{{wxAppId}}` 占位） |
| `pages.json.hbs` | 路由骨架（具体页面由编排引擎拼接） |
| `config.js.hbs` | **核心**：BaaS配置(`{{baasBaseURL}}`, `{{apiKey}}`, `{{projectId}}`) |
| `common/baas-sdk.js.hbs` | BaaS SDK（详见下方） |
| `common/request.js.hbs` | HTTP请求封装（自动附API Key） |
| `components/tabbar.vue.hbs` | 底部导航栏（菜单项由编排引擎注入） |

### 2.3 模块模板清单

#### 📝 内容模块 (`templates/modules/cms/`)
| 文件 | 说明 | 设计要点 |
|------|------|----------|
| `module.json` | 模块声明 | name:cms, 依赖base, 2页面, tabbar配置 |
| `pages/article-list.vue.hbs` | 文章列表 | 分类Tab切换、下拉刷新、触底加载更多、封面图+标题+摘要卡片 |
| `pages/article-detail.vue.hbs` | 文章详情 | 富文本渲染、阅读数、分享按钮、相关推荐 |
| `api/content.js.hbs` | 内容API | `getArticles(category, page)`, `getArticle(id)` |

**设计风格**: 母婴温馨风格，卡片式布局，封面图+标题+摘要+阅读量

#### 🤖 AI顾问模块 (`templates/modules/ai-advisor/`)
| 文件 | 说明 | 设计要点 |
|------|------|----------|
| `module.json` | 模块声明 | 依赖base, 4页面, tabbar配置 |
| `pages/chat.vue.hbs` | AI对话 | 聊天气泡(用户/AI双边)、SSE流式打字效果、快捷问题入口、免责声明 |
| `pages/symptom-check.vue.hbs` | 症状评估 | 症状选择+文字描述、评估结果卡片(风险等级颜色标识)、建议列表 |
| `pages/growth-record.vue.hbs` | 生长记录 | 身高体重折线图(echarts)、WHO标准曲线对比、记录表单 |
| `pages/health-report.vue.hbs` | 健康报告 | 综合评估报告卡片、生长趋势、AI建议、分享报告 |
| `api/ai.js.hbs` | AI API | `aiChat(message, childId, stream)`, `aiAssess(symptoms)` |

**设计风格**: 医疗专业感+亲和力，蓝绿色调，聊天气泡清晰区分医患角色

#### 🛒 电商模块 (`templates/modules/shop/`)
| 文件 | 说明 |
|------|------|
| `module.json` | 依赖base, 5页面 |
| `pages/product-list.vue.hbs` | 商品列表（分类筛选、搜索、网格布局） |
| `pages/product-detail.vue.hbs` | 商品详情（轮播图、价格、AI卖点标签、规格选择、加入购物车） |
| `pages/cart.vue.hbs` | 购物车（全选、数量调整、合计、去结算） |
| `pages/checkout.vue.hbs` | 结算页（地址选择、订单确认、微信支付按钮） |
| `pages/order-list.vue.hbs` | 订单列表（Tab:全部/待付款/已发货/已完成） |
| `api/shop.js.hbs` | 电商API |

#### 🎪 活动模块 (`templates/modules/activity/`)
| 文件 | 说明 |
|------|------|
| `module.json` | 依赖base, 4页面, tabbar配置 |
| `pages/activity-list.vue.hbs` | 活动列表（进行中/即将开始/已结束Tab） |
| `pages/activity-detail.vue.hbs` | 活动详情（海报头图、AI生成文案、报名按钮） |
| `pages/activity-signup.vue.hbs` | 报名表单 |
| `components/ai-copy-gen.vue.hbs` | AI文案生成组件 |
| `components/ai-poster.vue.hbs` | AI海报预览组件 |
| `api/activity.js.hbs` | 活动API |

#### 📅 预约模块 (`templates/modules/booking/`)
- 顾问列表 → 时段选择 → 预约确认 → 我的预约

#### 👑 会员模块 (`templates/modules/member/`)
- 会员中心（等级+积分）、优惠券列表、权益说明

#### 📊 数据模块 (`templates/modules/analytics/`)
- 儿童生长曲线图表、健康趋势、消费记录

### 2.4 BaaS SDK 设计

**文件**: `templates/base/common/baas-sdk.js.hbs`

```javascript
class BaaS {
  constructor(config) {
    this.baseURL = config.baseURL;
    this.apiKey = config.apiKey;
    this.projectId = config.projectId;
  }

  async _request(method, path, data, options = {}) {
    const url = `${this.baseURL}/baas/v1${path}`;
    const headers = {
      'Content-Type': 'application/json',
      'X-API-Key': this.apiKey,
    };
    // 如果已登录，附加用户token
    const userToken = uni.getStorageSync('baas_user_token');
    if (userToken) headers['Authorization'] = `Bearer ${userToken}`;

    const res = await uni.request({
      url, method, data, header: headers, ...options,
    });
    if (res.statusCode >= 400) throw new Error(res.data?.message || '请求失败');
    return res.data?.data;
  }

  // 内容
  getArticles(params) { return this._request('GET', '/articles', params); }
  getArticle(id) { return this._request('GET', `/articles/${id}`); }

  // AI对话 (SSE流式)
  aiChatStream(message, childId, onToken, onDone) {
    // 通过 wx.request + enableChunked 实现SSE流式接收
  }

  // 电商
  getProducts(params) { return this._request('GET', '/products', params); }
  createOrder(data) { return this._request('POST', '/orders', data); }
  wechatPay(orderId) { return this._request('POST', '/payment/wechat', { order_id: orderId }); }

  // 活动
  getActivities(params) { return this._request('GET', '/activities', params); }
  signupActivity(id, data) { return this._request('POST', `/activities/${id}/signup`, data); }

  // 用户
  wxLogin(code) { return this._request('POST', '/auth/wx-login', { code }); }
  getProfile() { return this._request('GET', '/user/profile'); }
}

export { BaaS };
```

---

## 三、前后端接口对接

### 管理后台调用的接口（需要JWT登录态）

| 页面 | 调用的API | 说明 |
|------|-----------|------|
| Login | `POST /api/v1/auth/login` | 手机号+密码 |
| Register | `POST /api/v1/auth/register` | 注册 |
| Dashboard | `GET /api/v1/projects` | 项目列表 |
| ProjectCreate | `POST /api/v1/projects` | 创建项目 |
| ProjectConfig | `PUT /api/v1/projects/:id` | 更新模块/品牌 |
| ProjectBuild | `POST /api/v1/projects/:id/build` | 触发构建 |
| ProjectBuild | `GET /api/v1/builds/:id/status` | 轮询状态 |
| ProjectBuild | `GET /api/v1/builds/:id/download` | 下载源码 |
| ProjectGuide | `GET /api/v1/projects/:id` | 获取API Key信息 |
| ContentManager | `GET/POST /api/v1/admin/articles` | 内容管理 |
| ActivityManager | `GET/POST /api/v1/admin/activities` | 活动管理 |
| AI写作面板 | `POST /api/v1/ai/generate-article` | AI写文章 |
| AI活动文案 | `POST /api/v1/ai/generate-activity-copy` | AI文案 |
| AI海报 | `POST /api/v1/ai/generate-poster` | AI海报 |

### 生成的UniApp调用的接口（X-API-Key鉴权）

参考 `baas-sdk.js` 封装，所有请求自动附带 `X-API-Key` header。

---

## 四、设计系统（Design Tokens）

### 色彩
```
主题色(健康绿): #4caf50 (可被租户品牌色覆盖)
辅助色(温暖橙): #ff9800
AI标识色(科技紫): #7c4dff
文字主色: #333333
文字辅色: #666666
背景色: #f5f5f5
卡片背景: #ffffff
危险/紧急: #f44336
```

### 组件规范
- **卡片**: 圆角 12px，阴影 `0 2px 8px rgba(0,0,0,0.08)`
- **按钮**: 主按钮填充、次按钮描边，圆角 8px
- **输入框**: 高度 44px，圆角 8px，聚焦边框主题色
- **间距**: 基础单位 8px，页面边距 16px

### 字体
- 中文: PingFang SC / Microsoft YaHei
- 数字: DIN Alternate
- 标题: 18px/700, 正文: 14px/400, 辅助: 12px/400

---

## 五、开发顺序

```
Week 1-2: 管理后台骨架 → 登录注册 → Dashboard → 项目CRUD
Week 2-3: 创建向导(3步骤) → 模块选择器 → 品牌配置 → 构建触发+进度
Week 3-4: UniApp基础骨架模板 → baa-sdk.js → 内容模块模板(cms)
Week 4-5: AI顾问模块模板 → 活动模块模板(含AI文案/海报组件)
Week 5-6: 电商模块模板(商品/购物车/订单/支付) → BaaS对接
Week 6-7: 接入指引页面 → 内容管理后台 → AI文案面板 → 活动管理后台
Week 7-8: 预约+会员+数据模块模板 → 移动端管理工具(UniApp) → 联调测试
```

## 六、Node.js模板引擎对接说明

前端开发的 `.vue.hbs` 模板文件中的 `{{变量名}}` 占位符由Node.js代码生成引擎在构建时替换为实际值：

| 占位符 | 示例值 | 注入来源 |
|--------|--------|----------|
| `{{appName}}` | 宝宝健康助手 | brand_config |
| `{{primaryColor}}` | #4caf50 | brand_config |
| `{{logo}}` | https://oss.../logo.png | brand_config |
| `{{baasBaseURL}}` | https://api.hxbaby.com | 环境变量 |
| `{{apiKey}}` | ak_abc123 | 项目自动生成 |
| `{{projectId}}` | 123 | 项目ID |
| `{{wxAppId}}` | wxabc123 | 客户配置 |
| `{{footer}}` | © 2026 宝宝健康助手 | brand_config |

**注意**: 模板文件使用 `.hbs` 扩展名，编排引擎会自动去除 `.hbs` 后缀。例如 `chat.vue.hbs` 生成后变为 `chat.vue`。
