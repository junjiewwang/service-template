package testutil

// ============================================
// 生成器特定的预设配置
// ============================================
// 这个文件只包含生成器测试专用的预设配置
// 通用的预设配置请使用 pkg/config/testutil 包

// 注意：所有通用预设已迁移到 pkg/config/testutil
// 请使用以下方式访问：
//   - testutil.MinimalConfig()
//   - testutil.GoServiceConfig()
//   - testutil.PythonServiceConfig()
//   - testutil.JavaServiceConfig()
//   - testutil.ConfigWithPlugins()
//   - testutil.ConfigWithCustomHealthcheck()
//   - testutil.ConfigWithMultiArchPlugin()

// 如果需要生成器特定的预设配置，可以在这里添加
// 示例：
// func NewGeneratorSpecificPreset() *config.ServiceConfig {
//     return configtestutil.NewConfigBuilder().
//         WithService("generator-specific", "Generator Specific Service").
//         // ... 生成器特定的配置
//         BuildWithDefaults()
// }
