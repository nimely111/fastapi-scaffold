package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// createDir ensures a directory exists
func createDir(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", path, err)
		os.Exit(1)
	}
}

// createFile creates a file with optional content
func createFile(path, content string) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", path, err)
		os.Exit(1)
	}
	defer file.Close()
	file.WriteString(content)
}

func main() {
	// Parse flags
	projectName := flag.String("name", "fastapi_project", "Project name")
	withAuth := flag.Bool("with-auth", false, "Include authentication module")
	flag.Parse()

	// Define project structure
	baseDir := *projectName
	folders := []string{
		filepath.Join(baseDir, "app"),
		filepath.Join(baseDir, "app", "routes"),
		filepath.Join(baseDir, "app", "models"),
	}

	if *withAuth {
		folders = append(folders, filepath.Join(baseDir, "app", "auth"))
	}

	// Create directories
	for _, folder := range folders {
		createDir(folder)
	}

	// Create essential files
	createFile(filepath.Join(baseDir, "main.py"),
		`
	from fastapi import FastAPI
	app = FastAPI()

	@app.get("/")
	def read_root():
    return {"message": "Hello, FastAPI!"}
`)

	createFile(filepath.Join(baseDir, "app", "routes", "__init__.py"), "")
	createFile(filepath.Join(baseDir, "app", "routes", "users.py"),
		`
	from fastapi import APIRouter
	router = APIRouter()

	@router.get("/users")
	def get_users():
		return {"users": []}
`)

	createFile(filepath.Join(baseDir, "app", "models", "__init__.py"), "")
	createFile(filepath.Join(baseDir, "app", "models", "user.py"),
		`
	from pydantic import BaseModel

class User(BaseModel):
    id: int
    username: str
    email: str
`)

	// Create optional auth module
	if *withAuth {
		createFile(filepath.Join(baseDir, "app", "auth", "__init__.py"), "")
		createFile(filepath.Join(baseDir, "app", "auth", "auth.py"),
			`
	from fastapi import APIRouter

	router = APIRouter()

	@router.post("/login")
	def login():
		return {"message": "Login successful"}
`)
	}

	// Create additional files
	createFile(filepath.Join(baseDir, "requirements.txt"), "fastapi\npydantic\njose\nuvicorn\n")
	createFile(filepath.Join(baseDir, ".gitignore"), "__pycache__/\n.env")
	createFile(filepath.Join(baseDir, "README.md"), fmt.Sprintf("# %s\n\nGenerated using FastAPI Scaffold CLI", *projectName))

	fmt.Println("FastAPI project structure created successfully!")
}
