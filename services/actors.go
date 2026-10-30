package services

import (
	"context"
	"errors"
	"io"
	"movies/dto"
	"movies/repositories"
	"movies/util"
)

type ActorService struct {
	AR *repositories.ActorRepo
	CS *CloudinaryService
}

func (AS *ActorService) UpdateActor(ctx context.Context, actorId int, ActorUpdateData dto.UpdateActorReq, image io.Reader) error {
	fields := util.UpdateActorReqToMap(ActorUpdateData)
	if image == nil {
		return AS.AR.UpdateActor(ctx, actorId, fields)
	}
	if AS.CS == nil {
		return errors.New("Cloudinary service is not initialized")
	}

	existingActor, err := AS.AR.GetActorById(ctx, actorId)
	if err != nil {
		return err
	}
	oldPublicID := existingActor.Actor.ImagePublicID

	newImage, err := AS.CS.UploadImage(ctx, image, "actors")
	if err != nil {
		return err
	}

	fields["image_url"] = newImage.URL
	fields["image_public_id"] = newImage.PublicID
	if err := AS.AR.UpdateActor(ctx, actorId, fields); err != nil {
		if newImage.PublicID != "" {
			_ = AS.CS.DeleteImage(ctx, newImage.PublicID)
		}
		return err
	}

	if oldPublicID != "" && oldPublicID != newImage.PublicID {
		return AS.CS.DeleteImage(ctx, oldPublicID)
	}
	return nil
}
