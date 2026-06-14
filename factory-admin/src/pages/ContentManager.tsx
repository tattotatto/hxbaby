import { useState } from 'react';
import { Card, Table, Button, Modal, Input, Tag, Space, message, Row, Col } from 'antd';
import { PlusOutlined, EyeOutlined } from '@ant-design/icons';
import AppLayout from '../components/Layout';
import AICopyPanel from '../components/AICopyPanel';

const { TextArea } = Input;

export default function ContentManager() {
  const [articles, setArticles] = useState<any[]>([]);
  const [showEditor, setShowEditor] = useState(false);
  const [editorContent, setEditorContent] = useState({ title: '', content: '', category: '' });

  const columns = [
    { title: '标题', dataIndex: 'title', key: 'title' },
    { title: '分类', dataIndex: 'category', key: 'category', render: (v: string) => v ? <Tag>{v}</Tag> : '-' },
    { title: '阅读量', dataIndex: 'view_count', key: 'views', width: 80 },
    {
      title: 'AI生成', dataIndex: 'ai_generated', key: 'ai', width: 80,
      render: (v: boolean) => v ? <Tag color="purple">AI</Tag> : <Tag>手动</Tag>,
    },
    { title: '状态', dataIndex: 'is_published', key: 'status', width: 80, render: (v: boolean) => v ? <Tag color="green">已发布</Tag> : <Tag>草稿</Tag> },
    {
      title: '操作', key: 'action', width: 100,
      render: () => <Button type="link" icon={<EyeOutlined />}>查看</Button>,
    },
  ];

  const handleSave = () => {
    message.success('文章已保存（API对接中）');
    setShowEditor(false);
  };

  return (
    <AppLayout>
      <Row gutter={24}>
        <Col span={16}>
          <Card
            title="内容管理"
            extra={<Button type="primary" icon={<PlusOutlined />} onClick={() => setShowEditor(true)}>新建文章</Button>}
          >
            <Table columns={columns} dataSource={articles} rowKey="id" locale={{ emptyText: '暂无文章，点击"新建文章"开始创作' }} />
          </Card>
        </Col>
        <Col span={8}>
          <AICopyPanel mode="article" />
        </Col>
      </Row>

      <Modal
        title="新建文章"
        open={showEditor}
        onCancel={() => setShowEditor(false)}
        onOk={handleSave}
        width={800}
        okText="保存"
        cancelText="取消"
      >
        <Space direction="vertical" style={{ width: '100%' }}>
          <Input
            placeholder="文章标题"
            value={editorContent.title}
            onChange={e => setEditorContent({ ...editorContent, title: e.target.value })}
            size="large"
          />
          <Input
            placeholder="分类"
            value={editorContent.category}
            onChange={e => setEditorContent({ ...editorContent, category: e.target.value })}
          />
          <TextArea
            placeholder="文章内容，或使用右侧AI助手生成"
            value={editorContent.content}
            onChange={e => setEditorContent({ ...editorContent, content: e.target.value })}
            rows={15}
          />
        </Space>
      </Modal>
    </AppLayout>
  );
}
