package log

import (
	"testing"
)

func TestLogger_Test(t *testing.T) {
	l := newLogger("asdf")
	l.Test("Test logger Test")
	l.Test("Test", "Multiple", "Arguments")
	l.Test("Test", 5, 0.15, "different", []int{3, 5})
}

func TestLogger_Debug(t *testing.T) {
	l := newLogger("asdf")
	l.Debug("Debug logger Test")
	test1("Debug1")
	test2("Debug2")
}

func test1(message string) {
	l := newLogger("asdf")
	l.Debug(message)
}

func test2(message string) {
	l := newLogger("asdf")
	l.Debug(message)
	test3("Debug2-1")
}

func test3(message string) {
	l := newLogger("asdf")
	l.Debug(message)
}

func TestLogger_Info(t *testing.T) {
	l := newLogger("asdf")
	l.Info("Info logger Test")
	l.Info("Marshaller", "Test", "?3")
}

func TestLogger_Warn(t *testing.T) {
	l := newLogger("asdf")
	l.Warn("Warn logger Test")
}

func TestLogger_Error(t *testing.T) {
	l := newLogger("asdf")
	l.Error("Error logger Test")
}

func TestLogger_Fatal(t *testing.T) {
	l := newLogger("asdf")
	l.Fatal("Fatal logger Test")
}
