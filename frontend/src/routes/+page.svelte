<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import Cashbox from '@/components/Cashbox.svelte';
	import Container from '@/components/Container.svelte';
	import Pocket from '@/components/Pocket.svelte';
	import Button from '@/components/ui/button/button.svelte';
</script>

<div class="space-y-4">
	<Header />

	<Container class="flex flex-row justify-between">
		<div>
			<h2 class="text-xs">Total</h2>
			<p class="font-semibold text-xl">$ 13500 THB</p>
		</div>
		<Button>อัพเดตยอดเงิน</Button>
	</Container>

	<Container class="space-y-8">
		<Cashbox amount={13000} />
		<div class="grid grid-cols-2 gap-8">
			<Pocket name="เงินเดือน" amount={11000} />
			<Pocket name="เงินออม" amount={2000} />
			<Pocket name="เที่ยวญี่ปุ่น" amount={5000} />
		</div>
	</Container>
</div>
