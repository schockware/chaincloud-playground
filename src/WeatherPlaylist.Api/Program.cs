using System.Security.Cryptography.X509Certificates;
using System.Text.Json;
using Microsoft.AspNetCore.Diagnostics.HealthChecks;
using WeatherPlaylist.Api;
using WeatherPlaylist.Api.Services;

var builder = WebApplication.CreateBuilder(args);

builder.AddServiceDefaults();
builder.Services.AddOpenApi();

builder.Services.ConfigureHttpJsonOptions(opts =>
    opts.SerializerOptions.PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower);

var owmBaseUrl = builder.Configuration["OWM_BASE_URL"] ?? "https://api.openweathermap.org";
var owmCaFile = builder.Configuration["OWM_TLS_CA_FILE"];

builder.Services.AddHttpClient("openweathermap",
    client => client.BaseAddress = new Uri(owmBaseUrl))
    .ConfigurePrimaryHttpMessageHandler(() => BuildHttpHandler(owmCaFile));

builder.Services.AddHttpClient("playlist-engine",
    client => client.BaseAddress = new Uri("http://playlist-engine"));

var owmApiKey = builder.Configuration["OPENWEATHERMAP_API_KEY"];
var useMockWeather = string.IsNullOrWhiteSpace(owmApiKey);

if (useMockWeather)
    builder.Services.AddScoped<IWeatherService, MockWeatherService>();
else
    builder.Services.AddScoped<IWeatherService, WeatherService>();

builder.Services.AddScoped<PlaylistEngineClient>();

var app = builder.Build();

if (app.Environment.IsDevelopment())
    app.MapOpenApi();

app.MapHealthChecks("/alive", new HealthCheckOptions
{
    Predicate = r => r.Tags.Contains("live")
});

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

        var mockParts = new List<string>();
        if (useMockWeather) mockParts.Add("weather");
        if (ctx.Items.ContainsKey("SpotifyMocked")) mockParts.Add("spotify");
        if (mockParts.Count > 0)
            ctx.Response.Headers["X-ARBITRARY-MOCK"] = string.Join(",", mockParts);

        return Task.CompletedTask;
    });

    await next(ctx);
});

app.MapWeatherEndpoints();
app.MapPlaylistEndpoints();
app.MapHealthEndpoints();

app.Run();

static HttpClientHandler BuildHttpHandler(string? caFile)
{
    var handler = new HttpClientHandler();
    if (string.IsNullOrWhiteSpace(caFile)) return handler;

    var caCert = X509CertificateLoader.LoadCertificateFromFile(caFile);
    handler.ServerCertificateCustomValidationCallback = (_, cert, chain, _) =>
    {
        if (chain is null || cert is null) return false;
        chain.ChainPolicy.TrustMode = X509ChainTrustMode.CustomRootTrust;
        chain.ChainPolicy.CustomTrustStore.Add(caCert);
        chain.ChainPolicy.RevocationMode = X509RevocationMode.NoCheck;
        return chain.Build(cert);
    };
    return handler;
}
