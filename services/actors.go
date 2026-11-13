package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"movies/apperrors"
	"movies/dto"
	"movies/models"
	"movies/repositories"
	"movies/util"
	"time"
)

type ActorService struct {
	AR *repositories.ActorRepo
	CS *CloudinaryService
}

func (AS *ActorService) AddActor(ctx context.Context, ActorData dto.AddActorReq, image io.Reader) error {
	if ActorData.DateOfBirth.After(time.Now()) {
		return apperrors.ErrInvalidBirthDate
	}
	uploadedImage, err := AS.CS.UploadImage(ctx, image, "Actors")

	if err != nil {
		return fmt.Errorf("upload actor photo: %w", err)
	}

	imageURL := uploadedImage.URL
	imagePublicID := uploadedImage.PublicID

	Actor := models.Actor{
		ImageURL:      imageURL,
		ImagePublicID: imagePublicID,
		Name:          ActorData.Name,
		BirthDate:     ActorData.DateOfBirth,
		HomeLand:      ActorData.HomeLand,
	}

	err = AS.AR.AddActor(ctx, Actor)

	if err != nil {
		DeleteErr := AS.CS.DeleteImage(ctx, imagePublicID)
		if DeleteErr != nil {
			return fmt.Errorf(
				"create movie: %w; cleanup image: %v", err, DeleteErr)
		}
	}

	return nil
}

func (AS *ActorService) GetActorById(ctx context.Context, actorId int) (dto.ActorDetailsResponse, error) {
	ActorData, err := AS.AR.GetActorById(ctx, actorId)
	if err != nil {
		return dto.ActorDetailsResponse{}, err
	}

	Movies := make([]dto.MovieInActorResponse, 0, len(ActorData.Movies))

	for _, movie := range ActorData.Movies {
		var Movie dto.MovieInActorResponse
		Movie.MovieId = movie.MovieId
		Movie.CharacterName = movie.CharacterName
		Movie.Role = movie.Role
		Movie.Title = movie.Title
		Movies = append(Movies, Movie)
	}

	resData := dto.ActorDetailsResponse{
		Id:        actorId,
		Name:      ActorData.Actor.Name,
		BirthDate: ActorData.Actor.BirthDate,
		HomeLand:  ActorData.Actor.HomeLand,
		ImageURL:  ActorData.Actor.ImageURL,
		Movies:    Movies,
	}

	return resData, nil

}

func (AS *ActorService) GetActorByName(ctx context.Context, actorName string) ([]dto.ActorSummaryResponse, error) {
	actors, err := AS.AR.GetActorByName(ctx, actorName)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ActorSummaryResponse, 0, len(actors))
	for _, actor := range actors {
		responses = append(responses, dto.ActorSummaryResponse{
			Id:       actor.Id,
			Name:     actor.Name,
			ImageURL: actor.ImageURL,
		})
	}

	return responses, nil
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

func (AS *ActorService) DeleteActor(ctx context.Context, actorID int) error {
	return AS.AR.DeleteActor(ctx, actorID)
}
