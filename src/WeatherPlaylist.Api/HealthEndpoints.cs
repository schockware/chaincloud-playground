using WeatherPlaylist.Api.Services;

namespace WeatherPlaylist.Api;

public static class HealthEndpoints
{
    public static void MapHealthEndpoints(this WebApplication app)
    {
        app.MapGet("/health", HandleAsync);
    }

    private static async Task<IResult> HandleAsync(
        IWeatherService weatherService,
        PlaylistEngineClient engineClient,
        CancellationToken ct)
    {
        string owmStatus;
        if (weatherService.IsMock)
        {
            owmStatus = "mock";
        }
        else
        {
            try
            {
                await weatherService.GetCurrentAsync(51.5074, -0.1278, ct);
                owmStatus = "ok";
            }
            catch
            {
                owmStatus = "unreachable";
            }
        }

        var engineStatus = await engineClient.IsHealthyAsync(ct) ? "ok" : "unreachable";
        var overall = owmStatus is "ok" or "mock" && engineStatus == "ok" ? "ok" : "degraded";

        return Results.Ok(new HealthResponse(overall, owmStatus, engineStatus));
    }
}
