var builder = DistributedApplication.CreateBuilder(args);

var owmApiKey = builder.AddParameter("OpenWeatherMapApiKey", secret: true);
var spotifyClientId = builder.AddParameter("SpotifyClientId", secret: true);
var spotifyClientSecret = builder.AddParameter("SpotifyClientSecret", secret: true);
var spotifyRefreshToken = builder.AddParameter("SpotifyRefreshToken", secret: true);

// Builds from ./src/playlist-engine/Dockerfile on every run via the Podman socket.
// For CVE experiment runs with a pre-built image, swap to AddContainer("playlist-engine", "playlist-engine", "chainguard").
var playlistEngine = builder.AddDockerfile("playlist-engine", "../playlist-engine")
    .WithHttpEndpoint(targetPort: 5100, name: "http")
    .WithEnvironment("SPOTIFY_CLIENT_ID", spotifyClientId)
    .WithEnvironment("SPOTIFY_CLIENT_SECRET", spotifyClientSecret)
    .WithEnvironment("SPOTIFY_REFRESH_TOKEN", spotifyRefreshToken);

builder.AddProject<Projects.WeatherPlaylist_Api>("weatherplaylist-api")
    .WithEnvironment("OPENWEATHERMAP_API_KEY", owmApiKey)
    .WithReference(playlistEngine.GetEndpoint("http"))
    .WaitFor(playlistEngine);

builder.Build().Run();
