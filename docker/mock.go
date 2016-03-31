package docker

// Mock a Docker interface.
type Mock struct {
	// Handler
	RunHandler         func(option RunOptions, stream LogStream) (int, error)
	PingHandler        func() error
	LogsHandler        func(id string, stream LogStream) error
	BuildHandler       func(option BuildOptions, stream LogStream) error
	TagHandler         func(option TagOptions) error
	PushHandler        func(option PushOptions, stream LogStream) error
	ImageIDHandler     func(name string) string
	RemoveImageHandler func(name string) error
	// Values
	RunCalled         bool
	RunOptions        RunOptions
	PingCalled        bool
	LogsCalled        bool
	LogsID            string
	BuildCalled       bool
	BuildOptions      BuildOptions
	TagCalled         bool
	TagOptions        TagOptions
	PushCalled        bool
	PushOptions       PushOptions
	ImageIDCalled     bool
	ImageIDName       string
	RemoveImageCalled bool
	RemoveImageName   string
}

// NewMock return a mock of Docker interface.
func NewMock() *Mock {

	d := &Mock{}

	d.RunHandler = func(option RunOptions, stream LogStream) (int, error) {
		return 0, nil
	}

	d.PingHandler = func() error {
		return nil
	}

	d.LogsHandler = func(id string, stream LogStream) error {
		return nil
	}

	d.BuildHandler = func(option BuildOptions, stream LogStream) error {
		return nil
	}

	d.TagHandler = func(option TagOptions) error {
		return nil
	}

	d.PushHandler = func(option PushOptions, stream LogStream) error {
		return nil
	}

	d.ImageIDHandler = func(name string) string {
		return ""
	}

	d.RemoveImageHandler = func(name string) error {
		return nil
	}

	return d
}

// Run is a mock function for Docker interface.
func (m *Mock) Run(option RunOptions, stream LogStream) (int, error) {
	m.RunCalled = true
	m.RunOptions = option
	return m.RunHandler(option, stream)
}

// Ping is a mock function for Docker interface.
func (m *Mock) Ping() error {
	m.PingCalled = true
	return m.PingHandler()
}

// Logs is a mock function for Docker interface.
func (m *Mock) Logs(id string, stream LogStream) error {
	m.LogsCalled = true
	m.LogsID = id
	return m.LogsHandler(id, stream)
}

// Build is a mock function for Docker interface.
func (m *Mock) Build(option BuildOptions, stream LogStream) error {
	m.BuildCalled = true
	m.BuildOptions = option
	return m.BuildHandler(option, stream)
}

// Tag is a mock function for Docker interface.
func (m *Mock) Tag(option TagOptions) error {
	m.TagCalled = true
	m.TagOptions = option
	return m.TagHandler(option)
}

// Push is a mock function for Docker interface.
func (m *Mock) Push(option PushOptions, stream LogStream) error {
	m.PushCalled = true
	m.PushOptions = option
	return m.PushHandler(option, stream)
}

// ImageID is a mock function for Docker interface.
func (m *Mock) ImageID(name string) string {
	m.ImageIDCalled = true
	m.ImageIDName = name
	return m.ImageIDHandler(name)
}

// RemoveImage is a mock function for Docker interface.
func (m *Mock) RemoveImage(name string) error {
	m.RemoveImageCalled = true
	m.RemoveImageName = name
	return m.RemoveImageHandler(name)
}
