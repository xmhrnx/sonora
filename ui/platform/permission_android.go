package platform

import (
	"errors"

	"gioui.org/app"
	_ "gioui.org/app/permission/microphone" // RECORD_AUDIO in the manifest
	"git.wow.st/gmp/jni"

	"github.com/mehran9675/sonora/ui/tr"
)

const recordAudio = "android.permission.RECORD_AUDIO"

// errMicDenied is returned when the user doesn't allow the microphone.
var errMicDenied = errors.New(tr.Msg("Sonora needs the microphone to record. Allow it when Android asks, or in the app's settings"))

// micGranted reports whether the app holds RECORD_AUDIO.
func micGranted(vm jni.JVM) bool {
	granted := false
	jni.Do(vm, func(env jni.Env) error {
		appCtx := jni.Object(app.AppContext())
		check := jni.GetMethodID(env, jni.GetObjectClass(env, appCtx), "checkSelfPermission", "(Ljava/lang/String;)I")
		res, err := jni.CallIntMethod(env, appCtx, check, jni.Value(jni.JavaString(env, recordAudio)))
		granted = err == nil && res == 0 // PackageManager.PERMISSION_GRANTED
		return err
	})
	return granted
}

// requestMic shows the permission prompt through the activity that
// hosts view. Run it on the main thread.
func requestMic(vm jni.JVM, view uintptr) error {
	return jni.Do(vm, func(env jni.Env) error {
		v := jni.Object(view)
		getContext := jni.GetMethodID(env, jni.GetObjectClass(env, v), "getContext", "()Landroid/content/Context;")
		activity, err := jni.CallObjectMethod(env, v, getContext)
		if err != nil {
			return err
		}
		perms := jni.NewObjectArray(env, 1, jni.FindClass(env, "java/lang/String"), jni.Object(jni.JavaString(env, recordAudio)))
		request := jni.GetMethodID(env, jni.GetObjectClass(env, activity), "requestPermissions", "([Ljava/lang/String;I)V")
		return jni.CallVoidMethod(env, activity, request, jni.Value(perms), jni.Value(1))
	})
}

// inputTypes names the AudioDeviceInfo types that are microphones;
// others (telephony, tuners, echo reference…) are left out.
var inputTypes = map[int32]string{
	3:  tr.Msg("Wired headset"),     // TYPE_WIRED_HEADSET
	5:  tr.Msg("Line in"),           // TYPE_LINE_ANALOG
	6:  tr.Msg("Digital line in"),   // TYPE_LINE_DIGITAL
	7:  tr.Msg("Bluetooth headset"), // TYPE_BLUETOOTH_SCO
	11: tr.Msg("USB microphone"),    // TYPE_USB_DEVICE
	12: tr.Msg("USB accessory"),     // TYPE_USB_ACCESSORY
	13: tr.Msg("Dock"),              // TYPE_DOCK
	15: tr.Msg("Phone microphone"),  // TYPE_BUILTIN_MIC
	22: tr.Msg("USB headset"),       // TYPE_USB_HEADSET
	26: tr.Msg("Bluetooth headset"), // TYPE_BLE_HEADSET
}
