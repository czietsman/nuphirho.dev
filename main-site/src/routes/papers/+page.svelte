<script lang="ts">
	import { papers, type Paper } from '$lib/papers.js';

	interface Erratum {
		date: string;
		location: string;
		correction: string;
		reason: string;
		type?: string;
	}

	interface PaperErrata {
		slug: string;
		linkedinUrl: string | null;
		vnextMarkdownUrl: string;
		vnextPdfUrl: string;
		errata: Erratum[];
	}

	const generatedErrata = import.meta.glob<{ errataBySlug: Record<string, PaperErrata> }>(
		'../../lib/errata.ts',
		{ eager: true },
	);
	const defaultErrata = generatedErrata['../../lib/errata.ts']?.errataBySlug ?? {};

	let {
		paperEntries = papers,
		paperErrata = defaultErrata,
	}: {
		paperEntries?: Paper[];
		paperErrata?: Record<string, PaperErrata>;
	} = $props();

	const programmePapers = $derived(paperEntries.filter((paper) => paper.kind === 'programme'));
	const externalPapers = $derived(paperEntries.filter((paper) => paper.kind === 'external'));
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
		<h2>Programme papers</h2>
		<ul class="paper-list">
			{#each programmePapers as paper}
				<li class="paper-entry">
					<span class="paper-year">{paper.date}</span>
					<div class="paper-body">
						<a href={paper.arxivUrl} class="paper-title" rel="noopener noreferrer" target="_blank">{paper.title}</a>
						<p class="paper-summary">{paper.summary}</p>
						{#if paperErrata[paper.slug]}
							<div class="paper-controls">
								<a href={paperErrata[paper.slug].vnextMarkdownUrl} rel="noopener noreferrer" target="_blank">Maintained next version (Markdown)</a>
								<a href={paperErrata[paper.slug].vnextPdfUrl} rel="noopener noreferrer" target="_blank">Maintained next version (PDF)</a>
								{#if paperErrata[paper.slug].linkedinUrl}
									<a href={paperErrata[paper.slug].linkedinUrl} rel="noopener noreferrer" target="_blank">Discuss on LinkedIn</a>
								{/if}
							</div>

							{#if paperErrata[paper.slug].errata.length > 0}
								<details class="errata">
									<summary>Errata ({paperErrata[paper.slug].errata.length})</summary>
									<ol>
										{#each paperErrata[paper.slug].errata as entry}
											<li>
												<p class="erratum-heading"><time>{entry.date}</time> · {entry.location}{entry.type ? ` · ${entry.type}` : ''}</p>
												<p><strong>Correction:</strong> {entry.correction}</p>
												<p><strong>Reason:</strong> {entry.reason}</p>
											</li>
										{/each}
									</ol>
									<p class="wayback-note">The corrected version is maintained here in the open so it remains discoverable through web archives if the arXiv version is never updated.</p>
								</details>
							{/if}
						{/if}
					</div>
				</li>
			{/each}
		</ul>

		{#if externalPapers.length > 0}
			<section class="prior-work">
				<h2>Prior work</h2>
				<ul class="paper-list">
					{#each externalPapers as paper}
						<li class="external-entry">
							<p class="paper-meta">{paper.author} · {paper.institution} · {paper.degree} · {paper.year}</p>
							<h3>{paper.title}</h3>
							<p class="paper-summary">{paper.abstract}</p>
							<a href={paper.externalUrl} rel="noopener noreferrer" target="_blank">View on SUNScholar</a>
						</li>
					{/each}
				</ul>
			</section>
		{/if}
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

	.paper-controls {
		display: flex;
		flex-wrap: wrap;
		gap: 0.5rem 1rem;
		margin-top: 0.75rem;
		font-size: 0.9rem;
	}

	.errata {
		margin-top: 1rem;
		border-left: 2px solid var(--color-border, #ddd);
		padding-left: 1rem;
	}

	.errata summary {
		cursor: pointer;
		font-weight: 600;
	}

	.errata ol {
		padding-left: 1.25rem;
	}

	.errata li + li {
		margin-top: 1rem;
	}

	.errata p {
		margin: 0.35rem 0;
	}

	.erratum-heading,
	.paper-meta,
	.wayback-note {
		color: var(--color-muted, #888);
		font-size: 0.85rem;
	}

	.prior-work {
		margin-top: 3rem;
	}

	.external-entry h3 {
		margin: 0.25rem 0 0.4rem;
	}

	.paper-meta {
		margin: 0;
	}

	@media (max-width: 600px) {
		.paper-entry {
			flex-direction: column;
			gap: 0.5rem;
		}
	}
</style>
