#!/usr/bin/env node
// Fetches rankings from the StateScore API and updates the Rankings section
// in README.md. Run while the backend is serving on the configured port.
//
// Usage: node scripts/update-rankings.mjs [base-url]
//   base-url defaults to http://127.0.0.1:8080

import { readFileSync, writeFileSync } from 'node:fs';
import { join, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';

const BASE = process.argv[2] || 'http://127.0.0.1:8080';
const API = `${BASE}/api/v1`;
const README = join(dirname(fileURLToPath(import.meta.url)), '..', 'README.md');

function formatPopulation(value) {
  if (value == null) return '—';
  if (value >= 1_000_000) return `${Number((value / 1_000_000).toFixed(1))}M`;
  if (value >= 1_000) return `${Math.round(value / 1_000)}k`;
  return String(value);
}

async function fetchJSON(path) {
  let res;
  try {
    res = await fetch(`${API}${path}`);
  } catch (err) {
    throw new Error(`GET ${path} failed — is the backend running on ${BASE}?\n  ${err.message}`);
  }
  if (!res.ok) {
    let msg = `${res.status} ${res.statusText}`;
    try { const body = await res.json(); msg = body.error?.message ?? msg; } catch {}
    throw new Error(`GET ${path}: ${msg}`);
  }
  return (await res.json()).data;
}

async function main() {
  const [states, scoreboard] = await Promise.all([
    fetchJSON('/states'),
    fetchJSON('/scores'),
  ]);

  const stateMap = {};
  for (const s of states) stateMap[s.id] = s;

  const ranked = scoreboard.scores
    .map(ss => ({ state: stateMap[ss.stateId], score: ss }))
    .filter(x => x.state)
    .sort((a, b) => b.score.overallScore - a.score.overallScore);

  const lines = [
    `| Rank | State | Code | Population | Score | Completeness |`,
    `|------|-------|------|-----------:|-------|--------------|`,
  ];
  for (let i = 0; i < ranked.length; i++) {
    const r = ranked[i];
    const score = r.score.overallScore.toFixed(1);
    const pct = (r.score.completeness * 100).toFixed(0);
    const population = formatPopulation(r.state.population);
    lines.push(`| ${i + 1} | ${r.state.name} | ${r.state.code} | ${population} | ${score} | ${pct}% |`);
  }

  const table = lines.join('\n');
  let readme = readFileSync(README, 'utf-8');
  const start = '<!-- RANKINGS_START -->';
  const end = '<!-- RANKINGS_END -->';

  const startIdx = readme.indexOf(start);
  const endIdx = readme.indexOf(end);

  if (startIdx === -1 || endIdx === -1) {
    throw new Error('Could not find RANKINGS markers in README.md');
  }

  const before = readme.slice(0, startIdx + start.length);
  const after = readme.slice(endIdx);
  readme = `${before}\n${table}\n${after}`;

  writeFileSync(README, readme, 'utf-8');
  console.log(`Updated README.md with ${ranked.length} rankings.`);
}

main().catch(err => {
  console.error('update-rankings failed:', err.message);
  process.exit(1);
});
