package handlers

import (
	"io"
	"movies/apperrors"
	"movies/dto"
	"movies/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ActorHandler struct {
	AS *services.ActorService
}

func (AH *ActorHandler) AddActor(c *gin.Context) {
	ctx := c.Request.Context()

	var actorData dto.AddActorReq

	if err := c.ShouldBind(&actorData); err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	image, err := actorData.Image.Open()
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	defer image.Close()

	err = AH.AS.AddActor(ctx, actorData, image)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusCreated)
}

func (AH *ActorHandler) GetActorBySearch(c *gin.Context) {

	ctx := c.Request.Context()

	name := c.Query("name")

	actors, err := AH.AS.GetActorByName(ctx, name)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, actors)
}

func (AH *ActorHandler) GetActorById(c *gin.Context) {
	ctx := c.Request.Context()
	ActorId := c.Param("id")

	actorId, err := strconv.Atoi(ActorId)

	if err != nil || actorId <= 0 {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	actorData, err := AH.AS.GetActorById(ctx, actorId)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, actorData)
}

func (AH *ActorHandler) UpdateActor(c *gin.Context) {
	ctx := c.Request.Context()

	ActorId := c.Param("id")

	actorId, err := strconv.Atoi(ActorId)

	if err != nil || actorId <= 0 {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	var updateActorData dto.UpdateActorReq

	if err := c.ShouldBind(&updateActorData); err != nil {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	var image io.Reader
	if updateActorData.Image != nil {
		file, err := updateActorData.Image.Open()
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		defer file.Close()
		image = file
	}

	err = AH.AS.UpdateActor(ctx, actorId, updateActorData, image)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}

func (AH *ActorHandler) DeleteActor(c *gin.Context) {
	ctx := c.Request.Context()

	ActorId := c.Param("id")

	actorId, err := strconv.Atoi(ActorId)

	if err != nil || actorId <= 0 {
		c.Error(apperrors.ErrBadRequestData)
		c.Abort()
		return
	}

	if err = AH.AS.DeleteActor(ctx, actorId); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}
