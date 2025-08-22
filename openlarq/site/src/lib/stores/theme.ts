import { writable } from 'svelte/store';
import { browser } from '$app/environment';

type Theme = 'light' | 'dark';

// get initial theme from localStorage or system preference
function getInitialTheme(): Theme {
	if (browser) {
		const stored = localStorage.getItem('theme') as Theme;
		if (stored) return stored;

		// Check system preference
		if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
			return 'dark';
		}
	}
	return 'light';
}

// create the store
export const theme = writable<Theme>(getInitialTheme());

// subscribe to theme changes and update localStorage and document
if (browser) {
	theme.subscribe((value) => {
		localStorage.setItem('theme', value);

		if (value === 'dark') {
			document.documentElement.classList.add('dark');
		} else {
			document.documentElement.classList.remove('dark');
		}
	});
}
