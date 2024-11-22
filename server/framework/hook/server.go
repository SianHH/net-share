package hook

var serverBeforeHookList []func()

func AddServerBeforeHookFunc(f func()) {
	serverBeforeHookList = append(serverBeforeHookList, f)
}
