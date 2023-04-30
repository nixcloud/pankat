package pankat

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func PandocMarkdown2HTML(articleMarkdown []byte) (string, error) {
	pandocProcess := exec.Command("pandoc", "-f", "markdown", "-t", "html5", "--highlight-style", "kate")
	stdin, err := pandocProcess.StdinPipe()
	if err != nil {
		fmt.Println("An error occurred: ", err)
		return "", err
	}
	buff := bytes.NewBufferString("")
	pandocProcess.Stdout = buff
	pandocProcess.Stderr = os.Stderr
	err1 := pandocProcess.Start()
	if err1 != nil {
		fmt.Println("An error occurred: ", err1)
		return "", err1
	}
	_, err2 := io.WriteString(stdin, string(articleMarkdown))
	if err2 != nil {
		fmt.Println("An error occurred: ", err2)
		return "", err2
	}
	err3 := stdin.Close()
	if err3 != nil {
		fmt.Println("An error occurred: ", err3)
		return "", err3
	}
	err4 := pandocProcess.Wait()
	if err4 != nil {
		fmt.Println("An error occurred during pandocProess wait: ", err4)
		fmt.Println("An error occurred: ", err4)
	}
	return string(buff.Bytes()), nil
}
