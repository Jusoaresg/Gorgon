<script>
	import Modal from "$lib/components/Modal.svelte";
	import { createEventDispatcher } from "svelte";
	import { json } from "@sveltejs/kit";
	import { toast } from "svelte-sonner";
	
	let dispatch = createEventDispatcher();
	
	export let show = false;
	export let episode;
	export let searchResults = [];
	export let isSearching = false;
	export let searchError = null;
	
	const startSearch = async () => {
		isSearching = true
		try {
		const res = await fetch("/api/v1/prowlarr/search/episode", {
			method: 'POST',
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify({
				id: episode.ID,
			})})

			const json = await res.json();
			searchResults = json.data

			if (searchResults.length === 0) {
				toast.info("No available episodes")
			}

		} catch(error) {
			console.log("Error while manual searching for episodes")
			searchError = true;
		}
		isSearching = false
	}
	
	const downloadRelease = async (release) => {
		try {
		const res = await fetch("/api/v1/qbittorrent/add/episode", {
			method: 'POST',
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify({
				EpisodeID: episode.ID,
				InfoHash: release.infoHash,
				magneticUrl: release.guid,
			})})

			const json = await res.json();
			searchResults = json.data

		} catch(error) {
			console.log("Error while manual searching for episodes")
		}
		dispatch("close", release);
	}
	
	const openInfoUrl = (infoUrl) => {
		if (infoUrl) {
			window.open(infoUrl, '_blank', 'noopener,noreferrer');
		}
	}
	
	const closeModal = () => {
		show = false;
		dispatch("close");
	}
	
	const formatSize = (bytes) => {
		if (!bytes) return "Unknown";
		const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
		const i = Math.floor(Math.log(bytes) / Math.log(1024));
		return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i];
	}
	
	const formatAge = (publishDate) => {
		if (!publishDate) return "Unknown";
		const now = new Date();
		const releaseDate = new Date(publishDate);
		const diffTime = Math.abs(now - releaseDate);
		const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
		
		if (diffDays === 0) return "Today";
		if (diffDays === 1) return "1 day ago";
		if (diffDays < 30) return `${diffDays} days ago`;
		if (diffDays < 365) return `${Math.floor(diffDays / 30)} months ago`;
		return `${Math.floor(diffDays / 365)} years ago`;
	}
	

</script>

<Modal show={show} on:close={closeModal} size="full">
	<div slot="title">
		<i class="fas fa-search"></i>
		Manual Search Results
	</div>
	<div slot="body">
		<div class="search-container">
			<!-- Episode Info Header -->
			<div class="episode-header">
				<div class="episode-info">
					<h3 class="episode-title">{episode?.Name}</h3>
					<p class="episode-details">Season {episode?.Season} • Episode {episode?.Number}</p>
				</div>
				{#if !isSearching && searchResults.length === 0 && !searchError}
					<button class="episode-btn search-btn" on:click={startSearch}>
						<i class="fas fa-search"></i>
						Start Search
					</button>
				{/if}
			</div>
			
			<!-- Search Status -->
			{#if isSearching}
				<div class="search-status searching">
					<div class="loading-spinner"></div>
					<p>Searching indexers...</p>
				</div>
			{:else if searchError}
				<div class="search-status error">
					<i class="fas fa-exclamation-triangle"></i>
					<p>Search failed: {searchError}</p>
					<button class="episode-btn retry-btn" on:click={startSearch}>
						<i class="fas fa-redo"></i>
						Retry Search
					</button>
				</div>
			{:else if searchResults.length === 0}
				<div class="search-status empty">
					<i class="fas fa-inbox"></i>
					<p>Click "Start Search" to find releases</p>
				</div>
			{/if}
			
			<!-- Search Results -->
			{#if searchResults.length > 0}
				<div class="results-header">
					<h4>Found {searchResults.length} release{searchResults.length !== 1 ? 's' : ''}</h4>
				</div>
				
				<div class="results-list">
					{#each searchResults as release, index}
						<div class="release-item">
							<div class="release-main">
								<div class="release-info">
									<div class="release-title">{release.title}</div>
									<div class="release-meta">
										<span class="indexer">
											<i class="fas fa-server"></i>
											{release.indexer}
										</span>

										<span class="age">
											<i class="fas fa-clock"></i>
											{formatAge(release.publishDate)}
										</span>
										<span class="size">
											<i class="fas fa-hdd"></i>
											{formatSize(release.size)}
										</span>
										{#if release.seeders !== undefined && release.seeders !== null}
											<span class="seeders" class:low-seeders={release.seeders < 5}>
												<i class="fas fa-arrow-up"></i>
												{release.seeders}
											</span>
										{/if}
										{#if release.leechers !== undefined && release.leechers !== null}
											<span class="leechers">
												<i class="fas fa-arrow-down"></i>
												{release.leechers}
											</span>
										{/if}
										{#if release.grabs !== undefined && release.grabs !== null && release.grabs > 0}
											<span class="grabs">
												<i class="fas fa-download"></i>
												{release.grabs} grabs
											</span>
										{/if}
									</div>
								</div>
								
								<div class="release-actions">
									{#if release.infoUrl}
										<button 
											aria-label="infoUrl"
											class="episode-btn info-btn" 
											on:click={() => openInfoUrl(release.infoUrl)}
											title="View on {release.indexer}"
										>
											<i class="fas fa-external-link-alt"></i>
										</button>
									{/if}
									<button 
										class="episode-btn download-btn" 
										on:click={() => downloadRelease(release)}
										disabled={release.downloading}
									>
										{#if release.downloading}
											<div class="mini-spinner"></div>
										{:else}
											<i class="fas fa-download"></i>
										{/if}
									</button>
								</div>
							</div>
							
							{#if release.rejected}
								<div class="rejection-reasons">
									<i class="fas fa-times-circle"></i>
									<span>Rejected: {release.rejectionReason}</span>
								</div>
							{/if}
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
	<div slot="footer">
		<div class="footer-buttons">
			<button class="episode-btn close-btn" on:click={closeModal}>
				<i class="fas fa-times"></i>
				Close
			</button>
			{#if searchResults.length > 0}
				<span class="results-count">{searchResults.length} results found</span>
			{/if}
		</div>
	</div>
</Modal>

<style>
	.search-container {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
		max-height: 70vh;
		overflow: hidden;
	}
	
	.episode-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		padding: 1rem;
		background-color: rgba(255, 255, 255, 0.05);
		border-radius: 8px;
		border-left: 3px solid var(--highlight, #3498db);
	}
	
	.episode-info {
		flex: 1;
	}
	
	.episode-title {
		margin: 0 0 0.5rem 0;
		color: var(--text-color);
		font-size: 1.2rem;
		font-weight: 600;
	}
	
	.episode-details {
		margin: 0;
		color: var(--text-color);
		opacity: 0.7;
		font-size: 0.9rem;
	}
	
	.search-status {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 2rem;
		text-align: center;
		border-radius: 8px;
	}
	
	.search-status.searching {
		background-color: rgba(52, 152, 219, 0.1);
		border: 1px solid rgba(52, 152, 219, 0.3);
	}
	
	.search-status.error {
		background-color: rgba(231, 76, 60, 0.1);
		border: 1px solid rgba(231, 76, 60, 0.3);
	}
	
	.search-status.empty {
		background-color: rgba(255, 255, 255, 0.02);
		border: 1px solid rgba(255, 255, 255, 0.1);
	}
	
	.search-status i {
		font-size: 2rem;
		margin-bottom: 1rem;
		opacity: 0.7;
	}
	
	.search-status p {
		margin: 0 0 1rem 0;
		color: var(--text-color);
		font-size: 1rem;
	}
	
	.loading-spinner {
		width: 32px;
		height: 32px;
		border: 3px solid rgba(52, 152, 219, 0.3);
		border-top: 3px solid #3498db;
		border-radius: 50%;
		animation: spin 1s linear infinite;
		margin-bottom: 1rem;
	}
	
	.mini-spinner {
		width: 16px;
		height: 16px;
		border: 2px solid rgba(255, 255, 255, 0.3);
		border-top: 2px solid white;
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}
	
	@keyframes spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
	}
	
	.results-header {
		padding: 0 0.5rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
		padding-bottom: 0.75rem;
	}
	
	.results-header h4 {
		margin: 0;
		color: var(--text-color);
		font-size: 1rem;
		font-weight: 600;
	}
	
	.results-list {
		flex: 1;
		overflow-y: auto;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}
	
	.release-item {
		background-color: rgba(255, 255, 255, 0.03);
		border: 1px solid rgba(255, 255, 255, 0.08);
		border-radius: 8px;
		transition: all 0.2s ease;
	}
	
	.release-item:hover {
		background-color: rgba(255, 255, 255, 0.06);
		border-color: rgba(255, 255, 255, 0.15);
	}
	
	.release-main {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem;
		gap: 1rem;
	}
	
	.release-info {
		flex: 1;
		min-width: 0;
	}
	
	.release-title {
		font-weight: 500;
		color: var(--text-color);
		margin-bottom: 0.5rem;
		word-break: break-word;
		line-height: 1.3;
	}
	
	.release-meta {
		display: flex;
		flex-wrap: wrap;
		gap: 1rem;
		font-size: 0.85rem;
		opacity: 0.8;
	}
	
	.release-meta span {
		display: flex;
		align-items: center;
		gap: 0.3rem;
		color: var(--text-color);
	}
	
	.release-meta i {
		width: 12px;
		text-align: center;
	}
	

	.low-seeders {
		color: #e74c3c !important;
	}
	
	.leechers {
		color: #e74c3c !important;
	}
	
	.grabs {
		color: #27ae60 !important;
	}
	
	.release-actions {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		flex-shrink: 0;
	}
	
	.rejection-reasons {
		padding: 0.75rem 1rem;
		background-color: rgba(231, 76, 60, 0.1);
		border-top: 1px solid rgba(231, 76, 60, 0.2);
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.85rem;
		color: #e74c3c;
	}
	
	/* Episode Button Styles - Consistent with Episode Card */
	.episode-btn {
		width: 32px;
		height: 32px;
		border-radius: 4px;
		background-color: rgba(255, 255, 255, 0.1);
		color: var(--text-color);
		border: none;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		transition: all 0.2s ease;
		font-size: 0.9rem;
	}
	
	.episode-btn:hover {
		background-color: rgba(255, 255, 255, 0.2);
		transform: translateY(-1px);
	}
	
	.episode-btn:active {
		transform: scale(0.95);
	}
	
	.episode-btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
		transform: none !important;
	}
	
	/* Search Button - Primary action */
	.search-btn {
		background-color: rgba(52, 152, 219, 0.2);
		color: #3498db;
		border: 1px solid rgba(52, 152, 219, 0.3);
		width: auto;
		padding: 0.5rem 1rem;
		gap: 0.5rem;
		font-weight: 500;
	}
	
	.search-btn:hover {
		background-color: rgba(52, 152, 219, 0.3);
		color: #fff;
		border-color: rgba(52, 152, 219, 0.5);
		box-shadow: 0 4px 8px rgba(52, 152, 219, 0.3);
	}
	
	/* Retry Button - Similar to search but for error state */
	.retry-btn {
		background-color: rgba(52, 152, 219, 0.2);
		color: #3498db;
		border: 1px solid rgba(52, 152, 219, 0.3);
		width: auto;
		padding: 0.5rem 1rem;
		gap: 0.5rem;
		font-weight: 500;
	}
	
	.retry-btn:hover {
		background-color: rgba(52, 152, 219, 0.3);
		color: #fff;
		border-color: rgba(52, 152, 219, 0.5);
		box-shadow: 0 4px 8px rgba(52, 152, 219, 0.3);
	}
	
	/* Info Button - Info action */
	.info-btn {
		background-color: rgba(52, 152, 219, 0.2);
		color: #3498db;
		border: 1px solid rgba(52, 152, 219, 0.3);
	}
	
	.info-btn:hover {
		background-color: rgba(52, 152, 219, 0.3);
		color: #fff;
		border-color: rgba(52, 152, 219, 0.5);
		box-shadow: 0 4px 8px rgba(52, 152, 219, 0.3);
	}
	
	/* Download Button - Success action */
	.download-btn {
		background-color: rgba(39, 174, 96, 0.2);
		color: #27ae60;
		border: 1px solid rgba(39, 174, 96, 0.3);
	}
	
	.download-btn:hover:not(:disabled) {
		background-color: rgba(39, 174, 96, 0.3);
		color: #fff;
		border-color: rgba(39, 174, 96, 0.5);
		box-shadow: 0 4px 8px rgba(39, 174, 96, 0.3);
	}
	
	/* Close Button - Secondary action */
	.close-btn {
		width: auto;
		padding: 0.5rem 1rem;
		gap: 0.5rem;
		font-weight: 500;
	}
	
	.footer-buttons {
		display: flex;
		justify-content: space-between;
		align-items: center;
		width: 100%;
	}
	
	.results-count {
		color: var(--text-color);
		opacity: 0.7;
		font-size: 0.9rem;
	}
	
	/* Mobile responsiveness */
	@media (max-width: 768px) {
		.episode-header {
			flex-direction: column;
			gap: 1rem;
			align-items: stretch;
		}
		
		.search-btn {
			width: 100%;
			justify-content: center;
		}
		
		.release-main {
			flex-direction: column;
			align-items: stretch;
			gap: 0.75rem;
		}
		
		.release-meta {
			gap: 0.75rem;
		}
		
		.release-actions {
			justify-content: flex-end;
		}
		
		.footer-buttons {
			flex-direction: column;
			gap: 0.5rem;
		}
		
		.close-btn {
			width: 100%;
		}
	}
</style>
