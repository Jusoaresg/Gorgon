<script>
import { createEventDispatcher } from 'svelte';
import { fade } from 'svelte/transition';

import '$lib/styles/components/configurations/configurations.css';

export let torrentClient = 'qbittorrent';
export let qBittorrentHost = '';
export let qBittorrentPort = '';
export let qBittorrentUsername = '';
export let qBittorrentPassword = '';

export let qBittorrentDownloadFolder = '';

const dispatch = createEventDispatcher()

async function handleConfigChange() {
	dispatch('change', {
		client: torrentClient,
		config: {
			host: qBittorrentHost,
			port: qBittorrentPort,
			username: qBittorrentUsername,
			password: qBittorrentPassword,
			downloadFolder: qBittorrentDownloadFolder
		}
	})
}

async function testQBittorrentConnection() {
	// dummy function
}

</script>

<div id="section-qbittorrent" class="config-section" in:fade={{ duration: 200 }}>
	<div class="section-header">
		<h2>
			<i class="fas fa-download"></i>
			qBittorrent Configuration
		</h2>
		<div class="section-actions">
			<button class="action-btn test-btn" on:click={testQBittorrentConnection}>
				<i class="fas fa-plug"></i> Test Connection
			</button>
		</div>
	</div>

	<div class="settings-grid">
		<div class="setting-group">
			<label for="qBittorrentHost">
				<i class="fas fa-server"></i>
				qBittorrent Host
			</label>
			<input 
				type="url" 
				id="qBittorrentHost" 
				bind:value={qBittorrentHost} 
				on:input={handleConfigChange}
				placeholder="http://localhost"
			>
		</div>

		<div class="setting-group">
			<label for="qBittorrentPort">
				<i class="fas fa-ethernet"></i>
				qBittorrent Port
			</label>
			<input 
				type="text" 
				id="qBittorrentPort" 
				bind:value={qBittorrentPort} 
				on:input={handleConfigChange}
				placeholder="9191"
			>
		</div>

		<div class="setting-group">
			<label for="qBittorrentUsername">
				<i class="fas fa-user"></i>
				Username
			</label>
			<input 
				type="text" 
				id="qBittorrentUsername" 
				bind:value={qBittorrentUsername} 
				on:input={handleConfigChange}
				placeholder="admin"
			>
		</div>

		<div class="setting-group">
			<label for="qBittorrentPassword">
				<i class="fas fa-lock"></i>
				Password
			</label>
			<input 
				type="password" 
				id="qBittorrentPassword" 
				bind:value={qBittorrentPassword} 
				on:input={handleConfigChange}
				placeholder="Enter password"
			>
		</div>

		<div class="setting-group full-width">
			<label for="qBittorrentDownloadFolder">
				<i class="fas fa-folder-open"></i>
				Download Folder
			</label>
			<input 
				type="text" 
				id="qBittorrentDownloadFolder" 
				bind:value={qBittorrentDownloadFolder} 
				on:input={handleConfigChange}
				placeholder="/downloads"
			>
		</div>
	</div>
</div>
