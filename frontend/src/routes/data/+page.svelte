<script lang="ts">
	import { api } from '$lib/api/client';
	import type { DataSource, DataImport, ImportIssue } from '$lib/api/types';
	let sources = $state<DataSource[]>([]),
		imports = $state<DataImport[]>([]),
		issues = $state<ImportIssue[]>([]),
		selected = $state<DataImport | null>(null),
		sourceId = $state(''),
		file = $state<File | null>(null),
		busy = $state(false),
		message = $state(''),
		error = $state('');
	let source = $state<Partial<DataSource>>({
		name: '',
		publisher: '',
		sourceUrl: '',
		license: '',
		format: 'csv',
		description: ''
	});
	const format = (value?: string) =>
		value
			? new Intl.DateTimeFormat(undefined, { dateStyle: 'medium', timeStyle: 'short' }).format(
					new Date(value)
				)
			: '—';
	async function load() {
		try {
			[sources, imports] = await Promise.all([api.getSources(), api.getImports()]);
			if (!sourceId && sources[0]) sourceId = String(sources[0].id);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Could not load data tools';
		}
	}
	async function addSource() {
		error = '';
		message = '';
		try {
			const saved = await api.saveSource(source);
			sources = [...sources, saved];
			sourceId = String(saved.id);
			source = {
				name: '',
				publisher: '',
				sourceUrl: '',
				license: '',
				format: 'csv',
				description: ''
			};
			message = 'Source saved.';
		} catch (e) {
			error = e instanceof Error ? e.message : 'Could not save source';
		}
	}
	async function upload() {
		if (!file || !sourceId) return;
		busy = true;
		error = '';
		message = '';
		try {
			await api.uploadCSV(Number(sourceId), file);
			message = 'Import queued. This page will refresh its progress.';
			file = null;
			await poll();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Import failed';
		} finally {
			busy = false;
		}
	}
	async function poll() {
		for (let i = 0; i < 30; i++) {
			await load();
			if (!imports.some((v) => v.status === 'pending' || v.status === 'running')) break;
			await new Promise((r) => setTimeout(r, 500));
		}
	}
	async function inspect(item: DataImport) {
		selected = item;
		const detail = await api.getImport(item.id);
		selected = detail.import;
		issues = detail.errors;
	}
	async function recalc() {
		busy = true;
		error = '';
		try {
			const r = await api.recalculate();
			message = `Recalculated ${r.statesCalculated} states for ${r.year}.`;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Recalculation failed';
		} finally {
			busy = false;
		}
	}
	$effect(() => {
		load();
	});
</script>

<svelte:head><title>Data workshop · StateScore</title></svelte:head>
<section class="page-head">
	<div>
		<p class="eyebrow">Dataset tools</p>
		<h1>Data workshop</h1>
		<p>
			Bring public statistics into StateScore, inspect every rejected row, and keep a traceable
			record of where the numbers came from.
		</p>
	</div>
	<button class="btn secondary" onclick={recalc} disabled={busy}>Recalculate scores</button>
</section>
{#if error}<p class="notice error" role="alert">{error}</p>{/if}{#if message}<p
		class="notice"
		role="status"
	>
		{message}
	</p>{/if}
<div class="workbench">
	<section class="card import-card">
		<p class="eyebrow">01 · Import</p>
		<h2>Load a CSV dataset</h2>
		<p class="muted">
			Required columns: <code>state_code</code>, <code>metric_slug</code>, <code>year</code>,
			<code>value</code>. An optional <code>source_record_id</code> preserves the publisher’s row key.
		</p>
		<div class="fields">
			<label class="field"
				><span>Data source</span><select bind:value={sourceId}
					><option value="">Select a source</option>{#each sources as s}<option value={s.id}
							>{s.name}</option
						>{/each}</select
				></label
			><label class="drop"
				><input
					type="file"
					accept=".csv,text/csv"
					onchange={(e) => (file = e.currentTarget.files?.[0] ?? null)}
				/><strong>{file?.name ?? 'Choose CSV file'}</strong><span>Maximum 10 MB</span></label
			>
		</div>
		<button class="btn" onclick={upload} disabled={!file || !sourceId || busy}
			>{busy ? 'Working…' : 'Start import'}</button
		>
	</section>
	<section class="card source-card">
		<p class="eyebrow">02 · Attribution</p>
		<h2>Register a source</h2>
		<div class="source-grid">
			<label class="field"
				><span>Name</span><input
					bind:value={source.name}
					placeholder="American Community Survey"
				/></label
			><label class="field"
				><span>Publisher</span><input
					bind:value={source.publisher}
					placeholder="U.S. Census Bureau"
				/></label
			><label class="field wide"
				><span>Source address</span><input
					type="url"
					bind:value={source.sourceUrl}
					placeholder="https://…"
				/></label
			><label class="field"
				><span>License</span><input
					bind:value={source.license}
					placeholder="Public domain"
				/></label
			><label class="field"><span>Format</span><input value="CSV" disabled /></label>
		</div>
		<button class="btn secondary" onclick={addSource} disabled={!source.name}>Save source</button>
	</section>
</div>
<section class="card history">
	<div class="section-title">
		<div>
			<p class="eyebrow">Import ledger</p>
			<h2>Validation history</h2>
		</div>
		<span>{imports.length} runs</span>
	</div>
	{#if imports.length === 0}<p class="empty">
			No imports yet. A completed upload will appear here with its validation report.
		</p>{:else}<div class="table-wrap">
			<table>
				<thead
					><tr
						><th>Status</th><th>Started</th><th>Read</th><th>Added</th><th>Rejected</th><th
						></th></tr
					></thead
				><tbody
					>{#each imports as item}<tr
							><td
								><span
									class:bad={item.status === 'failed'}
									class:warn={item.status === 'completed_with_errors'}
									class="status">{item.status.replaceAll('_', ' ')}</span
								></td
							><td>{format(item.startedAt)}</td><td class="score">{item.recordsRead}</td><td
								class="score">{item.recordsInserted}</td
							><td class="score">{item.recordsRejected}</td><td
								><button class="link" onclick={() => inspect(item)}>View report</button></td
							></tr
						>{/each}</tbody
				>
			</table>
		</div>{/if}
</section>
{#if selected}<div
		class="scrim"
		role="presentation"
		onclick={(e) => {
			if (e.target === e.currentTarget) selected = null;
		}}
	>
		<div class="report card" role="dialog" aria-modal="true" aria-labelledby="report-title">
			<button class="close" aria-label="Close report" onclick={() => (selected = null)}>×</button>
			<p class="eyebrow">Import #{selected.id}</p>
			<h2 id="report-title">Validation report</h2>
			<p>
				<span class="status">{selected.status.replaceAll('_', ' ')}</span> · {selected.recordsInserted}
				added · {selected.recordsRejected} rejected
			</p>
			{#if selected.errorSummary}<p class="notice error">
					{selected.errorSummary}
				</p>{/if}{#if issues.length}<div class="table-wrap">
					<table>
						<thead><tr><th>Row</th><th>Field</th><th>Raw value</th><th>Problem</th></tr></thead
						><tbody
							>{#each issues as issue}<tr
									><td>{issue.rowNumber ?? '—'}</td><td><code>{issue.fieldName}</code></td><td
										>{issue.rawValue || '—'}</td
									><td>{issue.errorMessage}</td></tr
								>{/each}</tbody
						>
					</table>
				</div>{:else}<p class="empty">Every row passed validation.</p>{/if}
		</div>
	</div>{/if}

<style>
	h2 {
		font: 400 1.7rem var(--font-display);
		margin: 0.2rem 0 1rem;
	}
	.workbench {
		display: grid;
		grid-template-columns: 1.15fr 0.85fr;
		gap: 1rem;
		margin-bottom: 1rem;
	}
	.workbench .card {
		position: relative;
		overflow: hidden;
	}
	.import-card:after {
		content: 'CSV';
		position: absolute;
		right: -0.08em;
		bottom: -0.25em;
		font: 700 8rem/1 var(--font-data);
		color: #287d8e0b;
		pointer-events: none;
	}
	.fields {
		display: grid;
		gap: 1rem;
		margin: 1.5rem 0;
	}
	.drop {
		position: relative;
		display: grid;
		place-items: center;
		min-height: 120px;
		border: 1px dashed var(--lake);
		border-radius: 12px;
		background: #287d8e08;
		cursor: pointer;
	}
	.drop input {
		position: absolute;
		inset: 0;
		opacity: 0;
		cursor: pointer;
	}
	.drop span {
		font-size: 0.75rem;
		color: var(--muted);
	}
	.source-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 0.8rem;
		margin-bottom: 1rem;
	}
	.wide {
		grid-column: 1/-1;
	}
	.history {
		margin-top: 1rem;
	}
	.section-title {
		display: flex;
		align-items: start;
		justify-content: space-between;
	}
	.section-title h2 {
		margin: 0;
	}
	.section-title > span {
		font: 700 0.75rem var(--font-data);
		color: var(--muted);
	}
	.status {
		display: inline-block;
		padding: 0.22rem 0.5rem;
		border-radius: 99px;
		background: #26715f18;
		color: var(--good);
		font-size: 0.72rem;
		font-weight: 800;
		text-transform: uppercase;
	}
	.status.bad {
		background: #f064491a;
		color: #a33d2a;
	}
	.status.warn {
		background: #d890451c;
		color: #85501b;
	}
	.link {
		border: 0;
		background: none;
		color: var(--lake);
		font-weight: 700;
		text-decoration: underline;
		cursor: pointer;
	}
	.notice {
		padding: 0.75rem 1rem;
		border-left: 4px solid var(--good);
		background: #26715f10;
	}
	.notice.error {
		border-color: var(--coral);
		background: #f0644910;
	}
	.scrim {
		position: fixed;
		z-index: 20;
		inset: 0;
		display: grid;
		place-items: center;
		padding: 1rem;
		background: #17242caa;
	}
	.report {
		position: relative;
		width: min(850px, 100%);
		max-height: 85vh;
		overflow: auto;
	}
	.close {
		position: absolute;
		right: 1rem;
		top: 0.7rem;
		border: 0;
		background: none;
		font-size: 2rem;
		cursor: pointer;
	}
	code {
		font-family: var(--font-data);
		font-size: 0.85em;
	}
	@media (max-width: 950px) {
		.workbench {
			grid-template-columns: 1fr;
		}
	}
	@media (max-width: 560px) {
		.source-grid {
			grid-template-columns: 1fr;
		}
		.wide {
			grid-column: auto;
		}
	}
</style>
