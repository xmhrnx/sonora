package platform

import (
	"fmt"
	"strconv"

	"gioui.org/app"
	_ "gioui.org/app/permission/microphone" // RECORD_AUDIO in the manifest
	"git.wow.st/gmp/jni"

	"github.com/mehran9675/sonora/internal/audio/record"
	"github.com/mehran9675/sonora/ui/tr"
)

// Microphones lists the inputs AudioManager reports; miniaudio sees
// only the default one on Android. It needs no permission.
func Microphones() ([]record.Device, error) {
	var out []record.Device
	seen := map[string]int{}
	err := jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		appCtx := jni.Object(app.AppContext())
		getService := jni.GetMethodID(env, jni.GetObjectClass(env, appCtx), "getSystemService", "(Ljava/lang/String;)Ljava/lang/Object;")
		am, err := jni.CallObjectMethod(env, appCtx, getService, jni.Value(jni.JavaString(env, "audio")))
		if err != nil {
			return err
		}
		getDevices := jni.GetMethodID(env, jni.GetObjectClass(env, am), "getDevices", "(I)[Landroid/media/AudioDeviceInfo;")
		arr, err := jni.CallObjectMethod(env, am, getDevices, jni.Value(1)) // GET_DEVICES_INPUTS
		if err != nil {
			return err
		}
		arrays := jni.FindClass(env, "java/lang/reflect/Array")
		getLength := jni.GetStaticMethodID(env, arrays, "getLength", "(Ljava/lang/Object;)I")
		n, err := jni.CallStaticIntMethod(env, arrays, getLength, jni.Value(arr))
		if err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			info, err := jni.GetObjectArrayElement(env, jni.ObjectArray(arr), jni.Size(i))
			if err != nil {
				return err
			}
			cls := jni.GetObjectClass(env, info)
			id, err := jni.CallIntMethod(env, info, jni.GetMethodID(env, cls, "getId", "()I"))
			if err != nil {
				return err
			}
			kind, err := jni.CallIntMethod(env, info, jni.GetMethodID(env, cls, "getType", "()I"))
			if err != nil {
				return err
			}
			label, ok := inputTypes[kind]
			label = tr.T(label)
			if !ok {
				continue
			}
			// The product name tells headsets apart; built-in mics all
			// carry the phone's model, so they get a number instead.
			if kind != 15 {
				if name := javaText(env, info, cls, "getProductName"); name != "" {
					label += " · " + name
				}
			}
			seen[label]++
			if seen[label] > 1 {
				label = fmt.Sprintf("%s %d", label, seen[label])
			}
			out = append(out, record.Device{ID: record.AndroidPrefix + strconv.Itoa(int(id)), Name: label})
		}
		return nil
	})
	return out, err
}

// javaText calls a no-argument method returning a CharSequence and
// returns it as a string, or "".
func javaText(env jni.Env, obj jni.Object, cls jni.Class, method string) string {
	cs, err := jni.CallObjectMethod(env, obj, jni.GetMethodID(env, cls, method, "()Ljava/lang/CharSequence;"))
	if err != nil || cs == 0 {
		return ""
	}
	toString := jni.GetMethodID(env, jni.GetObjectClass(env, cs), "toString", "()Ljava/lang/String;")
	s, err := jni.CallObjectMethod(env, cs, toString)
	if err != nil || s == 0 {
		return ""
	}
	return jni.GoString(env, jni.String(s))
}
