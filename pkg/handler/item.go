package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirktriplefive/test"
)

func (h *Handler) createItem(c *gin.Context){
	var input test.Item
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return 
	}

	
	rid, err := h.services.Item.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"rid": rid,
	})

}

type getAllItemResponse struct {
	Data []test.Item `json:"data"`
}

func (h *Handler) getAllItem(c *gin.Context){

	lists, err := h.services.Item.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 
	}

	c.JSON(http.StatusOK, getAllItemResponse{
		Data: lists, 
	})
}