# 1.0 用户配置与主题 YAML

2026-11-24 字段提案。已确认保留自定义主题，旧主题迁移独立于 `g` 可执行程序。主题转换工具另立任务。下文定义新文件的用户接口，行为基础见 [C02 与 C03](behavior-contract.md#c02--配置合并与显式覆盖)。

## 用户配置 `g.yaml`

2026-11-24 已确认：所有 1.0 CLI 运行选项和默认输入路径均可在 `g.yaml` 表达。每个选项使用一个规范键，CLI 别名和快捷写法映射到同一设置，例如 `-T` 对应 `format: tree`、`-S` 对应 `sort: ["size:desc"]`。以下键名为提案，配置只需写出要覆盖的值，例如：

```yaml
version: 1
icons: true
theme: themes/local.yaml
long:
  fields: [mode, size, mtime, name]
```

| 键 | 类型和取值 | 省略时的行为 |
| --- | --- | --- |
| `version` | 整数 `1`，新主配置必填 | 缺少配置文件时使用内置设置；已有配置文件通过版本识别语法 |
| `format` | `grid`、`lines`、`long`、`tree`、`json` | 交互终端用 grid，管道或重定向逐行输出；省略表达自动选择 |
| `color`、`hyperlink` | 各为 `auto`、`always` 或 `never` | 各自 `auto`，按输出终端能力启用 |
| `icons` | 布尔值 | `false` |
| `theme` | 主题文件路径字符串 | 使用内置主题。相对路径从 `g.yaml` 所在目录解析；命令行 `--theme` 的相对路径从调用目录解析 |
| `long.fields` | 非空、有序且无重复的字段名列表，使用 F05 的 31 个 CLI 字段名 | `mode, owner, group, size, mtime, name`；只对长列表生效 |
| `sort` | 非空字符串列表，如 `["mtime:desc", "name"]`；语法同 `--sort` | 默认目录优先、自然名称升序；显式排序时按列表顺序比较。命令行排序键整体替换配置中的列表 |
| `dir_first` | 布尔值 | 沿用 C10 的分组规则：默认排序目录优先，显式排序跨类型比较。显式 true/false 覆盖分组偏好 |
| `visibility` | `visible`、`all` 或 `hidden` | `visible` |
| `git_ignore` | 布尔值 | `false` |

| 其他规范键（提案） | 取值与 CLI 对应行为 |
| --- | --- |
| `paths` | 有序路径列表，对应位置参数。相对条目路径以调用目录为基准，显式 CLI 路径整体替换默认路径。 |
| `fields` | 有序字段列表，作用同 `--fields`。`format: json` 时选择 JSON 元数据；单独指定时选用长列表。`long.fields` 只设置长列表自身的默认列。 |
| `match`、`ignore`、`ext`、`only_mime` | 字符串列表，依次对应包含 glob、排除 glob、扩展名和 MIME 过滤。 |
| `kind`、`time_type`、`before`、`after` | `kind` 使用 `directory`、`file`、`non-dir`；时间字段及前后界限使用与 CLI 相同的语义和运行时区，日期时间写为字符串。 |
| `recurse`、`directory`、`follow`、`reverse` | 布尔值，依次对应 `-R`、`-d`、`-L`、`-r`。 |
| `depth`、`limit` | 非负整数。`depth: 0` 表示只列根，`limit: 0` 表示不限数量。`depth` 要求最终启用递归或树形模式。 |
| `checksum_algorithm` | 算法名列表，供 `checksum` 字段使用，默认 `sha256`。 |

新增 CLI 能力时同步定义 YAML 键、类型与示例；仍在讨论的功能随其 CLI 契约一起补充。配置发现选项（例如 `--no-config`）在打开文件前生效。

```yaml
version: 1
paths: [src, tests]
recurse: true
depth: 2
ext: [go, rs]
ignore: ["*_test.go"]
sort: ["mtime:desc"]
limit: 10
```

建议按内置值、用户配置、显式 CLI 的顺序合并。标量取最高优先级的显式值；同一键的列表按来源整体替换，同一 CLI 内的重复参数沿用各自的组合规则。显式 `false`、`0` 和允许为空的列表均生效。上例执行 `g -n 0 lib` 时，输入根改为 `lib`、数量不限，其余设置继续生效。合并后的请求共用 C01 的互斥与前置条件校验。未知键、重复键、类型或取值错误指出文件和键路径。

## 主题文件

主题以属性为单位覆盖内置样式。最小主题可以只改一种文件扩展名：

```yaml
extension:
  go: {color: cyan}
```

| 一级键 | 可用子键 | 作用 |
| --- | --- | --- |
| `version` | 可选整数 `1`，省略按 v1 读取 | 独立分享的主题可显式标注版本 |
| `kind` | `file`、`directory`、`symlink`、`fifo`、`socket`、`block_device`、`char_device`、`other` | 各类条目名称的基础样式，复用 `kind` 的既有枚举 |
| `extension` | 小写扩展名，如 `go`、`gz`，省略前导点 | 普通文件按最后一段扩展名匹配，ASCII 大小写等价；单个前导点的名称如 `.env` 走名称规则 |
| `name` | 精确文件名，如 `README.md`、`.git` | 对路径末段逐字匹配，区分大小写，可用于文件或目录 |
| `field` | 31 个 CLI 字段名，如 `mtime`、`size`、`owner` | 对应终端列的基础样式。`size` 即 CLI 字段名，JSON 另用 `size_bytes` |
| `permission` | `read`、`write`、`execute`、`setuid`、`setgid`、`sticky`、`unset` | `mode` 列对应的权限字符 |
| `size` | `B`、`KiB`、`MiB`、`GiB`、`TiB`、`PiB`、`EiB` | 按最终显示的 1024 进位单位细化大小列样式 |
| `git` | `clean`、`changed`、`conflicted`、`untracked`、`ignored` | 按 `git_status.state` 细化状态列样式 |
| `symlink` | `arrow`、`target`、`broken_target` | 链接箭头、目标文本、断链目标文本 |

每个选择器的值都是样式对象。`color` 接受 16 种 ANSI 名称（black、red、green、yellow、blue、magenta、cyan、white 及其 `bright-` 形式）、`default`、0–255 整数或带引号的 `"#RRGGBB"`。`bold`、`italic`、`underline`、`dim`、`blink` 均为布尔值。`icon` 是字符串，作用于名称选择器（`kind`、`extension`、`name`），由主配置的图标开关控制。遗漏属性继承，`bold: false` 取消加粗，`icon: ""` 清除该图标，`color: default` 恢复终端默认前景色。

完整示例片段：

```yaml
version: 1
kind:
  directory: {color: blue, bold: true}
extension:
  go: {color: cyan}
name:
  README.md: {color: "#D7AF5F", bold: true}
field:
  mtime: {dim: true}
permission:
  write: {color: yellow}
size:
  GiB: {bold: true}
git:
  conflicted: {color: red, bold: true}
symlink:
  broken_target: {color: red, underline: true}
```

先在同一选择器内合并内置和用户属性，再按具体程度应用。名称从 `field.name` 开始，依次应用 `kind`、`extension`、`name`；字段从 `field.<字段>` 开始，再应用权限字符、大小单位、Git 状态或链接片段的细化样式。只有同一属性冲突时后者优先。未知的固定选择器报出完整路径，例如 `git.confict`；`name` 和 `extension` 允许用户自定义匹配键。

## 兼容范围待选

旧主题还有按可执行文件、隐藏文件、空目录、挂载点、root 所有者、特定分支和 Git 两字符状态细分的样式。上述 v1 提案优先复用已确认的字段和状态；是否要求这些细分项一一可表达需要单独确认。独立迁移工具可复用新主题格式定义，并列出无法等价转换的旧键。旧 `Args:` 是主配置兼容议题，单独处理。
