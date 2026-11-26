<script lang="ts">
	import LiveFeed from '$lib/components/LiveFeed.svelte';
	import { isLive, newArticleCount, clearNewCount } from '$lib/stores/articles';
</script>

<LiveFeed />

<div class="shell">
	<header class="topbar">
		<div class="topbar-inner">
			<a href="/" class="logo">
				<span>⚾🏈</span>
				<span class="logo-text">Baltimore Sports Hub</span>
			</a>
			<div class="topbar-right">
				{#if $newArticleCount > 0}
					<button class="new-badge" on:click={clearNewCount}>
						+{$newArticleCount} new
					</button>
				{/if}
				<div class="live" class:active={$isLive} aria-label={$isLive ? 'Live' : 'Connecting'}>
					<span class="dot" />
					{$isLive ? 'Live' : 'Connecting…'}
				</div>
			</div>
		</div>
	</header>

	<main class="content">
		<slot />
	</main>

	<footer class="foot">
		Data from GDELT · ESPN · Baltimore Sun · MLB · NFL — Built with Go &amp; SvelteKit
	</footer>
</div>

<style>
	:global(*, *::before, *::after) { box-sizing: border-box; margin: 0; padding: 0; }
	:global(body) {
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
		background: #f4f5f7;
		color: #111;
		line-height: 1.5;
	}
	:global(a) { color: inherit; }
	.shell { display: flex; flex-direction: column; min-height: 100vh; }
	.topbar { background: #111; color: #fff; position: sticky; top: 0; z-index: 100; box-shadow: 0 2px 8px rgba(0,0,0,.25); }
	.topbar-inner { max-width: 1200px; margin: 0 auto; padding: 0.7rem 1.25rem; display: flex; justify-content: space-between; align-items: center; }
	.logo { display: flex; align-items: center; gap: 0.5rem; text-decoration: none; color: #fff; font-size: 1.05rem; font-weight: 700; }
	.logo-text { letter-spacing: -0.01em; }
	.topbar-right { display: flex; align-items: center; gap: 0.75rem; }
	.new-badge { background: #DF4601; color: #fff; border: none; border-radius: 999px; padding: 0.25rem 0.75rem; font-size: 0.78rem; font-weight: 600; cursor: pointer; animation: pulse 1.5s infinite; }
	@keyframes pulse { 0%,100% { opacity:1 } 50% { opacity:.65 } }
	.live { display: flex; align-items: center; gap: 0.35rem; font-size: 0.78rem; color: #888; }
	.live.active { color: #4ade80; }
	.dot { width: 7px; height: 7px; border-radius: 50%; background: currentColor; }
	.live.active .dot { animation: blink 1.5s infinite; }
	@keyframes blink { 0%,100% { opacity:1 } 50% { opacity:.3 } }
	.content { flex: 1; max-width: 1200px; margin: 0 auto; padding: 1.5rem 1.25rem; width: 100%; }
	.foot { text-align: center; padding: 1.5rem; font-size: 0.78rem; color: #aaa; border-top: 1px solid #e0e0e0; }
</style>
