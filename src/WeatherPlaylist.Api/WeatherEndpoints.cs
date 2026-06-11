using WeatherPlaylist.Api.Services;

namespace WeatherPlaylist.Api;

public static class WeatherEndpoints
{
    public static void MapWeatherEndpoints(this WebApplication app)
    {
        app.MapGet("/weather", HandleAsync);
    }

    private static async Task<IResult> HandleAsync(
        double lat,
        double lon,
        IWeatherService weatherService,
        HttpContext ctx,
        CancellationToken ct)
    {
        if (lat is < -90 or > 90 || lon is < -180 or > 180)
            return Results.Problem(
                title: "Invalid coordinates",
                detail: "lat must be in [-90, 90], lon in [-180, 180]",
                statusCode: 400);

        try
        {
            return Results.Ok(await weatherService.GetCurrentAsync(lat, lon, ct));
        }
        catch (HttpRequestException ex)
        {
            var correlationId = ctx.Items["CorrelationId"]?.ToString();
            return Results.Problem(
                title: "Upstream error from OpenWeatherMap",
                detail: ex.Message,
                statusCode: 502,
                extensions: correlationId is null ? null
                    : new Dictionary<string, object?> { ["correlation_id"] = correlationId });
        }
    }
}
