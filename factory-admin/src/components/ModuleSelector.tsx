import { Card, Row, Col, Tag, Checkbox } from 'antd';
import {
  FileTextOutlined, RobotOutlined, ShoppingCartOutlined,
  GiftOutlined, CalendarOutlined, CrownOutlined, BarChartOutlined, AppstoreOutlined,
} from '@ant-design/icons';

const MODULES = [
  { key: 'base', name: '基础框架', desc: '应用骨架、路由导航、BaaS SDK', icon: <AppstoreOutlined />, required: true, ai: false, deps: [] },
  { key: 'cms', name: '内容管理', desc: '文章发布、公告管理、轮播图', icon: <FileTextOutlined />, required: false, ai: false, deps: ['base'] },
  { key: 'ai-advisor', name: 'AI健康顾问', desc: '智能问答·症状评估·生长记录', icon: <RobotOutlined />, required: false, ai: true, deps: ['base'] },
  { key: 'shop', name: '电商商城', desc: '商品展示·购物车·订单管理', icon: <ShoppingCartOutlined />, required: false, ai: false, deps: ['base'] },
  { key: 'activity', name: '活动运营', desc: '活动发布·报名·AI文案·海报', icon: <GiftOutlined />, required: false, ai: true, deps: ['base'] },
  { key: 'booking', name: '预约服务', desc: '时段预约·顾问管理', icon: <CalendarOutlined />, required: false, ai: false, deps: ['base'] },
  { key: 'member', name: '会员体系', desc: '会员等级·积分·优惠券', icon: <CrownOutlined />, required: false, ai: false, deps: ['base'] },
  { key: 'analytics', name: '数据分析', desc: '用户分析·增长趋势', icon: <BarChartOutlined />, required: false, ai: false, deps: ['ai-advisor'] },
];

interface Props {
  value: string[];
  onChange: (modules: string[]) => void;
}

export default function ModuleSelector({ value, onChange }: Props) {
  const toggle = (key: string) => {
    if (key === 'base') return; // Cannot deselect base
    if (value.includes(key)) {
      onChange(value.filter(k => k !== key));
    } else {
      const mod = MODULES.find(m => m.key === key);
      const newModules = [...value, key];
      if (mod) {
        for (const dep of mod.deps) {
          if (!newModules.includes(dep)) {
            newModules.push(dep);
          }
        }
      }
      onChange([...new Set(newModules)]);
    }
  };

  return (
    <div>
      <h3 style={{ marginBottom: 16 }}>选择功能模块</h3>
      <p style={{ color: '#666', marginBottom: 24 }}>
        基础框架已默认选中。选择模块时自动补齐依赖。
      </p>
      <Row gutter={[16, 16]}>
        {MODULES.map(mod => {
          const selected = value.includes(mod.key);
          return (
            <Col key={mod.key} xs={24} sm={12} lg={6}>
              <Card
                hoverable
                onClick={() => toggle(mod.key)}
                style={{
                  border: selected ? '2px solid #4caf50' : '1px solid #f0f0f0',
                  cursor: mod.required ? 'default' : 'pointer',
                  opacity: mod.required ? 1 : (selected ? 1 : 0.85),
                }}
                styles={{ body: { padding: 16 } }}
              >
                <div style={{ fontSize: 28, marginBottom: 8, color: selected ? '#4caf50' : '#999' }}>
                  {mod.icon}
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 4 }}>
                  <strong>{mod.name}</strong>
                  {mod.ai && <Tag color="purple" style={{ fontSize: 10 }}>✨AI</Tag>}
                  {mod.required && <Tag color="blue" style={{ fontSize: 10 }}>必需</Tag>}
                </div>
                <p style={{ color: '#999', fontSize: 12, marginBottom: 8 }}>{mod.desc}</p>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <span style={{ fontSize: 11, color: '#bbb' }}>
                    依赖: {mod.deps.length > 0 ? mod.deps.join(', ') : '无'}
                  </span>
                  <Checkbox checked={selected} disabled={mod.required} />
                </div>
              </Card>
            </Col>
          );
        })}
      </Row>

      {/* Selected modules summary */}
      <Card size="small" style={{ marginTop: 16, background: '#f6ffed' }}>
        <strong>已选模块：</strong>
        {value.map((k, i) => (
          <Tag key={k} color="green" style={{ marginLeft: i > 0 ? 4 : 4 }}>{MODULES.find(m => m.key === k)?.name || k}</Tag>
        ))}
        <span style={{ marginLeft: 8, color: '#999' }}>共 {value.length} 个模块</span>
      </Card>
    </div>
  );
}
