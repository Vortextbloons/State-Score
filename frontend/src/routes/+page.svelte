<script lang="ts">
	import { onMount } from 'svelte';
	import { getStatus, type AppStatus } from '$lib/api';

	let status = $state<AppStatus | null>(null);
	let error = $state<string | null>(null);
	let loading = $state(true);

	onMount(() => {
		void (async () => {
			try {
				status = await getStatus();
			} catch (err) {
				error = err instanceof Error ? err.message : 'Failed to reach API';
			} finally {
				loading = false;
			}
		})();
	});
</script>

<main>
	<section class="hero">
		<h1>StateScore</h1>
		<p>
			Compare U.S. states across public metrics — rankings, profiles, and transparent scoring, all
			on your machine.
		</p>
	</section>

	<section class="panel" aria-live="polite">
		<h2>Backend status</h2>
		{#if loading}
			<p class="muted">Checking local API…</p>
		{:else if error}
			<p class="bad">{error}</p>
			<p class="muted">Start the Go server and ensure <code>/api</code> is reachable.</p>
		{:else if status}
			<dl>
				<div>
					<dt>Status</dt>
					<dd class:ok={status.status === 'ready'}>{status.status}</dd>
				</div>
				<div>
					<dt>Version</dt>
					<dd>{status.version}</dd>
				</div>
				<div>
					<dt>Database</dt>
					<dd class:ok={status.databaseReady}>
						{status.databaseReady ? 'ready' : 'unavailable'}
					</dd>
				</div>
				<div>
					<dt>Active imports</dt>
					<dd>{status.activeImports}</dd>
				</div>
			</dl>
		{/if}
	</section>
</main>

<style>
	main {
		display: grid;
		gap: 2rem;
		max-width: 40rem;
	}

	.hero h1 {
		margin: 0 0 0.5rem;
		font-size: clamp(2.5rem, 8vw, 3.75rem);
		font-weight: 700;
		letter-spacing: -0.03em;
		color: var(--accent);
	}

	.hero p {
		margin: 0;
		font-size: 1.125rem;
		line-height: 1.55;
		color: var(--muted);
	}

	.panel {
		border-top: 1px solid var(--line);
		padding-top: 1.25rem;
	}

	.panel h2 {
		margin: 0 0 0.75rem;
		font-size: 0.8rem;
		font-weight: 600;
		letter-spacing: 0.06em;
		text-transform: uppercase;
		color: var(--ink);
	}

	.muted {
		margin: 0;
		color: var(--muted);
	}

	.bad {
		margin: 0 0 0.5rem;
		color: var(--danger);
	}

	dl {
		margin: 0;
		display: grid;
		gap: 0.65rem;
	}

	dl div {
		display: grid;
		grid-template-columns: 9rem 1fr;
		gap: 0.75rem;
		align-items: baseline;
	}

	dt {
		margin: 0;
		color: var(--muted);
		font-size: 0.9rem;
	}

	dd {
		margin: 0;
		font-variant-numeric: tabular-nums;
	}

	dd.ok {
		color: var(--ok);
	}

	code {
		font-size: 0.9em;
	}
</style>
