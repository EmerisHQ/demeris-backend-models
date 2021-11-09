package validation

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin/binding"

	"github.com/go-playground/validator/v10"
)

const maxPort = 65535
const validScheme = "https"

func CosmosRPCURL(structValidator binding.StructValidator) {
	if v, ok := structValidator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("cosmosrpcurl", func(fl validator.FieldLevel) bool {
			path, ok := fl.Field().Interface().(string)
			if !ok {
				return false
			}

			if err := validateCosmosRpcUrl(path); err != nil {
				return false
			}

			return true
		}); err != nil {
			panic(err)
		}
	}
}

func validateCosmosRpcUrl(value string) error {

	if value == "" {
		return nil
	}

	url, err := url.Parse(value)
	if err != nil {
		return err
	}

	if url.Scheme != validScheme {
		return fmt.Errorf("unsupported URL scheme %s", url.Scheme)
	}

	if url.User != nil {
		return errors.New("URL cannot contain user")
	}

	_, port, _ := net.SplitHostPort(url.Host)
	x, err := strconv.Atoi(port)
	if err != nil || x > maxPort {
		return fmt.Errorf("invalid port %s", port)
	}

	if url.Path != "" || url.Fragment != "" || url.RawQuery != "" {
		return errors.New("URL cannot contain path info")
	}

	return nil
}
