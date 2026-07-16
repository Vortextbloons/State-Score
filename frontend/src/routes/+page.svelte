<script lang="ts">
	import { onMount } from 'svelte';
	import { loadScores, fmt, type ScoreData } from '$lib/scores';
	import ScoreStripe from '$lib/ScoreStripe.svelte';
	let data = $state<ScoreData | null>(null),
		error = $state('');
	onMount(async () => {
		try {
			data = await loadScores();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Could not load state data';
		}
	});
</script>

<svelte:head><title>Overview · StateScore</title></svelte:head>
<section class="hero"><img src="/favicon.png" alt="StateScore" class="hero-logo" /></section>
{#if error}<div class="card warning">{error}</div>{:else if !data}<div class="card muted">
		Loading the field guide…
	</div>{:else}<section class="overview">
		<div class="card summary">
			<p class="eyebrow">Current perspective</p>
			<h2>Balanced</h2>
			<dl>
				<div>
					<dt>States</dt>
					<dd>50</dd>
				</div>
				<div>
					<dt>Active metrics</dt>
					<dd>{data.metrics.length}</dd>
				</div>
				<div>
					<dt>As-of year</dt>
					<dd>{data.asOfYear ?? 'No data'}</dd>
				</div>
			</dl>
			<p class="muted">
				Scores are relative ranks among states (0–100), not absolute grades. Only categories with
				active data are included; equal category weight within that set.
			</p>
		</div>
		<div class="card leaders">
			<div class="section-title">
				<div>
					<p class="eyebrow">At a glance</p>
					<h2>Leading states</h2>
				</div>
				<a href="/rankings">See all 50 →</a>
			</div>
			{#each data.rows.slice(0, 5) as row, i}<a class="leader" href={`/states/${row.state.code}`}
					><b>{String(i + 1).padStart(2, '0')}</b><span
						><strong>{row.state.name}</strong><ScoreStripe
							scores={data.categories.map((c) => row.categories[c.id])}
						/></span
					><em class="score">{fmt(row.overall)}</em></a
				>{/each}{#if data.rows.every((r) => r.overall == null)}<div class="empty">
					Import metric values to calculate the first ranking. The state catalog is ready.
				</div>{/if}
		</div>
	</section>
	<section class="compare-call">
		<p class="eyebrow">Side by side</p>
		<h2>Different strengths.<br />Same measuring stick.</h2>
		<a class="btn" href="/compare">Compare states</a>
	</section>{/if}

<style>
	.hero {
		padding: 3vh 0 4vh;
		display: flex;
		justify-content: center;
	}
	.hero-logo {
		width: clamp(200px, 40vw, 400px);
		height: auto;
	}
	.overview {
		display: grid;
		grid-template-columns: minmax(250px, 0.7fr) 1.6fr;
		gap: 1rem;
	}
	.summary h2,
	.leaders h2 {
		font: 400 2rem var(--font-display);
		margin: 0;
	}
	.summary dl {
		display: grid;
		grid-template-columns: repeat(3, 1fr);
		margin: 2rem 0;
	}
	.summary dt {
		font-size: 0.7rem;
		color: var(--muted);
		text-transform: uppercase;
	}
	.summary dd {
		font: 700 1.25rem var(--font-data);
		margin: 0;
	}
	.section-title {
		display: flex;
		align-items: end;
		justify-content: space-between;
		margin-bottom: 1rem;
	}
	.section-title a {
		font-size: 0.85rem;
		color: var(--lake);
	}
	.leader {
		display: grid;
		grid-template-columns: 2rem 1fr 4rem;
		gap: 1rem;
		align-items: center;
		padding: 0.75rem 0;
		border-top: 1px solid var(--mist);
		text-decoration: none;
	}
	.leader > b {
		font: 400 0.75rem var(--font-data);
		color: var(--muted);
	}
	.leader span {
		display: grid;
		grid-template-columns: minmax(110px, 1fr) 1fr;
		align-items: center;
		gap: 1rem;
	}
	.leader em {
		text-align: right;
		font-style: normal;
		font-size: 1.15rem;
	}
	.compare-call {
		margin-top: 5rem;
		padding: 3rem;
		border-radius: 22px;
		background: var(--lake);
		color: white;
	}
	.compare-call h2 {
		font: 400 clamp(2rem, 4vw, 4rem)/1 var(--font-display);
		margin: 0.5rem 0 2rem;
	}
	.compare-call .eyebrow {
		color: #c9eef0;
	}
	@media (max-width: 850px) {
		.overview {
			grid-template-columns: 1fr;
		}
		.leader span {
			grid-template-columns: 1fr;
		}
		.leader .stripe {
			display: none;
		}
		.compare-call {
			padding: 2rem;
		}
	}
</style>
