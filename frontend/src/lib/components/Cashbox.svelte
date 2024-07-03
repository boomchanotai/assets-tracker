<script lang="ts">
	import type { Currency } from '@/types';
	import { cn } from '@/utils';

	export let amount: number;
	export let currency: Currency = 'THB';
</script>

<div class={cn('border border-black rounded-lg p-4', $$props.class)}>
	<h3>Cashbox</h3>
	<p>$ {amount} {currency}</p>
</div>
