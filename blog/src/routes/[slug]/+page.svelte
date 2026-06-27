<script lang="ts">
	import type { PageData } from './$types';
	import { formatDate } from '$lib/format';

	let { data }: { data: PageData } = $props();

	const canonicalUrl = $derived(`https://blog.nuphirho.dev/${data.post.slug}`);
	const description = $derived(data.post.subtitle ?? data.post.title);

	const jsonLd = $derived({
		'@context': 'https://schema.org',
		'@type': 'BlogPosting',
		headline: data.post.title,
		description,
		datePublished: data.post.publishDate,
		url: canonicalUrl,
		author: { '@type': 'Person', name: 'Christo Zietsman', url: 'https://nuphirho.dev/about' },
		publisher: { '@type': 'Organization', name: 'nuphirho', url: 'https://blog.nuphirho.dev' }
	});
</script>

<svelte:head>
	<title>{data.post.title} — nuphirho</title>
	<meta name="description" content={description} />
	<meta property="og:title" content={data.post.title} />
	<meta property="og:description" content={description} />
	<meta property="og:type" content="article" />
	<meta property="og:url" content={canonicalUrl} />
	<meta property="article:published_time" content={data.post.publishDate} />
	<link rel="canonical" href={canonicalUrl} />
	{@html `<script type="application/ld+json">${JSON.stringify(jsonLd)}</script>`}
</svelte:head>

<div class="container">
	<article>
		{#if data.post.coverImage}
			<img class="cover-image" src={data.post.coverImage} alt={data.post.title} />
		{/if}
		<header class="post-header">
			<h1 class="post-title">{data.post.title}</h1>
			{#if data.post.subtitle}<p class="post-subtitle">{data.post.subtitle}</p>{/if}
			<div class="post-meta">
				<time datetime={data.post.publishDate}>{formatDate(data.post.publishDate)} · {data.post.readingTimeMinutes} min read</time>
				{#if data.post.series}<span>{data.post.series}</span>{/if}
				{#if data.post.tags.length}
					<div class="tags">
						{#each data.post.tags as tag}
							<a href="/tags/{tag}" class="tag">{tag}</a>
						{/each}
					</div>
				{/if}
			</div>
		</header>

		{#if data.post.toc.length >= 3}
			<nav class="toc" aria-label="Table of contents">
				<h2>Contents</h2>
				<ul>
					{#each data.post.toc as entry}
						<li><a href="#{entry.id}">{entry.text}</a></li>
					{/each}
				</ul>
			</nav>
		{/if}

		<div class="post-content">
			{@html data.post.html}
		</div>

		{#if data.post.prevInSeries || data.post.nextInSeries}
			<nav class="series-nav" aria-label="Series navigation">
				{#if data.post.prevInSeries}
					<a href="/{data.post.prevInSeries.slug}" class="series-prev">← {data.post.prevInSeries.title}</a>
				{/if}
				{#if data.post.nextInSeries}
					<a href="/{data.post.nextInSeries.slug}" class="series-next">{data.post.nextInSeries.title} →</a>
				{/if}
			</nav>
		{/if}
	</article>
	<a href="/" class="back-link">← All posts</a>
</div>
