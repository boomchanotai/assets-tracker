<script lang="ts">
	import scb from '$lib/images/scb.png';
	import Icon from '@iconify/svelte';

	import Container from '$lib/components/Container.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import * as Select from '$lib/components/ui/select';

	import {
		bank,
		bankAccountTypes,
		financialAccountType,
		securitiesCompanies,
		mutualFundCompanies
	} from '$lib/constants/bank';
	import type { SelectOptions } from '@/types';

	let selectedAccountType: SelectOptions | undefined;
</script>

<Container class="py-4 space-y-4">
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
			<Dialog.Trigger class="h-full">
				<Button class="">
					<Icon icon="ph:plus-bold" />
				</Button>
			</Dialog.Trigger>
			<Dialog.Content class="sm:max-w-[425px]">
				<Dialog.Header>
					<Dialog.Title>Add New Account</Dialog.Title>
				</Dialog.Header>
				<div class="grid gap-4 py-4">
					<div class="grid grid-cols-4 items-center gap-4">
						<Label for="username" class="text-right">Type</Label>
						<Select.Root bind:selected={selectedAccountType}>
							<Select.Trigger class="col-span-3">
								<Select.Value placeholder="ประเภทบัญชี" />
							</Select.Trigger>
							<Select.Content class="max-h-64 overflow-y-auto">
								{#each [...bankAccountTypes, ...financialAccountType] as { label, value }}
									<Select.Item {value}>{label}</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					</div>

					{#if selectedAccountType !== undefined}
						{#if financialAccountType.filter(({ value }) => value === selectedAccountType?.value).length > 0}
							<div class="grid grid-cols-4 items-center gap-4">
								<Label for="name" class="text-right">Account Name</Label>
								<Input id="ชื่อบัญชี" placeholder="Account Name" class="col-span-3" />
							</div>
						{:else}
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
						{/if}

						{#if bankAccountTypes.filter(({ value }) => value === selectedAccountType?.value).length > 0}
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
						{:else if selectedAccountType?.value === 'mutual fund account'}
							<div class="grid grid-cols-4 items-center gap-4">
								<Label for="username" class="text-right">Mutual Fund</Label>
								<Select.Root>
									<Select.Trigger class="col-span-3">
										<Select.Value placeholder="บริษัทหลักทรัพย์จัดการกองทุน" />
									</Select.Trigger>
									<Select.Content class="max-h-64 overflow-y-auto">
										{#each mutualFundCompanies as { label, value }}
											<Select.Item {value}>{label}</Select.Item>
										{/each}
									</Select.Content>
								</Select.Root>
							</div>
						{:else}
							<div class="grid grid-cols-4 items-center gap-4">
								<Label for="username" class="text-right">Securities Companies</Label>
								<Select.Root>
									<Select.Trigger class="col-span-3">
										<Select.Value placeholder="บริษัทหลักทรัพย์" />
									</Select.Trigger>
									<Select.Content class="max-h-64 overflow-y-auto">
										{#each securitiesCompanies as { label, value }}
											<Select.Item {value}>{label}</Select.Item>
										{/each}
									</Select.Content>
								</Select.Root>
							</div>
						{/if}
					{/if}
				</div>
				<Dialog.Footer>
					<Button type="submit" class="gap-2">
						<Icon icon="ph:plus-bold" />
						<span>Add New Account</span>
					</Button>
				</Dialog.Footer>
			</Dialog.Content>
		</Dialog.Root>
	</div>
</Container>
