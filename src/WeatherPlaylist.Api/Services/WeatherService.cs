using System.Text.Json;
using System.Text.Json.Serialization;

namespace WeatherPlaylist.Api.Services;

public class WeatherService(IHttpClientFactory httpClientFactory, IConfiguration config) : IWeatherService
{
    private static readonly JsonSerializerOptions JsonOpts = new() { PropertyNameCaseInsensitive = true };

    public bool IsMock => false;

    public async Task<WeatherCondition> GetCurrentAsync(double lat, double lon, CancellationToken ct = default)
    {
        var apiKey = config["OPENWEATHERMAP_API_KEY"]
            ?? throw new InvalidOperationException("OPENWEATHERMAP_API_KEY not configured");

        var client = httpClientFactory.CreateClient("openweathermap");
        var response = await client.GetAsync(
            $"/data/2.5/weather?lat={lat}&lon={lon}&appid={apiKey}&units=metric", ct);

        if (!response.IsSuccessStatusCode)
            throw new HttpRequestException(
                $"OpenWeatherMap returned {(int)response.StatusCode}", null, response.StatusCode);

        var owm = await response.Content.ReadFromJsonAsync<OwmResponse>(JsonOpts, ct)
            ?? throw new InvalidOperationException("Empty response from OpenWeatherMap");

        return Map(owm);
    }

    private static WeatherCondition Map(OwmResponse owm)
    {
        var mainCode = owm.Weather.FirstOrDefault()?.Main ?? "Clear";
        var condition = mainCode.ToLowerInvariant() switch
        {
            "clouds" => "clouds",
            "rain" => "rain",
            "drizzle" => "drizzle",
            "thunderstorm" => "thunderstorm",
            "snow" => "snow",
            "mist" => "mist",
            "fog" => "fog",
            "haze" => "haze",
            _ => "clear"
        };

        var localUnix = owm.Dt + owm.Timezone;
        var localHour = (int)((localUnix % 86400) / 3600);
        if (localHour < 0) localHour += 24;
        var timeOfDay = localHour switch
        {
            >= 6 and < 12 => "morning",
            >= 12 and < 18 => "afternoon",
            >= 18 and < 21 => "evening",
            _ => "night"
        };

        return new WeatherCondition(
            condition,
            owm.Main.Temp,
            owm.Main.FeelsLike,
            owm.Main.Humidity,
            owm.Wind.Speed,
            owm.Weather.FirstOrDefault()?.Description ?? "",
            timeOfDay);
    }

    private record OwmResponse(
        [property: JsonPropertyName("weather")] OwmWeather[] Weather,
        [property: JsonPropertyName("main")] OwmMain Main,
        [property: JsonPropertyName("wind")] OwmWind Wind,
        [property: JsonPropertyName("timezone")] long Timezone,
        [property: JsonPropertyName("dt")] long Dt);

    private record OwmWeather(
        [property: JsonPropertyName("main")] string Main,
        [property: JsonPropertyName("description")] string Description);

    private record OwmMain(
        [property: JsonPropertyName("temp")] double Temp,
        [property: JsonPropertyName("feels_like")] double FeelsLike,
        [property: JsonPropertyName("humidity")] int Humidity);

    private record OwmWind([property: JsonPropertyName("speed")] double Speed);
}
