package api
import(
	"net/http"
	"encoding/json"
	"social-network/pkg/db/queries"
	"log"
	"social-network/pkg/utils"
	datamodels "social-network/pkg/dataModels"
	"fmt"
	 //"strconv"
	 //"strings"
)
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
	users, err:= queries.GetUsersList()
	if err != nil {
		log.Println(err)
			utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error22"})
			return
	}
	cookie, err := r.Cookie("session_id")
    if err != nil {
        fmt.Println("Error retrieving session_id cookie:", err)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusUnauthorized, Status: "Failed", ErrorMsg: "session ID not found"})
        return
    }
	CreatorID, err3 := queries.ValidateSession(cookie.Value)
	if err3 != nil {
		fmt.Println("Error retriving user id from session", err3)
        utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
        return
	}
	//remove the creator from the invitation list
	var users2 []datamodels.User
	for i:= 0; i<len(users);i++ {
		if users[i].ID != CreatorID {
			users2= append(users2,users[i])
		}
	}
	
	// users3, err4:= queries.GetAvailableUsersList(groupIDInt)
	// if err4 != nil {
	// 	fmt.Println("Error retriving the available userlist", err4)
    //     utils.SendResponse(w, datamodels.Response{Code: http.StatusInternalServerError, Status: "Failed", ErrorMsg: "Internal Server Error"})
    //     return
	// }
	// for i:= 0 ; i< len(users3);i++ {
	// 	fmt.Println(users3[i])
	// }
    // json.NewEncoder(w).Encode(users)
	w.Header().Set("Content-Type", "application/json")
if err := json.NewEncoder(w).Encode(users2); err != nil {
	log.Println("Error encoding users:", err)
	http.Error(w, "Failed to encode response", http.StatusInternalServerError)
}

}
}
