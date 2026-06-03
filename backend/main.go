package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/zero_trust"

	"github.com/creack/pty"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type ContainerInfo struct {
	ID         string   `json:"id"`
	Names      []string `json:"names"`
	Image      string   `json:"image"`
	State      string   `json:"state"`
	Status     string   `json:"status"`
	WorkingDir string   `json:"working_dir"`
}

type TunnelResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type RouteResponse struct {
	Hostname string `json:"hostname"`
	Service  string `json:"service"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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

	// Inisialisasi Cloudflare Client V4
	cfToken := os.Getenv("CF_API_TOKEN")
	cfAccountID := os.Getenv("CF_ACCOUNT_ID")
	cfClient := cloudflare.NewClient(option.WithAPIToken(cfToken))

	// ==========================================
	// DOCKER SUBSYSTEM ROUTING
	// ==========================================

	// 1. GET: List Container
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
			workingDir := item.Labels["com.docker.compose.project.working_dir"]

			var cleanNames []string
			for _, name := range item.Names {
				cleanNames = append(cleanNames, strings.TrimPrefix(name, "/"))
			}

			cleanContainers = append(cleanContainers, ContainerInfo{
				ID:         item.ID[:8],
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

		err = cli.ContainerStop(context.Background(), containerID, container.StopOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mematikan container: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Container " + containerID + " berhasil dimatikan"})
	})

	// 4. POST: Down Container
	r.POST("/api/docker/containers/:id/down", func(c *gin.Context) {
		containerID := c.Param("id")

		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cli.Close()

		ctx := context.Background()

		inspect, err := cli.ContainerInspect(ctx, containerID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal inspeksi container: " + err.Error()})
			return
		}

		if inspect.State.Running {
			log.Printf("Menghentikan container %s sebelum dihapus...", containerID)
			err = cli.ContainerStop(ctx, containerID, container.StopOptions{})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghentikan container: " + err.Error()})
				return
			}
		}

		log.Printf("Menghapus container %s...", containerID)
		err = cli.ContainerRemove(ctx, containerID, container.RemoveOptions{
			RemoveVolumes: false,
			Force:         true,
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

	// 5. GET/WS: Real-time Terminal Stream
	r.GET("/api/docker/containers/:id/terminal", func(c *gin.Context) {
		containerID := c.Param("id")

		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Gagal melakukan upgrade koneksi WebSocket:", err)
			return
		}
		defer ws.Close()

		cmd := exec.Command("bash", "--login")
		cmd.Env = os.Environ()

		if containerID == "local-machine" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				ws.WriteMessage(websocket.TextMessage, []byte("\r\nGagal mendeteksi Home Directory perangkat.\r\n"))
				return
			}
			cmd.Dir = homeDir
		} else {
			cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				ws.WriteMessage(websocket.TextMessage, []byte("\r\nGagal menghubungkan ke Docker Engine: " + err.Error() + "\r\n"))
				return
			}
			defer cli.Close()

			inspect, err := cli.ContainerInspect(context.Background(), containerID)
			if err != nil {
				ws.WriteMessage(websocket.TextMessage, []byte("\r\nGagal melakukan inspeksi kontainer: " + err.Error() + "\r\n"))
				return
			}

			workingDir := inspect.Config.Labels["com.docker.compose.project.working_dir"]
			if workingDir == "" {
				ws.WriteMessage(websocket.TextMessage, []byte("\r\n[WARN]: Label working_dir proyek tidak ditemukan. Menggunakan fallback root.\r\n"))
				homeDir, _ := os.UserHomeDir()
				cmd.Dir = homeDir
			} else {
				if _, err := os.Stat(workingDir); os.IsNotExist(err) {
					ws.WriteMessage(websocket.TextMessage, []byte("\r\n[ERROR]: Direktori proyek " + workingDir + " tidak ditemukan di host ini.\r\n"))
					return
				}
				cmd.Dir = workingDir
			}
		}

		ptmx, err := pty.Start(cmd)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("\r\nGagal memulai sesi PTY: " + err.Error() + "\r\n"))
			return
		}
		defer ptmx.Close()

		done := make(chan struct{})

		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := ptmx.Read(buf)
				if n > 0 {
					_ = ws.WriteMessage(websocket.TextMessage, buf[:n])
				}
				if err != nil {
					break
				}
			}
			close(done)
		}()

		go func() {
			for {
				_, message, err := ws.ReadMessage()
				if err != nil {
					break
				}
				_, _ = ptmx.Write(message)
			}
		}()

		<-done
		_ = cmd.Wait()
	})

	// ==========================================
	// CLOUDFLARE TUNNEL SUBSYSTEM ROUTING
	// ==========================================

	// 1. GET: List Cloudflare Tunnels
	r.GET("/api/cloudflare/tunnels", func(c *gin.Context) {
		if cfToken == "" || cfAccountID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "API Token atau Account ID Cloudflare belum di-set di env"})
			return
		}

		tunnels, err := cfClient.ZeroTrust.Tunnels.Cloudflared.List(
			context.Background(),
			zero_trust.TunnelCloudflaredListParams{
				AccountID: cloudflare.F(cfAccountID),
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data Cloudflare: " + err.Error()})
			return
		}

		var list []TunnelResponse
		for _, t := range tunnels.Result {
			list = append(list, TunnelResponse{
				ID:        t.ID,
				Name:      t.Name,
				Status:    string(t.Status),
				CreatedAt: t.CreatedAt.String(),
			})
		}

		c.JSON(http.StatusOK, list)
	})

	// 2. GET: List Published Ingress Routes
	r.GET("/api/cloudflare/tunnels/:id/routes", func(c *gin.Context) {
		tunnelID := c.Param("id")
		if cfToken == "" || cfAccountID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "API Token atau Account ID Cloudflare belum di-set di env"})
			return
		}

		config, err := cfClient.ZeroTrust.Tunnels.Cloudflared.Configurations.Get(
			context.Background(),
			tunnelID,
			zero_trust.TunnelCloudflaredConfigurationGetParams{
				AccountID: cloudflare.F(cfAccountID),
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil konfigurasi rute: " + err.Error()})
			return
		}

		var routes []RouteResponse
		for _, ingress := range config.Config.Ingress {
			if ingress.Hostname != "" {
				routes = append(routes, RouteResponse{
					Hostname: ingress.Hostname,
					Service:  ingress.Service,
				})
			}
		}

		c.JSON(http.StatusOK, routes)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}