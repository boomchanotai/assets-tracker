<script lang="ts">
	import type { Currency, Pocket } from '@/types';
	import { cn } from '@/utils';
	import Draggable from './Draggable.svelte';

	export let pockets: Pocket[];
	export let id: string;
	export let name: string;
	export let amount: number;
	export let currency: Currency = 'THB';
</script>

<Draggable
	draggableId={id}
	draggable={true}
	{pockets}
	class={cn(
		'bg-white w-full h-full flex flex-col justify-start items-start border border-black rounded-lg p-4 [&_*]:pointer-events-none select-none',
		$$props.class
	)}
>
	<h3>{name}</h3>
	<p>$ {amount} {currency}</p>
</Draggable>
