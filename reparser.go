package views

import "fmt"
import "path/filepath"
import "github.com/howeyc/fsnotify"

var r *reparser

type reparsable interface {
	Filenames() []string
	ReparseTemplate() View
}

type reparser struct {
	reparsables map[string]reparsable
	done        chan bool
	watcher     *fsnotify.Watcher
}

func newReparser() *reparser {
	watcher, e := fsnotify.NewWatcher()
	if e != nil {
		panic(e)
	}

	return &reparser{
		reparsables: make(map[string]reparsable),
		done:        make(chan bool),
		watcher:     watcher,
	}
}

func (r *reparser) addReparsable(view reparsable) {
	for _, name := range view.Filenames() {
		name = normalizeFilename(name)
		fmt.Println("add: ", name)
		r.reparsables[name] = view
		e := r.watcher.Watch(name)
		noError(e)
	}
}

func (r *reparser) start() {
	fmt.Println("reparser: start")

	done := false
	go func() {
		fmt.Println("entering loop")
		for !done {
			fmt.Println("iteration")
			select {
			case ev := <-r.watcher.Event:
				filename := normalizeFilename(ev.Name)
				view, ok := r.reparsables[filename]
				if ok {
					view.ReparseTemplate()

					// HACK: Some editor rename things instead of writing to the file directly.
					// if ev.IsRename() {
					// 	e := r.watcher.Watch(ev.Name)
					// 	noError(e)
					// }
				}

			case err := <-r.watcher.Error:
				fmt.Println("error")
				panic(err)
			case <-r.done:
				fmt.Println("done")
				done = true
			}
		}
		fmt.Println("stopped.")
	}()
}

func (r *reparser) stop() {
	r.done <- true
}

func registerAutoReparseView(view reparsable) {
	if r == nil {
		r = newReparser()
	}

	r.addReparsable(view)
}

func EnableAutoReparse() {
	if r == nil {
		r = newReparser()
	}

	r.start()
}

func StopAutoReparse() {
	if r != nil {
		r.stop()
	}
	r = nil
}

func normalizeFilename(name string) string {
	absName, e := filepath.Abs(name)
	noError(e)

	return filepath.Clean(absName)
}
