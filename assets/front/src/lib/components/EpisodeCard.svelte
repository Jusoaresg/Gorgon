<script>
import { createEventDispatcher, onMount, tick } from 'svelte';
import { slide, fade, fly, scale } from 'svelte/transition';
import EpisodeChangeTrackingModal from '$lib/components/EpisodeChangeTrackingModal.svelte';
import EpisodeDeleteModal from '$lib/components/EpisodeDeleteModal.svelte';
import EpisodeManualSearchModal from '$lib/components/EpisodeManualSearchModal.svelte';

//Styles
import '$lib/styles/components/episode-card.css';

export let episode;
export let seasonNumber;

const dispatch = createEventDispatcher();

let episodeModal = false;
let showSummary = false;
let showDeleteConfirm = false; // New state for delete confirmation
let episodeManualSearch = false;

$: status = episode.Tracking;

const statusOptions = ['wanted', 'skipped'];

function toggleSummary() {
	showSummary = !showSummary;
}

async function deleteEpisode(event) {
	dispatch('deleteEpisode', { id: episode.ID });
	showDeleteConfirm = false;
}

function getStatusClass(status) {
	switch (status.toLowerCase()) {
		case 'downloaded': return 'status-downloaded';
		case 'missing': return 'status-missing';
		case 'wanted': return 'status-wanted';
		case 'skipped': return 'status-skipped';
		case 'snatched': return 'status-snatched';
		default: return 'status-unknown';
	}
}

function getStatusIcon(status) {
	switch (status.toLowerCase()) {
		case 'downloaded': return 'fas fa-check';
		case 'missing': return 'fas fa-times';
		case 'wanted': return 'fas fa-clock';
		case 'skipped': return 'fas fa-ban';
		case 'snatched': return 'fas fa-magnet';
		default: return 'fas fa-question';
	}
}

function formatDate(dateStr) {
	if (!dateStr) return '';
	const date = new Date(dateStr * 1000);
	return date.toLocaleDateString('pt-BR', {
		year: 'numeric',
		month: 'long',
		day: 'numeric'
	});
}
</script>

<div class="episode-card" data-episode-id={episode.ID} data-status={status}>
	<div class={"episode-number " + getStatusClass(status)}>
		<span>{episode.Number}</span>
		<div class="absolute-number">
			{(seasonNumber - 1) * 10 + episode.Number}
		</div>
	</div>
	<div class="episode-content">
		<div class="episode-header">
			<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
			<!-- svelte-ignore a11y_click_events_have_key_events -->
			<h3 class="episode-title" on:click={toggleSummary} style="cursor: pointer;">
				{episode.Name}
				{#if episode.Summary}
					<i class="fas fa-chevron-{showSummary ? 'up' : 'down'}" style="margin-left: 8px; font-size: 0.8em; opacity: 0.7;"></i>
				{/if}
			</h3>
			<div class="episode-actions">

				<!-- Delete button - only show for downloaded episodes -->
				{#if status.toLowerCase() === 'downloaded'}
					<button 
						aria-label="delete episode"
						class="episode-btn delete-btn" 
						title="Delete Episode"
						on:click={async () => {
							showDeleteConfirm = false;
							await tick();
							showDeleteConfirm = true;
						}}
					>
						<i class="fas fa-trash"></i>
					</button>
				{/if}

				<button 
					aria-hidden="true" 
					class="episode-btn force-search-btn" 
					title="Force Search"
				>
					<i class="fas fa-search"></i>
				</button>
				<button 
					aria-hidden="true" 
					class="episode-btn manual-search-btn" 
					title="Manual Search"
					on:click={async () => {
						episodeManualSearch = false;
						await tick();
						episodeManualSearch = true;
					}}
				>
					<i class="fas fa-list"></i>
				</button>
				
				<button 
					aria-label="tracking"
					class={"episode-btn episode-status-toggle " + getStatusClass(status)} 
					title="Change Tracking Status"
					on:click={async () => {
						episodeModal = false;
						await tick();
						episodeModal = true;
					}}
				>
					<i class={getStatusIcon(status)}></i>
				</button>
			</div>
		</div>
		<div class="episode-info">
			{#if episode.AirStamp}
				<div class="episode-date mt-1">
					<span class="air-date">
						<i class="fas fa-calendar-alt"></i>
						{formatDate(episode.AirStamp)}
					</span>
				</div>
			{/if}
			{#if episode.Summary && showSummary}
				<div class="episode-summary mt-2" transition:fade={{duration: 200}}>
					{@html episode.Summary}
				</div>
			{/if}
		</div>
	</div>
</div>

<EpisodeManualSearchModal
	show={episodeManualSearch}
	episode={episode}
	on:close={() => { episodeManualSearch = false }}
/>

<EpisodeDeleteModal 
	show={showDeleteConfirm} 
	episode={episode} 
	on:deleteEpisode={deleteEpisode}
	on:close={() => { showDeleteConfirm = false }}
/>

<EpisodeChangeTrackingModal
	show={episodeModal}
	episode={episode}
	on:close={() => { episodeModal = false }}
/>

<style>
.episode-title:hover {
	color: var(--highlight);
}

.delete-btn {
	background-color: rgba(231, 76, 60, 0.2);
	color: #e74c3c;
}

.delete-btn:hover {
	background-color: rgba(231, 76, 60, 0.3);
	color: #fff;
}

.modal-overlay {
	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	background-color: rgba(0, 0, 0, 0.7);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 10000;
}

.delete-modal {
	background-color: var(--card-bg);
	border-radius: 12px;
	border: 2px solid rgba(255, 255, 255, 0.1);
	box-shadow: 0 20px 40px rgba(0, 0, 0, 0.6);
	max-width: 400px;
	width: 90%;
	max-height: 90vh;
	overflow: hidden;
}

.modal-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 1.5rem;
	border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-close {
	background: none;
	border: none;
	color: var(--text-color);
	font-size: 1.25rem;
	cursor: pointer;
	padding: 0.25rem;
	border-radius: 4px;
	transition: background-color 0.2s ease;
}

.modal-close:hover {
	background-color: rgba(255, 255, 255, 0.1);
}

.modal-body {
	padding: 1.5rem;
}

.episode-info-text {
	background-color: rgba(255, 255, 255, 0.05);
	padding: 0.75rem;
	border-radius: 6px;
	border-left: 3px solid var(--highlight);
}

.warning-text {
	color: #e74c3c;
	font-weight: 500;
	margin-bottom: 0 !important;
}

.modal-actions {
	display: flex;
	gap: 0.75rem;
	padding: 1.5rem;
	border-top: 1px solid rgba(255, 255, 255, 0.1);
	justify-content: flex-end;
}

.btn {
	padding: 0.75rem 1.5rem;
	border: none;
	border-radius: 6px;
	cursor: pointer;
	font-size: 0.9rem;
	font-weight: 500;
	transition: all 0.2s ease;
	min-width: 80px;
}

.btn-secondary {
	background-color: rgba(255, 255, 255, 0.1);
	color: var(--text-color);
}

.btn-secondary:hover {
	background-color: rgba(255, 255, 255, 0.2);
}

.btn-danger {
	background-color: #e74c3c;
	color: white;
}

.btn-danger:hover {
	background-color: #c0392b;
	transform: translateY(-1px);
}

@media (max-width: 768px) {
	.delete-modal {
		margin: 1rem;
		width: calc(100% - 2rem);
	}
	
	.modal-actions {
		flex-direction: column-reverse;
	}
	
	.btn {
		width: 100%;
	}
}
</style>
