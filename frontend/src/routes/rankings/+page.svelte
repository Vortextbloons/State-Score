<script lang="ts">
	import { onMount } from 'svelte';
	import { loadScores, fmt, formatPopulation, type ScoreData } from '$lib/scores';
	import ScoreStripe from '$lib/ScoreStripe.svelte';
	let data = $state<ScoreData | null>(null),
		search = $state(''),
		region = $state(''),
		sort = $state('overall'),
		direction = $state(-1),
		selected = $state<string[]>([]);
	onMount(async () => (data = await loadScores()));
	let regions = $derived([...new Set(data?.states.map((s) => s.region).filter(Boolean) ?? [])]);
	let rows = $derived(
		(data?.rows ?? [])
			.filter(
				(r) =>
					(!search ||
						`${r.state.name} ${r.state.code}`.toLowerCase().includes(search.toLowerCase())) &&
					(!region || r.state.region === region)
			)
			.toSorted((a, b) => {
				const av = sort === 'overall' ? a.overall : a.categories[Number(sort)],
					bv = sort === 'overall' ? b.overall : b.categories[Number(sort)];
				return ((av ?? -1) - (bv ?? -1)) * direction;
			})
	);
	function toggle(code: string) {
		selected = selected.includes(code)
			? selected.filter((x) => x !== code)
			: selected.length < 5
				? [...selected, code]
				: selected;
	}
	function exportCSV() {
		if (!data) return;
		const head = [
			'Rank',
			'State',
			'Code',
			'Population',
			'Population year',
			'Overall',
			...data.categories.map((c) => c.name)
		];
		const lines = [
			head,
			...rows.map((r, i) => [
				i + 1,
				r.state.name,
				r.state.code,
				r.state.population ?? '',
				r.state.populationYear ?? '',
				fmt(r.overall),
				...data!.categories.map((c) => fmt(r.categories[c.id]))
			])
		]
			.map((x) => x.join(','))
			.join('\n');
		const a = document.createElement('a');
		a.href = URL.createObjectURL(new Blob([lines], { type: 'text/csv' }));
		a.download = 'statescore-rankings.csv';
		a.click();
		URL.revokeObjectURL(a.href);
	}
</script>

<svelte:head><title>Rankings · StateScore</title></svelte:head>
<div class="page-head">
	<div>
		<p class="eyebrow">All 50 states</p>
		<h1>Rankings</h1>
		<p>
			Numbers are percentile ranks among states for the as-of year{data?.asOfYear
				? ` (${data.asOfYear})`
				: ''} — not absolute quality scores. A dotted completeness marker means some metric weight was
			redistributed.
		</p>
	</div>
	<button class="btn secondary" onclick={exportCSV}>Export CSV</button>
</div>
<section class="card filters">
	<div class="field">
		<label for="search">Find a state</label><input
			id="search"
			bind:value={search}
			placeholder="Name or code"
		/>
	</div>
	<div class="field">
		<label for="region">Region</label><select id="region" bind:value={region}
			><option value="">All regions</option>{#each regions as r (r)}<option>{r}</option
				>{/each}</select
		>
	</div>
	<div class="field">
		<label for="sort">Measure</label><select id="sort" bind:value={sort}
			><option value="overall">Overall</option>{#each data?.categories ?? [] as c (c.id)}<option
					value={c.id}>{c.name}</option
				>{/each}</select
		>
	</div>
	<button class="btn secondary" onclick={() => (direction *= -1)}
		>{direction < 0 ? 'Highest first' : 'Lowest first'}</button
	>
</section>
<section class="card table-wrap">
	{#if !data}<div class="empty">Loading rankings…</div>{:else}<table>
			<thead
				><tr
					><th>Rank</th><th>State</th><th>Population</th><th>Overall</th
					>{#each data.categories as c (c.id)}<th>{c.name}</th>{/each}<th>Compare</th></tr
				></thead
			><tbody
				>{#each rows as row, i (row.state.id)}<tr
						><td class="score">{i + 1}</td><td
							><a class="state" href={`/states/${row.state.code}`}
								><strong>{row.state.name}</strong><small
									>{row.state.code} · {row.state.region}</small
								><ScoreStripe scores={data.categories.map((c) => row.categories[c.id])} /></a
							></td
						><td class="population"
							><strong>{formatPopulation(row.state.population)}</strong><small
								>{row.state.populationYear ? `July ${row.state.populationYear}` : ''}</small
							></td
						><td class="score"
							>{fmt(row.overall)}
							{#if row.completeness < 1}<span
									title={`${Math.round(row.completeness * 100)}% complete`}>◌</span
								>{/if}</td
						>{#each data.categories as c (c.id)}<td class="score muted"
								>{fmt(row.categories[c.id])}</td
							>{/each}<td
							><button
								class="pick"
								class:selected={selected.includes(row.state.code)}
								onclick={() => toggle(row.state.code)}
								aria-label="Add to comparison"
								>{selected.includes(row.state.code) ? '✓' : '+'}</button
							></td
						></tr
					>{/each}</tbody
			>
		</table>
		{#if rows.every((r) => r.overall == null)}<div class="empty">
				No metric values have been imported yet. States will be scored here after an import.
			</div>{/if}{/if}
</section>
{#if selected.length >= 2}<a class="compare-dock" href={`/compare?states=${selected.join(',')}`}
		>Compare {selected.join(' · ')} <b>→</b></a
	>{/if}

<style>
	.filters {
		display: flex;
		align-items: end;
		gap: 0.8rem;
		margin-bottom: 1rem;
	}
	.filters .field:first-child {
		flex: 1;
	}
	.state {
		display: grid;
		min-width: 180px;
		text-decoration: none;
	}
	.population strong,
	.population small {
		display: block;
		white-space: nowrap;
	}
	.population small {
		color: var(--muted);
		font-family: var(--font-data);
	}
	.state small {
		color: var(--muted);
	}
	.pick {
		width: 32px;
		height: 32px;
		border: 1px solid var(--line);
		border-radius: 50%;
		background: white;
		color: var(--blue);
		cursor: pointer;
	}
	.pick.selected {
		background: var(--coral);
		color: white;
		border-color: var(--coral);
	}
	.compare-dock {
		position: fixed;
		right: 2rem;
		bottom: 2rem;
		padding: 0.9rem 1.2rem;
		border-radius: 99px;
		background: var(--coral);
		color: white;
		text-decoration: none;
		box-shadow: 0 8px 25px #17242c44;
	}
	.compare-dock b {
		margin-left: 1rem;
	}
	@media (max-width: 700px) {
		.filters {
			align-items: stretch;
			flex-direction: column;
		}
		.filters .field {
			width: 100%;
		}
	}
</style>
