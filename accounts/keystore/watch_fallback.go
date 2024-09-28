//go:build (darwin && !cgo) || ios || (linux && arm64) || windows || (!darwin && !freebsd && !linux && !netbsd && !solaris)
// +build darwin,!cgo ios linux,arm64 windows !darwin,!freebsd,!linux,!netbsd,!solaris

// This is the fallback implementation of directory watching.
// It is used on unsupported platforms.

package keystore

type watcher struct {
	running  bool
	runEnded bool
}

func newWatcher(*accountCache) *watcher { return new(watcher) }
func (*watcher) start()                 {}
func (*watcher) close()                 {}

// enabled returns false on systems not supported.
func (*watcher) enabled() bool { return false }
