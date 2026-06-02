package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Struct ringkas untuk dikirim ke frontend Vue
type ContainerInfo struct {
	ID         string   `json:"id"`
	Names      []string `json:"names"`
	Image      string   `json:"image"`
	State      string   `json:"state"` // running, exited, dll
	Status     string   `json:"status"` // Exited (0) 6 hours ago
	WorkingDir string   `json:"working_dir"` // Lokasi project docker-compose
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan")
	}

	r := gin.Default()

	// Middleware CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 1. GET: List Container yang sudah diformat ringkas
	r.GET("/api/docker/containers", func(c *gin.Context) {
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cli.Close()

		containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var cleanContainers []ContainerInfo
		for _, item := range containers {
			// Ambil info working directory docker-compose dari label jika ada
			workingDir := item.Labels["com.docker.compose.project.working_dir"]

			// Rapikan nama container (buang garing "/" di depan nama)
			var cleanNames []string
			for _, name := range item.Names {
				cleanNames = append(cleanNames, strings.TrimPrefix(name, "/"))
			}

			cleanContainers = append(cleanContainers, ContainerInfo{
				ID:         item.ID[:12], // potong 12 karakter aja biar rapi di UI
				Names:      cleanNames,
				Image:      item.Image,
				State:      item.State,
				Status:     item.Status,
				WorkingDir: workingDir,
			})
		}

		c.JSON(http.StatusOK, cleanContainers)
	})

	// 2. POST: Start Container
	r.POST("/api/docker/containers/:id/start", func(c *gin.Context) {
		containerID := c.Param("id")

		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cli.Close()

		err = cli.ContainerStart(context.Background(), containerID, container.StartOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyalakan container: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Container " + containerID + " berhasil dinyalakan"})
	})

	// 3. POST: Stop Container
	r.POST("/api/docker/containers/:id/stop", func(c *gin.Context) {
		containerID := c.Param("id")

		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cli.Close()

		// Stop dengan timeout default (nil berarti pakai default bawaan docker)
		err = cli.ContainerStop(context.Background(), containerID, container.StopOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mematikan container: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Container " + containerID + " berhasil dimatikan"})
	})

	// 4. POST: Down Container (Hanya Stop dan Remove Container, Tanpa Hapus Volume)
	r.POST("/api/docker/containers/:id/down", func(c *gin.Context) {
		containerID := c.Param("id")

		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cli.Close()

		ctx := context.Background()

		// 1. Cek status container terlebih dahulu
		inspect, err := cli.ContainerInspect(ctx, containerID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal inspeksi container: " + err.Error()})
			return
		}

		// 2. Jika container masih running, kita STOP dulu
		if inspect.State.Running {
			log.Printf("Menghentikan container %s sebelum dihapus...", containerID)
			err = cli.ContainerStop(ctx, containerID, container.StopOptions{})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghentikan container: " + err.Error()})
				return
			}
		}

		// 3. REMOVE Container (Tanpa menghapus volume pendukungnya)
		log.Printf("Menghapus container %s...", containerID)
		err = cli.ContainerRemove(ctx, containerID, container.RemoveOptions{
			RemoveVolumes: false, // <-- Ini dikunci di false agar volume tetap aman
			Force:         true,  // Paksa hapus jika ada kendala minor
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus container: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "Container berhasil di-down (dihentikan & dihapus) tanpa menghapus volume data.",
			"container_id": containerID,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}