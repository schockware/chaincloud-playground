namespace WeatherPlaylist.Api;

public record WeatherCondition(
    string Condition,
    double TempCelsius,
    double FeelsLikeCelsius,
    int Humidity,
    double WindSpeedMps,
    string Description,
    string TimeOfDay);

public record TempoRange(int MinBpm, int MaxBpm);

public record PlaylistRecipe(
    string Mood,
    string[] Genres,
    TempoRange TempoRange,
    double EnergyLevel,
    double Valence,
    int TrackCount);

public record GeneratePlaylistRequest(
    double? Lat,
    double? Lon,
    string? LocationLabel,
    string? ExperimentId);

public record PlaylistResponse(
    string PlaylistId,
    string SpotifyEmbedUrl,
    string PlaylistName,
    WeatherCondition WeatherSnapshot,
    PlaylistRecipe Recipe,
    string? ExperimentId,
    DateTimeOffset CreatedAt);

public record GoPlaylistResult(
    string PlaylistId,
    string SpotifyEmbedUrl,
    string PlaylistName,
    PlaylistRecipe Recipe,
    string? ExperimentId,
    int TracksAdded,
    DateTimeOffset CreatedAt);

public record HealthResponse(
    string Status,
    string? UpstreamWeather,
    string? UpstreamPlaylistEngine);
