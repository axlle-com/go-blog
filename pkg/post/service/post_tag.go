package service

import (
	"errors"
	contracts2 "github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/axlle-com/blog/pkg/post/repository"
	"gorm.io/gorm"
)

type PostTagService struct {
	postTagRepo  repository.PostTagRepository
	resourceRepo repository.PostTagResourceRepository
}

func NewPostTagService(
	postTagRepo repository.PostTagRepository,
	resourceRepo repository.PostTagResourceRepository,
) *PostTagService {
	return &PostTagService{
		postTagRepo:  postTagRepo,
		resourceRepo: resourceRepo,
	}
}

func (s *PostTagService) CreatePostTag(postTag *models.PostTag) (*models.PostTag, error) {
	if err := s.postTagRepo.Create(postTag); err != nil {
		return nil, err
	}

	return postTag, nil
}

func (s *PostTagService) UpdatePostTag(postTag *models.PostTag) (*models.PostTag, error) {
	if err := s.postTagRepo.Update(postTag); err != nil {
		return nil, err
	}

	return postTag, nil
}

func (s *PostTagService) Attach(resource contracts2.Resource, postTag contracts2.PostTag) error {
	hasRepo, err := s.resourceRepo.GetByParams(resource.GetUUID(), postTag.GetID())
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if hasRepo == nil {
		err = s.resourceRepo.Create(
			&models.PostTagHasResource{
				ResourceUUID: resource.GetUUID(),
				PostTagID:    postTag.GetID(),
			},
		)
	}
	return nil
}

func (s *PostTagService) DeleteForResource(resource contracts2.Resource) (err error) {
	byResource, err := s.resourceRepo.GetByResource(resource)
	if err != nil {
		return err
	}

	all := make(map[uint]*models.PostTagHasResource)
	only := make(map[uint]*models.PostTagHasResource)
	detach := make(map[uint]*models.PostTagHasResource)
	var postTagIDs []uint
	if byResource == nil {
		return nil
	}

	for _, r := range byResource {
		if r.ResourceUUID != resource.GetUUID() {
			all[r.PostTagID] = r
		} else {
			only[r.PostTagID] = r
		}
	}

	for id, _ := range only {
		if _, ok := all[id]; ok {
			detach[id] = all[id]
		} else {
			postTagIDs = append(postTagIDs, id)
		}
	}

	if len(detach) > 0 { // TODO need test
		for _, r := range detach {
			err = s.resourceRepo.DeleteByParams(r.ResourceUUID, r.PostTagID)
			if err != nil {
				return err
			}
		}
	}

	if len(postTagIDs) > 0 {
		postTags, err := s.postTagRepo.GetByIDs(postTagIDs)
		if err != nil {
			return err
		}
		err = s.DeletePostTags(postTags)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *PostTagService) DeletePostTags(postTags []*models.PostTag) (err error) {
	var ids []uint
	for _, postTag := range postTags {
		ids = append(ids, postTag.ID)
	}

	if len(ids) > 0 {
		if err = s.postTagRepo.DeleteByIDs(ids); err == nil {
			return nil
		}
	}
	return err
}
