<script>
	import { Toaster } from '$lib/components/ui/sonner';
	import { ModeWatcher, setMode } from 'mode-watcher';
	import './app.css';

	setMode('light');
</script>

<main>
	<ModeWatcher />
	<Toaster />
	<slot />
</main>
