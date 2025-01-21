package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type PrivateKeys struct {
	PrivateKeys []PrivateKey `json:"private_keys"`
}

type PrivateKey struct {
	ID         int    `json:"id"`
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

var privateKeysCmd = &cobra.Command{
	Use:   "private-keys",
	Short: "Private key related commands",
}

var listPrivateKeysCmd = &cobra.Command{
	Use:   "list",
	Short: "List all private keys",
	Run: func(cmd *cobra.Command, args []string) {
		CheckDefaultThings(nil)

		baseUrl := "security/keys"
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
		var jsondata []PrivateKey
		err = json.Unmarshal([]byte(data), &jsondata)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Fprintln(w, "Uuid\tName")
		for _, resource := range jsondata {
			if ShowSensitive {
				fmt.Fprintf(w, "%s\t%s\n", resource.UUID, resource.Name)
			} else {
				fmt.Fprintf(w, "%s\t%s\n", resource.UUID, resource.Name)
			}
		}
		w.Flush()
		fmt.Println("\nNote: Use -s to show sensitive information.")
	},
}
var onePrivateKeyCmd = &cobra.Command{
	Use:   "get [uuid]",
	Args:  cobra.ExactArgs(1),
	Short: "Get private key details by uuid",
	Run: func(cmd *cobra.Command, args []string) {
		CheckDefaultThings(nil)
		baseUrl := "security/keys/"

		uuid := args[0]
		var url = baseUrl + uuid

		data, err := Fetch(url)
		if err != nil {
			fmt.Println(err)
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
		var jsondata PrivateKey
		err = json.Unmarshal([]byte(data), &jsondata)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintln(w, "Uuid\tName\tPublicKey\tPrivateKey")
		if ShowSensitive {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", jsondata.UUID, jsondata.Name, jsondata.PublicKey, strings.ReplaceAll(jsondata.PrivateKey, "\n", "\\n"))

		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", jsondata.UUID, jsondata.Name, SensitiveInformationOverlay, SensitiveInformationOverlay)
		}
		w.Flush()
		fmt.Println("\nNote: Use -s to show sensitive information.")

	},
}

var addPrivateKeyCmd = &cobra.Command{
	Use:     "add",
	Example: `add <name> <private_key_or_file>`,
	Args:    cobra.ExactArgs(2),
	Short:   "Add a private key",
	Run: func(cmd *cobra.Command, args []string) {
		version := "4.0.0-beta.383"
		CheckDefaultThings(&version)
		baseUrl := "security/keys"
		name := args[0]
		privateKeyInput := args[1]

		var privateKey string
		// Check if input is a file path
		if _, err := os.Stat(privateKeyInput); err == nil {
			keyBytes, err := os.ReadFile(privateKeyInput)
			if err != nil {
				fmt.Printf("Error reading private key file: %v\n", err)
				return
			}
			privateKey = string(keyBytes)
		} else {
			privateKey = privateKeyInput
		}

		data := map[string]string{
			"name":        name,
			"private_key": privateKey,
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			return
		}
		_, err = Post(baseUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error adding private key: %v\n", err)
			return
		}
		fmt.Printf("Private key '%s' added successfully\n", name)
	},
}

var removePrivateKeyCmd = &cobra.Command{
	Use:   "remove [uuid]",
	Args:  cobra.ExactArgs(1),
	Short: "Remove a private key",
	Run: func(cmd *cobra.Command, args []string) {
		version := "4.0.0-beta.383"
		CheckDefaultThings(&version)
		baseUrl := "security/keys/"
		uuid := args[0]
		_, err := Delete(baseUrl + uuid)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Private key removed successfully")
	},
}

func init() {
	rootCmd.AddCommand(privateKeysCmd)
	privateKeysCmd.AddCommand(listPrivateKeysCmd)
	privateKeysCmd.AddCommand(onePrivateKeyCmd)
	privateKeysCmd.AddCommand(addPrivateKeyCmd)
	privateKeysCmd.AddCommand(removePrivateKeyCmd)
}
