<script>
import '$lib/styles/components/configurations/configurations.css';
import { fade } from 'svelte/transition';
import { createEventDispatcher } from 'svelte';

export let defaultShowInfoFolder = '';
export let showsFolder = '';

const dispatch = createEventDispatcher();

async function handleConfigChange() {
	dispatch("change", {
		showInfoFolder: defaultShowInfoFolder,
		showsFolder: showsFolder
	})
}
</script>


<div id="section-folders" class="config-section" in:fade={{ duration: 200 }}>
	<div class="section-header">
		<h2>
			<i class="fas fa-folder"></i>
			Folder Configuration
		</h2>
	</div>

	<div class="settings-grid">
		<div class="setting-group full-width">
			<label for="defaultShowInfoFolder">
				<i class="fas fa-info-circle"></i>
				Default Show Info Folder
			</label>
			<input 
				type="text" 
				id="defaultShowInfoFolder" 
				bind:value={defaultShowInfoFolder} 
				on:input={handleConfigChange}
				placeholder="assets/shows"
			>
			<small>Folder where show metadata and images are stored</small>
		</div>

		<div class="setting-group full-width">
			<label for="showsFolder">
				<i class="fas fa-video"></i>
				Shows Folder
			</label>
			<input 
				type="text" 
				id="showsFolder" 
				bind:value={showsFolder} 
				on:input={handleConfigChange}
				placeholder="/home/user/Videos/animes"
			>
			<small>Main folder where your shows/animes are stored</small>
		</div>
	</div>
</div>

