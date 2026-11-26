<script lang="ts">
	import type { Article } from '$lib/stores/articles';

	export let article: Article;

	const teamColor: Record<string, string> = {
		orioles: '#DF4601',
		ravens:  '#241773',
		both:    '#555',
	};
	const sourceLabel: Record<string, string> = {
		rss:     'RSS',
		gdelt:   'GDELT',
		reddit:  'Reddit',
		bluesky: 'Bluesky',
		youtube: 'YouTube',
	};

	$: accent = teamColor[article.team] ?? '#888';
	$: label  = sourceLabel[article.source] ?? article.source;
	$: timeAgo = formatTimeAgo(article.published_at);
	$: sentiment = sentimentLabel(article.sentiment_score);

	function formatTimeAgo(iso: string): string {
		const diff = Date.now() - new Date(iso).getTime();
		const mins = Math.floor(diff / 60000);
		if (mins < 60)  return `${mins}m ago`;
		const hrs = Math.floor(mins / 60);
		if (hrs < 24)   return `${hrs}h ago`;
		return `${Math.floor(hrs / 24)}d ago`;
	}

	function sentimentLabel(score?: number): string | null {
		if (score == null || score === 0) return null;
		if (score > 0.55) return '😊 Positive';
		if (score < 0.45) return '😟 Negative';
		return '😐 Neutral';
	}
</script>

<article class="card" style="--accent:{accent}">
	{#if article.image_url}
		<a href={article.url} target="_blank" rel="noopener noreferrer" class="img-link">
			<img src={article.image_url} alt={article.title} class="img" loading="lazy" />
		</a>
	{/if}
	<div class="body">
		<div class="meta">
			<span class="chip team">{article.team.toUpperCase()}</span>
			<span class="chip src">{label}</span>
			{#if sentiment}<span class="chip">{sentiment}</span>{/if}
		</div>
		<h2 class="title">
			<a href={article.url} target="_blank" rel="noopener noreferrer">{article.title}</a>
		</h2>
		{#if article.summary}
			<p class="summary">{article.summary}</p>
		{/if}
		<footer class="foot">
			{#if article.author}<span class="author">by {article.author}</span>{/if}
			<time datetime={article.published_at}>{timeAgo}</time>
		</footer>
	</div>
</article>

<style>
	.card { border-left: 4px solid var(--accent); background: #fff; border-radius: 8px; box-shadow: 0 1px 4px rgba(0,0,0,.07); overflow: hidden; transition: box-shadow .2s; }
	.card:hover { box-shadow: 0 4px 16px rgba(0,0,0,.12); }
	.img-link { display: block; }
	.img { width: 100%; height: 170px; object-fit: cover; display: block; }
	.body { padding: 1rem; }
	.meta { display: flex; flex-wrap: wrap; gap: .35rem; margin-bottom: .5rem; }
	.chip { font-size: .68rem; font-weight: 600; padding: 2px 8px; border-radius: 999px; background: #f0f0f0; color: #555; text-transform: uppercase; letter-spacing: .04em; }
	.chip.team { background: var(--accent); color: #fff; }
	.title { font-size: .97rem; font-weight: 600; margin: 0 0 .4rem; line-height: 1.35; }
	.title a { text-decoration: none; color: #111; }
	.title a:hover { text-decoration: underline; }
	.summary { font-size: .84rem; color: #555; margin: 0 0 .75rem; display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden; }
	.foot { display: flex; justify-content: space-between; font-size: .76rem; color: #999; }
	.author { font-style: italic; }
</style>
