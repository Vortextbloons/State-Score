<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import {
		loadScores,
		fmt,
		formatValue,
		formatPopulation,
		type ScoreData,
		type Row
	} from '$lib/scores';
	let data = $state<ScoreData | null>(null),
		row = $state<Row | null>(null);
	onMount(async () => {
		data = await loadScores();
		row = data.rows.find((r) => r.state.code === (page.params.code ?? '').toUpperCase()) ?? null;
	});
</script>

<svelte:head><title>{row?.state.name ?? 'State profile'} - StateScore</title></svelte:head>
{#if row && data}
	<div class="page-head">
		<div>
			<p class="eyebrow">{row.state.region} / {row.state.code}</p>
			<h1>{row.state.name}</h1>
			<p>
				{row.state.population ? formatPopulation(row.state.population) : 'Population unavailable'} residents
				{row.state.populationYear ? `(July ${row.state.populationYear} estimate)` : ''}. Scores are
				relative ranks among states as of {data.asOfYear ?? 'the latest available year'}. Each
				category below traces to its included metrics and their data years.
			</p>
		</div>
		<div class="overall">
			<small>Overall (relative)</small><strong>{fmt(row.overall)}</strong><span
				>#{data.rows.indexOf(row) + 1} of 50</span
			>
		</div>
	</div>
	{#if row.completeness < 1}<div class="notice">
			Completeness: {Math.round(row.completeness * 100)}%. Missing or quality-excluded metrics are
			omitted and their weight is redistributed.
		</div>{/if}
	<section class="categories">
		{#each data.categories as c, i (c.id)}<article class="card">
				<span class={`dot c${i}`}></span>
				<div>
					<p>{c.name}</p>
					<strong class="score">{fmt(row.categories[c.id])}</strong>
				</div>
				<small>{c.description}</small>
			</article>{/each}
	</section>
	<section class="card metrics">
		<div class="section">
			<div>
				<p class="eyebrow">Underlying measures</p>
				<h2>Metric ledger</h2>
			</div>
			<a class="btn secondary small" href={`/compare?states=${row.state.code}`}>Add to comparison</a
			>
		</div>
		<table>
			<thead
				><tr><th>Metric</th><th>Value</th><th>Year</th><th>Direction</th><th>Method</th></tr></thead
			><tbody
				>{#each data.metrics as m (m.id)}{@const v = row.values.find((x) => x.metricId === m.id)}<tr
						class:excluded={v?.quality && !v.quality.scoringEligible}
						><td
							><strong>{m.name}</strong><small>{m.description}</small
							>{#if v?.quality?.reportingCoverage != null}<small class="coverage"
									>FBI coverage {v.quality.reportingCoverage.toFixed(1)}%{v.quality
										.populationCovered != null
										? ` / ${formatPopulation(v.quality.populationCovered)} residents covered`
										: ''}{v.quality.dataRevision
										? ` / revision ${v.quality.dataRevision}`
										: ''}</small
								>{/if}{#if v?.quality && !v.quality.scoringEligible}<small class="quality-warning"
									>Excluded from score: {v.quality.exclusionReason ??
										'source quality threshold not met'}</small
								>{/if}</td
						><td class="score">{formatValue(v?.value ?? null, m.unit)}</td><td>{v?.year ?? '-'}</td
						><td>{m.higherIsBetter ? 'Higher is better' : 'Lower is better'}</td><td
							>{m.normalizationMethod}</td
						></tr
					>{/each}</tbody
			>
		</table>
	</section>
{:else}<div class="empty">Loading state profile...</div>{/if}

<style>
	.overall {
		min-width: 180px;
		border-left: 5px solid var(--coral);
		padding-left: 1rem;
		display: grid;
	}
	.overall small,
	.overall span {
		color: var(--muted);
	}
	.overall strong {
		font: 700 3.5rem/1 var(--font-data);
	}
	.notice {
		margin-bottom: 1rem;
		padding: 1rem;
		border: 1px solid #e5b48a;
		background: #fff3df;
		border-radius: 10px;
		color: #7b4b22;
	}
	.categories {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
		gap: 0.7rem;
	}
	.categories article {
		display: flex;
		flex-wrap: wrap;
		gap: 0.8rem;
	}
	.categories p {
		margin: 0;
		color: var(--muted);
	}
	.categories strong {
		font-size: 1.6rem;
	}
	.categories small {
		width: 100%;
		color: var(--muted);
	}
	.dot {
		width: 8px;
		border-radius: 5px;
		background: var(--lake);
	}
	.dot.c1 {
		background: #4c75a3;
	}
	.dot.c2 {
		background: #55a67a;
	}
	.dot.c3 {
		background: #d89045;
	}
	.dot.c4 {
		background: var(--coral);
	}
	.metrics {
		margin-top: 1rem;
	}
	.section {
		display: flex;
		justify-content: space-between;
		align-items: end;
	}
	.section h2 {
		font: 400 2rem var(--font-display);
		margin: 0;
	}
	.metrics td small {
		display: block;
		color: var(--muted);
		max-width: 30rem;
	}
	.metrics tr.excluded {
		background: #fff8ed;
	}
	.coverage {
		margin-top: 0.35rem;
		font-family: var(--font-data);
	}
	.quality-warning {
		color: #8a3f22 !important;
		font-weight: 700;
	}
	@media (max-width: 600px) {
		.categories {
			grid-template-columns: 1fr;
		}
	}
</style>
