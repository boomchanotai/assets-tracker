<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Select from '$lib/components/ui/select';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import type { SelectOptions } from '@/types';
	import Icon from '@iconify/svelte';

	import {
		banks,
		bankAccountTypes,
		financialAccountType,
		securitiesCompanies,
		mutualFundCompanies
	} from '$lib/constants/bank';
	import { useCreateAccountMutation } from '@/hook/mutation/account';
	import { toast } from 'svelte-sonner';

	let selectedAccountType: SelectOptions | undefined;
	let name: string;
	let selectedBank: SelectOptions | undefined;

	const createAccountMutation = useCreateAccountMutation();
	const handleSubmit = (event: SubmitEvent) => {
		event.preventDefault();
		if (!selectedAccountType) {
			toast.error('Please select account type');
			return;
		}

		if (!selectedBank) {
			toast.error('Please select bank');
			return;
		}

		$createAccountMutation.mutate({
			type: selectedAccountType.value,
			name,
			bank: selectedBank.value
		});
	};
</script>

<Dialog.Root>
	<Dialog.Trigger>
		<Button class="aspect-square">
			<Icon icon="ph:plus-bold" />
		</Button>
	</Dialog.Trigger>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>Add New Account</Dialog.Title>
		</Dialog.Header>
		<form on:submit={handleSubmit} class="grid gap-4 py-4">
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
						<Input id="ชื่อบัญชี" placeholder="Account Name" class="col-span-3" bind:value={name} />
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
							bind:value={name}
						/>
					</div>
				{/if}

				{#if bankAccountTypes.filter(({ value }) => value === selectedAccountType?.value).length > 0}
					<div class="grid grid-cols-4 items-center gap-4">
						<Label for="username" class="text-right">Bank</Label>
						<Select.Root bind:selected={selectedBank}>
							<Select.Trigger class="col-span-3">
								<Select.Value placeholder="ธนาคาร" />
							</Select.Trigger>
							<Select.Content>
								{#each banks as { label, value }}
									<Select.Item {value}>{label}</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					</div>
				{:else if selectedAccountType?.value === 'mutual fund account'}
					<div class="grid grid-cols-4 items-center gap-4">
						<Label for="username" class="text-right">Mutual Fund</Label>
						<Select.Root bind:selected={selectedBank}>
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
						<Select.Root bind:selected={selectedBank}>
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

			<Button type="submit" class="gap-2">
				<Icon icon="ph:plus-bold" />
				<span>Add New Account</span>
			</Button>
		</form>
	</Dialog.Content>
</Dialog.Root>
