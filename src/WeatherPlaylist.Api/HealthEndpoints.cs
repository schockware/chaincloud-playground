using WeatherPlaylist.Api.Services;

namespace WeatherPlaylist.Api;

public static class HealthEndpoints
{
    public static void MapHealthEndpoints(this WebApplication app)
    {
        app.MapGet("/health", HandleAsync);
    }

    private static async Task<IResult> HandleAsync(
        WeatherService weatherService,
        PlaylistEngineClient engineClient,
        IConfiguration config,
        CancellationToken ct)
    {
        string owmStatus;
        try
        {
            if (string.IsNullOrEmpty(config["OPENWEATHERMAP_API_KEY"]))
            {
                owmStatus = "unreachable";
            }
            else
            {
                await weatherService.GetCurrentAsync(51.5074, -0.1278, ct);
                owmStatus = "ok";
            }
        }
        catch
        {
            owmStatus = "unreachable";
        }

        var engineStatus = await engineClient.IsHealthyAsync(ct) ? "ok" : "unreachable";
        var overall = owmStatus == "ok" && engineStatus == "ok" ? "ok" : "degraded";

        return Results.Ok(new HealthResponse(overall, owmStatus, engineStatus));
    }
}
