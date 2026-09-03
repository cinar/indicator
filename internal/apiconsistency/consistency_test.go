// Copyright (c) 2021-2026 The Indicator Authors.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

// Package apiconsistency contains static-analysis tests that guard two
// library-wide conventions that are easy for a future contributor to
// silently regress when adding a new indicator/strategy type:
//
//  1. Every type with a ComputeWithContext method must also implement
//     fmt.Stringer (a String() string method) and have a package-level
//     New<TypeName> constructor function.
//  2. Every such type that is generic over helper.Number, rather than the
//     stricter helper.Float (float32|float64), is heuristically checked for
//     an in-body division of a value of its own generic type parameter -
//     which would silently truncate for integer type arguments.
//
// New types are discovered automatically by parsing the source with
// go/parser and go/ast (deliberately not golang.org/x/tools or go/types, to
// keep the module dependency-free and the check a pure syntax-level match)
// so nothing here needs to be updated when a new indicator is added - only
// when one is added that violates a convention.
package apiconsistency

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
)

// scannedDirs lists the library packages whose types must follow the
// ComputeWithContext conventions. cmd/ and examples/ are intentionally
// excluded - they are demo/CLI code, not library API surface.
var scannedDirs = []string{
	"trend",
	"momentum",
	"volatility",
	"volume",
	"strategy",
	"backtest",
	"asset",
}

// concreteNumericConversions are the builtin numeric conversion
// functions/type names that the division heuristic treats as producing a
// value that is definitely not the generic type parameter, e.g. float64(x)
// or `var x float64`.
var concreteNumericConversions = map[string]bool{
	"float64": true, "float32": true,
	"int": true, "int8": true, "int16": true, "int32": true, "int64": true,
	"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
	"uintptr": true, "byte": true, "rune": true,
}

// repoRoot locates the module root from this test file's own location, so
// the test works regardless of the working directory `go test` is invoked
// from.
func repoRoot(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to determine test file location")
	}

	// This file lives at <root>/internal/apiconsistency/consistency_test.go.
	return filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
}

// typeDecl records where a type was declared.
type typeDecl struct {
	spec *ast.TypeSpec
	file *ast.File
}

// pkgInfo holds the parsed AST data for one scanned package directory.
type pkgInfo struct {
	fset *token.FileSet

	// typeSpecs maps a type name to its declaration.
	typeSpecs map[string]*typeDecl

	// methods maps a type name to its method set, regardless of file.
	methods map[string]map[string]*ast.FuncDecl

	// methodFile records which file a given type's method was declared in.
	methodFile map[string]map[string]*ast.File

	// topFuncs maps a package-level (non-method) function name to its decl.
	topFuncs map[string]*ast.FuncDecl

	// topFuncsByFile maps a file to the package-level functions declared in it.
	topFuncsByFile map[*ast.File]map[string]*ast.FuncDecl

	// methodsByFile maps a file to type -> method -> decl, for methods
	// declared in that file.
	methodsByFile map[*ast.File]map[string]map[string]*ast.FuncDecl
}

// parsePackage parses every non-test .go file directly inside dir (it does
// not recurse - the scanned packages keep all their source at the top level
// and only have testdata/ subdirectories, which contain no .go files).
func parsePackage(t *testing.T, dir string) *pkgInfo {
	t.Helper()

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading dir %s: %v", dir, err)
	}

	info := &pkgInfo{
		fset:           token.NewFileSet(),
		typeSpecs:      map[string]*typeDecl{},
		methods:        map[string]map[string]*ast.FuncDecl{},
		methodFile:     map[string]map[string]*ast.File{},
		topFuncs:       map[string]*ast.FuncDecl{},
		topFuncsByFile: map[*ast.File]map[string]*ast.FuncDecl{},
		methodsByFile:  map[*ast.File]map[string]map[string]*ast.FuncDecl{},
	}

	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}

		path := filepath.Join(dir, name)

		file, err := parser.ParseFile(info.fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatalf("parsing %s: %v", path, err)
		}

		info.topFuncsByFile[file] = map[string]*ast.FuncDecl{}
		info.methodsByFile[file] = map[string]map[string]*ast.FuncDecl{}

		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				if d.Tok != token.TYPE {
					continue
				}
				for _, spec := range d.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						info.typeSpecs[ts.Name.Name] = &typeDecl{spec: ts, file: file}
					}
				}

			case *ast.FuncDecl:
				if d.Recv == nil {
					info.topFuncs[d.Name.Name] = d
					info.topFuncsByFile[file][d.Name.Name] = d
					continue
				}

				typeName, ok := receiverTypeName(d.Recv)
				if !ok {
					continue
				}

				if info.methods[typeName] == nil {
					info.methods[typeName] = map[string]*ast.FuncDecl{}
					info.methodFile[typeName] = map[string]*ast.File{}
				}
				info.methods[typeName][d.Name.Name] = d
				info.methodFile[typeName][d.Name.Name] = file

				if info.methodsByFile[file][typeName] == nil {
					info.methodsByFile[file][typeName] = map[string]*ast.FuncDecl{}
				}
				info.methodsByFile[file][typeName][d.Name.Name] = d
			}
		}
	}

	return info
}

// receiverTypeName extracts the base type name from a method receiver,
// stripping the pointer and any generic type-parameter instantiation, e.g.
// `(e *Ema[T])` -> "Ema".
func receiverTypeName(recv *ast.FieldList) (string, bool) {
	if recv == nil || len(recv.List) != 1 {
		return "", false
	}
	return baseTypeName(recv.List[0].Type)
}

func baseTypeName(expr ast.Expr) (string, bool) {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return baseTypeName(t.X)
	case *ast.IndexExpr:
		return baseTypeName(t.X)
	case *ast.IndexListExpr:
		return baseTypeName(t.X)
	case *ast.Ident:
		return t.Name, true
	default:
		return "", false
	}
}

// computeWithContextTypes returns the names of all types in the package that
// declare a method literally named ComputeWithContext, sorted for stable
// output.
func (p *pkgInfo) computeWithContextTypes() []string {
	var names []string
	for typeName, ms := range p.methods {
		if _, ok := ms["ComputeWithContext"]; ok {
			names = append(names, typeName)
		}
	}
	sort.Strings(names)
	return names
}

// TestComputeWithContextTypesFollowNamingConventions is the fully mechanical
// half of the check (part b of the task): every type with a
// ComputeWithContext method must also have a String() method and a
// New<TypeName> constructor function in the same package.
func TestComputeWithContextTypesFollowNamingConventions(t *testing.T) {
	root := repoRoot(t)

	for _, rel := range scannedDirs {
		rel := rel
		t.Run(rel, func(t *testing.T) {
			dir := filepath.Join(root, rel)
			pkg := parsePackage(t, dir)

			for _, typeName := range pkg.computeWithContextTypes() {
				typeName := typeName
				t.Run(typeName, func(t *testing.T) {
					file := "<unknown file>"
					if decl := pkg.typeSpecs[typeName]; decl != nil {
						file = pkg.fset.Position(decl.spec.Pos()).Filename
					}

					if _, ok := pkg.methods[typeName]["String"]; !ok {
						t.Errorf("%s: type %s has a ComputeWithContext method but no String() method; "+
							"every ComputeWithContext type must implement fmt.Stringer", file, typeName)
					}

					ctorName := "New" + typeName
					if _, ok := pkg.topFuncs[ctorName]; !ok {
						t.Errorf("%s: type %s has a ComputeWithContext method but no constructor function "+
							"named %s in package %s", file, typeName, ctorName, rel)
					}
				})
			}
		})
	}
}

// TestComputeWithContextTypesUseFloatWhenDividing is the heuristic half of
// the check (part a of the task). It is inherently approximate - "does this
// type divide" is a judgement call at the syntax level - so it is
// deliberately scoped to reduce noise:
//
//   - It only looks at types that are generic over helper.Number (the loose
//     constraint that also permits integers); types already using
//     helper.Float, or non-generic types, have nothing to flag.
//   - It only scans the ComputeWithContext method body, plus the bodies of
//     same-type methods and package-level helper functions it calls that are
//     declared in the *same file* (matching the task's own scoping - it does
//     not cross file or package boundaries, so calls into e.g. the helper
//     package, which is itself a set of generic Number-parameterized
//     utilities, deliberately are not traced into).
//   - A "/" division is only flagged if neither operand can be shown to be a
//     concrete (non-generic) value - a numeric literal, an explicit
//     conversion to a builtin numeric type (float64(x), int(x), ...), a
//     math.* call result, or a local variable traceable back to one of
//     those. Anything else (a bare identifier, a struct field, a channel
//     value) is conservatively treated as possibly being the generic type.
//
// This will have false negatives (division hidden behind a cross-file or
// cross-package helper call) and could in principle have false positives
// (a local variable whose concreteness the heuristic fails to trace) - that
// is acceptable; the point is an early-warning signal for the exact pattern
// the "helper.Number vs helper.Float" cleanup fixed, not a proof.
func TestComputeWithContextTypesUseFloatWhenDividing(t *testing.T) {
	root := repoRoot(t)

	for _, rel := range scannedDirs {
		rel := rel
		t.Run(rel, func(t *testing.T) {
			dir := filepath.Join(root, rel)
			pkg := parsePackage(t, dir)

			for _, typeName := range pkg.computeWithContextTypes() {
				decl, ok := pkg.typeSpecs[typeName]
				if !ok || decl.spec.TypeParams == nil || len(decl.spec.TypeParams.List) == 0 {
					// Not generic over a Number/Float-style type parameter.
					continue
				}

				constraint, typeParamName := classifyConstraint(decl.spec.TypeParams)
				if constraint != "Number" {
					// Already helper.Float, or a constraint this heuristic
					// does not recognize - nothing to flag either way.
					continue
				}

				compute := pkg.methods[typeName]["ComputeWithContext"]
				if compute == nil || compute.Body == nil {
					continue
				}
				computeFile := pkg.methodFile[typeName]["ComputeWithContext"]

				bodies := collectSameFileCallGraph(pkg, computeFile, typeName, compute)

				positions := findGenericDivision(bodies, typeParamName)
				if len(positions) == 0 {
					continue
				}

				file := pkg.fset.Position(decl.spec.Pos()).Filename
				var locs []string
				for _, pos := range positions {
					locs = append(locs, pkg.fset.Position(pos).String())
				}

				t.Errorf("%s: type %s[%s helper.Number] appears to divide a value of its own generic type "+
					"parameter in ComputeWithContext (at %s); if so it should be constrained to helper.Float "+
					"instead of helper.Number, otherwise an integer type argument would silently truncate",
					file, typeName, typeParamName, strings.Join(locs, ", "))
			}
		})
	}
}

// classifyConstraint returns the constraint name ("Number", "Float", or
// "unknown") and the type parameter's identifier (e.g. "T") for the first
// entry of a type parameter list.
func classifyConstraint(tp *ast.FieldList) (constraint, paramName string) {
	field := tp.List[0]
	if len(field.Names) == 0 {
		return "unknown", ""
	}
	paramName = field.Names[0].Name

	switch t := field.Type.(type) {
	case *ast.SelectorExpr:
		return t.Sel.Name, paramName
	case *ast.Ident:
		return t.Name, paramName
	default:
		return "unknown", paramName
	}
}

// collectSameFileCallGraph returns the ComputeWithContext method body, plus
// the bodies of same-type methods and package-level functions it
// (transitively) calls that are declared in the same file. It intentionally
// does not follow calls into other files or other packages.
func collectSameFileCallGraph(pkg *pkgInfo, file *ast.File, typeName string, root *ast.FuncDecl) []ast.Node {
	visited := map[string]bool{"method:ComputeWithContext": true}
	var bodies []ast.Node
	queue := []*ast.FuncDecl{root}

	fileMethods := pkg.methodsByFile[file][typeName]
	fileFuncs := pkg.topFuncsByFile[file]

	for len(queue) > 0 {
		fd := queue[0]
		queue = queue[1:]
		if fd.Body == nil {
			continue
		}
		bodies = append(bodies, fd.Body)

		ast.Inspect(fd.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			switch fun := unwrapIndex(call.Fun).(type) {
			case *ast.Ident:
				if callee, ok := fileFuncs[fun.Name]; ok && !visited["func:"+fun.Name] {
					visited["func:"+fun.Name] = true
					queue = append(queue, callee)
				}
			case *ast.SelectorExpr:
				if callee, ok := fileMethods[fun.Sel.Name]; ok && !visited["method:"+fun.Sel.Name] {
					visited["method:"+fun.Sel.Name] = true
					queue = append(queue, callee)
				}
			}
			return true
		})
	}

	return bodies
}

// unwrapIndex strips generic instantiation ("[T]" / "[T1, T2]") wrappers off
// a call target expression.
func unwrapIndex(expr ast.Expr) ast.Expr {
	switch t := expr.(type) {
	case *ast.IndexExpr:
		return unwrapIndex(t.X)
	case *ast.IndexListExpr:
		return unwrapIndex(t.X)
	default:
		return expr
	}
}

// findGenericDivision returns the positions of "/" binary expressions within
// bodies where neither operand could be shown to be concrete (see the
// isConcreteExpr documentation).
func findGenericDivision(bodies []ast.Node, _ string) []token.Pos {
	concrete := map[string]bool{}

	// Fixed-point pass over local variable declarations/assignments: a
	// handful of iterations is enough for the straight-line/simple-loop
	// style code this codebase uses.
	for i := 0; i < 5; i++ {
		changed := false
		for _, b := range bodies {
			ast.Inspect(b, func(n ast.Node) bool {
				switch s := n.(type) {
				case *ast.ValueSpec:
					if isConcreteTypeExpr(s.Type) {
						for _, name := range s.Names {
							if !concrete[name.Name] {
								concrete[name.Name] = true
								changed = true
							}
						}
					} else if len(s.Values) == len(s.Names) {
						for i, v := range s.Values {
							if isConcreteExpr(v, concrete) && !concrete[s.Names[i].Name] {
								concrete[s.Names[i].Name] = true
								changed = true
							}
						}
					}

				case *ast.AssignStmt:
					if len(s.Lhs) == len(s.Rhs) {
						for i, lhs := range s.Lhs {
							id, ok := lhs.(*ast.Ident)
							if !ok {
								continue
							}
							if isConcreteExpr(s.Rhs[i], concrete) && !concrete[id.Name] {
								concrete[id.Name] = true
								changed = true
							}
						}
					}
				}
				return true
			})
		}
		if !changed {
			break
		}
	}

	var positions []token.Pos
	for _, b := range bodies {
		ast.Inspect(b, func(n ast.Node) bool {
			bin, ok := n.(*ast.BinaryExpr)
			if !ok || bin.Op != token.QUO {
				return true
			}
			if !isConcreteExpr(bin.X, concrete) || !isConcreteExpr(bin.Y, concrete) {
				positions = append(positions, bin.Pos())
			}
			return true
		})
	}
	return positions
}

// isConcreteTypeExpr reports whether a declared variable type (e.g. in
// `var before float64`) is one of the builtin numeric types.
func isConcreteTypeExpr(expr ast.Expr) bool {
	id, ok := expr.(*ast.Ident)
	return ok && concreteNumericConversions[id.Name]
}

// isConcreteExpr conservatively reports whether expr can be shown, from
// syntax alone, to produce a value that is not the generic type parameter:
// a literal, an explicit conversion to (or arithmetic over) a builtin
// numeric type, a math.* call, or a local variable already known to be one
// of those. Anything else - a bare identifier, struct field, index
// expression, channel receive, etc. - is treated as possibly being the
// generic type, erring on the side of flagging rather than missing a real
// issue.
func isConcreteExpr(expr ast.Expr, concrete map[string]bool) bool {
	switch e := expr.(type) {
	case nil:
		return true
	case *ast.ParenExpr:
		return isConcreteExpr(e.X, concrete)
	case *ast.BasicLit:
		return true
	case *ast.Ident:
		return concrete[e.Name]
	case *ast.CallExpr:
		switch fun := unwrapIndex(e.Fun).(type) {
		case *ast.Ident:
			return concreteNumericConversions[fun.Name]
		case *ast.SelectorExpr:
			id, ok := fun.X.(*ast.Ident)
			return ok && id.Name == "math"
		default:
			return false
		}
	case *ast.BinaryExpr:
		return isConcreteExpr(e.X, concrete) && isConcreteExpr(e.Y, concrete)
	case *ast.UnaryExpr:
		return isConcreteExpr(e.X, concrete)
	default:
		return false
	}
}
