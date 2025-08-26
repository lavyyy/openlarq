<script lang="ts">
	import { theme } from '$lib/stores/theme';
	import Sun from '@lucide/svelte/icons/sun';
	import Moon from '@lucide/svelte/icons/moon';

	function toggleTheme() {
		$theme = $theme === 'light' ? 'dark' : 'light';
	}
</script>

<button
	on:click={toggleTheme}
	class="cursor-pointer rounded-lg p-2 transition-colors hover:bg-gray-100 dark:hover:bg-gray-800"
	aria-label="Toggle theme"
>
	{#if $theme === 'light'}
		<Moon class="h-5 w-5 text-gray-700 dark:text-gray-300" />
	{:else}
		<Sun class="h-5 w-5 text-gray-700 dark:text-gray-300" />
	{/if}
</button>
