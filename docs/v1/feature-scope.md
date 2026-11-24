# g 1.0 功能取舍表（草稿）

状态：范围提案，尚未整体批准。更新日期：2026-11-23。

代码阅读基线：`master@94ac309bbce6beca8628b933a8a0127d2a896ab5`，代码版本 `0.31.2`，`go.mod` 为 Go `1.26.1`。本表区分现有事实和 1.0 建议；“核心候选”不代表已经实现，“暂缓”也不授权删除旧功能。

## 已确认的产品约束

2026-11-23 维护者确认：1.0 机器输出仅保留 JSON（F11）。CSV、TSV、Markdown、NUL 路径输出在 1.0 退役，JSONL 排除在范围外。1.0 同时退役重复文件检测（F22）、持久化路径索引与模糊路径匹配（F23）。JSON 格式与兼容要求见 D06。

**1.0 不引入 preset、profile、命名场景模式、按目录类型自动推断场景，或借其他名字实现同等功能。** 普通结构化默认配置、自定义字段、独立参数、主题继续属于可讨论范围。`-l` 等现有便利参数可以规范化为明确字段和布局选项，不扩展成用户管理的命名预设系统。

建议定位：快速、可组合的目录查看器。日常浏览是主路径；用户明确请求的元数据和机器输出是扩展路径。

## 状态的含义

| 标记 | 含义 | Agent 应如何使用 |
| --- | --- | --- |
| 核心候选 | 建议进入 1.0 最小完整版本 | 在当前任务涉及该能力时实施，按行为契约验收 |
| 按需保留候选 | 建议保留，但只在请求时采集和初始化 | 拆成数据采集能力，不能增加普通浏览的隐含工作 |
| 后续候选 | 有价值的新能力，优先级低于核心迁移 | 不因架构预留顺手实现 |
| 范围待决 | 涉及现有功能去留或明显兼容变化 | 先记录影响；没有范围决定时保留现有代码/行为，不擅自退役 |
| 明确排除 | 用户明确排除 | 不实现，不预埋配置键、命令或自动推断逻辑 |

产品范围改变时，更新本表的状态、理由及关联任务/PR。实现中的常规技术选择可以由 Agent 在已授权任务内决定并记录；本表不要求每一项改动都重新请求确认。

## 功能取舍表

| ID | 能力 | 当前事实 | 1.0 建议 | 最小验收与边界 |
| --- | --- | --- | --- | --- |
| F01 | 路径输入与基本列表 | 已有无参数当前目录、多路径、文件输入、`-d`、`-1` [S1] [S2] | 普通路径输入、`-d` 和多根顺序已确认，见 C01 | 普通路径、`--`、以 `-` 开头的文件名、多根结果都有确定行为；保留位置参数为路径 |
| F02 | 普通终端列表布局 | 已有纵向网格、单列、长列表、横向网格、逗号列表和表格 [S1] | 保留纵向网格、单列、长列表，退役 across、comma、table | `--format` 只接受 `grid`、`lines`、`long`、`tree`、`json`。默认布局、快捷选项和取值校验见 [C03](behavior-contract.md#c03--终端与非终端的默认输出)。tree 单列于 F03 |
| F03 | tree、递归、深度控制 | 已有 `-T`、`-R`、`--depth` [S1] | 多根递归、深度取值、模式互斥与链接跟随已确认，见 C08。`--format=tree` 等价于 `-T/--tree`，旧值 `--format=T` 退役，暂时移除 `--tree-style` | 每个根从第 0 层计深度；数量按目录限制，祖先上下文和截断标记不占名额 |
| F04 | 大小、时间、权限、owner/group 等基础字段 | 已有 [S1] [S3] | 按 F05 总表保留；`-l` 默认六列及顺序见 [C03](behavior-contract.md#c03--终端与非终端的默认输出) | 原始值保留类型；时间概念准确；平台不适用字段可表达为不可用 |
| F05 | 自定义字段及顺序 | 当前内容由各 flag 拼装；自定义顺序是已有需求 [S3] [S4] | 2026-11-24 确认支持 `--fields` | 长列表按指定列顺序展示；JSON 保留条目身份并按需输出元数据；未知字段报参数错误。见 [C04](behavior-contract.md#c04--数据与字段模型) |
| F06 | 隐藏项、glob、扩展名、时间等筛选 | 已有 [S1] | 筛选入口及组合规则已确认，见 C09 | 筛选只判断数据；目录剪枝独立；无效条件报参数错误；同组合结果可复现 |
| F07 | 条目排序 | 已有名称、版本、MIME、大小、时间等排序 [S1] [S5] | 2026-11-24 确认退役完整 MIME、MIME 父类型和版本号排序；14 个排序键、`none` 和 5 个快捷键已确认，见 C10 | 比较器不做 I/O；昂贵排序键提前采集；同值排序有稳定后备规则；记录旧写法映射 |
| F08 | Git 文件状态、仓库分支与状态、提交详情 | 已有 [S1] [S6] | 按 F05 总表保留六个 Git 字段 | 单次运行按仓库批量获取；正确处理特殊路径；未安装 Git 和非仓库目录有不同结果 |
| F09 | Git ignore 过滤 | 帮助列出 `--git-ignore`，基线读取该值但未注册同名 flag [S1] | 显式启用，点名路径优先，已确认 | 列出显式目录的子项时按正常可见性规则处理 |
| F10 | 主题、颜色、图标、超链接 | 已有 [S1] | 保留自定义主题；颜色与超链接默认 `auto`，图标默认关闭，可经参数或配置开启。主题文件与覆盖规则见 C02 提案 | 默认样式见 C03；统一终端测宽与转义，机器字段独立于主题 |
| F11 | JSON | 已有，当前从展示字符串生成 `string` 值 [S7] | 确认在 1.0 直接切换到版本化格式；字段设计见 [C05](behavior-contract.md#c05--json) | 空结果、多根、部分失败仍为合法单文档；数值与时间保留语义；字段名独立于表头 |
| F12 | JSONL | 当前已列出的输出模式中未见 JSONL [S1] | 明确排除 | 机器输出采用 JSON |
| F13 | CSV、TSV、Markdown | 已有 [S1] [S7] | 1.0 确认退役 | 迁移说明列出三个导出格式的旧入口 |
| F14 | NUL/原始路径输出 | 已有 `--zero/-0` 文档；当前 `Zero.Print` 未写 NUL [S1] [S7] | 1.0 确认退役 | 迁移说明列出 `--zero/-0` |
| F15 | stdin 输入与 stdin0 | 已有按换行读取的 stdin；stdin0 为新增提案 [S1] [S2] | 输入范围待决 | 分别评估按行输入与 stdin0 的使用需求 |
| F16 | 默认配置、配置校验与迁移 | 已有 `Args`、主题和树样式配置 [S8] | 已确认结构化配置覆盖全部 CLI 运行选项与输入路径；字段设计见 C02 | `defaults < user config < explicit CLI`；显式 false/0/空列表可覆盖；旧配置在入口转换；不包含命名预设 |
| F17 | help、man、补全、shell init | 已有 [S1] [S9] | 核心候选 | 命令定义为单一来源；各 shell 的支持级别明确；生成内容能检查漂移 |
| F18 | 跨平台分发、安装入口、构建能力 | 已有 Linux/macOS/Windows 分发；普通构建与 release 的 tags 有差异 [S10] [S11] | 核心候选 | 保留根模块安装入口；正式分发能力集合明确；未编入的显式能力不能静默成功；长命令名 `g-ls` 为待评估交付项 |
| F19 | MIME、charset | 已有 [S1] | 保留 `mime`、`charset`，退役 MIME 父类型列 | 只对相关候选条目采集；避免重复读头；失败作为数据诊断，不写进正常字段字符串 |
| F20 | checksum、递归目录大小 | 已有 [S1] [S5] [S12] | 保留，按需采集 | checksum 流式读取；只哈希最终展示项（除非用于筛选排序）；目录大小有链接、硬链接、错误与深度规则 |
| F21 | mounts、扩展属性和平台特有字段 | 已有部分能力，mounts 受 build tag 控制 [S1] [S10] | 按 F05 总表保留，按需采集 | 平台适配单独组织；请求字段的 unsupported/failed 可区分；能力清单反映实际编译结果 |
| F22 | 重复文件检测 | 已有 `--duplicate` [S1] | 1.0 确认退役 | 迁移说明列出 `--duplicate` 的变化 |
| F23 | 持久化路径索引、模糊路径匹配 | 已有，受 fuzzy tag 控制 [S1] [S10] | 1.0 确认退役 | 迁移说明列出索引与模糊匹配入口、构建选项的变化 |
| F24 | preset/profile/场景模式 | 上一轮方案提出，用户已明确不想要 | **明确排除** | CLI、配置、文档、测试和架构中均没有命名预设或同等替代设计 |

## F05 可选字段总表

2026-11-24 维护者确认，1.0 当前保留下列 31 个 `--fields` 字段。已确认的 JSON 键名、类型、语义与示例见 [JSON 字段字典](json-fields.md)。`-l` 默认使用表中第一行的六列。`mime_parent` 退出字段范围，对应的旧参数为 `--mime-parent`、`--mime-parent-type`、`--mimetype-parent`。

| 能力 | CLI 字段 | 数据含义与 0.x 对应入口 | 使用方式 |
| --- | --- | --- | --- |
| 长列表默认六列 | `mode`、`owner`、`group`、`size`、`mtime`、`name` | 权限、所有者名、组名、条目自身的逻辑大小、最后修改时间、名称。`size` 对应 JSON `size_bytes` | 保留 |
| 路径与类型 | `path`、`kind`、`link_target` | 操作路径、条目类型、符号链接中保存的原始目标。旧版把目标随名称展示，1.0 可独立选择 `link_target` | 按需 |
| 其他时间 | `atime`、`ctime`、`birthtime` | 访问时间、元数据变更时间、创建时间。旧 `--create` 在 Linux/macOS 读 ctime，在 Windows 读 creation time。按实际语义映射 | 保留，按平台能力提供 |
| 归属与文件身份 | `owner_id`、`group_id`、`inode`、`nlink` | 前两项对应旧 `--uid/--gid`，Unix 使用 UID/GID，Windows 使用 SID。后两项为 inode 号、硬链接数（旧 `--link`） | 保留，按平台能力提供 |
| 空间占用 | `allocated_size`、`recursive_size` | 实际分配的磁盘空间、递归目录总量。旧 `--blocks` 在 Linux/macOS 读 `stat.Blocks`；`allocated_size` 以字节计。`recursive_size` 需要遍历子树 | 按需 |
| 类型与编码 | `mime`、`charset` | 完整 MIME 类型与字符编码，例如 `text/plain` 与 `UTF-8` | 按需 |
| 内容摘要 | `checksum` | 旧版可选 MD5、SHA-1、SHA-224、SHA-256、SHA-384、SHA-512、CRC32。1.0 默认 SHA-256，按算法标识返回摘要 | 按需 |
| Git 工作区与仓库 | `git_status`、`git_branch`、`git_repo_status` | 条目的暂存区与工作区状态、仓库分支、仓库整体状态。旧入口 `--git`、`--git-repo-branch`、`--git-repo-status` | 按需 |
| Git 文件历史 | `git_commit`、`git_author`、`git_author_time` | 最近涉及该路径的提交哈希、作者和作者时间；对应旧 `--git-detail` | 按需，独立于 Git 状态 |
| 系统属性 | `flags`、`has_xattrs`、`mount_device`、`mount_fstype` | macOS 文件标志、有无扩展属性、挂载点的设备与文件系统类型。旧版分别通过 `--flags`、权限后的 `@`、`--mounts` 提供 | 按需，按平台能力提供 |

**字段以外的既有输出。** 八进制权限、相对时间、大小单位、绝对或相对路径、引号、图标、超链接属于格式选项。行号 `#` 按结果顺序生成；`--total-size`、统计和页脚属于整份列表。JSON 的 `root_id` 与 `_b64` 伴随键按 C05 生成。duplicate、持久化索引和 fuzzy 的退役决定见 F22/F23。

源码基线：[CLI 字段入口][S3]、[Git 列](https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/content/git.go)、[权限与扩展属性](https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/content/permission.go)、[挂载点信息](https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/content/mounts.go)。

## 功能之外的 1.0 验收目标

- 一次执行的 Request、缓存、输出和错误状态独立；同进程连续执行互不污染。
- 默认浏览的工作量可解释；不会因为未请求的字段读取内容、递归统计或更新索引。
- 输出失败能传播到退出状态；部分读取失败仍能输出成功条目。
- 用固定 fixture 和相同能力集合比较旧、新版本；性能目标写成实测基线，避免引用代码注释中的百分比作为证据。
- 新增字段主要修改字段定义、采集器和对应验收，不穿透每一种布局。
- 用户可观察的兼容变化有迁移条目；不以“仍是 0.x”为理由默认取消全部兼容。

## 建议的首轮范围

首个可运行切片：F01 + F11 的最小字段集，随后加入 F04/F05/F06/F07。F03/F08/F09 在数据边界稳定后迁移。F15 输入范围单独讨论。

F19–F21 按 F05 总表迁移，采集随字段、筛选和排序需求启用。

尚未决定的具体政策列在 [行为契约的未决项](behavior-contract.md#未决项与推荐默认)，实施顺序见 [路线图](roadmap.md)。

## 资料与事实来源

下列源码链接固定到阅读基线；Issue 仅证明需求/反馈存在，不证明问题已复现或已批准。

[S1]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/docs/man.md
[S2]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/main.go
[S3]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/cli/view.go
[S4]: https://github.com/equationzhao/g/issues/91
[S5]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/sorter/sort.go
[S6]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/git/git_status.go
[S7]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/display/printer.go
[S8]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/config/load.go
[S9]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/shell/integration.go
[S10]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/docs/BuildOption.md
[S11]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/justfile
[S12]: https://github.com/equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/content/sum.go
