package controller

func SearchChats(w http.ResponseWriter, r *http.Request) {
  x, _ := ioutil.ReadAll(r.Body)
  state := model.UUID(vars["state"])
	room := vars["room"]

  _, isAuth := auth.CheckAuth(state)
	if !isAuth {
		errorHandler(w, r, 403, "Not Authorized!")
		return
	}
  searchString := string(x)
  
  //Call lambda python script to search chats

}
