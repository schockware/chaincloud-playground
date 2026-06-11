namespace WeatherPlaylist.Api.Services;

public interface IWeatherService
{
    bool IsMock { get; }
    Task<WeatherCondition> GetCurrentAsync(double lat, double lon, CancellationToken ct = default);
}
