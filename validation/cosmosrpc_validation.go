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

const (
	maxPort     = 65535
	schemeHttps = "https"
	schemeHttp  = "http"
)

func CosmosRPCURL(structValidator binding.StructValidator) {
	if v, ok := structValidator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("cosmosrpcurl", registerCosmosRpcValidation); err != nil {
			panic(err)
		}
	}
}

func registerCosmosRpcValidation(fl validator.FieldLevel) bool {
	path, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	if err := validateCosmosRpcUrl(path); err != nil {
		return false
	}

	return true
}

func validateCosmosRpcUrl(value string) error {

	if value == "" {
		return nil
	}

	parsedUrl, err := url.Parse(value)
	if err != nil {
		return err
	}

	// FIXME: Allow HTTP only on local/DEV
	if parsedUrl.Scheme != schemeHttp && parsedUrl.Scheme != schemeHttps {
		return fmt.Errorf("unsupported URL scheme %s", parsedUrl.Scheme)
	}

	// we allow basic auth for the edge case of private RPCs
	_, pwdSet := parsedUrl.User.Password()
	if parsedUrl.User.Username() != "" && !pwdSet {
		return errors.New("URL cannot contain user without password")
	}

	_, port, err := net.SplitHostPort(parsedUrl.Host)
	if err != nil {
		return fmt.Errorf("invalid host:port %s", parsedUrl.Host)
	}

	x, err := strconv.Atoi(port)
	if err != nil || x > maxPort {
		return fmt.Errorf("invalid port %s", port)
	}

	// we allow Paths for the edge case of the RPC being behind an API gateway
	if parsedUrl.Fragment != "" || parsedUrl.RawQuery != "" {
		return errors.New("URL cannot contain fragments or query params")
	}

	return nil
}
