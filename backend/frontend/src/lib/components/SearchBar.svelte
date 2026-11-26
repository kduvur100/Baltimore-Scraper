<script lang="ts">
	import { searchQuery, teamFilter } from '$lib/stores/articles';

	let inputValue = '';
	let debounceTimer: ReturnType<typeof setTimeout>;

	function onInput() {
		clearTimeout(debounceTimer);
		debounceTimer = setTimeout(() => searchQuery.set(inputValue.trim()), 300);
	}

	function clearSearch() {
		inputValue = '';
		searchQuery.set('');
	}
</script>

<div class="bar">
	<div class="search-wrap">
		<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
			<circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
		</svg>
		<input type="search" placeholder="Search articles…" bind:value={inputValue} on:input={onInput} class="input" />
		{#if inputValue}
			<button class="clear" on:click={clearSearch} aria-label="Clear search">✕</button>
		{/if}
	</div>
	<div class="filters">
		<button class="pill" class:active={$teamFilter === ''} on:click={() => teamFilter.set('')}>All</button>
		<button class="pill orioles" class:active={$teamFilter === 'orioles'} on:click={() => teamFilter.set('orioles')}>⚾ Orioles</button>
		<button class="pill ravens" class:active={$teamFilter === 'ravens'} on:click={() => teamFilter.set('ravens')}>🏈 Ravens</button>
	</div>
</div>

<style>
	.bar { display: flex; flex-wrap: wrap; gap: .75rem; align-items: center; margin-bottom: 1.5rem; }
	.search-wrap { position: relative; flex: 1; min-width: 200px; }
	.icon { position: absolute; left: .75rem; top: 50%; transform: translateY(-50%); width: 18px; height: 18px; color: #aaa; pointer-events: none; }
	.input { width: 100%; padding: .55rem .55rem .55rem 2.4rem; border: 1px solid #ddd; border-radius: 999px; font-size: .9rem; outline: none; transition: border-color .15s; box-sizing: border-box; }
	.input:focus { border-color: #888; }
	.clear { position: absolute; right: .75rem; top: 50%; transform: translateY(-50%); background: none; border: none; cursor: pointer; color: #aaa; font-size: .8rem; padding: 0; }
	.clear:hover { color: #555; }
	.filters { display: flex; gap: .4rem; }
	.pill { padding: .4rem .9rem; border: 1px solid #ddd; border-radius: 999px; background: #fff; font-size: .82rem; cursor: pointer; transition: background .15s, border-color .15s; }
	.pill:hover { background: #f5f5f5; }
	.pill.active { background: #222; color: #fff; border-color: #222; }
	.pill.orioles.active { background: #DF4601; border-color: #DF4601; }
	.pill.ravens.active  { background: #241773; border-color: #241773; }
</style>
