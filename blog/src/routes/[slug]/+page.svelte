<script lang="ts">
	import type { PageData } from './$types';
	import { formatDate } from '$lib/posts';

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>{data.post.title} — nuphirho</title>
	{#if data.post.subtitle}<meta name="description" content={data.post.subtitle} />{/if}
</svelte:head>

<div class="container">
	<article>
		<header class="post-header">
			<h1 class="post-title">{data.post.title}</h1>
			{#if data.post.subtitle}<p class="post-subtitle">{data.post.subtitle}</p>{/if}
			<div class="post-meta">
				<time datetime={data.post.publishDate}>{formatDate(data.post.publishDate)}</time>
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
