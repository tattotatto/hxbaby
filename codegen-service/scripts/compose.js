const fs = require('fs');
const path = require('path');
const { resolveModules, detectConflicts } = require('./resolver');
const { injectConfig, buildInjectionConfig } = require('./injector');

const TEMPLATES_DIR = path.join(__dirname, '..', 'templates');

async function compose(project) {
  const warnings = [];
  const modules = resolveModules(project.modules || []);
  const conflictWarnings = detectConflicts(modules);
  warnings.push(...conflictWarnings);

  const outputDir = path.join(__dirname, '..', 'output', project.build_task_id || 'latest');
  fs.mkdirSync(outputDir, { recursive: true });

  const config = buildInjectionConfig(project);

  // 1. 复制基础骨架
  const baseDir = path.join(TEMPLATES_DIR, 'base');
  copyAndInjectDir(baseDir, outputDir, config);

  // 2. 合并路由
  const pagesJSON = mergePages(modules);
  fs.writeFileSync(
    path.join(outputDir, 'pages.json'),
    JSON.stringify(pagesJSON, null, 2)
  );

  // 3. 复制各模块
  for (const mod of modules) {
    if (mod === 'base') continue;
    const modDir = path.join(TEMPLATES_DIR, 'modules', mod);
    if (!fs.existsSync(modDir)) {
      warnings.push(`模块 "${mod}" 模板不存在，已跳过`);
      continue;
    }
    const destDir = path.join(outputDir, 'src');
    copyAndInjectDir(modDir, destDir, config);
  }

  // 4. package.json
  const pkg = { name: (project.name || 'miniapp').replace(/\s+/g, '-').toLowerCase(), version: '1.0.0', scripts: { 'dev:mp-weixin': 'uni -p mp-weixin', 'build:mp-weixin': 'uni build -p mp-weixin' }, dependencies: { 'vue': '^3.4.0', 'uni-app': '^3.0.0', 'pinia': '^2.1.0' } };
  fs.writeFileSync(path.join(outputDir, 'package.json'), JSON.stringify(pkg, null, 2));

  // 5. README
  const readme = `# ${project.name || '小程序'}\n\n## 模块\n${modules.map(m => `- ${m}`).join('\n')}\n\n## 快速开始\n1. npm install\n2. 填写 manifest.json 中的微信 AppID\n3. npm run dev:mp-weixin\n4. 微信开发者工具导入 dist/dev/mp-weixin\n\n## BaaS配置\n- 地址: ${config.baasBaseURL}\n- API Key: ${config.apiKey}\n${warnings.length ? '\n## 注意事项\n' + warnings.map(w => `- ${w}`).join('\n') : ''}`;
  fs.writeFileSync(path.join(outputDir, 'README.md'), readme);

  return { outputDir, warnings };
}

function mergePages(modules) {
  const pages = [];
  for (const mod of modules) {
    const mf = path.join(TEMPLATES_DIR, 'modules', mod, 'module.json');
    if (!fs.existsSync(mf)) continue;
    const cfg = JSON.parse(fs.readFileSync(mf, 'utf-8'));
    if (cfg.pages) {
      for (const p of cfg.pages) pages.push({ path: p.path, style: p.style || {} });
    }
  }
  const tabBar = buildTabBar(modules);
  return { pages, globalStyle: { navigationBarTextStyle: 'black', navigationBarTitleText: '{{appName}}', navigationBarBackgroundColor: '#ffffff' }, ...(tabBar ? { tabBar } : {}) };
}

function buildTabBar(modules) {
  const items = [];
  const order = ['cms', 'ai-advisor', 'shop', 'activity', 'member'];
  for (const mod of order) {
    if (!modules.includes(mod)) continue;
    const mf = path.join(TEMPLATES_DIR, 'modules', mod, 'module.json');
    if (!fs.existsSync(mf)) continue;
    const cfg = JSON.parse(fs.readFileSync(mf, 'utf-8'));
    if (cfg.tabbar) items.push(cfg.tabbar);
  }
  if (items.length === 0) return null;
  return { list: items.slice(0, 5), color: '#999999', selectedColor: '{{primaryColor}}' };
}

function copyAndInjectDir(srcDir, destDir, config) {
  if (!fs.existsSync(srcDir)) return;
  for (const entry of fs.readdirSync(srcDir, { withFileTypes: true })) {
    const src = path.join(srcDir, entry.name);
    const dest = path.join(destDir, entry.name.replace('.hbs', ''));
    if (entry.isDirectory()) {
      fs.mkdirSync(dest, { recursive: true });
      copyAndInjectDir(src, dest, config);
    } else {
      let content = fs.readFileSync(src, 'utf-8');
      if (entry.name.endsWith('.hbs')) content = injectConfig(content, config);
      fs.writeFileSync(dest, content);
    }
  }
}

module.exports = { compose };
