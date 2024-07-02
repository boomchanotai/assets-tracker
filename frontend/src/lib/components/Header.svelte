<script lang="ts">
	import scb from '$lib/images/scb.png';
	import Icon from '@iconify/svelte';

	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import * as Select from '$lib/components/ui/select';

	import { bank, accountTypes, financialAccountType } from '$lib/constants/bank';

	let type: string | null = null;
</script>

<header class="py-4 px-6 space-y-4">
	<h1 class="text-xl font-semibold">Cloud Pocket</h1>
	<div class="flex items-center gap-4">
		<div
			class="flex flex-row gap-6 items-center border border-black px-4 py-2 rounded-lg hover:bg-black/10 transition-colors duration-150"
		>
			<div class="flex flex-row gap-2 items-center">
				<div><img src={scb} alt="" class="size-6 rounded-full" /></div>
				<div>
					<p class="text-sm font-medium">442-961089-7</p>
					<p class="text-xs">ธนาคารไทยพาณิชย์</p>
				</div>
			</div>
			<div><Icon icon="ph:caret-down-bold" /></div>
		</div>
		<Dialog.Root>
			<Dialog.Trigger class="border-black border p-4 rounded-lg">
				<Icon icon="ph:plus-bold" />
			</Dialog.Trigger>
			<Dialog.Content class="sm:max-w-[425px]">
				<Dialog.Header>
					<Dialog.Title>Add New Account</Dialog.Title>
				</Dialog.Header>
				<div class="grid gap-4 py-4">
					<div class="grid grid-cols-4 items-center gap-4">
						<Label for="username" class="text-right">Type</Label>
						<Select.Root>
							<Select.Trigger class="col-span-3">
								<Select.Value placeholder="ประเภทบัญชี" />
							</Select.Trigger>
							<Select.Content>
								{#each [...accountTypes, ...financialAccountType] as { label, value }}
									<Select.Item {value}>{label}</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					</div>
					{type}
					<div class="grid grid-cols-4 items-center gap-4">
						<Label for="name" class="text-right">Account No.</Label>
						<Input
							id="เลขบัญชี"
							placeholder="xxx-x-xxxxx-x"
							class="col-span-3"
							minlength={10}
							maxlength={10}
						/>
					</div>
					<div class="grid grid-cols-4 items-center gap-4">
						<Label for="username" class="text-right">Bank</Label>
						<Select.Root>
							<Select.Trigger class="col-span-3">
								<Select.Value placeholder="ธนาคาร" />
							</Select.Trigger>
							<Select.Content>
								{#each bank as { label, value }}
									<Select.Item {value}>{label}</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					</div>
				</div>
				<Dialog.Footer>
					<Button type="submit">Save changes</Button>
				</Dialog.Footer>
			</Dialog.Content>
		</Dialog.Root>
	</div>
</header>
