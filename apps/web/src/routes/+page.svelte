<script lang="ts">
	import Container from '@/components/Container.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { useLoginMutation } from '@/hook/auth';
	import { authStore } from '@/store/auth';
	import { readable } from 'svelte/store';
	import { goto } from '$app/navigation';

	let email = '';
	let password = '';

	authStore.subscribe((value) => {
		if (value.accessToken) {
			goto('/account');
		}
	});

	const loginMutation = useLoginMutation();
	const handleSubmit = async (event: SubmitEvent) => {
		event.preventDefault();
		$loginMutation.mutate({ email, password });
	};
</script>

<div class="space-y-4">
	<Container class="h-svh flex flex-col justify-center items-center">
		<Card.Root class="w-[350px]">
			<Card.Header>
				<Card.Title class="text-center">Sign in</Card.Title>
				<Card.Description class="text-center">Sign in to your account to continue</Card.Description>
			</Card.Header>
			<Card.Content>
				<form on:submit={handleSubmit}>
					<div class="grid w-full items-center gap-4">
						<div class="flex flex-col space-y-1.5">
							<Label for="name">Email</Label>
							<Input id="email" type="email" bind:value={email} placeholder="Email | อีเมล" />
						</div>
						<div class="flex flex-col space-y-1.5">
							<Label for="name">Password</Label>
							<Input
								id="password"
								type="password"
								bind:value={password}
								placeholder="Password | พาสเวิร์ด"
							/>
						</div>
						<div>
							<Button type="submit">Sign in</Button>
						</div>
					</div>
				</form>
			</Card.Content>
		</Card.Root>
	</Container>
</div>
