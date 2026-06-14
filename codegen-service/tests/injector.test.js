const { injectConfig, buildInjectionConfig } = require('../scripts/injector');
const { test } = require('node:test');
const assert = require('node:assert');

test('injectConfig replaces variables', () => {
  const tpl = 'const APP = "{{appName}}"; const COLOR = "{{primaryColor}}";';
  const result = injectConfig(tpl, { appName: '测试助手', primaryColor: '#ff0000' });
  assert(result.includes('测试助手'));
  assert(result.includes('#ff0000'));
  assert(!result.includes('{{'));
});

test('injectConfig leaves unmatched variables', () => {
  const result = injectConfig('{{appName}} {{unknown}}', { appName: 'Test' });
  assert(result.includes('Test'));
  assert(result.includes('{{unknown}}'));
});

test('buildInjectionConfig includes all required keys', () => {
  const project = { id: 123, name: '测试', api_key: 'ak_test', brand_config: {} };
  const config = buildInjectionConfig(project);
  assert(config.appName);
  assert(config.apiKey);
  assert.strictEqual(config.projectId, '123');
});

console.log('All injector tests passed');
