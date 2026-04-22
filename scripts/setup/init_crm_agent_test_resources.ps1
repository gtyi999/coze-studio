param(
    [string]$BaseUrl = "http://localhost:8888",
    [string]$MySQLContainer = "coze-mysql",
    [string]$MySQLDatabase = "opencoze",
    [string]$MySQLUser = "root",
    [string]$MySQLPassword = "",
    [string]$OutputFile = ""
)

$ErrorActionPreference = "Stop"

$script:RootDir = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)

if ([string]::IsNullOrWhiteSpace($OutputFile)) {
    $OutputFile = Join-Path $script:RootDir "output\crm-agent-test-resources.json"
}

$schemaSqlFile = Join-Path $script:RootDir "backend\types\ddl\crm_phase1.sql"
$demoSqlTemplateFile = Join-Path $script:RootDir "backend\types\ddl\crm_demo_data.sql.tpl"
$dockerEnvFile = Join-Path $script:RootDir "docker\.env"

$testUsers = @(
    [pscustomobject][ordered]@{
        key         = "crm-agent-owner"
        email       = "crm-agent-owner@example.com"
        password    = "Passw0rd!123"
        role        = "owner"
        description = "Primary CRM agent smoke-test account."
    },
    [pscustomobject][ordered]@{
        key         = "crm-agent-sales"
        email       = "crm-agent-sales@example.com"
        password    = "Passw0rd!123"
        role        = "sales"
        description = "Sales-oriented CRM agent validation account."
    },
    [pscustomobject][ordered]@{
        key         = "crm-agent-analyst"
        email       = "crm-agent-analyst@example.com"
        password    = "Passw0rd!123"
        role        = "analyst"
        description = "Read-heavy CRM dashboard and NL query account."
    }
)

$exampleQuestions = @(
    [pscustomobject][ordered]@{
        key             = "customer-count"
        language        = "en-US"
        title           = "Customer Count"
        question        = "How many customers do I have now?"
        expected_intent = "customer_count"
        note            = "Current frontend MVP question."
    },
    [pscustomobject][ordered]@{
        key             = "top-sales"
        language        = "en-US"
        title           = "Top Sales"
        question        = "Which sales rep has the best performance?"
        expected_intent = "top_sales_current_quarter"
        note            = "Exercises the ranking-style answer path."
    },
    [pscustomobject][ordered]@{
        key             = "sales-top-5"
        language        = "en-US"
        title           = "Sales Top 5"
        question        = "Show me the top five sales reps this quarter."
        expected_intent = "top_sales_current_quarter"
        note            = "Alternative ranking example for the same MVP path."
    },
    [pscustomobject][ordered]@{
        key             = "forecast"
        language        = "en-US"
        title           = "Hot Product Forecast"
        question        = "Which product will sell best next quarter?"
        expected_intent = "forecast_hot_product"
        note            = "Exercises the forecast-style answer path."
    },
    [pscustomobject][ordered]@{
        key             = "forecast-explainer"
        language        = "en-US"
        title           = "Forecast Explanation"
        question        = "Why do you think that product will sell best next quarter?"
        expected_intent = "forecast_hot_product"
        note            = "Useful for checking the answer plus disclaimer area."
    }
)

function Write-Step {
    param([string]$Message)
    Write-Host "==> $Message"
}

function Get-DefaultMySQLPassword {
    if (-not (Test-Path $dockerEnvFile)) {
        return "root"
    }

    $rootPasswordLine = Select-String -Path $dockerEnvFile -Pattern '^MYSQL_ROOT_PASSWORD=' | Select-Object -First 1
    if ($null -eq $rootPasswordLine) {
        return "root"
    }

    $value = ($rootPasswordLine.Line -split '=', 2)[1].Trim()
    if ([string]::IsNullOrWhiteSpace($value)) {
        return "root"
    }

    return $value
}

function Invoke-CurlJson {
    param(
        [Parameter(Mandatory = $true)][string]$Method,
        [Parameter(Mandatory = $true)][string]$Url,
        [string]$Body = "",
        [hashtable]$Headers = @{}
    )

    $headerFile = [System.IO.Path]::GetTempFileName()
    $bodyFile = [System.IO.Path]::GetTempFileName()
    $payloadFile = $null

    try {
        $curlArgs = @("-s", "-D", $headerFile, "-o", $bodyFile, "-X", $Method)
        foreach ($entry in $Headers.GetEnumerator()) {
            $curlArgs += @("-H", "$($entry.Key): $($entry.Value)")
        }
        if ($Method -ne "GET") {
            $payloadFile = [System.IO.Path]::GetTempFileName()
            Set-Content -Path $payloadFile -Value $Body -Encoding Ascii
            $curlArgs += @("--data-binary", "@$payloadFile")
        }
        $curlArgs += $Url

        & curl.exe @curlArgs | Out-Null

        $rawHeaders = Get-Content -Raw $headerFile
        $rawBody = Get-Content -Raw $bodyFile
        $statusMatches = [regex]::Matches($rawHeaders, 'HTTP/\d(?:\.\d)?\s+(\d{3})')
        if ($statusMatches.Count -eq 0) {
            throw "Unable to parse HTTP status for $Method $Url"
        }

        return [pscustomobject][ordered]@{
            StatusCode = [int]$statusMatches[$statusMatches.Count - 1].Groups[1].Value
            Headers    = $rawHeaders
            Body       = $rawBody
        }
    } finally {
        Remove-Item -Force $headerFile, $bodyFile -ErrorAction SilentlyContinue
        if ($null -ne $payloadFile) {
            Remove-Item -Force $payloadFile -ErrorAction SilentlyContinue
        }
    }
}

function Convert-ResponseToJson {
    param(
        [Parameter(Mandatory = $true)]$Response,
        [Parameter(Mandatory = $true)][string]$Context
    )

    if ($Response.StatusCode -lt 200 -or $Response.StatusCode -ge 300) {
        throw "$Context failed with HTTP $($Response.StatusCode): $($Response.Body)"
    }

    try {
        return ($Response.Body.Trim()) | ConvertFrom-Json
    } catch {
        throw "$Context returned non-JSON content: $($Response.Body)"
    }
}

function Get-SessionCookie {
    param([string]$RawHeaders)

    $match = [regex]::Match($RawHeaders, 'Set-Cookie:\s*session_key=([^;]+)', [System.Text.RegularExpressions.RegexOptions]::IgnoreCase)
    if (-not $match.Success) {
        throw "Missing session_key cookie in response headers."
    }

    return $match.Groups[1].Value.Trim()
}

function Test-LocalServiceReady {
    Write-Step "Checking web entry at $BaseUrl"
    $resp = Invoke-WebRequest -Uri $BaseUrl -UseBasicParsing
    if ($resp.StatusCode -ne 200) {
        throw "Expected $BaseUrl to return 200, got $($resp.StatusCode)"
    }
}

function Test-MySQLContainerReady {
    Write-Step "Checking Docker MySQL container '$MySQLContainer'"
    $names = docker ps --format '{{.Names}}'
    if (-not ($names -contains $MySQLContainer)) {
        throw "MySQL container '$MySQLContainer' is not running."
    }
}

function Invoke-MySQLSql {
    param(
        [Parameter(Mandatory = $true)][string]$SqlText,
        [string]$Context = "MySQL command"
    )

    $result = $SqlText | docker exec -e "MYSQL_PWD=$MySQLPassword" -i $MySQLContainer mysql "-u$MySQLUser" $MySQLDatabase 2>&1
    if ($LASTEXITCODE -ne 0) {
        throw "$Context failed: $result"
    }

    return $result
}

function Ensure-CrmSchema {
    Write-Step "Applying CRM schema if needed"
    $schemaSql = Get-Content -Raw $schemaSqlFile
    Invoke-MySQLSql -SqlText $schemaSql -Context "Apply CRM schema" | Out-Null
}

function Ensure-TestUser {
    param($UserSpec)

    $loginBody = (@{
        email    = $UserSpec.email
        password = $UserSpec.password
    } | ConvertTo-Json -Compress)

    $headers = @{ "Content-Type" = "application/json" }
    $loginResp = Invoke-CurlJson -Method "POST" -Url "$BaseUrl/api/passport/web/email/login/" -Body $loginBody -Headers $headers
    $loginJson = Convert-ResponseToJson -Response $loginResp -Context "Login user $($UserSpec.email)"

    if ($loginResp.StatusCode -ne 200 -or [int]$loginJson.code -ne 0) {
        Write-Step "Registering user $($UserSpec.email)"
        $registerResp = Invoke-CurlJson -Method "POST" -Url "$BaseUrl/api/passport/web/email/register/v2/" -Body $loginBody -Headers $headers
        $userJson = Convert-ResponseToJson -Response $registerResp -Context "Register user $($UserSpec.email)"
        if ([int]$userJson.code -ne 0) {
            throw "Register user $($UserSpec.email) failed: $($registerResp.Body)"
        }
        $cookie = Get-SessionCookie -RawHeaders $registerResp.Headers
        $userInfo = $userJson.data
    } else {
        Write-Step "Logging in with existing user $($UserSpec.email)"
        $userJson = $loginJson
        $cookie = Get-SessionCookie -RawHeaders $loginResp.Headers
        $userInfo = $userJson.data
    }

    $spaceResp = Invoke-CurlJson -Method "POST" -Url "$BaseUrl/api/playground_api/space/list" -Body "{}" -Headers @{
        "Content-Type" = "application/json"
        "Cookie"       = "session_key=$cookie"
    }
    $spaceJson = Convert-ResponseToJson -Response $spaceResp -Context "Fetch space list for $($UserSpec.email)"
    if ([int]$spaceJson.code -ne 0) {
        throw "Fetch space list for $($UserSpec.email) failed: $($spaceResp.Body)"
    }

    $spaceId = $spaceJson.data.bot_space_list[0].id
    if ([string]::IsNullOrWhiteSpace($spaceId)) {
        throw "No space_id returned for $($UserSpec.email)"
    }

    return [pscustomobject][ordered]@{
        key              = $UserSpec.key
        email            = $UserSpec.email
        password         = $UserSpec.password
        role             = $UserSpec.role
        description      = $UserSpec.description
        user_id          = $userInfo.user_id_str
        user_name        = $userInfo.name
        user_unique_name = $userInfo.user_unique_name
        space_id         = $spaceId
        session_cookie   = $cookie
        crm_url          = "$BaseUrl/space/$spaceId/crm"
    }
}

function Apply-CrmSeedIdOffset {
    param(
        [Parameter(Mandatory = $true)][string]$SqlText,
        [Parameter(Mandatory = $true)][string]$SpaceId
    )

    $suffix = if ($SpaceId.Length -gt 6) {
        $SpaceId.Substring($SpaceId.Length - 6)
    } else {
        $SpaceId
    }
    $offset = ([int64]$suffix) * 1000000

    $idGroups = @(
        (910001..910010),
        (920001..920020),
        (930001..930010),
        (940001..940020),
        (950001..950010),
        (960001..960020)
    )

    foreach ($group in $idGroups) {
        foreach ($id in $group) {
            $SqlText = $SqlText.Replace([string]$id, [string]([int64]$id + $offset))
        }
    }

    return $SqlText
}

function Seed-CrmDemoData {
    param(
        [Parameter(Mandatory = $true)][string]$SpaceId,
        [Parameter(Mandatory = $true)][string]$TenantId
    )

    Write-Step "Seeding CRM demo data for tenant_id=$TenantId space_id=$SpaceId"
    $demoSql = Get-Content -Raw $demoSqlTemplateFile
    $demoSql = $demoSql.Replace("__CRM_TENANT_ID__", $TenantId).Replace("__CRM_SPACE_ID__", $SpaceId)
    $demoSql = Apply-CrmSeedIdOffset -SqlText $demoSql -SpaceId $SpaceId
    Invoke-MySQLSql -SqlText $demoSql -Context "Seed CRM demo data for space $SpaceId" | Out-Null
}

function Verify-CrmAccess {
    param($UserResource)

    Write-Step "Verifying CRM dashboard for $($UserResource.email)"
    $dashboardUrl = "{0}/api/crm/dashboard/overview?space_id={1}" -f $BaseUrl, $UserResource.space_id
    $dashboardResp = Invoke-CurlJson -Method "GET" -Url $dashboardUrl -Headers @{
        "Cookie" = "session_key=$($UserResource.session_cookie)"
    }
    $dashboardJson = Convert-ResponseToJson -Response $dashboardResp -Context "Load CRM dashboard for $($UserResource.email)"
    if ([int]$dashboardJson.code -ne 0) {
        throw "CRM dashboard verification failed for $($UserResource.email): $($dashboardResp.Body)"
    }

    if ([int64]$dashboardJson.data.customer_total -le 0) {
        throw "CRM dashboard returned no customers for $($UserResource.email)"
    }

    $customerListUrl = "{0}/api/crm/customer/list?space_id={1}&page=1&page_size=5" -f $BaseUrl, $UserResource.space_id
    $customerResp = Invoke-CurlJson -Method "GET" -Url $customerListUrl -Headers @{
        "Cookie" = "session_key=$($UserResource.session_cookie)"
    }
    $customerJson = Convert-ResponseToJson -Response $customerResp -Context "Load CRM customer list for $($UserResource.email)"
    if ([int]$customerJson.code -ne 0) {
        throw "CRM customer list verification failed for $($UserResource.email): $($customerResp.Body)"
    }
}

if ([string]::IsNullOrWhiteSpace($MySQLPassword)) {
    $MySQLPassword = Get-DefaultMySQLPassword
}

Test-LocalServiceReady
Test-MySQLContainerReady
Ensure-CrmSchema

$userResources = New-Object System.Collections.Generic.List[object]

foreach ($userSpec in $testUsers) {
    $userResource = Ensure-TestUser -UserSpec $userSpec
    Seed-CrmDemoData -SpaceId $userResource.space_id -TenantId $userResource.space_id
    Verify-CrmAccess -UserResource $userResource

    $userResources.Add([pscustomobject][ordered]@{
        key              = $userResource.key
        email            = $userResource.email
        password         = $userResource.password
        role             = $userResource.role
        description      = $userResource.description
        user_id          = $userResource.user_id
        user_name        = $userResource.user_name
        user_unique_name = $userResource.user_unique_name
        space_id         = $userResource.space_id
        crm_url          = $userResource.crm_url
        seed_scope       = [pscustomobject][ordered]@{
            tenant_id = $userResource.space_id
            space_id  = $userResource.space_id
        }
    }) | Out-Null
}

$manifest = [pscustomobject][ordered]@{
    generated_at     = (Get-Date).ToString("yyyy-MM-dd HH:mm:ss")
    base_url         = $BaseUrl
    mysql_container  = $MySQLContainer
    default_password = $testUsers[0].password
    users            = $userResources
    examples         = $exampleQuestions
    notes            = @(
        "Each test user owns a separate personal space.",
        "CRM demo data is seeded into tenant_id == space_id for each personal space.",
        "The CRM NL query backend is still MVP-stage. The current frontend may fall back to built-in mock answers for supported example intents."
    )
}

$outputDir = Split-Path -Parent $OutputFile
if (-not (Test-Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir | Out-Null
}

$manifest | ConvertTo-Json -Depth 10 | Set-Content -Path $OutputFile -Encoding UTF8

Write-Step "CRM agent test resources are ready"
Write-Host "Manifest: $OutputFile"
Write-Host ""
Write-Host "Users:"
foreach ($userResource in $userResources) {
    Write-Host "- $($userResource.email) | password=$($userResource.password) | space_id=$($userResource.space_id)"
}
