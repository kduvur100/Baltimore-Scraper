<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { isLive, prependArticle, teamFilter } from '$lib/stores/articles';
	import type { Article } from '$lib/stores/articles';

	const API = import.meta.env.VITE_API_URL ?? '';
	let es: EventSource | null = null;
	let retryTimer: ReturnType<typeof setTimeout>;

	function disconnect() {
		es?.close();
		es = null;
		clearTimeout(retryTimer);
	}

	function connect() {
		disconnect();
		const team = $teamFilter;
		const url = `${API}/api/stream${team ? `?team=${team}` : ''}`;
		es = new EventSource(url);
		es.addEventListener('connected', () => isLive.set(true));
		es.addEventListener('article', (e: MessageEvent) => {
			try { prependArticle(JSON.parse(e.data) as Article); } catch {}
		});
		es.onerror = () => {
			isLive.set(false);
			disconnect();
			retryTimer = setTimeout(connect, 5000);
		};
	}

	$: $teamFilter, connect();

	onMount(connect);
	onDestroy(() => { disconnect(); isLive.set(false); });
</script>
