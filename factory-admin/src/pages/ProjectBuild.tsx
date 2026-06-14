import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { Card, Button, message, Descriptions, Space, Table, Tag, Typography } from 'antd';
import { DownloadOutlined, BuildOutlined } from '@ant-design/icons';
import { triggerBuild, getBuildStatus, getBuildHistory, getDownloadUrl, type BuildTask } from '../api/build';
import { getProject, type Project } from '../api/project';
import BuildProgress from '../components/BuildProgress';
import AppLayout from '../components/Layout';

const { Text } = Typography;

export default function ProjectBuild() {
  const { id } = useParams<{ id: string }>();
  const projectId = parseInt(id || '0');
  const [project, setProject] = useState<Project | null>(null);
  const [buildTask, setBuildTask] = useState<BuildTask | null>(null);
  const [history, setHistory] = useState<BuildTask[]>([]);
  const [polling, setPolling] = useState(false);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadProject();
    loadHistory();
  }, [projectId]);

  const loadProject = async () => {
    try {
      const p = await getProject(projectId);
      setProject(p);
    } catch { message.error('加载项目信息失败'); }
  };

  const loadHistory = async () => {
    try {
      const h = await getBuildHistory(projectId);
      setHistory(h);
    } catch {}
  };

  // Poll build status
  useEffect(() => {
    if (!buildTask || buildTask.status === 'done' || buildTask.status === 'failed') {
      setPolling(false);
      return;
    }
    setPolling(true);
    const timer = setInterval(async () => {
      try {
        const task = await getBuildStatus(buildTask.id);
        setBuildTask(task);
        if (task.status === 'done' || task.status === 'failed') {
          setPolling(false);
          if (task.status === 'done') message.success('构建完成！');
          loadHistory();
        }
      } catch {}
    }, 2000);
    return () => clearInterval(timer);
  }, [buildTask?.id, buildTask?.status]);

  const handleBuild = async () => {
    setLoading(true);
    try {
      const result = await triggerBuild(projectId);
      const task = await getBuildStatus(result.build_id);
      setBuildTask(task);
      message.success('构建已启动');
    } catch (err: any) {
      message.error(err.response?.data?.message || '构建失败');
    } finally {
      setLoading(false);
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
    {
      title: '状态', dataIndex: 'status', key: 'status', width: 100,
      render: (s: string) => {
        const colors: Record<string, string> = { done: 'green', failed: 'red', pending: 'blue' };
        const labels: Record<string, string> = { done: '完成', failed: '失败', pending: '处理中' };
        return <Tag color={colors[s] || 'default'}>{labels[s] || s}</Tag>;
      },
    },
    { title: '耗时', dataIndex: 'duration_ms', key: 'duration', width: 100, render: (v: number) => v ? `${(v / 1000).toFixed(1)}s` : '-' },
    { title: 'MD5', dataIndex: 'output_md5', key: 'md5', ellipsis: true },
    { title: '时间', dataIndex: 'created_at', key: 'time', render: (v: string) => new Date(v).toLocaleString() },
    {
      title: '操作', key: 'action', width: 100,
      render: (_: any, record: BuildTask) => record.status === 'done' ? (
        <Button type="link" href={getDownloadUrl(record.id)} icon={<DownloadOutlined />}>下载</Button>
      ) : null,
    },
  ];

  return (
    <AppLayout>
      <Card title="构建源码" style={{ marginBottom: 24 }}>
        {project && (
          <Descriptions column={3} size="small" style={{ marginBottom: 16 }}>
            <Descriptions.Item label="项目名称">{project.name}</Descriptions.Item>
            <Descriptions.Item label="状态">{project.status}</Descriptions.Item>
            <Descriptions.Item label="API Key">
              <Space>
                <Text copyable={{ text: project.api_key }}>{project.api_key.substring(0, 16)}...</Text>
              </Space>
            </Descriptions.Item>
          </Descriptions>
        )}

        <Button
          type="primary"
          size="large"
          icon={<BuildOutlined />}
          onClick={handleBuild}
          loading={loading}
          disabled={polling}
          style={{ marginBottom: 24 }}
        >
          {polling ? '构建中...' : '一键构建'}
        </Button>

        {buildTask && <BuildProgress status={buildTask.status} />}

        {buildTask?.status === 'done' && (
          <Card size="small" style={{ background: '#f6ffed', marginTop: 16 }}>
            <Space direction="vertical">
              <div><strong>MD5:</strong> {buildTask.output_md5}</div>
              <div><strong>耗时:</strong> {(buildTask.duration_ms / 1000).toFixed(1)}s</div>
              <Button type="primary" icon={<DownloadOutlined />} href={getDownloadUrl(buildTask.id)}>
                下载源码包
              </Button>
            </Space>
          </Card>
        )}
      </Card>

      {/* Build History */}
      <Card title="构建历史">
        <Table columns={columns} dataSource={history} rowKey="id" pagination={{ pageSize: 10 }} size="small" />
      </Card>
    </AppLayout>
  );
}
