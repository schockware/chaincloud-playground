namespace WeatherPlaylist.Api.Services;

public class MockWeatherService : IWeatherService
{
    public bool IsMock => true;

    private static readonly WeatherCondition[] Conditions =
    [
        new("clear",       22.0,  21.0, 40, 3.5, "clear sky",                        "afternoon"),
        new("clouds",      15.0,  14.0, 65, 5.2, "overcast clouds",                  "morning"),
        new("rain",        12.0,  10.5, 85, 7.1, "moderate rain",                    "evening"),
        new("thunderstorm",18.0,  17.0, 78, 9.3, "thunderstorm with rain",            "night"),
        new("snow",        -3.0,  -6.0, 90, 4.0, "light snow",                        "morning"),
        new("mist",        10.0,   9.0, 95, 1.5, "mist",                             "night"),
        new("drizzle",     14.0,  13.0, 80, 3.0, "light intensity drizzle",           "afternoon"),
    ];

    public Task<WeatherCondition> GetCurrentAsync(double lat, double lon, CancellationToken ct = default)
    {
        // Deterministic per location + UTC hour so results stay stable within a window.
        var idx = Math.Abs(HashCode.Combine((int)(lat * 10), (int)(lon * 10), DateTime.UtcNow.Hour)) % Conditions.Length;
        return Task.FromResult(Conditions[idx]);
    }
}
