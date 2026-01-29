# 数据库文件位置

## 配置

数据库文件默认放在项目根目录下：

```
promptManager/
├── promptmaster.db          <-- 数据库文件
├── main.go
├── backend/
├── frontend/
└── ...
```

## 修改位置

如需修改数据库位置，编辑 `backend/config/config.go`：

```go
func init() {
    // 修改这里来改变数据库位置
    DBPath = "./promptmaster.db"      // 当前目录
    // DBPath = "./data/promptmaster.db"  // 子目录
    // DBPath = "/path/to/db.sqlite"      // 绝对路径
}
```

## 备份

备份目录仍位于用户主目录：
- Windows: `%USERPROFILE%\.promptmaster\backups\`
- macOS/Linux: `~/.promptmaster/backups/`

## 数据迁移

如需将旧数据迁移到新位置：

```bash
# Windows
Copy-Item "$env:USERPROFILE\.promptmaster\promptmaster.db" "C:\Users\59634\GolandProjects\promptManager\promptmaster.db"

# macOS/Linux
cp ~/.promptmaster/promptmaster.db /path/to/project/promptmaster.db
```
