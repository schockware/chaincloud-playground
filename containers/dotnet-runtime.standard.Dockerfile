# Build context: src/
# podman build -f containers/dotnet-runtime.standard.Dockerfile -t weather-api:standard src/
#
# Standard alternative to cgr.dev/chainguard/dotnet-runtime (see IMAGE-DETAILS/dotnet-runtime/)
# Uses mcr.microsoft.com/dotnet/aspnet (ASP.NET Core runtime) rather than dotnet/runtime
# because WeatherPlaylist.Api requires the ASP.NET Core shared framework.
ARG BASE_IMAGE=mcr.microsoft.com/dotnet/aspnet:10.0

FROM mcr.microsoft.com/dotnet/sdk:10.0 AS builder
WORKDIR /build
COPY WeatherPlaylist.ServiceDefaults/ WeatherPlaylist.ServiceDefaults/
COPY WeatherPlaylist.Api/ WeatherPlaylist.Api/
WORKDIR /build/WeatherPlaylist.Api
RUN dotnet publish -c Release -o /publish

FROM ${BASE_IMAGE}
WORKDIR /app
COPY --from=builder /publish .
EXPOSE 8080
ENTRYPOINT ["dotnet", "WeatherPlaylist.Api.dll"]
