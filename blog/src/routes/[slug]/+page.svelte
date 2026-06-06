<script lang="ts">
	import type { PageData } from './$types';
	import { formatDate } from '$lib/posts';

	let { data }: { data: PageData } = $props();
	const post = $derived(data.post);
</script>

<svelte:head>
	<title>{post.title} — nuphirho</title>
	{#if post.subtitle}<meta name="description" content={post.subtitle} />{/if}
</svelte:head>

<div class="container">
	<article>
		<header class="post-header">
			<h1 class="post-title">{post.title}</h1>
			{#if post.subtitle}<p class="post-subtitle">{post.subtitle}</p>{/if}
			<div class="post-meta">
				<time datetime={post.publishDate}>{formatDate(post.publishDate)}</time>
				{#if post.series}<span>{post.series}</span>{/if}
				{#if post.tags.length}
					<div class="tags">
						{#each post.tags as tag}
							<span class="tag">{tag}</span>
						{/each}
					</div>
				{/if}
			</div>
		</header>
		<div class="post-content">
			{@html post.html}
		</div>
	</article>
	<a href="/" class="back-link">← All posts</a>
</div>
