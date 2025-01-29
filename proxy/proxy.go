package proxy

import (
	"os"

	"github.com/labstack/echo/v4"
)

func GetExistingObject(c echo.Context) error {
	localfs_path := os.Getenv("LOCALFS_PATH")

	bucket := c.Param("bucket")
	object := c.Param("object")

	return c.File(localfs_path + "/" + bucket + "/" + object)
}
