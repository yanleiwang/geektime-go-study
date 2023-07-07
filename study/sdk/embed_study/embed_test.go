package embed_study

/*
embed 是在 Go 1.16 中新加包。它通过
//go:embed 指令，可以在编译阶段将静态资源文件打包进编译好的程序中，并提供访问这些文件的能力
基本语法非常简单，首先导入embed 包，然后使用指令//go:embed 文件名 将对应的文件或目录结构导入到对应的变量上。
*/

import (
	_ "embed"
	"fmt"
	"testing"
)

//go:embed embed.txt
var txt string

func TestEmbed(t *testing.T) {
	fmt.Printf("embed.txt: %q\n", txt)
}
