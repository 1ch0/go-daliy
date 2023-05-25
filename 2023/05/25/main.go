package main

import (
	"fmt"

	"github.com/carlmjohnson/versioninfo"
)

func main() {
	fmt.Println("Shor:", versioninfo.Short())
	fmt.Println("Last Commit:", versioninfo.LastCommit)
	fmt.Println("Version:", versioninfo.Version)
	fmt.Println("Revision:", versioninfo.Revision)
	fmt.Println("DirtyBuild:", versioninfo.DirtyBuild)
	fmt.Println("LastCommit:", versioninfo.LastCommit)
}
