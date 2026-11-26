<script lang="ts">
	import { onMount } from 'svelte';
	import { articles, filteredArticles, clearNewCount } from '$lib/stores/articles';
	import ArticleCard from '$lib/components/ArticleCard.svelte';
	import SearchBar from '$lib/components/SearchBar.svelte';
	import type { PageData } from './$types';

	export let data: PageData;

	onMount(() => {
		articles.set(data.articles ?? []);
		clearNewCount();
	});
</script>

<svelte:head>
	<title>Baltimore Sports Hub — Live Orioles &amp; Ravens News</title>
	<meta name="description" content="Real-time Orioles and Ravens news aggregated from GDELT, ESPN, Baltimore Sun, MLB, NFL and more." />
</svelte:head>

<section class="page">
	<header class="page-header">
		<h1>Baltimore Sports Hub</h1>
		<p class="subtitle">Real-time news from GDELT · ESPN · Baltimore Sun · MLB · NFL</p>
	</header>

	<SearchBar />

	{#if data.error}
		<div class="banner error">⚠️ {data.error}</div>
	{/if}

	{#if $filteredArticles.length === 0 && !data.error}
		<div class="empty">
			<p>No articles yet — the scrapers run every 10–15 minutes. Check back soon!</p>
		</div>
	{:else}
		<div class="grid">
			{#each $filteredArticles as article (article.id)}
				<ArticleCard {article} />
			{/each}
		</div>
	{/if}
</section>

<style>
	.page { padding-bottom: 3rem; }
	.page-header { margin-bottom: 1.5rem; }
	h1 {
		font-size: clamp(1.4rem, 4vw, 2rem);
		font-weight: 800;
		letter-spacing: -0.02em;
		margin-bottom: 0.25rem;
	}
	.subtitle { color: #666; font-size: 0.9rem; }
	.banner { border-radius: 8px; padding: 0.75rem 1rem; margin-bottom: 1rem; font-size: 0.9rem; }
	.error { background: #fff3cd; border: 1px solid #ffc107; }
	.empty { text-align: center; padding: 4rem 1rem; color: #888; }
	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 1rem;
	}
</style>
