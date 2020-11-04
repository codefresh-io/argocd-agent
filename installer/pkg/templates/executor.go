package templates

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/codefresh-io/argocd-listener/installer/pkg/fs"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates/kubernetes"
	"html/template"
	"k8s.io/apimachinery/pkg/runtime"
	"regexp"
	"strconv"
	"strings"

	"k8s.io/client-go/kubernetes/scheme"
)


func SplitSingleTemplate(singleTemplate string) []string {
	return strings.Split(singleTemplate, "\n---\n")
}

func ResolveKubeTemplates(ManifestPath string) (error, map[string]string) {
	var kubeTemplates map[string]string
	if ManifestPath != "" {
		err, yml := fs.ReadFile(ManifestPath)
		if err != nil {
			return err, nil
		}
		kubeTemplates = make(map[string]string)
		for n, tpl := range SplitSingleTemplate(yml) {
			kubeTemplates["template_"+strconv.Itoa(n)+".yaml"] = tpl
		}
	} else {
		kubeTemplates = kubernetes.TemplatesMap()
	}
	return nil, kubeTemplates
}

func ExecuteTemplate(tplStr string, data interface{}) (string, error) {

	template, err := template.New("base").Funcs(sprig.FuncMap()).Parse(tplStr)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")
	err = template.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ParseTemplates - parses and exexute templates and return map of strings with obj data
func ParseTemplates(templatesMap map[string]string, data interface{}) (map[string]string, error) {
	parsedTemplates := make(map[string]string)
	nonEmptyParsedTemplateFunc := regexp.MustCompile(`[a-zA-Z0-9]`).MatchString
	for n, tpl := range templatesMap {
		tplEx, err := ExecuteTemplate(tpl, data)
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to parse and execute template %s", n))
			return nil, err
		}

		// we add only non-empty parsedTemplates
		if nonEmptyParsedTemplateFunc(tplEx) {
			parsedTemplates[n] = tplEx
		}

	}
	return parsedTemplates, nil
}

// KubeObjectsFromTemplates return map of runtime.Objects from templateMap
// see https://github.com/kubernetes/client-go/issues/193 for examples
func KubeObjectsFromTemplates(templatesMap map[string]string, data interface{}) (map[string]runtime.Object, error) {
	parsedTemplates, err := ParseTemplates(templatesMap, data)
	if err != nil {
		return nil, err
	}

	// Deserializing all kube objects from parsedTemplates
	// see https://github.com/kubernetes/client-go/issues/193 for examples
	kubeDecode := scheme.Codecs.UniversalDeserializer().Decode
	kubeObjects := make(map[string]runtime.Object)
	for n, objStr := range parsedTemplates {
		obj, _, err := kubeDecode([]byte(objStr), nil, nil)
		if err != nil {
			return nil, err
		}
		kubeObjects[n] = obj
	}
	return kubeObjects, nil
}

func GetKubeObjectsFromTemplate(values map[string]interface{}) (map[string]runtime.Object, error) {
	templatesMap := kubernetes.TemplatesMap()
	return KubeObjectsFromTemplates(templatesMap, values)
}
