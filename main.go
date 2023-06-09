package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// 获取当前日期
	now := time.Now()
	year := fmt.Sprintf("%d", now.Year())
	month := fmt.Sprintf("%02d", now.Month())
	day := fmt.Sprintf("%02d", now.Day())

	// 创建一级、二级、三级文件夹
	folders := []string{year, month, day}
	for i := 1; i <= 3; i++ {
		path := filepath.Join(folders[:i]...)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, os.ModePerm)
		}
		os.Chdir(path)

		// 在每个文件夹中创建README.md文件
		if _, err := os.Stat("README.md"); os.IsNotExist(err) {
			f, err := os.Create("README.md")
			if err != nil {
				panic(err)
			}
			defer f.Close()
			if i == 3 {
				// 如果是当天的文件夹，则在README.md中记录日期信息
				f.WriteString("# " + now.Format("2006") + " 日志记录\n\n## Reference\n - []()\n\n## Summary\n")
			} else {
				// 如果不是当天的文件夹，则在README.md中记录路径信息
				f.WriteString("# " + now.Format("2006-01") + " 日志记录\n\n## Reference\n - []()\n\n## Summary\n")
			}
		}

		// 如果是月份文件夹，则再创建当天的文件夹
		if i == 2 {
			os.Chdir(path)
			todayPath := filepath.Join(".", day)
			if _, err := os.Stat(todayPath); os.IsNotExist(err) {
				os.Mkdir(todayPath, os.ModePerm)
			}
			os.Chdir(todayPath)

			// 在当天的文件夹中创建README.md文件
			if _, err := os.Stat("README.md"); os.IsNotExist(err) {
				f, err := os.Create("README.md")
				if err != nil {
					panic(err)
				}
				defer f.Close()
				f.WriteString("# " + now.Format("2006-01-02") + " 日志记录\n\n## Reference\n - []()\n\n## Summary\n")
			}

			// 在当天的文件夹中创建main.go文件
			if _, err := os.Stat("main.go"); os.IsNotExist(err) {
				f, err := os.Create("main.go")
				if os.IsNotExist(err) {
					panic(err)
				}
				defer f.Close()
				f.WriteString("package main\n\nfunc main() {\n\n}")
			}

			// 将当天的文件夹作为一个Go工作区
			if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
				workspacePath := fmt.Sprintf("/%s/%s/%s", year, month, day)
				err := ioutil.WriteFile("go.mod", []byte(fmt.Sprintf("module %s\n", "github.com/1ch0/go-daliy"+workspacePath)), 0666)
				if os.IsNotExist(err) {
					panic(err)
				}
				cmd := exec.Command("go", "mod", "init", workspacePath)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err = cmd.Run()
				if !strings.Contains(err.Error(), "exit status 1") {
					panic(err)
				}
			}

			// 切回上一级目录
			os.Chdir("..")
		}

		// 切回上一级目录
		os.Chdir("..")
	}
	fmt.Printf("%s Done!", now.Format("2006-01-02"))
}
