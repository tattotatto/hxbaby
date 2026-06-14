// 模块依赖图
const DEPENDENCY_GRAPH = {
  'base':       [],
  'cms':        ['base'],
  'ai-advisor': ['base'],
  'shop':       ['base'],
  'activity':   ['base'],
  'booking':    ['base'],
  'member':     ['base'],
  'analytics':  ['ai-advisor'],
};

const WARNING_RULES = [
  {
    condition: (mods) => mods.includes('shop'),
    message: '电商模块需要配置微信支付商户号',
  },
];

/**
 * 解析模块列表：补齐依赖 + 拓扑排序
 */
function resolveModules(selected) {
  if (!selected.includes('base')) {
    selected = ['base', ...selected];
  }

  const resolved = new Set();
  const queue = [...selected];

  while (queue.length > 0) {
    const mod = queue.shift();
    if (resolved.has(mod)) continue;
    resolved.add(mod);

    const deps = DEPENDENCY_GRAPH[mod] || [];
    for (const dep of deps) {
      if (!resolved.has(dep)) queue.push(dep);
    }
  }

  return topologicalSort([...resolved]);
}

function topologicalSort(modules) {
  const sorted = [];
  const visited = new Set();
  const temp = new Set();

  function visit(mod) {
    if (temp.has(mod) || visited.has(mod)) return;
    temp.add(mod);
    for (const dep of (DEPENDENCY_GRAPH[mod] || [])) {
      if (!visited.has(dep)) visit(dep);
    }
    temp.delete(mod);
    visited.add(mod);
    sorted.push(mod);
  }

  for (const mod of modules) visit(mod);
  return sorted;
}

function detectConflicts(modules) {
  const warnings = [];
  for (const rule of WARNING_RULES) {
    if (rule.condition(modules)) warnings.push(rule.message);
  }
  return warnings;
}

module.exports = { resolveModules, detectConflicts, DEPENDENCY_GRAPH };
