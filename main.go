package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/inputs"
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

func getInputVolumeInDB(client *goobs.Client, inputName string) float64 {
	// First, determine what volume the input is currently set to
	getParams := &inputs.GetInputVolumeParams{
		InputName: inputName,
	}
	getResp, err := client.Inputs.GetInputVolume(getParams)
	if err != nil {
		log.Fatal(err)
	}
	return getResp.InputVolumeDb
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
				Name: "inputs",
				Subcommands: []*cli.Command{
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "list defined obs inputs",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("Listing obs inputs..")
							resp, err := client.Inputs.GetInputList()
							if err != nil {
								log.Fatal(err)
							}
							for _, v := range resp.Inputs {
								fmt.Println(v)
							}
							return nil
						},
					},
					{
						Name: "lower",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "input-name",
								Aliases:  []string{"i"},
								Usage:    "name of the input to lower",
								Required: true,
							},
						},
						Action: func(cCtx *cli.Context) error {

							inputName := cCtx.String("input-name")
							fmt.Printf("Lowering input %s\n", inputName)

							currentVolDb := getInputVolumeInDB(client, inputName)

							fmt.Printf("Current volume: %f\n", currentVolDb)

							//If it's too low, do nothing
							if currentVolDb <= -50 {
								fmt.Println("Volume is already at minimum")
								return nil
							}

							params := &inputs.SetInputVolumeParams{
								InputName:     inputName,
								InputVolumeDb: currentVolDb - 2,
							}
							_, err := client.Inputs.SetInputVolume(params)
							if err != nil {
								log.Fatal(err)
							}
							return nil
						},
					},
					{
						Name: "raise",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "input-name",
								Aliases:  []string{"i"},
								Usage:    "name of the input to raise",
								Required: true,
							},
						},
						Action: func(cCtx *cli.Context) error {

							inputName := cCtx.String("input-name")

							fmt.Printf("Raising input %s\n", inputName)

							currentVolDb := getInputVolumeInDB(client, inputName)

							fmt.Printf("Current volume: %f\n", currentVolDb)

							//If it's too high, do nothing
							if currentVolDb >= 26 {
								fmt.Println("Volume is already at max")
								return nil
							}

							params := &inputs.SetInputVolumeParams{
								InputName:     inputName,
								InputVolumeDb: currentVolDb + 2,
							}
							_, setVolErr := client.Inputs.SetInputVolume(params)
							if setVolErr != nil {
								log.Fatal(setVolErr)
							}
							return nil
						},
					},
				},
			},
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
