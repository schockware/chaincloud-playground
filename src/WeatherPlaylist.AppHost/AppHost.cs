var builder = DistributedApplication.CreateBuilder(args);

var owmApiKey = builder.AddParameter("OpenWeatherMapApiKey", secret: true);
var spotifyClientId = builder.AddParameter("SpotifyClientId", secret: true);
var spotifyClientSecret = builder.AddParameter("SpotifyClientSecret", secret: true);
var spotifyRefreshToken = builder.AddParameter("SpotifyRefreshToken", secret: true);

// Build playlist-engine with Chainguard base (default) or swap tag for standard-image experiments.
// Before first run: podman build -t playlist-engine:local ./src/playlist-engine
var playlistEngine = builder.AddContainer("playlist-engine", "playlist-engine", "local")
    .WithHttpEndpoint(targetPort: 5100, name: "http")
    .WithEnvironment("SPOTIFY_CLIENT_ID", spotifyClientId)
    .WithEnvironment("SPOTIFY_CLIENT_SECRET", spotifyClientSecret)
    .WithEnvironment("SPOTIFY_REFRESH_TOKEN", spotifyRefreshToken);

builder.AddProject<Projects.WeatherPlaylist_Api>("weatherplaylist-api")
    .WithEnvironment("OPENWEATHERMAP_API_KEY", owmApiKey)
    .WithReference(playlistEngine.GetEndpoint("http"))
    .WaitFor(playlistEngine);

builder.Build().Run();
