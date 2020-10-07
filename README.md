# About md-embed
把 `Markdown` 文件中的使用相对路径表示的图片编码成 `base64` 字符串，转化成`DataURL`，嵌入输出的 `markdown` 文件中。

语言：`go`

编译：

    go get gitee.com/rocket049/md-embed

用法：

```
	md-embed -o [Output filename] <Input filename>
	-o string
		Output file name (default "out.md")
```