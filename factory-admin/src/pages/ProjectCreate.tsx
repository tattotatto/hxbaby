import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Steps, Form, Input, Button, Card, message } from 'antd';
import ModuleSelector from '../components/ModuleSelector';
import BrandConfig from '../components/BrandConfig';
import AppLayout from '../components/Layout';
import { createProject, updateProject } from '../api/project';

const { TextArea } = Input;

export default function ProjectCreate() {
  const navigate = useNavigate();
  const [current, setCurrent] = useState(0);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    wx_app_id: '',
    modules: ['base'],
    brand_config: {
      appName: '',
      primaryColor: '#4caf50',
      secondaryColor: '#ff9800',
      logo: '',
      footer: '',
    },
  });
  const [createdId, setCreatedId] = useState<number | null>(null);
  const [loading, setLoading] = useState(false);

  const steps = [
    { title: '基本信息' },
    { title: '选择模块' },
    { title: '品牌配置' },
  ];

  const handleNext = () => setCurrent(current + 1);
  const handlePrev = () => setCurrent(current - 1);

  const handleFinish = async () => {
    setLoading(true);
    try {
      if (createdId) {
        // Update existing project with modules and brand
        await updateProject(createdId, {
          modules: formData.modules,
          brand_config: formData.brand_config,
        });
        message.success('项目配置已保存！');
        navigate(`/projects/${createdId}/build`);
      } else {
        // Create new project
        const project = await createProject({
          name: formData.name,
          description: formData.description,
          modules: formData.modules,
        });
        setCreatedId(project.id);
        // Update brand config
        await updateProject(project.id, {
          modules: formData.modules,
          brand_config: formData.brand_config,
        });
        message.success('项目创建成功！');
        navigate(`/projects/${project.id}/build`);
      }
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { message?: string } } })?.response?.data?.message || '操作失败';
      message.error(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <AppLayout>
      <Card>
        <Steps current={current} items={steps} style={{ marginBottom: 32 }} />

        {current === 0 && (
          <Form layout="vertical" style={{ maxWidth: 600 }}>
            <Form.Item label="小程序名称" required>
              <Input
                placeholder="例如：宝宝健康助手"
                value={formData.name}
                onChange={(e) => {
                  setFormData({ ...formData, name: e.target.value, brand_config: { ...formData.brand_config, appName: e.target.value } });
                }}
                maxLength={100}
              />
            </Form.Item>
            <Form.Item label="一句话描述">
              <TextArea
                placeholder="简单描述你的小程序用途"
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                rows={3}
                maxLength={200}
              />
            </Form.Item>
            <Form.Item label="微信AppID" help="可在微信开放平台获取，也可后续填写">
              <Input
                placeholder="wx..."
                value={formData.wx_app_id}
                onChange={(e) => setFormData({ ...formData, wx_app_id: e.target.value })}
              />
            </Form.Item>
            <Form.Item>
              <Button type="primary" onClick={handleNext} disabled={!formData.name.trim()}>
                下一步
              </Button>
            </Form.Item>
          </Form>
        )}

        {current === 1 && (
          <ModuleSelector
            value={formData.modules}
            onChange={(modules) => setFormData({ ...formData, modules })}
          />
        )}

        {current === 2 && (
          <BrandConfig
            value={formData.brand_config}
            onChange={(config) => setFormData({ ...formData, brand_config: { ...formData.brand_config, ...config } })}
          />
        )}

        {current > 0 && (
          <div style={{ marginTop: 24, display: 'flex', justifyContent: 'space-between' }}>
            <Button onClick={handlePrev}>上一步</Button>
            {current < 2 ? (
              <Button type="primary" onClick={handleNext}>下一步</Button>
            ) : (
              <Button type="primary" onClick={handleFinish} loading={loading}>
                完成创建
              </Button>
            )}
          </div>
        )}
      </Card>
    </AppLayout>
  );
}
