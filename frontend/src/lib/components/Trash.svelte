<script lang="ts">
	import { cn } from '@/utils';
	import Icon from '@iconify/svelte';
	import Draggable from './Draggable.svelte';

	const id = 'trash';
</script>

<Draggable
	draggableId={id}
	draggable={false}
	class={cn(
		'w-full flex flex-col justify-center items-center border border-black rounded-lg p-4 [&_*]:pointer-events-none',
		$$props.class
	)}
>
	<Icon icon="ph:trash-simple-fill" class="size-6" />
</Draggable>
