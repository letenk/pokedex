package usecase

import (
	"context"
	"log"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/letenk/pokedex/models/domain"
	"github.com/letenk/pokedex/models/web"
	"github.com/letenk/pokedex/repository"
	"github.com/letenk/pokedex/util"
)

type MonsterUsecase interface {
	Create(ctx context.Context, monster web.MonsterCreateRequest, file multipart.File, fileName string) (domain.Monster, error)
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
	}

	// Upload to aws S3
	imageLocationS3, err := UploadToAwsS3(ctx, file, fileName)
	if err != nil {
		return monster, err
	}
	// Passing image location in aws to object monster
	monster.Image = imageLocationS3

	// Create
	monster, err = u.repository.Create(ctx, monster)
	if err != nil {
		return monster, err
	}

	return monster, nil
}
