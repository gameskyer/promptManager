# 后端测试指南

## 测试文件说明

- `db_test.go` - 数据库和服务层测试

## 运行测试

### 运行所有测试

```bash
cd backend
go test ./tests/... -v
```

### 运行特定测试

```bash
# 运行分类服务测试
go test ./tests/... -v -run TestCategoryService

# 运行原子词服务测试
go test ./tests/... -v -run TestAtomService

# 运行预设服务测试
go test ./tests/... -v -run TestPresetService

# 运行版本服务测试
go test ./tests/... -v -run TestVersionService

# 运行数据导入测试
go test ./tests/... -v -run TestSeeder

# 运行集成测试
go test ./tests/... -v -run TestIntegration
```

### 运行基准测试

```bash
go test ./tests/... -bench=. -benchmem
```

## 测试覆盖范围

### CategoryService 测试
- ✅ CreateCategory - 创建分类
- ✅ GetCategoryByID - 根据ID获取分类
- ✅ GetCategoriesByParent - 获取子分类
- ✅ UpdateCategory - 更新分类
- ✅ DeleteCategory - 删除分类

### AtomService 测试
- ✅ CreateAtom - 创建原子词
- ✅ GetAtomByID - 根据ID获取原子词
- ✅ GetAtomsByCategory - 根据分类获取原子词
- ✅ FindAtomsBySynonym - 根据同义词搜索
- ✅ UpdateAtom - 更新原子词
- ✅ RecordUsage - 记录使用情况
- ✅ DeleteAtom - 删除原子词
- ✅ BatchImportAtoms - 批量导入原子词

### PresetService 测试
- ✅ CreatePreset - 创建预设（自动创建V1）
- ✅ GetPresetByID - 根据ID获取预设
- ✅ GetPresets - 获取预设列表（支持分页）
- ✅ SoftDeletePreset - 软删除预设
- ✅ RestorePreset - 恢复软删除的预设
- ✅ BuildPromptText - 构建提示词文本

### VersionService 测试
- ✅ CreateVersion - 创建新版本
- ✅ GetVersion - 获取特定版本
- ✅ GetVersionHistory - 获取版本历史
- ✅ StarVersion - 星标版本
- ✅ CompareVersions - 版本对比

### Seeder 测试
- ✅ SeedAll - 导入默认数据
- ✅ SeedFromJSON - 从JSON导入数据
- ✅ GetDefaultSeedData - 获取默认种子数据

### 集成测试
- ✅ EndToEndWorkflow - 端到端工作流测试

## 测试数据

测试使用内存中的SQLite数据库，每个测试独立运行，测试完成后自动清理。

## 添加新测试

1. 在 `db_test.go` 中添加测试函数
2. 使用 `setupTestDB(t)` 创建测试数据库
3. 运行服务方法并验证结果
4. 使用 `t.Errorf()` 报告错误

示例：

```go
func TestNewFeature(t *testing.T) {
    db := setupTestDB(t)
    service := services.NewXXXService(db)
    
    // 测试代码
    result, err := service.DoSomething()
    if err != nil {
        t.Errorf("Failed: %v", err)
    }
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```
