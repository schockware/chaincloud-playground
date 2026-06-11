using System.Text.Json;
using WeatherPlaylist.Api.Services;

namespace WeatherPlaylist.Api;

public static class PlaylistEndpoints
{
    private static readonly JsonSerializerOptions GoJsonOpts = new()
    {
        PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower,
        PropertyNameCaseInsensitive = true
    };

    public static void MapPlaylistEndpoints(this WebApplication app)
    {
        app.MapPost("/playlist/generate", HandleGenerateAsync);
    }

    private static async Task HandleGenerateAsync(
        GeneratePlaylistRequest req,
        IWeatherService weatherService,
        PlaylistEngineClient engineClient,
        HttpContext ctx,
        CancellationToken ct)
    {
        var correlationId = ctx.Items["CorrelationId"]?.ToString() ?? Guid.NewGuid().ToString();

        if (req.Lat is null or < -90 or > 90 || req.Lon is null or < -180 or > 180)
        {
            ctx.Response.StatusCode = 400;
            await ctx.Response.WriteAsJsonAsync(new
            {
                type = "https://tools.ietf.org/html/rfc9457",
                title = "Invalid or missing coordinates",
                status = 400,
                correlation_id = correlationId
            }, cancellationToken: ct);
            return;
        }

        WeatherCondition weather;
        try
        {
            weather = await weatherService.GetCurrentAsync(req.Lat.Value, req.Lon.Value, ct);
        }
        catch (HttpRequestException ex)
        {
            ctx.Response.StatusCode = 502;
            await ctx.Response.WriteAsJsonAsync(new
            {
                type = "https://tools.ietf.org/html/rfc9457",
                title = "Upstream error fetching weather",
                status = 502,
                detail = ex.Message,
                correlation_id = correlationId
            }, cancellationToken: ct);
            return;
        }

        var arbitraryMockHeader = ctx.Request.Headers["X-ARBITRARY-MOCK"].FirstOrDefault();

        var goResp = await engineClient.GenerateAsync(
            weather, req.LocationLabel, req.ExperimentId, correlationId, arbitraryMockHeader, ct);

        if (goResp.Headers.TryGetValues("X-ARBITRARY-MOCK", out var goMockVals)
            && string.Join(",", goMockVals).Contains("spotify"))
        {
            ctx.Items["SpotifyMocked"] = true;
        }

        if (!goResp.IsSuccessStatusCode)
        {
            ctx.Response.StatusCode = (int)goResp.StatusCode;
            ctx.Response.Headers["X-Api-Source"] = "GO-API-PASSTHROUGH";
            if (goResp.StatusCode == System.Net.HttpStatusCode.TooManyRequests
                && goResp.Headers.RetryAfter?.Delta is { } delta)
            {
                ctx.Response.Headers["Retry-After"] = ((int)delta.TotalSeconds).ToString();
            }
            ctx.Response.ContentType = "application/json";
            await goResp.Content.CopyToAsync(ctx.Response.Body, ct);
            return;
        }

        var goResult = await goResp.Content.ReadFromJsonAsync<GoPlaylistResult>(GoJsonOpts, ct)
            ?? throw new InvalidOperationException("Empty response from playlist-engine");

        var response = new PlaylistResponse(
            goResult.PlaylistId,
            goResult.SpotifyEmbedUrl,
            goResult.PlaylistName,
            weather,
            goResult.Recipe,
            goResult.ExperimentId,
            goResult.CreatedAt);

        ctx.Response.StatusCode = 200;
        await ctx.Response.WriteAsJsonAsync(response, GoJsonOpts, ct);
    }
}
