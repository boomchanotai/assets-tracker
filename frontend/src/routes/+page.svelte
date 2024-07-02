<script>
	import Header from '$lib/components/Header.svelte';
	import Container from '@/components/Container.svelte';
	import Pocket from '@/components/Pocket.svelte';
</script>

<div class="space-y-4">
	<Header />

	<Container class="grid grid-cols-2 gap-8">
		<Pocket name="เงินเดือน" amount={11000} />
		<Pocket name="เงินออม" amount={2000} />
		<Pocket name="เที่ยวญี่ปุ่น" amount={5000} />
	</Container>
</div>
