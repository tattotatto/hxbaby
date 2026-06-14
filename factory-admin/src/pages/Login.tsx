import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Form, Input, Button, Card, Typography, message, Space } from 'antd';
import { PhoneOutlined, LockOutlined, ShopOutlined } from '@ant-design/icons';
import { login } from '../api/auth';
import { useAuth } from '../hooks/useAuth';

const { Title, Text } = Typography;

interface LoginFormValues {
  phone: string;
  password: string;
}

export default function LoginPage() {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const auth = useAuth();

  const onFinish = async (values: LoginFormValues) => {
    setLoading(true);
    try {
      const result = await login(values.phone, values.password);
      auth.login(result.customer, result.token);
      message.success('登录成功');
      navigate('/dashboard');
    } catch (err: any) {
      const msg = err?.response?.data?.message ?? '登录失败，请检查手机号和密码';
      message.error(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%)',
      }}
    >
      <Card style={{ width: 400, boxShadow: '0 4px 24px rgba(0,0,0,0.1)' }}>
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          <div style={{ textAlign: 'center' }}>
            <ShopOutlined style={{ fontSize: 48, color: '#4caf50', marginBottom: 16 }} />
            <Title level={3} style={{ margin: 0 }}>
              小程序工厂
            </Title>
            <Text type="secondary">登录您的管理后台</Text>
          </div>

          <Form<LoginFormValues>
            name="login"
            onFinish={onFinish}
            layout="vertical"
            size="large"
            autoComplete="off"
          >
            <Form.Item
              name="phone"
              rules={[
                { required: true, message: '请输入手机号' },
                { pattern: /^1\d{10}$/, message: '请输入正确的11位手机号' },
              ]}
            >
              <Input prefix={<PhoneOutlined />} placeholder="手机号" maxLength={11} />
            </Form.Item>

            <Form.Item
              name="password"
              rules={[
                { required: true, message: '请输入密码' },
                { min: 6, message: '密码至少6位' },
              ]}
            >
              <Input.Password prefix={<LockOutlined />} placeholder="密码" />
            </Form.Item>

            <Form.Item>
              <Button type="primary" htmlType="submit" loading={loading} block>
                登录
              </Button>
            </Form.Item>
          </Form>

          <div style={{ textAlign: 'center' }}>
            <Text>
              还没有账号？
              <Link to="/register">立即注册</Link>
            </Text>
          </div>
        </Space>
      </Card>
    </div>
  );
}
