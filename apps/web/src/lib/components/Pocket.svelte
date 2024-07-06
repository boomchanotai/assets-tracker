<script lang="ts">
	import type { Currency } from '@/types';
	import { cn } from '@/utils';
	import Draggable from './Draggable.svelte';
	type Props = {
		id: string;
		name: string;
		amount: number;
		currency?: Currency;
		hidden?: boolean;
	};

	export let id: Props['id'];
	export let name: Props['name'];
	export let amount: Props['amount'];
	export let currency: Props['currency'] = 'THB';
</script>

<Draggable
	draggableId={id}
	draggable={true}
	class={cn(
		'bg-white w-full h-full flex flex-col justify-start items-start border border-black rounded-lg p-4 [&_*]:pointer-events-none select-none',
		$$props.class
	)}
>
	<h3>{name}</h3>
	<p>$ {amount} {currency}</p>
</Draggable>
