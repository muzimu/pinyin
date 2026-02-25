# py — 汉字转拼音首字母工具

读取标准输入中的汉字，输出拼音首字母，支持管道组合使用。

---

## 安装

### 前置条件

- Go 1.18+
- `$GOPATH/bin`（通常为 `~/go/bin`）已加入 `$PATH`

### 一键安装

在项目目录下执行：

```bash
go install .
```

安装后 `py` 即可在任意目录全局使用。

### 验证安装

```bash
which py        # 应输出 /Users/<你的用户名>/go/bin/py
echo "测试" | py
```

---

## 用法

```bash
echo "<汉字文本>" | py [选项]
```

### 选项

| 选项 | 简写 | 默认 | 说明 |
|------|------|------|------|
| `--lower` | `-l` | `true` | 首字母小写（默认） |
| `--upper` | `-u` | `false` | 首字母大写 |
| `--keep-non-han` | `-k` | `true` | 保留非汉字字符（数字、字母、符号等） |

> `-u` 与 `-l` 互斥，`-u` 显式指定时优先级更高。

---

## 示例

### 基础用法

```bash
# 默认小写
echo "张三" | py
# 输出: zs

# 显式大写
echo "张三" | py -u
# 输出: ZS

# 含非汉字字符（默认保留）
echo "张三 2024" | py
# 输出: zs 2024

# 过滤非汉字字符
echo "张三 2024" | py -k=false
# 输出: zs
```

### 配合 jq 使用

假设 `test.json` 内容如下：

```json
{
  "name": "张三",
  "city": "北京"
}
```

**读取单个字段的拼音首字母：**

```bash
jq -r '.name' test.json | py
# 输出: zs
```

**为JSON对象添加拼音首字母：**

```bash
jq -n '$initials | split("\n") as $a | inputs | . as $o | (input_line_number - 1) as $i | $a[$i] as $n | if (($o.name//"")!="" and $n!="") then $o + {name_py:$n} else $o end' --arg initials "$(jq -r '.name//""' test.json|py)" test.json
```

---

## 重新安装 / 更新

修改代码后重新执行即可覆盖：

```bash
go install .
```

---

## 帮助

```bash
py --help
```
