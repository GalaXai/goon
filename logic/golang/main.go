package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>My HTMX Test</title>
        <link rel="stylesheet" href="/static/style.css">
        <script src="https://unpkg.com/htmx.org@2.0.0"></script>
    </head>
    <body>
        <div class="container">
        `)
	/*
		Request to get topbar
		topbar(w, r)
		Request to get about me
		getAboutMe(w, r) # INSIPIRATON for now . vin.gg
	*/

	// Request to get Projects
	fmt.Fprintf(w, `
		<h1> projects </h1>
		<div id="projects-container" hx-trigger="load" hx-get="/projects?index=0" hx-target="this">
            Loading projects...
        </div>`)

	// Request get Tools
	getTools(w, r)

	fmt.Fprintf(w, `
           </div>
    </body>
    </html>
    `)
}

func getProjects(w http.ResponseWriter, r *http.Request) {
	projects := []struct {
		Name        string
		Description string
	}{
		{"personal page", "personal page build in GO and HTMX as a lightweight and fast webapp"},
		{"chat-go", "local llm inference with clausie artificats build and GO and htmx + Python for llms"},
		{"tiny-ai", "python project build during a fastai 2022 course part1&2, that allows for training a inference of nn's and genai models"},
		{"tiny-test", "a TODO project that uses tinygrad for something to learn it."},
	}
	index, _ := strconv.Atoi(r.URL.Query().Get("index"))
	if index >= len(projects) {
		index = 0
	}

	fmt.Fprintf(w, `
		<div class=project-section">
		<div id="project-list" hx-target="this" hx-swap="innerHTML">
		`)

	renderProjects(w, projects, index)
}

func renderProjects(w http.ResponseWriter, projects []struct{ Name, Description string }, startIndex int) {
	endIndex := startIndex + 3
	if endIndex > len(projects) {
		endIndex = len(projects)
	}

	prevIndex := (startIndex - 3 + len(projects)) % len(projects)
	nextIndex := endIndex % len(projects)

	fmt.Fprintf(w, `
		<div class="projects-grid">
			<div class="nav-button" hx-get="/projects?start=%d" hx-target="#projects-container">&lt;</div>
			<div class="project-container">
	`, prevIndex)

	for i := startIndex; i < endIndex; i++ {
		project := projects[i]
		fmt.Fprintf(w, `
			<div class="project">
				<img src="static/placeholder-tool.png" alt="%s">
				<div class="tool-description">
					<h3>%s</h3>
					<p>%s</p>
				</div>
			</div>
		`, project.Name, project.Name, project.Description)
	}

	fmt.Fprintf(w, `
			</div>
			<div class="nav-button" hx-get="/projects?start=%d" hx-target="#projects-container">&gt;</div>
		</div>
	`, nextIndex)
}

func getTools(w http.ResponseWriter, r *http.Request) {
	// Maybe pass tools as arguments?
	tools := []struct {
		Name        string
		Description string
	}{
		{"Go", "Efficient and fast programming language"},
		{"HTMX", "HTML-based AJAX for modern web apps"},
		{"SCSS", "CSS preprocessor for enhanced styling"},
		{"Git", "Version control system"},
		{"Docker", "Containerization platform"},
		{"Python", "Python"},
	}

	fmt.Fprintf(w, `
	<h2>tools i use</h2>
    <div class="tools-section">
        <h1>here are some of the tools and technologies i work with:</h1>
        <div class="tools-grid">`)

	for _, tool := range tools {
		fmt.Fprintf(w, `
        	<div class="tool-item">
                <img src="/static/placeholder-tool.png" alt="%s">
                <div class="tool-description"> %s </div>
            </div>`, tool.Name, tool.Description)
	}
	// Close the Div
	fmt.Fprintf(w, `
	</div></div>`)
}

func getAboutMe(w http.ResponseWriter, r *http.Request) {

}

func Bmain() {
	// Serve static files
	fileServer := http.FileServer(http.Dir("../../static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Handle the home page
	http.HandleFunc("/", homePage)
	http.HandleFunc("/projects", getProjects)

	port := ":8080"
	url := fmt.Sprintf("http://localhost%s", port)
	fmt.Printf("Server starting on %s\n", url)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}

func main() {
	server := NewAPIServer(":3000")
	server.Run()
}
