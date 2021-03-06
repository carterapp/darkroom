package contain

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

func CreateContain(input string, output string, postfix string, size int) error {
	if fi, err := os.Stat(input); err != nil {
		return err
	} else if fi.Mode().IsDir() {
		if files, err := ioutil.ReadDir(input); err == nil {
			for _, f := range files {
				if err := CreateContain(filepath.Join(input, f.Name()), "", postfix, size); err != nil {
					return err
				}
			}
			return nil
		} else {
			return fmt.Errorf("Could not read directory: %s (%s)", input, err)
		}
	}

	src, err := imaging.Open(input)
	srcSize := src.Bounds().Size()
	width := size
	height := 0
	if srcSize.Y > srcSize.X {
		height = size
		width = 0
	}
	if err != nil {
		return fmt.Errorf("Could not open file: %s (%s)", input, err)
	}
	thumbnail := imaging.Resize(src, width, height, imaging.Lanczos)

	dest := output
	if dest == "" {
		ts := strings.SplitN(input, ".", 2)
		dest = ts[0] + postfix
		if len(ts) > 1 {
			dest += "." + ts[1]
		}
	}

	if err := imaging.Save(thumbnail, dest); err != nil {
		return fmt.Errorf("Could not save file: %s (%s)", dest, err)
	} else {
		return nil
	}
}
