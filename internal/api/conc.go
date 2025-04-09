package api

import (
	"log"
	"sync"
)

const WORKER_COUNT = 5

func registerHook(projects []string, hookId string) {
	projectChan := make(chan string, len(projects))
	var wg sync.WaitGroup

	for i := 0; i < WORKER_COUNT; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for project := range projectChan {
				if err := registerHookForProject(project, hookId); err != nil {
					log.Printf("Error registering hook for project %s: %v", project, err)
				}
			}
		}()
	}
	for _, project := range projects {
		projectChan <- project
	}
	close(projectChan)
	wg.Wait()
	log.Printf("All hooks registered for projects: %v", projects)
}

func registerHookForProject(project string, hookId string) error {
	// Simulate the registration process
	log.Printf("Registering hook %s for project %s", hookId, project)
	return nil // Replace with actual registration logic
}
func main() {
	projects := []string{"project1", "project2", "project3"}
	hookId := "hook123"
	registerHook(projects, hookId)
}

// This is a placeholder for the actual registration logic.
