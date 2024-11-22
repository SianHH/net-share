package hook

func Run() {
	for _, hook := range serverBeforeHookList {
		hook()
	}
}
