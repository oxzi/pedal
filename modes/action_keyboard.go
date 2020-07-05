package modes

import (
	"runtime"
	"time"

	"github.com/micmonay/keybd_event"
)

// KeyboardAction simulates key presses.
type KeyboardAction struct {
	keyBonding keybd_event.KeyBonding
	startTime  time.Time

	fun func(keybd_event.KeyBonding) error
}

// NewKeyboardCustomAction allows to configure some specific keyboard action to be performed.
func NewKeyboardCustomAction(fun func(keybd_event.KeyBonding) error) (kbdAction KeyboardAction, err error) {
	if kbdBonding, kbdErr := keybd_event.NewKeyBonding(); kbdErr != nil {
		err = kbdErr
		return
	} else {
		kbdAction = KeyboardAction{keyBonding: kbdBonding}
	}

	if runtime.GOOS == "linux" {
		kbdAction.startTime = time.Now().Add(2 * time.Second)
	}

	return
}

// NewKeyboardPressAction presses the keys passed from the key codes on action.
func NewKeyboardPressAction(keys []int) (kbdAction KeyboardAction, err error) {
	if kbdAction, err = NewKeyboardCustomAction(nil); err != nil {
		return
	}

	kbdAction.fun = func(keyBonding keybd_event.KeyBonding) error {
		keyBonding.SetKeys(keys...)

		if err := keyBonding.Launching(); err != nil {
			return err
		}
		keyBonding.Clear()

		return nil
	}

	return
}

// Execute this KeyboardAction.
func (kbdAction KeyboardAction) Execute() error {
	for time.Now().Before(kbdAction.startTime) {
		// active wait for startup delay to pass
	}

	return kbdAction.fun(kbdAction.keyBonding)
}
