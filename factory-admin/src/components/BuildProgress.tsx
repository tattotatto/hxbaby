import { Steps, Spin } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';

const BUILD_STAGES = [
  '模块解析',
  '模板拼装',
  '配置注入',
  '依赖合并',
  '打包压缩',
  '完成',
];

interface Props {
  status: string; // pending/resolving/composing/injecting/packaging/done/failed
}

// Map status to current stage index
function getStageIndex(status: string): number {
  const stages: Record<string, number> = {
    pending: 0,
    resolving: 0,
    composing: 1,
    injecting: 2,
    packaging: 4,
    done: 6,
    failed: -1,
  };
  return stages[status] ?? 0;
}

export default function BuildProgress({ status }: Props) {
  if (status === 'failed') {
    return (
      <div style={{ textAlign: 'center', padding: 24 }}>
        <div style={{ color: '#f44336', fontSize: 16, marginBottom: 8 }}>构建失败</div>
        <p style={{ color: '#999' }}>请查看错误日志或稍后重试</p>
      </div>
    );
  }

  if (status === 'done') {
    return (
      <div style={{ textAlign: 'center', padding: 24 }}>
        <div style={{ color: '#4caf50', fontSize: 16, marginBottom: 8 }}>构建完成</div>
        <p style={{ color: '#999' }}>源码包已生成，可以下载</p>
        <Steps
          current={6}
          size="small"
          items={BUILD_STAGES.map(s => ({ title: s }))}
          style={{ marginTop: 16 }}
        />
      </div>
    );
  }

  const current = getStageIndex(status);
  return (
    <div style={{ textAlign: 'center', padding: 24 }}>
      <Spin indicator={<LoadingOutlined style={{ fontSize: 48, color: '#4caf50' }} spin />} />
      <div style={{ marginTop: 16, color: '#666' }}>正在生成小程序源码...</div>
      <Steps
        current={current}
        status={status === 'pending' ? 'wait' : 'process'}
        size="small"
        items={BUILD_STAGES.map(s => ({ title: s }))}
        style={{ marginTop: 24 }}
      />
    </div>
  );
}
