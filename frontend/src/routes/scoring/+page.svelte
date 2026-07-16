<script lang="ts">
	import { onMount } from 'svelte';
	import { loadScores, fmt, type ScoreData, type Row } from '$lib/scores';

	type SavedPerspective = { name?: string; weights?: Record<string, number> };

	let data = $state<ScoreData | null>(null);
	let weights = $state<Record<number, number>>({});
	let normalize = $state(true);
	let name = $state('My perspective');
	let error = $state('');
	let saved = $state(false);

	const total = $derived(Object.values(weights).reduce((sum, weight) => sum + weight, 0));
	const isValid = $derived(normalize ? total > 0 : total === 100);
	const effectiveWeights = $derived.by(() => {
		if (!isValid) return {} as Record<number, number>;
		const divisor = normalize ? total : 100;
		return Object.fromEntries(
			Object.entries(weights).map(([id, weight]) => [Number(id), weight / divisor])
		);
	});
	const rows = $derived.by(() => {
		if (!data || !isValid) return [] as Row[];
		return data.rows
			.map((row) => {
				const available = data!.categories.filter(
					(category) => row.categories[category.id] != null
				);
				const availableWeight = available.reduce(
					(sum, category) => sum + (effectiveWeights[category.id] ?? 0),
					0
				);
				const overall = availableWeight
					? available.reduce(
							(sum, category) =>
								sum + (row.categories[category.id] ?? 0) * effectiveWeights[category.id],
							0
						) / availableWeight
					: null;
				return { ...row, overall };
			})
			.toSorted((a, b) => (b.overall ?? -1) - (a.overall ?? -1));
	});

	onMount(async () => {
		try {
			data = await loadScores();
			const share = data.categories.length ? 100 / data.categories.length : 0;
			weights = Object.fromEntries(data.categories.map((category) => [category.id, share]));

			const raw = localStorage.getItem('statescore-profile');
			if (raw) {
				const profile = JSON.parse(raw) as SavedPerspective;
				const restored = Object.fromEntries(
					data.categories.map((category) => [
						category.id,
						clamp(Number(profile.weights?.[String(category.id)] ?? share))
					])
				);
				if (Object.values(restored).some((value) => Number.isFinite(value))) weights = restored;
				if (profile.name?.trim()) name = profile.name.trim();
			}
		} catch (cause) {
			error = cause instanceof Error ? cause.message : 'Could not load your priorities.';
		}
	});

	function clamp(value: number) {
		return Number.isFinite(value) ? Math.max(0, Math.min(100, value)) : 0;
	}

	function change(id: number, value: number) {
		weights = { ...weights, [id]: clamp(value) };
		saved = false;
	}

	function reset() {
		if (!data) return;
		const share = data.categories.length ? 100 / data.categories.length : 0;
		weights = Object.fromEntries(data.categories.map((category) => [category.id, share]));
		name = 'My perspective';
		saved = false;
	}

	function save() {
		if (!isValid) return;
		localStorage.setItem(
			'statescore-profile',
			JSON.stringify({ name: name.trim() || 'My perspective', weights })
		);
		saved = true;
	}
</script>

<svelte:head><title>Your priorities · StateScore</title></svelte:head>

<div class="page-head">
	<div>
		<p class="eyebrow">Build your perspective</p>
		<h1>Your priorities</h1>
		<p>
			Turn up what matters to you and watch the ranking respond. Your choices stay on this device.
		</p>
	</div>
</div>

{#if error}
	<div class="card warning" role="alert"><strong>Priorities are unavailable.</strong> {error}</div>
{:else if !data}
	<div class="card loading" aria-live="polite">Loading your priority controls…</div>
{:else}
	<section class="workspace">
		<div class="card editor">
			<div class="profile">
				<div class="field">
					<label for="name">Perspective name</label>
					<input id="name" bind:value={name} oninput={() => (saved = false)} />
				</div>
				<label class="check">
					<input type="checkbox" bind:checked={normalize} onchange={() => (saved = false)} />
					<span
						><strong>Keep the mix at 100%</strong><small
							>We’ll convert your numbers into a balanced share.</small
						></span
					>
				</label>
			</div>

			<div class="mix-heading">
				<div>
					<p class="eyebrow">Priority mix</p>
					<h2>What should count most?</h2>
				</div>
				<span class:bad={!isValid}
					>{normalize ? `${Math.round(total)} points entered` : `${Math.round(total)}% total`}</span
				>
			</div>

			{#each data.categories as category, i}
				<div class="weight">
					<span class={`swatch c${i}`}></span>
					<label for={`w${category.id}`}>
						<span
							><strong>{category.name}</strong><em
								>{isValid
									? `${Math.round((effectiveWeights[category.id] ?? 0) * 100)}% of score`
									: '—'}</em
							></span
						>
						<small>{category.description}</small>
					</label>
					<input
						id={`w${category.id}`}
						type="range"
						min="0"
						max="100"
						step="1"
						value={weights[category.id]}
						oninput={(event) =>
							change(category.id, +(event.currentTarget as HTMLInputElement).value)}
					/>
					<div class="number">
						<input
							aria-label={`${category.name} priority`}
							type="number"
							min="0"
							max="100"
							step="1"
							value={Math.round(weights[category.id])}
							oninput={(event) =>
								change(category.id, +(event.currentTarget as HTMLInputElement).value)}
						/><span>{normalize ? 'pts' : '%'}</span>
					</div>
				</div>
			{/each}

			{#if !isValid}
				<p class="validation" role="status">
					{total === 0
						? 'Add at least one priority to calculate a ranking.'
						: `Manual weights must total 100%. Adjust by ${Math.abs(100 - total)}%.`}
				</p>
			{/if}

			<footer>
				<button class="btn secondary" onclick={reset}>Reset to equal</button>
				<div>
					<span class="save-status" aria-live="polite">{saved ? 'Saved on this device' : ''}</span
					><button class="btn" onclick={save} disabled={!isValid}
						>{saved ? 'Saved' : 'Save perspective'}</button
					>
				</div>
			</footer>
		</div>

		<aside class="card preview">
			<div class="preview-head">
				<div>
					<p class="eyebrow">Live ranking</p>
					<h2>{name.trim() || 'My perspective'}</h2>
				</div>
				<span>{rows.length ? 'Updated' : 'Waiting'}</span>
			</div>
			{#if rows.length}
				<ol>
					{#each rows.slice(0, 8) as row, i}
						<li>
							<a href={`/states/${row.state.code}`}
								><b>{String(i + 1).padStart(2, '0')}</b><span>{row.state.name}</span><em
									class="score">{fmt(row.overall)}</em
								></a
							>
						</li>
					{/each}
				</ol>
				<a class="all-rankings" href="/rankings">View the standard ranking →</a>
			{:else}
				<div class="empty">Fix the priority total to preview your ranking.</div>
			{/if}
			<p class="muted">
				Only the mix changes here. The underlying public data and category scores stay the same.
			</p>
		</aside>
	</section>
{/if}

<style>
	.workspace {
		display: grid;
		grid-template-columns: minmax(0, 1.55fr) minmax(280px, 0.7fr);
		gap: 1rem;
		align-items: start;
	}
	.profile {
		display: grid;
		grid-template-columns: minmax(220px, 0.8fr) 1.2fr;
		gap: 2rem;
		align-items: end;
		padding-bottom: 1.5rem;
	}
	.check {
		display: flex;
		gap: 0.75rem;
		align-items: flex-start;
		padding: 0.75rem;
		border: 1px solid var(--line);
		border-radius: 12px;
		background: color-mix(in srgb, var(--mist) 55%, transparent);
		cursor: pointer;
	}
	.check input {
		margin-top: 0.2rem;
	}
	.check span {
		display: grid;
	}
	.check small {
		color: var(--muted);
	}
	.mix-heading {
		display: flex;
		justify-content: space-between;
		gap: 1rem;
		align-items: end;
		padding: 0.25rem 0 0.75rem;
		border-bottom: 2px solid var(--blue);
	}
	.mix-heading h2 {
		margin: 0;
		font: 400 1.65rem var(--font-display);
	}
	.mix-heading > span {
		font: 700 0.78rem var(--font-data);
		color: var(--lake);
	}
	.weight {
		display: grid;
		grid-template-columns: 8px minmax(170px, 1fr) 1.2fr 92px;
		gap: 1rem;
		align-items: center;
		padding: 1.05rem 0;
		border-bottom: 1px solid var(--mist);
	}
	.weight label > span {
		display: flex;
		justify-content: space-between;
		gap: 0.5rem;
	}
	.weight label em {
		font: 700 0.7rem var(--font-data);
		font-style: normal;
		color: var(--lake);
		white-space: nowrap;
	}
	.weight label small {
		display: block;
		color: var(--muted);
		margin-top: 0.15rem;
	}
	.swatch {
		height: 48px;
		border-radius: 5px;
		background: var(--lake);
	}
	.swatch.c1 {
		background: #4c75a3;
	}
	.swatch.c2 {
		background: #55a67a;
	}
	.swatch.c3 {
		background: #d89045;
	}
	.swatch.c4 {
		background: var(--coral);
	}
	input[type='range'] {
		width: 100%;
		accent-color: var(--lake);
		cursor: pointer;
	}
	.number {
		display: flex;
		align-items: center;
		gap: 0.35rem;
	}
	.number input {
		width: 62px;
		font-family: var(--font-data);
		text-align: right;
	}
	.number span {
		font-size: 0.72rem;
		color: var(--muted);
	}
	.validation {
		margin: 1rem 0 0;
		padding: 0.75rem 1rem;
		border-left: 3px solid var(--coral);
		background: color-mix(in srgb, var(--coral) 10%, transparent);
		color: var(--ink);
	}
	footer {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 1rem;
		padding-top: 1.5rem;
	}
	footer > div {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}
	.save-status {
		font-size: 0.78rem;
		color: var(--good);
	}
	button:disabled {
		opacity: 0.45;
		cursor: not-allowed;
	}
	.bad {
		color: var(--coral) !important;
	}
	.preview {
		position: sticky;
		top: 2rem;
		overflow: hidden;
	}
	.preview-head {
		display: flex;
		justify-content: space-between;
		gap: 1rem;
		align-items: start;
	}
	.preview-head h2 {
		font: 400 2rem var(--font-display);
		margin: 0;
		overflow-wrap: anywhere;
	}
	.preview-head > span {
		font: 700 0.65rem var(--font-data);
		text-transform: uppercase;
		letter-spacing: 0.08em;
		color: var(--good);
		padding: 0.3rem 0.5rem;
		border: 1px solid color-mix(in srgb, var(--good) 45%, transparent);
		border-radius: 99px;
	}
	.preview ol {
		list-style: none;
		padding: 0;
		margin: 1rem 0;
	}
	.preview li {
		border-top: 1px solid var(--mist);
	}
	.preview li:first-child {
		border-top: 2px solid var(--blue);
	}
	.preview a {
		display: grid;
		grid-template-columns: 2rem 1fr auto;
		gap: 0.6rem;
		padding: 0.78rem 0;
		text-decoration: none;
	}
	.preview a:hover span {
		color: var(--lake);
	}
	.preview a b {
		font-family: var(--font-data);
		color: var(--muted);
	}
	.preview a em {
		font-style: normal;
	}
	.preview .all-rankings {
		display: block;
		padding: 0.2rem 0 1rem;
		color: var(--lake);
		font-weight: 700;
	}
	.preview > p:last-child {
		font-size: 0.8rem;
		margin: 1.5rem 0 0;
	}
	.loading {
		color: var(--muted);
	}
	@media (max-width: 900px) {
		.workspace {
			grid-template-columns: 1fr;
		}
		.preview {
			position: static;
		}
	}
	@media (max-width: 650px) {
		.profile {
			grid-template-columns: 1fr;
			gap: 1rem;
		}
		.weight {
			grid-template-columns: 8px 1fr 88px;
		}
		.weight > input[type='range'] {
			grid-column: 2/4;
		}
		.mix-heading {
			align-items: start;
			flex-direction: column;
		}
		footer {
			align-items: stretch;
			flex-direction: column;
		}
		footer > div {
			justify-content: space-between;
		}
		.save-status {
			min-height: 1.2em;
		}
	}
</style>
