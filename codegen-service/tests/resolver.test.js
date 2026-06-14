const { resolveModules, detectConflicts } = require('../scripts/resolver');
const { test } = require('node:test');
const assert = require('node:assert');

test('resolveModules adds missing base', () => {
  const result = resolveModules(['cms']);
  assert(result.includes('base'));
  assert(result.includes('cms'));
});

test('resolveModules ensures dependencies first', () => {
  const result = resolveModules(['shop', 'analytics']);
  const baseIdx = result.indexOf('base');
  const shopIdx = result.indexOf('shop');
  assert(baseIdx < shopIdx);
});

test('resolveModules handles empty input', () => {
  const result = resolveModules([]);
  assert(result.includes('base'));
});

test('detectConflicts warns about shop without payment', () => {
  const warnings = detectConflicts(['base', 'shop']);
  assert(warnings.length > 0);
});

test('detectConflicts returns empty for valid config', () => {
  const warnings = detectConflicts(['base', 'cms']);
  assert.equal(warnings.length, 0);
});

console.log('All resolver tests passed');
