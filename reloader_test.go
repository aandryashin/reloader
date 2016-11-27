package reloader

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestWatch(t *testing.T) {
	dir, err := ioutil.TempDir("", "config")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dir)
	err = Watch(dir, func() {}, 1*time.Second)
	if err != nil {
		t.Errorf("uexpected error: %v\n", err)
	}
}

func TestWatchUexistentDir(t *testing.T) {
	dir, err := ioutil.TempDir("", "config")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(dir)
	err = Watch(dir, func() {}, 1*time.Second)
	if err == nil {
		t.Error("uexpected success")
	}
}

func TestNewWatcherError(t *testing.T) {
	fn := func() (*fsnotify.Watcher, error) {
		return nil, fmt.Errorf("error")
	}
	err := watch(fn, "", func() {}, 1*time.Second)
	if err == nil {
		t.Error("uexpected success")
	}
}

func TestTimerShouldNotTrigger(t *testing.T) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer watcher.Close()
	fn := func() (*fsnotify.Watcher, error) {
		return watcher, nil
	}
	dir, err := ioutil.TempDir("", "config")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dir)
	ch := make(chan bool)
	watch(fn, dir, func() { ch <- true }, 200*time.Millisecond)
	call := false
	watcher.Events <- fsnotify.Event{Op: fsnotify.Create}
	select {
	case <-time.After(100 * time.Millisecond):
	case call = <-ch:
	}
	if call {
		t.Error("timer should not be triggered after 100ms")
	}
}

func TestTimerShouldTrigger(t *testing.T) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer watcher.Close()
	fn := func() (*fsnotify.Watcher, error) {
		return watcher, nil
	}
	dir, err := ioutil.TempDir("", "config")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dir)
	ch := make(chan bool)
	watch(fn, dir, func() { ch <- true }, 100*time.Millisecond)
	call := false
	watcher.Events <- fsnotify.Event{Op: fsnotify.Create}
	select {
	case <-time.After(200 * time.Millisecond):
	case call = <-ch:
	}
	if !call {
		t.Error("timer should be triggered after 200ms")
	}
}

func TestTimerShouldTriggerOnce(t *testing.T) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer watcher.Close()
	fn := func() (*fsnotify.Watcher, error) {
		return watcher, nil
	}
	dir, err := ioutil.TempDir("", "config")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dir)
	ch := make(chan bool)
	watch(fn, dir, func() { ch <- true }, 300*time.Millisecond)
	call := false
	watcher.Events <- fsnotify.Event{Op: fsnotify.Create}
	select {
	case <-time.After(200 * time.Millisecond):
	case call = <-ch:
	}
	if call {
		t.Error("timer should not be triggered after 200ms")
	}
	watcher.Events <- fsnotify.Event{Op: fsnotify.Create}
	select {
	case <-time.After(200 * time.Millisecond):
	case call = <-ch:
	}
	if call {
		t.Error("timer still should not be triggered after 400ms")
	}
	select {
	case <-time.After(200 * time.Millisecond):
	case call = <-ch:
	}
	if !call {
		t.Error("timer should be triggered after 600ms")
	}
}
