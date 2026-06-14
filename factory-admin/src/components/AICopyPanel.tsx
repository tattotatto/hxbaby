import { useState } from 'react';
import { Card, Input, Button, Select, Space, Tabs, Spin, message, Typography } from 'antd';
import { RobotOutlined, CopyOutlined, LoadingOutlined } from '@ant-design/icons';
import { generateArticle, generateSummary, generateActivityCopy } from '../api/ai';

const { TextArea } = Input;
const { Text, Paragraph } = Typography;

interface Props {
  mode: 'article' | 'activity';
}

export default function AICopyPanel({ mode }: Props) {
  const [topic, setTopic] = useState('');
  const [category, setCategory] = useState('');
  const [generating, setGenerating] = useState(false);
  const [result, setResult] = useState<Record<string, string>>({});
  const [activeTab, setActiveTab] = useState('article');

  const handleGenerateArticle = async () => {
    if (!topic.trim()) { message.warning('请输入主题'); return; }
    setGenerating(true);
    try {
      const res = await generateArticle(topic, category);
      setResult(prev => ({ ...prev, article: res.content }));
      setActiveTab('article');
      message.success('文章生成成功');
    } catch { message.error('AI生成失败，请重试'); }
    finally { setGenerating(false); }
  };

  const handleGenerateSummary = async () => {
    const content = result.article;
    if (!content) { message.warning('请先生成文章'); return; }
    setGenerating(true);
    try {
      const res = await generateSummary(content);
      setResult(prev => ({ ...prev, summary: res.content }));
      setActiveTab('summary');
    } catch { message.error('摘要生成失败'); }
    finally { setGenerating(false); }
  };

  const handleGenerateCopy = async () => {
    if (!topic.trim()) { message.warning('请输入活动标题'); return; }
    setGenerating(true);
    try {
      const res = await generateActivityCopy(topic, category);
      setResult(prev => ({ ...prev, copy: res.content }));
      setActiveTab('copy');
    } catch { message.error('文案生成失败'); }
    finally { setGenerating(false); }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text).then(() => message.success('已复制到剪贴板'));
  };

  return (
    <Card
      title={<Space><RobotOutlined />AI 写作助手</Space>}
      size="small"
      style={{ background: '#fafafa' }}
    >
      <Space direction="vertical" style={{ width: '100%' }}>
        <Input
          placeholder={mode === 'article' ? '输入文章主题，如：婴儿辅食添加指南' : '输入活动标题'}
          value={topic}
          onChange={e => setTopic(e.target.value)}
          onPressEnter={mode === 'article' ? handleGenerateArticle : handleGenerateCopy}
        />
        {mode === 'article' && (
          <Select
            placeholder="选择分类（可选）"
            value={category || undefined}
            onChange={setCategory}
            allowClear
            style={{ width: '100%' }}
            options={[
              { value: '健康资讯', label: '健康资讯' },
              { value: '育儿知识', label: '育儿知识' },
              { value: '营养辅食', label: '营养辅食' },
              { value: '成长发育', label: '成长发育' },
            ]}
          />
        )}
        <Space>
          {mode === 'article' && (
            <>
              <Button type="primary" onClick={handleGenerateArticle} loading={generating} icon={<RobotOutlined />}>
                生成文章
              </Button>
              <Button onClick={handleGenerateSummary} loading={generating} disabled={!result.article}>
                生成摘要
              </Button>
            </>
          )}
          {mode === 'activity' && (
            <Button type="primary" onClick={handleGenerateCopy} loading={generating} icon={<RobotOutlined />}>
              生成文案
            </Button>
          )}
        </Space>
      </Space>

      {generating && (
        <div style={{ textAlign: 'center', padding: '24px 0' }}>
          <Spin indicator={<LoadingOutlined style={{ fontSize: 32, color: '#7c4dff' }} spin />} />
          <div style={{ marginTop: 8, color: '#999' }}>AI 正在创作中...</div>
        </div>
      )}

      {!generating && Object.keys(result).length > 0 && (
        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          style={{ marginTop: 16 }}
          items={[
            result.article && {
              key: 'article',
              label: '文章',
              children: (
                <div>
                  <div style={{ maxHeight: 400, overflow: 'auto', whiteSpace: 'pre-wrap', background: '#fff', padding: 12, borderRadius: 4, border: '1px solid #f0f0f0', fontSize: 13 }}>
                    {result.article}
                  </div>
                  <Button size="small" icon={<CopyOutlined />} onClick={() => copyToClipboard(result.article)} style={{ marginTop: 8 }}>
                    复制全文
                  </Button>
                </div>
              ),
            },
            result.summary && {
              key: 'summary',
              label: '摘要',
              children: (
                <div>
                  <Paragraph style={{ background: '#fff', padding: 12, borderRadius: 4 }}>{result.summary}</Paragraph>
                  <Button size="small" icon={<CopyOutlined />} onClick={() => copyToClipboard(result.summary)}>
                    复制摘要
                  </Button>
                </div>
              ),
            },
            result.copy && {
              key: 'copy',
              label: '活动文案',
              children: (
                <div>
                  <div style={{ maxHeight: 400, overflow: 'auto', whiteSpace: 'pre-wrap', background: '#fff', padding: 12, borderRadius: 4, border: '1px solid #f0f0f0', fontSize: 13 }}>
                    {result.copy}
                  </div>
                  <Button size="small" icon={<CopyOutlined />} onClick={() => copyToClipboard(result.copy)} style={{ marginTop: 8 }}>
                    复制文案
                  </Button>
                </div>
              ),
            },
          ].filter(Boolean) as any}
        />
      )}
    </Card>
  );
}
