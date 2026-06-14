import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Row, Col, Card, Tag, Button, Empty, Typography, Spin, message } from 'antd';
import { PlusOutlined, AppstoreOutlined } from '@ant-design/icons';
import AppLayout from '../components/Layout';
import { getProjects, type Project } from '../api/project';

const { Title, Text, Paragraph } = Typography;

const statusColorMap: Record<string, string> = {
  active: 'green',
  inactive: 'default',
  pending: 'orange',
  archived: 'red',
};

export default function DashboardPage() {
  const navigate = useNavigate();
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchProjects();
  }, []);

  const fetchProjects = async () => {
    setLoading(true);
    try {
      const data = await getProjects();
      setProjects(data);
    } catch (err: any) {
      const msg = err?.response?.data?.message ?? '获取项目列表失败';
      message.error(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <AppLayout>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Title level={4} style={{ margin: 0 }}>
          <AppstoreOutlined style={{ marginRight: 8 }} />
          我的小程序
        </Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/projects/new')}>
          创建小程序
        </Button>
      </div>

      <Spin spinning={loading}>
        {!loading && projects.length === 0 ? (
          <Empty
            description="还没有项目，点击创建第一个小程序"
            style={{ padding: 80 }}
          >
            <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/projects/new')}>
              创建小程序
            </Button>
          </Empty>
        ) : (
          <Row gutter={[16, 16]}>
            {projects.map((project) => (
              <Col key={project.id} xs={24} sm={12} lg={8} xl={6}>
                <Card
                  hoverable
                  onClick={() => navigate(`/projects/${project.id}/build`)}
                  title={
                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                      <AppstoreOutlined />
                      <Text strong ellipsis style={{ maxWidth: 140 }}>
                        {project.name}
                      </Text>
                    </div>
                  }
                  extra={
                    <Tag color={statusColorMap[project.status] ?? 'default'}>
                      {project.status}
                    </Tag>
                  }
                  style={{ height: '100%' }}
                >
                  <Paragraph
                    type="secondary"
                    ellipsis={{ rows: 2 }}
                    style={{ minHeight: 44, marginBottom: 12 }}
                  >
                    {project.description || '暂无描述'}
                  </Paragraph>
                  <div style={{ display: 'flex', flexWrap: 'wrap', gap: 4, marginBottom: 12 }}>
                    {project.modules?.map((mod) => (
                      <Tag key={mod} color="blue">
                        {mod}
                      </Tag>
                    ))}
                  </div>
                  <Text type="secondary" style={{ fontSize: 12 }}>
                    创建于 {new Date(project.created_at).toLocaleDateString('zh-CN')}
                  </Text>
                </Card>
              </Col>
            ))}
          </Row>
        )}
      </Spin>
    </AppLayout>
  );
}
