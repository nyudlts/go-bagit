package go_bagit_test

/*
func TestWithLogger(t *testing.T) {
	t.Run("It panics when a nil value is passed", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("WithLogger did not panic")
			}
		}()

		go_bagit.WithLogger(nil)
	})

	t.Run("It replaces the library logger", func(t *testing.T) {
		var buf bytes.Buffer
		tlog := stdlog.New(&buf, "", stdlog.LstdFlags)

		go_bagit.WithLogger(tlog)

		_ = go_bagit.ValidateBag("/dev/null", false, false)
		want := []byte("ERROR - open /dev/null/bag-info.txt: not a directory")
		if !bytes.Contains(buf.Bytes(), want) {
			t.Fatal("WithLogger did not replace the library logger")
		}
	})
}

func TestLogger(t *testing.T) {
	t.Run("It returns the current logger", func(t *testing.T) {
		go_bagit.WithLogger(stdlog.Default())

		logger := go_bagit.Logger()

		if logger != stdlog.Default() {
			t.Fatal("Logger returned an expected value")
		}
	})
}
*/
