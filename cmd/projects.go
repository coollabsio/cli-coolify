package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

type Application struct {
	ID          int    `json:"id"`
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
type Environment struct {
	ID           int           `json:"id"`
	Uuid         string        `json:"uuid"`
	Name         string        `json:"name"`
	CreatedAt    string        `json:"created_at"`
	UpdatedAt    string        `json:"updated_at"`
	Description  *string       `json:"description"`
	Applications []Application `json:"applications"`
}

type Project struct {
	Uuid         string        `json:"uuid"`
	Name         string        `json:"name"`
	Environments []Environment `json:"environments"`
}

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Project related commands",
}

var listProjectsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Run: func(cmd *cobra.Command, args []string) {
		CheckDefaultThings(nil)
		baseUrl := "projects"
		data, err := Fetch(baseUrl)
		if err != nil {
			log.Println(err)
			return
		}
		if PrettyMode {
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, []byte(data), "", "\t")
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(prettyJSON.String()))
			return
		}
		if JsonMode {
			fmt.Println(data)
			return
		}
		var jsondata []Project
		err = json.Unmarshal([]byte(data), &jsondata)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Fprintln(w, "Uuid\tName")
		for _, resource := range jsondata {
			fmt.Fprintf(w, "%s\t%s\n", resource.Uuid, resource.Name)
		}
		w.Flush()
	},
}

var oneProjectCmd = &cobra.Command{
	Use:   "get [uuid]",
	Short: "Get a project by uuid",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		CheckDefaultThings(nil)
		uuid := args[0]
		environment, _ := cmd.Flags().GetString("environment")
		if environment != "" {
			url := "projects/" + uuid + "/" + environment
			data, err := Fetch(url)
			if err != nil {
				log.Println(err)
				return
			}
			if PrettyMode {
				var prettyJSON bytes.Buffer
				err := json.Indent(&prettyJSON, []byte(data), "", "\t")
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println(string(prettyJSON.String()))
				return
			}
			if JsonMode {
				fmt.Println(data)
				return
			}
			var jsondata Environment
			err = json.Unmarshal([]byte(data), &jsondata)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Fprintln(w, "Uuid\tName\tStatus")
			for _, resource := range jsondata.Applications {
				fmt.Fprintf(w, "%s\t%s\t%s\n", resource.Uuid, resource.Name, resource.Status)
			}
			w.Flush()
			return
		}
		data, err := Fetch("projects/" + uuid)
		if err != nil {
			log.Println(err)
			return
		}
		if PrettyMode {
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, []byte(data), "", "\t")
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(prettyJSON.String()))
			return
		}
		if JsonMode {
			fmt.Println(data)
			return
		}
		var jsondata Project
		err = json.Unmarshal([]byte(data), &jsondata)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Fprintln(w, "Uuid\tName\tEnvironments")
		envNames := make([]string, len(jsondata.Environments))
		for i, env := range jsondata.Environments {
			envNames[i] = env.Name + " (" + env.Uuid + ")"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", jsondata.Uuid, jsondata.Name, strings.Join(envNames, ", "))
		w.Flush()
	},
}
var addProjectCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		CheckDefaultThings(nil)
		baseUrl := "projects"
		name := args[0]
		data := map[string]string{
			"name": name,
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		response, err := Post(baseUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println(err)
			return
		}
		msg := map[string]string{}
		json.Unmarshal([]byte(response), &msg)
		fmt.Println("Project added successfully with uuid " + msg["uuid"])
	},
}

var removeProjectCmd = &cobra.Command{
	Use:   "remove [uuid]",
	Short: "Remove a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		CheckDefaultThings(nil)
		baseUrl := "projects/"
		uuid := args[0]
		response, err := Delete(baseUrl + uuid)
		if err != nil {
			fmt.Println(err)
			return
		}
		msg := map[string]string{}
		json.Unmarshal([]byte(response), &msg)
		fmt.Println(msg["message"])
	},
}

func init() {
	rootCmd.AddCommand(projectsCmd)
	projectsCmd.AddCommand(listProjectsCmd)
	oneProjectCmd.Flags().StringP("environment", "e", "", "Environment")
	projectsCmd.AddCommand(oneProjectCmd)

	projectsCmd.AddCommand(addProjectCmd)
	projectsCmd.AddCommand(removeProjectCmd)
}
