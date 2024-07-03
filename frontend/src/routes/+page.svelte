<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import Cashbox from '@/components/Cashbox.svelte';
	import Container from '@/components/Container.svelte';
	import Pocket from '@/components/Pocket.svelte';
	import Button from '@/components/ui/button/button.svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import Balance from '@/components/Balance.svelte';
</script>

<div class="space-y-4">
	<Header />

	<Container class="flex flex-row justify-between">
		<Balance amount={14500} />
		<Dialog.Root>
			<Dialog.Trigger>
				<Button>อัพเดตยอดเงิน</Button>
			</Dialog.Trigger>
			<Dialog.Content class="sm:max-w-[425px]">
				<Dialog.Header>
					<Dialog.Title>Update Balance</Dialog.Title>
				</Dialog.Header>
				<div>
					<Input id="balance" type="number" placeholder="Balance" />
				</div>
				<Dialog.Footer>
					<Button type="submit">Save changes</Button>
				</Dialog.Footer>
			</Dialog.Content>
		</Dialog.Root>
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
