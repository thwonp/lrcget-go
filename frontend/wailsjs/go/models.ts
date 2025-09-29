export namespace audio {
	
	export class PlayerState {
	    status: number;
	    progress: number;
	    duration: number;
	    volume: number;
	    track?: database.PersistentTrack;
	
	    static createFrom(source: any = {}) {
	        return new PlayerState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.progress = source["progress"];
	        this.duration = source["duration"];
	        this.volume = source["volume"];
	        this.track = this.convertValues(source["track"], database.PersistentTrack);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace database {
	
	export class PersistentAlbum {
	    id: number;
	    name: string;
	    image_path?: string;
	    artist_name: string;
	    album_artist_name?: string;
	    name_lower?: string;
	    album_artist_name_lower?: string;
	    tracks_count: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new PersistentAlbum(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.image_path = source["image_path"];
	        this.artist_name = source["artist_name"];
	        this.album_artist_name = source["album_artist_name"];
	        this.name_lower = source["name_lower"];
	        this.album_artist_name_lower = source["album_artist_name_lower"];
	        this.tracks_count = source["tracks_count"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PersistentArtist {
	    id: number;
	    name: string;
	    name_lower?: string;
	    tracks_count: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new PersistentArtist(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.name_lower = source["name_lower"];
	        this.tracks_count = source["tracks_count"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PersistentConfig {
	    id: number;
	    skip_tracks_with_synced_lyrics: boolean;
	    skip_tracks_with_plain_lyrics: boolean;
	    show_line_count: boolean;
	    try_embed_lyrics: boolean;
	    theme_mode: string;
	    lrclib_instance: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new PersistentConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.skip_tracks_with_synced_lyrics = source["skip_tracks_with_synced_lyrics"];
	        this.skip_tracks_with_plain_lyrics = source["skip_tracks_with_plain_lyrics"];
	        this.show_line_count = source["show_line_count"];
	        this.try_embed_lyrics = source["try_embed_lyrics"];
	        this.theme_mode = source["theme_mode"];
	        this.lrclib_instance = source["lrclib_instance"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PersistentTrack {
	    id: number;
	    file_path: string;
	    file_name: string;
	    title: string;
	    album_name: string;
	    album_artist_name?: string;
	    album_id: number;
	    artist_name: string;
	    artist_id: number;
	    image_path?: string;
	    track_number?: number;
	    txt_lyrics?: string;
	    lrc_lyrics?: string;
	    duration: number;
	    instrumental: boolean;
	    title_lower?: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new PersistentTrack(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.file_path = source["file_path"];
	        this.file_name = source["file_name"];
	        this.title = source["title"];
	        this.album_name = source["album_name"];
	        this.album_artist_name = source["album_artist_name"];
	        this.album_id = source["album_id"];
	        this.artist_name = source["artist_name"];
	        this.artist_id = source["artist_id"];
	        this.image_path = source["image_path"];
	        this.track_number = source["track_number"];
	        this.txt_lyrics = source["txt_lyrics"];
	        this.lrc_lyrics = source["lrc_lyrics"];
	        this.duration = source["duration"];
	        this.instrumental = source["instrumental"];
	        this.title_lower = source["title_lower"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace lrclib {
	
	export class PublishResponse {
	    id: number;
	    token: string;
	
	    static createFrom(source: any = {}) {
	        return new PublishResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.token = source["token"];
	    }
	}
	export class SearchResult {
	    id: number;
	    trackName: string;
	    artistName: string;
	    albumName: string;
	    duration: number;
	    syncedLyrics?: string;
	    plainLyrics?: string;
	    instrumental: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.trackName = source["trackName"];
	        this.artistName = source["artistName"];
	        this.albumName = source["albumName"];
	        this.duration = source["duration"];
	        this.syncedLyrics = source["syncedLyrics"];
	        this.plainLyrics = source["plainLyrics"];
	        this.instrumental = source["instrumental"];
	    }
	}
	export class SearchResponse {
	    data: SearchResult[];
	
	    static createFrom(source: any = {}) {
	        return new SearchResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], SearchResult);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

