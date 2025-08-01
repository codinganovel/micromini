package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/pprof"
	"sort"
	"strconv"
	"syscall"
	"time"

	"github.com/codinganovel/autocd-go"
	"github.com/go-errors/errors"
	isatty "github.com/mattn/go-isatty"
	"github.com/micro-editor/tcell/v2"
	"github.com/zyedidia/micro/v2/internal/action"
	"github.com/zyedidia/micro/v2/internal/buffer"
	"github.com/zyedidia/micro/v2/internal/clipboard"
	"github.com/zyedidia/micro/v2/internal/config"
	"github.com/zyedidia/micro/v2/internal/screen"
	"github.com/zyedidia/micro/v2/internal/shell"
	"github.com/zyedidia/micro/v2/internal/util"
)

var (
	// Command line flags
	flagVersion   = flag.Bool("version", false, "Show the version number and information")
	flagConfigDir = flag.String("config-dir", "", "Specify a custom location for the configuration directory")
	flagOptions   = flag.Bool("options", false, "Show all option help")
	flagDebug     = flag.Bool("debug", false, "Enable debug mode (prints debug info to ./log.txt)")
	flagProfile   = flag.Bool("profile", false, "Enable CPU profiling (writes profile info to ./micro.prof)")
	flagClean     = flag.Bool("clean", false, "Clean configuration directory")
	flagAutoCD    = flag.Bool("autocd", false, "Change to file directory on exit")
	optionFlags   map[string]*string

	sighup chan os.Signal

	timerChan chan func()
)

func InitFlags() {
	flag.Usage = func() {
		fmt.Println("Usage: micro [OPTIONS] [FILE]...")
		fmt.Println("-autocd")
		fmt.Println("    \tChange to file directory on exit")
		fmt.Println("-clean")
		fmt.Println("    \tCleans the configuration directory")
		fmt.Println("-config-dir dir")
		fmt.Println("    \tSpecify a custom location for the configuration directory")
		fmt.Println("[FILE]:LINE:COL (if the `parsecursor` option is enabled)")
		fmt.Println("+LINE:COL")
		fmt.Println("    \tSpecify a line and column to start the cursor at when opening a buffer")
		fmt.Println("-options")
		fmt.Println("    \tShow all option help")
		fmt.Println("-debug")
		fmt.Println("    \tEnable debug mode (enables logging to ./log.txt)")
		fmt.Println("-profile")
		fmt.Println("    \tEnable CPU profiling (writes profile info to ./micro.prof")
		fmt.Println("    \tso it can be analyzed later with \"go tool pprof micro.prof\")")
		fmt.Println("-version")
		fmt.Println("    \tShow the version number and information")

		fmt.Print("\nMicro's options can also be set via command line arguments for quick\nadjustments. For real configuration, please use the settings.json\nfile (see 'help options').\n\n")
		fmt.Println("-option value")
		fmt.Println("    \tSet `option` to `value` for this session")
		fmt.Println("    \tFor example: `micro -syntax off file.c`")
		fmt.Println("\nUse `micro -options` to see the full list of configuration options")
	}

	optionFlags = make(map[string]*string)

	for k, v := range config.DefaultAllSettings() {
		optionFlags[k] = flag.String(k, "", fmt.Sprintf("The %s option. Default value: '%v'.", k, v))
	}

	flag.Parse()

	if *flagVersion {
		// If -version was passed
		fmt.Println("Version:", util.Version)
		fmt.Println("Commit hash:", util.CommitHash)
		fmt.Println("Compiled on", util.CompileDate)
		exit(0)
	}

	if *flagOptions {
		// If -options was passed
		var keys []string
		m := config.DefaultAllSettings()
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := m[k]
			fmt.Printf("-%s value\n", k)
			fmt.Printf("    \tDefault value: '%v'\n", v)
		}
		exit(0)
	}

	if util.Debug == "OFF" && *flagDebug {
		util.Debug = "ON"
	}
}

// DoCleanFlag handles the -clean flag
func DoCleanFlag() {
	if *flagClean {
		CleanConfig()
		exit(0)
	}
}

// LoadInput determines which files should be loaded into buffers
// based on the input stored in flag.Args()
func LoadInput(args []string) []*buffer.Buffer {
	// There are a number of ways micro should start given its input

	// 1. If it is given a files in flag.Args(), it should open those

	// 2. If there is no input file and the input is not a terminal, that means
	// something is being piped in and the stdin should be opened in an
	// empty buffer

	// 3. If there is no input file and the input is a terminal, an empty buffer
	// should be opened

	var filename string
	var input []byte
	var err error
	buffers := make([]*buffer.Buffer, 0, len(args))

	btype := buffer.BTDefault
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		btype = buffer.BTStdout
	}

	files := make([]string, 0, len(args))
	flagStartPos := buffer.Loc{-1, -1}
	flagr := regexp.MustCompile(`^\+(\d+)(?::(\d+))?$`)
	for _, a := range args {
		match := flagr.FindStringSubmatch(a)
		if len(match) == 3 && match[2] != "" {
			line, err := strconv.Atoi(match[1])
			if err != nil {
				screen.TermMessage(err)
				continue
			}
			col, err := strconv.Atoi(match[2])
			if err != nil {
				screen.TermMessage(err)
				continue
			}
			flagStartPos = buffer.Loc{col - 1, line - 1}
		} else if len(match) == 3 && match[2] == "" {
			line, err := strconv.Atoi(match[1])
			if err != nil {
				screen.TermMessage(err)
				continue
			}
			flagStartPos = buffer.Loc{0, line - 1}
		} else {
			files = append(files, a)
		}
	}

	if len(files) > 0 {
		// Option 1
		// We go through each file and load it
		for i := 0; i < len(files); i++ {
			buf, err := buffer.NewBufferFromFileAtLoc(files[i], btype, flagStartPos)
			if err != nil {
				screen.TermMessage(err)
				continue
			}
			// If the file didn't exist, input will be empty, and we'll open an empty buffer
			buffers = append(buffers, buf)
		}
	} else if !isatty.IsTerminal(os.Stdin.Fd()) {
		// Option 2
		// The input is not a terminal, so something is being piped in
		// and we should read from stdin
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			screen.TermMessage("Error reading from stdin: ", err)
			input = []byte{}
		}
		buffers = append(buffers, buffer.NewBufferFromStringAtLoc(string(input), filename, btype, flagStartPos))
	} else {
		// Option 3, just open an empty buffer
		buffers = append(buffers, buffer.NewBufferFromStringAtLoc(string(input), filename, btype, flagStartPos))
	}

	return buffers
}

func checkBackup(name string) error {
	target := filepath.Join(config.ConfigDir, name)
	backup := util.AppendBackupSuffix(target)
	if info, err := os.Stat(backup); err == nil {
		input, err := os.ReadFile(backup)
		if err == nil {
			t := info.ModTime()
			msg := fmt.Sprintf(buffer.BackupMsg, target, t.Format("Mon Jan _2 at 15:04, 2006"), backup)
			choice := screen.TermPrompt(msg, []string{"r", "i", "a", "recover", "ignore", "abort"}, true)

			if choice%3 == 0 {
				// recover
				err := os.WriteFile(target, input, util.FileMode)
				if err != nil {
					return err
				}
				return os.Remove(backup)
			} else if choice%3 == 1 {
				// delete
				return os.Remove(backup)
			} else if choice%3 == 2 {
				// abort
				return errors.New("Aborted")
			}
		}
	}
	return nil
}

func exit(rc int) {
	for _, b := range buffer.OpenBuffers {
		if !b.Modified() {
			b.Fini()
		}
	}

	if screen.Screen != nil {
		screen.Screen.Fini()
	}

	// Try to use autocd to change to the directory of the current file
	// This allows seamless directory navigation when the editor exits
	if *flagAutoCD && action.Tabs != nil && len(action.Tabs.List) > 0 {
		tab := action.Tabs.List[action.Tabs.Active()]
		if tab != nil && len(tab.Panes) > 0 {
			for _, pane := range tab.Panes {
				if bp, ok := pane.(*action.BufPane); ok && bp.Buf.Path != "" {
					targetDir := filepath.Dir(bp.Buf.Path)
					// Try autocd, but fallback to normal exit if it fails
					autocd.ExitWithDirectoryOrFallback(targetDir, func() {
						os.Exit(rc)
					})
					return
				}
			}
		}
	}

	os.Exit(rc)
}

func main() {
	defer func() {
		if util.Stdout.Len() > 0 {
			fmt.Fprint(os.Stdout, util.Stdout.String())
		}
		exit(0)
	}()

	var err error

	InitFlags()

	if *flagProfile {
		f, err := os.Create("micro.prof")
		if err != nil {
			log.Fatal("error creating CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("error starting CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	InitLog()

	err = config.InitConfigDir(*flagConfigDir)
	if err != nil {
		screen.TermMessage(err)
	}

	config.InitRuntimeFiles(true)

	err = checkBackup("settings.json")
	if err != nil {
		screen.TermMessage(err)
		exit(1)
	}

	err = config.ReadSettings()
	if err != nil {
		screen.TermMessage(err)
	}
	err = config.InitGlobalSettings()
	if err != nil {
		screen.TermMessage(err)
	}

	// flag options
	for k, v := range optionFlags {
		if *v != "" {
			nativeValue, err := config.GetNativeValue(k, *v)
			if err != nil {
				screen.TermMessage(err)
				continue
			}
			if err = config.OptionIsValid(k, nativeValue); err != nil {
				screen.TermMessage(err)
				continue
			}
			config.GlobalSettings[k] = nativeValue
			config.VolatileSettings[k] = true
		}
	}

	DoCleanFlag()

	err = screen.Init()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Fatal: Micro could not initialize a Screen.")
		exit(1)
	}
	m := clipboard.SetMethod(config.GetGlobalOption("clipboard").(string))
	clipErr := clipboard.Initialize(m)

	defer func() {
		if err := recover(); err != nil {
			if screen.Screen != nil {
				screen.Screen.Fini()
			}
			fmt.Println("Micro encountered an error:", errors.Wrap(err, 2).ErrorStack(), "\nIf you can reproduce this error, please report it at https://github.com/zyedidia/micro/issues")
			// backup all open buffers
			for _, b := range buffer.OpenBuffers {
				b.Backup()
			}
			exit(1)
		}
	}()

	err = checkBackup("bindings.json")
	if err != nil {
		screen.TermMessage(err)
		exit(1)
	}

	action.InitBindings()
	action.InitCommands()

	action.InitGlobals()
	buffer.SetMessager(action.InfoBar)
	args := flag.Args()
	b := LoadInput(args)

	if len(b) == 0 {
		// No buffers to open
		screen.Screen.Fini()
		runtime.Goexit()
	}

	action.InitTabs(b)

	err = config.InitColorscheme()
	if err != nil {
		screen.TermMessage(err)
	}

	if clipErr != nil {
		log.Println(clipErr, " or change 'clipboard' option")
	}

	config.StartAutoSave()
	if a := config.GetGlobalOption("autosave").(float64); a > 0 {
		config.SetAutoTime(a)
	}

	screen.Events = make(chan tcell.Event)

	util.Sigterm = make(chan os.Signal, 1)
	sighup = make(chan os.Signal, 1)
	signal.Notify(util.Sigterm, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT)
	signal.Notify(sighup, syscall.SIGHUP)

	timerChan = make(chan func())

	// Here is the event loop which runs in a separate thread
	go func() {
		for {
			screen.Lock()
			e := screen.Screen.PollEvent()
			screen.Unlock()
			if e != nil {
				screen.Events <- e
			}
		}
	}()

	// clear the drawchan so we don't redraw excessively
	// if someone requested a redraw before we started displaying
	for len(screen.DrawChan()) > 0 {
		<-screen.DrawChan()
	}

	// wait for initial resize event
	select {
	case event := <-screen.Events:
		action.Tabs.HandleEvent(event)
	case <-time.After(10 * time.Millisecond):
		// time out after 10ms
	}

	for {
		DoEvent()
	}
}

// DoEvent runs the main action loop of the editor
func DoEvent() {
	var event tcell.Event

	// Display everything
	screen.Screen.Fill(' ', config.DefStyle)
	screen.Screen.HideCursor()
	action.Tabs.Display()
	for _, ep := range action.MainTab().Panes {
		ep.Display()
	}
	action.MainTab().Display()
	action.InfoBar.Display()
	screen.Screen.Show()

	// Check for new events
	select {
	case f := <-shell.Jobs:
		// If a new job has finished while running in the background we should execute the callback
		f.Function(f.Output, f.Args)
	case <-config.Autosave:
		for _, b := range buffer.OpenBuffers {
			b.AutoSave()
		}
	case <-shell.CloseTerms:
		action.Tabs.CloseTerms()
	case event = <-screen.Events:
	case <-screen.DrawChan():
		for len(screen.DrawChan()) > 0 {
			<-screen.DrawChan()
		}
	case f := <-timerChan:
		f()
	case b := <-buffer.BackupCompleteChan:
		b.RequestedBackup = false
	case <-sighup:
		exit(0)
	case <-util.Sigterm:
		exit(0)
	}

	if e, ok := event.(*tcell.EventError); ok {
		log.Println("tcell event error: ", e.Error())

		if e.Err() == io.EOF {
			// shutdown due to terminal closing/becoming inaccessible
			exit(0)
		}
		return
	}

	if event != nil {
		_, resize := event.(*tcell.EventResize)
		if resize {
			action.InfoBar.HandleEvent(event)
			action.Tabs.HandleEvent(event)
		} else if action.InfoBar.HasPrompt {
			action.InfoBar.HandleEvent(event)
		} else {
			action.Tabs.HandleEvent(event)
		}
	}

}
