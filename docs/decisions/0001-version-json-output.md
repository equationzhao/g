# 1.0 直接切换到带版本号的 JSON 输出

状态：Accepted，2026-11-23。关联 F11、C05、D06。

0.31.2 的 JSON printer 从展示字段取得字符串，并按表头生成键。脚本需要解析可读单位和日期才能计算或比较。[基线实现](https://github.com/Equationzhao/g/blob/94ac309bbce6beca8628b933a8a0127d2a896ab5/internal/display/printer.go)

维护者确认 `g --json` 在 1.0 直接切换到带版本号的新格式。大小使用字节整数，修改时间使用带时区的 RFC 3339 字符串。原有脚本随版本升级调整。具体字段和路径规则见 [C05 格式设计](../v1/behavior-contract.md#c05--json)。
