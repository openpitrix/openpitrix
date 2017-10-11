package repos

import (
	//	"net/http"
	//	log "github.com/Sirupsen/logrus"
	middleware "github.com/go-openapi/runtime/middleware"

	reposapi "openpitrix.io/openpitrix/pkg/swagger/restapi/operations/repos"
)

// GetRepos returns all the repositories
func GetRepos(params reposapi.GetReposParams) middleware.Responder {
	/*
		reposCollection, err := data.GetRepos()
		if err != nil {
			log.Error("unable to get Repos collection: ", err)
			return reposapi.NewGetAllReposDefault(http.StatusInternalServerError).WithPayload(internalServerErrorPayload())
		}
		var repos []*data.Repo
		reposCollection.FindAll(&repos)
		resources := helpers.MakeRepoResources(repos)

		payload := handlers.DataResourcesBody(resources)
		return reposapi.NewGetAllReposOK().WithPayload(payload)
	*/
	return middleware.NotImplemented("operation handlers.repos.GetRepos has not yet been implemented")
}
