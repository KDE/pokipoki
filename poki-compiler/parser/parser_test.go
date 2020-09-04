package parser

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func Reexec(t *testing.T, test string, expectedCode int) bool {
	if os.Getenv("POKIPOKI_TEST") == "1" {
		return true
	}

	cmd := exec.Command(os.Args[0], "-test.run="+test)
	cmd.Env = append(os.Environ(), "POKIPOKI_TEST=1")
	err := cmd.Run()

	if expectedCode == 0 && err == nil {
		return false
	} else if err == nil {
		t.Fatalf("process ran with zero exit code, want exit code %d", expectedCode)
	}

	exitError, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("other error encountered: %+v", err)
	}

	if exitError.ExitCode() != expectedCode {
		t.Fatalf("process ran with %d exit code, want exit code %d", exitError.ExitCode(), expectedCode)
	}

	return false
}

func TestCustomScannerNoEOF(t *testing.T) {
	if Reexec(t, "TestCustomScannerNoEOF", 1) {
		s := customScanner{}
		s.Init(strings.NewReader(`yeet`))
		s.ScanNoEOF()
		s.ScanNoEOF()
	}
}

func TestCustomScannerIdent(t *testing.T) {
	s := customScanner{}
	s.Init(strings.NewReader(`ident iDent iEdent iDenT`))
	s.ScanIdent()
	s.ScanIdent()
	s.ScanIdent()
	s.ScanIdent()
}

func TestCustomScannerIdentXFail(t *testing.T) {
	if Reexec(t, "TestCustomScannerIdentXFail", 1) {
		s := customScanner{}
		s.Init(strings.NewReader(`Ident`))
		s.ScanIdent()
	}
}

func TestCustomScannerName(t *testing.T) {
	s := customScanner{}
	s.Init(strings.NewReader(`Ident IDent`))
	s.ScanName()
	s.ScanName()
}

func TestCustomScannerNameXFail(t *testing.T) {
	if Reexec(t, "TestCustomScannerNameXFail", 1) {
		s := customScanner{}
		s.Init(strings.NewReader(`iDent`))
		s.ScanName()
	}
}

func TestCustomScannerExpecting(t *testing.T) {
	if Reexec(t, "TestCustomScannerExpecting", 1) {
		s := customScanner{}
		s.Init(strings.NewReader(`{ } |`))
		s.ScanExpecting(`{`)
		s.ScanExpecting(`}`)
		s.ScanExpecting(`|`)
	}
}

func TestCustomScannerExpectingXFail(t *testing.T) {
	if Reexec(t, "TestCustomScannerExpectingXFail", 1) {
		s := customScanner{}
		s.Init(strings.NewReader(`|`))
		s.ScanExpecting(`{`)
	}
}

func TestCustomScannerNumber(t *testing.T) {
	if Reexec(t, "TestCustomScannerNumber", 1) {
		s := customScanner{}
		s.Init(strings.NewReader(`123`))
		s.Init(strings.NewReader(`456`))
		s.Init(strings.NewReader(`789`))
		s.Init(strings.NewReader(`410297319824612897361`))
		s.ScanNumber()
		s.ScanNumber()
		s.ScanNumber()
		s.ScanNumber()
	}
}

func TestCustomScannerNumberXFail(t *testing.T) {
	if Reexec(t, "TestCustomScannerNumberXFail", 1) {
		s := customScanner{}
		s.Init(strings.NewReader(`a123`))
		s.ScanNumber()
	}
}
