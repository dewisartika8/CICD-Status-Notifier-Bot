# CICD Status Notifier Bot - Test Runner
# PowerShell script for running all bot-related tests

Write-Host "🚀 Running CICD Status Notifier Bot Tests" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

# Test counters
$TotalTests = 0
$PassedTests = 0
$FailedTests = 0

# Function to run tests and track results
function Run-TestSuite {
    param(
        [string]$TestPath,
        [string]$TestName
    )
    
    Write-Host "`n📋 Running $TestName..." -ForegroundColor Yellow
    
    $result = & go test $TestPath -v
    $exitCode = $LASTEXITCODE
    
    if ($exitCode -eq 0) {
        Write-Host "✅ $TestName PASSED" -ForegroundColor Green
        $script:PassedTests++
    } else {
        Write-Host "❌ $TestName FAILED" -ForegroundColor Red
        $script:FailedTests++
    }
    $script:TotalTests++
}

# Get backend directory
$BackendDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $BackendDir

Write-Host "📁 Current directory: $(Get-Location)"
Write-Host ""

# Run all test suites
Write-Host "🧪 Starting Test Execution..."

# Unit Tests - Bot Domain
Run-TestSuite "./tests/unit/bot/domain/..." "Bot Domain Layer Tests"

# Unit Tests - Bot Service  
Run-TestSuite "./tests/unit/bot/service/..." "Bot Service Layer Tests"

# Integration Tests - Bot
Run-TestSuite "./tests/integration/bot/..." "Bot Integration Tests"

# Unit Tests - Core (existing)
Run-TestSuite "./tests/unit/..." "Core Unit Tests"

# Integration Tests - Webhook (existing)
Run-TestSuite "./tests/integration/..." "Webhook Integration Tests"

# Summary
Write-Host ""
Write-Host "📊 TEST SUMMARY" -ForegroundColor Cyan
Write-Host "===============" -ForegroundColor Cyan
Write-Host "Total Test Suites: $TotalTests"
Write-Host "Passed: $PassedTests" -ForegroundColor Green

if ($FailedTests -gt 0) {
    Write-Host "Failed: $FailedTests" -ForegroundColor Red
    Write-Host ""
    Write-Host "❌ Some tests failed. Please check the output above." -ForegroundColor Red
    exit 1
} else {
    Write-Host "Failed: $FailedTests" -ForegroundColor Green
    Write-Host ""
    Write-Host "🎉 All tests passed successfully!" -ForegroundColor Green
    Write-Host ""
    Write-Host "✨ Bot implementation is ready for deployment!" -ForegroundColor Cyan
}
}
