package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"grd0.net/proxy/s3/database"
	"grd0.net/proxy/s3/proxy"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	s3_origin_endpoint := os.Getenv("S3_ORIGIN_ENDPOINT")
	s3_origin_use_ssl := os.Getenv("S3_ORIGIN_USE_SSL")
	s3_origin_protocol := "https"

	if s3_origin_use_ssl == "false" {
		s3_origin_protocol = "http"
	}

	// Create a file to store data and downloaded content
	localfs_path := os.Getenv("LOCALFS_PATH")
	err = os.MkdirAll(localfs_path, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("failed to create directories: %w", err))
	}

	database.InitDatabase()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8443", "https://cloud.grd0.net"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.Logger())

	e.GET("/:bucket/:object", proxy.GetExistingObject)
	// Must register the path params, otherwise middleware will fail to obtain the data
	e.HEAD("/:bucket/:object", func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "")
	})

	url1, err := url.Parse(s3_origin_protocol + "://" + s3_origin_endpoint)
	if err != nil {
		e.Logger.Fatal(err)
	}

	balancer := middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
		{
			URL: url1,
		},
	})

	e.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer: balancer,
		Skipper:  proxy.SkipperLogic,
	}))

	e.Use(middleware.Decompress())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.Logger.Fatal(e.Start(":80"))
}
