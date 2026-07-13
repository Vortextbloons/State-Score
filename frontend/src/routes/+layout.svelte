<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import '../app.css';

	let { children } = $props();
</script>

<svelte:head>
	<title>StateScore</title>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="shell">
	<header class="top">
		<a class="brand" href="/">StateScore</a>
		<p class="tag">Local U.S. state comparison</p>
	</header>
	{@render children()}
</div>

<style>
	.shell {
		min-height: 100vh;
		padding: 2rem clamp(1.25rem, 4vw, 3rem) 3rem;
		display: flex;
		flex-direction: column;
		gap: 2rem;
	}

	.top {
		display: flex;
		flex-wrap: wrap;
		align-items: baseline;
		gap: 0.75rem 1.25rem;
		border-bottom: 1px solid var(--line);
		padding-bottom: 1rem;
	}

	.brand {
		font-family: var(--font-display);
		font-size: 1.5rem;
		font-weight: 700;
		letter-spacing: -0.03em;
		color: var(--accent);
		text-decoration: none;
	}

	.tag {
		margin: 0;
		color: var(--muted);
		font-size: 0.95rem;
	}
</style>
