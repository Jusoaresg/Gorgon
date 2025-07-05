<script>
import { createEventDispatcher, tick } from 'svelte';
import { fade } from 'svelte/transition';
import EpisodeChangeTrackingModal from '$lib/components/EpisodeChangeTrackingModal.svelte';
import EpisodeDeleteModal from '$lib/components/EpisodeDeleteModal.svelte';
import EpisodeManualSearchModal from '$lib/components/EpisodeManualSearchModal.svelte';

//Styles
import '$lib/styles/components/episode-card.css';

export let episode;

const dispatch = createEventDispatcher();

let episodeModal = false;
let showSummary = false;
let showDeleteConfirm = false;
let episodeManualSearch = false;

$: status = episode.Tracking;
$: isUpcoming = episode.AirStamp && (episode.AirStamp * 1000) > Date.now();

function toggleSummary() {
	showSummary = !showSummary;
}

async function deleteEpisode() {
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

function formatDateRelative(dateStr) {
	if (!dateStr) return '';
	const date = new Date(dateStr * 1000);
	const now = new Date();
	const diffTime = date - now;
	const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
	
	if (diffDays === 0) return 'Today';
	if (diffDays === 1) return 'Tomorrow';
	if (diffDays === -1) return 'Yesterday';
	if (diffDays > 0) return `In ${diffDays} days`;
	if (diffDays < 0) return `${Math.abs(diffDays)} days ago`;
	
	return '';
}

function getDownloadPath(episode) {
	return episode.DownloadPath || episode.FilePath || '';
}
</script>

<div class="episode-card" data-episode-id={episode.ID} data-status={status}>
	<div class={"episode-number " + getStatusClass(status)}>
		<span>{episode.Number}</span>
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
				<div class="episode-date">
					<span class="air-date">
						<i class="fas fa-calendar-alt"></i>
						{formatDate(episode.AirStamp)}
						{#if formatDateRelative(episode.AirStamp)}
							<span class="relative-date">({formatDateRelative(episode.AirStamp)})</span>
						{/if}
						{#if isUpcoming}
							<span class="upcoming-badge">Upcoming</span>
						{/if}
					</span>
				</div>
			{/if}
			
			{#if status.toLowerCase() === 'downloaded' && getDownloadPath(episode)}
				<div class="download-location">
					<i class="fas fa-folder"></i>
					<span class="download-path" title={getDownloadPath(episode)}>
						{getDownloadPath(episode)}
					</span>
				</div>
			{/if}
			
			{#if episode.Summary && showSummary}
				<hr class="summary-divider" />
				<div class="episode-summary" transition:fade={{duration: 200}}>
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
