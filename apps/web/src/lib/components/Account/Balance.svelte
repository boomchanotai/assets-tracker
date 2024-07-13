<script lang="ts">
	import type { Currency } from '@/types';
	import Container from '../Container.svelte';
	import UpdateBalance from '../dialog/AddBalance.svelte';

	export let accountId: string;
	export let amount: number = 0;
	export let currency: Currency = 'THB';
</script>

<Container class="flex flex-row justify-between">
	<div>
		<h2 class="text-xs">Total</h2>
		<p class="font-semibold text-xl">$ {amount} {currency}</p>
	</div>
	<div>
		<UpdateBalance {accountId} />
	</div>
</Container>
