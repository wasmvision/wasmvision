package runtime

import (
	"log"

	"github.com/orsinium-labs/wypes"
)

func hostedModules(logging bool) wypes.Modules {
	return wypes.Modules{
		"hosted": wypes.Module{
			"println": wypes.H1(hostPrintln),
			"log":     wypes.H1(hostLogFunc(logging)),
		},
	}
}

func hostPrintln(msg wypes.String) wypes.Void {
	println(msg.Unwrap())
	return wypes.Void{}
}

func hostLogFunc(logging bool) func(wypes.String) wypes.Void {
	return func(msg wypes.String) wypes.Void {
		if logging {
			log.Println(msg.Unwrap())
		}
		return wypes.Void{}
	}
}
