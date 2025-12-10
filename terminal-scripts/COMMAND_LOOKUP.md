# Bash 命令解析与自定义命令生效原理

本文档解释为什么在终端里输入 `switch_go` 这类自定义命令可以被识别并执行，以及相关的底层机制（以 Bash 为例）。

## 命令解析顺序（从高到低优先级）

1. **别名 alias 展开**：词法分析阶段先把输入单词按别名替换。仅在交互或 `shopt expand_aliases` 开启时有效。  
2. **保留字检查**：`if`、`for`、`case` 等关键字在这一层被识别。  
3. **shell 函数**：如果存在同名函数，直接执行函数体（运行在当前 shell，不会启动新进程）。  
4. **内建命令 builtin**：`cd`、`echo` 等内建命令在这一层匹配。  
5. **哈希表缓存**：Bash 维护“命令名 → 绝对路径”的缓存，避免每次都遍历 PATH。缓存失效或 PATH 变化时会刷新。  
6. **PATH 搜索**：按 PATH 中目录顺序自左到右查找同名可执行文件，找到即执行。

> 只要你的自定义命令出现在上述任一层，它就能被终端识别。最常见的是“函数”或“PATH 中的可执行脚本”。

## PATH 如何决定命令版本

- PATH 是由多个目录组成的有序列表。解析到步骤 6 时，Bash 会按顺序查找命令名。  
- 目录越靠前，优先级越高；第一个命中的可执行文件即被执行。  
- 因此把 `~/go-current/bin` 放在 PATH 最前面，就能让 `go` 永远使用当前符号链接指向的版本。

## `switch_go` 如何被识别

- 脚本文件 `switch_go` 被放在 PATH 中的某个目录（如 `~/Desktop/program/gorder2/terminal-scripts`）。  
- 终端启动时加载 `~/.bashrc`，其中的 `export PATH="...:/home/jws/Desktop/program/gorder2/terminal-scripts:$PATH"` 把该目录加入 PATH。  
- 当你输入 `switch_go` 时，Bash 在前述解析顺序中到达 PATH 搜索阶段，找到这个可执行脚本并执行。

## 符号链接与版本切换

- `switch_go` 内部用 `ln -sfn` 更新符号链接 `~/go-current`，令其指向指定的 GOROOT。  
- 由于 PATH 把 `~/go-current/bin` 放在最前，`go` 命令解析时会命中该链接下的 `go` 可执行文件，实现版本切换。  
- 这一步是“改指针，不改 PATH”，因此切换迅速且不会累积多余路径。

## 子 shell 与当前 shell 的区别

- 以“直接执行脚本”的方式（例如 `switch_go ...`）会在子进程运行；它修改文件系统（符号链接）是持久的，但对子 shell 的 `PATH`/`GOROOT` 环境修改不会回写到父 shell。  
- 为了让父 shell 立刻看到新版本，只需确保父 shell 的 PATH 固定指向 `~/go-current/bin`。`switch_go` 修改符号链接后，父 shell 再次执行 `go` 就会命中新版本。  
- 如果需要修改父 shell 的环境变量（非本例），需要 `source script.sh` 或定义为函数。

## 与命令缓存的交互

- Bash 的哈希表可能缓存旧路径。更新符号链接后，如需强制刷新，可运行 `hash -r` 清空缓存。  
- 实际上，当 PATH 前缀不变、只改符号链接指向时，缓存命中仍然会跳到 `~/go-current/bin/go`，无需额外操作。

## 小结

- 把脚本所在目录加入 PATH，可执行脚本即成为终端可识别的命令。  
- 把 `~/go-current/bin` 放在 PATH 最前，通过更新符号链接即可切换 `go` 版本。  
- Bash 的解析顺序确保自定义函数/脚本在 PATH 中可被找到并执行，必要时可用 `hash -r` 刷新缓存。
