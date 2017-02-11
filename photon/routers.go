package photon

// Contains functionality for routers API.
type RoutersAPI struct {
	client *Client
}

var routerUrl string = "/routers/"

// Gets a router with the specified ID.
func (api *RoutersAPI) Get(id string) (router *Router, err error) {
	res, err := api.client.restClient.Get(api.client.Endpoint+routerUrl+id, api.client.options.TokenOptions)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	var result Router
	err = json.NewDecoder(res.Body).Decode(&result)
	return &result, nil
}

// Sets router's name.
func (api *RoutersAPI) SetName(id string, routerName *RouterSetNameOperation) (task *Task, err error) {
	body, err := json.Marshal(routerName)
	if err != nil {
		return
	}

	res, err := api.client.restClient.Post(
		api.client.Endpoint+routerUrl+id+"/set_router_name",
		"application/json",
		bytes.NewReader(body),
		api.client.options.TokenOptions)

	if err != nil {
		return
	}

	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}