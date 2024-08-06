package tools

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
)

const path = "./tmp_testing"

func TestFSScan(t *testing.T) {
	type testcase struct {
		a, b, c int
	}
	tests := []testcase{{0, 0, 0}, {1, 1, 1}, {2, 2, 2}, {3, 3, 3}, {4, 4, 4}, {5, 5, 0}, {5, 5, 5}}
	for _, i := range tests {
		a, b, c := i.a, i.b, i.c
		os.RemoveAll(path)
		fcount := populateFS(t, path, a, b, c)
		var counter atomic.Int32
		scanFS(path, func(de fs.DirEntry, s string) { counter.Add(1) })
		if counter.Load() != int32(fcount) {
			t.Errorf("counter.Load(): %v, fcount: %v, testCase: %v\n", counter.Load(), fcount, i)
		}
	}
	os.RemoveAll(path)
}

func BenchmarkFSScan(b *testing.B) {
	b.Log("setting fs...")
	fcount := populateFS(b, path, 5, 6, 40)
	defer os.RemoveAll(path)
	b.Log("running bench...")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var counter atomic.Int64
		scanFS(path, func(de fs.DirEntry, s string) { counter.Add(1) })
		if counter.Load() != int64(fcount) {
			b.Errorf("counter.Load(): %v, fcount: %v, testCase: %v\n", counter.Load(), fcount, i)
		}

	}

}

func populateFS(t testing.TB, dir string, depth, dirWidth, fileWidth int) int {
	dirsToPopulate := []string{dir}
	err := os.Mkdir(dir, 0755)
	if err != nil {
		t.Error(err, dir)
	}
	count := 0
	for len(dirsToPopulate) > 0 {
		d := dirsToPopulate[len(dirsToPopulate)-1]
		dirsToPopulate = dirsToPopulate[:len(dirsToPopulate)-1]
		for i := range dirWidth {
			fp := filepath.Join(d, fmt.Sprint(i))
			err := os.Mkdir(fp, 0755)
			if err != nil {
				t.Error(err, d)
			}
			if spl := strings.Split(fp, "/"); len(spl) > depth {
				continue
			}
			dirsToPopulate = append(dirsToPopulate, fp)
		}
		for range fileWidth {
			fName, err := os.CreateTemp(d, "file_for_test_*")
			if err != nil {
				t.Error(err, d)
			}
			fName.Close()
			count++
		}

	}
	return count
}
