package api 
import("net/http"
"time"
"social-network/pkg/utils"
datamodels "social-network/pkg/dataModels"
"social-network/pkg/db/queries"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	response := datamodels.Response{}
	if r.Method == http.MethodPost {
		
	cookie, err1 := r.Cookie("session_id")
	if err1 != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusBadRequest, Status: "Failed", ErrorMsg: "no active session"})
		return
	}

	err2:= queries.DeleteSession(cookie.Value)
	if err2 != nil {
		utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "failed to logout"})
		return
	}

	// Expire the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(-time.Hour),
	})
	response.Code = 200
	response.Status = "OK"
	utils.SendResponse(w, response) //send the response
	}else{
		utils.SendResponse(w, datamodels.Response{Code: http.StatusMethodNotAllowed, Status: "Failed", ErrorMsg: "Method Not Allowed"})
		return
	}

}
