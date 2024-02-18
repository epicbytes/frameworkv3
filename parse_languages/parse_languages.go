package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go/ast"
	"go/parser"
	"go/token"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatal(errors.New("must specify exactly tree arguments: templatesPath, localizationsPath, languages(en,fr,zh...)"))
	}
	if collectedLanguages, err := ExtractTemplatesDirectory(os.Args[1], "_templ.go"); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(collectedLanguages)
		fmt.Println(ExtractLocalizationsDirectory(os.Args[2]))
	}
}

type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}

type Language = map[string][]string

type file struct {
	fset         *token.FileSet
	astFile      *ast.File
	src          []byte
	filename     string
	main         bool
	translations Language
}

func (f *file) walk(fn func(ast.Node) bool) {
	ast.Walk(walker(fn), f.astFile)
}

func (f *file) find() Language {
	f.findTranslationCalls()
	return f.translations
}

func isPkgDot(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
}

func isIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}

func (f *file) findTranslationCalls() {
	f.walk(func(node ast.Node) bool {
		ce, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		isLibT := isPkgDot(ce.Fun, "lib", "T")
		if !isLibT {
			return true
		}
		if fun, ok := ce.Fun.(*ast.SelectorExpr); ok {
			switch fun.Sel.Name {
			case "T":
				var scope = ""
				var name = ""
				if ar, ok := ce.Args[1].(*ast.BasicLit); ok {
					scope = ar.Value
				}
				if ar, ok := ce.Args[2].(*ast.BasicLit); ok {
					name = ar.Value
				}
				f.translations[scope] = append(f.translations[scope], name)
			default:
				return true
			}
		}
		return true
	})
}

func ExtractTemplatesDirectory(dir string, filenameSuffix string) (Language, error) {
	all := make(Language)
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, filenameSuffix) {
				f, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				lang, err := CollectTranslations(info.Name(), f)
				if err != nil {
					return err
				}
				if len(lang) != 0 {
					all = mergeLanguages(all, lang)
				}
				return nil
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return all, nil
}

func ExtractLocalizationsDirectory(dir string) (Language, error) {
	all := make(Language)
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, ".locale.yml") {
				/*f, err := os.ReadFile(path)
				if err != nil {
					return err
				}*/
				/*lang, err := CollectTranslationsFromLocalizations(info.Name(), f)
				if err != nil {
					return err
				}
				if len(lang) != 0 {
					all = mergeLocLanguages(all, lang)
				}*/
				return nil
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return all, nil
}

func CollectTranslations(filename string, content []byte) (Language, error) {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	f := &file{fset: fset, astFile: astFile, src: content, filename: filename, translations: make(Language)}
	return f.find(), nil
}

func CollectTranslationsFromLocalizations(filename string, content []byte) (map[string]interface{}, error) {
	var config map[string]interface{}
	err := yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func mergeLanguages(all Language, lang Language) Language {
	for l1, v1 := range lang {
		if _, ok := all[l1]; ok {
			all[l1] = lo.Uniq(lo.Union(all[l1], v1))
		} else {
			all[l1] = lo.Uniq(v1)
		}
	}

	return all
}

/*func mergeLocLanguages(all map[string]interface{}, lang map[string]interface{}) map[string]interface{} {
	for l1, v1 := range lang {
		if _, ok := all[l1]; ok {
			all[l1] = lo.Uniq(lo.Union(all[l1], v1))
		} else {
			all[l1] = lo.Uniq(v1)
		}
	}

	return all
}
*/
