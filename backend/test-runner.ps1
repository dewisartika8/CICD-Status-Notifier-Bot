# CICD Status Notifier Bot - Test Runner
Write-Host "Running CICD Status Notifier Bot Tests" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan

$TotalTests = 0
$PassedTests = 0
$FailedTests = 0

function Run-TestSuite {
    param([string]$TestPath, [string]$TestName)
    
    Write-Host "Running $TestName..." -ForegroundColor Yellow
    
    $result = & go test $TestPath -v
    if ($LASTEXITCODE -eq 0) {
        Write-Host "$TestName PASSED" -ForegroundColor Green
        $script:PassedTests++
    } else {
        Write-Host "$TestName FAILED" -ForegroundColor Red
        $script:FailedTests++
    }
    $script:TotalTests++
}

Set-Location (Split-Path -Parent $MyInvocation.MyCommand.Path)
Write-Host "Current directory: $(Get-Location)"

Write-Host "Starting Test Execution..."

Run-TestSuite "./tests/unit/bot/domain/..." "Bot Domain Layer Tests"
Run-TestSuite "./tests/unit/bot/service/..." "Bot Service Layer Tests"
Run-TestSuite "./tests/integration/bot/..." "Bot Integration Tests"
Run-TestSuite "./tests/unit/..." "Core Unit Tests"
Run-TestSuite "./tests/integration/..." "Webhook Integration Tests"

Write-Host ""
Write-Host "TEST SUMMARY" -ForegroundColor Cyan
Write-Host "Total Test Suites: $TotalTests"
Write-Host "Passed: $PassedTests" -ForegroundColor Green
Write-Host "Failed: $FailedTests" -ForegroundColor $(if ($FailedTests -gt 0) { "Red" } else { "Green" })

if ($FailedTests -gt 0) {
    Write-Host "Some tests failed." -ForegroundColor Red
    exit 1
} else {
    Write-Host "All tests passed successfully!" -ForegroundColor Green
}
