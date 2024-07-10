<script lang="ts">
	import { QueryClientProvider, QueryClient } from '@tanstack/svelte-query';
	import { Toaster } from '$lib/components/ui/sonner';
	import { ModeWatcher, setMode } from 'mode-watcher';
	import './app.css';

	setMode('light');

	const queryClient = new QueryClient();
</script>

<QueryClientProvider client={queryClient}>
	<ModeWatcher />
	<Toaster />
	<slot />
</QueryClientProvider>
