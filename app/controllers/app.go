package controllers

import (
	"bitbucket.org/evard/evardbugs/app/models"
	"bitbucket.org/evard/evardbugs/app/routes"
	"github.com/revel/revel"
	"github.com/ottob/go-semver/semver"
)

type App struct {
	GorpController
	PageTitle string
}

func (c *App) PageLoad() revel.Result {
	c.PageTitle = "Correct page title"
	return nil
}

func (c App) Index() revel.Result {
	c.RenderArgs["pageTitle"] = c.PageTitle
	return c.Render()
}

func (c App) IndexPost(message string) revel.Result {
	s := revel.Config.StringDefault("app.version", "0.1")
	ver, err := semver.NewVersion(s)
	if err != nil {
		panic(err)
	}

	newCase := models.Case{Message: message, GuideVersion: ver}

	if err := c.Txn.Insert(&newCase); err != nil {
		panic(err)
	}

	c.Flash.Success("Success")
	return c.Redirect(routes.App.Index())
}
