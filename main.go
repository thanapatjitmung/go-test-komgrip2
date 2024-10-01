package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Menu struct {
	ID         int64    `json:"id"`
	Name       string   `json:"name"`
	Route      string   `json:"route"`
	Icon       string   `json:"icon"`
	IsChildren bool     `json:"is_children"`
	Children   []*Child `json:"children"`
}

type Child struct {
	ID       int64  `json:"id"`
	ParentID int64  `json:"parent_id"`
	Name     string `json:"name"`
	Route    string `json:"route"`
}

var menuData = []*Menu{
	{
		ID:         1,
		Name:       "ข้อมูลสรุป",
		Route:      "/dashboards",
		Icon:       "fi-rr-chart-pie-alt",
		IsChildren: false,
	},
	{
		ID:         2,
		Name:       "ข้อมูลอื่นๆ",
		Route:      "/dashboards/others",
		Icon:       "fi-rr-chart-pie-alt",
		IsChildren: false,
	},
}

var childData = []*Child{
	{ID: 16, ParentID: 1, Name: "สินทรัพย์", Route: "/dashboards/asset"},
	{ID: 17, ParentID: 1, Name: "สภาพทาง", Route: "/dashboards/condition"},
	{ID: 18, ParentID: 2, Name: "ตัวเลือก 1", Route: "/dashboards/option1"},
	{ID: 19, ParentID: 2, Name: "ตัวเลือก 2", Route: "/dashboards/option2"},
}

func getMenu(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	var selectedMenu []*Menu

	if id != "" && idInt > 0 {
		for _, menu := range menuData {
			if fmt.Sprintf("%d", menu.ID) == id {
				selectedMenu = append(selectedMenu, menu)
				break
			}
		}
	} else {
		selectedMenu = append(selectedMenu, menuData...)
	}

	if selectedMenu == nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status": false,
			"code":   404,
			"data":   "Menu not found",
		})
	}

	var children []*Child
	for _, menu := range selectedMenu {
		for _, child := range childData {
			if child.ParentID == menu.ID {
				children = append(children, child)
			}
		}
		menu.Children = children
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": true,
		"code":   200,
		"data":   selectedMenu,
	})
}

func main() {
	e := echo.New()
	e.GET("/api/menus/:id", getMenu)
	e.Logger.Fatal(e.Start(":8080"))
}
