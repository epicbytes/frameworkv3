package main

import (
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatal(errors.New("must specify exactly tree arguments: templatesPath, localizationsPath, languages(en,fr,zh...)"))
	}
	if collectedLanguages, err := ExtractTemplatesDirectory(os.Args[1], strings.Split(os.Args[3], ","), "_templ.go"); err != nil {
		log.Fatal(err)
	} else {
		collectedLocalizations, err := ExtractLocalizationsDirectory(os.Args[2])
		if err != nil {
			log.Println(err)
		}

		ready := AssignLanguages(collectedLocalizations, collectedLanguages)
		for lang, data := range ready {
			d, err := yaml.Marshal(Language{lang: data})
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			f, err := os.Create(os.Args[2] + "/" + lang + ".locale.yml")
			if err != nil {
				panic(err)
			}
			defer f.Close()

			_, err = io.WriteString(f, string(d))
			if err != nil {
				panic(err)
			}
		}
	}
}

type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}

type Language = map[string]map[string]map[string]string

type file struct {
	fset         *token.FileSet
	astFile      *ast.File
	src          []byte
	filename     string
	main         bool
	translations Language
	langs        []string
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
				var scope string
				var name string
				if ar, ok := ce.Args[1].(*ast.BasicLit); ok {
					scope, _ = strconv.Unquote(ar.Value)
				}
				if ar, ok := ce.Args[2].(*ast.BasicLit); ok {
					name, _ = strconv.Unquote(ar.Value)
				}
				for _, l := range f.langs {
					_, ok := f.translations[l][scope]
					if !ok {
						f.translations[l][scope] = make(map[string]string)
					}
					f.translations[l][scope][name] = ""
				}
			default:
				return true
			}
		}
		return true
	})
}

func AssignLanguages(maps ...Language) Language {
	out := make(Language)

	for _, m := range maps {
		for k, lngs := range m {
			if len(lngs) > 0 {
				if _, ok := out[k]; !ok {
					out[k] = map[string]map[string]string{}
				}
				for k2, v2 := range lngs {
					if _, ok := out[k][k2]; !ok {
						out[k][k2] = map[string]string{}
					}
					if len(v2) > 0 {
						for k3, v3 := range v2 {
							val, ok := out[k][k2][k3]
							if ok {
								out[k][k2][k3] = val
							} else {
								out[k][k2][k3] = v3
							}
						}
					}
				}
			}
		}
	}
	return out
}

func ExtractTemplatesDirectory(dir string, langs []string, filenameSuffix string) (Language, error) {
	all := make(Language)
	for _, l := range langs {
		all[l] = make(map[string]map[string]string)
	}
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
				lang, err := CollectTranslations(info.Name(), langs, f)
				if err != nil {
					return err
				}
				all = AssignLanguages(all, lang)
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
			if strings.HasSuffix(path, "locale.yml") {
				f, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				lang, err := CollectTranslationsFromLocalizations(info.Name(), f)
				if err != nil {
					return err
				}
				if len(lang) != 0 {
					all = AssignLanguages(all, lang)
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

func CollectTranslations(filename string, langs []string, content []byte) (Language, error) {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var translations = make(Language)

	for _, l := range langs {
		translations[l] = make(map[string]map[string]string)
	}

	f := &file{fset: fset, astFile: astFile, src: content, filename: filename, langs: langs, translations: translations}
	return f.find(), nil
}

func CollectTranslationsFromLocalizations(filename string, content []byte) (Language, error) {
	var config Language
	err := yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
