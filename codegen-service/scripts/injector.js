/**
 * 配置注入器 — 将客户配置替换到模板占位符
 */
function injectConfig(template, config) {
  let result = template;
  for (const [key, value] of Object.entries(config)) {
    const regex = new RegExp(`\\{\\{${key}\\}\\}`, 'g');
    result = result.replace(regex, String(value ?? ''));
  }
  return result;
}

function buildInjectionConfig(project) {
  const brand = project.brand_config || {};
  return {
    appName: brand.appName || project.name || '健康助手',
    primaryColor: brand.primaryColor || '#4caf50',
    secondaryColor: brand.secondaryColor || '#ff9800',
    logo: brand.logo || '/static/logo.png',
    footer: brand.footer || `© ${new Date().getFullYear()} ${project.name || ''}`,
    baasBaseURL: process.env.BAAS_BASE_URL || 'http://localhost:8080',
    apiKey: project.api_key || 'YOUR_API_KEY',
    projectId: String(project.id || ''),
    wxAppId: project.wx_app_id || 'YOUR_WX_APP_ID',
  };
}

module.exports = { injectConfig, buildInjectionConfig };
