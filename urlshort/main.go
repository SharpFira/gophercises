package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type PageData struct {
	InMemoryMap map[string]string
}

type SimpleHandler struct {
	inMemMap map[string]string
}

type YamlHandler struct {
	yamlPath string
}

type RedirectPath struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func (handler SimpleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling redirect...")

	urlKey := r.PathValue("key")

	url, exists := handler.inMemMap[urlKey]

	if !exists {
		http.Error(w, "We did not find the redirect...", http.StatusNotFound)
		return
	}

	fmt.Println("Redirect to: ", url)

	http.Redirect(w, r, url, http.StatusFound)
}

func (handler YamlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Invoking yaml handler...")

	data, err := os.ReadFile(handler.yamlPath)

	if err != nil {
		log.Fatal("Could not read yaml file: ", handler.yamlPath)
	}

	var paths []RedirectPath

	yamlParsingError := yaml.Unmarshal(data, &paths)

	if yamlParsingError != nil {
		log.Fatal("Could not parse yaml file: ", yamlParsingError.Error())
	}

	urlKey := strings.TrimPrefix(r.URL.Path, "/yaml")
	fmt.Println("KEY IS: ", urlKey)

	for _, v := range paths {
		fmt.Println("PATH ARE: ", v.Path)
		if v.Path == urlKey {
			http.Redirect(w, r, v.URL, http.StatusFound)
			return
		}
	}

	http.Error(w, "Could not find key in yaml", http.StatusBadRequest)
}

func main() {
	templateData := PageData{
		InMemoryMap: make(map[string]string),
	}

	template := template.Must(template.ParseFiles("./index.html"))

	redirectHandler := SimpleHandler{
		inMemMap: templateData.InMemoryMap,
	}

	http.Handle("/redirect/{key}", redirectHandler)

	yamlHandler := YamlHandler{
		yamlPath: "./redirect.yaml",
	}

	http.Handle("/yaml/{key}", yamlHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var errorMessage string

		if r.Method == "GET" {
			template.Execute(w, templateData)
			return
		}

		if r.Method == "POST" {
			fmt.Println("Recieved Form Content")

			urlKey := r.FormValue("key")
			urlValue := r.FormValue("url")

			if len(urlKey) == 0 {
				errorMessage = "key value is empty"
				log.Println(errorMessage)
				http.Error(w, errorMessage, http.StatusBadRequest)
				return
			}

			if len(urlValue) == 0 {
				errorMessage = "url value is empty"
				log.Println(errorMessage)
				http.Error(w, errorMessage, http.StatusBadRequest)
				return
			}

			fmt.Printf("Recieved values key: %s and url: %s \n", urlKey, urlValue)

			templateData.InMemoryMap[urlKey] = urlValue

			log.Println(len(templateData.InMemoryMap))

			template.Execute(w, templateData)
			return
		}
	})

	fmt.Println("Running on port: http://localhost:8080")

	http.ListenAndServe(":8080", nil)

}
