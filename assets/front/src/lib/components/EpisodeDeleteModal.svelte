<script>
	import Modal from "$lib/components/Modal.svelte";
	import { createEventDispatcher } from "svelte";
	
	let dispatch = createEventDispatcher();
	
	export let show = false;
	export let episode;
	
	const confirmDelete = () => {
		dispatch("deleteEpisode", { 
			id: episode.ID, 
			name: episode.Name 
		});
		show = false;
	}
	
	const cancelDelete = () => {
		show = false;
	}
</script>

<Modal show={show} on:close={cancelDelete}>
	<div slot="title">
		<i class="fas fa-trash"></i>
		Delete Episode
	</div>
	<div slot="body">
		<div class="delete-section">
			<div class="warning-message">
				<i class="fas fa-exclamation-triangle"></i>
				<p>Are you sure you want to delete this episode?</p>
			</div>
			
			<div class="episode-info">
				<h4 class="episode-title">{episode?.Name}</h4>
				<p class="episode-details">Episode {episode?.Number}</p>
			</div>
			
			<div class="warning-notice">
				<i class="fas fa-info-circle"></i>
				<span>This action cannot be undone.</span>
			</div>
		</div>
	</div>
	<div slot="footer">
		<div class="footer-buttons">
			<button class="btn btn-secondary" on:click={cancelDelete}>
				<i class="fas fa-times"></i>
				Cancel
			</button>
			<button class="btn btn-danger" on:click={confirmDelete}>
				<i class="fas fa-trash"></i>
				Delete Episode
			</button>
		</div>
	</div>
</Modal>

<style>
	.delete-section {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}
	
	.warning-message {
		display: flex;
		align-items: flex-start;
		gap: 0.75rem;
		padding: 1rem;
		background-color: rgba(231, 76, 60, 0.1);
		border: 1px solid rgba(231, 76, 60, 0.3);
		border-radius: 8px;
	}
	
	.warning-message i {
		color: #e74c3c;
		font-size: 1.2rem;
		margin-top: 0.1rem;
		flex-shrink: 0;
	}
	
	.warning-message p {
		margin: 0;
		color: var(--text-color);
		font-weight: 500;
		line-height: 1.4;
	}
	
	.episode-info {
		padding: 1rem;
		background-color: rgba(255, 255, 255, 0.05);
		border-radius: 8px;
		border-left: 3px solid var(--highlight, #3498db);
	}
	
	.episode-title {
		margin: 0 0 0.5rem 0;
		color: var(--text-color);
		font-size: 1.1rem;
		font-weight: 600;
	}
	
	.episode-details {
		margin: 0;
		color: var(--text-color);
		opacity: 0.8;
		font-size: 0.9rem;
	}
	
	.warning-notice {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.75rem;
		background-color: rgba(241, 196, 15, 0.1);
		border: 1px solid rgba(241, 196, 15, 0.3);
		border-radius: 6px;
		font-size: 0.9rem;
	}
	
	.warning-notice i {
		color: #f1c40f;
		flex-shrink: 0;
	}
	
	.warning-notice span {
		color: var(--text-color);
		font-weight: 500;
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
	
	.btn-danger {
		background-color: #e74c3c;
		color: white;
		border: 2px solid #e74c3c;
	}
	
	.btn-danger:hover:not(:disabled) {
		background-color: #c0392b;
		border-color: #c0392b;
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(231, 76, 60, 0.3);
	}
	
	/* Mobile responsiveness */
	@media (max-width: 768px) {
		.footer-buttons {
			flex-direction: column-reverse;
		}
		
		.btn {
			width: 100%;
		}
		
		.delete-section {
			gap: 1rem;
		}
		
		.warning-message, .episode-info, .warning-notice {
			padding: 0.75rem;
		}
	}
</style>
