# `g --json` v1 字段字典

状态：2026-11-24 维护者已确认 schema v1 字段契约。字段范围见 [F05](feature-scope.md#f05-可选字段总表)，顶层对象、路径基准与缺值规则见 [C05](behavior-contract.md#c05--json)。本文定义 `entries[]` 中的键。

## 投影和取值

每条记录必带 `root_id`、`path`、`name`、`kind`。`root_id` 是指向 `roots` 的非负整数，属于协议字段。默认投影另含 `size_bytes`、`mtime`。显式 `--fields` 替换可选部分：`--fields owner,mode` 生成四个身份键加 `owner`、`mode`。JSON 对象用键名访问。

以下类型描述成功取值。省略表示未请求；已请求但无可用值时保留键并写 `null`。字段级诊断的 `field` 使用 **JSON 键名**，例如 `recursive_size_bytes`。数值 0、布尔值 false、空数组均为成功结果。所有时间遵循 [C03 统一时区规则](behavior-contract.md#c03--终端与非终端的默认输出)，按本次运行的时区编码为 RFC 3339，使用该时刻的实际偏移。UTC 环境使用 `Z`，上海示例使用 `+08:00`。保留可用的小数秒（最多九位）。字节数和计数使用精确的 JSON 非负整数；处理超大文件的消费者应使用整数解码。

## 31 个 CLI 字段到 JSON 的映射

| CLI 字段 | JSON 键 | JSON 类型 | 成功值的含义 |
| --- | --- | --- | --- |
| `mode` | `mode` | string | Unix 权限及特殊位的四位八进制文本，如 `0644`、`4755`。类型见 `kind`；平台无法提供对应语义时为 `null`。 |
| `owner` | `owner` | string | 文件所有者账号名，Windows 可含域限定。账号名解析失败时给出诊断，底层标识可用 `owner_id` 取得。 |
| `group` | `group` | string | 文件所属组名，Windows 使用文件安全描述符的主组。 |
| `size` | `size_bytes` | integer | 条目自身的逻辑字节数，使用条目自身的 stat 信息。目录记录目录项自身大小，符号链接记录链接自身大小。 |
| `mtime` | `mtime` | string | 条目最后修改时间，如 `2026-11-23T09:41:00+08:00`。 |
| `name` | `name` | string | 路径末段；文件系统根使用根路径。 |
| `path` | `path` | string | 从调用时工作目录可操作的路径。相对输入产生相对路径，绝对输入产生绝对路径，使用平台分隔符。 |
| `kind` | `kind` | string | 自身类型。枚举为 `file`、`directory`、`symlink`、`fifo`、`socket`、`block_device`、`char_device`、`other`。 |
| `link_target` | `link_target` | string | 符号链接中存储的原始目标，如 `../lib`，相对目标以链接所在目录为基准。断链仍可返回该值。 |
| `atime` | `atime` | string | 文件系统记录的最后访问时间。 |
| `ctime` | `ctime` | string | 最近一次元数据状态变化的时间，如改权限或属主。 |
| `birthtime` | `birthtime` | string | 文件系统提供的创建时间。缺少该能力时返回 `null`。 |
| `owner_id` | `owner_id` | string | Unix UID 的十进制文本，如 `501`；Windows 使用规范的 SID 字符串。 |
| `group_id` | `group_id` | string | Unix GID 的十进制文本或 Windows 主组 SID，与 `owner_id` 采用同一编码规则。 |
| `inode` | `inode` | string | 文件系统内 inode 号的完整十进制文本，无前导零（零值为 `0`）。作用域是所属文件系统。 |
| `nlink` | `nlink` | integer | 文件系统返回的硬链接数。 |
| `allocated_size` | `allocated_size_bytes` | integer | 文件系统计入本条目的已分配字节数，按平台报告的分配单位换算。稀疏文件可小于逻辑大小。 |
| `recursive_size` | `recursive_size_bytes` | integer | 目录子树中普通文件的逻辑字节数之和。空目录为 0；普通文件及其他种类对此字段自然不适用。遍历规则见下。 |
| `mime` | `mime` | string | 普通文件的 MIME 类型，使用小写的 `type/subtype`，如 `text/plain`。编码另存于 `charset`。 |
| `charset` | `charset` | string | 可识别的文本编码标签，以小写输出，如 `utf-8`、`us-ascii`。二进制或无法识别的编码返回 `null`。 |
| `checksum` | `checksum` | object | 普通文件内容摘要。键为算法名，值为定长小写十六进制文本，保留前导零。单算法和多算法均用对象。 |
| `git_status` | `git_status` | object | 当前工作树内该路径的状态。对象含 `state` 和 `xy`；正常、冲突与目录的规则见下。 |
| `git_branch` | `git_branch` | string | Git 工作树根条目的当前本地分支名，如 `main`。尚无提交但 HEAD 已指向分支时仍返回分支名；分离 HEAD 时为 `null`。 |
| `git_repo_status` | `git_repo_status` | string | Git 工作树根的整体状态，枚举 `clean`、`dirty`。暂存、未暂存、冲突、未跟踪项均计入 `dirty`；忽略项不计。 |
| `git_commit` | `git_commit` | string | 从 HEAD 查询最近涉及该路径的提交，输出完整的小写十六进制对象 ID；长度由仓库对象格式决定。 |
| `git_author` | `git_author` | string | 上述提交记录的原始作者姓名，保留原拼写。 |
| `git_author_time` | `git_author_time` | string | 上述提交的作者时间，按运行环境时区输出。三项历史字段共享同一提交。 |
| `flags` | `flags` | array of string | 平台文件标志名，去重并按字典序排列，如 `["hidden","uchg"]`。首版使用 macOS 标志名称；支持且无标志时为 `[]`。 |
| `has_xattrs` | `has_xattrs` | boolean | 当前进程能列出的扩展属性中至少有一项时为 `true`；成功列举为空时为 `false`。 |
| `mount_device` | `mount_device` | string | 条目自身为挂载点时，返回挂载表中的设备或来源标识，如 `/dev/disk3s1`；保留系统原值。 |
| `mount_fstype` | `mount_fstype` | string | 同一挂载记录的文件系统类型，如 `apfs`、`ext4`。 |

## 边界与样例

**文件与路径。** 自身属性读取链接本身，读取 `link_target` 无需目标存在。`name`、`path`、根 `input`、`link_target` 的无效 UTF-8 字节按 C05 生成同名 `_b64` 伴随键。`flags` 中已置位且没有名字的位用 `0x` 加小写十六进制单比特掩码表示，避免把未知位误报为空数组。

**目录合计。** 从所列目录向下遍历同一文件系统的普通文件，每个目录合计内按文件身份去重硬链接；每个请求合计的目录独立计数。跳过符号链接及下级挂载点，遍历包含隐藏项，列表的筛选、显示深度和 limit 只改变哪些目录出现在结果中。权限错误、身份不可得或中途读取失败使该目录的合计为 `null`，诊断关联该目录与 `recursive_size_bytes`，说明中给出失败的子路径。

**摘要。** 1.0 已确认默认使用 SHA-256。仅请求 `checksum` 时，对象只含 `sha256`。可选键为 `md5`、`sha1`、`sha224`、`sha256`、`sha384`、`sha512`、`crc32`（IEEE）。仅计算请求算法；空文件也有有效摘要。同一文件内容读取失败时整项为 `null`。例如，空文件选两种算法：

```json
{"checksum":{"sha256":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","crc32":"00000000"}}
```

**Git。** `git_status.state` 取 `clean`、`changed`、`conflicted`、`untracked`、`ignored`。文件的 `xy` 是两个字符，采用 Git porcelain v1 的 XY 编码，空格保留；干净文件补为两个空格。冲突保留原始组合（如 `UU`、`AA`），未跟踪和忽略分别为 `??`、`!!`。普通目录汇总子项时 `xy` 为 `null`：优先级依次为冲突、已跟踪变动、未跟踪；其余为干净，目录自身被忽略才为 `ignored`。整目录仅有未跟踪子项时为 `untracked`。子模块在父仓库中使用 Git 对该条目的直接状态，仓库字段则描述它自己的工作树。仓库外字段自然不适用。历史字段以 HEAD 下当前路径查询，目录查询覆盖其子树，按路径查询到的重命名提交可返回，历史查询止于改名前的路径边界。分离 HEAD 仍可读取提交历史与仓库状态。[Git 状态编码](https://git-scm.com/docs/git-status#_short_format)

以下为 `--fields mode,owner_id,inode,git_status` 的一条示例记录。它包含身份键及显式请求的字段：

```json
{
  "root_id": 0, "path": "src/main.go", "name": "main.go", "kind": "file",
  "mode": "0644", "owner_id": "501", "inode": "9007199254740993",
  "git_status": {"state": "changed", "xy": " M"}
}
```

**缺值诊断（已确认）。** 选中的字段自然不适用时仅返回 `null`。例如普通文件的 `link_target`、非挂载点的 `mount_device`、平台未提供的创建时间，以及成功读取后仍无法识别的 `charset`。字段读取或采集失败时返回 `null`，并在 `diagnostics` 记录对应路径、JSON 键和原因。例如已识别为符号链接，读取目标时遇到权限错误，诊断码为 `permission_denied`。

## 2026-11-24 已确认的决定

| 决定 | 约定 |
| --- | --- |
| 权限与标识符的类型 | `mode` 固定四位八进制文本；`owner_id`、`group_id`、`inode` 都用字符串，统一 Unix/Windows 身份表示并允许精确比对大 inode。大小与数量继续用整数。 |
| 目录合计口径 | 同一文件系统、只计普通文件、硬链接在单个合计内去重、跳过链接与下级挂载点；包含隐藏项，独立于显示筛选，失败返回 `null`。 |
| Git 字段的作用域和形状 | 状态给出 `state` 与 `xy`；普通目录聚合状态；分支和仓库状态仅对工作树根返回；历史记录用完整提交 ID，并按当前路径查询。 |
