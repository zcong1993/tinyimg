package main

import (
	"flag"
	"fmt"
	"github.com/apex/log"
	"github.com/discordapp/lilliput"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
)

const (
	// Name is cli name
	Name = "tinyimg"
	// Version is cli current version
	Version = "v0.0.2"
	// ExitCodeOK is exit code 0
	ExitCodeOK = 0
	// ExitCodeError is exit code 1
	ExitCodeError = 1
	// MaxParallel is parallel workers
	MaxParallel = 10
)

// GitCommit is cli current git commit hash
var GitCommit string

// WorkerLog is msg struct for worker
type WorkerLog struct {
	// IsError is msg type
	IsError bool
	// Msg is log msg
	Msg string
	// Input is handle filename
	Input string
	// Output is output filename
	Output string
}

// CLI is the command line object
type CLI struct {
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments
func (cli *CLI) Run(args []string) int {
	var (
		version  bool
		help     bool
		quantity int
		outdir   string
	)

	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&help, "help", false, "show help")
	flags.BoolVar(&help, "h", false, "show help(short)")

	flags.BoolVar(&version, "version", false, "show version")
	flags.BoolVar(&version, "v", false, "show version(short)")

	flags.IntVar(&quantity, "quantity", 85, "out put image quantity, 0-100")
	flags.IntVar(&quantity, "q", 85, "out put image quantity, 0-100(short)")

	flags.StringVar(&outdir, "outdir", ".", "output dir")
	flags.StringVar(&outdir, "o", ".", "output dir(short)")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	if version {
		ShowVersion()
		return ExitCodeOK
	}

	if help {
		fmt.Fprint(cli.outStream, helpText)
		fmt.Println()
		return ExitCodeOK
	}

	if quantity < 0 || quantity > 100 {
		LogError("Invalid argument: quantity should in 0 - 100")
		return ExitCodeError
	}

	if outdir != "" {
		if _, err := os.Stat(outdir); os.IsNotExist(err) {
			err := os.Mkdir(outdir, 0700)
			if err != nil {
				LogError(fmt.Sprintf("create dir error %s\n", err.Error()))
				return ExitCodeError
			}
		}
	}

	parsedArgs := flags.Args()
	if len(parsedArgs) < 1 {
		LogError("Invalid argument: you should provide at least 1 input file")
		return ExitCodeError
	}
	files := parsedArgs
	workers := minInt(len(files), MaxParallel)
	ch := make(chan WorkerLog, workers)
	for _, file := range files {
		go handler(file, quantity, outdir, ch)
	}

	for i := 0; i < len(files); i++ {
		out := <-ch
		if out.IsError {
			errLogger := CreateLogger(log.Fields{
				"file": out.Input,
				"msg":  out.Msg,
			})
			errLogger.Error("failed")
		} else {
			LogSuccess(log.Fields{
				"file":     out.Input,
				"output":   out.Output,
				"quantity": quantity,
			})
		}
	}
	return ExitCodeOK
}

func handler(file string, quantity int, outdir string, ch chan<- WorkerLog) {
	inputBuf, err := ioutil.ReadFile(file)
	if err != nil {
		ch <- WorkerLog{true, err.Error(), file, ""}
		return
	}
	decoder, err := lilliput.NewDecoder(inputBuf)
	if err != nil {
		ch <- WorkerLog{true, fmt.Sprintf("error decoding image, %s\n", err), file, ""}
		return
	}
	defer decoder.Close()
	header, err := decoder.Header()
	if err != nil {
		ch <- WorkerLog{true, fmt.Sprintf("error reading image header, %s\n", err), file, ""}
		return
	}
	ops := lilliput.NewImageOps(8192)
	defer ops.Close()
	outputImg := make([]byte, 20*1024*1024)
	outputType := "." + strings.ToLower(decoder.Description())
	outputWidth := header.Width()
	outputHeight := header.Height()
	resizeMethod := lilliput.ImageOpsFit
	encodeOptions := createOpts(quantity)
	opts := &lilliput.ImageOptions{
		FileType:             outputType,
		Width:                outputWidth,
		Height:               outputHeight,
		ResizeMethod:         resizeMethod,
		NormalizeOrientation: true,
		EncodeOptions:        encodeOptions[outputType],
	}
	outputImg, err = ops.Transform(decoder, opts, outputImg)
	if err != nil {
		ch <- WorkerLog{true, fmt.Sprintf("error transforming image, %s\n", err), file, ""}
		return
	}
	outputFilename := filepath.Join(outdir, fmt.Sprintf("tiny-%d-%s", quantity, file))
	if _, err := os.Stat(outputFilename); !os.IsNotExist(err) {
		ch <- WorkerLog{true, fmt.Sprintf("output filename %s exists, not replace it\n", outputFilename), file, ""}
		return
	}
	err = ioutil.WriteFile(outputFilename, outputImg, 0400)
	if err != nil {
		ch <- WorkerLog{true, fmt.Sprintf("error writing out resized image, %s\n", err), file, ""}
		return
	}

	ch <- WorkerLog{false, fmt.Sprintf("image written to %s\n", outputFilename), file, outputFilename}
}

func createOpts(quantity int) map[string]map[int]int {
	n := float64(quantity) / 10
	q := int(math.Abs(n))
	return map[string]map[int]int{
		".jpeg": map[int]int{lilliput.JpegQuality: quantity},
		".png":  map[int]int{lilliput.PngCompression: q},
		".webp": map[int]int{lilliput.WebpQuality: quantity},
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ShowVersion is handler for version command
func ShowVersion() {
	version := fmt.Sprintf("%s version %s", Name, Version)
	if len(GitCommit) != 0 {
		version += fmt.Sprintf(" (%s)", GitCommit)
	}
	fmt.Println(version)
}

var helpText = `
tinyimg is cli tool to compress jpeg, png or webp image.

 Usage:
 	tinyimg [options] SOURCE...
 Options:
	-quantity, -q               Output quantity, 0 - 100 int, default is 85
 	-outdir, -o                 Output dir, default ./
 Example:
 	tinyimg -q=30 -o=out *.jpg
`
