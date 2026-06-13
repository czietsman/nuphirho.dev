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
		<header class="post-header">
			<h1 class="post-title">{data.post.title}</h1>
			{#if data.post.subtitle}<p class="post-subtitle">{data.post.subtitle}</p>{/if}
			<div class="post-meta">
				<time datetime={data.post.publishDate}>{formatDate(data.post.publishDate)}</time>
				<span class="reading-time">{data.post.readingTimeMinutes} min read</span>
				{#if data.post.series}<span>{data.post.series}</span>{/if}
				{#if data.post.tags.length}
					<div class="tags">
						{#each data.post.tags as tag}
							<span class="tag">{tag}</span>
						{/each}
					</div>
				{/if}
			</div>
		</header>
		<div class="post-content">
			{@html data.post.html}
		</div>
	</article>
	<a href="/" class="back-link">← All posts</a>
</div>
