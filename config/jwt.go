package config

import (
	"os"
	"time"
)

var JWTSecret = []byte(os.Getenv("JWT_SECRET")) // Replace with a strong secret key
var TokenExpiration = time.Hour * 24            // Token valid for 24 hours
