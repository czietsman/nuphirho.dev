<script lang="ts">
	import type { PageData } from './$types';
	import { formatDate } from '$lib/format';

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>nuphirho — Posts</title>
	<meta name="description" content="Writing on software engineering, process, and AI." />
	<meta property="og:title" content="nuphirho — Posts" />
	<meta property="og:description" content="Writing on software engineering, process, and AI." />
	<meta property="og:type" content="website" />
	<meta property="og:url" content="https://blog.nuphirho.dev/" />
</svelte:head>

<div class="container">
	<div class="page-intro">
		<h1>Posts</h1>
		<p>Writing on software engineering, process, and AI.</p>
	</div>
	<ul class="post-list">
		{#each data.posts as post (post.slug)}
			<li class="post-item">
				<h2 class="post-item-title">
					<a href="/{post.slug}">{post.title}</a>
				</h2>
				<div class="post-item-meta">
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
				{#if post.subtitle}<p>{post.subtitle}</p>{/if}
			</li>
		{/each}
	</ul>
</div>
