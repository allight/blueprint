// Copyright 2019 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// bpdoc docs.
package bpdoc

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/google/blueprint"
)

// foo docs.
func fooFactory() (blueprint.Module, []interface{}) {
	return nil, []interface{}{&props{}}
}

// props docs.
type props struct {
	// A docs.
	A string
}

var pkgPath string
var pkgFiles map[string][]string

func init() {
	pc, filename, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(pc)

	var err error
	pkgPath, err = funcNameToPkgPath(fn.Name())
	if err != nil {
		panic(err)
	}

	pkgFiles = map[string][]string{
		pkgPath: {filename},
	}
}

func TestModuleTypeDocs(t *testing.T) {
	r := NewReader(pkgFiles)
	mt, err := r.ModuleType("foo_module", reflect.ValueOf(fooFactory))
	if err != nil {
		t.Fatal(err)
	}

	if mt.Text != "foo docs.\n" {
		t.Errorf("unexpected docs %q", mt.Text)
	}

	if mt.PkgPath != pkgPath {
		t.Errorf("expected pkgpath %q, got %q", pkgPath, mt.PkgPath)
	}
}

func TestPropertyStruct(t *testing.T) {
	r := NewReader(pkgFiles)
	ps, err := r.PropertyStruct(pkgPath, "props", reflect.ValueOf(props{A: "B"}))
	if err != nil {
		t.Fatal(err)
	}

	if ps.Text != "props docs.\n" {
		t.Errorf("unexpected docs %q", ps.Text)
	}
	if len(ps.Properties) != 1 {
		t.Fatalf("want 1 property, got %d", len(ps.Properties))
	}

	if ps.Properties[0].Name != "a" || ps.Properties[0].Text != "A docs.\n\n" || ps.Properties[0].Default != "B" {
		t.Errorf("unexpected property docs %q %q %q",
			ps.Properties[0].Name, ps.Properties[0].Text, ps.Properties[0].Default)
	}
}

func TestPackage(t *testing.T) {
	r := NewReader(pkgFiles)
	pkg, err := r.Package(pkgPath)
	if err != nil {
		t.Fatal(err)
	}

	if pkg.Text != "bpdoc docs.\n" {
		t.Errorf("unexpected docs %q", pkg.Text)
	}
}

func TestFuncToPkgPath(t *testing.T) {
	tests := []struct {
		f    string
		want string
	}{
		{
			f:    "github.com/google/blueprint/bootstrap.Main",
			want: "github.com/google/blueprint/bootstrap",
		},
		{
			f:    "android/soong/android.GenruleFactory",
			want: "android/soong/android",
		},
		{
			f:    "android/soong/android.ModuleFactoryAdapter.func1",
			want: "android/soong/android",
		},
		{
			f:    "main.Main",
			want: "main",
		},
	}
	for _, tt := range tests {
		t.Run(tt.f, func(t *testing.T) {
			got, err := funcNameToPkgPath(tt.f)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("funcNameToPkgPath(%v) = %v, want %v", tt.f, got, tt.want)
			}
		})
	}
}
