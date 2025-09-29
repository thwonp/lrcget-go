(function(){const i=document.createElement("link").relList;if(i&&i.supports&&i.supports("modulepreload"))return;for(const a of document.querySelectorAll('link[rel="modulepreload"]'))r(a);new MutationObserver(a=>{for(const n of a)if(n.type==="childList")for(const l of n.addedNodes)l.tagName==="LINK"&&l.rel==="modulepreload"&&r(l)}).observe(document,{childList:!0,subtree:!0});function t(a){const n={};return a.integrity&&(n.integrity=a.integrity),a.referrerpolicy&&(n.referrerPolicy=a.referrerpolicy),a.crossorigin==="use-credentials"?n.credentials="include":a.crossorigin==="anonymous"?n.credentials="omit":n.credentials="same-origin",n}function r(a){if(a.ep)return;a.ep=!0;const n=t(a);fetch(a.href,n)}})();const d="/assets/logo-universal.157a874a.png";document.querySelector("#app").innerHTML=`
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
`;document.getElementById("logo").src=d;async function u(){try{await app.GetInit()?(c(),await o()):p()}catch(e){console.error("Failed to initialize app:",e),s("Failed to initialize application")}}function p(){document.getElementById("setup-section").style.display="block",document.getElementById("library-section").style.display="none",document.getElementById("player-section").style.display="none"}function c(){document.getElementById("setup-section").style.display="none",document.getElementById("library-section").style.display="block",document.getElementById("player-section").style.display="block"}async function o(){try{s("Loading library data...");const[e,i,t]=await Promise.all([app.GetTracks(),app.GetAlbums(),app.GetArtists()]);document.getElementById("track-count").textContent=e.length,document.getElementById("album-count").textContent=i.length,document.getElementById("artist-count").textContent=t.length,y(e),m(i),b(t),s("Library loaded successfully")}catch(e){console.error("Failed to load library data:",e),s("Failed to load library data")}}function y(e){const i=document.getElementById("tracks-list");i.innerHTML=e.map(t=>`
        <div class="list-item track-item" data-track-id="${t.id}">
            <div class="item-info">
                <div class="item-title">${t.title}</div>
                <div class="item-subtitle">${t.artist_name} - ${t.album_name}</div>
            </div>
            <div class="item-actions">
                <button class="btn btn-sm" onclick="playTrack(${t.id})">Play</button>
                <button class="btn btn-sm" onclick="downloadLyrics(${t.id})">Download Lyrics</button>
            </div>
        </div>
    `).join("")}function m(e){const i=document.getElementById("albums-list");i.innerHTML=e.map(t=>`
        <div class="list-item album-item" data-album-id="${t.id}">
            <div class="item-info">
                <div class="item-title">${t.name}</div>
                <div class="item-subtitle">${t.artist_name} (${t.tracks_count} tracks)</div>
            </div>
        </div>
    `).join("")}function b(e){const i=document.getElementById("artists-list");i.innerHTML=e.map(t=>`
        <div class="list-item artist-item" data-artist-id="${t.id}">
            <div class="item-info">
                <div class="item-title">${t.name}</div>
                <div class="item-subtitle">${t.tracks_count} tracks</div>
            </div>
        </div>
    `).join("")}window.playTrack=async function(e){try{await app.PlayTrack(e),s("Playing track")}catch(i){console.error("Failed to play track:",i),s("Failed to play track")}};window.downloadLyrics=async function(e){try{s("Downloading lyrics...");const i=await app.DownloadLyrics(e);s(i)}catch(i){console.error("Failed to download lyrics:",i),s("Failed to download lyrics")}};document.getElementById("initialize-library").addEventListener("click",async function(){try{s("Initializing library..."),await app.InitializeLibrary(),c(),await o()}catch(e){console.error("Failed to initialize library:",e),s("Failed to initialize library")}});document.getElementById("set-directories").addEventListener("click",async function(){const i=document.getElementById("directories").value.split(",").map(t=>t.trim()).filter(t=>t);if(i.length===0){s("Please enter at least one directory");return}try{await app.SetDirectories(i),s("Directories set successfully")}catch(t){console.error("Failed to set directories:",t),s("Failed to set directories")}});document.querySelectorAll(".tab-btn").forEach(e=>{e.addEventListener("click",function(){const i=this.dataset.tab;document.querySelectorAll(".tab-btn").forEach(t=>t.classList.remove("active")),this.classList.add("active"),document.querySelectorAll(".tab-pane").forEach(t=>t.classList.remove("active")),document.getElementById(`${i}-tab`).classList.add("active")})});function s(e){document.getElementById("status-text").textContent=e}u();
