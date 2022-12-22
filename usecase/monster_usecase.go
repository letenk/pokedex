package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/letenk/pokedex/models/domain"
	"github.com/letenk/pokedex/models/web"
	"github.com/letenk/pokedex/repository"
	"github.com/letenk/pokedex/util"
)

type MonsterUsecase interface {
	FindAll(ctx context.Context, reqQuery web.MonsterQueryRequest) ([]domain.Monster, error)
	FindByID(ctx context.Context, ID string) (domain.Monster, error)
	Create(ctx context.Context, monster web.MonsterCreateRequest, file multipart.File, fileName string) (domain.Monster, error)
	Update(ctx context.Context, ID string, reqUpdate web.MonsterUpdateRequest, file multipart.File, fileName string) (domain.Monster, error)
	UpdateMarkMonsterCaptured(ctx context.Context, ID string, reqUpdate web.MonsterUpdateRequestMonsterCapture) (bool, error)
	Delete(ctx context.Context, ID string) (bool, error)
}

type monsterUsecase struct {
	repository repository.MonsterRepository
}

func NewUsecaseMonster(repository repository.MonsterRepository) *monsterUsecase {
	return &monsterUsecase{repository}
}

func UploadToAwsS3(ctx context.Context, file multipart.File, fileName string) (string, error) {
	// Load Config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// Config
	s3Config := &aws.Config{
		Region:      aws.String(config.AWS_REGION), // set region aws
		Credentials: credentials.NewStaticCredentials(config.AWS_KEY_ID, config.AWS_SECRET_KEY, ""),
	}

	// Create new instance session
	sess := session.New(s3Config)

	// Bucket name
	bucketName := config.AWS_BUCKET_NAME
	// Create new uploadert with session
	uploader := s3manager.NewUploader(sess)

	// Create context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Create object input
	input := &s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		// ACL:         aws.String("public-read"),
		Body:        file,
		ContentType: aws.String("image/jpg"),
	}

	// Upload to aws with context
	res, err := uploader.UploadWithContext(ctx, input)
	if err != nil {
		return "", err
	}

	return res.Location, nil
}

func DeleteItemFromAwsS3(ctx context.Context, keyName string) error {
	// Load Config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// Config
	s3Config := &aws.Config{
		Region:      aws.String(config.AWS_REGION), // set region aws
		Credentials: credentials.NewStaticCredentials(config.AWS_KEY_ID, config.AWS_SECRET_KEY, ""),
	}

	// Create new instance session
	sess := session.New(s3Config)

	// Creates a new instance of the S3 client with a session.
	svc := s3.New(sess)

	// Delete Object
	_, err = svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(config.AWS_BUCKET_NAME),
		Key:    aws.String(keyName),
	})

	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(config.AWS_BUCKET_NAME),
		Key:    aws.String(keyName),
	})

	if err != nil {
		return err
	}

	return nil
}

func (u *monsterUsecase) Create(ctx context.Context, req web.MonsterCreateRequest, file multipart.File, fileName string) (domain.Monster, error) {
	// Passing data request into object monster
	monster := domain.Monster{
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		Length:      req.Length,
		Weight:      req.Weight,
		Hp:          req.Hp,
		Attack:      req.Attack,
		Defends:     req.Defends,
		Speed:       req.Speed,
		TypeID:      req.TypeID,
		ImageName:   fileName,
		ImageURL:    fileName,
	}

	// Create
	monster, err := u.repository.Create(ctx, monster)
	if err != nil {
		return monster, err
	}

	// Upload to aws S3
	imageLocationS3, err := UploadToAwsS3(ctx, file, fileName)
	if err != nil {
		return monster, err
	}

	// Passing image location in aws to object monster
	monster.ImageURL = imageLocationS3

	// Update image
	_, err = u.repository.Update(ctx, monster)
	if err != nil {
		return monster, err
	}

	return monster, nil
}

func (u *monsterUsecase) FindAll(ctx context.Context, reqQuery web.MonsterQueryRequest) ([]domain.Monster, error) {
	// Find all
	monsters, err := u.repository.FindAll(ctx, reqQuery)
	if err != nil {
		return monsters, err
	}

	return monsters, nil
}

func (u *monsterUsecase) FindByID(ctx context.Context, ID string) (domain.Monster, error) {
	// Find by id
	monster, err := u.repository.FindByID(ctx, ID)
	if err != nil {
		return monster, err
	}

	return monster, nil
}

func (u *monsterUsecase) Update(ctx context.Context, ID string, reqUpdate web.MonsterUpdateRequest, file multipart.File, fileName string) (domain.Monster, error) {

	// Find by id
	currentMonster, err := u.repository.FindByID(ctx, ID)
	if err != nil {
		return currentMonster, err
	}
	// Parse reqUpdate form when not empty
	if reqUpdate.Name != "" {
		currentMonster.Name = reqUpdate.Name
	}
	if reqUpdate.CategoryID != "" {
		currentMonster.CategoryID = reqUpdate.CategoryID
	}
	if reqUpdate.Description != "" {
		currentMonster.Description = reqUpdate.Description
	}
	if reqUpdate.Length != "" {
		floatLength, err := strconv.ParseFloat(reqUpdate.Length, 32)
		if err != nil {
			return currentMonster, err
		}
		currentMonster.Length = float32(floatLength)
	}
	if reqUpdate.Weight != "" {
		intWeight, err := strconv.Atoi(reqUpdate.Weight)
		if err != nil {
			return currentMonster, err
		}
		currentMonster.Weight = uint16(intWeight)
	}
	if reqUpdate.Hp != "" {
		intHp, err := strconv.Atoi(reqUpdate.Hp)
		if err != nil {
			return currentMonster, err
		}
		currentMonster.Hp = uint16(intHp)
	}
	if reqUpdate.Attack != "" {
		intAttack, err := strconv.Atoi(reqUpdate.Attack)
		if err != nil {
			return currentMonster, err
		}
		currentMonster.Attack = uint16(intAttack)
	}
	if reqUpdate.Defends != "" {
		intDefends, err := strconv.Atoi(reqUpdate.Defends)
		if err != nil {
			return currentMonster, err
		}
		currentMonster.Defends = uint16(intDefends)
	}
	if reqUpdate.Speed != "" {
		intSpeed, err := strconv.Atoi(reqUpdate.Speed)
		if err != nil {
			return currentMonster, err
		}
		currentMonster.Speed = uint16(intSpeed)
	}
	if reqUpdate.Catched != "" {
		catchedBool, err := strconv.ParseBool(reqUpdate.Catched)
		if err != nil {
			return currentMonster, err
		}
		currentMonster.Catched = catchedBool
	}

	if len(reqUpdate.TypeID) != 0 {
		currentMonster.TypeID = reqUpdate.TypeID
	}

	dataUpdate := domain.Monster{
		ID:          currentMonster.ID,
		Name:        currentMonster.Name,
		CategoryID:  currentMonster.CategoryID,
		Description: currentMonster.Description,
		Length:      currentMonster.Length,
		Weight:      currentMonster.Weight,
		Hp:          currentMonster.Hp,
		Attack:      currentMonster.Attack,
		Defends:     currentMonster.Defends,
		Speed:       currentMonster.Speed,
		Catched:     currentMonster.Catched,
		ImageName:   currentMonster.ImageName,
		ImageURL:    currentMonster.ImageURL,
		TypeID:      reqUpdate.TypeID,
	}

	// Update monster image
	monsterUpdated, err := u.repository.Update(ctx, dataUpdate)
	if err != nil {
		return currentMonster, err
	}

	if fileName != "" {
		ctxToAws, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// Remove item from aws
		err := DeleteItemFromAwsS3(ctxToAws, monsterUpdated.ImageName)

		if err != nil {
			return currentMonster, err
		}

		// Upload new image to aws S3
		newImageLocationS3, err := UploadToAwsS3(ctxToAws, file, fileName)
		if err != nil {
			return currentMonster, err
		}

		// Update current imageName to new image and new url image
		monsterUpdated.ImageName = fileName
		monsterUpdated.ImageURL = newImageLocationS3

		// Update monster image
		_, err = u.repository.Update(ctx, monsterUpdated)
		if err != nil {
			return currentMonster, err
		}
	}

	return monsterUpdated, nil
}

func (u *monsterUsecase) UpdateMarkMonsterCaptured(ctx context.Context, ID string, reqUpdate web.MonsterUpdateRequestMonsterCapture) (bool, error) {
	// Find by id
	currentMonster, err := u.repository.FindByID(ctx, ID)
	if err != nil {
		return false, err
	}

	currentMonster.Catched = reqUpdate.Catched

	// Update monster
	_, err = u.repository.Update(ctx, currentMonster)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (u *monsterUsecase) Delete(ctx context.Context, ID string) (bool, error) {
	// Find monster
	monster, err := u.repository.FindByID(ctx, ID)

	if monster.ID == "" {
		var msgNotfound error
		msg := fmt.Sprintf("monster with id %s not found", ID)
		msgNotfound = errors.New(msg)

		return false, msgNotfound
	}

	if err != nil {
		return false, err
	}

	// Delete
	ok, err := u.repository.Delete(ctx, monster)
	if err != nil {
		return false, err
	}

	// Remove image in aws
	if ok {
		ctxToAws, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// Remove item from aws
		err := DeleteItemFromAwsS3(ctxToAws, monster.ImageName)

		if err != nil {
			return false, err
		}
	}

	return true, nil
}
