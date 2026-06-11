var builder = DistributedApplication.CreateBuilder(args);

var owmApiKey = builder.AddParameter("OpenWeatherMapApiKey", secret: true);
var spotifyClientId = builder.AddParameter("SpotifyClientId", secret: true);
var spotifyClientSecret = builder.AddParameter("SpotifyClientSecret", secret: true);
var spotifyRefreshToken = builder.AddParameter("SpotifyRefreshToken", secret: true);

var useMockServers = builder.Configuration["USE_MOCK_SERVERS"] == "true";

// Builds from src/playlist-engine/Dockerfile on every run via the Podman socket.
// For CVE experiment runs with a pre-built image, swap to AddContainer("playlist-engine", "playlist-engine", "chainguard").
var playlistEngine = builder.AddDockerfile("playlist-engine", "../playlist-engine")
    .WithHttpEndpoint(targetPort: 5100, name: "http")
    .WithEnvironment("SPOTIFY_CLIENT_ID", spotifyClientId)
    .WithEnvironment("SPOTIFY_CLIENT_SECRET", spotifyClientSecret)
    .WithEnvironment("SPOTIFY_REFRESH_TOKEN", spotifyRefreshToken);

var weatherApi = builder.AddProject<Projects.WeatherPlaylist_Api>("weatherplaylist-api")
    .WithEnvironment("OPENWEATHERMAP_API_KEY", owmApiKey)
    .WithReference(playlistEngine.GetEndpoint("http"))
    .WaitFor(playlistEngine);

if (useMockServers)
{
    var certsPath = Path.GetFullPath(
        Path.Combine(builder.AppHostDirectory, "..", "..", "containers", "certs"));

    var mockSpotify = builder.AddDockerfile(
            "mock-spotify", "../../src/mock-spotify", "../containers/mock-spotify.Dockerfile")
        .WithHttpsEndpoint(targetPort: 5200, name: "https")
        .WithBindMount(certsPath, "/certs", isReadOnly: true)
        .WithEnvironment("MOCK_TLS_CERT_FILE", "/certs/mock-spotify.crt")
        .WithEnvironment("MOCK_TLS_KEY_FILE", "/certs/mock-spotify.key");

    var mockOwm = builder.AddDockerfile(
            "mock-owm", "../../src/mock-owm", "../containers/mock-owm.Dockerfile")
        .WithHttpsEndpoint(targetPort: 5300, name: "https")
        .WithBindMount(certsPath, "/certs", isReadOnly: true)
        .WithEnvironment("MOCK_TLS_CERT_FILE", "/certs/mock-owm.crt")
        .WithEnvironment("MOCK_TLS_KEY_FILE", "/certs/mock-owm.key");

    playlistEngine
        .WaitFor(mockSpotify)
        .WithBindMount(certsPath, "/certs", isReadOnly: true)
        .WithEnvironment("SPOTIFY_MOCK_BASE_URL", mockSpotify.GetEndpoint("https"))
        .WithEnvironment("SPOTIFY_TLS_CA_FILE", "/certs/ca.crt");

    weatherApi
        .WaitFor(mockOwm)
        .WithEnvironment("OWM_BASE_URL", mockOwm.GetEndpoint("https"))
        .WithEnvironment("OWM_TLS_CA_FILE", Path.Combine(certsPath, "ca.crt"));
}

builder.Build().Run();
