package pigae

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Deleplace/programming-idioms/pig"
)

//
// Language names exist in 3 forms : nice, standard, lowercase
// Ex : "C++", "Cpp", "cpp"
//

var mainStreamLangs = [...]string{"C", "Cpp", "Csharp", "Go", "Java", "JS", "Obj-C", "PHP", "Python", "Ruby"}

// Return alpha codes for each language (no encoding problems).
// See printNiceLang to display them more fancy.
func mainStreamLanguages() []string {
	return mainStreamLangs[:]
}

var moreLangs = [...]string{"Ada", "Caml", "Clojure", "Cobol", "D", "Dart", "Elixir", "Erlang", "Fortran", "Haskell", "Lua", "Lisp", "Pascal", "Perl", "Prolog", "Rust", "Scala", "Scheme", "VB"}

func moreLanguages() []string {
	// These do *not* include the mainStreamLanguages()
	return moreLangs[:]
}

var synonymLangs = map[string]string{
	"Javascript":   "JS",
	"Objective C":  "Obj-C",
	"Visual Basic": "VB",
}

var allLangs []string
var allNiceLangs []string

func allLanguages() []string {
	if allLangs == nil {
		mainstream := mainStreamLanguages()
		more := moreLanguages()
		allLangs = make([]string, len(mainstream)+len(more))
		copy(allLangs, mainstream)
		copy(allLangs[len(mainstream):], more)
		sort.Strings(allLangs)
		allNiceLangs = make([]string, len(allLangs))
		for i, lg := range allLangs {
			allNiceLangs[i] = printNiceLang(lg)
		}
	}
	return allLangs
}

// autocompletions is a map[string][]string
var autocompletions = precomputeAutocompletions()

func languageAutoComplete(fragment string) []string {
	fragment = strings.ToLower(fragment)

	// Dynamic search (slow)
	// options := []string{}
	// for _, lg := range allLanguages() {
	// 	if strings.Contains(strings.ToLower(lg), fragment) || strings.Contains(strings.ToLower(printNiceLang(lg)), fragment) {
	// 		options = append(options, printNiceLang(lg))
	// 	}
	// }
	// return options

	// Precomputed search (fast)
	return autocompletions[fragment]
}

func printNiceLang(lang string) string {
	switch strings.TrimSpace(strings.ToLower(lang)) {
	case "cpp":
		return "C++"
	case "csharp":
		return "C#"
	default:
		return lang
	}
}

func printNiceLangs(langs []string) []string {
	nice := make([]string, len(langs))
	for i, lang := range langs {
		nice[i] = printNiceLang(lang)
	}
	return nice
}

func printShortLang(lang string) string {
	switch strings.TrimSpace(strings.ToLower(lang)) {
	case "clojure":
		return "Clj"
	case "cobol":
		return "Co bol"
	case "cpp":
		return "C++"
	case "csharp":
		return "C#"
	case "erlang":
		return "Er lang"
	case "elixir":
		return "Eli xir"
	case "fortran":
		return "For tran"
	case "haskell":
		return "Has kell"
	case "obj-c":
		return "Obj C"
	case "pascal":
		return "Pas"
	case "python":
		return "Py"
	case "scheme":
		return "scm"
	case "prolog":
		return "Pro log"
	default:
		return lang
	}
}

func indexByLowerCase(langs []string) map[string]string {
	m := map[string]string{}
	for _, lg := range langs {
		m[strings.ToLower(lg)] = lg
	}
	return m
}

var langLowerCaseIndex = indexByLowerCase(allLanguages())

func normLang(lang string) string {
	lg := strings.TrimSpace(strings.ToLower(lang))
	switch lg {
	case "c++":
		return "Cpp"
	case "c#":
		return "Csharp"
	case "javascript":
		return "JS"
	case "golang":
		return "Go"
	case "objective c":
		return "Obj-C"
	default:
		return langLowerCaseIndex[lg]
	}
}

func precomputeAutocompletions() map[string][]string {
	m := make(map[string][]string, 100)

	// Crazy shadowing of variable "lg" is allowed in go...
	for _, trueLg := range allLanguages() {
		niceLg := printNiceLang(trueLg)
		for _, lg := range []string{trueLg, niceLg} {
			lg := strings.ToLower(lg)
			fragments := substrings(lg)
			for _, frag := range fragments {
				if !pig.StringSliceContains(m[frag], niceLg) {
					m[frag] = append(m[frag], niceLg)
				}
			}
		}
	}

	for lg, trueLg := range synonymLangs {
		niceLg := printNiceLang(trueLg)
		lg := strings.ToLower(lg)
		fragments := substrings(lg)
		for _, frag := range fragments {
			if !pig.StringSliceContains(m[frag], niceLg) {
				m[frag] = append(m[frag], niceLg)
			}
		}
	}

	fmt.Fprintf(os.Stderr, "---\n")
	return m
}

func substrings(s string) []string {
	L := len(s)
	seen := make(map[string]bool, L*L)
	fragments := make([]string, L*L)
	// This works well for language names with only 1-byte chars, not for any string
	for i := 0; i < L; i++ {
		for j := i + 1; j <= L; j++ {
			frag := s[i:j]
			if seen[frag] {
				continue
			}
			seen[frag] = true
			fragments = append(fragments, frag)
		}
	}
	return fragments
}
