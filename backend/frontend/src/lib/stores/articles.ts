import { writable, derived } from 'svelte/store';

export type Team = 'orioles' | 'ravens' | 'both';
export type Source = 'rss' | 'gdelt' | 'reddit' | 'bluesky' | 'youtube';

export interface Article {
	id: string;
	title: string;
	url: string;
	summary: string;
	source: Source;
	team: Team;
	author?: string;
	image_url?: string;
	sentiment_score?: number;
	published_at: string;
	created_at: string;
}

export const articles     = writable<Article[]>([]);
export const teamFilter   = writable<Team | ''>('');
export const searchQuery  = writable('');
export const isLive       = writable(false);
export const newArticleCount = writable(0);

export const filteredArticles = derived(
	[articles, teamFilter, searchQuery],
	([$articles, $team, $q]) => {
		let list = $articles;
		if ($team) {
			list = list.filter(a => a.team === $team || a.team === 'both');
		}
		if ($q.trim()) {
			const lq = $q.toLowerCase();
			list = list.filter(a =>
				a.title.toLowerCase().includes(lq) ||
				(a.summary ?? '').toLowerCase().includes(lq)
			);
		}
		return list;
	}
);

export function prependArticle(article: Article) {
	articles.update(current => {
		if (current.some(a => a.id === article.id)) return current;
		newArticleCount.update(n => n + 1);
		return [article, ...current];
	});
}

export function clearNewCount() {
	newArticleCount.set(0);
}
