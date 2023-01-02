package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mr-keppy/bookings/internal/config"
	"github.com/mr-keppy/bookings/internal/models"
)


var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M){
	gob.Register(models.Reservation{})
	testApp.InProduction = false

	session = scs.New()

	session.Lifetime = 24*time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction

	testApp.Session = session
	app = &testApp
	
	os.Exit(m.Run())
}
func TestDefaultData(t *testing.T){
	
}