# 重构基线观察

记录日期：2026-11-23。参考提交：`94ac309bbce6beca8628b933a8a0127d2a896ab5`。观察方式：读取源码、测试、文档与 workflow；本资料包不附带该版本的构建通过或性能实测结论。

此文件用于定位。实施前与当前分支核对，已修复的条目更新为历史记录。用户反馈与代码事实分开使用。

## 1. 结构与执行路径

| ID | 观察到的事实 | 直接改造切口 | 来源 |
| --- | --- | --- | --- |
| B01 | CLI 的 renderer、printer、sorter、enabler、hooks、返回码等处于包级状态 | Run 生命周期局部化 | [cli/g.go][cli] |
| B02 | flag Action 会替换 printer、改变主题等 | 输入解析与 Resolve/Plan 分离 | [cli/display.go][display-flags] |
| B03 | FileInfo 同时有 os.FileInfo、展示 Meta 和 byte cache | typed Entry 与显示 Cell 分离 | [item/fileinfo.go][item] |
| B04 | JSON 由 `OrderedMap[string,string]` 编码展示字符串 | typed schema 和专门机器输出 | [printer.go][printer] |
| B05 | Zero.Print 写 OrderedContent，没有在该循环写 NUL 分隔字节 | 按 F14 记录旧参数的退役变化 | [printer.go][printer] |
| B06 | MIME 与递归大小排序比较器中会采集数据和写缓存 | 先采集排序键，再纯比较 | [sorter/sort.go][sorter] |
| B07 | 自适应选择先 ReadDir 估计数量，处理器再次 ReadDir；两者逐项 Info | 一次枚举、按需属性查询，重新建立基准 | [adaptive][adaptive]、[processors][processors] |
| B08 | checksum 读完整文件并保留在 Entry 的内容缓存 | 流式计算和有界内存 | [sum.go][sum] |
| B09 | Dive 的 Match=false 会在考虑进入目录前 continue | 展示匹配与遍历剪枝分离 | [dive.go][dive]、[filter][filter] |
| B10 | 配置是 Args 列表并写全局 Default；入口插入/重排 os.Args | typed config 与兼容转换边界 | [config][config]、[main.go][main] |
| B11 | 多输入路径的主循环直接打印文本标题，再对每根调用 printer | 多根由统一输出生命周期处理 | [cli/g.go][cli] |

上表为静态代码依据；更改后用对应行为 fixture 验证。它没有给出哪一处占总耗时多少或能够提速多少的测量结论。

## 2. 构建、测试和交付

| 事实 | 后续事项 | 来源 |
| --- | --- | --- |
| go.mod 指定 Go 1.26.1；版本文件为 0.31.2 | 使用当前分支版本来源，避免多处复制 | [go.mod][gomod]、[version.go][version] |
| docs/man.md 的版本仍写 0.30.0 | 帮助、man、版本同源生成 | [man][man] |
| 普通构建没有 fuzzy/mounts tags，发布 build recipe 带这两个 tag | 官方能力集合待定并统一 | [BuildOption][buildoption]、[justfile][just] |
| Darwin alias 实现导入 C，平台发布使用 CGO=1 | 保留真实平台依赖或单独改变实现 | [Darwin][darwin]、[justfile][just] |
| 旧 run_test.sh 在开头 exit 0，just test 仍调用 | 建立真实黑盒验证入口 | [脚本][testscript]、[justfile][just] |
| just test 禁内联，普通 CI go test 没有同样参数 | 记录 gomonkey 的验证条件，逐步减少耦合 | [justfile][just]、[CI][ci] |
| Windows CI job 构建；平台行为测试覆盖不能由它推导 | 平台验证分类记录 | [CI][ci] |
| 发布脚本有 prerelease 布尔标记，版本计算使用 X.Y.Z | beta 需要真实后缀和对应解析 | [release][release]、[publish][publish] |

## 3. 用户需求线索

- [#91](https://github.com/equationzhao/g/issues/91)：自定义列及列顺序。
- [#221](https://github.com/equationzhao/g/issues/221)：短命令 `g` 与 shell 插件冲突。

Issue 是需求线索，状态可能变化。实施任务核对当前 issue 与代码；不能仅凭旧报告断言某个问题仍然存在。

## 4. 本次用户约束

功能范围和决策状态汇总在[功能取舍表](feature-scope.md#已确认的产品约束)。

[cli]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/cli/g.go
[display-flags]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/cli/display.go
[item]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/item/fileinfo.go
[printer]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/display/printer.go
[sorter]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/sorter/sort.go
[adaptive]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/util/adaptive_strategy.go
[processors]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/util/processors.go
[sum]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/content/sum.go
[dive]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/cli/dive.go
[filter]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/filter/itemfliter.go
[config]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/config/load.go
[main]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/main.go
[gomod]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/go.mod
[version]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/cli/version.go
[man]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/docs/man.md
[buildoption]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/docs/BuildOption.md
[just]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/justfile
[darwin]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/osbased/filedetail_darwin.go
[testscript]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/script/run_test.sh
[ci]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/.github/workflows/go.yml
[release]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/.github/workflows/release.yml
[publish]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/.github/workflows/publish-release.yml
