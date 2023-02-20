package editor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func OpenAndRead() []byte {
	f, err := createTempFile()
	if err != nil {
		panic(err)
	}
	defer remove(f)

	// close the file and open in the default editor
	f.Close()
	openFileInDefaultEditor(f)

	o, err := os.ReadFile(f.Name())
	if err != nil {
		panic(err)
	}
	return o
}

func remove(f *os.File) {
	err := os.Remove(f.Name())
	if err != nil {
		fmt.Printf("%v", err)
	}
}

func createTempFile() (*os.File, error) {
	f, err := os.CreateTemp("", "tip-k8ctx-*.yaml")
	if err != nil {
		return nil, err
	}

	return f, nil
}

func openFileInDefaultEditor(f *os.File) {
	fmt.Println(f.Name())
	cmd := exec.Command(getEditor(), f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic("oops")
	}
	fmt.Println(cmd.String())
}

func getEditor() string {
	editor, found := os.LookupEnv("EDITOR")
	if found {
		return editor
	}

	if runtime.GOOS == "windows" {
		return "notepad"
	}

	return ""
}
