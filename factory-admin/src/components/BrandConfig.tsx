import { Form, Input, Row, Col, Card } from 'antd';

interface Props {
  value: {
    appName: string;
    primaryColor: string;
    secondaryColor: string;
    logo: string;
    footer: string;
  };
  onChange: (config: Record<string, string>) => void;
}

const PRESET_COLORS = ['#4caf50', '#2196f3', '#ff9800', '#e91e63', '#9c27b0', '#00bcd4', '#ff5722', '#607d8b'];

export default function BrandConfig({ value, onChange }: Props) {
  const update = (key: string, val: string) => {
    onChange({ ...value, [key]: val });
  };

  return (
    <Row gutter={24}>
      <Col span={14}>
        <h3 style={{ marginBottom: 16 }}>品牌配置</h3>
        <Form layout="vertical">
          <Form.Item label="应用名称">
            <Input value={value.appName} onChange={(e) => update('appName', e.target.value)} placeholder="显示在导航栏和启动页" />
          </Form.Item>
          <Form.Item label="主题色">
            <div style={{ display: 'flex', gap: 8, marginBottom: 8, flexWrap: 'wrap' }}>
              {PRESET_COLORS.map(color => (
                <div
                  key={color}
                  onClick={() => update('primaryColor', color)}
                  style={{
                    width: 32, height: 32, borderRadius: 6, background: color,
                    cursor: 'pointer', border: value.primaryColor === color ? '3px solid #333' : '1px solid #ddd',
                  }}
                />
              ))}
            </div>
            <Input value={value.primaryColor} onChange={(e) => update('primaryColor', e.target.value)} placeholder="#4caf50" style={{ width: 120 }} />
          </Form.Item>
          <Form.Item label="辅助色">
            <Input value={value.secondaryColor} onChange={(e) => update('secondaryColor', e.target.value)} placeholder="#ff9800" style={{ width: 120 }} />
          </Form.Item>
          <Form.Item label="Logo URL">
            <Input value={value.logo} onChange={(e) => update('logo', e.target.value)} placeholder="https://example.com/logo.png" />
            {value.logo && <img src={value.logo} alt="Logo preview" style={{ marginTop: 8, maxHeight: 64 }} />}
          </Form.Item>
          <Form.Item label="底部版权信息">
            <Input value={value.footer} onChange={(e) => update('footer', e.target.value)} placeholder={`© ${new Date().getFullYear()} 版权所有`} />
          </Form.Item>
        </Form>
      </Col>
      <Col span={10}>
        {/* Phone mockup preview */}
        <Card title="手机预览" size="small" style={{ position: 'sticky', top: 16 }}>
          <div style={{
            width: 280, margin: '0 auto', border: '2px solid #333', borderRadius: 24,
            padding: '20px 12px', background: '#fff',
          }}>
            {/* Status bar */}
            <div style={{ textAlign: 'center', fontSize: 12, color: '#333', marginBottom: 8 }}>
              {value.appName || '应用名称'}
            </div>
            {/* Content preview */}
            <div style={{ background: '#f5f5f5', borderRadius: 8, padding: 16, minHeight: 200 }}>
              <div style={{
                background: value.primaryColor, color: '#fff', padding: '8px 12px',
                borderRadius: 6, marginBottom: 12, textAlign: 'center', fontSize: 13,
              }}>
                {value.appName || '应用名称'}
              </div>
              <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap' }}>
                {[value.primaryColor, value.secondaryColor].map((c, i) => (
                  <div key={i} style={{ width: 40, height: 40, background: c, borderRadius: 8 }} />
                ))}
              </div>
            </div>
            {/* Bottom bar */}
            <div style={{ textAlign: 'center', fontSize: 10, color: '#999', marginTop: 12 }}>
              {value.footer || '© 版权所有'}
            </div>
          </div>
        </Card>
      </Col>
    </Row>
  );
}
