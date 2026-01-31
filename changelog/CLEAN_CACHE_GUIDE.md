# Wails 缓存清理指南

> 解决构建问题、热重载失效、代码修改不生效等缓存相关问题

---

## 🚨 常见问题

遇到以下情况时需要清理缓存：

| 问题现象 | 原因 |
|---------|------|
| `index.html: file does not exist` | `frontend/dist` 被删除但未重建 |
| 前端代码修改不生效 | Vite 缓存未刷新 |
| Go 代码修改不生效 | Go 编译缓存未清理 |
| 绑定方法找不到 | Wails 绑定未重新生成 |
| 奇怪的编译错误 | 各类缓存混合导致 |

---

## 🧹 清理方式

### 方式一：标准清理（推荐）

保留 `frontend/dist`，适合大多数情况：

```powershell
# 1. 停止 Wails 进程
taskkill /F /IM "wails.exe" 2>$null
taskkill /F /IM "PromptMaster-dev.exe" 2>$null

# 2. 清理缓存（保留 dist）
Remove-Item -Path "frontend/node_modules/.vite" -Recurse -Force
Remove-Item -Path "build/bin/*.exe" -Force
go clean -cache

# 3. 重新启动
cd D:\Project\ProjectGo\src\promptManager
wails dev
```

### 方式二：深度清理

完全清理，需要重新构建前端：

```powershell
# 1. 停止 Wails 进程
taskkill /F /IM "wails.exe" 2>$null
taskkill /F /IM "PromptMaster-dev.exe" 2>$null

# 2. 深度清理
cd D:\Project\ProjectGo\src\promptManager
Remove-Item -Path "frontend/dist" -Recurse -Force
Remove-Item -Path "frontend/node_modules/.vite" -Recurse -Force
Remove-Item -Path "build/bin/*" -Recurse -Force
go clean -cache

# ⚠️ 关键步骤：重新构建前端
# 如果没有 dist/index.html，Wails 会报错！
cd frontend
npm run build

# 3. 返回项目根目录并启动
cd ..
wails dev
```

---

## 📁 缓存位置说明

```
project/
├── frontend/
│   ├── dist/              ← 前端构建输出 ⭐ 重要！
│   │   └── index.html     ← Wails 入口文件
│   ├── node_modules/
│   │   └── .vite/         ← Vite 开发缓存
│   └── src/
│       └── lib/
│           └── wailsjs/   ← Go 绑定生成代码
├── build/
│   └── bin/               ← 可执行文件输出
└── go-build/              ← Go 编译缓存（系统临时目录）
```

| 路径 | 作用 | 能否删除 |
|------|------|---------|
| `frontend/dist` | 前端生产构建 | ⚠️ 删除后必须 `npm run build` |
| `frontend/node_modules/.vite` | Vite 开发缓存 | ✅ 可安全删除 |
| `frontend/src/lib/wailsjs` | Go 绑定代码 | ✅ 会自动重新生成 |
| `build/bin/` | 应用可执行文件 | ✅ 可安全删除 |
| `go-build/` | Go 编译缓存 | ✅ 执行 `go clean -cache` |

---

## 🔧 一键清理脚本

保存为 `clean-cache.ps1` 在项目根目录：

```powershell
# clean-cache.ps1
param(
    [switch]$Deep  # 深度清理模式
)

Write-Host "🧹 Wails 缓存清理工具" -ForegroundColor Cyan
Write-Host "======================" -ForegroundColor Cyan

# 停止进程
Write-Host "`n1. 停止 Wails 进程..." -ForegroundColor Yellow
taskkill /F /IM "wails.exe" 2>$null | Out-Null
taskkill /F /IM "PromptMaster*.exe" 2>$null | Out-Null
Start-Sleep -Seconds 1
Write-Host "   ✅ 进程已停止" -ForegroundColor Green

# 清理缓存
Write-Host "`n2. 清理缓存文件..." -ForegroundColor Yellow

# 标准清理
Remove-Item -Path "frontend/node_modules/.vite" -Recurse -Force -ErrorAction SilentlyContinue
Remove-Item -Path "build/bin/*" -Recurse -Force -ErrorAction SilentlyContinue
Write-Host "   ✅ Vite 缓存已清理" -ForegroundColor Green
Write-Host "   ✅ 构建输出已清理" -ForegroundColor Green

# Go 缓存
go clean -cache 2>$null | Out-Null
go clean -testcache 2>$null | Out-Null
Write-Host "   ✅ Go 缓存已清理" -ForegroundColor Green

# 深度清理
if ($Deep) {
    Write-Host "`n   🔥 深度清理模式..." -ForegroundColor Magenta
    
    # 删除 dist（必须重建）
    Remove-Item -Path "frontend/dist" -Recurse -Force -ErrorAction SilentlyContinue
    Write-Host "   ✅ frontend/dist 已删除" -ForegroundColor Yellow
    
    # 删除 node_modules（需要重新 npm install）
    $cleanNodeModules = Read-Host "   是否删除 node_modules？(y/N)"
    if ($cleanNodeModules -eq 'y' -or $cleanNodeModules -eq 'Y') {
        Remove-Item -Path "frontend/node_modules" -Recurse -Force -ErrorAction SilentlyContinue
        Write-Host "   ✅ node_modules 已删除" -ForegroundColor Yellow
        Write-Host "   ⚠️  需要重新执行: npm install" -ForegroundColor Red
    }
    
    # 重建前端
    Write-Host "`n3. 重新构建前端..." -ForegroundColor Yellow
    cd frontend
    npm run build
    if ($LASTEXITCODE -ne 0) {
        Write-Host "   ❌ 构建失败！请检查错误信息" -ForegroundColor Red
        exit 1
    }
    cd ..
    Write-Host "   ✅ 前端构建完成" -ForegroundColor Green
} else {
    # 标准模式下检查 dist 是否存在
    if (-not (Test-Path "frontend/dist/index.html")) {
        Write-Host "`n   ⚠️  frontend/dist 不存在，需要重新构建..." -ForegroundColor Yellow
        cd frontend
        npm run build
        cd ..
        Write-Host "   ✅ 前端构建完成" -ForegroundColor Green
    }
}

Write-Host "`n✨ 清理完成！" -ForegroundColor Cyan
Write-Host "`n可以运行以下命令启动开发服务器：" -ForegroundColor White
Write-Host "   wails dev" -ForegroundColor Cyan
Write-Host ""
```

### 使用脚本

```powershell
# 标准清理（推荐）
.\clean-cache.ps1

# 深度清理（遇到顽固问题时）
.\clean-cache.ps1 -Deep

# 然后启动 Wails
wails dev
```

---

## 📋 快速命令参考

### Windows PowerShell

```powershell
# 最快清理（只清理构建输出）
Remove-Item "build/bin/*" -Force; wails dev

# 标准清理（解决大部分问题）
Remove-Item "frontend/node_modules/.vite" -Recurse -Force; go clean -cache; wails dev

# 完整清理（深度问题）
Remove-Item "frontend/dist","frontend/node_modules/.vite" -Recurse -Force; go clean -cache; cd frontend; npm run build; cd ..; wails dev
```

### Mac/Linux Bash

```bash
# 最快清理
rm -rf build/bin/* && wails dev

# 标准清理
rm -rf frontend/node_modules/.vite && go clean -cache && wails dev

# 完整清理
rm -rf frontend/dist frontend/node_modules/.vite && go clean -cache && cd frontend && npm run build && cd .. && wails dev
```

---

## ⚠️ 重要提示

### 1. 关于 `frontend/dist`

| 情况 | 处理方式 |
|------|---------|
| 删除 `dist` 后直接用 `wails dev` | ❌ 报错：`index.html: file does not exist` |
| 删除 `dist` 后先 `npm run build` | ✅ 正常 |
| 保留 `dist` 直接 `wails dev` | ✅ 正常（Vite 会接管开发服务器） |

### 2. 清理优先级

遇到问题时，按以下顺序尝试：

1. **最简单**：只重启 `wails dev`（自动刷新部分缓存）
2. **标准清理**：清理 `.vite` 和 Go 缓存
3. **深度清理**：删除 `dist` 并重新 `npm run build`
4. **终极清理**：删除 `node_modules` 重新 `npm install`

### 3. 常见错误解决

```
Error: unable to infer the AssetDir from your Assets fs.FS: index.html: file does not exist
```
**解决**：`cd frontend && npm run build`

```
vite: not found
```
**解决**：`cd frontend && npm install`

```
error: undefined reference to...
```
**解决**：`go clean -cache` 或删除 `build/bin/*`

---

## 🔄 清理后标准启动流程

```powershell
# 1. 进入项目目录
cd D:\Project\ProjectGo\src\promptManager

# 2. 清理缓存（选择一种方式）
# 方式 A：标准
Remove-Item "frontend/node_modules/.vite" -Recurse -Force
go clean -cache

# 方式 B：深度（如果方式 A 无效）
Remove-Item "frontend/dist","frontend/node_modules/.vite" -Recurse -Force
go clean -cache
cd frontend && npm run build && cd ..

# 3. 启动开发服务器
wails dev
```

---

**文档版本**: v1.0  
**最后更新**: 2026-01-29
