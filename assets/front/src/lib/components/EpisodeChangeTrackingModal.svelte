<script>
	import Modal from "$lib/components/Modal.svelte";
	import { createEventDispatcher } from "svelte";
	import { changeTrackingStatus } from "$lib/js/EpisodeFetchs";

	let dispatch = createEventDispatcher();
	
	export let show = false;
	export let episode;

	let newStatus = ''

	const statusOptions = [
		{ value: 'wanted', label: 'Wanted', icon: 'fas fa-clock', color: '#f39c12' },
		{ value: 'skipped', label: 'Skipped', icon: 'fas fa-ban', color: '#7f8c8d' }
	];

	const applyChanges = () => {
		changeTrackingStatus(episode.ID, newStatus)
		dispatch("close")
	}

	$: currentStatus = episode?.Tracking ?? '';

	$: if (show) {
		newStatus = '';
		currentStatus = episode?.Tracking;
	}

</script>

<Modal show={show} on:close={() => {
	show = false;
	newStatus = '';
}
}>
	<div slot="title">
		<i class="fas fa-tasks"></i>
		Edit Episode Status
	</div>

	<div slot="body">
		<!-- Status Selection -->
		<div class="status-section">
			<h3 class="section-title">
				<i class="fas fa-flag"></i>
				Change Status To:
			</h3>
			<div class="status-selector">
				{#each statusOptions as option}
					<button
						type="button"
						class="status-option {newStatus === option.value ? 'selected' : ''}"
						style="--status-color: {option.color}"
						on:click={() => newStatus = option.value}
						disabled={currentStatus === option.value || currentStatus === 'downloaded'}
					>
						<i class={option.icon}></i>
						{option.label}
					</button>
				{/each}
			</div>
		</div>
	</div>

	<div slot="footer">
		<div class="footer-buttons">
			<button class="btn btn-secondary" on:click={() => {show = false}}>
				<i class="fas fa-times"></i>
				Cancel
			</button>
			<button 
				class="btn btn-primary" 
				on:click={applyChanges}
				disabled={newStatus === ''}
			>
				<i class="fas fa-check"></i>
				Apply Changes
			</button>
		</div>

	</div>
</Modal>

<style>
	.section-title {
		margin: 0 0 1rem 0;
		color: var(--text-color);
		font-size: 1.1rem;
		font-weight: 500;
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}
	
	.section-title i {
		color: var(--highlight, #3498db);
		width: 16px;
		text-align: center;
	}

	.status-section {
		margin-bottom: 2rem;
	}

	.status-selector {
		display: flex;
		gap: 1rem;
	}

	.status-option {
		flex: 1;
		padding: 1rem;
		border: 2px solid rgba(255, 255, 255, 0.1);
		border-radius: 8px;
		background-color: rgba(255, 255, 255, 0.05);
		color: var(--text-color);
		cursor: pointer;
		transition: all 0.2s ease;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		font-size: 1rem;
		font-weight: 500;
		text-transform: capitalize;
	}
	
	.status-option:hover {
		background-color: rgba(255, 255, 255, 0.1);
		transform: translateY(-2px);
	}
	
	.status-option.selected {
		border-color: var(--status-color);
		background-color: var(--status-color);
		color: white;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
	}

	.footer-buttons {
		display: flex;
		justify-content: flex-end;
		gap: 1rem;
	}

	.btn {
		padding: 0.75rem 1.5rem;
		border: none;
		border-radius: 8px;
		font-size: 0.9rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
		display: flex;
		align-items: center;
		gap: 0.5rem;
		min-width: 120px;
		justify-content: center;
	}
	
	.btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
		transform: none !important;
	}
	
	.btn-secondary {
		background-color: rgba(255, 255, 255, 0.1);
		color: var(--text-color);
		border: 2px solid rgba(255, 255, 255, 0.1);
	}
	
	.btn-secondary:hover:not(:disabled) {
		background-color: rgba(255, 255, 255, 0.2);
		transform: translateY(-1px);
	}
	
	.btn-primary {
		background-color: var(--highlight, #3498db);
		color: white;
		border: 2px solid var(--highlight, #3498db);
	}
	
	.btn-primary:hover:not(:disabled) {
		background-color: #2980b9;
		border-color: #2980b9;
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(52, 152, 219, 0.3);
	}


</style>
