import './style.css';
import './app.css';

import logo from './assets/images/logo-universal.png';
import * as App from '../wailsjs/go/app/App';

// Create the main UI
document.querySelector('#app').innerHTML = `
    <div class="app">
        <header class="header">
            <img id="logo" class="logo" alt="LRCGET Logo">
            <h1>LRCGET</h1>
            <p>Mass-download LRC synced lyrics for your music library</p>
        </header>
        
        <main class="main">
            <div class="setup-section" id="setup-section">
                <h2>Setup</h2>
                <div class="form-group">
                    <label for="directories">Music Directories:</label>
                    <input type="text" id="directories" placeholder="Enter directory paths separated by commas">
                    <button id="set-directories" class="btn">Set Directories</button>
                </div>
                <button id="initialize-library" class="btn btn-primary">Initialize Library</button>
            </div>
            
            <div class="library-section" id="library-section" style="display: none;">
                <h2>Library</h2>
                <div class="stats">
                    <div class="stat">
                        <span class="stat-label">Tracks:</span>
                        <span class="stat-value" id="track-count">0</span>
                    </div>
                    <div class="stat">
                        <span class="stat-label">Albums:</span>
                        <span class="stat-value" id="album-count">0</span>
                    </div>
                    <div class="stat">
                        <span class="stat-label">Artists:</span>
                        <span class="stat-value" id="artist-count">0</span>
                    </div>
                </div>
                
                <div class="tabs">
                    <button class="tab-btn active" data-tab="tracks">Tracks</button>
                    <button class="tab-btn" data-tab="albums">Albums</button>
                    <button class="tab-btn" data-tab="artists">Artists</button>
                </div>
                
                <div class="tab-content">
                    <div id="tracks-tab" class="tab-pane active">
                        <div class="search-box">
                            <input type="text" id="track-search" placeholder="Search tracks...">
                        </div>
                        <div id="tracks-list" class="list"></div>
                    </div>
                    <div id="albums-tab" class="tab-pane">
                        <div class="search-box">
                            <input type="text" id="album-search" placeholder="Search albums...">
                        </div>
                        <div id="albums-list" class="list"></div>
                    </div>
                    <div id="artists-tab" class="tab-pane">
                        <div class="search-box">
                            <input type="text" id="artist-search" placeholder="Search artists...">
                        </div>
                        <div id="artists-list" class="list"></div>
                    </div>
                </div>
            </div>
            
            <div class="player-section" id="player-section" style="display: none;">
                <h2>Now Playing</h2>
                <div class="player-info">
                    <div class="track-info">
                        <div class="track-title" id="current-track-title">No track selected</div>
                        <div class="track-artist" id="current-track-artist"></div>
                        <div class="track-album" id="current-track-album"></div>
                    </div>
                    <div class="player-controls">
                        <button id="play-btn" class="btn">Play</button>
                        <button id="pause-btn" class="btn">Pause</button>
                        <button id="stop-btn" class="btn">Stop</button>
                        <div class="volume-control">
                            <label for="volume">Volume:</label>
                            <input type="range" id="volume" min="0" max="100" value="100">
                        </div>
                    </div>
                    <div class="progress">
                        <div class="progress-bar">
                            <div class="progress-fill" id="progress-fill"></div>
                        </div>
                        <div class="time-info">
                            <span id="current-time">0:00</span>
                            <span id="total-time">0:00</span>
                        </div>
                    </div>
                </div>
            </div>
        </main>
        
        <div class="status-bar" id="status-bar">
            <span id="status-text">Ready</span>
        </div>
    </div>
`;

// Set logo
document.getElementById('logo').src = logo;

// Initialize the application
async function initializeApp() {
    try {
        // Check if library is initialized
        const isInitialized = await app.GetInit();
        
        if (isInitialized) {
            showLibrarySection();
            await loadLibraryData();
        } else {
            showSetupSection();
        }
    } catch (error) {
        console.error('Failed to initialize app:', error);
        updateStatus('Failed to initialize application');
    }
}

// Show setup section
function showSetupSection() {
    document.getElementById('setup-section').style.display = 'block';
    document.getElementById('library-section').style.display = 'none';
    document.getElementById('player-section').style.display = 'none';
}

// Show library section
function showLibrarySection() {
    document.getElementById('setup-section').style.display = 'none';
    document.getElementById('library-section').style.display = 'block';
    document.getElementById('player-section').style.display = 'block';
}

// Load library data
async function loadLibraryData() {
    try {
        updateStatus('Loading library data...');
        
        // Load tracks, albums, and artists
        const [tracks, albums, artists] = await Promise.all([
            app.GetTracks(),
            app.GetAlbums(),
            app.GetArtists()
        ]);
        
        // Update stats
        document.getElementById('track-count').textContent = tracks.length;
        document.getElementById('album-count').textContent = albums.length;
        document.getElementById('artist-count').textContent = artists.length;
        
        // Render lists
        renderTracks(tracks);
        renderAlbums(albums);
        renderArtists(artists);
        
        updateStatus('Library loaded successfully');
    } catch (error) {
        console.error('Failed to load library data:', error);
        updateStatus('Failed to load library data');
    }
}

// Render tracks
function renderTracks(tracks) {
    const container = document.getElementById('tracks-list');
    container.innerHTML = tracks.map(track => `
        <div class="list-item track-item" data-track-id="${track.id}">
            <div class="item-info">
                <div class="item-title">${track.title}</div>
                <div class="item-subtitle">${track.artist_name} - ${track.album_name}</div>
            </div>
            <div class="item-actions">
                <button class="btn btn-sm" onclick="playTrack(${track.id})">Play</button>
                <button class="btn btn-sm" onclick="downloadLyrics(${track.id})">Download Lyrics</button>
            </div>
        </div>
    `).join('');
}

// Render albums
function renderAlbums(albums) {
    const container = document.getElementById('albums-list');
    container.innerHTML = albums.map(album => `
        <div class="list-item album-item" data-album-id="${album.id}">
            <div class="item-info">
                <div class="item-title">${album.name}</div>
                <div class="item-subtitle">${album.artist_name} (${album.tracks_count} tracks)</div>
            </div>
        </div>
    `).join('');
}

// Render artists
function renderArtists(artists) {
    const container = document.getElementById('artists-list');
    container.innerHTML = artists.map(artist => `
        <div class="list-item artist-item" data-artist-id="${artist.id}">
            <div class="item-info">
                <div class="item-title">${artist.name}</div>
                <div class="item-subtitle">${artist.tracks_count} tracks</div>
            </div>
        </div>
    `).join('');
}

// Play track
window.playTrack = async function(trackId) {
    try {
        await app.PlayTrack(trackId);
        updateStatus('Playing track');
    } catch (error) {
        console.error('Failed to play track:', error);
        updateStatus('Failed to play track');
    }
};

// Download lyrics
window.downloadLyrics = async function(trackId) {
    try {
        updateStatus('Downloading lyrics...');
        const result = await app.DownloadLyrics(trackId);
        updateStatus(result);
    } catch (error) {
        console.error('Failed to download lyrics:', error);
        updateStatus('Failed to download lyrics');
    }
};

// Initialize library
document.getElementById('initialize-library').addEventListener('click', async function() {
    try {
        updateStatus('Initializing library...');
        await app.InitializeLibrary();
        showLibrarySection();
        await loadLibraryData();
    } catch (error) {
        console.error('Failed to initialize library:', error);
        updateStatus('Failed to initialize library');
    }
});

// Set directories
document.getElementById('set-directories').addEventListener('click', async function() {
    const directoriesInput = document.getElementById('directories');
    const directories = directoriesInput.value.split(',').map(dir => dir.trim()).filter(dir => dir);
    
    if (directories.length === 0) {
        updateStatus('Please enter at least one directory');
        return;
    }
    
    try {
        await app.SetDirectories(directories);
        updateStatus('Directories set successfully');
    } catch (error) {
        console.error('Failed to set directories:', error);
        updateStatus('Failed to set directories');
    }
});

// Tab switching
document.querySelectorAll('.tab-btn').forEach(btn => {
    btn.addEventListener('click', function() {
        const tab = this.dataset.tab;
        
        // Update active tab button
        document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
        this.classList.add('active');
        
        // Update active tab pane
        document.querySelectorAll('.tab-pane').forEach(pane => pane.classList.remove('active'));
        document.getElementById(`${tab}-tab`).classList.add('active');
    });
});

// Update status
function updateStatus(message) {
    document.getElementById('status-text').textContent = message;
}

// Initialize the app when the page loads
initializeApp();
