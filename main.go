package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/urfave/cli/v2"
)

func getScenes() {

}

func main() {

	// Setup obs-websockets client
	client, err := goobs.New("localhost:4455", goobs.WithPassword("goodpassword"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()

	app := &cli.App{
		Name: "obs-cli",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"a"},
				Usage:   "list defined obs scenes",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("Listing obs scenes..")
					resp, err := client.Scenes.GetSceneList()
					if err != nil {
						fmt.Println("Error getting scene list: ", err)
					}
					for _, v := range resp.Scenes {
						fmt.Printf("%2d %s\n", v.SceneIndex, v.SceneName)
					}

					return nil
				},
			},
			{
				Name:  "change",
				Usage: "switch to a different scene",
				Action: func(cCtx *cli.Context) error {
					selectedScene := cCtx.String("scene-name")
					fmt.Printf("Switching to scene %s", selectedScene)
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
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
