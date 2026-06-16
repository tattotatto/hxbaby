# AI 知识库管理页面 — 设计文档

> **日期**: 2026-06-17
> **范围**: 最小可用方案 — 文档上传 + 基础管理
> **影响**: ai-service + biz-service + factory-admin

---

## 1. 架构

```
浏览器 → Go Biz (:9090) → Python AI (:8001)
           ↑ JWT认证        ↑ loader→chunker→embedder→vector_store
```

## 2. API 设计

### Python AI Service (`ai-service/api/knowledge.py`)

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/ai/knowledge/upload` | multipart file → 入库 |
| GET | `/ai/knowledge/documents` | 已入库文档列表 |
| DELETE | `/ai/knowledge/documents/{source}` | 按 source 删除 |
| GET | `/ai/knowledge/stats` | 统计信息 |

### Go Biz Service (`biz-service/internal/handler/knowledge.go`)

代理层，JWT 中间件保护：
- `POST /api/v1/knowledge/upload`
- `GET /api/v1/knowledge/documents`
- `DELETE /api/v1/knowledge/documents/:source`
- `GET /api/v1/knowledge/stats`

## 3. 前端

- **页面**: `KnowledgeManager.tsx` — 统计卡片 + Upload.Dragger + Table
- **API**: `knowledge.ts` — 4 个 API 函数
- **路由**: `/knowledge`
- **导航**: 侧边栏新增 "知识库管理" 菜单项

## 4. 文件清单

| 文件 | 操作 |
|------|------|
| `ai-service/api/knowledge.py` | 新增 |
| `ai-service/main.py` | 修改(注册路由) |
| `biz-service/internal/handler/knowledge.go` | 新增 |
| `biz-service/internal/router/router.go` | 修改(注册路由) |
| `factory-admin/src/api/knowledge.ts` | 新增 |
| `factory-admin/src/pages/KnowledgeManager.tsx` | 新增 |
| `factory-admin/src/components/Layout.tsx` | 修改(菜单) |
| `factory-admin/src/App.tsx` | 修改(路由) |

## 5. 前端布局

```
┌─ 统计卡片行 (3 Card: 文档数/Chunk数/集合名) ─┐
│  Upload.Dragger (PDF/DOCX/TXT/MD/CSV/XLSX)  │
│  Table (文件名/格式/Chunks/时间/删除)          │
└──────────────────────────────────────────────┘
```
