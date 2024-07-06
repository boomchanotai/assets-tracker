<script lang="ts">
	import type { Currency } from '@/types';

	export let amount: number;
	export let currency: Currency = 'THB';
</script>

<div>
	<h2 class="text-xs">Total</h2>
	<p class="font-semibold text-xl">$ {amount} {currency}</p>
</div>
