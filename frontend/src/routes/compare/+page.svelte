<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { loadScores, fmt, formatValue, type ScoreData } from '$lib/scores';
	let data = $state<ScoreData | null>(null),
		selected = $state<string[]>([]),
		copyLabel = $state('Copy summary');
	onMount(async () => {
		data = await loadScores();
		const q = (page.url.searchParams.get('states') ?? '').split(',').filter(Boolean);
		selected = q.length ? q.slice(0, 5) : data.states.slice(0, 2).map((s) => s.code);
	});
	let rows = $derived(data?.rows.filter((r) => selected.includes(r.state.code)) ?? []);
	function add(e: Event) {
		const code = (e.currentTarget as HTMLSelectElement).value;
		if (code && !selected.includes(code) && selected.length < 5) selected = [...selected, code];
	}
	async function copy() {
		const text = rows.map((r) => `${r.state.name}: ${fmt(r.overall)}`).join('\n');
		await navigator.clipboard.writeText(text);
		copyLabel = 'Copied';
		setTimeout(() => (copyLabel = 'Copy summary'), 1500);
	}
</script>

<svelte:head><title>Compare states · StateScore</title></svelte:head>
<div class="page-head">
	<div>
		<p class="eyebrow">Two to five states</p>
		<h1>Compare</h1>
		<p>
			Read strengths across the page. Bar length always means a better normalized score, even when a
			raw metric is better when lower.
		</p>
	</div>
	<div class="controls"><button class="btn secondary" onclick={copy}>{copyLabel}</button></div>
</div>
{#if data}<section class="chooser card">
		<div class="chips">
			{#each selected as code}<button
					onclick={() => (selected = selected.filter((x) => x !== code))}
					>{data.states.find((s) => s.code === code)?.name} <span>×</span></button
				>{/each}
		</div>
		{#if selected.length < 5}<select onchange={add} value=""
				><option value="">Add a state…</option
				>{#each data.states.filter((s) => !selected.includes(s.code)) as s}<option value={s.code}
						>{s.name}</option
					>{/each}</select
			>{/if}
	</section>
	<section class="scoreboard card">
		<div class="grid head">
			<span>Overall</span>{#each rows as row}<div>
					<strong>{row.state.code}</strong><b class="score">{fmt(row.overall)}</b><small
						>#{data.rows.indexOf(row) + 1}</small
					>
				</div>{/each}
		</div>
		{#each data.categories as c}<div class="grid category">
				<strong>{c.name}</strong>{#each rows as row}<div class="bar">
						<i style:width={`${row.categories[c.id] ?? 0}%`}></i><span class="score"
							>{fmt(row.categories[c.id])}</span
						>
					</div>{/each}
			</div>{/each}
	</section>
	<section class="card metrics">
		<p class="eyebrow">Raw measures</p>
		<h2>Metric-by-metric</h2>
		{#each data.metrics as m}<div class="grid metric">
				<span
					><strong>{m.name}</strong><small
						>{m.higherIsBetter ? 'Higher is better ↑' : 'Lower is better ↓'}</small
					></span
				>{#each rows as row}{@const v = row.values.find((x) => x.metricId === m.id)}<span
						class="score"
						>{formatValue(v?.value ?? null, m.unit)}<small>{v?.year ?? 'Missing'}</small></span
					>{/each}
			</div>{/each}
	</section>
	{#if selected.length < 2}<div class="empty">
			Add at least two states for a useful comparison.
		</div>{/if}{/if}

<style>
	.chooser {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 1rem;
	}
	.chips {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem;
	}
	.chips button {
		border: 0;
		border-radius: 99px;
		background: var(--mist);
		padding: 0.5rem 0.8rem;
		color: var(--blue);
		font-weight: 700;
		cursor: pointer;
	}
	.chips span {
		margin-left: 0.3rem;
	}
	.grid {
		display: grid;
		grid-template-columns: minmax(150px, 1fr) repeat(var(--count, 5), minmax(110px, 1fr));
		gap: 1rem;
		align-items: center;
	}
	.head {
		--count: 5;
		padding-bottom: 1rem;
		border-bottom: 2px solid var(--blue);
	}
	.head > div {
		display: grid;
	}
	.head b {
		font-size: 1.7rem;
	}
	.head small {
		color: var(--muted);
	}
	.category {
		--count: 5;
		padding: 1rem 0;
		border-bottom: 1px solid var(--mist);
	}
	.bar {
		position: relative;
		height: 36px;
		background: var(--mist);
		border-radius: 6px;
		overflow: hidden;
	}
	.bar i {
		display: block;
		height: 100%;
		background: var(--lake);
		opacity: 0.55;
	}
	.bar span {
		position: absolute;
		inset: 7px 8px;
		text-align: right;
	}
	.metrics {
		margin-top: 1rem;
	}
	.metrics h2 {
		font: 400 2rem var(--font-display);
		margin-top: 0;
	}
	.metric {
		--count: 5;
		padding: 0.8rem 0;
		border-top: 1px solid var(--mist);
	}
	.metric small {
		display: block;
		color: var(--muted);
		font: 400 0.72rem var(--font-body);
	}
	@media (max-width: 900px) {
		.scoreboard,
		.metrics {
			overflow: auto;
		}
		.grid {
			width: 900px;
		}
	}
</style>
