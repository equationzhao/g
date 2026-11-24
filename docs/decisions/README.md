# 决策记录

只有影响跨模块边界、公共 CLI/配置/输出、功能退役或长期维护成本的决定才需要 ADR。局部函数实现、普通重命名等写在任务记录即可。

## 使用方式

1. 复制 [0000-template.md](0000-template.md)，命名为 `0001-short-topic.md` 等顺序文件。
2. 写出问题、约束、可选方案、推荐决定及影响的功能/契约 ID。
3. 未决定时保持 `Proposed`；真正决定后改为 `Accepted` 并写明任务/讨论/PR 依据。
4. 被新决定替代时改为 `Superseded`，链接后续记录；保留历史理由。

## 已接受的决定

- [ADR 0001 · 1.0 直接切换到带版本号的 JSON 输出](0001-version-json-output.md)

v1 的其他决策状态见[功能取舍表](../v1/feature-scope.md)和[行为契约](../v1/behavior-contract.md)。
