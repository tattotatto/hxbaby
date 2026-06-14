import { Card, Steps, Typography, Alert } from 'antd';
import AppLayout from '../components/Layout';

const { Text, Paragraph } = Typography;

export default function ProjectGuide() {
  return (
    <AppLayout>
      <Card title="接入指引">
        <Steps
          direction="vertical"
          current={-1}
          items={[
            {
              title: '1. 下载源码并解压',
              description: '在"构建源码"页面点击"一键构建"，等待构建完成后下载ZIP包并解压到本地目录。',
            },
            {
              title: '2. 安装依赖',
              description: (
                <div>
                  <Text>打开终端，进入解压后的目录，运行：</Text>
                  <pre style={{ background: '#f5f5f5', padding: '8px 12px', borderRadius: 4, marginTop: 8 }}>
                    npm install
                  </pre>
                </div>
              ),
            },
            {
              title: '3. 配置微信AppID',
              description: (
                <div>
                  <Text>在 </Text>
                  <Text code>src/manifest.json</Text>
                  <Text> 中填写你的微信小程序AppID。在微信开放平台 开发管理 开发设置中获取。</Text>
                </div>
              ),
            },
            {
              title: '4. 打开微信开发者工具',
              description: (
                <div>
                  <Text>下载并安装</Text>
                  <Text strong> 微信开发者工具</Text>
                  <Text>，选择"导入项目"，选择解压后的目录，填写AppID，点击"确定"即可预览。</Text>
                </div>
              ),
            },
          ]}
        />

        <Alert
          type="info"
          message="BaaS API 配置"
          description={
            <div>
              <Paragraph>
                API Key 已自动内置在源码中，无需手动配置。
                如需查看或修改 API 地址，请编辑 <Text code>src/config.js</Text> 文件。
              </Paragraph>
              <Paragraph>
                BaaS API 文档：<a href="https://api.hxbaby.com/docs" target="_blank" rel="noreferrer">https://api.hxbaby.com/docs</a>
              </Paragraph>
            </div>
          }
          style={{ marginTop: 24 }}
        />

        <Alert
          type="warning"
          message="注意事项"
          description={
            <ul style={{ margin: 0, paddingLeft: 20 }}>
              <li>生成的小程序基于 UniApp (Vue 3)，如需使用其他平台（支付宝、抖音等），请参考 UniApp 官方文档</li>
              <li>API Key 请妥善保管，泄露后可在后台重新生成</li>
              <li>修改模块配置后需要重新构建</li>
            </ul>
          }
          style={{ marginTop: 16 }}
        />
      </Card>
    </AppLayout>
  );
}
