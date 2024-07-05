<script lang="ts">
	import type { Currency } from '@/types';
	import { cn } from '@/utils';
	import { dragstart, dragover, dragenter, dragleave, drop } from '@/dragndrop';

	type Props = {
		id: string;
		name: string;
		amount: number;
		currency?: Currency;
	};

	export let id: Props['id'];
	export let name: Props['name'];
	export let amount: Props['amount'];
	export let currency: Props['currency'] = 'THB';
</script>

<button
	draggable={true}
	aria-grabbed={true}
	class={cn(
		'w-full h-full flex flex-col justify-start items-start border border-black aspect-square rounded-lg p-4 [&_*]:pointer-events-none',
		$$props.class
	)}
	on:dragstart={(e) => dragstart(e, id)}
	on:dragover={dragover}
	on:dragenter={dragenter}
	on:dragleave={dragleave}
	on:drop={(e) => drop(e, id)}
>
	<h3>{name}</h3>
	<p>$ {amount} {currency}</p>
</button>
