import { useState } from 'react';
import { Card, Table, Button, Modal, Input, DatePicker, Tag, Space, message, Row, Col } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import AppLayout from '../components/Layout';
import AICopyPanel from '../components/AICopyPanel';

const { TextArea } = Input;
const { RangePicker } = DatePicker;

export default function ActivityManager() {
  const [activities, setActivities] = useState<any[]>([]);
  const [showEditor, setShowEditor] = useState(false);
  const [editorContent, setEditorContent] = useState({
    title: '', description: '', location: '',
    ai_copy: '',
  });

  const columns = [
    { title: '标题', dataIndex: 'title', key: 'title' },
    { title: '开始时间', dataIndex: 'start_time', key: 'start', render: (v: string) => v ? new Date(v).toLocaleString() : '-' },
    { title: '结束时间', dataIndex: 'end_time', key: 'end', render: (v: string) => v ? new Date(v).toLocaleString() : '-' },
    { title: '报名数', dataIndex: 'current_count', key: 'count', width: 80 },
    {
      title: '状态', dataIndex: 'status', key: 'status', width: 100,
      render: (s: string) => {
        const labels: Record<string, string> = { draft: '草稿', published: '已发布', ended: '已结束' };
        const colors: Record<string, string> = { draft: 'default', published: 'green', ended: 'blue' };
        return <Tag color={colors[s] || 'default'}>{labels[s] || s}</Tag>;
      },
    },
    {
      title: '操作', key: 'action', width: 100,
      render: () => <Button type="link">查看</Button>,
    },
  ];

  const handleSave = () => {
    message.success('活动已保存（API对接中）');
    setShowEditor(false);
  };

  const handleApplyAICopy = (copy: string) => {
    setEditorContent(prev => ({ ...prev, description: copy }));
    message.success('AI文案已填入');
  };

  return (
    <AppLayout>
      <Row gutter={24}>
        <Col span={16}>
          <Card
            title="活动管理"
            extra={<Button type="primary" icon={<PlusOutlined />} onClick={() => setShowEditor(true)}>新建活动</Button>}
          >
            <Table columns={columns} dataSource={activities} rowKey="id" locale={{ emptyText: '暂无活动，点击"新建活动"开始' }} />
          </Card>
        </Col>
        <Col span={8}>
          <AICopyPanel mode="activity" />
        </Col>
      </Row>

      <Modal
        title="新建活动"
        open={showEditor}
        onCancel={() => setShowEditor(false)}
        onOk={handleSave}
        width={800}
        okText="保存"
        cancelText="取消"
      >
        <Space direction="vertical" style={{ width: '100%' }}>
          <Input
            placeholder="活动标题"
            value={editorContent.title}
            onChange={e => setEditorContent({ ...editorContent, title: e.target.value })}
            size="large"
          />
          <RangePicker showTime style={{ width: '100%' }} placeholder={['开始时间', '结束时间']} />
          <Input
            placeholder="活动地点"
            value={editorContent.location}
            onChange={e => setEditorContent({ ...editorContent, location: e.target.value })}
          />
          <TextArea
            placeholder="活动详情（可使用右侧AI助手生成文案）"
            value={editorContent.description}
            onChange={e => setEditorContent({ ...editorContent, description: e.target.value })}
            rows={8}
          />
        </Space>
      </Modal>
    </AppLayout>
  );
}
