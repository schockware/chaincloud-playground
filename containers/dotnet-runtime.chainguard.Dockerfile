# Build context: src/
# podman build -f containers/dotnet-runtime.chainguard.Dockerfile -t weather-api:chainguard src/
#
# Researched image: cgr.dev/chainguard/dotnet-runtime (see IMAGE-DETAILS/dotnet-runtime/)
# If ASP.NET Core runtime components are missing at startup, switch BASE_IMAGE to
# cgr.dev/chainguard/aspnet-runtime:latest
ARG BASE_IMAGE=cgr.dev/chainguard/dotnet-runtime:latest

FROM cgr.dev/chainguard/dotnet-sdk:latest AS builder
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
