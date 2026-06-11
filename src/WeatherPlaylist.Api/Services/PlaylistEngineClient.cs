namespace WeatherPlaylist.Api.Services;

public class PlaylistEngineClient(IHttpClientFactory httpClientFactory)
{
    public async Task<HttpResponseMessage> GenerateAsync(
        WeatherCondition weather,
        string? locationLabel,
        string? experimentId,
        string correlationId,
        string? arbitraryMockHeader = null,
        CancellationToken ct = default)
    {
        var client = httpClientFactory.CreateClient("playlist-engine");

        var body = new
        {
            weather = new
            {
                condition = weather.Condition,
                temp_celsius = weather.TempCelsius,
                feels_like_celsius = weather.FeelsLikeCelsius,
                humidity = weather.Humidity,
                wind_speed_mps = weather.WindSpeedMps,
                description = weather.Description,
                time_of_day = weather.TimeOfDay
            },
            location_label = locationLabel,
            experiment_id = experimentId
        };

        var req = new HttpRequestMessage(HttpMethod.Post, "/playlist/generate")
        {
            Content = JsonContent.Create(body)
        };
        req.Headers.TryAddWithoutValidation("X-Correlation-Id", correlationId);
        if (!string.IsNullOrEmpty(experimentId))
            req.Headers.TryAddWithoutValidation("X-Experiment-Id", experimentId);
        if (!string.IsNullOrEmpty(arbitraryMockHeader))
            req.Headers.TryAddWithoutValidation("X-ARBITRARY-MOCK", arbitraryMockHeader);

        return await client.SendAsync(req, HttpCompletionOption.ResponseContentRead, ct);
    }

    public async Task<bool> IsHealthyAsync(CancellationToken ct = default)
    {
        try
        {
            var client = httpClientFactory.CreateClient("playlist-engine");
            var response = await client.GetAsync("/health", ct);
            return response.IsSuccessStatusCode;
        }
        catch
        {
            return false;
        }
    }
}
