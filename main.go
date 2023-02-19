package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/urfave/cli/v2"
)

const defaultObsAddress = "localhost:4455"
const defaultObsPassword = "goodpassword"

func createObsClient() *goobs.Client {
	// Setup obs-websockets client
	client, err := goobs.New(defaultObsAddress, goobs.WithPassword(defaultObsPassword))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func getScenes(client *goobs.Client) []string {
	sceneNames := []string{}

	resp, err := client.Scenes.GetSceneList()
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range resp.Scenes {
		sceneNames = append(sceneNames, v.SceneName)
	}
	return sceneNames
}

func sceneExists(client *goobs.Client, sceneName string) bool {
	sceneNames := getScenes(client)
	for _, v := range sceneNames {
		if v == sceneName {
			return true
		}
	}
	return false
}

func ensureSceneExists(client *goobs.Client, sceneName string) {
	if !sceneExists(client, sceneName) {
		log.Fatalf("Scene %s does not exist", sceneName)
	}
}

func main() {

	client := createObsClient()
	defer client.Disconnect()

	app := &cli.App{
		Name:        "obs-cli",
		Description: "A command line interface to drive Streamdeck actions related to Open Broadcaster Studio (OBS)",
		Commands: []*cli.Command{
			{
				Name: "scene",
				Subcommands: []*cli.Command{
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "list defined obs scenes",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("Listing obs scenes..")
							sceneNames := getScenes(client)
							for _, sceneName := range sceneNames {
								fmt.Println(sceneName)
							}
							return nil
						},
					},
					{
						Name:    "change",
						Aliases: []string{"c"},
						Usage:   "switch to a different scene",
						Action: func(cCtx *cli.Context) error {
							selectedScene := cCtx.String("scene-name")
							fmt.Printf("Switching to scene %s\n", selectedScene)

							ensureSceneExists(client, selectedScene)

							params := &scenes.SetCurrentProgramSceneParams{
								SceneName: selectedScene,
							}
							_, switchSceneErr := client.Scenes.SetCurrentProgramScene(params)
							if switchSceneErr != nil {
								fmt.Println("Error switching scene: ", switchSceneErr)
							}

							return nil
						},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "scene-name",
								Usage:    "name of the scene to switch to",
								Required: true,
							},
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
