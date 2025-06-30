<svelte:head>
	<link href="/css/pages/add-show-edit.css" rel="stylesheet" />
</svelte:head>

<script>
    import { goto } from "$app/navigation";
    import { error } from "@sveltejs/kit";
    import { onMount } from "svelte";
    import { fade, slide, fly, scale, blur } from "svelte/transition";
    import { toast } from "svelte-sonner";

    // Styles
    import '$lib/styles/pages/add-show-edit.css'

    let show = null;
    let isAdded = false;

    let loading = false;
    let pageReady = false;

    let tracking = "all";

    onMount(async () => {
	    const stored = sessionStorage.getItem('showToEdit');

	    const json = JSON.parse(stored);
	    show = json.Show;
	    isAdded = json.isAdded;
	    
	    // Add a small delay for smoother page transitions
	    setTimeout(() => {
	        pageReady = true;
	    }, 100);
    });

    async function addShow() {
	const progressIndicator = document.getElementById('progress');
	    loading = true;

	    try {
		    const res = await fetch("/api/v1/database/show", {
			    method: 'POST',
			    headers: {
				    "Content-Type": "application/json"
			    },
			    body: JSON.stringify({
				    id: show.id,
				    tracking_type: tracking
			    })
		    })

		    if(!res.ok) {
			    return error
		    }

		    goto("/")

	    } catch(e) {
		    return error;
	    } finally {
		    toast.success("Show added successfuly.");
		    loading = false;
	    }
    }

    function cancel() {
        if (confirm('Are you sure you want to cancel? Your configuration will be lost.')) {
		    goto("/");
        }
    }
</script>

{#if show && pageReady}
	<div class="add-show-edit-page" in:fade={{ duration: 300, delay: 100 }}>
		<div class="edit-container">
			<div class="edit-header" in:fly={{ y: -30, duration: 500, delay: 200 }}>
				<h1>Configure Show Settings</h1>
				<p class="edit-subtitle">Set up tracking preferences before adding this show to your collection</p>
			</div>

			<!-- Show Preview -->
			<div class="show-preview" in:fly={{ y: 30, duration: 500, delay: 300 }}>
				<div class="show-poster">
					<div class="placeholder-poster">
						<img src={show.image.medium} alt="show banner" in:scale={{ duration: 400, delay: 400 }}/>
					</div>
				</div>
				<div class="show-info" in:fade={{ duration: 400, delay: 500 }}>
					<h2 class="show-title">{show.name}</h2>
					<div class="show-meta" in:slide={{ duration: 300, delay: 600 }}>
						<span class="show-status">{show.status}</span>
						<!-- TODO: Change this -->
						<span class="show-date">{show.premiered}</span>
						<span class="show-language">{show.language}</span>
						<!-- <span class="show-network">AMC</span> -->
					</div>
					<div class="show-summary" in:fade={{ duration: 400, delay: 700 }}>
						{@html show.summary}
					</div>
				</div>
			</div>

			<!-- Tracking Configuration -->
			<div class="configuration-section" in:fly={{ x: -30, duration: 500, delay: 800 }}>
				<h3 class="section-title"><i class="fas fa-eye"></i> Episode Tracking</h3>
				<div class="form-group">
					<label class="form-label" for="">How should episodes be tracked?</label>
					<p class="form-description">Choose how you want to handle episode monitoring for this show.</p>
					<div class="radio-group">
						<label class="radio-option" in:slide={{ duration: 300, delay: 900 }}>
							<input type="radio" bind:group={tracking} name="tracking" value="all" checked>
							<div class="radio-content">
								<div class="radio-title">Track All Episodes</div>
								<div class="radio-description">Monitor and download all episodes, including already aired ones</div>
							</div>
						</label>
						<label class="radio-option" in:slide={{ duration: 300, delay: 1000 }}>
							<input type="radio" bind:group={tracking} name="tracking" value="future">
							<div class="radio-content">
								<div class="radio-title">Track Future Episodes Only</div>
								<div class="radio-description">Only monitor episodes that haven't aired yet</div>
							</div>
						</label>
						<label class="radio-option" in:slide={{ duration: 300, delay: 1100 }}>
							<input type="radio" bind:group={tracking} name="tracking" value="none">
							<div class="radio-content">
								<div class="radio-title">No Automatic Tracking</div>
								<div class="radio-description">Add show to library but don't automatically search for episodes</div>
							</div>
						</label>
					</div>
				</div>
			</div>

			<!-- Quality Settings -->
			<!--
			<div class="configuration-section">
				<h3 class="section-title"><i class="fas fa-hd-video"></i> Quality Preferences</h3>
				<div class="form-group">
					<label class="form-label">Preferred video quality</label>
					<p class="form-description">Select which video qualities are acceptable for downloads.</p>
					<div class="quality-grid">
						<label class="quality-option">
							<input type="checkbox" name="quality" value="480p">
							<span>480p</span>
						</label>
						<label class="quality-option">
							<input type="checkbox" name="quality" value="720p" checked>
							<span>720p HD</span>
						</label>
						<label class="quality-option">
							<input type="checkbox" name="quality" value="1080p" checked>
							<span>1080p FHD</span>
						</label>
						<label class="quality-option">
							<input type="checkbox" name="quality" value="4k">
							<span>4K UHD</span>
						</label>
					</div>
				</div>
			</div>
			-->

			<!-- Advanced Settings -->
			<!--
			<div class="configuration-section">
				<h3 class="section-title"><i class="fas fa-cog"></i> Advanced Options</h3>
				<div class="form-group">
					<label class="form-label">Additional settings</label>
					<div class="checkbox-group">
						<label class="checkbox-option">
							<input type="checkbox" name="options" value="search-backlog" checked>
							<span>Search for existing episodes in backlog</span>
						</label>
						<label class="checkbox-option">
							<input type="checkbox" name="options" value="season-folders" checked>
							<span>Use season folders for organization</span>
						</label>
						<label class="checkbox-option">
							<input type="checkbox" name="options" value="anime-detection">
							<span>Enable anime detection and naming</span>
						</label>
						<label class="checkbox-option">
							<input type="checkbox" name="options" value="notifications" checked>
							<span>Send notifications when episodes are downloaded</span>
						</label>
					</div>
				</div>

				<div class="form-group">
					<label class="form-label" for="root-folder">Root Folder</label>
					<p class="form-description">Choose where this show's files will be stored.</p>
					<select class="form-control" id="root-folder">
						<option value="/media/tv">TV Shows (/media/tv)</option>
						<option value="/media/series">Series (/media/series)</option>
						<option value="/downloads/complete">Downloads (/downloads/complete)</option>
					</select>
				</div>

				<div class="form-group">
					<label class="form-label" for="language">Language Profile</label>
					<p class="form-description">Preferred language for releases.</p>
					<select class="form-control" id="language">
						<option value="english" selected>English</option>
						<option value="portuguese">Portuguese</option>
						<option value="spanish">Spanish</option>
						<option value="french">French</option>
						<option value="any">Any Language</option>
					</select>
				</div>
			</div>
			-->

			{#if loading}
				<div class="progress-indicator" id="progress" 
				     in:scale={{ duration: 300 }} 
				     out:fade={{ duration: 200 }}>
					<div class="spinner"></div>
					<span>Adding show to your library...</span>
				</div>
			{:else}
				<!-- Action Buttons -->
				<div class="action-buttons" 
				     in:fly={{ y: 30, duration: 400, delay: 1200 }} 
				     out:fade={{ duration: 200 }}>
					<button class="btn btn-primary" on:click={addShow} disabled={loading}>
						{#if loading}
							<i class="fas fa-spinner fa-spin"></i> 
							Adding Show...
						{:else}
							<i class="fas fa-plus"></i>
							Add Show to Library
						{/if}
					</button>
					<button class="btn btn-danger" on:click={cancel}>
						<i class="fas fa-times"></i>
						Cancel
					</button>
				</div>
			{/if}

		</div>
	</div>
{:else if show}
	<!-- Loading state while page is preparing -->
	<div class="add-show-edit-page" in:blur={{ duration: 300 }}>
		<div class="edit-container">
			<div class="loading-placeholder" in:fade={{ duration: 200 }}>
				<div class="spinner"></div>
				<span>Loading show details...</span>
			</div>
		</div>
	</div>
{/if}
