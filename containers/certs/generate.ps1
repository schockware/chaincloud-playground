# Generate local CA and leaf certs for mock services.
# Idempotent — skips files that already exist.
# Requires: openssl in PATH
# Run once from the repo root: .\containers\certs\generate.ps1

$dir = $PSScriptRoot

function Show-Status($path) {
    if (Test-Path $path) { Write-Host "EXISTS  $path"; return $true }
    return $false
}

# CA key + self-signed cert
if (-not (Show-Status "$dir\ca.key")) {
    & openssl genrsa -out "$dir\ca.key" 2048
}
if (-not (Show-Status "$dir\ca.crt")) {
    & openssl req -new -x509 -key "$dir\ca.key" -out "$dir\ca.crt" -days 3650 `
        -subj "/CN=MockLocalCA/O=ChaincloudPlayground"
    Write-Host "CREATED $dir\ca.crt"
}

function New-LeafCert {
    param([string]$Name, [string]$Port)

    $key = "$dir\$Name.key"
    $csr = "$dir\$Name.csr"
    $crt = "$dir\$Name.crt"
    $ext = "$dir\$Name.ext"

    if ((Test-Path $key) -and (Test-Path $crt)) {
        Write-Host "EXISTS  $key + $crt"
        return
    }

    @"
[v3_req]
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
IP.1 = 127.0.0.1
"@ | Out-File -Encoding utf8 $ext

    & openssl genrsa -out $key 2048
    & openssl req -new -key $key -out $csr -subj "/CN=localhost/O=ChaincloudPlayground"
    & openssl x509 -req -in $csr -CA "$dir\ca.crt" -CAkey "$dir\ca.key" `
        -CAcreateserial -out $crt -days 3650 -extfile $ext -extensions v3_req
    Remove-Item $csr, $ext -ErrorAction SilentlyContinue
    Write-Host "CREATED $crt"
}

New-LeafCert -Name "mock-spotify" -Port "5200"
New-LeafCert -Name "mock-owm"     -Port "5300"

Write-Host "`nDone. Commit *.crt files. Private *.key files are gitignored."
