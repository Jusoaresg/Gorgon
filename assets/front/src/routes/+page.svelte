<script>
import { goto } from '$app/navigation';
import { onMount } from 'svelte';
import { fade, fly, scale } from 'svelte/transition';

// Styles
import '$lib/styles/shows.css';

let allShows = [];
let shows = [];
let currentView = 'grid';
let searchQuery = '';

// Stats calculadas
$: totalShows = shows.length;
$: continuingShows = shows.filter(show => show.Status === 'Running').length;
$: endedShows = shows.filter(show => show.Status === 'Ended').length;
$: missingShows = shows.filter(show => !show.Status || show.Status === 'Missing').length;

function navigateToShow(showId) {
	goto(`/show/${showId}`);
}

function handleSearch() {
	const query = searchQuery.toLowerCase();
	shows = allShows.filter(show => 
		show.Name.toLowerCase().includes(query)
	);
}

function switchView(view) {
	currentView = view;
}

function getStatusColor(status) {
	switch(status) {
		case 'Running': return '#00FF7F';
		case 'Ended': return '#FF6B6B';
		case 'Missing': return '#FFD93D';
		default: return '#999';
	}
}

onMount(async () => {
	try {
		const res = await fetch("/api/v1/database/show");
		const json = await res.json();
		allShows = json.data || [];
		shows = json.data || [];
		console.log(shows)
	} catch (error) {
		console.error('Erro ao carregar shows:', error);
		allShows = [];
		shows = [];
	}
});
</script>

<svelte:head>
	<title>My Shows</title>
</svelte:head>

<div class="shows-page">
	<div class="shows-page container">
		<!-- Header -->
		<div class="header">
			<h1>My Shows</h1>
			<div class="header-controls">
				<input 
					type="text" 
					class="search-input" 
					placeholder="Filter shows..." 
					bind:value={searchQuery}
					on:input={handleSearch}
				>
				<div class="view-controls">
					<button 
						class="view-btn" 
						class:active={currentView === 'grid'}
						on:click={() => switchView('grid')}
					>
						Grid
					</button>
					<button 
						class="view-btn" 
						class:active={currentView === 'list'}
						on:click={() => switchView('list')}
					>
						List
					</button>
				</div>
			</div>
		</div>

		<!-- Stats -->
		<div class="stats-bar">
			<div class="stats-left">
				<div class="stat-item">
					<span>Total:</span>
					<span class="stat-value">{totalShows}</span>
				</div>
				<div class="stat-item">
					<span>Continuing:</span>
					<span class="stat-value">{continuingShows}</span>
				</div>
				<div class="stat-item">
					<span>Ended:</span>
					<span class="stat-value">{endedShows}</span>
				</div>
				<div class="stat-item">
					<span>Missing:</span>
					<span class="stat-value">{missingShows}</span>
				</div>
			</div>
			<div class="stats-right">
				<span>Last updated: now</span>
			</div>
		</div>

		<!-- Content -->
		<div class="shows-container">
			{#if shows.length === 0}
				<div class="empty-state">
					<h3>No shows found</h3>
					<p>Try adjusting your search or add some shows to your collection.</p>
				</div>
			{:else}
				<!-- Grid View -->
				{#if currentView === 'grid'}
					<div class="grid-view active">
						<div class="shows-grid">
							{#each shows as show}
								<div
									class="show-card"
									aria-hidden=true
									on:click={() => navigateToShow(show.ID)}
									transition:fly={{ y:20, duration: 200 }}
								>
									<div class="show-poster">
										{#if show.ImageOriginal}
											<img src={show.ImageMedium} alt={show.Name} loading="lazy">
										{:else}
											<div class="no-image">
												<span>No Image</span>
											</div>
										{/if}
										<div 
											class="show-status" 
											style="background-color: {getStatusColor(show.Status)}"
										></div>
									</div>
									<div class="show-info">
										<div class="show-title">{show.Name}</div>
										<div class="show-meta">
											<span>{show.Premiered ? show.Premiered : 'N/A'}</span>
											<span>{show.Rating || 'N/A'}</span>
										</div>
									</div>
								</div>
							{/each}
						</div>
					</div>
				{/if}

				<!-- List View -->
				{#if currentView === 'list'}
					<div class="list-view active">
						<div class="shows-list">
							<div class="list-header">
								<div></div>
								<div>Title</div>
								<div>Year</div>
								<div>Status</div>
								<div>Rating</div>
								<div></div>
							</div>
							<div class="shows-list-content">
								{#each shows as show}
									<div
										class="show-row"
										aria-hidden=true
										on:click={() => navigateToShow(show.ID)}
										transition:fly={{ y:20, duration: 200 }}
									>
										<div class="row-poster">
											{#if show.ImageMedium}
												<img src={show.ImageMedium} alt={show.Name} loading="lazy">
											{:else}
												<div class="no-image-small"></div>
											{/if}
										</div>
										<div class="row-title">{show.Name}</div>
										<div class="row-year">{show.FirstAired ? new Date(show.FirstAired).getFullYear() : 'N/A'}</div>
										<div class="row-status" style="background-color: {getStatusColor(show.Status)}">{show.Status || 'Unknown'}</div>
										<div class="row-rating">{show.Rating || 'N/A'}</div>
										<div class="row-actions">
											<div class="action-icon" title="Edit">✎</div>
											<div class="action-icon" title="Delete">×</div>
										</div>
									</div>
								{/each}
							</div>
						</div>
					</div>
				{/if}
			{/if}
		</div>
	</div>
</div>
