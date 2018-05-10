package test

func get() *Ywaz {
	return &Ywaz{
		YW: &mime{i: "test"},
	}
}

func test() {
	get().Get("1")
}
