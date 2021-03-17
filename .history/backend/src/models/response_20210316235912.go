package models
import (
	"backend/function"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"github.com/gorilla/mux"
	"backend/models"
)

type LoginResponse struct {
	Response  string `json:"response,omitempty"`
	Error error  `json:"error,omitempty"`
}