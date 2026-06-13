<script lang="ts">
	import type { PageData } from './$types';
	import { formatDate } from '$lib/format';

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>#{data.tag} — nuphirho</title>
	<meta name="description" content="Posts tagged {data.tag} on nuphirho." />
	<link rel="canonical" href="https://blog.nuphirho.dev/tags/{data.tag}" />
</svelte:head>

<div class="container">
	<div class="page-intro">
		<h1>#{data.tag}</h1>
		<p>{data.posts.length} post{data.posts.length === 1 ? '' : 's'}</p>
	</div>
	<ul class="post-list">
		{#each data.posts as post (post.slug)}
			<li class="post-item">
				<h2 class="post-item-title">
					<a href="/{post.slug}">{post.title}</a>
				</h2>
				<div class="post-item-meta">
					<time datetime={post.publishDate}>{formatDate(post.publishDate)} · {post.readingTimeMinutes} min read</time>
					{#if post.series}<span>{post.series}</span>{/if}
				</div>
				{#if post.subtitle}<p>{post.subtitle}</p>{/if}
			</li>
		{/each}
	</ul>
	<a href="/" class="back-link">← All posts</a>
</div>
