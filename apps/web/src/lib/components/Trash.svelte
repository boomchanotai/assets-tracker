<script lang="ts">
	import { cn } from '@/utils';
	import Icon from '@iconify/svelte';
	import Draggable from './Draggable.svelte';
	import type { Pocket } from '$lib/types';

	export let accountId: string;
	export let pockets: Pocket[];
	const id = 'TRASH';
</script>

<Draggable
	{accountId}
	draggableId={id}
	draggable={false}
	{pockets}
	class={cn(
		'w-full h-full flex flex-col justify-center items-center border border-black rounded-lg p-4 [&_*]:pointer-events-none select-none',
		$$props.class
	)}
>
	<Icon icon="ph:trash-simple-fill" class="size-6" />
</Draggable>
