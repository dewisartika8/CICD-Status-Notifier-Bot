#!/usr/bin/env powershell

# Code Quality Check Script for CI/CD Status Notifier Bot
# This script runs comprehensive quality checks on the Go codebase

Write-Host "🔍 Running Code Quality Checks..." -ForegroundColor Cyan

# Change to backend directory
Set-Location "backend"

$success = $true

# 1. Go Module Verification
Write-Host "`n📦 Checking Go modules..." -ForegroundColor Yellow
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Go mod tidy failed" -ForegroundColor Red
    $success = $false
} else {
    Write-Host "✅ Go modules are clean" -ForegroundColor Green
}

# 2. Build Check
Write-Host "`n🏗️ Building project..." -ForegroundColor Yellow
go build ./...
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Build failed" -ForegroundColor Red
    $success = $false
} else {
    Write-Host "✅ Build successful" -ForegroundColor Green
}

# 3. Go Vet Check
Write-Host "`n🔎 Running go vet..." -ForegroundColor Yellow
go vet ./...
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Go vet found issues" -ForegroundColor Red
    $success = $false
} else {
    Write-Host "✅ Go vet passed" -ForegroundColor Green
}

# 4. Format Check
Write-Host "`n📝 Checking code formatting..." -ForegroundColor Yellow
$unformatted = go fmt ./...
if ($unformatted) {
    Write-Host "⚠️ Code was reformatted. Please commit the changes." -ForegroundColor Yellow
} else {
    Write-Host "✅ Code is properly formatted" -ForegroundColor Green
}

# 5. Test Run (if tests exist)
Write-Host "`n🧪 Running tests..." -ForegroundColor Yellow
$testFiles = Get-ChildItem -Recurse -Name "*_test.go"
if ($testFiles.Count -gt 0) {
    go test ./... -v
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Tests failed" -ForegroundColor Red
        $success = $false
    } else {
        Write-Host "✅ All tests passed" -ForegroundColor Green
    }
} else {
    Write-Host "⚠️ No test files found" -ForegroundColor Yellow
}

# Summary
Write-Host "`n📊 Quality Check Summary:" -ForegroundColor Cyan
if ($success) {
    Write-Host "✅ All quality checks passed!" -ForegroundColor Green
    Write-Host "🚀 Code is ready for deployment" -ForegroundColor Green
    exit 0
} else {
    Write-Host "❌ Some quality checks failed" -ForegroundColor Red
    Write-Host "🔧 Please fix the issues before proceeding" -ForegroundColor Yellow
    exit 1
}
