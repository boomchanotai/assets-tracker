<script lang="ts">
	import type { Currency } from '@/types';
	import { cn } from '@/utils';
	import { dragstart, dragover, dragenter, dragleave, drop } from '@/dragndrop';

	export let id: string;
	export let amount: number;
	export let currency: Currency = 'THB';
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<button
	draggable={true}
	class={cn(
		'w-full flex flex-col justify-start items-start border border-black rounded-lg p-4 [&_*]:pointer-events-none',
		$$props.class
	)}
	on:dragstart={(e) => dragstart(e, id)}
	on:dragover={dragover}
	on:dragenter={dragenter}
	on:dragleave={dragleave}
	on:drop={(e) => drop(e, id)}
>
	<h3>Cashbox</h3>
	<p>$ {amount} {currency}</p>
</button>
