import { useState, useEffect, useCallback } from 'react';
import { Card, Upload, Table, Button, Modal, message, Space, Tag, Spin, Empty, Statistic, Row, Col } from 'antd';
import { InboxOutlined, DeleteOutlined, ReloadOutlined, FileTextOutlined, DatabaseOutlined } from '@ant-design/icons';
import type { UploadProps } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import AppLayout from '../components/Layout';
import {
  getDocuments,
  deleteDocument,
  getKnowledgeStats,
  uploadDocument,
  type DocumentItem,
  type KnowledgeStats,
} from '../api/knowledge';

const ACCEPT = '.pdf,.docx,.txt,.md,.csv,.xlsx';

export default function KnowledgeManager() {
  const [documents, setDocuments] = useState<DocumentItem[]>([]);
  const [stats, setStats] = useState<KnowledgeStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [uploading, setUploading] = useState(false);

  const fetchData = useCallback(async () => {
    try {
      const [docsRes, statsRes] = await Promise.all([
        getDocuments(),
        getKnowledgeStats(),
      ]);
      setDocuments(docsRes.documents ?? []);
      setStats(statsRes);
    } catch (err: any) {
      message.error(err?.response?.data?.message ?? '获取知识库数据失败');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleUpload: UploadProps['customRequest'] = async (options) => {
    const { file, onSuccess, onError } = options as any;
    setUploading(true);
    try {
      const result = await uploadDocument(file as File);
      message.success(`文档 "${result.source}" 上传成功，${result.chunks} 个分块已入库`);
      onSuccess?.(result, file);
      fetchData();
    } catch (err: any) {
      const msg = err?.response?.data?.message ?? err?.response?.data?.detail ?? '上传失败';
      message.error(msg);
      onError?.(err);
    } finally {
      setUploading(false);
    }
  };

  const handleDelete = (source: string) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除文档 "${source}" 的所有知识分块吗？此操作不可撤销。`,
      okText: '确定删除',
      okType: 'danger',
      cancelText: '取消',
      onOk: async () => {
        try {
          const result = await deleteDocument(source);
          message.success(`已删除 ${result.deleted_chunks} 个分块`);
          fetchData();
        } catch (err: any) {
          message.error(err?.response?.data?.message ?? '删除失败');
        }
      },
    });
  };

  const formatIcons: Record<string, string> = {
    pdf: '📄',
    docx: '📝',
    txt: '📃',
    md: '📋',
    csv: '📊',
    xlsx: '📈',
  };

  const columns: ColumnsType<DocumentItem> = [
    {
      title: '文件名',
      dataIndex: 'source',
      key: 'source',
      render: (text: string, record: DocumentItem) => (
        <Space>
          <span>{formatIcons[record.format] ?? '📁'}</span>
          <span>{text}</span>
        </Space>
      ),
    },
    {
      title: '格式',
      dataIndex: 'format',
      key: 'format',
      width: 80,
      render: (fmt: string) => <Tag>{fmt.toUpperCase()}</Tag>,
    },
    {
      title: '分块数',
      dataIndex: 'chunk_count',
      key: 'chunk_count',
      width: 100,
    },
    {
      title: '入库时间',
      dataIndex: 'indexed_at',
      key: 'indexed_at',
      width: 180,
      render: (t: string) => (t ? new Date(t).toLocaleString('zh-CN') : '-'),
    },
    {
      title: '操作',
      key: 'actions',
      width: 80,
      render: (_: unknown, record: DocumentItem) => (
        <Button
          type="text"
          danger
          icon={<DeleteOutlined />}
          onClick={() => handleDelete(record.source)}
        />
      ),
    },
  ];

  return (
    <AppLayout>
      <Spin spinning={loading}>
        {/* 统计卡片 */}
        <Row gutter={16} style={{ marginBottom: 24 }}>
          <Col span={8}>
            <Card>
              <Statistic
                title="文档总数"
                value={stats?.total_documents ?? 0}
                prefix={<FileTextOutlined />}
              />
            </Card>
          </Col>
          <Col span={8}>
            <Card>
              <Statistic
                title="分块总数"
                value={stats?.total_chunks ?? 0}
                prefix={<DatabaseOutlined />}
              />
            </Card>
          </Col>
          <Col span={8}>
            <Card>
              <Statistic
                title="集合名称"
                value={stats?.collection_name ?? '-'}
                prefix={<DatabaseOutlined />}
              />
            </Card>
          </Col>
        </Row>

        {/* 上传区 */}
        <Card title="📤 文档上传" style={{ marginBottom: 24 }}>
          <Upload.Dragger
            accept={ACCEPT}
            customRequest={handleUpload}
            showUploadList={false}
            disabled={uploading}
          >
            <p className="ant-upload-drag-icon">
              <InboxOutlined />
            </p>
            <p className="ant-upload-text">点击或拖拽文件到此区域上传</p>
            <p className="ant-upload-hint">
              支持 PDF / DOCX / TXT / Markdown / CSV / Excel，单文件不超过 50MB
            </p>
            {uploading && <p style={{ color: '#4caf50' }}>正在处理文档...</p>}
          </Upload.Dragger>
        </Card>

        {/* 文档列表 */}
        <Card
          title="📋 已入库文档"
          extra={
            <Button icon={<ReloadOutlined />} onClick={fetchData}>
              刷新
            </Button>
          }
        >
          {documents.length > 0 ? (
            <Table
              columns={columns}
              dataSource={documents}
              rowKey="source"
              pagination={{ pageSize: 20, showSizeChanger: false }}
            />
          ) : (
            <Empty description="暂无文档，请上传知识库文档" />
          )}
        </Card>
      </Spin>
    </AppLayout>
  );
}
