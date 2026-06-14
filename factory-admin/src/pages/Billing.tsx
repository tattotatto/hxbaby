import { Card, Row, Col, Button, Tag, Typography, Space, List } from 'antd';
import { CheckCircleOutlined, CrownOutlined } from '@ant-design/icons';
import AppLayout from '../components/Layout';

const { Title, Text } = Typography;

const PLANS = [
  {
    name: '免费版',
    price: '¥0',
    period: '永久免费',
    color: '#4caf50',
    projects: 1,
    features: [
      '1个小程序项目',
      '8个功能模块可选',
      '基础BaaS API（3000次/月）',
      'AI健康顾问（100次/月）',
      'AI文案生成（20次/月）',
      '源码下载',
      '社区支持',
    ],
    highlighted: false,
  },
  {
    name: '基础版',
    price: '¥299',
    period: '/月',
    color: '#2196f3',
    projects: 3,
    features: [
      '3个小程序项目',
      '所有功能模块',
      'BaaS API（3万次/月）',
      'AI健康顾问（1000次/月）',
      'AI文案生成（200次/月）',
      'AI海报生成',
      '自定义品牌配置',
      '邮件支持',
    ],
    highlighted: true,
    tag: '推荐',
  },
  {
    name: '专业版',
    price: '¥999',
    period: '/月',
    color: '#ff9800',
    projects: 10,
    features: [
      '10个小程序项目',
      '所有功能模块',
      'BaaS API（15万次/月）',
      'AI健康顾问（5000次/月）',
      'AI文案生成（1000次/月）',
      'AI海报生成',
      '微信支付集成',
      '数据分析看板',
      '优先技术支持',
      '专属客户经理',
    ],
    highlighted: false,
    icon: <CrownOutlined />,
  },
  {
    name: '企业版',
    price: '联系我们',
    period: '',
    color: '#7c4dff',
    projects: -1,
    features: [
      '不限项目数',
      '所有功能',
      '不限API调用',
      '不限AI调用',
      '私有化部署可选',
      '定制化开发',
      'SLA保障 99.9%',
      '7×24小时专属支持',
      '培训与咨询',
    ],
    highlighted: false,
  },
];

export default function Billing() {
  return (
    <AppLayout>
      <div style={{ textAlign: 'center', marginBottom: 32 }}>
        <Title level={3}>选择适合你的套餐</Title>
        <Text type="secondary">所有套餐均包含8大功能模块和源码生成能力</Text>
      </div>
      <Row gutter={[16, 16]}>
        {PLANS.map((plan) => (
          <Col key={plan.name} xs={24} sm={12} lg={6}>
            <Card
              hoverable
              style={{
                border: plan.highlighted ? `2px solid ${plan.color}` : '1px solid #f0f0f0',
                height: '100%',
                position: 'relative',
              }}
              bodyStyle={{ padding: 24, display: 'flex', flexDirection: 'column', height: '100%' }}
            >
              {plan.tag && (
                <Tag
                  color={plan.color}
                  style={{ position: 'absolute', top: -12, right: 16, fontSize: 13, padding: '2px 12px' }}
                >
                  {plan.tag}
                </Tag>
              )}
              {plan.icon && <div style={{ fontSize: 28, color: plan.color, marginBottom: 8 }}>{plan.icon}</div>}
              <Title level={4}>{plan.name}</Title>
              <div style={{ marginBottom: 16 }}>
                <Text style={{ fontSize: 28, fontWeight: 700, color: plan.color }}>{plan.price}</Text>
                {plan.period && <Text type="secondary">{plan.period}</Text>}
              </div>
              <List
                size="small"
                dataSource={plan.features}
                renderItem={(item: string) => (
                  <List.Item style={{ border: 'none', padding: '4px 0' }}>
                    <Space>
                      <CheckCircleOutlined style={{ color: plan.color, fontSize: 12 }} />
                      <Text style={{ fontSize: 13 }}>{item}</Text>
                    </Space>
                  </List.Item>
                )}
                style={{ flex: 1, marginBottom: 16 }}
              />
              <Button
                type={plan.highlighted ? 'primary' : 'default'}
                block
                size="large"
                style={{
                  background: plan.highlighted ? plan.color : undefined,
                  borderColor: plan.highlighted ? plan.color : undefined,
                }}
              >
                {plan.price === '联系我们' ? '联系销售' : '立即开通'}
              </Button>
            </Card>
          </Col>
        ))}
      </Row>

      <Card style={{ marginTop: 32, textAlign: 'center', background: '#fafafa' }}>
        <Title level={5}>常见问题</Title>
        <Space direction="vertical" style={{ maxWidth: 600, margin: '0 auto', textAlign: 'left' }}>
          <div>
            <Text strong>Q: 可以随时升级套餐吗？</Text>
            <br /><Text type="secondary">A: 可以，升级后立即生效，差价按剩余天数折算。</Text>
          </div>
          <div>
            <Text strong>Q: 免费版有什么限制？</Text>
            <br /><Text type="secondary">A: 免费版限制1个项目，API和AI调用有月度限额，适合个人体验和小规模使用。</Text>
          </div>
          <div>
            <Text strong>Q: 生成的源码可以商用吗？</Text>
            <br /><Text type="secondary">A: 可以，生成的源码归您所有，可自由修改和商用。</Text>
          </div>
        </Space>
      </Card>
    </AppLayout>
  );
}
