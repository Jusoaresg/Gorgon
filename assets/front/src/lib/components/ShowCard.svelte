<script>
import { createEventDispatcher } from "svelte";
import { fly, scale, fade } from 'svelte/transition';
import { quintOut, quintIn } from 'svelte/easing';

export let show;
export let nextEpisode = null; // New prop for upcoming episode data
export let showUpcoming = false;
export let i = 0; // Animation index

const dispatch = createEventDispatcher();

function getStatusColor(status) {
	switch(status) {
		case 'Running': return '#00FF7F';
		case 'Ended': return '#FF6B6B';
		case 'Missing': return '#FFD93D';
		default: return '#999';
	}
}

function formatNextEpisode(episode) {
	if (!episode) return null;
	
	const airDate = new Date(episode.AirStamp);
	const now = new Date();
	const diffTime = airDate - now;
	const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
	const hasAlreadyAired = airDate < now;
	
	// Format the date
	const dateOptions = { 
		month: 'short', 
		day: 'numeric',
		hour: '2-digit',
		minute: '2-digit'
	};
	
	let timeText = '';
	if (diffDays === 0 && !hasAlreadyAired) {
		timeText = 'Today';
	} else if (diffDays <= 0 && hasAlreadyAired) {
		timeText = 'Aired';
	} else if (diffDays === 1) {
		timeText = 'Tomorrow';
	} else if (diffDays <= 7) {
		timeText = `${diffDays} days`;
	} else {
		//TODO: Change this Locale Date String to the user config
		timeText = airDate.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
	}
	
	return {
		episode: `S${episode.Season?.toString().padStart(2, '0')}E${episode.Number?.toString().padStart(2, '0')}`,
		name: episode.Name,
		timeText,
		airDate: airDate.toLocaleDateString('en-US', dateOptions),
		isToday: diffDays === 0 && !hasAlreadyAired,
		isAired: hasAlreadyAired && diffDays <= 0,
		isTomorrow: diffDays === 1,
		isThisWeek: diffDays > 1 && diffDays <= 7
	};
}

async function onClick() {
	dispatch("select", show.ID);
}

$: episodeInfo = showUpcoming && nextEpisode ? formatNextEpisode(nextEpisode) : null;
</script>

<div
	class="show-card"
	class:upcoming-mode={showUpcoming}
	aria-hidden="true"
	on:click={onClick}
	in:fly={{ y: 30, duration: 250, delay: i * 20, easing: quintOut }}
	out:fly={{ x: -200, duration: 250, easing: quintIn }}
>
	<div class="show-poster">
		{#if show.ImageOriginal}
			<img src={show.ImageMedium} alt={show.Name} loading="lazy"
				in:scale={{ duration: 200, delay: 100 }}>
		{:else}
			<div class="no-image" in:fade={{ duration: 200 }}>
				<span>No Image</span>
			</div>
		{/if}
		
		<div 
			class="show-status" 
			style="background-color: {getStatusColor(show.Status)}"
			in:scale={{ duration: 150, delay: 200 }}
		></div>

		{#if episodeInfo}
			<div class="next-episode-badge" 
				class:today={episodeInfo.isToday}
				class:tomorrow={episodeInfo.isTomorrow}
				class:this-week={episodeInfo.isThisWeek}
				class:aired={episodeInfo.isAired}
				in:scale={{ duration: 200, delay: 250 }}>
				{episodeInfo.timeText}
			</div>
		{/if}
	</div>
	
	<div class="show-info" in:fly={{ y: 10, duration: 200, delay: 150 }}>
		<div class="show-title">{show.Name}</div>
		
		{#if showUpcoming && episodeInfo}
			<div class="episode-info" in:fade={{ duration: 200, delay: 200 }}>
				<div class="episode-number">{episodeInfo.episode}</div>
				<div class="episode-name">{episodeInfo.name || 'Episode ' + episodeInfo.episode}</div>
				<div class="episode-time">{episodeInfo.airDate}</div>
			</div>
		{:else}
			<div class="show-meta">
				<span>{show.Premiered ? show.Premiered : 'N/A'}</span>
				<span>{show.Rating || 'N/A'}</span>
			</div>
		{/if}
	</div>
</div>

<style>
	/* Enhanced card for upcoming mode */
	.show-card.upcoming-mode {
		border: 1px solid rgba(125, 0, 163, 0.3);
	}

	.show-card.upcoming-mode:hover {
		border-color: var(--highlight);
		box-shadow: 0 6px 20px rgba(125, 0, 163, 0.4);
	}

	.next-episode-badge {
		position: absolute;
		top: 6px;
		left: 6px;
		padding: 2px 6px;
		border-radius: 4px;
		font-size: 10px;
		font-weight: 600;
		background-color: rgba(0, 0, 0, 0.8);
		color: white;
		backdrop-filter: blur(4px);
		border: 1px solid rgba(255, 255, 255, 0.2);
	}

	.next-episode-badge.today {
		background-color: rgba(255, 107, 107, 0.9);
		color: white;
		animation: pulse 2s infinite;
	}

	.next-episode-badge.tomorrow {
		background-color: rgba(255, 217, 61, 0.9);
		color: #333;
	}

	.next-episode-badge.this-week {
		background-color: rgba(0, 255, 127, 0.9);
		color: #333;
	}

	.next-episode-badge.aired {
		background-color: rgba(138, 43, 226, 0.9);
		color: white;
	}

	@keyframes pulse {
		0%, 100% { opacity: 1; }
		50% { opacity: 0.7; }
	}

	/* Episode info section */
	.episode-info {
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.episode-number {
		font-size: 11px;
		font-weight: 600;
		color: var(--highlight);
		letter-spacing: 0.5px;
	}

	.episode-name {
		font-size: 10px;
		color: var(--text-color);
		font-weight: 500;
		line-height: 1.2;
		max-height: 2.4em;
		overflow: hidden;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
	}

	.episode-time {
		font-size: 10px;
		color: var(--text-muted);
		margin-top: 1px;
	}

	.show-card.upcoming-mode .show-info {
		padding: 6px 10px 8px;
	}

	.show-meta {
		font-size: 11px;
		color: var(--text-muted);
		display: flex;
		justify-content: space-between;
		align-items: center;
	}
</style>
