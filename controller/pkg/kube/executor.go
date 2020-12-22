package kube

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"html/template"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"regexp"

	"k8s.io/client-go/kubernetes/scheme"
)

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

// ParseTemplates - parses and exexute templates and return map of strings with kubeobj data
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

func GenerateSingleManifest(parsedTemplates map[string]string) string {
	singleTemplate := ""
	for _, tpl := range parsedTemplates {
		singleTemplate += tpl
		singleTemplate += "\n---\n"
	}
	return singleTemplate
}

// KubeObjectsFromTemplates return map of runtime.Objects from templateMap
// see https://github.com/kubernetes/client-go/issues/193 for examples
func BuildObjectsFromTemplates(templatesMap map[string]string, data interface{}) (map[string]runtime.Object, map[string]string, error) {
	parsedTemplates, err := ParseTemplates(templatesMap, data)
	if err != nil {
		return nil, nil, err
	}

	_ = apiextv1beta1.AddToScheme(scheme.Scheme)

	// Deserializing all kube objects from parsedTemplates
	// see https://github.com/kubernetes/client-go/issues/193 for examples
	kubeDecode := scheme.Codecs.UniversalDeserializer().Decode
	kubeObjects := make(map[string]runtime.Object)
	for n, objStr := range parsedTemplates {
		obj, _, err := kubeDecode([]byte(objStr), nil, nil)
		if err != nil {
			return nil, parsedTemplates, err
		}
		kubeObjects[n] = obj
	}
	return kubeObjects, parsedTemplates, nil
}

//func GetKubeObjectsFromTemplate(values map[string]interface{}) (map[string]runtime.Object, map[string]string, error) {
//templatesMap := kubernetes.TemplatesMap()
//return KubeObjectsFromTemplates(templatesMap, values)
//}
