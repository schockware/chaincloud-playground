using System.Text.Json;
using Microsoft.AspNetCore.Diagnostics.HealthChecks;
using WeatherPlaylist.Api;
using WeatherPlaylist.Api.Services;

var builder = WebApplication.CreateBuilder(args);

builder.AddServiceDefaults();
builder.Services.AddOpenApi();

builder.Services.ConfigureHttpJsonOptions(opts =>
    opts.SerializerOptions.PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower);

builder.Services.AddHttpClient("openweathermap",
    client => client.BaseAddress = new Uri("https://api.openweathermap.org"));

builder.Services.AddHttpClient("playlist-engine",
    client => client.BaseAddress = new Uri("http://playlist-engine"));

builder.Services.AddScoped<WeatherService>();
builder.Services.AddScoped<PlaylistEngineClient>();

var app = builder.Build();

if (app.Environment.IsDevelopment())
    app.MapOpenApi();

// Aspire liveness probe (separate from our application /health endpoint)
app.MapHealthChecks("/alive", new HealthCheckOptions
{
    Predicate = r => r.Tags.Contains("live")
});

// Generate or inherit X-Correlation-Id; echo it and X-Experiment-Id on every response
app.Use(async (ctx, next) =>
{
    var correlationId = ctx.Request.Headers["X-Correlation-Id"].FirstOrDefault()
        ?? Guid.NewGuid().ToString();
    ctx.Items["CorrelationId"] = correlationId;

    var experimentId = ctx.Request.Headers["X-Experiment-Id"].FirstOrDefault();
    if (!string.IsNullOrEmpty(experimentId))
        ctx.Items["ExperimentId"] = experimentId;

    ctx.Response.OnStarting(() =>
    {
        ctx.Response.Headers["X-Correlation-Id"] = correlationId;
        if (!string.IsNullOrEmpty(experimentId))
            ctx.Response.Headers["X-Experiment-Id"] = experimentId;
        return Task.CompletedTask;
    });

    await next(ctx);
});

app.MapWeatherEndpoints();
app.MapPlaylistEndpoints();
app.MapHealthEndpoints();

app.Run();
