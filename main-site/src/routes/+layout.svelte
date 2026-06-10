<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

	const STORAGE_KEY = 'nuphirho-theme';
	let theme = $state<'light' | 'dark'>('light');
	let menuOpen = $state(false);

	let { children } = $props();

	function initTheme() {
		const applied = document.documentElement.getAttribute('data-theme') as 'light' | 'dark' | null;
		theme = applied ?? 'light';
	}

	function toggleTheme() {
		theme = theme === 'dark' ? 'light' : 'dark';
		document.documentElement.setAttribute('data-theme', theme);
		try { localStorage.setItem(STORAGE_KEY, theme); } catch {}
	}

	onMount(initTheme);

	function closeMenu() { menuOpen = false; }
</script>

<a href="#main" class="skip-link">Skip to content</a>

<header class="site-header">
	<div class="container">
		<a href="/" class="site-name" aria-label="nuphirho home">nuphirho</a>

		<button
			class="burger-menu"
			class:active={menuOpen}
			id="burger-menu"
			aria-label="Toggle navigation"
			aria-expanded={menuOpen}
			onclick={() => { menuOpen = !menuOpen; }}
		>
			<span></span>
			<span></span>
			<span></span>
		</button>

		<nav class="site-nav" class:active={menuOpen} id="site-nav" aria-label="Main navigation">
			<a href="/about" aria-current={$page.url.pathname === '/about' ? 'page' : undefined} onclick={closeMenu}>Who am I</a>
			<a href="/words-of-meaning" aria-current={$page.url.pathname === '/words-of-meaning' ? 'page' : undefined} onclick={closeMenu}>Words</a>
			<a href="https://promptq.ai" target="_blank" rel="noopener">PromptQ</a>
			<a href="https://blog.nuphirho.dev">Blog</a>
			<a href="/novel-findings" aria-current={$page.url.pathname === '/novel-findings' ? 'page' : undefined} onclick={closeMenu}>Novel Findings</a>
		</nav>

		<button
			id="theme-toggle"
			class="theme-toggle"
			type="button"
			aria-label={theme === 'dark' ? 'Switch to light theme' : 'Switch to dark theme'}
			onclick={toggleTheme}
		>{theme === 'dark' ? 'Light' : 'Dark'}</button>
	</div>
</header>

<main id="main" class="site-main" class:landing={$page.url.pathname === '/'}>
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
			<a href="/privacy">Privacy</a>
			<a href="/cookies">Cookies</a>
		</span>
	</div>
</footer>
