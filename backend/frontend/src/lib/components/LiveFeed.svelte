<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { isLive, prependArticle, teamFilter } from '$lib/stores/articles';
	import type { Article } from '$lib/stores/articles';

	const API = import.meta.env.VITE_API_URL ?? '';
	let es: EventSource | null = null;
	let retryTimer: ReturnType<typeof setTimeout>;

	function connect() {
		const team = $teamFilter;
		const url = `${API}/api/stream${team ? `?team=${team}` : ''}`;
		es = new EventSource(url);
		es.addEventListener('connected', () => isLive.set(true));
		es.addEventListener('article', (e: MessageEvent) => {
			try { prependArticle(JSON.parse(e.data) as Article); } catch {}
		});
		es.onerror = () => {
			isLive.set(false);
			es?.close();
			retryTimer = setTimeout(connect, 5000);
		};
	}

	onMount(connect);
	onDestroy(() => { es?.close(); clearTimeout(retryTimer); isLive.set(false); });
</script>
