<script lang="ts">
	import { papers } from '$lib/papers.js';
</script>

<svelte:head>
	<title>Papers — nuphirho</title>
	<meta name="description" content="Published research papers by Christo Zietsman." />
</svelte:head>

<div class="container">
	<div class="page-header">
		<h1 class="page-title">Papers</h1>
	</div>

	<div class="manifesto-content">
		<ul class="paper-list">
			{#each papers as paper}
				<li class="paper-entry">
					<span class="paper-year">{paper.year}</span>
					<div class="paper-body">
						<a href={paper.url} class="paper-title" rel="noopener noreferrer" target="_blank">{paper.title}</a>
						<p class="paper-summary">{paper.summary}</p>
					</div>
				</li>
			{/each}
		</ul>
	</div>
</div>

<style>
	.paper-list {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 2rem;
	}

	.paper-entry {
		display: flex;
		gap: 1.5rem;
		align-items: flex-start;
	}

	.paper-year {
		flex-shrink: 0;
		font-variant-numeric: tabular-nums;
		color: var(--color-muted, #888);
		font-size: 0.9rem;
		padding-top: 0.15rem;
	}

	.paper-title {
		display: block;
		font-weight: 600;
		margin-bottom: 0.4rem;
		line-height: 1.4;
	}

	.paper-summary {
		margin: 0;
		color: var(--color-muted, #888);
		font-size: 0.95rem;
		line-height: 1.6;
	}
</style>
