<script>
	import { page } from '$app/stores';
</script>

<div class="error-page">
	<h1>{$page.status}</h1>
	<p>{$page.error?.message ?? 'Something went wrong.'}</p>
	<a href="/">← Back to home</a>
</div>

<style>
	.error-page {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		min-height: 60vh;
		gap: 1rem;
		color: #555;
	}
	h1 { font-size: 4rem; font-weight: 800; color: #111; }
	a { color: #241773; text-decoration: underline; }
</style>
