<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();

	let theme = $state('light');
	let menuOpen = $state(false);

	function initTheme() {
		const stored = localStorage.getItem('nuphirho-theme');
		theme = stored ?? (matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
		document.documentElement.setAttribute('data-theme', theme);
		if (!stored) {
			matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
				if (!localStorage.getItem('nuphirho-theme')) {
					theme = e.matches ? 'dark' : 'light';
					document.documentElement.setAttribute('data-theme', theme);
				}
			});
		}
	}

	function toggle() {
		theme = theme === 'dark' ? 'light' : 'dark';
		document.documentElement.setAttribute('data-theme', theme);
		try { localStorage.setItem('nuphirho-theme', theme); } catch {}
	}

	function closeMenu() { menuOpen = false; }

	onMount(() => {
		initTheme();
	});

	$effect(() => {
		const path = $page.url.pathname;
		fetch('/api/stats', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ path })
		}).catch(() => {});
	});
</script>

<a href="#main" class="skip-link">Skip to content</a>

<header class="site-header">
	<div class="container">
		<a href="https://nuphirho.dev" class="site-name" aria-label="nuphirho home">nuphirho</a>

		<button
			class="burger-menu"
			class:active={menuOpen}
			aria-label="Toggle navigation"
			aria-expanded={menuOpen}
			onclick={() => { menuOpen = !menuOpen; }}
		>
			<span></span>
			<span></span>
			<span></span>
		</button>

		<nav class="site-nav" class:active={menuOpen} aria-label="Main navigation">
			<a href="/about" aria-current={$page.url.pathname === '/about' ? 'page' : undefined} onclick={closeMenu}>This Blog</a>
			<a href="https://nuphirho.dev/words-of-meaning" onclick={closeMenu}>Words</a>
			<a href="/" aria-current={$page.url.pathname === '/' ? 'page' : undefined} onclick={closeMenu}>Posts</a>
		</nav>

		<button id="theme-toggle" class="theme-toggle" type="button" onclick={toggle}
			aria-label={theme === 'dark' ? 'Switch to light theme' : 'Switch to dark theme'}>
			{theme === 'dark' ? 'Light' : 'Dark'}
		</button>
	</div>
</header>

<main id="main" class="site-main">
	{@render children()}
</main>

<footer class="site-footer">
	<div class="container">
		<span>&copy; 2026 Christo Zietsman</span>
		<span class="footer-links">
			<a href="https://doi.org/10.48550/arXiv.2603.25773">arXiv Paper 1</a>
			<a href="https://doi.org/10.48550/arXiv.2604.21090">arXiv Paper 2</a>
			<a href="https://github.com/czietsman">GitHub</a>
			<a href="https://www.linkedin.com/in/christo-zietsman/">LinkedIn</a>
			<a href="https://nuphirho.dev/privacy">Privacy</a>
			<a href="https://nuphirho.dev/cookies">Cookies</a>
		</span>
	</div>
</footer>

<style>
	@import '../app.css';
</style>
