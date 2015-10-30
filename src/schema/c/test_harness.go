package c

import (
	"C"
	"unsafe"
	"schema/browse"
	"fmt"
	"strings"
	"bytes"
)

var readExpectations = map[string]string {
	"read one a" : `{"e":"a"}`,
	"read one b" : `{"e":"b"}`,
	"read two a" : `{"e":["a"]}`,
}

type testHarness struct {
	browser browse.Browser
	failed []string
	passed []string
}

//export conf2_testharness_new
func conf2_testharness_new(browser_hnd_id unsafe.Pointer) unsafe.Pointer {
	browser_hnd, found := GoHandles()[browser_hnd_id]
	if ! found {
		panic(fmt.Sprint("Browser not found", browser_hnd_id))
	}
	browser := browser_hnd.Data.(browse.Browser)
	harness := &testHarness{browser:browser}
	return NewGoHandle(harness).ID
}

func harness_from_handle_id(harness_hnd_id unsafe.Pointer) *testHarness {
	harness_hnd, found := GoHandles()[harness_hnd_id]
	if ! found {
		panic(fmt.Sprint("Test harness not found", harness_hnd_id))
	}
	return harness_hnd.Data.(*testHarness)
}

//export conf2_testharness_test_run
func conf2_testharness_test_run(harness_hnd_id unsafe.Pointer, c_testname *C.char) (passed C.short) {
	var err error
	harness := harness_from_handle_id(harness_hnd_id)

	testname := C.GoString(c_testname)
	details := strings.Split(testname, " ")
	var root browse.Node
	var state *browse.Selection
	if root, state, err = harness.browser.RootSelector(); err != nil {
		harness.failure(testname, err.Error())
		return FALSE_SHORT
	}
	var s browse.Node
	var path *browse.Path
	if path, err = browse.ParsePath(details[1]); err != nil {
		harness.failure(testname, err.Error())
		return FALSE_SHORT
	}
	if s, state, err = browse.WalkPath(state, root, path); err != nil {
		harness.failure(testname, err.Error())
		return FALSE_SHORT
	}
	var actual string
	switch details[0] {
	case "read":
		actual, err = tojson(state, s)
		if err != nil {
			harness.failure(testname, err.Error())
		} else {
			expected := readExpectations[testname]
			if expected != actual {
				failure := fmt.Sprintf("Expected\"%s\" Actual \"%s\"", expected, actual)
				harness.failure(testname, failure)
				return FALSE_SHORT
			}
		}
	default:
		harness.failure(testname, "Not a valid test")
		return FALSE_SHORT
	}

	harness.passed = append(harness.passed, testname)
	return TRUE_SHORT
}

//export conf2_testharness_report
func conf2_testharness_report(harness_hnd_id unsafe.Pointer) *C.char {
	harness := harness_from_handle_id(harness_hnd_id)
	var reportBuff bytes.Buffer
	reportBuff.WriteString(fmt.Sprintf("Passed : %d\n", len(harness.passed)))
	for _, pass := range harness.passed {
		reportBuff.WriteString(pass)
		reportBuff.WriteRune('\n')
	}

	reportBuff.WriteString(fmt.Sprintf("Failed : %d\n", len(harness.failed)))
	for _, fail := range harness.failed {
		reportBuff.WriteString(fail)
		reportBuff.WriteRune('\n')
	}

	return C.CString(string(reportBuff.Bytes()))
}

func (h *testHarness) failure(testname string, reason string) {
	failure := fmt.Sprintf("%s - %s", testname, reason)
	h.failed = append(h.failed, failure)
}

func tojson(state *browse.Selection, s browse.Node) (json string, err error) {
	var actual bytes.Buffer
	w := browse.NewJsonWriter(&actual)
	out, _ := w.GetSelector()
	if err = browse.SelectionInsert(state, s, out, browse.WalkAll()); err != nil {
		return
	}
	json = string(actual.Bytes())
	return
}
