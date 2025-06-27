<script>
import { goto } from '$app/navigation';
import { onMount } from 'svelte';
import { fade, fly, scale, slide } from 'svelte/transition';
import { quintOut } from 'svelte/easing';

// Styles
import '$lib/styles/shows.css';

let allShows = [];
let shows = [];
let currentView = 'grid';
let searchQuery = '';
let isLoading = true;

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
	} catch (error) {
		console.error('Erro ao carregar shows:', error);
		allShows = [];
		shows = [];
	} finally {
		isLoading = false;
	}
});
</script>

<svelte:head>
	<title>My Shows</title>
</svelte:head>

{#if isLoading}
	<div class="loading-container" in:fade={{ duration: 200 }}>
		<div class="spinner"></div>
		<p>Loading shows...</p>
	</div>
{:else}
	<div class="shows-page" in:fade={{ duration: 250 }}>
		<div class="shows-page container">
			<!-- Header -->
			<div class="header" in:fly={{ y: -20, duration: 300, delay: 50, easing: quintOut }}>
				<h1>My Shows</h1>
				<div class="header-controls">
					<input 
						type="text" 
						class="search-input" 
						placeholder="Filter shows..." 
						bind:value={searchQuery}
						on:input={handleSearch}
						in:scale={{ duration: 200, delay: 150 }}
					>
					<div class="view-controls" in:fly={{ x: 20, duration: 250, delay: 200 }}>
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
			<div class="stats-bar" in:fly={{ y: 20, duration: 300, delay: 100, easing: quintOut }}>
				<div class="stats-left">
					<div class="stat-item" in:scale={{ duration: 200, delay: 250 }}>
						<span>Total:</span>
						<span class="stat-value">{totalShows}</span>
					</div>
					<div class="stat-item" in:scale={{ duration: 200, delay: 280 }}>
						<span>Continuing:</span>
						<span class="stat-value">{continuingShows}</span>
					</div>
					<div class="stat-item" in:scale={{ duration: 200, delay: 310 }}>
						<span>Ended:</span>
						<span class="stat-value">{endedShows}</span>
					</div>
					<div class="stat-item" in:scale={{ duration: 200, delay: 340 }}>
						<span>Missing:</span>
						<span class="stat-value">{missingShows}</span>
					</div>
				</div>
				<div class="stats-right" in:fade={{ duration: 200, delay: 400 }}>
					<span>Last updated: now</span>
				</div>
			</div>

			<!-- Content -->
			<div class="shows-container">
				{#if shows.length === 0}
					<div class="empty-state" in:fade={{ duration: 300, delay: 200 }}>
						<h3>No shows found</h3>
						<p>Try adjusting your search or add some shows to your collection.</p>
					</div>
				{:else}
					<!-- Grid View -->
					{#if currentView === 'grid'}
						<div class="grid-view active" in:fade={{ duration: 200 }}>
							<div class="shows-grid">
								{#each shows as show, i (show.ID)}
									<div
										class="show-card"
										aria-hidden=true
										on:click={() => navigateToShow(show.ID)}
										in:fly={{ y: 30, duration: 250, delay: i * 20, easing: quintOut }}
										out:scale={{ duration: 150 }}
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
										</div>
										<div class="show-info" in:fly={{ y: 10, duration: 200, delay: 150 }}>
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
						<div class="list-view active" in:fade={{ duration: 200 }}>
							<div class="shows-list">
								<div class="list-header" in:fly={{ y: -10, duration: 200, delay: 50 }}>
									<div></div>
									<div>Title</div>
									<div>Year</div>
									<div>Status</div>
									<div>Rating</div>
									<div></div>
								</div>
								<div class="shows-list-content">
									{#each shows as show, i (show.ID)}
										<div
											class="show-row"
											aria-hidden=true
											on:click={() => navigateToShow(show.ID)}
											in:fly={{ x: -20, duration: 200, delay: i * 15, easing: quintOut }}
											out:fly={{ x: 20, duration: 150 }}
										>
											<div class="row-poster">
												{#if show.ImageMedium}
													<img src={show.ImageMedium} alt={show.Name} loading="lazy"
														 in:scale={{ duration: 150, delay: 50 }}>
												{:else}
													<div class="no-image-small" in:fade={{ duration: 150 }}></div>
												{/if}
											</div>
											<div class="row-title">{show.Name}</div>
											<div class="row-year">{show.FirstAired ? new Date(show.FirstAired).getFullYear() : 'N/A'}</div>
											<div class="row-status" 
												 style="background-color: {getStatusColor(show.Status)}"
												 in:scale={{ duration: 150, delay: 100 }}>{show.Status || 'Unknown'}</div>
											<div class="row-rating">{show.Rating || 'N/A'}</div>
											<div class="row-actions" in:fade={{ duration: 150, delay: 150 }}>
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
{/if}

<style>
.loading-container {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	height: 50vh;
	gap: 1rem;
}

.spinner {
	width: 40px;
	height: 40px;
	border: 4px solid #f3f3f3;
	border-top: 4px solid #3498db;
	border-radius: 50%;
	animation: spin 1s linear infinite;
}

@keyframes spin {
	0% { transform: rotate(0deg); }
	100% { transform: rotate(360deg); }
}
</style>
