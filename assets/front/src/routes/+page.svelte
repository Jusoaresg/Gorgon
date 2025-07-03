<script>
import ShowCard from '$lib/components/ShowCard.svelte';
import { goto } from '$app/navigation';
import { onMount } from 'svelte';
import { fade, fly, scale, slide } from 'svelte/transition';
import { quintOut, quintIn } from 'svelte/easing';

// Styles
import '$lib/styles/shows.css';
    import Loading from '$lib/components/Loading.svelte';

let allShows = [];
let shows = [];
let showFilter = 'all';
let searchQuery = '';
let isLoading = true;

// Stats calculation
$: totalShows = shows.length;
$: continuingShows = shows.filter(show => show.Status === 'Running').length;
$: endedShows = shows.filter(show => show.Status === 'Ended').length;
$: missingShows = shows.filter(show => !show.Status || show.Status === 'Missing').length;

let upcomingEpisodes = [];

$: if (!isLoading && allShows.length > 0) {
	const now = new Date();

	upcomingEpisodes = allShows
		.map(showAgg => {
			const { Show, Episodes } = showAgg;

			// Get episodes from today onwards, including already aired today
			const relevantEps = Episodes
			.filter(ep => {
				const airDate = new Date(ep.AirStamp * 1000);
				const today = new Date();
				today.setHours(0, 0, 0, 0); // Start of today
				return airDate >= today; // Include episodes from today onwards
			})
			.sort((a, b) => new Date(a.AirStamp * 1000) - new Date(b.AirStamp * 1000));

			if (relevantEps.length === 0) return null;

			const nextEpisode = relevantEps[0];
			const airDate = new Date(nextEpisode.AirStamp * 1000);

			// Calculate days difference and check if already aired
			const diffTime = airDate - now;
			const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
			const hasAlreadyAired = airDate < now; // Episode time has passed
			const isToday = diffDays === 0 || hasAlreadyAired && diffDays <= 0;

			// Determine priority for sorting
			let priority = 999;
			if (isToday && hasAlreadyAired) {
				priority = 0; // Already aired today - show first with different styling
			} else if (diffDays === 0) {
				priority = 1; // Airing today - show second
			} else if (diffDays === 1) {
				priority = 2; // Tomorrow
			} else if (diffDays > 1 && diffDays <= 7) {
				priority = 3; // This week
			} else if (diffDays <= 30) {
				priority = 4; // This month
			}

			return {
				show: Show,
				nextEpisode: nextEpisode,
				priority: priority,
				daysUntil: diffDays,
				hasAlreadyAired: hasAlreadyAired,
				isToday: isToday
			};
		})
		.filter(Boolean)
		// Sort by priority first, then by air date
		.sort((a, b) => {
			if (a.priority !== b.priority) {
				return a.priority - b.priority;
			}
			return new Date(a.nextEpisode.AirStamp * 1000) - new Date(b.nextEpisode.AirStamp * 1000);
		});
}

$: todayEpisodes = upcomingEpisodes.filter(item => item.category === 'today');
$: tomorrowEpisodes = upcomingEpisodes.filter(item => item.category === 'tomorrow');
$: thisWeekEpisodes = upcomingEpisodes.filter(item => item.category === 'this-week');
$: thisMonthEpisodes = upcomingEpisodes.filter(item => item.category === 'this-month');

function handleSearch() {
	const query = searchQuery.toLowerCase();

	if (showFilter === 'upcoming') {
		const now = new Date();

		upcomingEpisodes = allShows
			.map(showAgg => {
				const { Show, Episodes } = showAgg;

				// Filter by search query first
				if (query && !Show.Name.toLowerCase().includes(query)) {
					return null;
				}

				// Get episodes from today onwards, including already aired today
				const relevantEps = Episodes
				.filter(ep => {
					const airDate = new Date(ep.AirStamp * 1000);
					const today = new Date();
					today.setHours(0, 0, 0, 0); // Start of today
					return airDate >= today; // Include episodes from today onwards
				})
				.sort((a, b) => new Date(a.AirStamp * 1000) - new Date(b.AirStamp * 1000));

				if (relevantEps.length === 0) return null;

				const nextEpisode = relevantEps[0];
				const airDate = new Date(nextEpisode.AirStamp * 1000);

				// Calculate days difference and check if already aired
				const diffTime = airDate - now;
				const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
				const hasAlreadyAired = airDate < now; // Episode time has passed
				const isToday = diffDays === 0 || hasAlreadyAired && diffDays <= 0;

				// Determine priority for sorting
				let priority = 999;
				if (isToday && hasAlreadyAired) {
					priority = 0; // Already aired today - show first with different styling
				} else if (diffDays === 0) {
					priority = 1; // Airing today - show second
				} else if (diffDays === 1) {
					priority = 2; // Tomorrow
				} else if (diffDays > 1 && diffDays <= 7) {
					priority = 3; // This week
				} else if (diffDays <= 30) {
					priority = 4; // This month
				}

				return {
					show: Show,
					nextEpisode: nextEpisode,
					priority: priority,
					daysUntil: diffDays,
					hasAlreadyAired: hasAlreadyAired,
					isToday: isToday
				};
			})
			.filter(Boolean)
			.sort((a, b) => {
				if (a.priority !== b.priority) {
					return a.priority - b.priority;
				}
				return new Date(a.nextEpisode.AirStamp * 1000) - new Date(b.nextEpisode.AirStamp * 1000);
			});
	} else if (showFilter === 'all') {
		shows = allShows.filter(showAgg => {
			const showName = showAgg.Show.Name.toLowerCase();
			return showName.includes(query);
		});
	};
}

function navigateToShow(event) {
	const showID = event.detail
	goto(`/show/${showID}`);
}

function switchFilter(filter) {
	showFilter = filter;
	localStorage.setItem('sortMode', showFilter)
}

onMount(async () => {
	const sortMode = localStorage.getItem('sortMode');
	if (sortMode) {
		showFilter = sortMode;
	}

	try {
		const res = await fetch("/api/v1/database/show/full");
		const json = await res.json();
		console.log(json.data)
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
	<Loading text="Loading shows..."/>
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
							class:active={showFilter === 'all'}
							on:click={() => switchFilter('all')}
						>
							All Shows
						</button>
						<button 
							class="view-btn" 
							class:active={showFilter === 'upcoming'}
							on:click={() => switchFilter('upcoming')}
						>
							Next Episodes
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
					<!-- All Shows -->
					{#if showFilter === 'all'}
						<div class="grid-view active" in:fade={{ duration: 200 }}>
							<div class="shows-grid">
								{#each shows as showAggregated, i (showAggregated.Show.ID)}
									<ShowCard
										show={showAggregated.Show}
										showUpcoming={false}
										{i}
										on:select={navigateToShow}
									/>
								{/each}
							</div>
						</div>
					{/if}


					<!-- Next Episode -->
					{#if showFilter === 'upcoming'}
						<div class="grid-view active" in:fade={{ duration: 200 }}>
							<div class="shows-grid">
								{#each upcomingEpisodes as { show, nextEpisode, category, daysUtil }, i (show.ID)}
									<ShowCard
										show={show}
										showUpcoming={true}
										nextEpisode={nextEpisode}
										{i}
										on:select={navigateToShow}
									/>
								{/each}
							</div>
						</div>
					{/if}
				{/if}
			</div>
		</div>
	</div>
{/if}
