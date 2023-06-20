package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// albums file path
var albumsPath string

func main() {
	albumsPath = *flag.String("path", "/tmp/albums.json", "path to store the albums")
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	albums = loadFromFile(albumsPath)

	router.Run("localhost:8080")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	writeToFile(albums, albumsPath)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func writeToFile(albs []album, path string) {
	file, _ := json.MarshalIndent(albs, "", " ")

	_ = ioutil.WriteFile(path, file, 0644)
}

func loadFromFile(path string) []album {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to load albums file. Using new file", path)
		return albums
	}

	var albs []album

	err = json.Unmarshal(body, &albs)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to process albums file. Using new file", path)
		return albums
	}

	fmt.Println("Using exiting albums file at", path)
	// fmt.Println(fmt.Sprintln("Using exiting albums file at", path))
	return albs
}
