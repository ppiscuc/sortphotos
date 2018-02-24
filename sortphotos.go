package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/xiam/exif"
)

// brew install libexif

var (
	src           *string
	dst           *string
	inplace       *bool
	info          *bool
	infodatetaken *bool
	sort          *bool
)

func main() {
	src = flag.String("src", "", "src directory of photos")
	dst = flag.String("dst", "", "dst directory for photos")
	info = flag.Bool("info", false, "show only information")
	infodatetaken = flag.Bool("infodate", false, "info about the date taken")
	sort = flag.Bool("sort", false, "sort the photos")
	flag.Parse()
	if *src == "" && *dst == "" {
		fmt.Println("src and dst must be specified")
		os.Exit(1)
	}
	if *info == false && *infodatetaken == false && *sort == false {
		fmt.Println("either info, infodatataken or sort must be specified")
		os.Exit(1)
	}
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

// run performs the filewalk of src
// it skips files that are not .jpg or .jpeg
func run() error {
	err := filepath.Walk(*src, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		lowerext := strings.ToLower(filepath.Ext(info.Name()))
		if lowerext != ".jpg" && lowerext != ".jpeg" {
			fmt.Printf("skipping\t%s\n", info.Name())
			return nil
		}
		return runOnFile(path)
	})
	return err
}

// runOnFile processes a single file
func runOnFile(fph string) error {
	data, err := exif.Read(fph)
	if err != nil {
		return errors.Wrap(err, "failed to read filename")
	}
	if *info == true {
		showInfo(data)
		return nil
	}

	if *infodatetaken == true || *sort == true {
		t, err := processData(data)
		if err != nil {
			return errors.Wrap(err, "failed to process")
		}
		if *infodatetaken == true {
			fmt.Printf("%s\t%s\n", filepath.Base(fph), t)
			return nil
		}
		//we need to sort in format: 2017/12-25
		year := t.Format("2006")
		monthday := t.Format("01-02")
		adst, err4 := filepath.Abs(*dst)
		if err4 != nil {
			return err
		}
		newpath := filepath.Join(adst, year, monthday)
		if _, err1 := os.Stat(newpath); os.IsNotExist(err1) {
			if err2 := os.MkdirAll(newpath, 0770); err2 != nil {
				return errors.Wrap(err2, "failed to create dir")
			}
		}
		ext := filepath.Ext(fph)
		newbasefilename := t.Format(time.RFC3339) + ext
		absnewfilename := filepath.Join(newpath, newbasefilename)
		if err5 := CopyFile(absnewfilename, fph, 0600); err5 != nil {
			return errors.Wrap(err5, "failed to copy")
		}
	}
	return nil
}

func showInfo(data *exif.Data) {
	for key, val := range data.Tags {
		fmt.Printf("%s = %s\n", key, val)
	}
}

func processData(data *exif.Data) (time.Time, error) {
	// format: 2017:05:24 21:03:35
	format := "2006:01:02 15:04:05"
	if odate, ok := data.Tags["Date and Time (Original)"]; ok {
		return time.Parse(format, odate)
	}
	if digitizedDate, ok := data.Tags["Date and Time (Digitized)"]; ok {
		return time.Parse(format, digitizedDate)
	}
	if dt, ok := data.Tags["Date and Time"]; ok {
		return time.Parse(format, dt)
	}
	return time.Time{}, fmt.Errorf("no data and time found")
}

// CopyFile copies the contents from src to dst using io.Copy.
// If dst does not exist, CopyFile creates it with permissions perm;
// otherwise CopyFile truncates it before writing.
func CopyFile(dst, src string, perm os.FileMode) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()
	_, err = io.Copy(out, in)
	return
}
