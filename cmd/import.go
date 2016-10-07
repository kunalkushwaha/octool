package cmd

import (
	"container/list"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/docker/distribution/digest"
	"github.com/heroku/docker-registry-client/registry"
)

// DockerHub is url for dockerhub registry
const DockerHub = "https://registry-1.docker.io"

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:     "import",
	Short:   "Imports container image from remote registery and convert it to runc's rootfs",
	Run:     importImage,
	Example: "octool import docker://<repository-name>/<image-name>:<tag>",
}

func init() {
	RootCmd.AddCommand(importCmd)

	importCmd.Flags().StringP("user", "u", "", "username for registry")
	importCmd.Flags().StringP("password", "p", "", "password for registry")
	importCmd.Flags().StringP("target", "t", "", "target folder to download image")

}

func importImage(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Printf("Error: image url required!\n\n")
		cmd.Help()
		return
	}
	var url string
	urlString := args[0]
	source, image, err := getRegistrySource(urlString)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch source {
	case "docker":
		url = DockerHub
	default:
		fmt.Println("Invaild or unsupported source registry : ", source)
	}

	username, _ := cmd.Flags().GetString("user")
	password, _ := cmd.Flags().GetString("password")
	target, _ := cmd.Flags().GetString("target")

	hub, err := registry.New(url, username, password)
	if err != nil {
		fmt.Println("Error while creating registry connection")
		return
	}

	imageDetails := strings.Split(image, ":")
	if len(imageDetails) < 2 {
		imageDetails = append(imageDetails, "latest")
	}

	manifest, err := hub.Manifest(imageDetails[0], imageDetails[1])
	if err != nil {
		if strings.Contains(err.Error(), "status=404") {
			fmt.Printf("Error: image %s not found\n", imageDetails[0])
		}
		if strings.Contains(err.Error(), "status=401") {
			fmt.Printf("Authentication required for registry\n")
		}
		return
	}

	if target == "" {
		target = imageDetails[0] + "." + imageDetails[1]
	}
	os.MkdirAll(target, 0777)

	// Write Manifiest
	bManifest, _ := manifest.MarshalJSON()
	manifestFile, err := os.Create(filepath.Join(target, "manifest.json"))
	if err != nil {
		fmt.Println("Unable to create Manifest file")
		return
	}
	manifestFile.Write(bManifest)
	layerList := list.New()
	// Download layes within that folder.
	for _, layer := range manifest.Manifest.FSLayers {
		// or obtain the digest from an existing manifest's FSLayer list
		digest, _ := digest.ParseDigest(layer.BlobSum.String())
		reader, err := hub.DownloadLayer(imageDetails[0], digest)

		if reader != nil {
			defer reader.Close()
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		tarFile, err := os.Create(filepath.Join(target, layer.BlobSum.String()+".tar"))
		if err != nil {
			fmt.Println("Unable to create tar file")
			return
		}

		_, err = io.Copy(tarFile, reader)
		if err != nil {
			fmt.Println("Unable to create tar file")
			return
		}
		layerList.PushFront(layer.BlobSum.String())

	}
	tgt := filepath.Join(target, "rootfs")
	err = os.MkdirAll(tgt, 0777)
	if err != nil {
		fmt.Printf("Unable to creare folder :%s \n", tgt)
		fmt.Println(err)
		os.Exit(1)
	}
	// Extract layers to build rootfs
	for l := layerList.Front(); l != nil; l = l.Next() {
		source := filepath.Join(target, l.Value.(string))
		commandArgs := []string{"-xf", source + ".tar", "-C", tgt}
		err = exec.Command("tar", commandArgs...).Run()
		if err != nil {
			fmt.Printf("Error while extracting source: %s, Err: %v\n", source, err)
			//os.Exit(1)
		}
	}
	fmt.Printf("rootfs is prepared at : %s \n", tgt)
	for l := layerList.Front(); l != nil; l = l.Next() {
		source := filepath.Join(target, l.Value.(string)) + ".tar"
		os.Remove(source)
	}

}

func getRegistrySource(url string) (string, string, error) {
	source := strings.Split(url, "://")
	if len(source) < 2 {
		return "", "", fmt.Errorf("Error: invalid url : %s", url)
	}
	return source[0], source[1], nil
}
